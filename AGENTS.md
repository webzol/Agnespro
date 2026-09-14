# AGENTS.md — 项目协作规则

> 本文件由 Codex 自动维护,记录与本项目协作时必须遵守的硬性规则。
> 任何为本项目工作的 AI 助手(包括 Codex)开工前请先阅读本文件。

## 铁律:每次开发完成必须推送到 GitHub

**所有开发任务结束后,必须执行以下推送流程,不接受"只 commit 不 push"或"本地留底"**:

1. `git status` 确认工作区干净或变更符合预期
2. `git add` 暂存所有相关变更
3. `git commit` 提交(信息要写清楚做了什么)
4. `git push` 推送到 `origin` 的对应分支
5. 验证推送成功(`git log origin/<branch> --oneline | head -3`)

## 推送凭据

- 优先使用 SSH: `git@github.com:webzol/agnesai.git`
- SSH 不可用时回退到 HTTPS + Personal Access Token(已配置在 Windows Credential Manager)
- **禁止**把任何 token 写入代码、commit message、或 `git config`

## 仓库信息

- Owner: `webzol`
- Repo: `agnesai`
- 默认分支: `main`
- 远程: `origin`
