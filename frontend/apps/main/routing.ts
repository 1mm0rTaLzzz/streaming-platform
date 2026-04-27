import { defineRouting } from 'next-intl/routing';

export const routing = defineRouting({
  locales: ['en', 'ru', 'es', 'pt', 'ar', 'fr', 'de', 'zh', 'ja', 'ko'],
  defaultLocale: 'en',
  localeDetection: true,
});
