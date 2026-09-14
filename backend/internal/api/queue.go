package api

import "net/http"

func (s *Server) queueStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"queue": s.Queue.Stats(),
	})
}
