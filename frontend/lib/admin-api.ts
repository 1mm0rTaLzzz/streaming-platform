const BASE = typeof window !== 'undefined' ? window.location.origin : '';

function getToken() {
  return typeof window !== 'undefined' ? localStorage.getItem('wc26_admin_token') ?? '' : '';
}

function authHeaders(): Record<string, string> {
  return { 'Content-Type': 'application/json', Authorization: `Bearer ${getToken()}` };
}

export const adminApi = {
  login: (username: string, password: string) =>
    fetch(`${BASE}/api/admin/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    }).then((r) => r.json()),

  matches: {
    list: () => fetch(`${BASE}/api/admin/matches`, { headers: authHeaders() }).then((r) => r.json()),
    create: (data: Record<string, unknown>) =>
      fetch(`${BASE}/api/admin/matches`, { method: 'POST', headers: authHeaders(), body: JSON.stringify(data) }).then((r) => r.json()),
    update: (id: number, data: Record<string, unknown>) =>
      fetch(`${BASE}/api/admin/matches/${id}`, { method: 'PUT', headers: authHeaders(), body: JSON.stringify(data) }).then((r) => r.json()),
    delete: (id: number) =>
      fetch(`${BASE}/api/admin/matches/${id}`, { method: 'DELETE', headers: authHeaders() }),
    updateScore: (id: number, homeScore: number, awayScore: number, status: string, minute?: number) =>
      fetch(`${BASE}/api/admin/matches/${id}/score`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({
          home_score: homeScore,
          away_score: awayScore,
          status,
          ...(typeof minute === 'number' ? { minute } : {}),
        }),
      }).then((r) => r.json()),
  },

  streams: {
    list: () => fetch(`${BASE}/api/admin/streams`, { headers: authHeaders() }).then((r) => r.json()),
    create: (data: Record<string, unknown>) =>
      fetch(`${BASE}/api/admin/streams`, { method: 'POST', headers: authHeaders(), body: JSON.stringify(data) }).then((r) => r.json()),
    update: (id: number, data: Record<string, unknown>) =>
      fetch(`${BASE}/api/admin/streams/${id}`, { method: 'PUT', headers: authHeaders(), body: JSON.stringify(data) }).then((r) => r.json()),
    delete: (id: number) =>
      fetch(`${BASE}/api/admin/streams/${id}`, { method: 'DELETE', headers: authHeaders() }),
  },

  streamLaunch: {
    start: (data: { url: string; key: string; match_id?: number; width?: number; height?: number; fps?: number }) =>
      fetch(`${BASE}/api/admin/stream/launch`, { method: 'POST', headers: authHeaders(), body: JSON.stringify(data) }).then((r) => r.json()),
    stop: (key?: string) =>
      fetch(`${BASE}/api/admin/stream/stop`, { method: 'POST', headers: authHeaders(), body: JSON.stringify({ key }) }).then((r) => r.json()),
    status: () =>
      fetch(`${BASE}/api/admin/stream/status`, { headers: authHeaders() }).then((r) => r.json()),
    debug: (enable: boolean, url?: string) =>
      fetch(`${BASE}/api/admin/stream/debug`, { method: 'POST', headers: authHeaders(), body: JSON.stringify({ enable, url }) }).then((r) => r.json()),
  },

  mirrors: {
    list: () => fetch(`${BASE}/api/admin/mirrors`, { headers: authHeaders() }).then((r) => r.json()),
    create: (data: Record<string, unknown>) =>
      fetch(`${BASE}/api/admin/mirrors`, { method: 'POST', headers: authHeaders(), body: JSON.stringify(data) }).then((r) => r.json()),
    update: (id: number, data: Record<string, unknown>) =>
      fetch(`${BASE}/api/admin/mirrors/${id}`, { method: 'PUT', headers: authHeaders(), body: JSON.stringify(data) }).then((r) => r.json()),
    delete: (id: number) =>
      fetch(`${BASE}/api/admin/mirrors/${id}`, { method: 'DELETE', headers: authHeaders() }),
    activate: (id: number) =>
      fetch(`${BASE}/api/admin/mirrors/${id}/activate`, { method: 'POST', headers: authHeaders() }).then((r) => r.json()),
  },
};
