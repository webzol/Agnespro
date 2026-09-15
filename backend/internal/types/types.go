// Package types contains shared types used across the studio backend.
package types

import "time"

// JobStatus is the lifecycle state of a generation job.
type JobStatus string

const (
	StatusPending     JobStatus = "pending"
	StatusPlanning    JobStatus = "planning"
	StatusCharacters  JobStatus = "characters"
	StatusProps       JobStatus = "props"
	StatusScenes      JobStatus = "scenes"
	StatusVideo       JobStatus = "video"
	StatusDone        JobStatus = "done"
	StatusFailed      JobStatus = "failed"
	StatusCancelled   JobStatus = "cancelled"
)

// EpisodeState tracks generation progress for a single episode.
type EpisodeState string

const (
	EpisodePending    EpisodeState = "pending"
	EpisodeRunning    EpisodeState = "running"
	EpisodeCharacters EpisodeState = "characters"
	EpisodeProps      EpisodeState = "props"
	EpisodeScenes     EpisodeState = "scenes"
	EpisodeVideo      EpisodeState = "video"
	EpisodeDone       EpisodeState = "done"
	EpisodeFailed     EpisodeState = "failed"
	EpisodeSkipped    EpisodeState = "skipped"
)

// Job is the top-level record for a script submitted to the studio.
type Job struct {
	ID         string         `json:"id"`
	Title      string         `json:"title"`
	Script     string         `json:"script"`
	Style       string         `json:"style"`         // global visual style hint
	AspectRatio  string         `json:"aspect_ratio,omitempty"`  // 9:16 / 16:9 / 1:1 / 4:3 / 3:4 / 21:9
	EpisodeCount int            `json:"episode_count,omitempty"` // 用户选择的目标集数(parser 拆集后按此截断)
	GenreID      string         `json:"genre_id,omitempty"`      // 叙事题材 id(来自风格库)
	GenreName    string         `json:"genre_name,omitempty"`    // 叙事题材中文名
	ScriptModel  string         `json:"script_model,omitempty"`   // AI 剧本模型(chat)
	ImageModel   string         `json:"image_model,omitempty"`    // AI 绘图模型
	VideoModel   string         `json:"video_model,omitempty"`    // AI 视频模型
	Status     JobStatus      `json:"status"`
	Error      string         `json:"error,omitempty"`
	Progress   int            `json:"progress"`       // 0-100
	Episodes   []Episode      `json:"episodes"`
	Characters []Character    `json:"characters"`
	Props      []Prop         `json:"props"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	StartedAt  *time.Time     `json:"started_at,omitempty"`
	FinishedAt *time.Time     `json:"finished_at,omitempty"`
	// Stats
	NumEpisodes  int `json:"num_episodes"`
	NumScenes    int `json:"num_scenes"`
	NumCharacters int `json:"num_characters"`
	NumProps     int `json:"num_props"`
	NumVideos    int `json:"num_videos"`
}

// Episode is one chunk of a Job's script, generated independently.
type Episode struct {
	ID         string        `json:"id"`
	Index      int           `json:"index"`
	Title      string        `json:"title"`
	Body       string        `json:"body"`            // raw script for this episode
	Summary    string        `json:"summary"`         // one-line summary
	State      EpisodeState  `json:"state"`
	Progress   int           `json:"progress"`        // 0-100
	Error      string        `json:"error,omitempty"`
	Characters []string      `json:"characters"`      // character names appearing
	Props      []string      `json:"props"`           // props referenced
	Scenes     []Scene       `json:"scenes"`
	VideoURL   string        `json:"video_url,omitempty"`
	VideoID    string        `json:"video_id,omitempty"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

// Scene is one visual unit of an episode (typically a paragraph or beat).
type Scene struct {
	ID         string    `json:"id"`
	Index      int       `json:"index"`
	Heading    string    `json:"heading"`     // e.g. "INT. CAFE - DAY"
	Narration  string    `json:"narration"`   // descriptive paragraph
	ImagePrompt string   `json:"image_prompt"`// refined prompt used for image gen
	ImageURL   string    `json:"image_url,omitempty"`
	DurationS  float64   `json:"duration_s"`   // target video seconds
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Character is a reusable persona across the whole job.
type Character struct {
	ID         string    `json:"id"`
	JobID      string    `json:"job_id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`            // protagonist, supporting, etc.
	Appearance string    `json:"appearance"`      // description used for image gen
	ImageURL   string    `json:"image_url,omitempty"`
	ImagePrompt string   `json:"image_prompt,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Prop is a reusable object/scene element.
type Prop struct {
	ID         string    `json:"id"`
	JobID      string    `json:"job_id"`
	Name       string    `json:"name"`
	Kind       string    `json:"kind"`            // object, location, vehicle, etc.
	Description string   `json:"description"`
	ImageURL   string    `json:"image_url,omitempty"`
	ImagePrompt string   `json:"image_prompt,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Asset is a generated binary file (image, video) stored on disk.
type Asset struct {
	ID        string    `json:"id"`
	JobID     string    `json:"job_id"`
	Kind      string    `json:"kind"`            // image, video
	RefType   string    `json:"ref_type"`        // character, prop, scene, episode
	RefID     string    `json:"ref_id"`
	Filename  string    `json:"filename"`
	MimeType  string    `json:"mime_type"`
	SizeBytes int64     `json:"size_bytes"`
	CreatedAt time.Time `json:"created_at"`
}

// Settings is the persisted studio configuration.
type Settings struct {
	APIKeySet     bool   `json:"api_key_set"`     // never expose the actual key
	APIRoute      string `json:"api_route"`       // international / international-alt / china
	BaseURL       string `json:"base_url"`
	Concurrency   int    `json:"concurrency"`
	PollIntervalS int    `json:"poll_interval_s"`
	ScriptModel string `json:"script_model,omitempty"` // 默认 AI 剧本模型
	ImageModel  string `json:"image_model,omitempty"`  // 默认 AI 绘图模型
	VideoModel  string `json:"video_model,omitempty"`  // 默认 AI 视频模型
	UpdatedAt     time.Time `json:"updated_at"`
}
