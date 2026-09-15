# DEV_NOTES — Agnes AI Studio

> 给后续大模型接手用的开发日志。
> 关键决策 / 踩坑 / 待验证都集中在这里，避免重复踩坑。

## 1. TD 的硬性铁律（每次开工先确认）

- **每次开发完成必须 commit + push 到 GitHub**（不接受"只 commit 不 push"）。
- **绝不提交任何密钥 / key / token**。包括：
  - 代码里的 API key、access token、private key
  - `.env`、`.pem`、`.key`、`*.pem`、`*.key`、`secrets.yaml`、`.masterkey`、`.apikey`
  - 生成资产 `data/assets/*`、运行时数据 `data/*`
- 仓库已有 `.gitignore` 防护以上大部分文件，commit 前再自查一次。

## 2. 本地启动流程（已验证可行）

```
# 1. 设置环境变量（解决沙箱写权限问题，见第 3 节）
$env:GOTELEMETRY = "off"
$env:GOCACHE     = "E:\AI\.gocache"
$env:GOTMPDIR    = "E:\AI\.gotmp"

# 2. 编译
cd E:\AI\backend
go build -o server.exe ./cmd/server

# 3. 启动（监听 :8080）
.\server.exe -addr :8080 -data ./data -route international
```

启动后浏览器打开 `http://localhost:8080`，先到 **Settings** 填 Agnes API key，再 **New Job** 粘贴剧本。

## 3. 沙箱环境踩坑（Windows + 受限用户）

| 问题 | 现象 | 解决方案 |
|---|---|---|
| `git clone https://github.com/...` 报 `SEC_E_NO_CREDENTIALS` | 沙箱内 git 走 schannel 认证失败 | 升级到 sandbox 外面执行（`require_escalated`） |
| 私有仓库 401 Unauthorized | 没有 PAT / PAT 失效 | 用 classic PAT + 用户名认证（详见第 4 节） |
| GitHub HTTPS basic auth 报 "Password authentication is not supported" | PAT 写错 / 失效 / 复制截断 | 先 `Invoke-WebRequest -Headers @{Authorization="Bearer $tok"} https://api.github.com/user` 验证 token 状态 |
| `~/.gitconfig` Permission denied | 沙箱用户对 home 目录下的 `.gitconfig` 没写权限 | 把 credential helper 写到仓库本地 `.git/config`，绕过 global config |
| `~/.git-credentials` 没法创建 | 同上 | 用 `git config credential.helper "store --file=E:/AI/.git-credentials"` 把凭据文件重定向到仓库根 |
| `AppData\Local\go-build` 写入 Access is denied | Go 默认 build cache 目录没写权限 | 设置 `$env:GOCACHE = "E:\AI\.gocache"`、`$env:GOTMPDIR = "E:\AI\.gotmp"` |
| Go telemetry `upload.token` 写入失败 | 沙箱用户对 `AppData\Roaming\go\telemetry` 无写权限 | 设置 `$env:GOTELEMETRY = "off"`（关闭 telemetry） |
| `AppData\Roaming\go\telemetry\...` Permission denied | 同上 | 同上 |

## 4. Git 认证（已配好，不要改）

仓库 `E:\AI\.git\config` 已配置：

```
[remote "origin"]
    url = https://github.com/webzol/agnesai.git
    fetch = +refs/heads/*:refs/remotes/origin/*
[credential]
    helper = store --file=E:/AI/.git-credentials
```

- Token 存在 `E:\AI\.git-credentials`（明文，但已被 `.gitignore` 排除，不会被 commit）。
- **绝对不要把 token 写回 `remote.origin.url` 或 commit message 或任何代码文件**。
- 这个 PAT 在 30 天后（2026-10-15）会过期 —— 届时在 https://github.com/settings/tokens 重新生成 classic PAT + `repo` scope，更新 `.git-credentials` 即可。

## 5. AGENTS.md 铁律摘要（仓库自带的）

- 提交前自查清单已在仓库 `AGENTS.md` 里。
- Go 后端 `go build` 必须通过；前端静态资源语法必须有效。

## 6. 待验证 / 未做事项

- [ ] 浏览器实际访问 `http://localhost:8080`、填 API key、跑通 New Job → 视频生成全流程（沙箱内没法开浏览器，需要 TD 自己点一下）
- [ ] PAT 过期提醒（2026-10-15）—— 到期前主动重发
- [ ] 是否切到 SSH（`git@github.com:webzol/agnesai.git`）？当前 HTTPS + PAT 已稳定，暂无必要
- [ ] 是否升级到 fine-grained token？classic + `repo` scope 已满足需求，fine-grained 更精细但配置繁琐

## 7. 变更记录

- 2026-09-15：首次拉取仓库 + 本地启动成功（PID 见 `data/server.pid`，如果存在）。处理 PAT 泄露：`.git/config` URL 抹掉 token，凭据重定向到仓库本地文件。`.gitignore` 补 `.gocache/`、`.gotmp/`、`.git-credentials`。

## 8. 风格库 v2 + AI 模型选择 + 视频比例 + 剧集数（2026-09-15）

### 新增 endpoints
- `GET /api/styles/library` — 一次返回视觉画风库(24)+ 叙事题材(8)
  + 比例(6) + 剧集预设(6) 全部数据
- `GET /api/settings/models` — 用已设置 key 调 Agnes /v1/models,
  按 ID 前缀(`-image-` / `-video-`)分类返回 chat/image/video 列表

### 数据模型扩展（types.go）
Job 加 7 个字段:
- AspectRatio / EpisodeCount / GenreID / GenreName
- ScriptModel / ImageModel / VideoModel
Settings 同步加 3 个默认 model 字段

### 三级 model fallback（pipeline.go）
`pickModel(jobModel, cfgModel, fallback)`:
1. Job 字段（用户在 NewJob 选的）
2. Pipeline.Cfg（用户在 Settings 设的默认）
3. 代码硬编码（agnes-2.5-flash / agnes-image-2.1-flash / agnes-video-v2.0）

修改 pipeline 启动时从 settings 注入默认 model;
通过 POST /api/settings 改 model 后实时更新 Pipeline.Cfg（无需重启）。

### 比例 size 映射
- 6 个比例 → image size + video size 两套
- Agnes 视频原生支持 720x1280 / 1280x720 / 1024x1024
- 映射表在 internal/types/aspect.go,跨 api/pipeline 包共享

### 关键决策
- `decodeJSON` 默认 DisallowUnknownFields(strict)，
  generateScriptReq 必须显式接收所有前端传的字段，否则 400
- 风格库预览图：用 API key 批量调 image API 生成 24 张(每张 1-2.5 MB)
- API key 仅在内存中使用，data/.apikey 已加密存盘;不写入任何文件/commit

### 关键踩坑
1. **json DisallowUnknownFields**：前端加了新字段 (aspect_ratio) 但后端
   struct 没声明 → 直接 400。修法是给 generateScriptReq 也声明这些字段
   (即使后端暂不使用,只是收下避免报错)
2. **跨包 helper 调用**：resolveAspectSize 原本在 api 包,
   pipeline 包也要用 → 移到 types 公共包
3. **遗漏的 Size 参数**：原本 pipeline.go 调 CreateVideo 时**没传 Size**,
   全部视频都用默认 1280x720;这次顺手补上了
4. **types/aspect.go 初次创建失败**：大 Python 脚本里其中一个
   replace 没生效(类型残留),编译报"imported and not used";
   单独跑一个小 Python 创建后立刻成功

### 前端改动
- index.html: 替换 ai-genre 10→8 项,两种模式都加 chip-btn x 3 + 模型下拉 x 3
- app.js: loadStyleLibrary/loadModels + modal 控制 + 字段透传
- style.css: +6KB 样式(chip-btn/modal/style-library-grid)
- web/assets/img/library/: 24 张 AI 生成预览图(全部能 serve 200)

### 待验证（用户做）
- 浏览器实际体验风格库 modal UI
- 跑一个完整端到端 Job 生成(视频生成耗时长+token 多,沙箱里未跑)

## 9. 后台配置中心 (2026-09-15)

### 需求
"设置密钥"从前端 Settings 页面挪到独立的后台,新增后端配置 API。

### 新增 endpoints (internal/api/admin.go)
- `GET  /admin/api/config`         — 一站式 GET (key 状态/route/默认模型/可用模型列表)
- `POST /admin/api/config`         — 一站式 POST (子集更新)
- `POST /admin/api/config/test-key` — 测试连接

### 关键设计
- 独立路径前缀 /admin/api/,与 /api/settings/* 解耦
- GET 一次性聚合 + 实时拉取可用模型(10s timeout)
- POST 支持 partial update,空字段不动
- route 切换需重启(restart_required=true,前端用 toast 提示)
- model 切换立即生效(直接改 Pipeline.Cfg)
- clear_api_key:true 提供清除 key 能力

### 前端 admin view
- 导航栏加 "⚙ 后台" tab
- 4 个 section: API Key / API Route / 默认模型 / 可用模型列表
- route radio 切换即自动保存(无需点按钮)
- 模型 chips 点击自动填入对应 select(flash 反馈)

### 复用与边界
- 复用 testKeyResp / writeJSON / coalesceStr 等既有 helper
- admin handler 不重复造轮子:setAPIKey / SaveSettings 逻辑仍由 settings.go 提供

### 待验证(用户)
- 浏览器实际体验 admin view UI
- 实际切换 route + 重启服务验证生效

