# Agnes AI Studio — 开发文档

> 面向开发者的架构说明、扩展指南与运维手册。

---

## 1. 架构总览

```
┌────────────────────────────────────────────────────────────┐
│                    Browser (Frontend)                       │
│  HTML5 + CSS3 (iOS Liquid Glass) + Vanilla JS               │
│  - 单页应用,无构建步骤                                       │
│  - 鼠标跟随高光、毛玻璃背景、动态进度条                        │
└─────────────────────────────┬──────────────────────────────┘
                              │ fetch /api/*
┌─────────────────────────────▼──────────────────────────────┐
│              Go HTTP Server (net/http, Go 1.23+)            │
│  - net/http ServeMux with method routing (Go 1.22+ pattern)│
│  - All API routes mounted under /api/*                      │
│  - Frontend static files served at /                        │
└────────┬──────────────┬─────────────┬─────────────┬─────────┘
         │              │             │             │
   ┌─────▼─────┐  ┌─────▼──────┐ ┌────▼─────┐  ┌────▼─────┐
   │  Store    │  │ Pipeline   │ │  Queue   │  │ AgnesAI  │
   │ (JSON FS) │  │ Orchestr.  │ │ (Worker) │  │  Client  │
   └───────────┘  └────────────┘ └──────────┘  └────┬─────┘
   data/                                          HTTP
   - jobs.json                                    │
   - .apikey (encrypted)                          ▼
   - .masterkey                          ┌────────────────────┐
   - settings.json                       │   apihub.agnes-ai  │
   - assets/                             │     .com / .cn     │
                                         │  (OpenAI-compat)   │
                                         └────────────────────┘
```

### 1.1 关键模块

| 包 | 职责 |
|---|---|
| `internal/config` | 加载环境变量,解析路由 → base URL,管理主密钥路径 |
| `internal/cryptox` | AES-256-GCM 加密 API key(32 字节主密钥) |
| `internal/store` | 任务/设置的 JSON 持久化、Asset 文件管理 |
| `internal/types` | 所有跨包共享的结构体 |
| `internal/agnes` | 文本/图像/视频三个 API 客户端 + 视频轮询 |
| `internal/pipeline` | 任务编排:解析 → 角色 → 道具 → 场景 → 视频 |
| `internal/queue` | 内存任务队列,支持取消与并发 |
| `internal/api` | HTTP handlers(router / middleware / jobs / settings / assets / queue / docs) |
| `cmd/server` | 入口;加载 .env、解析 flag、组装各组件、启动 HTTP server |

### 1.2 数据流(单个剧集)

```
User clicks "Run all" on /api/jobs/{id}
       │
       ▼
queue.Enqueue(Job{Run: pipeline.Run})
       │
       ▼ (when worker free)
pipeline.Run(ctx, jobID, 0)
       │
       ├── Status: planning → parse script → 5%
       ├── Status: characters
       │     ├── chat: extract characters & props (10–20%)
       │     └── for each: image gen (20–35%)
       ├── Status: scenes (per episode, sequential)
       │     ├── split into scenes
       │     ├── for each scene: image gen (10–60% per ep)
       │     ├── POST /v1/videos with input_reference (65%)
       │     ├── poll /agnesapi?video_id=… (70–99%)
       │     └── download final mp4 to data/assets
       └── Status: done (100%)
```

---

## 2. 项目结构

```
agnesai/
├── backend/
│   ├── cmd/server/main.go            # 入口
│   ├── internal/
│   │   ├── agnes/                    # AgnesAI 客户端
│   │   │   ├── client.go             # HTTP 客户端基类
│   │   │   ├── chat.go               # 文本补全
│   │   │   ├── image.go              # 图像生成
│   │   │   └── video.go              # 视频生成 + 轮询
│   │   ├── api/                      # HTTP handlers
│   │   │   ├── router.go
│   │   │   ├── middleware.go
│   │   │   ├── jobs.go
│   │   │   ├── settings.go
│   │   │   ├── assets.go
│   │   │   ├── queue.go
│   │   │   └── docs.go
│   │   ├── pipeline/                 # 生成流程编排
│   │   │   ├── parser.go             # 剧本解析、角色/道具抽取
│   │   │   └── pipeline.go           # 主流程
│   │   ├── queue/queue.go            # 内存任务队列
│   │   ├── store/store.go            # JSON 持久化
│   │   ├── cryptox/cryptox.go        # AES-256-GCM
│   │   ├── config/config.go          # 环境变量
│   │   └── types/types.go            # 共享结构
│   ├── web/                          # 嵌入式前端
│   │   ├── index.html
│   │   └── assets/
│   │       ├── style.css             # iOS 26 Liquid Glass
│   │       ├── markdown.js           # 极简 MD 渲染器
│   │       └── app.js                # 路由、API、状态
│   ├── go.mod
│   └── .env.example
├── docs/
│   ├── DEVELOPMENT.md                # 本文件
│   └── USER_GUIDE.md                 # 用户文档(也内置于前端)
├── .env.example
├── .gitignore
├── README.md
└── AGENTS.md
```

---

## 3. 本地开发

### 3.1 前置要求

- Go 1.23+(用到 `http.ServeMux` 的方法路由语法)
- Agnes AI API key — 申请:https://platform.agnes-ai.com/
- 可选:Go 1.22+ 标准库即可,无第三方依赖

### 3.2 构建

```bash
cd backend
go build -o server.exe ./cmd/server
./server.exe
```

第一次构建可能需要 ~30 秒(标准库下载/编译缓存)。

### 3.3 配置

复制 `.env.example` 为 `.env` 并按需修改,或直接用 flag / 环境变量:

```bash
./server.exe -addr :9000 -data ./mydata -route international-alt
```

### 3.4 前端开发

前端是纯静态文件,直接修改 `backend/web/` 下的文件即可,刷新浏览器即生效。

推荐浏览器 DevTools Network 面板:
- `GET /api/jobs` — 任务列表
- `GET /api/jobs/{id}` — 任务详情(轮询每 3s)
- `GET /api/queue` — 队列状态(轮询每 4s)
- `GET /api/assets/{name}` — 图片/视频文件

### 3.5 修改后重新部署

```bash
# 1. 停止旧服务
Ctrl+C

# 2. 重新构建
go build -o server.exe ./cmd/server

# 3. 启动新服务
./server.exe
```

数据(data/)会被保留。

---

## 4. 扩展指南

### 4.1 添加新模型(例如 `agnes-llm-pro`)

修改 `internal/agnes/chat.go`:

```go
func (c *Client) ProChat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
    req.Model = "agnes-llm-pro"
    return c.Chat(ctx, req)
}
```

### 4.2 添加新的 API 端点

在 `internal/api/router.go` 注册:

```go
mux.HandleFunc("GET /api/foo", s.handleFoo)
```

在对应文件实现 handler,用 `writeJSON` / `writeError` 输出,`decodeJSON` 读取 body。

### 4.3 自定义剧集拆分规则

`internal/pipeline/parser.go` 的 `ParseEpisodes` 内:

```go
epRe := regexp.MustCompile(`(?i)^\\s*(?:#{1,6}\\s+)?...`)
```

按需修改正则。

### 4.4 修改前端液态玻璃效果

`backend/web/assets/style.css` 的 `:root` 变量:

```css
--glass-tint: rgba(255, 255, 255, 0.06);
--glass-border: rgba(255, 255, 255, 0.16);
--grad: linear-gradient(135deg, #7c8cff 0%, #38d6ff 50%, #ff5fa2 100%);
```

修改后刷新浏览器即可。

### 4.5 切换数据库(JSON → SQLite)

把 `internal/store/store.go` 的 `readJobs / writeJobs` 替换为 SQLite 调用即可。推荐用 `modernc.org/sqlite`(纯 Go,无 CGO)。其余 handler 不需改动。

### 4.6 添加鉴权(防止公网暴露)

最简单的办法 — 在 `internal/api/router.go` 给 mux 包一层:

```go
func (s *Server) authMW(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("X-Studio-Token")
        if token == "" || token != os.Getenv("AGNES_STUDIO_AUTH_TOKEN") {
            http.Error(w, "forbidden", http.StatusForbidden); return
        }
        next.ServeHTTP(w, r)
    })
}
```

但更推荐:把后端放在 nginx / cloudflare 后面,只暴露 `/` 给最终用户,把 `/api` 限制为内网。

---

## 5. API 客户端示例

`internal/agnes` 是薄封装,任何 Go 程序都可以独立 import:

```go
import "agnesai/studio/internal/agnes"

c := agnes.New("https://apihub.agnes-ai.com/v1", os.Getenv("AGNES_API_KEY"))

// 1. 文本
resp, err := c.Chat(ctx, agnes.ChatRequest{
    Model: "agnes-2.5-flash",
    Messages: []agnes.ChatMessage{
        {Role: "system", Content: "You are concise."},
        {Role: "user", Content: "Say hi"},
    },
})
fmt.Println(resp.Choices[0].Message.Content)

// 2. 图像
img, err := c.GenerateImage(ctx, agnes.ImageRequest{
    Prompt: "a red apple, watercolor",
    Size:   "1024x1024",
})
fmt.Println(img.Data[0].URL)

// 3. 视频
vid, err := c.CreateVideo(ctx, agnes.VideoRequest{
    Prompt:  "a leaf falling in slow motion",
    Seconds: "5",
})
final, err := c.WaitForVideo(ctx, vid.VideoID, 5*time.Second)
fmt.Println(final.URL) // 视频的最终 URL
```

---

## 6. 安全模型

| 风险 | 缓解措施 |
|---|---|
| API key 泄露 | AES-256-GCM 加密落盘;主密钥可由 `AGNES_STUDIO_MASTER_KEY` 提供(否则自动生成在 `data/.masterkey`,0600 权限) |
| 任意脚本执行 | 不解释用户脚本,只调用 LLM 提取结构化字段 |
| 路径遍历 | `internal/api/assets.go` 的 `serveAsset` 拒绝含 `/\\` 或 `..` 的文件名 |
| CORS | 启用宽松 CORS(开发友好);生产建议改严格白名单 |
| 速率限制 | 自由档 1 RPM 视频;已内置 1–2s 抖动 |
| 公网暴露 | 无内置鉴权;生产请用反向代理或加 `authMW`(见 4.6) |

---

## 7. 路线图

- [ ] 多任务并发(目前 `AGNES_STUDIO_MAX_CONCURRENCY=1`)
- [ ] 视频拼接(多集合成一个长视频)
- [ ] 实时进度推送(WebSocket)
- [ ] 用户/多账号系统
- [ ] 字幕/配音集成
- [ ] 移动端原生 App

---

## 8. 故障排查

| 现象 | 原因 | 解决 |
|---|---|---|
| 启动报 `master key file: ...` | 首次启动需写 `data/.masterkey` | 检查 `data/` 是否可写 |
| 视频引用图失败 | 本地服务不可公网访问 | 设置 `AGNES_STUDIO_PUBLIC_BASE_URL` 或用隧道 |
| `400 invalid_request_error` | prompt 触发平台过滤 | 调整 prompt,避免敏感词 |
| `429 rate_limit_exceeded` | RPM 超限 | 降低 `AGNES_STUDIO_MAX_CONCURRENCY` 或升级套餐 |
| `502 bad_gateway` | 平台侧故障 | 等待后重试,或切到 `international-alt` 路由 |
| 前端路由跳到 `/jobs` 显示空白 | `data/jobs.json` 损坏 | 停止服务,删除 `data/jobs.json`,重启 |

---

## 9. 许可证

仅供学习和个人项目使用。商用请联系 AgnesAI 官方(platform.agnes-ai.com)。