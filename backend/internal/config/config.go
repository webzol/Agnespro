// Package config loads runtime configuration from environment variables.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Config struct {
	Addr           string
	DataDir        string
	BaseURL        string        // resolved Agnes API base URL
	APIRoute       string        // international / international-alt / china
	Concurrency    int
	PollInterval   time.Duration
	VideoTimeout   time.Duration
	MasterKeyPath  string        // file holding the master encryption key
	MasterKey      []byte        // 32-byte key, populated after Load
}

var routeURLs = map[string]string{
	"international":     "https://apihub.agnes-ai.com/v1",
	"international-alt": "https://apihub.agnes-ai.cn/v1",
	"china":             "https://api.agnes-ai.cn/v1",
}

// Load reads environment variables, applies defaults, and resolves the
// Agnes API base URL according to the route setting.
func Load() (*Config, error) {
	c := &Config{
		Addr:         getEnv("AGNES_STUDIO_ADDR", ":8080"),
		DataDir:      getEnv("AGNES_STUDIO_DATA_DIR", "./data"),
		APIRoute:     getEnv("AGNES_STUDIO_ROUTE", "international"),
		PollInterval: 5 * time.Second,
		VideoTimeout: 15 * time.Minute,
	}

	if v := os.Getenv("AGNES_STUDIO_VIDEO_POLL_INTERVAL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid AGNES_STUDIO_VIDEO_POLL_INTERVAL: %w", err)
		}
		c.PollInterval = d
	}
	if v := os.Getenv("AGNES_STUDIO_VIDEO_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid AGNES_STUDIO_VIDEO_TIMEOUT: %w", err)
		}
		c.VideoTimeout = d
	}

	if v := os.Getenv("AGNES_STUDIO_MAX_CONCURRENCY"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return nil, fmt.Errorf("invalid AGNES_STUDIO_MAX_CONCURRENCY: %v", v)
		}
		c.Concurrency = n
	} else {
		c.Concurrency = 1
	}

	base, ok := routeURLs[c.APIRoute]
	if !ok {
		return nil, fmt.Errorf("unknown AGNES_STUDIO_ROUTE %q (valid: international, international-alt, china)", c.APIRoute)
	}
	c.BaseURL = base

	if err := os.MkdirAll(c.DataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	abs, err := filepath.Abs(c.DataDir)
	if err != nil {
		return nil, err
	}
	c.DataDir = abs
	c.MasterKeyPath = filepath.Join(c.DataDir, ".masterkey")

	return c, nil
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
