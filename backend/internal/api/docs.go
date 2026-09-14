package api

import (
	"net/http"
	"os"
	"path/filepath"
)

func (s *Server) serveDocs(w http.ResponseWriter, r *http.Request) {
	// The frontend fetches /api/docs for the embedded usage guide.
	// Serve docs/USER_GUIDE.md from the repo root if present.
	candidates := []string{
		filepath.Join("docs", "USER_GUIDE.md"),
		filepath.Join("..", "docs", "USER_GUIDE.md"),
		filepath.Join("..", "..", "docs", "USER_GUIDE.md"),
	}
	for _, c := range candidates {
		if b, err := os.ReadFile(c); err == nil {
			w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
			w.Header().Set("Cache-Control", "no-cache")
			_, _ = w.Write(b)
			return
		}
	}
	writeError(w, http.StatusNotFound, "documentation not found")
}
