// Package store persists jobs, episodes, scenes, characters, props,
// assets metadata, and settings to JSON files under the data directory.
//
// The store is safe for concurrent use within a single process. Each
// record is persisted atomically by writing a temp file and renaming.
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"agnesai/studio/internal/types"
)

type Store struct {
	mu       sync.RWMutex
	dataDir  string
	assetsDir string
}

func New(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, err
	}
	assets := filepath.Join(dataDir, "assets")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		return nil, err
	}
	return &Store{dataDir: dataDir, assetsDir: assets}, nil
}

func (s *Store) DataDir() string      { return s.dataDir }
func (s *Store) AssetsDir() string    { return s.assetsDir }

// --- Jobs ---

func (s *Store) CreateJob(j *types.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	jobs, err := s.readJobs()
	if err != nil {
		return err
	}
	jobs[j.ID] = j
	return s.writeJobs(jobs)
}

func (s *Store) GetJob(id string) (*types.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	jobs, err := s.readJobs()
	if err != nil {
		return nil, err
	}
	j, ok := jobs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return j, nil
}

func (s *Store) UpdateJob(j *types.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	jobs, err := s.readJobs()
	if err != nil {
		return err
	}
	if _, ok := jobs[j.ID]; !ok {
		return ErrNotFound
	}
	j.UpdatedAt = time.Now()
	jobs[j.ID] = j
	return s.writeJobs(jobs)
}

func (s *Store) ListJobs() ([]*types.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	jobs, err := s.readJobs()
	if err != nil {
		return nil, err
	}
	out := make([]*types.Job, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, j)
	}
	sort.Slice(out, func(i, k int) bool { return out[i].CreatedAt.After(out[k].CreatedAt) })
	return out, nil
}

func (s *Store) DeleteJob(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	jobs, err := s.readJobs()
	if err != nil {
		return err
	}
	j, ok := jobs[id]
	if !ok {
		return ErrNotFound
	}
	// remove asset files referenced by this job
	for _, fn := range collectAssetFilenames(j) {
		_ = os.Remove(filepath.Join(s.assetsDir, fn))
	}
	delete(jobs, id)
	return s.writeJobs(jobs)
}

// --- Settings ---

func (s *Store) LoadSettings() (*types.Settings, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.readSettings()
}

func (s *Store) SaveSettings(set *types.Settings) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	set.UpdatedAt = time.Now()
	return s.writeSettings(set)
}

// --- Encrypted API key ---

func (s *Store) GetAPIKey() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, err := os.ReadFile(s.apiKeyPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	return string(b), nil
}

func (s *Store) SetAPIKey(enc string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return os.WriteFile(s.apiKeyPath(), []byte(enc), 0o600)
}

// --- Asset files ---

// SaveAsset writes src (e.g. downloaded image bytes) to the assets dir
// and returns the metadata. The caller is responsible for storing
// a.record in a job via UpdateJob if it wants to track the asset on
// the job; otherwise the file may become orphaned.
func (s *Store) SaveAsset(jobID, kind, refType, refID, filename, mime string, src io.Reader) (*types.Asset, error) {
	full := filepath.Join(s.assetsDir, filename)
	f, err := os.Create(full)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	n, err := io.Copy(f, src)
	if err != nil {
		return nil, err
	}
	return &types.Asset{
		ID:        filename,
		JobID:     jobID,
		Kind:      kind,
		RefType:   refType,
		RefID:     refID,
		Filename:  filename,
		MimeType:  mime,
		SizeBytes: n,
		CreatedAt: time.Now(),
	}, nil
}

func (s *Store) OpenAsset(filename string) (*os.File, error) {
	return os.Open(filepath.Join(s.assetsDir, filename))
}

// --- internal ---

var ErrNotFound = errors.New("not found")

func (s *Store) apiKeyPath() string     { return filepath.Join(s.dataDir, ".apikey") }
func (s *Store) jobsPath() string       { return filepath.Join(s.dataDir, "jobs.json") }
func (s *Store) settingsPath() string   { return filepath.Join(s.dataDir, "settings.json") }

func (s *Store) readJobs() (map[string]*types.Job, error) {
	return readJSONMap[*types.Job](s.jobsPath())
}

func (s *Store) writeJobs(m map[string]*types.Job) error {
	return writeJSONAtomic(s.jobsPath(), m)
}

func (s *Store) readSettings() (*types.Settings, error) {
	b, err := os.ReadFile(s.settingsPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &types.Settings{
				APIRoute:       "international",
				BaseURL:        "https://apihub.agnes-ai.com/v1",
				Concurrency:    1,
				PollIntervalS:  5,
				UpdatedAt:      time.Now(),
			}, nil
		}
		return nil, err
	}
	var set types.Settings
	if err := json.Unmarshal(b, &set); err != nil {
		return nil, err
	}
	return &set, nil
}

func (s *Store) writeSettings(set *types.Settings) error {
	return writeJSONAtomic(s.settingsPath(), set)
}

func readJSONMap[V any](path string) (map[string]V, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]V{}, nil
		}
		return nil, err
	}
	if len(b) == 0 {
		return map[string]V{}, nil
	}
	m := map[string]V{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return m, nil
}

func writeJSONAtomic(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}


// collectAssetFilenames walks a job and returns every asset filename
// referenced by its characters, props, scenes, or final episode video.
func collectAssetFilenames(j *types.Job) []string {
	seen := map[string]struct{}{}
	var out []string
	add := func(u string) {
		if u == "" {
			return
		}
		fn := filepath.Base(u)
		if _, ok := seen[fn]; ok {
			return
		}
		seen[fn] = struct{}{}
		out = append(out, fn)
	}
	for i := range j.Characters { add(j.Characters[i].ImageURL) }
	for i := range j.Props { add(j.Props[i].ImageURL) }
	for _, ep := range j.Episodes {
		for _, sc := range ep.Scenes { add(sc.ImageURL) }
		add(ep.VideoURL)
	}
	return out
}
