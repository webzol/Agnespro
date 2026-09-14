# Agnes AI Studio

> 把一段剧本变成短剧的工作台。后端 Go + 前端 iOS 26 液态玻璃。

## 它能做什么

- **粘贴剧本** → 自动拆出剧集、抽取角色与道具、生成首帧图、合成 5 秒短剧
- **支持选择**:一次生成全部剧集,或指定只生成某一集
- **任务排队**:提交后自动排队,可以在 UI 看实时进度
- **数据本地化**:任务、API key、生成资产全部保存在本机
- **官方集成**:基于 [AgnesAI](https://www.agnes-ai.com/) 三个模型(`agnes-2.5-flash` / `agnes-image-2.1-flash` / `agnes-video-v2.0`)

## 快速开始

```bash
cd backend
go build -o server.exe ./cmd/server
./server.exe
```

打开 http://localhost:8080,到 **Settings** 填入 API key,然后 **New Job** 粘贴剧本。

详细使用:[`docs/USER_GUIDE.md`](docs/USER_GUIDE.md)
架构与扩展:[`docs/DEVELOPMENT.md`](docs/DEVELOPMENT.md)

## 截图

UI 截图请参见 `docs/USER_GUIDE.md`。

## 项目结构

```
.
├── backend/             # Go 后端 + 内嵌前端
│   ├── cmd/server/      # 入口
│   ├── internal/        # 业务包
│   └── web/             # 前端(液态玻璃)
├── docs/                # 用户文档 + 开发文档
├── AGENTS.md            # 与 AI 协作的铁律
└── README.md            # 本文件
```

## 致谢

- 官方 Skills 仓库: https://github.com/AgnesAI-Labs/skills
- 官方文档: https://www.agnes-ai.com/zh-Hans/docs/overview
- 平台: https://platform.agnes-ai.com/