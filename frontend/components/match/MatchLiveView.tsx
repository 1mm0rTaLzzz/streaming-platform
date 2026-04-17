'use client';

import { useEffect, useState } from 'react';
import type { Match } from '@/lib/api';
import { ReconnectingWS } from '@/lib/ws';
import VideoPlayer from '@/components/player/VideoPlayer';
import ChatPanel from '@/components/chat/ChatPanel';
import LiveCountdown from '@/components/match/LiveCountdown';
import FlagDisplay from '@/components/match/FlagDisplay';
import { useTranslations } from 'next-intl';

interface Props {
  match: Match;
  locale: string;
  lang: string;
  region: string;
}

interface ScoreUpdateMessage {
  type?: string;
  match_id?: number;
  home_score?: number;
  away_score?: number;
  minute?: number;
  status?: Match['status'];
}

export default function MatchLiveView({ match, locale, lang, region }: Props) {
  const tMatch = useTranslations('match');
  const [status,    setStatus]    = useState<Match['status']>(match.status);
  const [homeScore, setHomeScore] = useState(match.home_score);
  const [awayScore, setAwayScore] = useState(match.away_score);
  const [minute,    setMinute]    = useState<number | undefined>(match.minute);

  useEffect(() => {
    setStatus(match.status);
    setHomeScore(match.home_score);
    setAwayScore(match.away_score);
    setMinute(match.minute);
  }, [match.id, match.status, match.home_score, match.away_score, match.minute]);

  useEffect(() => {
    if (typeof window === 'undefined') return;
    const ws = new ReconnectingWS(`${window.location.origin}/api/ws/live`, (data) => {
      const msg = data as ScoreUpdateMessage;
      if (msg.type !== 'score_update' || msg.match_id !== match.id) return;
      if (typeof msg.home_score === 'number') setHomeScore(msg.home_score);
      if (typeof msg.away_score === 'number') setAwayScore(msg.away_score);
      if (typeof msg.minute    === 'number') setMinute(msg.minute);
      if (msg.status) setStatus(msg.status);
    });
    return () => ws.close();
  }, [match.id]);

  const nameKey  = locale === 'ru' ? 'name_ru' : 'name_en';
  const homeName = match.home_team?.[nameKey] ?? tMatch('tbd');
  const awayName = match.away_team?.[nameKey] ?? tMatch('tbd');
  const homeFlag = match.home_team?.flag_url ?? '🏳';
  const awayFlag = match.away_team?.flag_url ?? '🏳';
  const isLive   = status === 'live' || status === 'half_time';

  return (
    <div className="space-y-6">

      {/* Score card */}
      <div
        className="rounded-2xl p-6 md:p-8 relative overflow-hidden"
        style={{
          background: 'var(--bg-low)',
          border: `1px solid ${isLive ? 'var(--live)' : 'var(--outline-subtle)'}`,
        }}
      >
        {/* Grid pattern */}
        <div
          className="absolute inset-0 pointer-events-none"
          style={{
            backgroundImage: `
              linear-gradient(var(--outline-dim) 1px, transparent 1px),
              linear-gradient(90deg, var(--outline-dim) 1px, transparent 1px)
            `,
            backgroundSize: '40px 40px',
            maskImage: 'radial-gradient(ellipse 70% 70% at 50% 50%, black 30%, transparent 100%)',
          }}
        />

        <div className="relative flex items-center justify-between gap-4 md:gap-8">
          {/* Home team */}
          <div className="flex-1 text-center md:text-right">
            <div className="mb-3 flex justify-center md:justify-end">
              <FlagDisplay src={homeFlag} name={homeName} size="xl" />
            </div>
            <div
              className="text-base md:text-2xl font-bold leading-tight"
              style={{ color: 'var(--text-hi)' }}
            >
              {homeName}
            </div>
          </div>

          {/* Score / time */}
          <div className="text-center shrink-0 px-2">
            {(isLive || status === 'finished') ? (
              <>
                <div
                  className="font-score text-4xl md:text-5xl font-medium tabular-nums"
                  style={{ color: 'var(--primary)' }}
                >
                  {homeScore}
                  <span
                    className="mx-2 text-3xl md:text-4xl"
                    style={{ color: 'var(--text-lo)' }}
                  >
                    –
                  </span>
                  {awayScore}
                </div>
                {isLive && (
                  <div
                    className="mt-2 flex items-center justify-center gap-1.5"
                    style={{ color: 'var(--live)' }}
                  >
                    <span className="live-dot w-1.5 h-1.5 rounded-full" style={{ background: 'var(--live)' }} />
                    <span className="font-score text-xs font-medium">
                      {status === 'half_time' ? tMatch('halfTime') : `${minute ?? 0}'`}
                    </span>
                  </div>
                )}
                {status === 'finished' && (
                  <div
                    className="mt-2 font-score text-[11px] tracking-widest uppercase"
                    style={{ color: 'var(--text-lo)' }}
                  >
                    {tMatch('fullTime')}
                  </div>
                )}
              </>
            ) : (
              <div>
                <div
                  className="font-score text-2xl tracking-widest"
                  style={{ color: 'var(--text-lo)' }}
                >
                  ·–·
                </div>
                <div className="mt-2 text-xs" style={{ color: 'var(--text-mid)' }}>
                  {new Date(match.scheduled_at).toLocaleString(locale, {
                    month: 'short',
                    day:   'numeric',
                    hour:  '2-digit',
                    minute:'2-digit',
                  })}
                </div>
                <LiveCountdown
                  scheduledAt={match.scheduled_at}
                  className="font-score text-xs mt-0.5"
                  style={{ color: 'var(--primary)' }}
                  startingNowLabel={tMatch('startingNow')}
                />
              </div>
            )}
          </div>

          {/* Away team */}
          <div className="flex-1 text-center md:text-left">
            <div className="mb-3 flex justify-center md:justify-start">
              <FlagDisplay src={awayFlag} name={awayName} size="xl" />
            </div>
            <div
              className="text-base md:text-2xl font-bold leading-tight"
              style={{ color: 'var(--text-hi)' }}
            >
              {awayName}
            </div>
          </div>
        </div>

        {match.venue && (
          <p
            className="relative mt-5 text-xs text-center tracking-wide"
            style={{ color: 'var(--text-lo)' }}
          >
            {match.venue}, {match.city}
          </p>
        )}
      </div>

      {/* Player + Chat */}
      <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
        <div className="lg:col-span-2 space-y-3">
          {match.streams && match.streams.length > 0 ? (
            <VideoPlayer streams={match.streams} preferredLang={lang} preferredRegion={region} />
          ) : (
            <div
              className="flex aspect-video items-center justify-center rounded-xl"
              style={{ background: 'var(--bg-card)', border: '1px solid var(--outline-subtle)' }}
            >
              <div className="text-center">
                <div className="text-4xl mb-3">📡</div>
                <p className="text-sm" style={{ color: 'var(--text-mid)' }}>{tMatch('noStreams')}</p>
              </div>
            </div>
          )}

          {match.streams && match.streams.length > 0 && (
            <div className="flex flex-wrap gap-2">
              {Array.from(new Set(match.streams.map((s) => s.language_code))).map((lc) => (
                <a
                  key={lc}
                  href={`?lang=${lc}&region=${region}`}
                  className="rounded-lg px-3 py-1 text-xs font-score font-medium transition-colors duration-200"
                  style={
                    lc === lang
                      ? { background: 'var(--primary-dim)', color: 'var(--primary)', border: '1px solid var(--primary-dim)' }
                      : { background: 'transparent', color: 'var(--text-mid)', border: '1px solid var(--outline-subtle)' }
                  }
                >
                  {lc.toUpperCase()}
                </a>
              ))}
            </div>
          )}
        </div>

        <div className="h-[480px] lg:h-[560px]">
          <ChatPanel matchId={match.id} />
        </div>
      </div>
    </div>
  );
}
