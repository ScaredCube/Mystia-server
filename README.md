# Mystia Voice Backend

Mystia 是一个类似 Discord/TeamSpeak 的语音服务器后端，专为高性能语音通话设计。

## 特性

- **高性能语音**: 基于 [LiveKit](https://livekit.io/) 协议实现，支持大规模低延迟语音通话。
- **自动拉起服务**: 启动 Mystia 时会自动启动同目录下的 `livekit-server.exe`，简化部署流程。
- **gRPC 接口**: 前后端采用 gRPC 协议交互，提供强类型的接口定义和高效的通信。
- **账号系统**: 完整的用户注册与登录功能，密码使用 Bcrypt 加密存储。
- **权限管理**:
  - **角色模型**: 支持 `USER` 和 `ADMIN` 角色。
  - **自动管理员**: 系统首个注册用户自动获得超级管理员权限。
  - **动态管理**: 管理员可以创建频道、列出用户并授予/撤销其他用户的管理员身份。
- **快速部署**: 使用 SQLite 作为轻量级数据库，无需复杂的环境搭建。
- **灵活配置**: 支持通过环境变量及命令行参数自定义运行端口。

---

## 技术栈

- **语言**: Go (Golang)
- **协议**: gRPC, WebRTC (via LiveKit)
- **数据库**: SQLite
- **认证**: JWT (JSON Web Token)
- **加密**: Bcrypt

---

## 快速开始

### 1. 环境准备
确保已安装：
- [Go](https://go.dev/)
- [LiveKit Server](https://github.com/livekit/livekit/)
在release界面下载打包好的livekit-server放在项目根目录下（与编译好的mystia-server同一目录）


### 2. Proto 文件转换
如果你修改了 `proto/voice.proto`，需要重新生成 Go 代码。

**转换命令**:
在项目根目录下运行：
```bash
protoc --go_out=. --go-grpc_out=. proto/voice.proto
```

### 3. 编译项目
在项目根目录下运行：
```bash
go build -o mystia-server.exe cmd/server/main.go
```



### 4. 配置与运行
你可以通过环境变量或命令行参数来运行服务。

#### 环境变量配置 (可选)
| 变量名 | 默认值 | 说明 |
| :--- | :--- | :--- |
| `DB_PATH` | `voice.db` | SQLite 数据库文件路径 |
| `JWT_SECRET` | `super-secret-key` | JWT 签名密钥 |
| `LIVEKIT_API_KEY` | `devkey` | LiveKit API Key |
| `LIVEKIT_API_SECRET` | `secret` | LiveKit API Secret |
| `LIVEKIT_HOST` | `http://localhost:7880` | LiveKit 服务器地址 |
| `PORT` | `50051` | Mystia 服务端口 |

#### 直接运行 (推荐)
```bash
# 使用默认配置启动
.\mystia-server.exe

# 自定义端口启动
.\mystia-server.exe -port 60061 -lk-port 8880
```

---

## 文档

- **接口定义**: [proto/voice.proto](./proto/voice.proto)
- **开发者文档 (中)**: [API_DOCS.md](./API_DOCS.md) - 包含详细的接口参数、认证流程及网络配置说明。

---

## 项目结构

```text
.
├── cmd/
│   ├── server/           # 服务端入口 (main.go)
│   └── test_client/      # 测试客户端 (用于验证功能)
├── internal/
│   ├── db/               # 数据库操作 (SQLite)
│   ├── livekit/          # LiveKit Token 生成逻辑
│   └── service/          # gRPC 服务实现及中间件
├── proto/                # gRPC Protobuf 定义文件
├── API_DOCS.md           # 面向前端的接口说明文档
└── README.md             # 项目说明
```

## 贡献

欢迎提交 Issue 或 Pull Request 来完善这个项目。

## 许可证

本项目采用 MIT 许可证。
