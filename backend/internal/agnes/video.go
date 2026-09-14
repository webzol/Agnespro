package agnes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type VideoRequest struct {
	Model          string         `json:"model"`
	Prompt         string         `json:"prompt"`
	Seconds        string         `json:"seconds,omitempty"`
	Size           string         `json:"size,omitempty"`
	InputReference []string       `json:"input_reference,omitempty"`
	Extra          map[string]any `json:"-"`
}

type VideoCreateResponse struct {
	ID        string `json:"id"`
	VideoID   string `json:"video_id"`
	TaskID    string `json:"task_id"`
	Object    string `json:"object"`
	Model     string `json:"model"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	CreatedAt int64  `json:"created_at"`
	Seconds   string `json:"seconds"`
	Size      string `json:"size"`
}

type VideoPollResponse struct {
	ID             string `json:"id"`
	VideoID        string `json:"video_id"`
	TaskID         string `json:"task_id"`
	Object         string `json:"object"`
	Model          string `json:"model"`
	Status         string `json:"status"`
	InternalStatus string `json:"internal_status"`
	Progress       int    `json:"progress"`
	URL            string `json:"url"`
	Error          *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
	Seconds     string `json:"seconds"`
	Size        string `json:"size"`
	CreatedAt   int64  `json:"created_at"`
	CompletedAt *int64 `json:"completed_at"`
	ExpiresAt   *int64 `json:"expires_at"`
}

func (c *Client) CreateVideo(ctx context.Context, req VideoRequest) (*VideoCreateResponse, error) {
	if req.Model == "" {
		req.Model = "agnes-video-v2.0"
	}
	if req.Seconds == "" {
		req.Seconds = "5"
	}
	payload := map[string]any{
		"model":   req.Model,
		"prompt":  req.Prompt,
		"seconds": req.Seconds,
	}
	if req.Size != "" {
		payload["size"] = req.Size
	}
	if len(req.InputReference) > 0 {
		payload["input_reference"] = req.InputReference
	}
	for k, v := range req.Extra {
		payload[k] = v
	}
	var out VideoCreateResponse
	_, err := c.do(ctx, http.MethodPost, "/videos", payload, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) PollVideo(ctx context.Context, videoID string) (*VideoPollResponse, error) {
	if videoID == "" {
		return nil, errors.New("agnes: empty video id")
	}
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}
	host := u.Scheme + "://" + u.Host
	apiPath := strings.TrimSuffix(u.Path, "/v1")
	pollURL := fmt.Sprintf("%s%s/agnesapi?video_id=%s", host, apiPath, url.QueryEscape(videoID))
	var out VideoPollResponse
	_, err = c.doAbsolute(ctx, http.MethodGet, pollURL, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) WaitForVideo(ctx context.Context, videoID string, interval time.Duration) (*VideoPollResponse, error) {
	if interval <= 0 {
		interval = 5 * time.Second
	}
	for {
		resp, err := c.PollVideo(ctx, videoID)
		if err != nil {
			return nil, err
		}
		switch resp.Status {
		case "completed", "succeeded", "success":
			return resp, nil
		case "failed", "cancelled", "canceled":
			if resp.Error != nil {
				return resp, fmt.Errorf("agnes video %s: %s", resp.Status, resp.Error.Message)
			}
			return resp, fmt.Errorf("agnes video %s", resp.Status)
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(interval):
		}
	}
}
