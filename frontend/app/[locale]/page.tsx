import { api, type Match } from '@/lib/api';
import MatchCard from '@/components/match/MatchCard';
import { getTranslations } from 'next-intl/server';

export const dynamic = 'force-dynamic';
export const revalidate = 0;

async function getMatches(status?: string): Promise<Match[]> {
  try {
    const today = new Date().toISOString().slice(0, 10);
    const params: Record<string, string> = status ? { status } : { date: today };
    const res = await api.matches.list(params);
    return res.matches ?? [];
  } catch {
    return [];
  }
}

export default async function HomePage({ params }: { params: { locale: string } }) {
  const { locale } = params;
  const tHome  = await getTranslations({ locale, namespace: 'home' });
  const tMatch = await getTranslations({ locale, namespace: 'match' });
  const [liveMatches, halfTimeMatches, todayMatches] = await Promise.all([
    getMatches('live'),
    getMatches('half_time'),
    getMatches(),
  ]);
  const liveNow  = [...liveMatches, ...halfTimeMatches];
  const upcoming = todayMatches.filter((m) => m.status === 'scheduled');

  const matchLabels = {
    tbd:         tMatch('tbd'),
    halfTime:    tMatch('halfTime'),
    fullTime:    tMatch('finished'),
    versus:      tMatch('versus'),
    startingNow: tMatch('startingNow'),
  };

  return (
    <div className="space-y-10">

      {/* Hero */}
      <div
        className="relative rounded-2xl overflow-hidden"
        style={{ background: 'var(--bg-low)' }}
      >
        {/* Subtle grid */}
        <div
          className="absolute inset-0 pointer-events-none"
          style={{
            backgroundImage: `
              linear-gradient(rgba(255,255,255,0.02) 1px, transparent 1px),
              linear-gradient(90deg, rgba(255,255,255,0.02) 1px, transparent 1px)
            `,
            backgroundSize: '48px 48px',
            maskImage: 'radial-gradient(ellipse 80% 80% at 50% 50%, black 30%, transparent 100%)',
          }}
        />
        {/* Primary glow */}
        <div
          className="absolute inset-x-0 bottom-0 h-48 pointer-events-none"
          style={{ background: 'radial-gradient(ellipse 60% 100% at 50% 100%, var(--primary-glow), transparent)' }}
        />

        <div className="relative text-center px-6 py-12 md:py-16">
          {/* Host flags */}
          <div className="flex items-center justify-center gap-3 mb-4 text-3xl md:text-4xl">
            🇺🇸 🇨🇦 🇲🇽
          </div>

          <p
            className="text-xs font-bold uppercase tracking-widest mb-3"
            style={{ color: 'var(--secondary)', fontFamily: 'var(--font-body)' }}
          >
            North America Journey
          </p>

          <h1
            className="font-display font-bold italic leading-none mb-6"
            style={{
              fontSize: 'clamp(3.5rem, 12vw, 8rem)',
              color:     'var(--primary)',
              textShadow: '0 0 60px rgba(149,170,255,0.35)',
            }}
          >
            WORLD CUP 2026
          </h1>

          <div className="flex items-center justify-center gap-3">
            <a
              href={`/${locale}/schedule`}
              className="btn-primary inline-flex items-center gap-2 px-5 py-2.5 rounded-lg text-sm font-bold transition-opacity hover:opacity-90"
            >
              {tHome('viewFullSchedule')}
            </a>
            <a
              href={`/${locale}/groups`}
              className="inline-flex items-center gap-2 px-5 py-2.5 rounded-lg text-sm font-bold transition-all"
              style={{
                background: 'var(--bg-bright)',
                color:      'var(--text-hi)',
                fontFamily: 'var(--font-body)',
              }}
            >
              {tHome('groups')}
            </a>
          </div>
        </div>
      </div>

      {/* Live */}
      {liveNow.length > 0 && (
        <section>
          <div className="flex items-center gap-2 mb-4">
            <span className="live-dot w-2 h-2 rounded-full shrink-0" style={{ background: 'var(--secondary)' }} />
            <span
              className="text-sm font-bold"
              style={{ color: 'var(--secondary)', fontFamily: 'var(--font-body)' }}
            >
              {liveNow.length === 1 ? '1 match in progress' : `${liveNow.length} matches in progress`}
            </span>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
            {liveNow.map((m, i) => (
              <div key={m.id} className="fade-up" style={{ animationDelay: `${i * 0.06}s` }}>
                <MatchCard match={m} locale={locale} labels={matchLabels} />
              </div>
            ))}
          </div>
        </section>
      )}

      {/* Today */}
      {upcoming.length > 0 && (
        <section>
          <h2
            className="text-sm font-bold mb-4"
            style={{ color: 'var(--text-mid)', fontFamily: 'var(--font-body)' }}
          >
            {tHome('todayMatches')}
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
            {upcoming.map((m, i) => (
              <div key={m.id} className="fade-up" style={{ animationDelay: `${i * 0.06}s` }}>
                <MatchCard match={m} locale={locale} labels={matchLabels} />
              </div>
            ))}
          </div>
        </section>
      )}

      {/* Empty */}
      {liveNow.length === 0 && upcoming.length === 0 && (
        <div className="text-center py-24">
          <div className="text-5xl mb-4">⚽</div>
          <p className="text-base font-semibold" style={{ color: 'var(--text-mid)', fontFamily: 'var(--font-body)' }}>
            {tHome('noMatchesToday')}
          </p>
          <a
            href={`/${locale}/schedule`}
            className="inline-block mt-3 text-sm font-bold transition-opacity hover:opacity-75"
            style={{ color: 'var(--primary)', fontFamily: 'var(--font-body)' }}
          >
            {tHome('viewFullSchedule')} →
          </a>
        </div>
      )}
    </div>
  );
}
