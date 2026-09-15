// Command server is the Agnes AI Studio backend.
//
// It exposes a REST API and serves the bundled frontend (web/).
// All persistence is local (JSON files under the data directory) and
// the user's Agnes API key is encrypted at rest with AES-256-GCM.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"agnesai/studio/internal/api"
	"agnesai/studio/internal/config"
	"agnesai/studio/internal/cryptox"
	"agnesai/studio/internal/pipeline"
	"agnesai/studio/internal/queue"
	"agnesai/studio/internal/store"
)

func main() {
	var (
		showVersion = flag.Bool("version", false, "print version and exit")
		addr        = flag.String("addr", "", "override listen address (e.g. :8080)")
		dataDir     = flag.String("data", "", "override data directory")
		route       = flag.String("route", "", "override API route: international | international-alt | china")
		webDir      = flag.String("web", "", "override web assets directory")
		masterKey   = flag.String("master-key", "", "override master encryption key (base64, 32 bytes)")
	)
	flag.Parse()

	if *showVersion {
		fmt.Println("agnesai-studio 1.0.0")
		return
	}

	loadDotEnv(".env")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if *addr != "" {
		cfg.Addr = *addr
	}
	if *dataDir != "" {
		cfg.DataDir = *dataDir
		_ = os.MkdirAll(cfg.DataDir, 0o755)
	}
	if *route != "" {
		cfg.APIRoute = *route
		switch *route {
		case "international":
			cfg.BaseURL = "https://apihub.agnes-ai.com/v1"
		case "international-alt":
			cfg.BaseURL = "https://apihub.agnes-ai.cn/v1"
		case "china":
			cfg.BaseURL = "https://api.agnes-ai.cn/v1"
		}
	}

	webDirValue := *webDir
	if webDirValue != "" {
		abs, _ := filepath.Abs(webDirValue)
		webDirValue = abs
	} else {
		candidates := []string{"web", filepath.Join("..", "..", "web"), filepath.Join("..", "web")}
		for _, c := range candidates {
			if st, err := os.Stat(c); err == nil && st.IsDir() {
				abs, _ := filepath.Abs(c)
				webDirValue = abs
				break
			}
		}
	}

	var cipher *cryptox.Cipher
	if *masterKey != "" {
		cipher, err = cryptox.NewFromBase64(*masterKey)
		if err != nil {
			log.Fatalf("master key: %v", err)
		}
	} else if env := os.Getenv("AGNES_STUDIO_MASTER_KEY"); env != "" {
		cipher, _ = cryptox.New(cryptox.Derive(env))
	} else {
		cipher, err = cryptox.NewFromFile(cfg.MasterKeyPath)
		if err != nil {
			log.Fatalf("master key file: %v", err)
		}
	}

	st, err := store.New(cfg.DataDir)
	if err != nil {
		log.Fatalf("store: %v", err)
	}
	// Load persisted settings (incl. default model selections) before
	// constructing the pipeline so they become the initial defaults.
	savedSet, _ := st.LoadSettings()
	pl := pipeline.New(st, cipher, pipeline.Config{
		BaseURL:      cfg.BaseURL,
		PollInterval: cfg.PollInterval,
		VideoTimeout: cfg.VideoTimeout,
		VideoSeconds: "5",
		ImageSize:    "1024x1024",
		ScriptModel:  ifNonEmpty(savedSet.ScriptModel, "agnes-2.5-flash"),
		ImageModel:   ifNonEmpty(savedSet.ImageModel, "agnes-image-2.1-flash"),
		VideoModel:   ifNonEmpty(savedSet.VideoModel, "agnes-video-v2.0"),
	})
	q := queue.New(cfg.Concurrency)
	srv := api.NewServer(cfg, st, cipher, pl, q, webDirValue)
	q.OnJobUpdate = func(jobID, state string, err error) {
		log.Printf("queue: job=%s state=%s err=%v", jobID, state, err)
	}

	httpSrv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 30 * time.Second,
	}

	idle := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Printf("shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = httpSrv.Shutdown(ctx)
		q.Stop()
		close(idle)
	}()

	if webDirValue != "" {
		log.Printf("serving frontend from %s", webDirValue)
	}
	log.Printf("Agnes AI Studio listening on http://localhost%s", cfg.Addr)
	log.Printf("Agnes route=%s base=%s", cfg.APIRoute, cfg.BaseURL)
	log.Printf("Data dir: %s", cfg.DataDir)

	if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %v", err)
	}
	<-idle
	log.Printf("bye")
}

func loadDotEnv(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range rangeLines(string(b)) {
		line = trim(line)
		if line == "" || line[0] == '#' {
			continue
		}
		eq := indexByte(line, '=')
		if eq < 0 {
			continue
		}
		k := trim(line[:eq])
		v := trim(line[eq+1:])
		if len(v) >= 2 {
			if v[0] == '"' && v[len(v)-1] == '"' {
				v = v[1 : len(v)-1]
			} else if v[0] == '\'' && v[len(v)-1] == '\'' {
				v = v[1 : len(v)-1]
			}
		}
		if _, exists := os.LookupEnv(k); !exists {
			_ = os.Setenv(k, v)
		}
	}
}

func rangeLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

func trim(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\r') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}


// ifNonEmpty returns fallback when s is empty.
func ifNonEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}
