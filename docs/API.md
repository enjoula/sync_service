# API 文档

## 目录

- [基本信息](#基本信息)
- [统一响应格式](#统一响应格式)
- [错误码](#错误码)
- [公开API（无需认证）](#公开api无需认证)
  - [用户注册](#1-用户注册)
  - [用户登录](#2-用户登录)
  - [健康检查](#3-健康检查)
  - [IP信息查询](#4-ip信息查询)
- [需要认证的API](#需要认证的api)
  - [获取当前用户信息](#1-获取当前用户信息)
- [监控API](#监控api)
- [使用示例](#使用示例)
- [Token管理](#token管理)
- [追踪ID](#追踪idtrace-id)
- [安全性说明](#安全性说明)
- [性能监控](#性能监控)
- [数据库设计](#数据库设计)
- [常见问题](#常见问题)

## 基本信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **响应格式**: 所有响应均为JSON格式

## 统一响应格式

所有响应均采用统一的 JSON 格式，同时在响应头中包含 `X-Trace-ID` 用于请求追踪。

**响应体**:
```json
{
  "code": 0,              // 0表示成功，其他值表示错误
  "message": "Success",   // 响应消息
  "data": {},             // 响应数据
  "trace_id": "uuid..."   // 请求追踪ID
}
```

**响应头**:
- `X-Trace-ID`: 请求追踪ID（与响应体中的 trace_id 字段相同）

## 错误码

| 错误码 | 说明 |
|--------|------|
| 0      | 成功 |
| 400    | 请求参数错误 |
| 401    | 未授权（需要登录或token无效） |
| 409    | 资源冲突（如用户已存在） |
| 500    | 服务器内部错误 |

## 公开API（无需认证）

### 1. 用户注册

**接口**: `POST /user/register`

**描述**: 注册新用户，注册成功后自动登录并返回token

**请求体**:
```json
{
  "username": "testuser",
  "password": "password123",
  "device": "iOS-iPhone15Pro"
}
```

**参数说明**:
- `username` (string, 必需): 用户名
- `password` (string, 必需): 密码
- `device` (string, 必需): 设备信息，格式：平台-设备型号，示例：
  - `iOS-iPhone15Pro`
  - `android-SamsungS24`
  - `TV-XiaomiTV5`
  - `PC-Windows11`

**响应示例**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "id": 244770533047660544,
    "username": "testuser",
    "nickname": "喜羊羊",
    "avatar": "https://img.pngsucai.com/00/53/87/d9bbfc94bea22c16.webp",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "trace_id": "abc-123-def-456"
}
```

**错误响应示例**:

1. 参数错误（用户名或密码为空）:
```json
{
  "code": 400,
  "message": "username/password required",
  "data": null,
  "trace_id": "abc-123-def-456"
}
```

2. 用户名已存在:
```json
{
  "code": 409,
  "message": "用户名重复",
  "data": null,
  "trace_id": "abc-123-def-456"
}
```

### 2. 用户登录

**接口**: `POST /user/login`

**描述**: 用户登录，返回JWT token

**请求体**:
```json
{
  "username": "testuser",
  "password": "password123",
  "device": "iOS-iPhone15Pro"
}
```

**参数说明**:
- `username` (string, 必需): 用户名
- `password` (string, 必需): 密码
- `device` (string, 必需): 设备信息，格式：平台-设备型号，示例：
  - `iOS-iPhone15Pro`
  - `android-SamsungS24`
  - `TV-XiaomiTV5`
  - `PC-Windows11`

**响应示例**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "id": 244770533047660544,
    "username": "testuser",
    "nickname": "喜羊羊",
    "avatar": "https://img.pngsucai.com/00/53/87/d9bbfc94bea22c16.webp",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "trace_id": "abc-123-def-456"
}
```

**错误响应示例**:

1. 用户名或密码错误:
```json
{
  "code": 401,
  "message": "invalid credentials",
  "data": null,
  "trace_id": "abc-123-def-456"
}
```

2. 请求参数无效:
```json
{
  "code": 400,
  "message": "invalid request",
  "data": null,
  "trace_id": "abc-123-def-456"
}
```

### 3. 健康检查

**接口**: `GET /ping`

**描述**: 检查服务是否正常运行

**响应示例**:
```json
{
  "code": 0,
  "message": "pong",
  "data": {
    "time": "ok"
  },
  "trace_id": "abc-123-def-456"
}
```

### 4. IP信息查询

**接口**: `GET /ip-info`

**描述**: 获取客户端IP信息（用于调试）

**响应示例**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "ip_info": {
      "real_ip": "192.168.1.100",
      "is_private": "true",
      "x_real_ip": "192.168.1.100",
      "x_forwarded_for": "",
      "cf_connecting_ip": "",
      "remote_addr": "192.168.1.100:54321"
    },
    "all_headers": {
      "User-Agent": "Mozilla/5.0...",
      "Accept": "application/json"
    },
    "remote_addr": "192.168.1.100:54321",
    "client_ip": "192.168.1.100"
  },
  "trace_id": "abc-123-def-456"
}
```

## 需要认证的API

### 认证方式

在请求头中添加JWT token：

```
Authorization: Bearer <your-token-here>
```

### 1. 获取当前用户信息

**接口**: `GET /user/me`

**描述**: 获取当前登录用户信息

**请求头**:
- `Authorization`: `Bearer <token>` (必需)

**响应示例**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "user": "testuser"
  },
  "trace_id": "abc-123-def-456"
}
```

**错误响应示例**:

1. 缺少认证头:
```json
{
  "code": 401,
  "message": "missing authorization header",
  "data": null,
  "trace_id": "abc-123-def-456"
}
```

2. 认证头格式错误（缺少 Bearer 前缀）:
```json
{
  "code": 401,
  "message": "invalid authorization header format",
  "data": null,
  "trace_id": "abc-123-def-456"
}
```

3. Token已过期或失效:
```json
{
  "code": 401,
  "message": "token已过期",
  "data": null,
  "trace_id": "abc-123-def-456"
}
```

4. Token无效（签名错误或格式错误）:
```json
{
  "code": 401,
  "message": "invalid token: token signature is invalid",
  "data": null,
  "trace_id": "abc-123-def-456"
}
```

5. Token验证失败（数据库错误）:
```json
{
  "code": 500,
  "message": "token验证失败",
  "data": null,
  "trace_id": "abc-123-def-456"
}
```

## 监控API

### Prometheus指标

**接口**: `GET /metrics`

**描述**: Prometheus监控指标

**响应格式**: Prometheus文本格式

**示例**:
```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="POST",endpoint="/login",status="200"} 42
```

## 使用示例

```bash
# 注册
curl -X POST http://localhost:8080/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pwd123","device":"PC-MacBookPro"}'

# 登录
curl -X POST http://localhost:8080/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pwd123","device":"PC-MacBookPro"}'

# 获取用户信息（需要先登录获取token）
curl http://localhost:8080/user/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# 健康检查
curl http://localhost:8080/ping

# IP信息
curl http://localhost:8080/ip-info
```

## Token管理

### Token结构

JWT Token 采用标准的 JWT 格式，包含以下声明（Claims）：

- `username`: 用户名
- `jti` (JWT ID): 唯一标识符（UUID），确保每个 token 唯一
- `iat` (Issued At): 签发时间戳
- `exp` (Expiration): 过期时间戳

**示例 Token Payload**:
```json
{
  "username": "testuser",
  "jti": "a3f2c1b5-4d3e-2f1a-9b8c-7d6e5f4a3b2c",
  "iat": 1699660800,
  "exp": 1699747200
}
```

### Token生命周期

- **有效期**: 365天（从签发时间开始计算）
- **刷新**: Token 过期后需要重新登录获取新token
- **多设备**: 最多3个活跃token，超过后最早的 token 会被设置为失效（is_active=0）
- **存储**: Token 在数据库中记录，包含设备信息、IP地址、活跃状态等

### Token验证流程

1. **验证顺序**: 先检查数据库中的 `is_active` 状态，再验证 token 签名和过期时间
2. **数据库检查**: 确保 token 存在且 `is_active=1`
3. **签名验证**: 使用 JWT 密钥验证 token 签名
4. **过期检查**: 验证 token 是否在有效期内

### Token失效场景

1. **自然过期**: Token 超过365天有效期
2. **多设备限制**: 用户登录新设备超过3个活跃 token 限制时，最早的 token 被标记为失效
3. **主动登出**: 用户主动登出时，token 的 `is_active` 被设置为 0（待实现）
4. **数据库状态**: Token 在数据库中被标记为失效（`is_active=0`）

## 追踪ID（Trace ID）

每个请求都会生成唯一的追踪ID（UUID格式），用于日志关联和问题排查。

### 获取追踪ID

追踪ID会同时在响应体和响应头中返回：

1. **响应体**: 从响应JSON的 `trace_id` 字段获取
2. **响应头**: 从响应头的 `X-Trace-ID` 获取

**示例**:
```bash
# 使用 curl 查看响应头
curl -i http://localhost:8080/ping

# 响应头中包含
X-Trace-ID: 550e8400-e29b-41d4-a716-446655440000
```

### 使用追踪ID

在日志系统中搜索追踪ID，可以查看该请求的完整处理流程，包括：

- 请求接收时间
- 请求参数和头部信息
- 数据库查询记录
- 业务逻辑执行过程
- 错误堆栈信息（如有）
- 响应时间

### 中间件处理顺序

所有请求都会经过以下中间件处理（按顺序）：

1. **Trace 中间件**: 生成追踪ID并存入 context
2. **Recovery 中间件**: 捕获 panic 并记录日志
3. **Logger 中间件**: 记录请求日志（包含追踪ID）
4. **Prometheus 中间件**: 收集监控指标
5. **JWTAuth 中间件**: JWT 认证（仅需要认证的路由）

## 安全性说明

### JWT密钥管理

- **配置方式**: JWT 密钥通过 etcd 配置中心管理，支持热更新
- **默认密钥**: 开发环境使用默认密钥 `change-me-default`
- **生产环境**: 必须在 etcd 中配置强密钥

**配置示例**:
```bash
# 在 etcd 中设置 JWT 密钥
etcdctl put /video-service/secret '{
  "jwt_key": "your-strong-secret-key-here",
  "mysql_dsn": "..."
}'
```

### 密码安全

- 用户密码使用 bcrypt 算法加密存储
- 加密强度：bcrypt cost factor = 10（默认）
- 不存储明文密码，不可逆向解密

### IP获取策略

服务支持从多个来源获取真实客户端IP，优先级如下：

1. `CF-Connecting-IP` (Cloudflare)
2. `True-Client-IP`
3. `X-Real-IP`
4. `X-Forwarded-For`（取第一个IP）
5. `RemoteAddr`

**信任的代理网段**:
- `127.0.0.1` (本地回环)
- `10.0.0.0/8` (私有网络 A类)
- `172.16.0.0/12` (私有网络 B类)
- `192.168.0.0/16` (私有网络 C类)

### 请求安全建议

1. **HTTPS**: 生产环境必须使用 HTTPS 传输，防止 token 被窃取
2. **Token存储**: 客户端应安全存储 token（如使用 httpOnly cookie 或安全的本地存储）
3. **Token刷新**: Token 过期后需要重新登录，建议实现自动刷新机制
4. **错误处理**: 不要在客户端暴露详细的错误堆栈信息

## 限流（待实现）

当前版本暂未实现限流功能，未来版本将添加：

- 基于IP的限流
- 基于用户的限流
- 自适应限流
- Redis实现的分布式限流

## 版本管理

当前API版本：v1（未在URL中体现）

未来可能的版本管理方案：
- URL版本：`/v1/register`
- Header版本：`Accept: application/vnd.api.v1+json`

## 性能监控

### Prometheus指标

服务自动暴露以下监控指标：

- **HTTP请求指标**:
  - `http_requests_total`: 请求总数（按方法、路径、状态码分组）
  - `http_request_duration_seconds`: 请求耗时直方图
  - `http_request_size_bytes`: 请求大小
  - `http_response_size_bytes`: 响应大小

- **Go运行时指标**:
  - `go_goroutines`: 当前 goroutine 数量
  - `go_memstats_alloc_bytes`: 内存分配量
  - `go_gc_duration_seconds`: GC耗时

**访问指标**:
```bash
curl http://localhost:8080/metrics
```

### 日志系统

- **日志格式**: JSON格式，便于日志分析
- **日志级别**: Debug, Info, Warn, Error
- **日志输出**: 同时输出到控制台和文件（`./logs/app.log`）
- **日志内容**: 包含追踪ID、用户信息、IP地址、请求参数等

## 数据库设计

### 用户表（users）

主要字段：
- `id`: 用户ID（使用雪花算法生成）
- `username`: 用户名（唯一索引）
- `password`: 加密后的密码（bcrypt）
- `nickname`: 昵称（自动生成）
- `avatar`: 头像URL（自动生成）
- `created_at`: 创建时间
- `updated_at`: 更新时间

### Token表（user_tokens）

主要字段：
- `id`: 自增主键
- `user_id`: 用户ID（外键）
- `token`: JWT token字符串（唯一索引）
- `device`: 设备类型
- `ip_address`: 登录IP地址
- `is_active`: 是否活跃（0=失效, 1=活跃）
- `expires_at`: 过期时间
- `created_at`: 创建时间

## 常见问题

### Q1: Token过期后如何处理？

A: Token有效期为365天（1年），过期后客户端需要重新调用登录接口获取新的token。当前版本不支持刷新token机制，需要用户重新输入用户名和密码。

**建议**：
- 在token即将过期时（如剩余7天）提示用户
- 检测到401错误时，引导用户重新登录
- 使用本地存储保存用户名（需用户同意），便于快速重新登录
- 由于有效期较长，一般无需担心频繁过期的问题

### Q2: 如何处理多设备登录？

A: 系统支持最多3个设备同时登录。当用户在第4个设备登录时，最早的token会被自动失效。

**场景示例**：
1. 用户在设备A登录 → 活跃token数: 1
2. 用户在设备B登录 → 活跃token数: 2
3. 用户在设备C登录 → 活跃token数: 3
4. 用户在设备D登录 → 设备A的token失效，活跃token数: 3

**提示**：
- 被踢下线的设备会收到401错误
- 可以通过查看`ip_address`和`device`字段识别不同设备

### Q3: 为什么我的请求返回401错误？

可能的原因：

1. **未提供Authorization头**
   ```json
   {"code": 401, "message": "missing authorization header"}
   ```
   解决：在请求头中添加 `Authorization: Bearer <token>`

2. **Authorization头格式错误**
   ```json
   {"code": 401, "message": "invalid authorization header format"}
   ```
   解决：确保格式为 `Bearer <token>`，注意Bearer后有空格

3. **Token已过期或失效**
   ```json
   {"code": 401, "message": "token已过期"}
   ```
   解决：重新登录获取新token

4. **Token签名无效**
   ```json
   {"code": 401, "message": "invalid token: ..."}
   ```
   解决：确认token未被篡改，重新登录

### Q4: 如何调试API请求？

推荐步骤：

1. **查看响应的trace_id**：
   ```bash
   curl -i http://localhost:8080/your-endpoint
   ```
   
2. **在日志中搜索trace_id**：
   ```bash
   grep "trace-id-here" logs/app.log
   ```

3. **使用IP信息接口验证网络配置**：
   ```bash
   curl http://localhost:8080/ip-info
   ```

4. **检查Prometheus指标**：
   ```bash
   curl http://localhost:8080/metrics
   ```

### Q5: 注册时昵称和头像是如何生成的？

- **昵称**：从预定义的昵称列表中随机选择（如"喜羊羊"、"美羊羊"等）
- **头像**：从预定义的头像URL列表中随机选择

如需自定义昵称和头像，需要调用用户信息更新接口（待实现）。

### Q6: 如何在生产环境部署？

关键步骤：

1. **配置JWT密钥**：
   ```bash
   etcdctl put /video-service/secret '{
     "jwt_key": "your-production-strong-key",
     "mysql_dsn": "production-db-connection"
   }'
   ```

2. **使用HTTPS**：配置Nginx或其他反向代理，启用TLS

3. **配置信任的代理**：确保正确获取客户端真实IP

4. **监控配置**：
   - 配置Prometheus采集
   - 配置日志收集系统
   - 设置告警规则

5. **数据库优化**：
   - 创建必要的索引
   - 配置数据库连接池
   - 启用慢查询日志

### Q7: 如何限制用户密码强度？

当前版本仅要求用户名和密码不为空。建议在客户端实施密码强度检查：

**推荐规则**：
- 最小长度：8位
- 包含大小写字母、数字和特殊字符
- 不包含用户名
- 不是常见弱密码

**服务端增强**（未来版本）：
- 密码强度验证
- 密码历史检查
- 防止暴力破解

### Q8: API响应时间慢怎么办？

排查步骤：

1. **查看Prometheus指标**：
   - `http_request_duration_seconds`：查看请求耗时分布
   - 识别慢接口

2. **分析日志**：
   - 搜索耗时超过阈值的请求
   - 查看数据库查询时间

3. **优化措施**：
   - 添加数据库索引
   - 启用Redis缓存
   - 优化查询语句
   - 增加数据库连接池大小

4. **检查基础设施**：
   - 网络延迟
   - 数据库性能
   - Redis性能

### Q9: 如何实现用户登出功能？

当前版本可以通过数据库操作实现：

```sql
-- 将用户的token标记为失效
UPDATE user_tokens 
SET is_active = 0 
WHERE token = 'your-token-here';
```

**建议**：实现一个登出接口（待开发）：
- `POST /logout`
- 需要认证
- 将当前token设置为失效

### Q10: 支持哪些数据库？

当前版本使用 MySQL，理论上支持 GORM 兼容的所有数据库：

- ✅ MySQL 5.7+
- ✅ MariaDB 10.3+
- ⚠️ PostgreSQL（需修改部分SQL）
- ⚠️ SQLite（仅用于开发测试）

**注意**：生产环境强烈推荐使用 MySQL 或 PostgreSQL。

---

## 更新日志

### v1.0.0 (当前版本)

**功能**：
- ✅ 用户注册、登录
- ✅ JWT认证
- ✅ 多设备管理（最多3个）
- ✅ 请求追踪（Trace ID）
- ✅ Prometheus监控
- ✅ 结构化日志
- ✅ IP获取和记录

**待开发**：
- ⏳ 用户登出接口
- ⏳ 用户信息修改
- ⏳ Token刷新机制
- ⏳ 限流功能
- ⏳ 视频相关功能

---

**文档版本**: v1.4  
**最后更新**: 2024-11-12  
**维护者**: Video Service Team

### 更新内容 (v1.4)
- ✅ 精简使用示例：移除JavaScript和Python示例，仅保留curl示例

### 更新内容 (v1.3)
- ✅ 修正 device 参数位置：从请求头改为请求体
- ✅ 修正 Token 有效期说明：从24小时更新为365天
- ✅ 更新所有代码示例以反映正确的API调用方式
- ✅ 补充 device 参数格式说明和示例

