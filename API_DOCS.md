# Video Service API 文档

## 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **响应格式**: 统一JSON格式

## 统一响应格式

所有接口返回统一的JSON格式：

```json
{
  "code": 0,              // 业务状态码：0表示成功，其他值表示错误
  "message": "Success",   // 响应消息
  "data": {},            // 响应数据（成功时包含业务数据，失败时为null）
  "trace_id": "..."      // 请求追踪ID（用于分布式追踪）
}
```

## 错误码说明

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 0 | 200 | 成功 |
| 400 | 200 | 请求参数错误 |
| 401 | 200 | 未授权（需要登录或token无效） |
| 409 | 200 | 资源冲突（如用户已存在） |
| 500 | 200 | 服务器内部错误 |

---

## 1. 用户注册接口

### 1.1 接口信息

- **接口路径**: `/register`
- **请求方法**: `POST`
- **接口描述**: 用户注册接口，注册成功后自动生成token实现注册即登录
- **是否需要认证**: 否

### 1.2 请求参数

#### 请求头 (Headers)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Content-Type | string | 是 | application/json |
| X-Device | string | 否 | 设备类型（如"web"、"tv"、"mobile"等） |

#### 请求体 (Body)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**请求示例**:

```json
{
  "username": "testuser",
  "password": "test123456"
}
```

### 1.3 响应数据

#### 成功响应 (code: 0)

```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "id": 244812787503398912,        // 用户ID
    "username": "testuser",          // 用户名
    "nickname": "孙悟空",             // 用户昵称（自动生成）
    "token": "eyJhbGciOiJIUzI1NiIs..." // JWT token（过期时间1年）
  },
  "trace_id": "e865ddb1-e441-4b49-84d1-cd1fae0e2434"
}
```

#### 错误响应

**用户名重复** (code: 409):

```json
{
  "code": 409,
  "message": "用户名重复",
  "data": null,
  "trace_id": "..."
}
```

**请求参数无效** (code: 400):

```json
{
  "code": 400,
  "message": "请求参数无效",
  "data": null,
  "trace_id": "..."
}
```

**账号密码不能为空** (code: 400):

```json
{
  "code": 400,
  "message": "账号密码不能为空",
  "data": null,
  "trace_id": "..."
}
```

### 1.4 功能说明

1. 验证用户名和密码
2. 使用bcrypt加密密码
3. 自动生成用户昵称
4. 创建用户记录
5. 自动生成JWT token（过期时间1年）
6. 保存token到user_tokens表
7. 如果用户已有3个或更多活跃token，自动停用最早创建的token

### 1.5 调用示例

**cURL**:

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -H "X-Device: web" \
  -d '{
    "username": "testuser",
    "password": "test123456"
  }'
```

**JavaScript (fetch)**:

```javascript
fetch('http://localhost:8080/register', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-Device': 'web'
  },
  body: JSON.stringify({
    username: 'testuser',
    password: 'test123456'
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

---

## 2. 用户登录接口

### 2.1 接口信息

- **接口路径**: `/login`
- **请求方法**: `POST`
- **接口描述**: 用户登录接口，验证账号密码后生成JWT token
- **是否需要认证**: 否

### 2.2 请求参数

#### 请求头 (Headers)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Content-Type | string | 是 | application/json |

#### 请求体 (Body)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |
| device | string | 否 | 设备类型（如"web"、"tv"、"mobile"等） |

**请求示例**:

```json
{
  "username": "testuser",
  "password": "test123456",
  "device": "web"
}
```

### 2.3 响应数据

#### 成功响应 (code: 0)

```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "username": "testuser",          // 用户名（来自users表）
    "avatar": "",                    // 用户头像URL（来自users表）
    "nickname": "孙悟空",             // 用户昵称（来自users表）
    "token": "eyJhbGciOiJIUzI1NiIs..." // JWT token（来自user_tokens表，过期时间1年）
  },
  "trace_id": "3d0d6422-3d8c-4dda-990a-4aab140925c5"
}
```

#### 错误响应

**账号密码不能为空** (code: 400):

```json
{
  "code": 400,
  "message": "账号密码不能为空",
  "data": null,
  "trace_id": "..."
}
```

**用户名或密码错误** (code: 401):

```json
{
  "code": 401,
  "message": "用户名或密码错误",
  "data": null,
  "trace_id": "..."
}
```

**请求参数无效** (code: 400):

```json
{
  "code": 400,
  "message": "请求参数无效",
  "data": null,
  "trace_id": "..."
}
```

### 2.4 功能说明

1. 验证username和password参数不能为空
2. 验证用户名和密码是否正确
3. 生成JWT token（过期时间1年）
4. 检查该用户的token数量
5. 如果用户已有3个或更多活跃token，自动停用最早创建的token
6. 保存新token到user_tokens表
7. 返回用户信息和token

### 2.5 Token管理规则

- Token过期时间：1年（365天）
- 每个用户最多保留3个活跃token
- 当超过3个时，自动停用最早创建的token（is_active置为0）
- Token存储在user_tokens表中

### 2.6 调用示例

**cURL**:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "test123456",
    "device": "web"
  }'
```

**JavaScript (fetch)**:

```javascript
fetch('http://localhost:8080/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    username: 'testuser',
    password: 'test123456',
    device: 'web'
  })
})
.then(response => response.json())
.then(data => {
  if (data.code === 0) {
    console.log('登录成功，token:', data.data.token);
    // 保存token用于后续请求
    localStorage.setItem('token', data.data.token);
  } else {
    console.error('登录失败:', data.message);
  }
});
```

**Python (requests)**:

```python
import requests

url = "http://localhost:8080/login"
headers = {"Content-Type": "application/json"}
data = {
    "username": "testuser",
    "password": "test123456",
    "device": "web"
}

response = requests.post(url, json=data, headers=headers)
result = response.json()

if result["code"] == 0:
    print("登录成功")
    token = result["data"]["token"]
    print(f"Token: {token}")
else:
    print(f"登录失败: {result['message']}")
```

---

## 3. 使用Token进行认证

登录成功后，需要在后续请求的请求头中携带token：

```
Authorization: <token>
```

**示例**:

```bash
curl -X GET http://localhost:8080/user/me \
  -H "Authorization: eyJhbGciOiJIUzI1NiIs..."
```

---

## 4. 常见错误处理

### 4.1 参数验证错误

- **错误码**: 400
- **常见原因**:
  - 缺少必填参数（username或password）
  - 参数为空字符串
  - JSON格式错误
- **解决方案**: 检查请求参数是否完整且格式正确

### 4.2 认证错误

- **错误码**: 401
- **常见原因**:
  - 用户名或密码错误
  - Token已过期
  - Token格式错误
- **解决方案**: 
  - 检查用户名和密码是否正确
  - 重新登录获取新token

### 4.3 资源冲突

- **错误码**: 409
- **常见原因**:
  - 用户名已存在（注册时）
  - Token重复（极少发生）
- **解决方案**: 使用不同的用户名注册

### 4.4 服务器错误

- **错误码**: 500
- **常见原因**: 服务器内部错误
- **解决方案**: 联系技术支持或稍后重试

---

## 5. 测试用例

### 5.1 注册接口测试

```bash
# 正常注册
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123456"}'

# 用户名重复
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123456"}'

# 缺少参数
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser"}'
```

### 5.2 登录接口测试

```bash
# 正常登录
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123456","device":"web"}'

# 缺少username
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"password":"test123456","device":"web"}'

# 缺少password
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","device":"web"}'

# 错误密码
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"wrong","device":"web"}'

# 不存在的用户
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"nonexist","password":"test123456","device":"web"}'
```

---

## 6. 注意事项

1. **Token安全**:
   - Token应妥善保管，不要泄露给他人
   - Token过期时间为1年，过期后需要重新登录
   - 建议在HTTPS环境下使用

2. **多设备登录**:
   - 支持多设备同时登录
   - 每个用户最多保留3个活跃token
   - 超过限制时，最早创建的token会被自动停用

3. **密码安全**:
   - 密码使用bcrypt加密存储
   - 建议使用强密码（包含字母、数字、特殊字符）

4. **请求追踪**:
   - 所有响应都包含trace_id字段
   - trace_id也包含在响应头X-Trace-Id中
   - 可用于问题排查和日志追踪

---

## 7. 更新日志

- **2025-11-06**: 
  - 登录接口增加username字段返回
  - 登录接口增加参数验证（username和password不能为空）
  - 错误信息优化："账号密码不能为空"

