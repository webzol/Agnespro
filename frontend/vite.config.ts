import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import path from "node:path";

// Build output goes into the Go backend's static directory so the existing
// server can serve the SPA directly. During `vite dev` we proxy /api to
// the Go backend on :8080.
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: { "@": path.resolve(__dirname, "src") },
  },
  build: {
    outDir: path.resolve(__dirname, "../backend/web"),
    emptyOutDir: false,
    target: "es2020",
    sourcemap: false,
    rollupOptions: {
      // Keep absolute /assets/* paths as runtime strings (served by backend directly)
      external: [/^\/assets\//],
    },
  },
  server: {
    port: 5173,
    proxy: {
      "/api":   { target: "http://localhost:8080", changeOrigin: true },
      "/admin": { target: "http://localhost:8080", changeOrigin: true },
      "/assets/img/library": { target: "http://localhost:8080", changeOrigin: true },
    },
  },
});
