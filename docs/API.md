# API 文档

## 基本信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **响应格式**: 所有响应均为JSON格式

## 统一响应格式

```json
{
  "code": 0,              // 0表示成功，其他值表示错误
  "message": "Success",   // 响应消息
  "data": {},             // 响应数据
  "trace_id": "uuid..."   // 请求追踪ID
}
```

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

**接口**: `POST /register`

**描述**: 注册新用户，注册成功后自动登录并返回token

**请求头**:
- `X-Device` (可选): 设备类型，如`web`、`tv`、`mobile`

**请求体**:
```json
{
  "username": "testuser",
  "password": "password123"
}
```

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

**错误响应**:
```json
{
  "code": 400,
  "message": "username/password required",
  "data": null,
  "trace_id": "abc-123-def-456"
}
```

### 2. 用户登录

**接口**: `POST /login`

**描述**: 用户登录，返回JWT token

**请求头**:
- `X-Device` (可选): 设备类型

**请求体**:
```json
{
  "username": "testuser",
  "password": "password123"
}
```

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

**错误响应**:
```json
{
  "code": 401,
  "message": "invalid credentials",
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

**未认证响应**:
```json
{
  "code": 401,
  "message": "missing authorization header",
  "data": null,
  "trace_id": "abc-123-def-456"
}
```

**Token过期响应**:
```json
{
  "code": 401,
  "message": "token已过期",
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

### curl

```bash
# 注册
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -H "X-Device: web" \
  -d '{"username":"alice","password":"pwd123"}'

# 登录
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pwd123"}'

# 获取用户信息（需要先登录获取token）
curl http://localhost:8080/user/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# 健康检查
curl http://localhost:8080/ping

# IP信息
curl http://localhost:8080/ip-info
```

### JavaScript (fetch)

```javascript
// 注册
const register = async (username, password) => {
  const response = await fetch('http://localhost:8080/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Device': 'web'
    },
    body: JSON.stringify({ username, password })
  });
  return await response.json();
};

// 登录
const login = async (username, password) => {
  const response = await fetch('http://localhost:8080/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ username, password })
  });
  return await response.json();
};

// 获取用户信息
const getMe = async (token) => {
  const response = await fetch('http://localhost:8080/user/me', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return await response.json();
};
```

### Python (requests)

```python
import requests

# 注册
def register(username, password):
    response = requests.post(
        'http://localhost:8080/register',
        headers={'X-Device': 'web'},
        json={'username': username, 'password': password}
    )
    return response.json()

# 登录
def login(username, password):
    response = requests.post(
        'http://localhost:8080/login',
        json={'username': username, 'password': password}
    )
    return response.json()

# 获取用户信息
def get_me(token):
    response = requests.get(
        'http://localhost:8080/user/me',
        headers={'Authorization': f'Bearer {token}'}
    )
    return response.json()
```

## Token管理

### Token生命周期

- **有效期**: 24小时
- **刷新**: 需要重新登录获取新token
- **多设备**: 最多3个活跃token，超过后旧设备会被踢下线

### Token失效场景

1. Token过期（24小时后）
2. 用户登录新设备，超过3个活跃token限制
3. 用户主动登出（设置is_active=0）
4. 服务器强制下线（待实现）

## 追踪ID

每个请求都会生成唯一的追踪ID，用于日志关联和问题排查。

### 获取追踪ID

1. 从响应JSON的`trace_id`字段获取
2. 从响应头的`X-Trace-ID`获取

### 使用追踪ID

在日志中搜索追踪ID，可以查看该请求的完整处理流程。

## 限流（待实现）

当前版本暂未实现限流功能，未来版本将添加：

- 基于IP的限流
- 基于用户的限流
- 自适应限流

## 版本管理

当前API版本：v1（未在URL中体现）

未来可能的版本管理方案：
- URL版本：`/v1/register`
- Header版本：`Accept: application/vnd.api.v1+json`

