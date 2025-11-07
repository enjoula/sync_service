# 系统架构文档

## 📐 架构概述

本项目采用经典的三层架构设计，遵循Go语言项目标准布局，提供清晰的代码组织和良好的可维护性。

## 🏗️ 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         客户端                               │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTP/HTTPS
┌──────────────────────▼──────────────────────────────────────┐
│                     中间件层                                 │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────┐  │
│  │ Trace ID │ Recovery │  Logger  │   JWT    │Prometheus│  │
│  └──────────┴──────────┴──────────┴──────────┴──────────┘  │
└──────────────────────┬──────────────────────────────────────┘
┌──────────────────────▼──────────────────────────────────────┐
│                    Handler层 (API)                          │
│  ┌──────────┬──────────┬──────────┬──────────┐             │
│  │  auth    │   user   │  health  │  debug   │             │
│  └──────────┴──────────┴──────────┴──────────┘             │
└──────────────────────┬──────────────────────────────────────┘
┌──────────────────────▼──────────────────────────────────────┐
│                   Service层 (业务逻辑)                       │
│  ┌───────────────────────────────────────┐                  │
│  │       UserService (用户业务)          │                  │
│  └───────────────────────────────────────┘                  │
└──────────────────────┬──────────────────────────────────────┘
┌──────────────────────▼──────────────────────────────────────┐
│                Repository层 (数据访问)                       │
│  ┌─────────────────┬─────────────────┐                      │
│  │ UserRepository  │TokenRepository  │                      │
│  └─────────────────┴─────────────────┘                      │
└──────────────────────┬──────────────────────────────────────┘
┌──────────────────────▼──────────────────────────────────────┐
│                  基础设施层 (Infrastructure)                 │
│  ┌──────┬──────┬──────┬──────┬──────┬──────────┐           │
│  │MySQL │Redis │Config│Logger│Metrics│Scheduler│           │
│  └──────┴──────┴──────┴──────┴──────┴──────────┘           │
└─────────────────────────────────────────────────────────────┘
```

## 📦 目录结构说明

### cmd/ - 应用程序入口

```
cmd/
└── server/
    └── main.go          # 主程序入口，负责初始化和启动
```

**职责**: 
- 初始化各个基础设施组件
- 启动HTTP服务器
- 优雅关闭

### internal/ - 私有应用程序代码

#### handler/ - HTTP处理层
```
internal/handler/
├── auth.go             # 认证相关（注册、登录）
├── user.go             # 用户相关（获取用户信息）
├── health.go           # 健康检查
└── debug.go            # 调试接口（IP信息）
```

**职责**:
- 处理HTTP请求
- 参数验证和绑定
- 调用Service层
- 返回响应

#### service/ - 业务逻辑层
```
internal/service/
└── user_service.go     # 用户业务逻辑
```

**职责**:
- 实现核心业务逻辑
- 协调多个Repository
- 事务管理
- 业务规则验证

#### repository/ - 数据访问层
```
internal/repository/
└── user_repository.go  # 用户数据访问
```

**职责**:
- 封装数据库操作
- 提供统一的数据访问接口
- 处理数据库错误

#### model/ - 数据模型
```
internal/model/
└── model.go            # 所有数据模型定义
```

**职责**:
- 定义数据结构
- GORM标签配置
- JSON序列化规则

#### middleware/ - HTTP中间件
```
internal/middleware/
├── trace.go            # 追踪ID生成
├── recovery.go         # Panic恢复
├── logger.go           # 请求日志
└── auth.go             # JWT认证
```

**职责**:
- 请求预处理
- 通用功能封装
- 认证授权

#### router/ - 路由配置
```
internal/router/
└── router.go           # 路由注册和配置
```

**职责**:
- 注册路由
- 配置中间件
- 路由分组

#### pkg/ - 内部公共包
```
internal/pkg/
├── auth/               # JWT认证工具
│   └── jwt.go
├── errors/             # 错误定义
│   └── errors.go
├── response/           # 响应封装
│   └── response.go
└── utils/              # 工具函数
    ├── idgen.go        # ID生成器
    ├── avatar.go       # 头像生成器
    ├── nickname.go     # 昵称生成器
    └── ip.go           # IP工具
```

**职责**:
- 提供可复用的工具函数
- 统一的错误处理
- 统一的响应格式

### pkg/ - 公共库（可被外部引用）

#### infrastructure/ - 基础设施层
```
pkg/infrastructure/
├── cache/              # Redis缓存
│   └── redis.go
├── database/           # MySQL数据库
│   └── mysql.go
├── config/             # 配置管理
│   └── config.go
├── logger/             # 日志系统
│   └── logger.go
├── metrics/            # 监控指标
│   └── metrics.go
└── scheduler/          # 定时任务
    └── scheduler.go
```

**职责**:
- 提供基础设施服务
- 管理外部依赖
- 生命周期管理

### deployments/ - 部署文件

```
deployments/
└── docker/
    ├── Dockerfile
    └── docker-compose.yml
```

### scripts/ - 脚本文件

```
scripts/
├── init_etcd.sh        # 初始化Etcd配置
└── build.sh            # 编译脚本
```

### docs/ - 文档

```
docs/
├── ARCHITECTURE.md     # 架构文档（本文档）
├── API.md              # API文档
└── REFACTOR_PLAN.md    # 重构计划
```

## 🔄 请求处理流程

### 1. 用户注册流程

```
Client → POST /register
    ↓
Middleware: Trace (生成追踪ID)
    ↓
Middleware: Recovery (捕获panic)
    ↓
Middleware: Logger (记录请求)
    ↓
Handler: Register
    ├─ 参数验证
    ├─ 获取IP和设备信息
    └─ Service: Register
        ├─ 验证用户名密码
        ├─ 检查用户是否存在
        ├─ bcrypt加密密码
        ├─ 生成用户ID（雪花算法）
        ├─ 生成随机昵称和头像
        ├─ Repository: Create User
        ├─ 生成JWT Token
        ├─ Repository: Create Token
        └─ 管理活跃Token
    ↓
返回响应（用户信息 + Token）
```

### 2. 用户登录流程

```
Client → POST /login
    ↓
中间件处理（同上）
    ↓
Handler: Login
    ├─ 参数验证
    ├─ 获取IP和设备信息
    └─ Service: Login
        ├─ Repository: FindByUsername
        ├─ bcrypt验证密码
        ├─ 生成JWT Token
        ├─ Repository: Create Token
        └─ 管理活跃Token
    ↓
返回响应（用户信息 + Token）
```

### 3. 获取用户信息流程

```
Client → GET /user/me (Header: Authorization: Bearer <token>)
    ↓
中间件处理
    ↓
Middleware: JWTAuth
    ├─ 提取Token
    ├─ Repository: IsTokenActive (检查is_active)
    ├─ 验证Token签名和过期时间
    └─ 提取用户信息存入Context
    ↓
Handler: Me
    ├─ 从Context获取用户信息
    └─ 返回用户数据
    ↓
返回响应
```

## 🔑 核心设计模式

### 1. 依赖注入

```go
// Service依赖Repository
type userService struct {
    userRepo      repository.UserRepository
    userTokenRepo repository.UserTokenRepository
}

func NewUserService() UserService {
    return &userService{
        userRepo:      repository.NewUserRepository(),
        userTokenRepo: repository.NewUserTokenRepository(),
    }
}
```

### 2. 接口设计

```go
// 定义接口，便于测试和扩展
type UserService interface {
    Register(username, password, device, ipAddress string) (string, *model.User, error)
    Login(username, password, device, ipAddress string) (string, *model.User, error)
    GetCurrentUser(c *gin.Context) (string, error)
}
```

### 3. 中间件模式

```go
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理
        // ...
        c.Next()  // 继续处理
        // 后置处理（如果需要）
    }
}
```

## 📊 数据流向

```
HTTP请求 → 中间件 → Handler → Service → Repository → Database
                                                        ↓
HTTP响应 ← 中间件 ← Handler ← Service ← Repository ← Database
```

## 🛡️ 安全设计

### 1. 密码安全
- 使用bcrypt加密
- 成本因子：bcrypt.DefaultCost (10)
- 密码哈希不返回给客户端

### 2. JWT Token管理
- 24小时过期时间
- 双重验证（数据库is_active + JWT签名）
- 支持强制下线（设置is_active=0）
- 多设备管理（最多3个活跃token）

### 3. IP追踪
- 智能获取真实IP（支持反向代理）
- 记录登录IP
- 支持IP黑白名单（待实现）

## 🎯 设计原则

1. **单一职责**：每个层、每个模块只负责一件事
2. **依赖倒置**：高层模块不依赖低层模块，都依赖抽象（接口）
3. **接口隔离**：客户端不应该依赖它不需要的接口
4. **开闭原则**：对扩展开放，对修改关闭
5. **DRY原则**：不重复自己（Don't Repeat Yourself）

## 🚀 扩展性

### 添加新业务模块

1. 在`internal/model/`添加数据模型
2. 在`internal/repository/`添加数据访问层
3. 在`internal/service/`添加业务逻辑层
4. 在`internal/handler/`添加HTTP处理器
5. 在`internal/router/`注册路由

### 添加新中间件

1. 在`internal/middleware/`创建中间件文件
2. 在`internal/router/router.go`中注册

### 添加新基础设施

1. 在`pkg/infrastructure/`创建相应目录
2. 实现初始化和生命周期管理
3. 在`cmd/server/main.go`中调用初始化

## 📈 性能优化点

1. **数据库连接池**：配置合理的连接池大小
2. **Redis缓存**：缓存热点数据（待实现）
3. **批量查询**：使用IN查询减少数据库访问
4. **索引优化**：为常用查询添加索引
5. **限流**：添加API限流（待实现）

## 🔍 监控和可观测性

1. **日志**：
   - 结构化日志（JSON）
   - 追踪ID关联
   - 分级日志（Debug/Info/Warn/Error）

2. **指标**：
   - Prometheus指标
   - HTTP请求计数
   - 响应时间分布

3. **追踪**：
   - 请求追踪ID
   - 分布式追踪（待实现）

## 📝 待改进项

1. 添加单元测试和集成测试
2. 实现Redis缓存策略
3. 添加API限流
4. 实现分布式追踪（Jaeger/Zipkin）
5. 添加视频相关API实现
6. 优化数据库查询性能
7. 添加配置热重载
8. 实现优雅关闭

## 🔗 相关文档

- [API文档](./API.md)
- [重构计划](./REFACTOR_PLAN.md)
- [README](../README.md)

