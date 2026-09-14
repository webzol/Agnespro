package api

import (
	"net/http"
	"path/filepath"
	"strings"
)

func (s *Server) serveAsset(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" || strings.ContainsAny(name, "/\\") {
		writeError(w, http.StatusBadRequest, "invalid asset name")
		return
	}
	// Defence in depth: only serve from assets dir; path traversal blocked
	clean := filepath.Clean(name)
	if clean != name || strings.HasPrefix(clean, ".") {
		writeError(w, http.StatusBadRequest, "invalid asset name")
		return
	}
	f, err := s.Store.OpenAsset(clean)
	if err != nil {
		writeError(w, http.StatusNotFound, "asset not found")
		return
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Try to set a sensible content type from the extension
	ext := strings.ToLower(filepath.Ext(clean))
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".webp":
		w.Header().Set("Content-Type", "image/webp")
	case ".mp4":
		w.Header().Set("Content-Type", "video/mp4")
	case ".bin":
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeContent(w, r, clean, stat.ModTime(), f)
}
