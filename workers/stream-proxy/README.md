# stream-proxy Worker

Cloudflare Worker edge proxy for HLS playlists (`.m3u8`) and media segments with:
- CORS headers for browser playback
- preflight (`OPTIONS`) handling
- upstream passthrough for `GET` and `HEAD`
- open-proxy protection via `ALLOWED_UPSTREAM_HOSTS`

## Supported proxy formats

1. Query parameter:

```text
/proxy?url=https%3A%2F%2Fcdn.example.com%2Flive%2Findex.m3u8
```

2. Encoded path:

```text
/proxy/https%3A%2F%2Fcdn.example.com%2Flive%2Findex.m3u8
```

## Host whitelist (`ALLOWED_UPSTREAM_HOSTS`)

Set allowed upstream hosts as a comma-separated list.

- Exact host: `cdn.example.com`
- Wildcard suffix: `*.example-cdn.net`

Example:

```text
ALLOWED_UPSTREAM_HOSTS=cdn.example.com,edge-1.example.com,*.video.example.net
```

Notes:
- `*` (allow all) is intentionally ignored.
- `localhost` is blocked.
- local/private IP literals are blocked (for example `127.0.0.1`, `10.x.x.x`, `192.168.x.x`, `::1`).
- Only `http://` and `https://` upstream URLs are accepted.

## Example `wrangler.toml` (single env)

Create `workers/stream-proxy/wrangler.toml`:

```toml
name = "stream-proxy"
main = "worker.js"
compatibility_date = "2026-04-13"

[vars]
ALLOWED_UPSTREAM_HOSTS = "cdn.example.com,*.video.example.net"
```

## Example `wrangler.toml` (staging + production)

```toml
name = "stream-proxy"
main = "worker.js"
compatibility_date = "2026-04-13"

[vars]
ALLOWED_UPSTREAM_HOSTS = "staging-cdn.example.com"

[env.production]
name = "stream-proxy-prod"

[env.production.vars]
ALLOWED_UPSTREAM_HOSTS = "cdn.example.com,*.video.example.net"
```

## Deploy commands

From repository root:

```bash
cd workers/stream-proxy
npx wrangler deploy
```

Deploy production env:

```bash
cd workers/stream-proxy
npx wrangler deploy --env production
```

Local run:

```bash
cd workers/stream-proxy
npx wrangler dev
```

## Request examples

Proxy playlist via query param:

```bash
curl "https://stream-proxy.<your-subdomain>.workers.dev/proxy?url=https%3A%2F%2Fcdn.example.com%2Flive%2Findex.m3u8"
```

Proxy segment via query param:

```bash
curl -H "Range: bytes=0-1023" \
  "https://stream-proxy.<your-subdomain>.workers.dev/proxy?url=https%3A%2F%2Fcdn.example.com%2Flive%2Fsegment00001.ts"
```

Proxy via encoded path:

```bash
curl "https://stream-proxy.<your-subdomain>.workers.dev/proxy/https%3A%2F%2Fcdn.example.com%2Flive%2Findex.m3u8"
```

Preflight check:

```bash
curl -i -X OPTIONS \
  -H "Origin: https://example.com" \
  -H "Access-Control-Request-Method: GET" \
  "https://stream-proxy.<your-subdomain>.workers.dev/proxy?url=https%3A%2F%2Fcdn.example.com%2Flive%2Findex.m3u8"
```
