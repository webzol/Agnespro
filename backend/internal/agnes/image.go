package agnes

import (
	"context"
	"net/http"
)

type ImageRequest struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	N       int    `json:"n,omitempty"`
	Size    string `json:"size,omitempty"`
	Quality string `json:"quality,omitempty"`
}

type ImageResponse struct {
	Created int64 `json:"created"`
	Data    []struct {
		URL           string `json:"url"`
		B64JSON       string `json:"b64_json"`
		RevisedPrompt string `json:"revised_prompt"`
	} `json:"data"`
}

func (c *Client) GenerateImage(ctx context.Context, req ImageRequest) (*ImageResponse, error) {
	if req.Model == "" {
		req.Model = "agnes-image-2.1-flash"
	}
	if req.N == 0 {
		req.N = 1
	}
	if req.Size == "" {
		req.Size = "1024x1024"
	}
	var out ImageResponse
	_, err := c.do(ctx, http.MethodPost, "/images/generations", req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
