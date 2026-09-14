# AGENTS.md — 项目协作规则

> 本文件记录与本项目协作时必须遵守的硬性规则。
> 任何为本项目工作的 AI 助手(包括 Codex)开工前请先阅读本文件。

## 铁律:每次开发完成必须推送到 GitHub

**所有开发任务结束后,必须执行以下推送流程,不接受"只 commit 不 push"或"本地留底"**:

1. `git status` 确认工作区干净或变更符合预期
2. `git add` 暂存所有相关变更
3. `git commit` 提交(信息要写清楚做了什么)
4. `git push` 推送到 `origin` 的对应分支
5. 验证推送成功(`git log origin/main --oneline | head -3`)

## 推送凭据

- 优先使用 SSH: `git@github.com:webzol/agnesai.git`
- SSH 不可用时回退到 HTTPS + Personal Access Token(已配置在 Windows Credential Manager)
- **禁止**把任何 token 写入代码、commit message 或 `git config`

## 项目范围

本项目(`agnesai`)是一个**视频生成工作台**:

- **后端**: Go 1.23+,使用标准库 `net/http`(Go 1.22+ 方法路由)+ `crypto/aes` 加密敏感数据 + JSON 文件持久化
- **前端**: 纯 HTML5 + CSS3 + 原生 JS(零构建),iOS 26 液态玻璃风格
- **集成**: Agnes AI 官方 OpenAI 兼容 API — 文本/图像/视频三个模型
- **功能**: 粘贴剧本 → 自动拆剧集 → 抽取角色/道具 → 生成首帧图 → 合成 5 秒短剧

详见:
- 用户文档: `docs/USER_GUIDE.md`(同时内置于前端 `/api/docs`)
- 开发文档: `docs/DEVELOPMENT.md`

## 仓库信息

- Owner: `webzol`
- Repo: `agnesai`
- 默认分支: `main`
- 远程: `origin`

## 提交前自查清单

- [ ] 没有把任何 `.env`、API key、`.masterkey`、`.apikey` 文件加入提交(已通过 `.gitignore` 防护)
- [ ] 没有把生成的素材(`data/assets/*`)加入提交
- [ ] 没有把构建产物(`backend/server.exe`)加入提交
- [ ] 代码改动有清晰 commit message(说明改了什么、为什么)
- [ ] Go 后端 `go build` 通过
- [ ] 前端静态资源(HTML/CSS/JS)语法有效

## 安全守则

- API key 一律走 Settings 页面录入,加密落盘(`data/.apikey` 0600)
- 主加密密钥在 `data/.masterkey` 0600,**丢失则需要重新输入 API key**
- 生产环境切勿把 `data/` 目录暴露到公网