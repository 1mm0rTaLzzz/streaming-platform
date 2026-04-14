const CORS_EXPOSE_HEADERS = [
  'Accept-Ranges',
  'Cache-Control',
  'Content-Length',
  'Content-Range',
  'Content-Type',
  'ETag',
  'Last-Modified',
];

const FORWARDED_REQUEST_HEADERS = [
  'accept',
  'accept-encoding',
  'accept-language',
  'cache-control',
  'if-match',
  'if-modified-since',
  'if-none-match',
  'if-range',
  'if-unmodified-since',
  'range',
  'referer',
  'user-agent',
];

export default {
  async fetch(request, env) {
    if (request.method === 'OPTIONS') {
      return handlePreflight(request);
    }

    if (request.method !== 'GET' && request.method !== 'HEAD') {
      return jsonError(405, 'Method not allowed. Use GET, HEAD, or OPTIONS.');
    }

    const route = extractUpstreamFromRequest(request);
    if (!route.upstreamUrl) {
      return jsonError(
        400,
        'Missing upstream URL. Use /proxy?url=https%3A%2F%2Fexample.com%2Findex.m3u8 or /proxy/<encoded-url>.'
      );
    }

    let upstreamUrl;
    try {
      upstreamUrl = new URL(route.upstreamUrl);
    } catch {
      return jsonError(400, 'Invalid upstream URL.');
    }

    if (!isHttpProtocol(upstreamUrl)) {
      return jsonError(400, 'Only http:// and https:// upstream URLs are allowed.');
    }

    const allowedHosts = parseAllowedHosts(env.ALLOWED_UPSTREAM_HOSTS);
    if (allowedHosts.length === 0) {
      return jsonError(500, 'Worker is not configured: ALLOWED_UPSTREAM_HOSTS is empty.');
    }

    const normalizedHost = upstreamUrl.hostname.toLowerCase();
    if (isBlockedLocalHost(normalizedHost) || isBlockedAddressLiteral(normalizedHost)) {
      return jsonError(403, 'Upstream host is blocked.');
    }

    if (!isHostAllowed(normalizedHost, allowedHosts)) {
      return jsonError(403, 'Upstream host is not in ALLOWED_UPSTREAM_HOSTS.');
    }

    const upstreamRequest = buildUpstreamRequest(request, upstreamUrl);

    let upstreamResponse;
    try {
      upstreamResponse = await fetch(upstreamRequest, {
        redirect: 'follow',
      });
    } catch {
      return jsonError(502, 'Failed to reach upstream.');
    }

    return withCors(upstreamResponse);
  },
};

function extractUpstreamFromRequest(request) {
  const requestUrl = new URL(request.url);

  const queryUrl = requestUrl.searchParams.get('url');
  if (queryUrl) {
    return { upstreamUrl: queryUrl };
  }

  const proxyPrefix = '/proxy/';
  if (!requestUrl.pathname.startsWith(proxyPrefix)) {
    return { upstreamUrl: null };
  }

  const encodedPart = requestUrl.pathname.slice(proxyPrefix.length);
  if (!encodedPart) {
    return { upstreamUrl: null };
  }

  const decoded = safeDecodeURIComponent(encodedPart);
  if (!decoded) {
    return { upstreamUrl: null };
  }

  // Optional helper: if caller puts query params on the proxy URL while
  // using path format, append them to upstream.
  if (!requestUrl.search || decoded.includes('?')) {
    return { upstreamUrl: decoded };
  }

  return { upstreamUrl: `${decoded}${requestUrl.search}` };
}

function buildUpstreamRequest(request, upstreamUrl) {
  const headers = new Headers();

  for (const key of FORWARDED_REQUEST_HEADERS) {
    const value = request.headers.get(key);
    if (value) {
      headers.set(key, value);
    }
  }

  return new Request(upstreamUrl.toString(), {
    method: request.method,
    headers,
    redirect: 'follow',
  });
}

function handlePreflight(request) {
  const requestedHeaders = request.headers.get('Access-Control-Request-Headers');

  const headers = new Headers({
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Methods': 'GET,HEAD,OPTIONS',
    'Access-Control-Allow-Headers': requestedHeaders || 'Range,Content-Type,Accept,Origin',
    'Access-Control-Max-Age': '86400',
    Vary: 'Access-Control-Request-Headers',
  });

  return new Response(null, {
    status: 204,
    headers,
  });
}

function withCors(response) {
  const headers = new Headers(response.headers);

  headers.set('Access-Control-Allow-Origin', '*');
  headers.set('Access-Control-Allow-Methods', 'GET,HEAD,OPTIONS');
  headers.set('Access-Control-Expose-Headers', CORS_EXPOSE_HEADERS.join(', '));
  headers.set('Timing-Allow-Origin', '*');

  // Prevent passing cookies from upstream through a public media proxy.
  headers.delete('set-cookie');
  headers.delete('set-cookie2');

  return new Response(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers,
  });
}

function parseAllowedHosts(rawValue) {
  if (!rawValue || typeof rawValue !== 'string') {
    return [];
  }

  return rawValue
    .split(',')
    .map((item) => item.trim().toLowerCase())
    .filter((item) => item.length > 0 && item !== '*');
}

function isHostAllowed(hostname, allowedHosts) {
  for (const allowed of allowedHosts) {
    if (allowed.startsWith('*.')) {
      const root = allowed.slice(2);
      if (hostname === root || hostname.endsWith(`.${root}`)) {
        return true;
      }
      continue;
    }

    if (hostname === allowed) {
      return true;
    }
  }

  return false;
}

function isBlockedLocalHost(hostname) {
  return hostname === 'localhost' || hostname.endsWith('.localhost');
}

function isBlockedAddressLiteral(hostname) {
  const ipv6 = hostname.replace(/^\[|\]$/g, '');
  if (ipv6.includes(':')) {
    const normalized = ipv6.toLowerCase();
    return (
      normalized === '::1' ||
      normalized.startsWith('fe80:') ||
      normalized.startsWith('fc') ||
      normalized.startsWith('fd')
    );
  }

  const match = hostname.match(/^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$/);
  if (!match) {
    return false;
  }

  const octets = match.slice(1).map(Number);
  if (octets.some((value) => Number.isNaN(value) || value < 0 || value > 255)) {
    return true;
  }

  const [a, b] = octets;
  if (a === 10 || a === 127 || a === 0) return true;
  if (a === 169 && b === 254) return true;
  if (a === 172 && b >= 16 && b <= 31) return true;
  if (a === 192 && b === 168) return true;

  return false;
}

function safeDecodeURIComponent(value) {
  try {
    return decodeURIComponent(value);
  } catch {
    return null;
  }
}

function isHttpProtocol(url) {
  return url.protocol === 'http:' || url.protocol === 'https:';
}

function jsonError(status, message) {
  const headers = new Headers({
    'Content-Type': 'application/json; charset=utf-8',
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Methods': 'GET,HEAD,OPTIONS',
  });

  return new Response(JSON.stringify({ error: message }), {
    status,
    headers,
  });
}
