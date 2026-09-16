import { createRouter, createWebHashHistory } from "vue-router";
import HomeView from "@/views/HomeView.vue";
import NewJobView from "@/views/NewJobView.vue";
import JobsView from "@/views/JobsView.vue";
import JobDetailView from "@/views/JobDetailView.vue";
import AdminView from "@/views/AdminView.vue";
import DocsView from "@/views/DocsView.vue";
import ScriptTemplatesView from "@/views/ScriptTemplatesView.vue";

// Eager imports: all views baked into the main bundle so route switches are
// instant (no per-navigation chunk fetch). Total payload is ~120 KB which is
// acceptable for an internal SaaS tool.
const routes = [
  { path: "/",          name: "home",      component: HomeView },
  { path: "/new",       name: "new",       component: NewJobView },
  { path: "/jobs",      name: "jobs",      component: JobsView },
  { path: "/detail/:id", name: "detail",   component: JobDetailView },
  { path: "/admin",     name: "admin",     component: AdminView },
  { path: "/docs",      name: "docs",      component: DocsView },
  { path: "/templates", name: "templates", component: ScriptTemplatesView },
  { path: "/:catchAll(.*)", redirect: "/" },
];

export const router = createRouter({
  history: createWebHashHistory(),
  routes,
});
