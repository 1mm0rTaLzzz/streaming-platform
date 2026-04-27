import { NextRequest, NextResponse } from 'next/server';
import createMiddleware from 'next-intl/middleware';
import { routing } from './routing';
import { resolveLocaleFromCookieOrCountry } from './lib/geo';

const intlMiddleware = createMiddleware(routing);

export default function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const localePattern = routing.locales.join('|');
  const hasLocalePrefixRe = new RegExp(`^/(${localePattern})(?:/|$)`);

  if (pathname === '/' || !hasLocalePrefixRe.test(pathname)) {
    const locale = resolveLocaleFromCookieOrCountry({
      cookieLocale: request.cookies.get('NEXT_LOCALE')?.value,
      countryCode: request.headers.get('CF-IPCountry'),
      supportedLocales: routing.locales,
      fallbackLocale: routing.defaultLocale,
    });

    const url = request.nextUrl.clone();
    url.pathname = pathname === '/' ? `/${locale}` : `/${locale}${pathname}`;
    return NextResponse.redirect(url);
  }

  return intlMiddleware(request);
}

export const config = {
  matcher: ['/((?!api|_next|_vercel|.*\\..*).*)'],
};
