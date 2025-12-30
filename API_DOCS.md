# Mystia Voice Backend 接口文档

本文档详细说明了语音后端服务的 gRPC 接口及使用方法，供前端开发参考。

## 1. 基本信息

- **协议**: gRPC
- **通信方式**: Unary (简单请求-响应)
- **认证方式**: JWT (Bearer Token)
- **LiveKit 服务**: 后端提供 Access Token，前端需配合 [LiveKit SDK](https://docs.livekit.io/realtime/client-sdk/) 使用。

---

## 2. 网络与端口配置 (Networking)

为了让前端正常连接，你需要确保服务端开放以下端口：

### 2.1 Mystia 业务后端 (本服务)
- **50051 (TCP)**: gRPC 服务默认端口。前端通过此端口进行登录、获取频道列表等业务操作。
- **启动参数**:
    - `-port <number>`: 自定义 Mystia gRPC 端口。
    - `-lk-port <number>`: 自定义连接 LiveKit 的 TCP 端口（会覆盖环境变量中的 Host 端口）。

### 2.2 LiveKit 媒体服务器
LiveKit 需要开放多个端口以保证 WebRTC 通信的顺畅：
- **7880 (TCP)**: HTTP API 和 WebSocket 信令端口。前端连接 LiveKit 时主要访问此端口。
- **443 (UDP)** 或 **50000-60000 (UDP)**: WebRTC 媒体传输端口（用于传输语音流）。
- **3478 (UDP)**: STUN 服务（可选，用于内网穿透）。

> [!IMPORTANT]
> 如果是在云服务器上运行，请务必在安全组/防火墙中开放上述 TCP 和 UDP 端口。

---

## 3. 身份认证与用户服务 (AuthService)

所有的请求（除 `Register` 和 `Login` 外）都必须在请求头中携带身份令牌。

### 2.1 用户注册 (Register)
- **方法**: `rpc Register(RegisterRequest) returns (RegisterResponse)`
- **说明**: 注册新用户。系统第一个注册的用户将自动获得 `SUPER_ADMIN` 角色，后续注册的用户默认为普通用户。
- **请求参数**:
    - `username`: 用户名 (登录凭据)
    - `password`: 密码
    - `nickname`: 显示昵称

### 2.2 用户登录 (Login)
- **方法**: `rpc Login(LoginRequest) returns (LoginResponse)`
- **说明**: 登录并获取 JWT Token。
- **返回参数**:
    - `token`: JWT 字符串 (后续请求需在 Header 中带上 `Authorization: Bearer <token>`)
    - `user`: 用户详情（包含 ID、昵称、角色）

---

## 3. 频道服务 (ChannelService)

用于获取语音频道列表及加入语音房间。

### 3.1 获取频道列表 (ListChannels)
- **方法**: `rpc ListChannels(ListChannelsRequest) returns (ListChannelsResponse)`
- **说明**: 获取当前所有可用的语音频道。

### 3.2 加入频道 (JoinChannel)
- **方法**: `rpc JoinChannel(JoinChannelRequest) returns (JoinChannelResponse)`
- **说明**: 请求加入指定频道，由于是语音服务器，返回的是 LiveKit 连接所需的 Token。
- **返回参数**:
    - `token`: LiveKit Access Token
    - `url`: LiveKit 服务器地址 (例如 `http://localhost:7880`)
- **后续操作**: 前端调用 LiveKit SDK 的 `room.connect(url, token)` 即可进入语音通话。

### 3.3 创建频道 (CreateChannel) - **管理员专用**
- **方法**: `rpc CreateChannel(CreateChannelRequest) returns (Channel)`
- **说明**: 只有具有 `ADMIN` 或 `SUPER_ADMIN` 角色的用户可以调用。

---

## 4. 管理员服务 (AdminService)

仅限 `ADMIN` 或 `SUPER_ADMIN` 角色调用的接口。

### 4.0 角色与权限等级 (Role Hierarchy)

系统目前支持以下三种角色：
1. **SUPER_ADMIN (超级管理员)**：
   - 系统第一个注册的用户自动获得。
   - **最高权限**：拥有所有管理权限。
   - **不可变更**：任何管理员（包括他自己）都无法取消其权限或将其降级。
2. **ADMIN (管理员)**：
   - 由其他管理员提升而来。
   - 拥有日常管理权限（创建频道、管理普通用户）。
   - **受限**：无法修改 `SUPER_ADMIN` 的状态。
3. **USER (普通用户)**：
   - 默认角色，仅能进行语音聊天和查看频道。

### 4.1 用户管理列表 (ListUsers)
- **方法**: `rpc ListUsers(ListUsersRequest) returns (ListUsersResponse)`
- **说明**: 获取系统内所有用户的列表。`ADMIN` 和 `SUPER_ADMIN` 均可调用。

### 4.2 设置/取消管理员 (SetAdminStatus)
- **方法**: `rpc SetAdminStatus(SetAdminStatusRequest) returns (SetAdminStatusResponse)`
- **说明**: 修改指定用户的角色。
- **限制**: 无法修改目标用户为 `SUPER_ADMIN` 的状态。如果尝试修改 `SUPER_ADMIN`，将返回 `PERMISSION_DENIED` 错误。

---

## 5. 错误码说明

后端遵循标准的 gRPC 状态码：
- `OK (0)`: 成功
- `UNAUTHENTICATED (16)`: 未登录或 Token 已过期
- `PERMISSION_DENIED (7)`: 权限不足（如非管理员尝试创建频道）
- `ALREADY_EXISTS (6)`: 用户名已存在
- `INVALID_ARGUMENT (3)`: 请求参数缺失或格式错误
- `INTERNAL (13)`: 服务器内部错误
