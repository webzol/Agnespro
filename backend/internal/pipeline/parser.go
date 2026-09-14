// Package pipeline orchestrates the full generation flow:
// script -> episodes -> characters & props -> per-episode scenes & video.
package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"agnesai/studio/internal/agnes"
)

// ParseEpisodes splits a raw script into episodes. Detection rules, in
// order:
//   1. Markdown headings starting with "## " or "### " whose text looks
//      like an episode marker (e.g. "第一集", "Episode 2", "EP03").
//   2. Plain-text lines matching patterns like "第N集", "Episode N",
//      "EP N", "Ep.N", "Scene N".
//   3. If no markers are found, the whole script is treated as a
//      single episode.
func ParseEpisodes(script, fallbackTitle string) []ParsedEpisode {
	script = strings.TrimSpace(script)
	if script == "" {
		return []ParsedEpisode{{
			Index: 1,
			Title: fallbackTitle,
			Body:  "",
		}}
	}

	lines := strings.Split(script, "\n")
	type marker struct {
		line int
		text string
	}
	var markers []marker
	epRe := regexp.MustCompile(`(?im)^\s*(?:#{1,6}\s+)?(?:第?\s*([0-9一二三四五六七八九十百千]+)\s*[集回章]|episode\s*([0-9]+)|ep\s*[.:-]?\s*([0-9]+)|scene\s*([0-9]+))`)
	for i, ln := range lines {
		if m := epRe.FindStringSubmatch(ln); m != nil {
			title := strings.TrimSpace(strings.TrimLeft(ln, "#"))
			title = epRe.ReplaceAllString(title, "")
			title = strings.TrimSpace(strings.TrimLeft(title, ":：-—"))
			if title == "" {
				title = fmt.Sprintf("Episode %d", len(markers)+1)
			}
			markers = append(markers, marker{line: i, text: title})
		}
	}

	if len(markers) == 0 {
		return []ParsedEpisode{{
			Index: 1,
			Title: fallbackTitle,
			Body:  script,
		}}
	}

	var eps []ParsedEpisode
	for idx, mk := range markers {
		start := mk.line + 1
		end := len(lines)
		if idx+1 < len(markers) {
			end = markers[idx+1].line
		}
		body := strings.TrimSpace(strings.Join(lines[start:end], "\n"))
		eps = append(eps, ParsedEpisode{
			Index: idx + 1,
			Title: mk.text,
			Body:  body,
		})
	}
	return eps
}

type ParsedEpisode struct {
	Index int
	Title string
	Body  string
}

type ScriptAnalysis struct {
	Characters []CharacterHint `json:"characters"`
	Props      []PropHint      `json:"props"`
}

type CharacterHint struct {
	Name       string `json:"name"`
	Role       string `json:"role"`
	Appearance string `json:"appearance"`
}

type PropHint struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
}

func AnalyzeScript(ctx context.Context, client *agnes.Client, title, script, style string) (*ScriptAnalysis, string, error) {
	prompt := fmt.Sprintf(`You are a script supervisor preparing assets for a short video adaptation.

Title: %s
Style hint: %s

Script:
"""
%s
"""

Tasks:
1. Extract every recurring CHARACTER. For each, output: name, role (e.g. protagonist, supporting, antagonist), and a one-paragraph visual appearance description suitable for an image generator.
2. Extract every recurring PROP or LOCATION (objects, vehicles, key places). For each, output: name, kind (object|location|vehicle|other), and a one-paragraph visual description.
3. Reply ONLY with a JSON object of the form:
{"characters":[{"name":"...","role":"...","appearance":"..."}],"props":[{"name":"...","kind":"...","description":"..."}]}

Do not include any prose, markdown fences, or commentary. The response must be valid JSON.`, title, style, script)

	var resp *agnes.ChatResponse
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		resp, err = client.Chat(ctx, agnes.ChatRequest{
			Model: "agnes-2.5-flash",
			Messages: []agnes.ChatMessage{
				{Role: "system", Content: "You reply with strict JSON only."},
				{Role: "user", Content: prompt},
			},
		})
		if err == nil {
			break
		}
		if ctx.Err() != nil {
			return nil, "", ctx.Err()
		}
		select {
		case <-ctx.Done():
			return nil, "", ctx.Err()
		case <-time.After(time.Duration(attempt) * time.Second):
		}
	}
	if err != nil {
		return nil, "", err
	}
	if len(resp.Choices) == 0 {
		return nil, "", fmt.Errorf("agnes chat returned no choices")
	}
	raw := resp.Choices[0].Message.Content
	raw = stripCodeFence(raw)
	var out ScriptAnalysis
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return &ScriptAnalysis{}, strings.TrimSpace(raw), nil
	}
	summary := strings.TrimSpace(raw)
	if len(summary) > 240 {
		summary = summary[:240] + "..."
	}
	return &out, summary, nil
}

func stripCodeFence(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		if i := strings.Index(s, "\n"); i >= 0 {
			s = s[i+1:]
		}
		if j := strings.LastIndex(s, "```"); j >= 0 {
			s = s[:j]
		}
	}
	return strings.TrimSpace(s)
}

func SplitScenes(body string) []ParsedScene {
	body = strings.TrimSpace(body)
	if body == "" {
		return []ParsedScene{{
			Index:     1,
			Heading:   "Opening",
			Narration: body,
		}}
	}
	paragraphs := regexp.MustCompile(`\n\s*\n`).Split(body, -1)
	var scenes []ParsedScene
	idx := 0
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		idx++
		heading := firstSentence(p)
		if heading == "" {
			heading = fmt.Sprintf("Scene %d", idx)
		}
		scenes = append(scenes, ParsedScene{
			Index:     idx,
			Heading:   heading,
			Narration: p,
		})
	}
	if len(scenes) == 0 {
		scenes = append(scenes, ParsedScene{Index: 1, Heading: "Opening", Narration: body})
	}
	return scenes
}

type ParsedScene struct {
	Index     int
	Heading   string
	Narration string
}

func firstSentence(p string) string {
	stop := -1
	for i, r := range p {
		if r == '.' || r == '!' || r == '?' || r == '。' {
			stop = i
			break
		}
		if i > 120 {
			stop = i
			break
		}
	}
	if stop < 0 {
		if len(p) > 80 {
			return strings.TrimSpace(p[:80]) + "..."
		}
		return strings.TrimSpace(p)
	}
	return strings.TrimSpace(p[:stop+1])
}