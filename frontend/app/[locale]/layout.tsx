import type { Metadata } from 'next';
import { NextIntlClientProvider } from 'next-intl';
import { getMessages, getTranslations } from 'next-intl/server';
import { Lexend, Manrope } from 'next/font/google';
import Header from '@/components/layout/Header';
import { routing } from '@/routing';

const lexend  = Lexend({ subsets: ['latin'], weight: ['300','400','500','600','700','800'], variable: '--font-display' });
const manrope = Manrope({ subsets: ['latin', 'cyrillic'], weight: ['300','400','500','600','700','800'], variable: '--font-body' });

export async function generateMetadata({
  params,
}: {
  params: { locale: string };
}): Promise<Metadata> {
  const locale = routing.locales.includes(params.locale as (typeof routing.locales)[number])
    ? params.locale
    : routing.defaultLocale;

  const baseUrl = process.env.NEXT_PUBLIC_CANONICAL_URL ?? 'https://wc26.live';
  const languageAlternates = Object.fromEntries(
    routing.locales.map((supportedLocale) => [supportedLocale, `${baseUrl}/${supportedLocale}`]),
  );
  const tMeta = await getTranslations({ locale, namespace: 'meta' });

  return {
    title: tMeta('title'),
    description: tMeta('description'),
    alternates: {
      canonical: `${baseUrl}/${locale}`,
      languages: {
        ...languageAlternates,
        'x-default': `${baseUrl}/${routing.defaultLocale}`,
      },
    },
  };
}

export default async function LocaleLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: { locale: string };
}) {
  const { locale } = params;
  const messages = await getMessages();
  const tFooter = await getTranslations({ locale, namespace: 'footer' });
  const isRTL = ['ar', 'he', 'fa', 'ur'].includes(locale);

  return (
    <div lang={locale} dir={isRTL ? 'rtl' : 'ltr'} className={`${lexend.variable} ${manrope.variable}`}>
        <NextIntlClientProvider messages={messages}>
          <Header locale={locale} />
          <main className="max-w-7xl mx-auto px-4 py-8 md:py-10">{children}</main>
          <footer className="mt-20 py-8 px-4 text-center" style={{ borderTop: '1px solid var(--border-faint)' }}>
            <p className="text-[11px] tracking-wide" style={{ color: 'var(--text-lo)', maxWidth: '640px', margin: '0 auto' }}>
              {tFooter('disclaimer')}
            </p>
          </footer>
        </NextIntlClientProvider>
    </div>
  );
}
