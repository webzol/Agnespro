package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"agnesai/studio/internal/agnes"
)

type generateScriptReq struct {
	Title        string `json:"title"`
	Style        string `json:"style"`         // free-text style hint (legacy)
	VisualStyle  string `json:"visual_style"`   // "cinematic" | "anime" | "ink" | "cyber" | "oil" | "warm" | "noir" | "scifi" | "kid" | "docu"
	VisualStyleName string `json:"visual_style_name"`
	Idea         string `json:"idea"`
	Genre        string `json:"genre"`
	Length       string `json:"length"` // "short" | "medium" | "long"
	Lang         string `json:"lang"`    // "zh" | "en"
}

var visualStyleMap = map[string]string{
	"cinematic": "电影级色调,戏剧光影,浅景深,电影感构图",
	"anime":     "日式动漫风格,明亮色块,清晰线条,二次元",
	"ink":       "中国水墨画风,留白意境,水墨晕染",
	"cyber":     "赛博朋克,霓虹紫蓝,机械城市,未来感",
	"oil":       "油画质感,厚重笔触,古典色彩,博物馆级",
	"warm":      "暖色胶片,怀旧黄绿,胶片颗粒,80年代质感",
	"noir":      "黑白电影,高对比,光影强烈,侦探片",
	"scifi":     "未来科幻,冷蓝紫,金属质感,太空感",
	"kid":       "童趣插画,明亮色彩,可爱风格,儿童绘本",
	"docu":      "纪录片,自然色调,真实朴素,纪实摄影",
}

type generateScriptResp struct {
	Script string `json:"script"`
	Title  string `json:"title"`
}

func (s *Server) generateScript(w http.ResponseWriter, r *http.Request) {
	var req generateScriptReq
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}
	idea := strings.TrimSpace(req.Idea)
	if idea == "" {
		writeError(w, http.StatusBadRequest, "idea is required")
		return
	}
	if req.Length == "" {
		req.Length = "medium"
	}
	if req.Lang == "" {
		req.Lang = "zh"
	}
	if req.Genre == "" {
		req.Genre = "温情短剧"
	}
	if req.Title == "" {
		req.Title = "未命名剧本"
	}
	// Map the visual style id to a concrete Chinese description that
	// the LLM can use to drive the image prompts downstream.
	if req.VisualStyle == "" {
		req.VisualStyle = "cinematic"
	}
	if hint, ok := visualStyleMap[req.VisualStyle]; ok {
		req.Style = hint
	} else if req.Style == "" {
		req.Style = "电影感"
	}

	// Decrypt the API key. We need it to call the LLM.
	enc, err := s.Store.GetAPIKey()
	if err != nil || enc == "" {
		writeError(w, http.StatusBadRequest, "no API key configured (set one in Settings)")
		return
	}
	raw, err := s.Crypt.Decrypt(enc)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "decrypt api key: "+err.Error())
		return
	}
	client := agnes.New(s.Cfg.BaseURL, string(raw))

	// Length hint maps to approximate scene count and total character budget.
	lengthHint := map[string]string{
		"short":  "约 3 场戏,总长 400-600 字",
		"medium": "约 5 场戏,总长 800-1200 字",
		"long":   "约 8 场戏,总长 1500-2000 字",
	}[req.Length]
	if lengthHint == "" {
		lengthHint = "约 5 场戏,总长 800-1200 字"
	}

	langHint := "中文"
	if req.Lang == "en" {
		langHint = "English"
	}

	prompt := strings.Join([]string{
		"你是一位经验丰富的短剧编剧。请根据用户的想法,创作一份可直接用于视频生成的剧本。",
		"",
		"## 任务",
		"- 类型: " + req.Genre,
		"- 视觉风格: " + orDefault(req.Style, "电影感"),
		"- 长度: " + lengthHint,
		"- 语言: " + langHint,
		"- 标题: " + req.Title,
		"",
		"## 用户想法",
		idea,
		"",
		"## 写作要求",
		"1. 开头一句话交代时间、地点、人物。",
		"2. 主体分若干集(使用「第N集」或 `## 第N集` 开头),每集 1-3 个段落。",
		"3. 每段以画面/动作为主,适合直接喂给图像模型生成首帧。",
		"4. 出现的人物请给出可视化的外貌描述(便于后续抽取角色)。",
		"5. 关键道具/地点请在首次出现时显式命名(便于后续抽取道具)。",
		"6. 结尾留下一个反转或钩子,让观众想看下一集。",
		"7. 直接输出剧本正文,不要解释、不要 markdown 标题(除第N集外)、不要前言后记。",
	}, "\n")

	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	resp, err := client.Chat(ctx, agnes.ChatRequest{
		Model: "agnes-2.5-flash",
		Messages: []agnes.ChatMessage{
			{Role: "system", Content: "You are a concise screenwriter. Reply in the requested language only, with the script body."},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, "LLM call failed: "+err.Error())
		return
	}
	if len(resp.Choices) == 0 {
		writeError(w, http.StatusBadGateway, "LLM returned no choices")
		return
	}
	script := strings.TrimSpace(resp.Choices[0].Message.Content)
	writeJSON(w, http.StatusOK, generateScriptResp{Script: script, Title: req.Title})
}

func orDefault(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}