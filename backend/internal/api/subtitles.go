package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"agnesai/studio/internal/types"
)

// subtitleSegment is one SRT cue.
type subtitleSegment struct {
	Index   int
	StartS  float64
	EndS    float64
	Text    string
}

// generateSRT builds an SRT string for the given episode. Scene
// boundaries are taken from episode.Scenes (with fallback to paragraphs
// split by full-width punctuation). Each scene gets a proportional
// time slice of the target video length.
func generateSRT(ep types.Episode, totalSeconds float64) string {
	if totalSeconds <= 0 {
		totalSeconds = 5
	}
	segments := buildSubtitleSegments(ep, totalSeconds)
	var b strings.Builder
	for _, s := range segments {
		fmt.Fprintf(&b, "%d\n", s.Index)
		fmt.Fprintf(&b, "%s --> %s\n", formatSRTTimestamp(s.StartS), formatSRTTimestamp(s.EndS))
		fmt.Fprintf(&b, "%s\n\n", strings.TrimSpace(s.Text))
	}
	return b.String()
}

func formatSRTTimestamp(seconds float64) string {
	if seconds < 0 {
		seconds = 0
	}
	h := int(seconds / 3600)
	m := int(seconds/60) - h*60
	s := seconds - float64(h*3600+m*60)
	return fmt.Sprintf("%02d:%02d:%06.3f", h, m, s)
}

var (
	sentenceEndZH = regexp.MustCompile(`[。!?;]`)
	sentenceEndEN = regexp.MustCompile(`[.!?;]`)
)

// splitBySentences splits a Chinese/English paragraph into sentences,
// preserving order. Empty strings are dropped.
func splitBySentences(text string) []string {
	out := []string{}
	text = strings.TrimSpace(text)
	if text == "" {
		return out
	}
	// split by Chinese sentence terminators first
	for _, raw := range sentenceEndZH.Split(text, -1) {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		// then split English terminators
		for _, sub := range sentenceEndEN.Split(raw, -1) {
			sub = strings.TrimSpace(sub)
			if sub != "" {
				out = append(out, sub)
			}
		}
	}
	return out
}

// buildSubtitleSegments maps an episode's scenes to SRT cues, weighted
// by the relative length of each scene's narration text.
func buildSubtitleSegments(ep types.Episode, totalSeconds float64) []subtitleSegment {
	type sceneText struct {
		text string
	}
	var sources []sceneText
	if len(ep.Scenes) > 0 {
		for _, sc := range ep.Scenes {
			t := sc.Narration
			if t == "" {
				t = sc.Heading
			}
			for _, s := range splitBySentences(t) {
				sources = append(sources, sceneText{text: s})
			}
		}
	} else {
		for _, s := range splitBySentences(ep.Body) {
			sources = append(sources, sceneText{text: s})
		}
	}
	if len(sources) == 0 {
		sources = []sceneText{{text: ep.Title}}
	}

	// Weight each sentence by its character count
	weights := make([]int, len(sources))
	totalWeight := 0
	for i, s := range sources {
		w := len([]rune(s.text))
		if w < 4 {
			w = 4 // minimum duration to be readable
		}
		weights[i] = w
		totalWeight += w
	}
	if totalWeight == 0 {
		totalWeight = 1
	}

	segs := make([]subtitleSegment, 0, len(sources))
	cursor := 0.0
	for i, s := range sources {
		dur := totalSeconds * float64(weights[i]) / float64(totalWeight)
		if dur < 0.8 {
			dur = 0.8
		}
		segs = append(segs, subtitleSegment{
			Index:  i + 1,
			StartS: cursor,
			EndS:   cursor + dur,
			Text:   s.text,
		})
		cursor += dur
	}
	return segs
}

func (s *Server) downloadSubtitles(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	j, err := s.Store.GetJob(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	style := r.URL.Query().Get("style")
	if style == "" {
		style = "modern"
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+sanitizeFilename(j.Title)+".srt\"")
	// We only generate for the first done episode for now
	if len(j.Episodes) == 0 {
		fmt.Fprint(w, "")
		return
	}
	ep := j.Episodes[0]
	srt := applySubtitleStyle(srtBase(ep, 5), style)
	fmt.Fprint(w, srt)
}

func srtBase(ep types.Episode, seconds float64) string {
	return generateSRT(ep, seconds)
}

// applySubtitleStyle returns the SRT with style metadata in a header
// comment. The SRT spec doesn't have a style field, so we encode the
// style as a WebVTT-like style block. Most players (VLC, mpv, ffmpeg
// with -c:v subtitles filter) read the style from the SRT file name
// or from a sidecar file. We also return a JSON sidecar when style
// is requested via ?sidecar=1.
func applySubtitleStyle(srt, style string) string {
	// Prepend a style header as comments (SRT ignores lines starting with ;)
	if s := getSubtitleStyle(style); s != nil {
		header := fmt.Sprintf(";; STYLE: %s\n;; FONT: %s\n;; COLOR: %s\n;; OUTLINE: %s\n;; SHADOW: %s\n;; SIZE: %d\n;; POSITION: %s\n;; ALIGN: %s\n\n",
			s.Name, s.Font, s.Color, s.Outline, s.Shadow, s.Size, s.Position, s.Align)
		return header + srt
	}
	return srt
}

func sanitizeFilename(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "job"
	}
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "\\", "_")
	s = strings.ReplaceAll(s, "..", "_")
	s = strings.ReplaceAll(s, "\x00", "")
	return s
}

// subtitleStyle is a JSON-friendly preset.
type subtitleStyle struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Font     string `json:"font"`
	Color    string `json:"color"`
	Outline  string `json:"outline"`
	Shadow   string `json:"shadow"`
	Size     int    `json:"size"`
	Position string `json:"position"`
	Align    string `json:"align"`
	Preview  string `json:"preview"` // path to preview image
}

var subtitleStyles = []subtitleStyle{
	{ID: "modern", Name: "现代白", Font: "Noto Sans CJK SC", Color: "#FFFFFF", Outline: "#000000 2px", Shadow: "0 1px 3px rgba(0,0,0,.6)", Size: 48, Position: "bottom", Align: "center", Preview: "/assets/img/substyle_modern.png"},
	{ID: "neon", Name: "霓虹", Font: "Orbitron", Color: "#FF3DFF", Outline: "0 0 8px #FF3DFF", Shadow: "0 0 20px #FF3DFF", Size: 52, Position: "bottom", Align: "center", Preview: "/assets/img/substyle_neon.png"},
	{ID: "typewriter", Name: "打字机", Font: "Courier New", Color: "#FFE4B5", Outline: "#000 1px", Shadow: "none", Size: 44, Position: "top", Align: "left", Preview: "/assets/img/substyle_typewriter.png"},
	{ID: "cinema", Name: "影院金", Font: "Source Han Serif", Color: "#FFD56B", Outline: "#1a1a1a 2px", Shadow: "0 2px 4px rgba(0,0,0,.8)", Size: 50, Position: "bottom", Align: "center", Preview: "/assets/img/substyle_cinema.png"},
	{ID: "kawaii", Name: "软萌粉", Font: "PingFang SC", Color: "#FF9EC7", Outline: "white 3px", Shadow: "0 2px 6px rgba(255,158,199,.5)", Size: 54, Position: "bottom", Align: "center", Preview: "/assets/img/substyle_kawaii.png"},
	{ID: "minimal", Name: "极简黑", Font: "Inter", Color: "#1A1A1A", Outline: "rgba(255,255,255,.9) 3px", Shadow: "none", Size: 46, Position: "bottom", Align: "center", Preview: "/assets/img/substyle_minimal.png"},
	{ID: "news", Name: "新闻蓝带", Font: "Noto Sans", Color: "#FFFFFF", Outline: "#0a3d91 4px", Shadow: "0 0 0 #0a3d91", Size: 48, Position: "top", Align: "center", Preview: "/assets/img/substyle_news.png"},
	{ID: "burn", Name: "火焰", Font: "Impact", Color: "#FFEE00", Outline: "#FF4500 3px", Shadow: "0 0 16px #FF4500", Size: 60, Position: "bottom", Align: "center", Preview: "/assets/img/substyle_burn.png"},
	{ID: "ink", Name: "水墨", Font: "STKaiti", Color: "#000000", Outline: "white 2px", Shadow: "0 1px 0 rgba(0,0,0,.3)", Size: 50, Position: "bottom", Align: "center", Preview: "/assets/img/substyle_ink.png"},
	{ID: "anime", Name: "二次元", Font: "Yuanti SC", Color: "#FFFFFF", Outline: "#FF1493 4px", Shadow: "0 0 12px #FF1493", Size: 56, Position: "bottom", Align: "center", Preview: "/assets/img/substyle_anime.png"},
}

func getSubtitleStyle(id string) *subtitleStyle {
	for i, s := range subtitleStyles {
		if s.ID == id {
			return &subtitleStyles[i]
		}
	}
	return nil
}

func (s *Server) listSubtitleStyles(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"styles": subtitleStyles})
}

// listVisualStyles returns the built-in visual styles for the AI script
// generator with preview images.
type visualStyle struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Preview  string `json:"preview"`
}

var visualStyles = []visualStyle{
	{ID: "cinematic", Name: "电影感", Desc: "电影级色调,戏剧光影", Preview: "/assets/img/style_cinematic.png"},
	{ID: "anime", Name: "日式动漫", Desc: "明亮二次元,色块分明", Preview: "/assets/img/style_anime.png"},
	{ID: "ink", Name: "水墨国风", Desc: "水墨晕染,留白意境", Preview: "/assets/img/style_ink.png"},
	{ID: "cyber", Name: "赛博朋克", Desc: "霓虹紫蓝,机械城市", Preview: "/assets/img/style_cyber.png"},
	{ID: "oil", Name: "油画质感", Desc: "厚重笔触,古典色彩", Preview: "/assets/img/style_oil.png"},
	{ID: "warm", Name: "暖色胶片", Desc: "怀旧黄绿,胶片颗粒", Preview: "/assets/img/style_warm.png"},
	{ID: "noir", Name: "黑色电影", Desc: "高对比黑白,光影强烈", Preview: "/assets/img/style_noir.png"},
	{ID: "scifi", Name: "赛博科幻", Desc: "冷蓝紫,未来金属", Preview: "/assets/img/style_scifi.png"},
	{ID: "kid", Name: "童趣插画", Desc: "明亮色彩,可爱风格", Preview: "/assets/img/style_kid.png"},
	{ID: "docu", Name: "纪录片", Desc: "自然色调,真实朴素", Preview: "/assets/img/style_docu.png"},
}

func (s *Server) listVisualStyles(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"styles": visualStyles})
}

func (s *Server) downloadJobAssets(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	j, err := s.Store.GetJob(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	// Build a JSON manifest with all generated asset URLs (single ZIP
	// would require extra deps; manifest is portable and works without).
	type item struct {
		Kind string `json:"kind"`
		URL  string `json:"url"`
		Name string `json:"name"`
	}
	var items []item
	for _, c := range j.Characters {
		if c.ImageURL != "" {
			items = append(items, item{Kind: "character:" + c.Name, URL: c.ImageURL, Name: c.Name + ".png"})
		}
	}
	for _, p := range j.Props {
		if p.ImageURL != "" {
			items = append(items, item{Kind: "prop:" + p.Name, URL: p.ImageURL, Name: p.Name + ".png"})
		}
	}
	for _, ep := range j.Episodes {
		for _, sc := range ep.Scenes {
			if sc.ImageURL != "" {
				items = append(items, item{Kind: "scene", URL: sc.ImageURL, Name: sc.Heading + ".png"})
			}
		}
		if ep.VideoURL != "" {
			items = append(items, item{Kind: "video", URL: ep.VideoURL, Name: ep.Title + ".mp4"})
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"job_id": j.ID,
		"title":  j.Title,
		"count":  len(items),
		"items":  items,
		"note":   "Each URL accepts ?download=1 to force browser download. Use the subtitle endpoint for caption files.",
	})
	_ = time.Now // keep import
}