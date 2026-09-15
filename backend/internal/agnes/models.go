// Package agnes — ListModels + ClassifyModel
package agnes

import (
	"context"
	"net/http"
	"strings"
)

// ModelInfo is one entry from the OpenAI-compatible GET /models endpoint.
type ModelInfo struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by,omitempty"`
}

// ModelsListResponse is the OpenAI-compatible wrapper.
type ModelsListResponse struct {
	Object string      `json:"object"`
	Data   []ModelInfo `json:"data"`
}

// ListModels calls GET /models and returns the available models.
// The Agnes gateway returns all model IDs (chat/image/video) together;
// callers should classify by ID prefix.
func (c *Client) ListModels(ctx context.Context) (*ModelsListResponse, error) {
	var out ModelsListResponse
	_, err := c.do(ctx, http.MethodGet, "/models", nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ClassifyModel returns one of "chat" / "image" / "video" based on the
// model ID prefix. Unknown IDs default to "chat".
func ClassifyModel(id string) string {
	lower := strings.ToLower(id)
	switch {
	case strings.Contains(lower, "-image-"):
		return "image"
	case strings.Contains(lower, "-video-"):
		return "video"
	default:
		return "chat"
	}
}
