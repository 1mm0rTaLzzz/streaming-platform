# WC 2026 Streaming Platform

Live streaming platform for FIFA World Cup 2026. Aggregates HLS streams, supports 10+ languages, mirror system for geo-resilience, automated browser-based streaming via headless Chromium.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         VIEWERS (global)                            │
└──────────────────────────────┬──────────────────────────────────────┘
                               │ HTTPS
                ┌──────────────▼──────────────┐
                │       Cloudflare CDN        │
                │  DDoS protection, GeoDNS,   │
                │  origin IP hidden           │
                └──────────────┬──────────────┘
                               │
                ┌──────────────▼──────────────┐
                │         nginx :80           │  ← static frontend
                │   /api/*   → backend:8080   │  ← API proxy
                │   /live/*  → rtmp:80/hls/   │  ← HLS stream proxy
                │   /hls/*   → hls_proxy:8089 │  ← restream proxy
                └──┬──────┬──────┬────────────┘
                   │      │      │
        ┌──────────▼─┐  ┌─▼──────▼────────┐  ┌──────────────────┐
        │  backend   │  │   rtmp service  │  │   hls_proxy      │
        │  Go + Gin  │  │  alfg/nginx-rtmp│  │  restream proxy  │
        │  :8080     │  │  RTMP :1935     │  │  :8089           │
        │            │  │  HLS  :80       │  └──────────────────┘
        │  REST API  │  └────────▲────────┘
        │  WebSocket │           │ rtmp push
        │  JWT auth  │   ┌────────┴─────────┐
        └──┬──────┬──┘   │    streamer      │
           │      │      │  Chromium+Xvfb   │
     ┌─────▼─┐  ┌─▼──┐   │  ffmpeg capture  │
     │  PG   │  │Redis│  │  :8090 HTTP API  │
     │  :5432│  │:6379│  └────────▲─────────┘
     └───────┘  └────┘            │ CDP :9222
                                  │ (Google login setup)
                         ┌────────┴─────────┐
                         │  Admin Browser   │
                         │  SSH tunnel      │
                         └──────────────────┘
```

## Stream Flow

### OBS / External Stream
```
OBS → rtmp://server:1935/live/{key} → nginx-rtmp → HLS segments → /live/{key}.m3u8 → viewers
```

### Automated Browser Stream (headless)
```
Admin UI (/admin/streams/launch)
  → POST /api/admin/stream/launch
  → streamer service
  → Xvfb :99 + Chromium --kiosk {url}
  → ffmpeg x11grab :99 → rtmp://rtmp:1935/live/{key}
  → nginx-rtmp → HLS → viewers
```

### External HLS Aggregation
```
3rd-party m3u8 URL → stored in DB (streams table) → served via admin-added stream records → HLS.js player
```

## Services

| Service | Image | Port | Purpose |
|---------|-------|------|---------|
| `nginx` | nginx:alpine | 80, 443 | Reverse proxy + static frontend |
| `backend` | custom Go | 8080 | REST API, WebSocket hub, JWT auth |
| `postgres` | postgres:16 | 5432 | Match/team/stream/mirror data |
| `redis` | redis:7 | 6379 | Rate limiting, chat pub/sub |
| `rtmp` | alfg/nginx-rtmp | 1935, 8088 | RTMP ingest → HLS segments |
| `streamer` | custom | 8090 (internal) | Headless Chromium → RTMP |
| `hls_proxy` | pcruz1905/hls-restream | 8089 | Restream external HLS sources |

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/health` | Healthcheck |
| GET | `/api/mirrors` | Active mirror domains |
| GET | `/api/matches` | Matches (filter: date, status, stage) |
| GET | `/api/matches/:id` | Match detail + streams |
| GET | `/api/groups` | Group standings |
| WS | `/api/ws/chat/:matchId` | Live chat |
| WS | `/api/ws/live` | Score updates |
| POST | `/api/admin/login` | JWT login |
| CRUD | `/api/admin/matches` | Match management |
| CRUD | `/api/admin/streams` | Stream URL management |
| CRUD | `/api/admin/mirrors` | Mirror domain management |
| POST | `/api/admin/stream/launch` | Start headless browser stream |
| POST | `/api/admin/stream/stop` | Stop headless stream |
| GET | `/api/admin/stream/status` | Active stream status |
| POST | `/api/admin/stream/debug` | Open Chromium for Google login |

## Quick Start

```bash
cd docker
docker compose up -d

# Verify
curl http://localhost/api/health
# → {"status":"ok"}
```

Admin panel: `http://localhost:8080/admin`
Default credentials: `admin / admin123`

## Mirror System

Multiple domains on different TLDs and registrars. Backend goroutine health-checks every 15s. Frontend auto-switches on failure.

```
primary.tv  ──┐
backup.live ──┼──→ Cloudflare → origin VPS
cold.stream ──┘
```

Telegram bot sends alerts on mirror switch. Admin panel at `/admin/mirrors`.
