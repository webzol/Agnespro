// Thin fetch wrapper. All paths are relative so Vite proxy works in dev
// and the Go backend serves the SPA in production.
export class ApiError extends Error {
  status: number;
  body: any;
  constructor(status: number, msg: string, body: any) {
    super(msg); this.status = status; this.body = body;
  }
}

async function req<T>(method: string, path: string, body?: any): Promise<T> {
  const opts: RequestInit = {
    method,
    headers: body !== undefined ? { "Content-Type": "application/json" } : {},
  };
  if (body !== undefined) opts.body = JSON.stringify(body);
  const r = await fetch(path, opts);
  if (!r.ok) {
    let msg = `${r.status} ${r.statusText}`;
    let payload: any = null;
    try { payload = await r.json(); if (payload?.error) msg = payload.error; } catch {}
    throw new ApiError(r.status, msg, payload);
  }
  if (r.status === 204) return undefined as any;
  const ct = r.headers.get("Content-Type") || "";
  if (ct.includes("application/json")) return r.json();
  return (await r.text()) as any;
}

export const api = {
  get:    <T = any>(path: string)            => req<T>("GET", path),
  post:   <T = any>(path: string, body?: any) => req<T>("POST", path, body),
  put:    <T = any>(path: string, body?: any) => req<T>("PUT", path, body),
  del:    <T = any>(path: string)            => req<T>("DELETE", path),
};
