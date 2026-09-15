package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"agnesai/studio/internal/pipeline"
	"agnesai/studio/internal/queue"
	"agnesai/studio/internal/store"
	"agnesai/studio/internal/types"
)

type createJobReq struct {
	Title           string `json:"title"`
	Script          string `json:"script"`
	Style           string `json:"style"`
	VisualStyle     string `json:"visual_style"`
	VisualStyleName string `json:"visual_style_name"`
}

type localEpisode struct {
	Index int    `json:"index"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func shortUUID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *Server) createJob(w http.ResponseWriter, r *http.Request) {
	var req createJobReq
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}
	req.Script = strings.TrimSpace(req.Script)
	if req.Script == "" {
		writeError(w, http.StatusBadRequest, "script is required")
		return
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = "Untitled"
	}
	now := time.Now()
	parsed := pipeline.ParseEpisodes(req.Script, title)
	// If a visual_style id was provided, use its Chinese description as the
	// style hint that the pipeline will use for prompts.
	if req.VisualStyle != "" {
		if name := req.VisualStyleName; name != "" {
			req.Style = name
		} else {
			req.Style = "电影感"
		}
	}
	j := &types.Job{
		ID:         "job_" + shortUUID(),
		Title:      title,
		Script:     req.Script,
		Style:      req.Style,
		Status:     types.StatusPending,
		Progress:   0,
		Episodes:   make([]types.Episode, len(parsed)),
		NumEpisodes: len(parsed),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	for i, ep := range parsed {
		j.Episodes[i] = types.Episode{
			ID:        "ep_" + shortUUID(),
			Index:     ep.Index,
			Title:     ep.Title,
			Body:      ep.Body,
			State:     types.EpisodePending,
			Progress:  0,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	if err := s.Store.CreateJob(j); err != nil {
		writeError(w, http.StatusInternalServerError, "save job: "+err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, j)
}

func (s *Server) listJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := s.Store.ListJobs()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"jobs": jobs})
}

func (s *Server) getJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	j, err := s.Store.GetJob(id)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, j)
}

func (s *Server) deleteJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.Store.DeleteJob(id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "job not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) parseJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	j, err := s.Store.GetJob(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	parsed := pipeline.ParseEpisodes(j.Script, j.Title)
	out := make([]localEpisode, len(parsed))
	for i, p := range parsed {
		out[i] = localEpisode{Index: p.Index, Title: p.Title, Body: p.Body}
	}
	writeJSON(w, http.StatusOK, map[string]any{"episodes": out, "count": len(out)})
}

func (s *Server) runJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := s.Store.GetJob(id); err != nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	s.enqueueRun(id, 0)
	writeJSON(w, http.StatusAccepted, map[string]any{"ok": true, "queued": true, "job_id": id})
}

func (s *Server) runEpisode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	nStr := r.PathValue("n")
	n, err := strconv.Atoi(nStr)
	if err != nil || n < 1 {
		writeError(w, http.StatusBadRequest, "invalid episode number")
		return
	}
	if _, err := s.Store.GetJob(id); err != nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	s.enqueueRun(id, n)
	writeJSON(w, http.StatusAccepted, map[string]any{"ok": true, "queued": true, "job_id": id, "episode": n})
}

func (s *Server) cancelJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.Queue.Cancel(id)
	if j, err := s.Store.GetJob(id); err == nil {
		j.Status = types.StatusCancelled
		_ = s.Store.UpdateJob(j)
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) enqueueRun(jobID string, onlyEpisode int) {
	pl := s.Pipeline
	q := s.Queue
	storeRef := s.Store
	q.Enqueue(&queue.Job{
		ID: jobID,
		Run: func(ctx context.Context) error {
			err := pl.Run(ctx, jobID, onlyEpisode)
			if err != nil {
				if j, gerr := storeRef.GetJob(jobID); gerr == nil {
					if j.Status != types.StatusDone {
						j.Status = types.StatusFailed
						j.Error = err.Error()
						_ = storeRef.UpdateJob(j)
					}
				}
			}
			return err
		},
	})
}
