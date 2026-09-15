// admin.go - 后台管理 API（独立的 key / 路由 / 默认模型配置入口）
//
// 设计:集中返回所有可配置项,前端可以一次性加载并展示;
// 同样支持一次性 POST 修改任意子集(key / route / 模型)。
//
// 与 /api/settings/* 的关系:
//   - /api/settings/* 是给前端用户 Settings 页面用的(只读 view + 部分更新)
//   - /admin/api/config 是给"后台"专门用的一站式 API,聚合了 key + route + models + 列表

package api

import (
	"context"
	"net/http"
	"time"

	"agnesai/studio/internal/agnes"
)

// adminConfigResp is the aggregated response for /admin/api/config
type adminConfigResp struct {
	APIKeySet      bool        `json:"api_key_set"`        // true if key is configured (never the actual key)
	APIRoute       string      `json:"api_route"`          // international / international-alt / china
	BaseURL        string      `json:"base_url"`           // https://apihub.agnes-ai.com/v1 etc
	ScriptModel    string      `json:"script_model"`       // default chat model
	ImageModel     string      `json:"image_model"`        // default image model
	VideoModel     string      `json:"video_model"`        // default video model
	Concurrency    int         `json:"concurrency"`
	AvailableModels []modelLite `json:"available_models"`  // populated when API key is set
	ModelsError    string      `json:"models_error,omitempty"` // populated when model fetch failed
	FetchedAt      time.Time   `json:"fetched_at"`
}

type modelLite struct {
	ID      string `json:"id"`
	Type    string `json:"type"`              // chat / image / video
	OwnedBy string `json:"owned_by,omitempty"`
}

// GET /admin/api/config
//
// Returns the full admin view in one shot: key state, route, default models,
// and (when key is set) the available models from the Agnes gateway.
func (s *Server) adminGetConfig(w http.ResponseWriter, r *http.Request) {
	set, err := s.Store.LoadSettings()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	enc, err := s.Store.GetAPIKey()
	keySet := err == nil && enc != ""

	resp := adminConfigResp{
		APIKeySet:   keySet,
		APIRoute:    s.Cfg.APIRoute,
		BaseURL:     s.Cfg.BaseURL,
		ScriptModel: coalesceStr(set.ScriptModel, "agnes-2.5-flash"),
		ImageModel:  coalesceStr(set.ImageModel, "agnes-image-2.1-flash"),
		VideoModel:  coalesceStr(set.VideoModel, "agnes-video-v2.0"),
		Concurrency: s.Cfg.Concurrency,
		FetchedAt:   time.Now(),
	}

	if keySet {
		raw, err := s.Crypt.Decrypt(enc)
		if err == nil {
			client := agnes.New(s.Cfg.BaseURL, string(raw))
			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
			defer cancel()
			if ml, err := client.ListModels(ctx); err != nil {
				resp.ModelsError = err.Error()
			} else {
				for _, m := range ml.Data {
					resp.AvailableModels = append(resp.AvailableModels, modelLite{
						ID:      m.ID,
						Type:    agnes.ClassifyModel(m.ID),
						OwnedBy: m.OwnedBy,
					})
				}
				if resp.AvailableModels == nil {
					resp.AvailableModels = []modelLite{}
				}
			}
		} else {
			resp.ModelsError = "decrypt: " + err.Error()
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

func coalesceStr(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// adminConfigReq is the request body for POST /admin/api/config.
// All fields are optional; only non-empty fields are applied.
type adminConfigReq struct {
	APIKey      string `json:"api_key,omitempty"`       // new key (encrypt + store); empty = no change
	ClearAPIKey bool   `json:"clear_api_key,omitempty"` // true = clear stored key
	Route       string `json:"route,omitempty"`          // international / international-alt / china
	ScriptModel string `json:"script_model,omitempty"`  // default chat model
	ImageModel  string `json:"image_model,omitempty"`   // default image model
	VideoModel  string `json:"video_model,omitempty"`   // default video model
	Concurrency *int   `json:"concurrency,omitempty"`   // pipeline concurrency (1-N)
}

// POST /admin/api/config
//
// Updates any subset of: api_key, route, default models, concurrency.
// Route changes require a server restart (active base URL stays in cfg until restart).
// Model changes take effect immediately for new jobs (Pipeline.Cfg is mutated in place).
func (s *Server) adminUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var req adminConfigReq
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}

	set, err := s.Store.LoadSettings()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "load settings: "+err.Error())
		return
	}

	routeChanged := false
	modelsChanged := false

	// 1. API key handling
	if req.ClearAPIKey {
		if err := s.Store.SetAPIKey(""); err != nil {
			writeError(w, http.StatusInternalServerError, "clear key: "+err.Error())
			return
		}
	} else if req.APIKey != "" {
		enc, err := s.Crypt.Encrypt([]byte(req.APIKey))
		if err != nil {
			writeError(w, http.StatusInternalServerError, "encrypt: "+err.Error())
			return
		}
		if err := s.Store.SetAPIKey(enc); err != nil {
			writeError(w, http.StatusInternalServerError, "store key: "+err.Error())
			return
		}
	}

	// 2. Route
	if req.Route != "" {
		valid := map[string]bool{"international": true, "international-alt": true, "china": true}
		if !valid[req.Route] {
			writeError(w, http.StatusBadRequest, "invalid route (use international, international-alt, or china)")
			return
		}
		if req.Route != s.Cfg.APIRoute {
			routeChanged = true
		}
		set.APIRoute = req.Route
		switch req.Route {
		case "international":
			set.BaseURL = "https://apihub.agnes-ai.com/v1"
		case "international-alt":
			set.BaseURL = "https://apihub.agnes-ai.cn/v1"
		case "china":
			set.BaseURL = "https://api.agnes-ai.cn/v1"
		}
	}

	// 3. Default models
	if req.ScriptModel != "" {
		set.ScriptModel = req.ScriptModel
		if s.Pipeline != nil {
			s.Pipeline.Cfg.ScriptModel = req.ScriptModel
		}
		modelsChanged = true
	}
	if req.ImageModel != "" {
		set.ImageModel = req.ImageModel
		if s.Pipeline != nil {
			s.Pipeline.Cfg.ImageModel = req.ImageModel
		}
		modelsChanged = true
	}
	if req.VideoModel != "" {
		set.VideoModel = req.VideoModel
		if s.Pipeline != nil {
			s.Pipeline.Cfg.VideoModel = req.VideoModel
		}
		modelsChanged = true
	}

	// 4. Concurrency
	if req.Concurrency != nil && *req.Concurrency > 0 {
		set.Concurrency = *req.Concurrency
	}

	if err := s.Store.SaveSettings(set); err != nil {
		writeError(w, http.StatusInternalServerError, "save settings: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":               true,
		"restart_required": routeChanged, // route change needs restart (in-memory cfg)
		"models_updated":   modelsChanged,
	})
}

// POST /admin/api/config/test-key
//
// Tests the stored (or provided) API key against the configured base URL.
// Useful for verifying a new key before saving it.
func (s *Server) adminTestKey(w http.ResponseWriter, r *http.Request) {
	var body struct {
		APIKey string `json:"api_key"`
	}
	_ = decodeJSON(r, &body)

	apiKey := ""
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
		writeJSON(w, http.StatusOK, testKeyResp{
			OK:       false,
			Message:  err.Error(),
			BaseURL:  s.Cfg.BaseURL,
			Route:    s.Cfg.APIRoute,
			Duration: time.Since(start).Milliseconds(),
		})
		return
	}
	writeJSON(w, http.StatusOK, testKeyResp{
		OK:       true,
		Message:  "connected",
		BaseURL:  s.Cfg.BaseURL,
		Route:    s.Cfg.APIRoute,
		Duration: time.Since(start).Milliseconds(),
	})
}
