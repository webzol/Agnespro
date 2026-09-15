package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"agnesai/studio/internal/config"
	"agnesai/studio/internal/cryptox"
	"agnesai/studio/internal/pipeline"
	"agnesai/studio/internal/queue"
	"agnesai/studio/internal/store"
)

type Server struct {
	Cfg      *config.Config
	Store    *store.Store
	Crypt    *cryptox.Cipher
	Pipeline *pipeline.Pipeline
	Queue    *queue.Queue
	WebDir   string
}

func NewServer(cfg *config.Config, st *store.Store, cr *cryptox.Cipher, pl *pipeline.Pipeline, q *queue.Queue, webDir string) *Server {
	return &Server{Cfg: cfg, Store: st, Crypt: cr, Pipeline: pl, Queue: q, WebDir: webDir}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	// API
	mux.HandleFunc("GET /api/health", s.health)
	mux.HandleFunc("GET /api/queue", s.queueStatus)
	mux.HandleFunc("GET /api/settings", s.getSettings)
	mux.HandleFunc("POST /api/settings", s.updateSettings)
	mux.HandleFunc("GET /api/settings/models", s.listModels)
	mux.HandleFunc("POST /api/settings/key", s.setAPIKey)
	mux.HandleFunc("POST /api/settings/key/test", s.testAPIKey)
	mux.HandleFunc("DELETE /api/settings/key", s.clearAPIKey)
	mux.HandleFunc("POST /api/settings/route", s.setRoute)

	mux.HandleFunc("GET /api/jobs", s.listJobs)
	mux.HandleFunc("POST /api/jobs", s.createJob)
	mux.HandleFunc("GET /api/jobs/{id}", s.getJob)
	mux.HandleFunc("DELETE /api/jobs/{id}", s.deleteJob)
	mux.HandleFunc("POST /api/jobs/{id}/run", s.runJob)
	mux.HandleFunc("POST /api/jobs/{id}/cancel", s.cancelJob)
	mux.HandleFunc("POST /api/jobs/{id}/episodes/{n}/run", s.runEpisode)
	mux.HandleFunc("POST /api/jobs/{id}/parse", s.parseJob)

	mux.HandleFunc("GET /api/assets/{name}", s.serveAsset)
	mux.HandleFunc("GET /api/jobs/{id}/subtitles", s.downloadSubtitles)
	mux.HandleFunc("GET /api/jobs/{id}/download", s.downloadJobAssets)
	mux.HandleFunc("GET /api/styles/visual", s.listVisualStyles)
	mux.HandleFunc("GET /api/styles/library", s.listStyleLibrary)
	mux.HandleFunc("GET /api/styles/subtitle", s.listSubtitleStyles)
	mux.HandleFunc("POST /api/scripts/generate", s.generateScript)
	mux.HandleFunc("GET /api/docs", s.serveDocs)

	// Static frontend
	if s.WebDir != "" {
		mux.Handle("/", s.staticHandler())
	}

	return Chain(mux, Logger, CORS)
}

func (s *Server) staticHandler() http.Handler {
	fs := http.FileServer(http.Dir(s.WebDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the request looks like an API call, fall through (mux will 404)
		p := r.URL.Path
		if strings.HasPrefix(p, "/api/") {
			http.NotFound(w, r)
			return
		}
		// SPA fallback: serve index.html for paths that don't have an extension
		if filepath.Ext(p) == "" && p != "/" {
			full := filepath.Join(s.WebDir, p)
			if _, err := os.Stat(full); os.IsNotExist(err) {
				http.ServeFile(w, r, filepath.Join(s.WebDir, "index.html"))
				return
			}
		}
		// Normal file server
		fs.ServeHTTP(w, r)
	})
}

func newID(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func backgroundContext() context.Context {
	return context.Background()
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":     true,
		"now":    time.Now(),
		"addr":   s.Cfg.Addr,
		"route":  s.Cfg.APIRoute,
		"base":   s.Cfg.BaseURL,
		"app":    "agnesai-studio",
		"build":  "1.0.0",
	})
}
