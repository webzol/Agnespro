package pipeline

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"agnesai/studio/internal/agnes"
	"agnesai/studio/internal/cryptox"
	"agnesai/studio/internal/store"
	"agnesai/studio/internal/types"
)

type Pipeline struct {
	Store *store.Store
	Crypt *cryptox.Cipher
	Cfg   Config
}

type Config struct {
	BaseURL      string
	PollInterval time.Duration
	VideoTimeout time.Duration
	VideoSeconds string
	ImageSize    string
	Style        string
	ScriptModel  string // 默认 AI 剧本模型 id(chat);空则用代码 fallback
	ImageModel   string // 默认 AI 绘图模型 id;空则用代码 fallback
	VideoModel   string // 默认 AI 视频模型 id;空则用代码 fallback
}

func New(s *store.Store, c *cryptox.Cipher, cfg Config) *Pipeline {
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = 5 * time.Second
	}
	if cfg.VideoTimeout <= 0 {
		cfg.VideoTimeout = 15 * time.Minute
	}
	if cfg.VideoSeconds == "" {
		cfg.VideoSeconds = "5"
	}
	if cfg.ImageSize == "" {
		cfg.ImageSize = "1024x1024"
	}
	return &Pipeline{Store: s, Crypt: c, Cfg: cfg}
}

// resolveImageSize: 根据 Job.AspectRatio 解析 image size,fallback 到 Cfg.ImageSize
func (p *Pipeline) resolveImageSize(j *types.Job) string {
	if j != nil && j.AspectRatio != "" {
		if img, _ := types.ResolveAspectSize(j.AspectRatio); img != "" {
			return img
		}
	}
	return p.Cfg.ImageSize
}

// resolveVideoSize: 根据 Job.AspectRatio 解析 video size,fallback 到 1280x720
func (p *Pipeline) resolveVideoSize(j *types.Job) string {
	if j != nil && j.AspectRatio != "" {
		if _, vid := types.ResolveAspectSize(j.AspectRatio); vid != "" {
			return vid
		}
	}
	return "1280x720"
}

// pickModel: 优先用 Job 的 model,fallback 到 cfg,最后用硬编码 fallback
func pickModel(jobModel, cfgModel, fallback string) string {
	if jobModel != "" {
		return jobModel
	}
	if cfgModel != "" {
		return cfgModel
	}
	return fallback
}

func (p *Pipeline) client() (*agnes.Client, error) {
	enc, err := p.Store.GetAPIKey()
	if err != nil {
		return nil, fmt.Errorf("read api key: %w", err)
	}
	if enc == "" {
		return nil, errors.New("no Agnes API key configured (set one in Settings)")
	}
	raw, err := p.Crypt.Decrypt(enc)
	if err != nil {
		return nil, fmt.Errorf("decrypt api key: %w", err)
	}
	return agnes.New(p.Cfg.BaseURL, string(raw)), nil
}

func (p *Pipeline) updateJob(j *types.Job) error {
	return p.Store.UpdateJob(j)
}

func (p *Pipeline) Run(ctx context.Context, jobID string, onlyEpisode int) error {
	j, err := p.Store.GetJob(jobID)
	if err != nil {
		return err
	}
	now := time.Now()
	j.StartedAt = &now
	j.Status = types.StatusPlanning
	j.Progress = 5
	if err := p.updateJob(j); err != nil {
		return err
	}

	client, err := p.client()
	if err != nil {
		j.Status = types.StatusFailed
		j.Error = err.Error()
		_ = p.updateJob(j)
		return err
	}

	parsed := ParseEpisodes(j.Script, j.Title)
	j.Episodes = make([]types.Episode, len(parsed))
	for i, ep := range parsed {
		j.Episodes[i] = types.Episode{
			ID:        newID("ep"),
			Index:     ep.Index,
			Title:     ep.Title,
			Body:      ep.Body,
			State:     types.EpisodePending,
			Progress:  0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	j.NumEpisodes = len(j.Episodes)
	j.Characters = nil
	j.Props = nil
	j.Progress = 10
	if err := p.updateJob(j); err != nil {
		return err
	}

	j.Status = types.StatusCharacters
	analysis, summary, err := AnalyzeScript(ctx, client, j.Title, j.Script, j.Style)
	if err != nil {
		j.Status = types.StatusFailed
		j.Error = "analyze: " + err.Error()
		_ = p.updateJob(j)
		return err
	}
	_ = summary
	j.Characters = make([]types.Character, 0, len(analysis.Characters))
	for _, c := range analysis.Characters {
		if strings.TrimSpace(c.Name) == "" {
			continue
		}
		j.Characters = append(j.Characters, types.Character{
			ID:         newID("ch"),
			JobID:      j.ID,
			Name:       strings.TrimSpace(c.Name),
			Role:       strings.TrimSpace(c.Role),
			Appearance: strings.TrimSpace(c.Appearance),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}
	j.Props = make([]types.Prop, 0, len(analysis.Props))
	for _, pr := range analysis.Props {
		if strings.TrimSpace(pr.Name) == "" {
			continue
		}
		j.Props = append(j.Props, types.Prop{
			ID:          newID("pp"),
			JobID:       j.ID,
			Name:        strings.TrimSpace(pr.Name),
			Kind:        strings.TrimSpace(pr.Kind),
			Description: strings.TrimSpace(pr.Description),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}
	j.NumCharacters = len(j.Characters)
	j.NumProps = len(j.Props)
	j.Progress = 20
	if err := p.updateJob(j); err != nil {
		return err
	}

	if len(j.Characters)+len(j.Props) > 0 {
		j.Status = types.StatusCharacters
		items := make([]imageJob, 0, len(j.Characters)+len(j.Props))
		for i := range j.Characters {
			items = append(items, imageJob{
				RefType: "character",
				RefID:   j.Characters[i].ID,
				Prompt:  buildCharacterPrompt(j.Characters[i], p.Cfg.Style),
				Width:   p.Cfg.ImageSize,
			})
		}
		for i := range j.Props {
			items = append(items, imageJob{
				RefType: "prop",
				RefID:   j.Props[i].ID,
				Prompt:  buildPropPrompt(j.Props[i], p.Cfg.Style),
				Width:   p.Cfg.ImageSize,
			})
		}
		if err := p.runImageBatch(ctx, client, j, items, 20, 35); err != nil {
			j.Status = types.StatusFailed
			j.Error = "character images: " + err.Error()
			_ = p.updateJob(j)
			return err
		}
	}

	j.Status = types.StatusScenes
	if err := p.updateJob(j); err != nil {
		return err
	}

	epStart := 0
	epEnd := len(j.Episodes)
	if onlyEpisode > 0 && onlyEpisode <= len(j.Episodes) {
		epStart = onlyEpisode - 1
		epEnd = onlyEpisode
	}
	for i := epStart; i < epEnd; i++ {
		if ctx.Err() != nil {
			j.Status = types.StatusCancelled
			_ = p.updateJob(j)
			return ctx.Err()
		}
		ep := &j.Episodes[i]
		if err := p.runEpisode(ctx, client, j, ep); err != nil {
			ep.State = types.EpisodeFailed
			ep.Error = err.Error()
			j.Status = types.StatusFailed
			j.Error = fmt.Sprintf("episode %d: %s", ep.Index, err.Error())
			_ = p.updateJob(j)
			return err
		}
	}

	j.Status = types.StatusDone
	j.Progress = 100
	done := time.Now()
	j.FinishedAt = &done
	return p.updateJob(j)
}

type imageJob struct {
	RefType string
	RefID   string
	Prompt  string
	Width   string
}

func (p *Pipeline) runImageBatch(ctx context.Context, client *agnes.Client, j *types.Job, items []imageJob, startProg, endProg int) error {
	total := len(items)
	for i, it := range items {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		prog := startProg + (endProg-startProg)*(i+1)/total
		j.Progress = prog
		_ = p.updateJob(j)

		jitter, _ := rand.Int(rand.Reader, big.NewInt(2))
		time.Sleep(time.Duration(jitter.Int64()+1) * time.Second)

		img, err := client.GenerateImage(ctx, agnes.ImageRequest{
			Model:  pickModel(j.ImageModel, p.Cfg.ImageModel, "agnes-image-2.1-flash"),
			Prompt: it.Prompt,
			Size:   p.resolveImageSize(j),
		})
		if err != nil {
			return fmt.Errorf("image %s: %w", it.RefType, err)
		}
		if len(img.Data) == 0 || img.Data[0].URL == "" {
			return fmt.Errorf("image %s: empty response", it.RefType)
		}
		url := img.Data[0].URL
		fn, err := p.downloadAsset(ctx, j.ID, it.RefType, it.RefID, url)
		if err != nil {
			return fmt.Errorf("download %s image: %w", it.RefType, err)
		}
		attachImageToRef(j, it.RefType, it.RefID, fn, it.Prompt)
	}
	return nil
}

func attachImageToRef(j *types.Job, refType, refID, filename, prompt string) {
	switch refType {
	case "character":
		for i := range j.Characters {
			if j.Characters[i].ID == refID {
				j.Characters[i].ImageURL = "/api/assets/" + filename
				j.Characters[i].ImagePrompt = prompt
				j.Characters[i].UpdatedAt = time.Now()
				return
			}
		}
	case "prop":
		for i := range j.Props {
			if j.Props[i].ID == refID {
				j.Props[i].ImageURL = "/api/assets/" + filename
				j.Props[i].ImagePrompt = prompt
				j.Props[i].UpdatedAt = time.Now()
				return
			}
		}
	}
}

func (p *Pipeline) downloadAsset(ctx context.Context, jobID, refType, refID, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	// Use a longer timeout client for asset downloads; the default
	// client has no timeout which can hang on slow upstream.
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("download asset: status %d", resp.StatusCode)
	}
	ext := extensionFromContentType(resp.Header.Get("Content-Type"), url)
	fn := fmt.Sprintf("%s_%s_%s%s", jobID, refType, shortID(refID), ext)
	a, err := p.Store.SaveAsset(jobID, "image", refType, refID, fn, resp.Header.Get("Content-Type"), resp.Body)
	if err != nil {
		return "", err
	}
	return a.Filename, nil
}

func extensionFromContentType(ct, url string) string {
	ct = strings.ToLower(strings.TrimSpace(ct))
	switch {
	case strings.Contains(ct, "png"):
		return ".png"
	case strings.Contains(ct, "jpeg"), strings.Contains(ct, "jpg"):
		return ".jpg"
	case strings.Contains(ct, "webp"):
		return ".webp"
	case strings.Contains(ct, "mp4"):
		return ".mp4"
	}
	lu := strings.ToLower(url)
	switch {
	case strings.HasSuffix(lu, ".png"):
		return ".png"
	case strings.HasSuffix(lu, ".jpg"), strings.HasSuffix(lu, ".jpeg"):
		return ".jpg"
	case strings.HasSuffix(lu, ".webp"):
		return ".webp"
	case strings.HasSuffix(lu, ".mp4"):
		return ".mp4"
	}
	return ".bin"
}

func (p *Pipeline) runEpisode(ctx context.Context, client *agnes.Client, j *types.Job, ep *types.Episode) error {
	ep.State = types.EpisodeRunning
	ep.UpdatedAt = time.Now()
	_ = p.updateJob(j)

	scenes := SplitScenes(ep.Body)
	ep.Scenes = make([]types.Scene, len(scenes))
	for i, sc := range scenes {
		ep.Scenes[i] = types.Scene{
			ID:          newID("sc"),
			Index:       sc.Index,
			Heading:     sc.Heading,
			Narration:   sc.Narration,
			ImagePrompt: buildScenePrompt(ep.Title, sc.Heading, sc.Narration, p.Cfg.Style),
			DurationS:   5.0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	}
	ep.Progress = 10
	_ = p.updateJob(j)

	ep.State = types.EpisodeScenes
	items := make([]imageJob, 0, len(ep.Scenes))
	for i := range ep.Scenes {
		items = append(items, imageJob{
			RefType: "scene",
			RefID:   ep.Scenes[i].ID,
			Prompt:  ep.Scenes[i].ImagePrompt,
			Width:   p.Cfg.ImageSize,
		})
	}
	if err := p.runEpisodeImageBatch(ctx, client, j, ep, items, 10, 60); err != nil {
		return err
	}

	ep.State = types.EpisodeVideo
	ep.Progress = 65
	_ = p.updateJob(j)

	imageURLs := make([]string, 0, len(ep.Scenes))
	for _, sc := range ep.Scenes {
		abs := absAssetURL(sc.ImageURL)
		if abs != "" {
			imageURLs = append(imageURLs, abs)
		}
	}

	videoPrompt := buildVideoPrompt(ep.Title, ep.Scenes, p.Cfg.Style)
	vid, err := client.CreateVideo(ctx, agnes.VideoRequest{
		Model:          pickModel(j.VideoModel, p.Cfg.VideoModel, "agnes-video-v2.0"),
		Prompt:         videoPrompt,
		Seconds:        p.Cfg.VideoSeconds,
		Size:           p.resolveVideoSize(j),
		InputReference: imageURLs,
	})
	if err != nil {
		return fmt.Errorf("create video: %w", err)
	}
	ep.VideoID = vid.VideoID
	ep.Progress = 70
	_ = p.updateJob(j)

	pollCtx, cancel := context.WithTimeout(ctx, p.Cfg.VideoTimeout)
	defer cancel()
	final, err := client.WaitForVideo(pollCtx, vid.VideoID, p.Cfg.PollInterval)
	if err != nil {
		return fmt.Errorf("wait video: %w", err)
	}
	if final.URL == "" {
		return fmt.Errorf("video completed but no URL")
	}
	fn, err := p.downloadAsset(ctx, j.ID, "episode", ep.ID, final.URL)
	if err != nil {
		ep.VideoURL = final.URL
	} else {
		ep.VideoURL = "/api/assets/" + fn
	}
	ep.Progress = 100
	ep.State = types.EpisodeDone
	ep.UpdatedAt = time.Now()
	j.NumVideos++
	return p.updateJob(j)
}

func (p *Pipeline) runEpisodeImageBatch(ctx context.Context, client *agnes.Client, j *types.Job, ep *types.Episode, items []imageJob, startProg, endProg int) error {
	total := len(items)
	if total == 0 {
		return nil
	}
	for i, it := range items {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		prog := startProg + (endProg-startProg)*(i+1)/total
		ep.Progress = prog
		_ = p.updateJob(j)

		jitter, _ := rand.Int(rand.Reader, big.NewInt(2))
		time.Sleep(time.Duration(jitter.Int64()+1) * time.Second)

		img, err := client.GenerateImage(ctx, agnes.ImageRequest{
			Model:  pickModel(j.ImageModel, p.Cfg.ImageModel, "agnes-image-2.1-flash"),
			Prompt: it.Prompt,
			Size:   p.resolveImageSize(j),
		})
		if err != nil {
			return fmt.Errorf("image scene %d: %w", i+1, err)
		}
		if len(img.Data) == 0 || img.Data[0].URL == "" {
			return fmt.Errorf("image scene %d: empty response", i+1)
		}
		url := img.Data[0].URL
		fn, err := p.downloadAsset(ctx, j.ID, "scene", it.RefID, url)
		if err != nil {
			return fmt.Errorf("download scene %d: %w", i+1, err)
		}
		for k := range ep.Scenes {
			if ep.Scenes[k].ID == it.RefID {
				ep.Scenes[k].ImageURL = "/api/assets/" + fn
				ep.Scenes[k].UpdatedAt = time.Now()
				break
			}
		}
	}
	return nil
}

func absAssetURL(p string) string {
	if base := os.Getenv("AGNES_STUDIO_PUBLIC_BASE_URL"); base != "" {
		return strings.TrimRight(base, "/") + p
	}
	return ""
}

func newID(prefix string) string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return prefix + "_" + hex.EncodeToString(b)
}

func shortID(id string) string {
	if len(id) > 12 {
		return id[len(id)-12:]
	}
	return id
}

func buildCharacterPrompt(c types.Character, style string) string {
	parts := []string{}
	if c.Appearance != "" {
		parts = append(parts, c.Appearance)
	} else {
		parts = append(parts, c.Name)
	}
	if c.Role != "" {
		parts = append(parts, "role: "+c.Role)
	}
	if style != "" {
		parts = append(parts, "style: "+style)
	}
	parts = append(parts, "single character, full body, clean background, character design sheet, high detail")
	return strings.Join(parts, ". ")
}

func buildPropPrompt(pr types.Prop, style string) string {
	parts := []string{}
	if pr.Description != "" {
		parts = append(parts, pr.Description)
	} else {
		parts = append(parts, pr.Name)
	}
	if pr.Kind != "" {
		parts = append(parts, "kind: "+pr.Kind)
	}
	if style != "" {
		parts = append(parts, "style: "+style)
	}
	parts = append(parts, "single object, isolated, clean background, high detail")
	return strings.Join(parts, ". ")
}

func buildScenePrompt(epTitle, heading, narration, style string) string {
	parts := []string{}
	if epTitle != "" {
		parts = append(parts, "from: "+epTitle)
	}
	if heading != "" {
		parts = append(parts, heading)
	}
	if narration != "" {
		trim := narration
		if len(trim) > 480 {
			trim = trim[:480] + "..."
		}
		parts = append(parts, trim)
	}
	if style != "" {
		parts = append(parts, "style: "+style)
	}
	parts = append(parts, "cinematic, high detail, no text overlay, no watermark")
	return strings.Join(parts, ". ")
}

func buildVideoPrompt(epTitle string, scenes []types.Scene, style string) string {
	parts := []string{}
	if epTitle != "" {
		parts = append(parts, "Episode: "+epTitle)
	}
	for _, sc := range scenes {
		if sc.Heading != "" {
			parts = append(parts, sc.Heading)
		}
	}
	if style != "" {
		parts = append(parts, "Visual style: "+style)
	}
	parts = append(parts, "Smooth cinematic camera motion, consistent character design, no text overlay, no watermark")
	return strings.Join(parts, "\n")
}

var _ io.Reader = nil