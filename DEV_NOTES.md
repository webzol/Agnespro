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
