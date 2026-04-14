# Graph Report - /Users/aleksandrterehin/Programming/DS_ML/stream_site/graphify-out/..  (2026-04-14)

## Corpus Check
- Corpus is ~21,707 words - fits in a single context window. You may not need a graph.

## Summary
- 252 nodes · 298 edges · 44 communities detected
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 27 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Deploy & Live Match Pipeline|Deploy & Live Match Pipeline]]
- [[_COMMUNITY_Mirror Health & Config|Mirror Health & Config]]
- [[_COMMUNITY_Match Handler & Models|Match Handler & Models]]
- [[_COMMUNITY_Admin API & Stream Data Model|Admin API & Stream Data Model]]
- [[_COMMUNITY_Admin Mirrors UI|Admin Mirrors UI]]
- [[_COMMUNITY_Admin Pages|Admin Pages]]
- [[_COMMUNITY_API Fetch & Mirror Discovery|API Fetch & Mirror Discovery]]
- [[_COMMUNITY_Health Handler & Repository|Health Handler & Repository]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 17|Community 17]]
- [[_COMMUNITY_Community 18|Community 18]]
- [[_COMMUNITY_Community 19|Community 19]]
- [[_COMMUNITY_Community 20|Community 20]]
- [[_COMMUNITY_Community 21|Community 21]]
- [[_COMMUNITY_Community 22|Community 22]]
- [[_COMMUNITY_Community 23|Community 23]]
- [[_COMMUNITY_Community 24|Community 24]]
- [[_COMMUNITY_Community 25|Community 25]]
- [[_COMMUNITY_Community 26|Community 26]]
- [[_COMMUNITY_Community 27|Community 27]]
- [[_COMMUNITY_Community 28|Community 28]]
- [[_COMMUNITY_Community 29|Community 29]]
- [[_COMMUNITY_Community 30|Community 30]]
- [[_COMMUNITY_Community 31|Community 31]]
- [[_COMMUNITY_Community 32|Community 32]]
- [[_COMMUNITY_Community 33|Community 33]]
- [[_COMMUNITY_Community 34|Community 34]]
- [[_COMMUNITY_Community 35|Community 35]]
- [[_COMMUNITY_Community 36|Community 36]]
- [[_COMMUNITY_Community 37|Community 37]]
- [[_COMMUNITY_Community 38|Community 38]]
- [[_COMMUNITY_Community 39|Community 39]]
- [[_COMMUNITY_Community 40|Community 40]]
- [[_COMMUNITY_Community 41|Community 41]]
- [[_COMMUNITY_Community 42|Community 42]]
- [[_COMMUNITY_Community 43|Community 43]]

## God Nodes (most connected - your core abstractions)
1. `MirrorHandler` - 12 edges
2. `MatchHandler` - 12 edges
3. `apiFetch — Mirror-Aware Fetch with DoH Fallback` - 11 edges
4. `MirrorHealthChecker Service` - 11 edges
5. `fetch()` - 11 edges
6. `Server Main Entrypoint` - 9 edges
7. `TelegramNotifier Service` - 9 edges
8. `HealthHandler` - 8 edges
9. `StreamHandler` - 7 edges
10. `Anti-Blocking Architecture` - 7 edges

## Surprising Connections (you probably didn't know these)
- `handlePassthrough Origin Proxy` --conceptually_related_to--> `Server Main Entrypoint`  [INFERRED]
  scripts/cloudflare-worker.js → backend/cmd/server/main.go
- `find_streams.py Stream Discovery Utility` --conceptually_related_to--> `HLS.js Video Player`  [INFERRED]
  scripts/find_streams.py → PLAN.md
- `Mirror Activate Button (Standby to Primary Promotion)` --implements--> `Mirror System Anti-Blocking`  [INFERRED]
  frontend/app/admin/mirrors/page.tsx → PLAN.md
- `Admin Mirrors Management Page` --implements--> `Rationale: Cold standby domains for instant activation`  [INFERRED]
  frontend/app/admin/mirrors/page.tsx → PLAN.md
- `VideoPlayer Language Filter and Preferred Lang Sort` --implements--> `i18n Global Audience Localization`  [INFERRED]
  frontend/components/player/VideoPlayer.tsx → PLAN.md
- `CF-IPCountry Header Regional Stream Auto-Selection` --implements--> `Geolocation-Based Regional Routing`  [INFERRED]
  frontend/app/[locale]/matches/[id]/page.tsx → PLAN.md
- `VideoPlayer Region-Based Stream Sorting` --implements--> `Geolocation-Based Regional Routing`  [INFERRED]
  frontend/components/player/VideoPlayer.tsx → PLAN.md
- `Admin Live Score Push via WebSocket` --references--> `WebSocket Hub Redis PubSub Fan-Out`  [INFERRED]
  frontend/app/admin/matches/page.tsx → PLAN.md
- `MatchLiveView WebSocket Live Score Updates` --references--> `WebSocket Hub Redis PubSub Fan-Out`  [INFERRED]
  frontend/components/match/MatchLiveView.tsx → PLAN.md
- `Telegram Alert on Mirror Activation` --references--> `Telegram Bot Mirror Alerts`  [INFERRED]
  frontend/app/admin/mirrors/page.tsx → PLAN.md

## Hyperedges (group relationships)
- **Geolocation Locale Routing Pipeline** — ext_cfipipcountry, middleware_middleware, geo_resolvelocale, geo_maplocale, ext_nextintl [EXTRACTED 0.95]
- **Mirror Resilience and DoH Fallback Chain** — api_apifetch, api_discovermirrordns, ext_cloudflare_doh, concept_mirrorchain, api_mirror [EXTRACTED 0.95]
- **FIFA WC26 Stream Data Model** — api_match, api_team, api_stream, api_group, api_standing, api_mirror [EXTRACTED 0.98]
- **Admin Management API Surface** — adminapi_adminapi, api_match, api_stream, api_mirror [INFERRED 0.85]
- **Anti-Piracy Scanner Stealth Pattern** — robots_robots, concept_stealthseo, middleware_middleware [INFERRED 0.80]
- **Telegram Alert Pipeline** — config_telegramtoken, config_telegramchatid, service_telegramnotifier, handler_mirroractivate, service_handlestatuschange [EXTRACTED 0.97]
- **Auto-Promotion Flow** — service_mirrorhealthchecker, service_handlestatuschange, service_promotenexthealthy, service_telegramnotifier, model_mirror [EXTRACTED 0.98]
- **Cloudflare Edge HLS Proxy** — scripts_cfworker, scripts_cfworker_validateupstream, scripts_cfworker_handlestream, scripts_cfworker_rewritem3u8, scripts_cfworker_handlesegment [EXTRACTED 0.98]
- **k6 Load Testing Scenarios** — scripts_loadtest, scripts_loadtest_smoke, scripts_loadtest_load, scripts_loadtest_stress [EXTRACTED 1.00]
- **MatchFull Response Model** — model_matchfull, model_match, model_team, model_stream, handler_matchhandler [EXTRACTED 0.95]
- **Anti-Blocking Design: Cloudflare + iptables + diverse registrars + hidden admin + cold standby** — plan_anti_blocking, plan_cloudflare, plan_rationale_cloudflare_origin_hide, plan_rationale_iptables_cloudflare_only, plan_rationale_diverse_registrars, plan_rationale_hidden_admin, plan_rationale_cold_standby [EXTRACTED 0.95]
- **Mirror Management: Admin UI + Activate + Telegram Alert + Health Checker** — admin_mirrors_page, admin_mirrors_activate_button, admin_mirrors_telegram_alert, plan_mirror_health_checker, plan_telegram_bot [EXTRACTED 0.90]
- **Live Match Experience: MatchLiveView + VideoPlayer + ChatPanel + WebSocket + Countdown** — matchliveview_component, videoplayer_component, matchliveview_chatpanel, matchliveview_ws_score, matchliveview_livecountdown [EXTRACTED 0.95]
- **Geo-Based Stream Selection: CF-IPCountry + mapCountryToRegion + VideoPlayer Region Sort** — match_id_cf_ipcountry, match_id_geo_lib, videoplayer_region_sort, plan_geo_routing [EXTRACTED 0.90]
- **Core Tech Stack** — plan_nextjs14_frontend, plan_go_gin_backend, plan_postgresql16, plan_redis7, plan_hlsjs, plan_cloudflare [EXTRACTED 1.00]

## Communities

### Community 0 - "Deploy & Live Match Pipeline"
Cohesion: 0.08
Nodes (26): Admin Live Score Push via WebSocket, VPS Deploy Script (REGION-based SSH Docker Compose), Deploy Health Check Loop (backend /api/health), REGION Environment Variable for Deploy Targeting, Playwright Browser Automation for m3u8 Capture, find_streams.py Stream Discovery Utility, CF-IPCountry Header Regional Stream Auto-Selection, mapCountryToRegion Geo Utility (+18 more)

### Community 1 - "Mirror Health & Config"
Cohesion: 0.14
Nodes (13): Config Struct, getEnv(), Load(), TelegramChatID Config Field, TelegramToken Config Field, MirrorHandler.Activate Endpoint, Server Main Entrypoint, Mirror Model (+5 more)

### Community 2 - "Match Handler & Models"
Cohesion: 0.13
Nodes (10): MatchHandler, MatchHandler.UpdateScore, scanMatches(), AdminUser Model, Group, Match Model, MatchFull Composite Model, Standing Model (+2 more)

### Community 3 - "Admin API & Stream Data Model"
Cohesion: 0.12
Nodes (15): Admin API Client, Match Interface, Stream Interface, Team Interface, Geolocation-Based Locale Routing Pattern, Anti-Piracy Scanner Stealth SEO Pattern, CF-IPCountry Cloudflare Header, next-intl Middleware (+7 more)

### Community 4 - "Admin Mirrors UI"
Cohesion: 0.16
Nodes (14): adminApi.mirrors Client, Mirror Activate Button (Standby to Primary Promotion), Admin Mirrors Management Page, Telegram Alert on Mirror Activation, Anti-Blocking Architecture, Cloudflare CDN and GeoDNS, Mirror Health Checker Background Goroutine, Mirror System Anti-Blocking (+6 more)

### Community 5 - "Admin Pages"
Cohesion: 0.31
Nodes (10): extractApiError(), handleActivate(), handleDelete(), handleQuickStatus(), handleScorePush(), handleSubmit(), load(), resetForm() (+2 more)

### Community 6 - "API Fetch & Mirror Discovery"
Cohesion: 0.2
Nodes (13): Public API Object (matches / groups / mirrors), apiFetch — Mirror-Aware Fetch with DoH Fallback, discoverMirrorViaDns — Cloudflare DoH Mirror Discovery, discoverMirrorViaDns(), fetchWithTimeout(), getStoredMirrors(), Mirror Interface, storeMirrors() (+5 more)

### Community 7 - "Health Handler & Repository"
Cohesion: 0.16
Nodes (10): MirrorRepository Interface, HealthHandler, HealthHandler.Mirrors Geo-Sort, MirrorRepository, countryToRegion(), sortMirrorsByRegion(), k6 Load Test Script, k6 Load Scenario (+2 more)

### Community 8 - "Community 8"
Cohesion: 0.29
Nodes (12): buildUpstreamRequest(), extractUpstreamFromRequest(), fetch(), handlePreflight(), isBlockedAddressLiteral(), isBlockedLocalHost(), isHostAllowed(), isHttpProtocol() (+4 more)

### Community 9 - "Community 9"
Cohesion: 0.22
Nodes (1): MirrorHandler

### Community 10 - "Community 10"
Cohesion: 0.57
Nodes (7): corsResponse(), fetch(), handlePassthrough(), handleSegment(), handleStream(), rewriteM3u8(), validateUpstreamUrl()

### Community 11 - "Community 11"
Cohesion: 0.29
Nodes (1): StreamHandler

### Community 12 - "Community 12"
Cohesion: 0.4
Nodes (6): Cloudflare Worker HLS Proxy, handleSegment TS Proxy, handleStream m3u8 Proxy, handlePassthrough Origin Proxy, rewriteM3u8 Segment URL Rewriter, validateUpstreamUrl Domain Allowlist

### Community 13 - "Community 13"
Cohesion: 0.6
Nodes (5): ensureSupportedLocale(), mapCountryToLocale(), mapCountryToRegion(), normalizeCountryCode(), resolveLocaleFromCookieOrCountry()

### Community 14 - "Community 14"
Cohesion: 0.67
Nodes (3): adminApi.matches Client, Admin Matches Management Page, Admin Quick Match Status (Start/Halftime/End)

### Community 15 - "Community 15"
Cohesion: 0.67
Nodes (0): 

### Community 16 - "Community 16"
Cohesion: 0.67
Nodes (0): 

### Community 17 - "Community 17"
Cohesion: 0.67
Nodes (0): 

### Community 18 - "Community 18"
Cohesion: 0.67
Nodes (0): 

### Community 19 - "Community 19"
Cohesion: 0.67
Nodes (0): 

### Community 20 - "Community 20"
Cohesion: 0.67
Nodes (0): 

### Community 21 - "Community 21"
Cohesion: 1.0
Nodes (2): authHeaders(), getToken()

### Community 22 - "Community 22"
Cohesion: 1.0
Nodes (2): find_m3u8(), main()

### Community 23 - "Community 23"
Cohesion: 1.0
Nodes (2): Group Interface, Standing Interface

### Community 24 - "Community 24"
Cohesion: 1.0
Nodes (2): Database Schema: groups teams matches streams mirrors, PostgreSQL 16 Database

### Community 25 - "Community 25"
Cohesion: 1.0
Nodes (0): 

### Community 26 - "Community 26"
Cohesion: 1.0
Nodes (0): 

### Community 27 - "Community 27"
Cohesion: 1.0
Nodes (0): 

### Community 28 - "Community 28"
Cohesion: 1.0
Nodes (0): 

### Community 29 - "Community 29"
Cohesion: 1.0
Nodes (0): 

### Community 30 - "Community 30"
Cohesion: 1.0
Nodes (0): 

### Community 31 - "Community 31"
Cohesion: 1.0
Nodes (0): 

### Community 32 - "Community 32"
Cohesion: 1.0
Nodes (1): FIFA WC2026 Live Streaming Aggregator Platform

### Community 33 - "Community 33"
Cohesion: 1.0
Nodes (1): Next.js 14 Frontend

### Community 34 - "Community 34"
Cohesion: 1.0
Nodes (1): Go + Gin Backend

### Community 35 - "Community 35"
Cohesion: 1.0
Nodes (1): Nginx Reverse Proxy

### Community 36 - "Community 36"
Cohesion: 1.0
Nodes (1): Phase 4: Resilience (Planned)

### Community 37 - "Community 37"
Cohesion: 1.0
Nodes (1): Phase 5: Production (Planned)

### Community 38 - "Community 38"
Cohesion: 1.0
Nodes (1): Locale Homepage: Live Now + Todays Matches

### Community 39 - "Community 39"
Cohesion: 1.0
Nodes (0): 

### Community 40 - "Community 40"
Cohesion: 1.0
Nodes (0): 

### Community 41 - "Community 41"
Cohesion: 1.0
Nodes (0): 

### Community 42 - "Community 42"
Cohesion: 1.0
Nodes (0): 

### Community 43 - "Community 43"
Cohesion: 1.0
Nodes (0): 

## Knowledge Gaps
- **47 isolated node(s):** `Next.js Config (Rewrites & next-intl)`, `Public API Object (matches / groups / mirrors)`, `Team Interface`, `Group Interface`, `Standing Interface` (+42 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Community 23`** (2 nodes): `Group Interface`, `Standing Interface`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 24`** (2 nodes): `Database Schema: groups teams matches streams mirrors`, `PostgreSQL 16 Database`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 25`** (2 nodes): `layout.tsx`, `RootLayout()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 26`** (2 nodes): `page.tsx`, `getMatches()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 27`** (2 nodes): `page.tsx`, `getGroups()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 28`** (2 nodes): `page.tsx`, `AdminLoginPage()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 29`** (2 nodes): `Header.tsx`, `switchLocale()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 30`** (2 nodes): `ServiceWorkerRegister.tsx`, `ServiceWorkerRegister()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 31`** (2 nodes): `main.go`, `main()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 32`** (1 nodes): `FIFA WC2026 Live Streaming Aggregator Platform`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 33`** (1 nodes): `Next.js 14 Frontend`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 34`** (1 nodes): `Go + Gin Backend`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 35`** (1 nodes): `Nginx Reverse Proxy`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 36`** (1 nodes): `Phase 4: Resilience (Planned)`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 37`** (1 nodes): `Phase 5: Production (Planned)`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 38`** (1 nodes): `Locale Homepage: Live Now + Todays Matches`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 39`** (1 nodes): `next.config.js`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 40`** (1 nodes): `page.tsx`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 41`** (1 nodes): `layout.tsx`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 42`** (1 nodes): `MatchLiveView.tsx`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 43`** (1 nodes): `load-test.k6.js`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Server Main Entrypoint` connect `Community 1 (24 nodes)` to `Community 2 (20 nodes)`, `Community 7 (14 nodes)`, `Community 9 (10 nodes)`, `Community 11 (7 nodes)`, `Community 12 (6 nodes)`?**
  _High betweenness centrality (0.055) - this node is a cross-community bridge._
- **Why does `MatchHandler` connect `Community 2 (20 nodes)` to `Community 1 (24 nodes)`, `Community 7 (14 nodes)`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Why does `HealthHandler` connect `Community 7 (14 nodes)` to `Community 1 (24 nodes)`, `Community 12 (6 nodes)`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `apiFetch — Mirror-Aware Fetch with DoH Fallback` (e.g. with `Next.js Config (Rewrites & next-intl)` and `ReconnectingWS — Auto-Reconnect WebSocket Client`) actually correct?**
  _`apiFetch — Mirror-Aware Fetch with DoH Fallback` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Next.js Config (Rewrites & next-intl)`, `Public API Object (matches / groups / mirrors)`, `Team Interface` to the rest of the system?**
  _47 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 0 (26 nodes)` be split into smaller, more focused modules?**
  _Cohesion score 0.08 - nodes in this community are weakly interconnected._
- **Should `Community 1 (24 nodes)` be split into smaller, more focused modules?**
  _Cohesion score 0.14 - nodes in this community are weakly interconnected._