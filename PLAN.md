# Streaming Site - FIFA World Cup 2026 - Architecture Plan

## Context

Строим ГЛОБАЛЬНУЮ live-стриминг платформу для трансляции матчей ЧМ-2026 (48 команд, США/Канада/Мексика, июнь-июль 2026). Сайт агрегирует m3u8 потоки (не хостит видео сам). Нужна устойчивость к блокировкам через систему зеркал. **Аудитория — весь мир.** Стримы на множестве языков, UI с i18n, геолокационная маршрутизация.

**Обновление требований (апрель 2026):** сайт должен обслуживать глобальную аудиторию:
- Стримы с комментарием на 20+ языках (EN, RU, ES, PT, AR, FR, DE, ZH, JA, KO и др.)
- UI с полной локализацией (next-intl)
- Геолокация пользователя → показ стримов на родном языке первыми
- Региональные зеркала (EU, Americas, Asia, MENA)
- RTL поддержка для арабского/иврита
- SEO на всех языках (alternates hreflang)

---

## Tech Stack

| Layer | Technology | Why |
|-------|-----------|-----|
| Frontend | **Next.js 14** (standalone) | App Router, next-intl i18n, легко зеркалить |
| Backend | **Go + Gin** | Один бинарник, тысячи WebSocket соединений, низкое потребление памяти |
| Database | **PostgreSQL 16** | Реляционные данные (матчи/команды/группы), логическая репликация |
| Cache/PubSub | **Redis 7** | Чат (Pub/Sub), rate limiting, кеш зеркал |
| Video | **HLS.js** | Клиентский плеер m3u8 потоков, нет собственного видеохостинга |
| Proxy | **Nginx** | Отдает статику + проксирует /api к Go, WebSocket upgrade |
| CDN/Shield | **Cloudflare** | Скрывает origin IP, DDoS защита, SSL, GeoDNS |
| Deploy | **Docker Compose** | Все сервисы в контейнерах |

---

## Architecture Overview

```
Global Users
    |
[Cloudflare Global CDN + GeoDNS]
    |
    ├── EU Region     --> [VPS EU: Nginx -> Go + PG replica + Redis]
    ├── Americas      --> [VPS US: Nginx -> Go + PG primary + Redis]
    ├── Asia          --> [VPS SG: Nginx -> Go + PG replica + Redis]
    └── MENA          --> [VPS AE: Nginx -> Go + PG replica + Redis]
                                 ↑
                        PG logical replication
```

---

## File Structure (реальная, как построено)

```
stream_site/
  backend/
    cmd/server/main.go              -- Entry point, Gin router, graceful shutdown
    internal/
      config/config.go              -- Env-based config
      handler/
        auth.go                     -- POST /api/admin/login (JWT)
        chat.go                     -- WS /api/ws/chat/:matchId + /api/ws/live
        group.go                    -- GET /api/groups, /api/groups/:id (standings)
        health.go                   -- GET /api/health, /api/mirrors
        match.go                    -- GET/POST/PUT/DELETE /api/matches
        mirror.go                   -- Admin CRUD mirrors
        stream.go                   -- Admin CRUD streams
      middleware/
        auth.go                     -- JWT Bearer validation
        ratelimit.go                -- Redis sliding window rate limiter
      model/models.go               -- Team, Match, Stream, Mirror, AdminUser structs
      repository/
        postgres.go                 -- sqlx connection pool
        redis.go                    -- go-redis client
      service/mirror_health.go      -- Background goroutine: ping mirrors every 15s
      ws/
        hub.go                      -- WebSocket hub, Redis PubSub fan-out
        client.go                   -- Per-connection read/write pumps
    migrations/001_initial.sql      -- Schema: groups/teams/matches/streams/mirrors/admin_users
    Dockerfile                      -- Multistage: golang:1.21 → alpine:3.19 (~15MB)
  frontend/
    app/
      page.tsx                      -- Redirect / → /en
      [locale]/
        layout.tsx                  -- Root layout, next-intl, Bebas Neue + Inter fonts
        page.tsx                    -- Homepage: Live Now + Today's Matches
        groups/page.tsx             -- 12 groups, standings tables
        schedule/page.tsx           -- Full schedule grouped by date
        matches/[id]/page.tsx       -- Player page: HLS.js + chat + match info
    components/
      player/VideoPlayer.tsx        -- HLS.js, multi-source fallback, quality/lang selector
      match/MatchCard.tsx           -- Match card: flags, score, LIVE badge, countdown
      chat/ChatPanel.tsx            -- WebSocket chat, auto-reconnect, username
      layout/Header.tsx             -- Nav + language selector (10 locales)
    lib/
      api.ts                        -- Mirror-aware fetch, typed helpers, localStorage cache
      ws.ts                         -- ReconnectingWS class with exponential backoff
    messages/
      en.json / ru.json / es.json   -- UI translations (10 languages)
    i18n.ts                         -- next-intl config
    next.config.js                  -- standalone output, /api/* rewrite to backend
  telegram-bot/main.go              -- Polling bot: /mirrors command, mirror alerts
  docker/
    docker-compose.yml              -- PG 16 + Redis 7 + backend + nginx
    nginx/nginx.conf                -- Reverse proxy, WS upgrade, static files, gzip
    .env.example                    -- JWT_SECRET, ADMIN_PASS, TELEGRAM_BOT_TOKEN
  scripts/
    seed-teams.sql                  -- 48 команд ЧМ-2026, 12 групп A-L
    seed-admin.sql                  -- Default admin user + localhost mirror
  graphify-out/
    graph.html                      -- Интерактивный граф архитектуры (89 узлов)
    GRAPH_REPORT.md                 -- Audit report
    graph.json                      -- Raw graph data
```

---

## Database Schema

```sql
groups        (id, name A-L)
teams         (id, code, name_en, name_ru, flag_url, group_id)
matches       (id, home_team_id, away_team_id, group_id, stage, venue, city,
               scheduled_at, status, home_score, away_score, minute)
streams       (id, match_id, url, label, language_code BCP47, region, 
               commentary_type, quality, is_active, priority)
mirrors       (id, domain, is_active, is_primary, last_health_check, health_status, region)
admin_users   (id, username, password_hash, role)
```

Stream languages: `en ru es pt ar fr de zh ja ko` + любые другие BCP47  
Stream regions: `global eu us asia mena`  
Stream commentary: `full no_commentary home_team`

---

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/health` | — | Healthcheck + mirrors list |
| GET | `/api/mirrors` | — | Активные зеркала |
| GET | `/api/matches` | — | Матчи (?status=live, ?date=2026-06-11, ?stage=group) |
| GET | `/api/matches/:id` | — | Матч + стримы (?lang=ru&region=eu) |
| GET | `/api/groups` | — | Все группы со standings |
| GET | `/api/groups/:id` | — | Группа + standings |
| WS | `/api/ws/chat/:matchId` | — | Чат (?username=xxx) |
| WS | `/api/ws/live` | — | Live score updates |
| POST | `/api/admin/login` | — | JWT login |
| CRUD | `/api/admin/matches` | JWT | Управление матчами |
| CRUD | `/api/admin/streams` | JWT | Управление потоками |
| CRUD | `/api/admin/mirrors` | JWT | Управление зеркалами |

---

## Key Features

### HLS Video Player
- HLS.js с fallback на нативный HLS (Safari/iOS)
- Массив источников — при `NETWORK_ERROR` автоматически пробует следующий
- Фильтр по языку: показывает стримы на языке интерфейса первыми
- Quality selector из HLS manifest levels (Auto/1080p/720p/480p)

### Mirror System
- 5–10 доменов на разных TLD (.com .tv .live .stream) и регистраторах
- Каждый домен → отдельный Cloudflare zone → один origin VPS
- Background goroutine: ping каждые 15 сек, обновляет `health_status` в БД
- Frontend API client: пробует `window.origin` → localStorage mirrors → hardcoded fallback
- При смене домена: Telegram бот оповещает канал

### Anti-Blocking
- Cloudflare proxy: origin IP никогда не раскрывается клиентам
- `iptables` на VPS: принимает только трафик с IP Cloudflare
- Разные регистраторы + privacy WHOIS
- Admin panel на отдельном скрытом домене с IP whitelist
- Cold standby домены готовы к мгновенной активации

### i18n (Global Audience)
- Locale routing: `/en/`, `/ru/`, `/es/`, `/ar/`, ...
- Переключение языка в хедере — меняет URL prefix
- RTL автоматически для ar/he/fa/ur (`dir="rtl"` на `<html>`)
- Stream selector фильтрует по `language_code` стрима

---

## Deployment

```bash
cd docker/
cp .env.example .env        # заполнить JWT_SECRET, ADMIN_PASS
docker compose up -d        # поднимает PG + Redis + backend + nginx
```

- **VPS**: BuyVM / FlokiNET / 1984 Hosting (DMCA-tolerant)
- **Spec**: 2 vCPU, 4GB RAM, 80GB SSD
- **Мониторинг**: Uptime Kuma (зеркала) + Telegram алерты

---

## Implementation Status

### Phase 1 — Foundation ✅ DONE
- [x] Go backend: все хендлеры, WebSocket hub, mirror health checker
- [x] Next.js 14 frontend: 10 языков, HLS.js плеер, чат, группы, расписание
- [x] Docker Compose: PG + Redis + Go + Nginx
- [x] DB migrations + seed 48 команд

### Phase 2 — Core (следующий)
- [x] Admin panel UI (Next.js) для управления матчами/потоками
- [x] Seed полного расписания (104 матча с датами/стадионами)
- [x] Страница `/admin/login` + JWT в localStorage
- [x] Live score push от admin → WebSocket → клиенты

### Phase 3 — Live Features
- [x] Admin: кнопки start/halftime/end match, обновление счёта
- [x] Frontend: live countdown до начала матча
- [x] Service Worker: offline schedule cache

### Phase 4 — Resilience
- [ ] Cloudflare Workers script (edge proxy для m3u8 стримов)
- [ ] Mirror discovery улучшение: CF-IPCountry → regional redirect
- [ ] hreflang SEO meta tags
- [ ] Геолокация → дефолтный язык

### Phase 5 — Production
- [ ] Deploy на VPS (EU + US минимум)
- [ ] Cloudflare setup: 3+ домена, GeoDNS
- [ ] Telegram bot активация с реальным токеном
- [ ] Load testing (k6) перед ЧМ
- [ ] Uptime Kuma мониторинг всех зеркал

---

## Quick Start (локальная разработка)

```bash
# Backend
cd backend && go run ./cmd/server

# Frontend
cd frontend && npm run dev

# Вся инфраструктура
cd docker && docker compose up -d postgres redis

# Проверка
curl http://localhost:8080/api/health
open http://localhost:3000/en
```
