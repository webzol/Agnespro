import { api } from "./client";
import type {
  Job, StyleLibrary, AdminConfig, ModelInfo,
} from "@/types";

export const jobsApi = {
  list:    ()        => api.get<{ jobs: Job[] }>("/api/jobs").then(r => r.jobs),
  get:     (id: string) => api.get<Job>(`/api/jobs/${id}`),
  create:  (body: any) => api.post<Job>("/api/jobs", body),
  del:     (id: string) => api.del(`/api/jobs/${id}`),
  run:     (id: string) => api.post(`/api/jobs/${id}/run`),
  cancel:  (id: string) => api.post(`/api/jobs/${id}/cancel`),
  parse:   (id: string) => api.post(`/api/jobs/${id}/parse`),
};

export const scriptApi = {
  generate: (body: any) => api.post<{ script: string; title: string }>("/api/scripts/generate", body),
};

export const stylesApi = {
  library: () => api.get<StyleLibrary>("/api/styles/library"),
  visual:  () => api.get<{ styles: any[] }>("/api/styles/visual"),
};

export const adminApi = {
  getConfig:     ()        => api.get<AdminConfig>("/admin/api/config"),
  updateConfig:  (body: any) => api.post("/admin/api/config", body),
  testKey:       (api_key?: string) => api.post<{ ok: boolean; message?: string; base_url?: string; route?: string; duration_ms?: number }>("/admin/api/config/test-key", api_key ? { api_key } : {}),
  listModels:    ()        => api.get<{ models: ModelInfo[] }>("/api/settings/models"),
};

export const settingsApi = {
  get:    ()        => api.get<any>("/api/settings"),
  update: (body: any) => api.post("/api/settings", body),
  setKey: (api_key: string) => api.post("/api/settings/key", { api_key }),
  testKey: (api_key: string) => api.post<any>("/api/settings/key/test", { api_key }),
  clearKey: () => api.del("/api/settings/key"),
  setRoute: (route: string) => api.post("/api/settings/route", { route }),
};

export const docsApi = {
  md: () => fetch("/api/docs").then(r => r.text()),
};
