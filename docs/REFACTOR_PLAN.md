# 目录结构重构方案

## 📋 重构目标
将当前扁平化的目录结构重构为更清晰、更符合Go标准项目布局的层次化结构。

## 🎯 设计原则
1. **按功能分层**：清晰的三层架构（Handler → Service → Repository）
2. **按业务分组**：相关功能放在同一目录下
3. **职责单一**：每个目录只负责一类功能
4. **标准布局**：遵循Go社区标准项目布局

## 📂 新目录结构

```
video_service/
├── cmd/                          # 应用程序入口
│   └── server/
│       └── main.go              # 主程序入口
│
├── internal/                     # 私有应用程序代码（不可被外部引用）
│   ├── handler/                 # HTTP处理层（原api目录）
│   │   ├── auth.go             # 认证处理器（注册、登录）
│   │   ├── user.go             # 用户处理器（用户信息）
│   │   ├── health.go           # 健康检查
│   │   └── debug.go            # 调试接口（IP信息）
│   │
│   ├── service/                 # 业务逻辑层
│   │   └── user_service.go     # 用户业务逻辑
│   │
│   ├── repository/              # 数据访问层
│   │   └── user_repository.go  # 用户数据访问
│   │
│   ├── model/                   # 数据模型（原models目录）
│   │   └── model.go            # 所有数据模型定义
│   │
│   ├── middleware/              # HTTP中间件
│   │   ├── auth.go             # JWT认证中间件
│   │   ├── logger.go           # 日志中间件
│   │   ├── recovery.go         # 恢复中间件
│   │   └── trace.go            # 追踪中间件
│   │
│   ├── router/                  # 路由配置
│   │   └── router.go           # 路由注册
│   │
│   └── pkg/                     # 内部公共包
│       ├── auth/               # JWT认证工具
│       │   └── jwt.go
│       ├── errors/             # 错误定义
│       │   └── errors.go
│       ├── response/           # 响应封装
│       │   └── response.go
│       └── utils/              # 工具函数
│           ├── idgen.go        # ID生成器（原idgen目录）
│           ├── avatar.go       # 头像生成器（原avatar目录）
│           ├── nickname.go     # 昵称生成器（原nickname目录）
│           └── ip.go           # IP工具（原utils/ip.go）
│
├── pkg/                         # 公共库（可被外部项目引用）
│   └── infrastructure/          # 基础设施层
│       ├── cache/              # Redis缓存
│       │   └── redis.go
│       ├── database/           # MySQL数据库
│       │   └── mysql.go
│       ├── config/             # 配置管理
│       │   └── config.go
│       ├── logger/             # 日志系统
│       │   └── logger.go
│       ├── metrics/            # 监控指标
│       │   └── metrics.go
│       └── scheduler/          # 定时任务
│           └── scheduler.go
│
├── configs/                     # 配置文件
│   ├── config.yaml             # 应用配置
│   └── prometheus.yml          # Prometheus配置
│
├── migrations/                  # 数据库迁移文件
│   └── init.sql                # 初始化SQL
│
├── scripts/                     # 脚本文件（新增）
│   ├── init_etcd.sh           # 初始化etcd配置
│   └── build.sh               # 编译脚本
│
├── docs/                        # 文档（新增）
│   ├── API.md                  # API文档
│   ├── ARCHITECTURE.md         # 架构文档
│   └── REFACTOR_PLAN.md        # 本文档
│
├── deployments/                 # 部署文件（新增）
│   └── docker/
│       ├── Dockerfile          # Docker镜像构建
│       └── docker-compose.yml  # Docker编排
│
├── logs/                        # 日志输出目录
│   └── app.log                 # 应用日志
│
├── go.mod                       # Go模块定义
├── go.sum                       # 依赖版本锁定
├── Makefile                     # 构建脚本（新增）
└── README.md                    # 项目说明

```

## 🔄 主要变更

### 1. API层重命名为Handler层
- `internal/api/` → `internal/handler/`
- 文件按功能合并：
  - `auth_handler.go` → `auth.go`（认证相关）
  - `user_handler.go` → `user.go`（用户相关）
  - `health_handler.go` → `health.go`（健康检查）
  - `ip_handler.go` → `debug.go`（调试接口）

### 2. 工具类整合
将分散的小工具合并到 `internal/pkg/utils/`：
- `internal/idgen/` → `internal/pkg/utils/idgen.go`
- `internal/avatar/` → `internal/pkg/utils/avatar.go`
- `internal/nickname/` → `internal/pkg/utils/nickname.go`
- `internal/utils/ip.go` → `internal/pkg/utils/ip.go`

### 3. 基础设施层提取
将基础设施相关代码移到 `pkg/infrastructure/`：
- `internal/cache/` → `pkg/infrastructure/cache/`
- `internal/db/` → `pkg/infrastructure/database/`
- `internal/config/` → `pkg/infrastructure/config/`
- `internal/logger/` → `pkg/infrastructure/logger/`
- `internal/metrics/` → `pkg/infrastructure/metrics/`
- `internal/scheduler/` → `pkg/infrastructure/scheduler/`

### 4. 核心包移到内部pkg
- `internal/auth/` → `internal/pkg/auth/`
- `internal/errors/` → `internal/pkg/errors/`
- `internal/response/` → `internal/pkg/response/`

### 5. 数据模型重命名
- `internal/models/` → `internal/model/`（单数形式更符合Go习惯）

### 6. 新增目录
- `scripts/`：存放初始化、构建等脚本
- `docs/`：存放项目文档
- `deployments/docker/`：存放Docker相关文件
- `Makefile`：统一构建命令

## 📝 导入路径变更对照表

| 旧路径 | 新路径 |
|--------|--------|
| `video-service/internal/api` | `video-service/internal/handler` |
| `video-service/internal/auth` | `video-service/internal/pkg/auth` |
| `video-service/internal/avatar` | `video-service/internal/pkg/utils` |
| `video-service/internal/cache` | `video-service/pkg/infrastructure/cache` |
| `video-service/internal/config` | `video-service/pkg/infrastructure/config` |
| `video-service/internal/db` | `video-service/pkg/infrastructure/database` |
| `video-service/internal/errors` | `video-service/internal/pkg/errors` |
| `video-service/internal/idgen` | `video-service/internal/pkg/utils` |
| `video-service/internal/logger` | `video-service/pkg/infrastructure/logger` |
| `video-service/internal/metrics` | `video-service/pkg/infrastructure/metrics` |
| `video-service/internal/models` | `video-service/internal/model` |
| `video-service/internal/nickname` | `video-service/internal/pkg/utils` |
| `video-service/internal/response` | `video-service/internal/pkg/response` |
| `video-service/internal/scheduler` | `video-service/pkg/infrastructure/scheduler` |
| `video-service/internal/utils` | `video-service/internal/pkg/utils` |

## ✅ 重构优势

### 1. 更清晰的层次结构
- **业务层**：handler → service → repository
- **支撑层**：infrastructure（基础设施）
- **工具层**：utils（纯函数工具）

### 2. 更好的可维护性
- 相关代码聚合在一起
- 减少目录深度，降低查找成本
- 职责清晰，易于定位问题

### 3. 更符合Go标准
- `internal/`：私有代码，不可被外部引用
- `pkg/`：公共库，可被外部项目引用
- `cmd/`：应用程序入口

### 4. 便于扩展
- 添加新业务模块时结构清晰
- 基础设施代码可复用
- 便于团队协作

## 🚀 迁移步骤

1. ✅ 创建新目录结构
2. ✅ 移动文件到新位置
3. ✅ 更新所有import路径
4. ✅ 移动Docker文件到deployments/
5. ✅ 创建Makefile
6. ✅ 创建初始化脚本
7. ✅ 更新README.md
8. ✅ 验证编译通过
9. ✅ 验证Docker构建通过
10. ✅ 删除旧目录

## 📌 注意事项

1. **保持功能不变**：重构只改变目录结构，不改变业务逻辑
2. **逐步迁移**：按模块逐个迁移，每次迁移后验证编译
3. **更新文档**：同步更新README和其他文档
4. **Git提交**：每个阶段完成后提交代码，便于回滚

