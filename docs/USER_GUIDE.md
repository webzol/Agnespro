# Agnes AI Studio — 使用文档

> 一个把剧本自动变成短剧的工作台。粘贴剧本,后端先用文本模型抽取角色与道具,再用图像模型生成首帧图,最后用视频模型合成短剧。所有操作都在浏览器里完成。

---

## 1. 快速开始

### 1.1 启动服务

```bash
# 进入后端目录
cd backend

# 拉取依赖并构建
go build -o server.exe ./cmd/server   # Windows
# go build -o server ./cmd/server      # macOS / Linux

# 运行
./server.exe
```

默认监听 `http://localhost:8080`。首次启动会在 `data/` 下创建:
- `data/.masterkey` — 用于加密 API key 的主密钥(自动生成,妥善保管)
- `data/.apikey` — 加密后的 Agnes API key
- `data/jobs.json` — 任务持久化
- `data/settings.json` — 设置持久化
- `data/assets/` — 生成的图片与视频

### 1.2 配置 API Key

1. 打开 `http://localhost:8080`,进入 **Settings** 页
2. 填入从 [platform.agnes-ai.com](https://platform.agnes-ai.com/) 申请的 API key(以 `sk-` 开头)
3. 点击 **测试连接**,确认可达后点击 **保存**
4. 选择合适的服务路由:
   - **国际(主)**: `apihub.agnes-ai.com` — 海外首选
   - **国际(备)**: `apihub.agnes-ai.cn` — 备用,主线路不可达时使用
   - **中国大陆**: `api.agnes-ai.cn` — 中国大陆服务
5. 修改路由后需重启服务生效

> ⚠️ 路由变更需重启服务,API key 已加密保存无需重启。

### 1.3 创建第一个任务

1. 进入 **New Job** 页
2. 填入标题与(可选)视觉风格
3. 粘贴剧本,支持标记剧集:
   ```
   第1集 雨夜
   咖啡店里,林夏正看着窗外发呆...

   第2集 重逢
   ...
   ```
   标记识别规则:中文 `第N集 / 第N回 / 第N章`,英文 `Episode N / EP N / Ep.N / Scene N`。
4. 点击 **仅解析剧集** 预览拆分结果(不调用 API)
5. 确认无误后点击 **创建任务**,跳转至任务详情
6. 在详情页点击 **Run all** 开始生成,或点击某一集下方的 **Generate this episode** 只生成单集

---

## 2. 生成流程

每个任务按以下顺序执行:

1. **Planning(规划)** — 拆解剧集,生成预览
2. **Characters(角色)** — 调用 `agnes-2.5-flash` 抽取所有角色,再调用 `agnes-image-2.1-flash` 为每个角色生成首帧图
3. **Props(道具)** — 同上,抽取并生成道具图
4. **Scenes(场景)** — 为每集按段落拆解为若干 Scene,为每个 Scene 生成首帧图
5. **Video(视频)** — 以所有 Scene 图片为 `input_reference`,调用 `agnes-video-v2.0` 合成 5 秒视频,并轮询 `agnesapi?video_id=...` 直至完成
6. **Done** — 完成

> 整个过程**通常需要 5–20 分钟**,取决于剧集数与 Scenes 数。

---

## 3. 队列与限流

- 服务内置单 worker 任务队列(`AGNES_STUDIO_MAX_CONCURRENCY`,默认 1)
- 每个图像/视频请求之间会**随机 sleep 1–2 秒**,避免触发免费档 RPM 限制
- 视频轮询间隔默认 5 秒,超时 15 分钟(`AGNES_STUDIO_VIDEO_TIMEOUT`)
- 任务会被持久化,服务重启后**进度不会丢失**
- 同一时间只能跑一个任务,新任务会自动排队;可在任务详情页 **取消** 当前任务

免费档速率参考(摘自官方文档):

| 资源 | 实际可执行 RPM |
|---|---|
| 文本(自由档) | 20 |
| 图像 1K(自由档) | 20 |
| 视频(自由档) | 1 |

---

## 4. 常见问题

### 4.1 测试连接失败,显示 `401 Unauthorized`

- API key 错误或已失效,请到 [platform.agnes-ai.com](https://platform.agnes-ai.com/) 重新生成
- 选了错误的服务路由(国际账号不能访问中国大陆服务)

### 4.2 视频生成一直 queued

- 正常现象,视频模型需要排队;在 RPM 紧张时可能持续 5–10 分钟
- 如果超过 15 分钟仍未完成,会自动超时失败,可在详情页重新 **Run all**

### 4.3 生成的图片在视频里无法引用

- 视频模型会从公网下载 `input_reference` 图片。如果你在本地运行,本地 `http://localhost:8080` 不可被 Agnes 后端访问
- 解决:把本服务部署到公网,或使用 ngrok / cloudflared 等隧道把端口暴露到公网,然后设置 `AGNES_STUDIO_PUBLIC_BASE_URL` 环境变量
- 例如:`AGNES_STUDIO_PUBLIC_BASE_URL=https://your-tunnel.example.com`

### 4.4 想清空所有数据

```bash
# 停止服务,然后:
rm -rf data/jobs.json data/assets/*
# data/.masterkey 和 data/.apikey 是加密密钥与 key,删除后需要重新输入 API key
```

### 4.5 移动端能用吗?

可以。前端响应式布局,iOS Safari / Android Chrome 均支持。毛玻璃效果在不支持 `backdrop-filter` 的浏览器上会退化为半透明背景。

---

## 5. 快捷参考

### 5.1 环境变量

| 变量 | 默认值 | 说明 |
|---|---|---|
| `AGNES_STUDIO_ADDR` | `:8080` | HTTP 监听地址 |
| `AGNES_STUDIO_DATA_DIR` | `./data` | 数据目录 |
| `AGNES_STUDIO_ROUTE` | `international` | 服务路由 |
| `AGNES_STUDIO_MAX_CONCURRENCY` | `1` | 同时运行的任务数 |
| `AGNES_STUDIO_VIDEO_POLL_INTERVAL` | `5s` | 视频轮询间隔 |
| `AGNES_STUDIO_VIDEO_TIMEOUT` | `15m` | 视频超时 |
| `AGNES_STUDIO_MASTER_KEY` | _空_ | 主加密密钥(base64 32 字节),留空自动生成 |
| `AGNES_STUDIO_PUBLIC_BASE_URL` | _空_ | 公网 base URL,用于视频引用图 |

### 5.2 API 速查

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/api/health` | 健康检查 |
| `GET` | `/api/queue` | 队列状态 |
| `GET` | `/api/settings` | 读取设置 |
| `POST` | `/api/settings/key` | 保存 API key(`{api_key:"..."}`) |
| `DELETE` | `/api/settings/key` | 清除 API key |
| `POST` | `/api/settings/key/test` | 测试 API key |
| `POST` | `/api/settings/route` | 切换服务路由 |
| `GET` | `/api/jobs` | 列出所有任务 |
| `POST` | `/api/jobs` | 创建任务(`{title, script, style}`) |
| `GET` | `/api/jobs/{id}` | 任务详情 |
| `DELETE` | `/api/jobs/{id}` | 删除任务 |
| `POST` | `/api/jobs/{id}/run` | 运行全部剧集 |
| `POST` | `/api/jobs/{id}/episodes/{n}/run` | 只运行第 n 集 |
| `POST` | `/api/jobs/{id}/cancel` | 取消运行 |
| `GET` | `/api/assets/{name}` | 下载生成的图片/视频 |
| `GET` | `/api/docs` | 获取本用户文档(Markdown) |

### 5.3 支持的剧集标记

| 写法 | 示例 |
|---|---|
| 中文 | `第1集`、`第二回`、`第3章` |
| 英文 | `Episode 2`、`Ep 3`、`Ep.4`、`Scene 5` |
| Markdown 标题 | `## 第1集`、`### Episode 2` |

> 检测不到任何标记时,整段剧本作为**单集**处理。

---

## 6. 联系与反馈

- 仓库: [github.com/webzol/agnesai](https://github.com/webzol/agnesai)
- 官方文档: [agnes-ai.com/docs/overview](https://www.agnes-ai.com/zh-Hans/docs/overview)
- 官方 Skills 仓库: [github.com/AgnesAI-Labs/skills](https://github.com/AgnesAI-Labs/skills/tree/main)
- 平台(申请 API key): [platform.agnes-ai.com](https://platform.agnes-ai.com/)