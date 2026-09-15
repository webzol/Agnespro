import { createRouter, createWebHashHistory } from "vue-router";
import HomeView from "@/views/HomeView.vue";

const routes = [
  { path: "/",         name: "home",     component: () => import("@/views/HomeView.vue") },
  { path: "/new",      name: "new",      component: () => import("@/views/NewJobView.vue") },
  { path: "/jobs",     name: "jobs",     component: () => import("@/views/JobsView.vue") },
  { path: "/detail/:id", name: "detail",   component: () => import("@/views/JobDetailView.vue") },
  { path: "/settings", name: "settings", component: () => import("@/views/SettingsView.vue") },
  { path: "/admin",    name: "admin",    component: () => import("@/views/AdminView.vue") },
  { path: "/docs",     name: "docs",     component: () => import("@/views/DocsView.vue") },
  { path: "/:catchAll(.*)", redirect: "/" },
];

export const router = createRouter({
  history: createWebHashHistory(),
  routes,
});
