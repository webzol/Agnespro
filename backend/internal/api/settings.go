package api

import (
	"context"
	"net/http"
	"time"

	"agnesai/studio/internal/agnes"
	"agnesai/studio/internal/types"
)

type setKeyReq struct {
	APIKey string `json:"api_key"`
}

type testKeyResp struct {
	OK       bool   `json:"ok"`
	Message  string `json:"message,omitempty"`
	BaseURL  string `json:"base_url,omitempty"`
	Route    string `json:"route,omitempty"`
	Duration int64  `json:"duration_ms,omitempty"`
}

func (s *Server) getSettings(w http.ResponseWriter, r *http.Request) {
	set, err := s.Store.LoadSettings()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Mark whether key is set without revealing it
	enc, err := s.Store.GetAPIKey()
	if err == nil && enc != "" {
		set.APIKeySet = true
	}
	// Always report the configured base URL/route for the UI
	set.BaseURL = s.Cfg.BaseURL
	set.APIRoute = s.Cfg.APIRoute
	set.Concurrency = s.Cfg.Concurrency
	set.PollIntervalS = int(s.Cfg.PollInterval / time.Second)
	writeJSON(w, http.StatusOK, set)
}

func (s *Server) setAPIKey(w http.ResponseWriter, r *http.Request) {
	var req setKeyReq
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}
	if req.APIKey == "" {
		writeError(w, http.StatusBadRequest, "api_key is required")
		return
	}
	enc, err := s.Crypt.Encrypt([]byte(req.APIKey))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "encrypt: "+err.Error())
		return
	}
	if err := s.Store.SetAPIKey(enc); err != nil {
		writeError(w, http.StatusInternalServerError, "store key: "+err.Error())
		return
	}
	// also flush the temporary key from this turn's memory; settings
	// will reload from disk on next request.
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "saved": true})
}

func (s *Server) clearAPIKey(w http.ResponseWriter, r *http.Request) {
	// Save an empty file to clear.
	if err := s.Store.SetAPIKey(""); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "cleared": true})
}

func (s *Server) testAPIKey(w http.ResponseWriter, r *http.Request) {
	// Use the just-submitted key if provided; otherwise read encrypted key
	apiKey := ""
	var body struct {
		APIKey string `json:"api_key"`
	}
	_ = decodeJSON(r, &body)
	if body.APIKey != "" {
		apiKey = body.APIKey
	} else {
		enc, err := s.Store.GetAPIKey()
		if err != nil || enc == "" {
			writeJSON(w, http.StatusOK, testKeyResp{OK: false, Message: "no API key configured"})
			return
		}
		raw, err := s.Crypt.Decrypt(enc)
		if err != nil {
			writeJSON(w, http.StatusOK, testKeyResp{OK: false, Message: "decrypt failed: " + err.Error()})
			return
		}
		apiKey = string(raw)
	}
	c := agnes.New(s.Cfg.BaseURL, apiKey)
	start := time.Now()
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	if err := c.Ping(ctx); err != nil {
		writeJSON(w, http.StatusOK, testKeyResp{OK: false, Message: err.Error(), BaseURL: s.Cfg.BaseURL, Route: s.Cfg.APIRoute, Duration: time.Since(start).Milliseconds()})
		return
	}
	writeJSON(w, http.StatusOK, testKeyResp{OK: true, Message: "connected", BaseURL: s.Cfg.BaseURL, Route: s.Cfg.APIRoute, Duration: time.Since(start).Milliseconds()})
}

type setRouteReq struct {
	Route string `json:"route"`
}

func (s *Server) setRoute(w http.ResponseWriter, r *http.Request) {
	var req setRouteReq
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}
	valid := map[string]bool{"international": true, "international-alt": true, "china": true}
	if !valid[req.Route] {
		writeError(w, http.StatusBadRequest, "invalid route (use international, international-alt, or china)")
		return
	}
	set, err := s.Store.LoadSettings()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	set.APIRoute = req.Route
	// update base URL mapping
	switch req.Route {
	case "international":
		set.BaseURL = "https://apihub.agnes-ai.com/v1"
	case "international-alt":
		set.BaseURL = "https://apihub.agnes-ai.cn/v1"
	case "china":
		set.BaseURL = "https://api.agnes-ai.cn/v1"
	}
	if err := s.Store.SaveSettings(set); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Note: changing route at runtime requires a server restart to
	// update the active base URL. The UI will warn the user.
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "route": req.Route, "base_url": set.BaseURL, "restart_required": true})
}

type updateSettingsReq struct {
	Concurrency *int `json:"concurrency,omitempty"`
	Style       string `json:"style,omitempty"`
}

func (s *Server) updateSettings(w http.ResponseWriter, r *http.Request) {
	var req updateSettingsReq
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}
	set, err := s.Store.LoadSettings()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if req.Concurrency != nil && *req.Concurrency > 0 {
		set.Concurrency = *req.Concurrency
	}
	if err := s.Store.SaveSettings(set); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, set)
}

var _ = types.StatusPending
