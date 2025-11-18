# 🎬 Sync Service — 视频数据同步服务

一个基于 Golang + Gin + GORM + Redis + Prometheus 的视频数据同步服务，支持从豆瓣等多个数据源自动同步影视信息，包含定时任务、JWT 认证、TraceID 链路追踪、日志系统等完整功能。

## ✨ 核心功能

- 🎯 **豆瓣电影同步**：自动同步最新电影信息（每8小时）
- 🕒 **定时任务调度**：基于cron的灵活定时任务系统
- 🔄 **手动触发接口**：支持通过API手动触发同步任务
- 📊 **监控指标**：集成Prometheus监控
- 🔍 **链路追踪**：每个请求自动生成TraceID
- 📝 **日志系统**：结构化日志，文件+控制台双输出
- 🐳 **Docker部署**：完整的Docker Compose配置

## 🎬 豆瓣电影同步功能

### 自动同步内容

**第一阶段：电影列表**
- 获取豆瓣最新80部电影
- 保存基本信息：标题、类型、评分、封面等
- 自动去重，避免重复保存

**第二阶段：电影详情**
- 导演、主演（多个逗号分隔）
- 类型标签（多个逗号分隔）
- 制片国家/地区
- 上映日期（YYYY-MM-DD格式）
- 片长（分钟）
- IMDb ID
- 电影简介

### 同步方式

1. **自动同步**：每8小时自动执行（服务启动后自动注册）
2. **手动触发**：通过HTTP接口手动触发
   ```bash
   curl -X POST http://localhost:5500/api/sync/douban/movies
   ```

### 查看同步日志

```bash
# 实时查看同步日志
tail -f logs/app.log | grep "豆瓣"

# 或使用测试脚本
./scripts/test_douban_sync.sh
```

## 🚀 快速启动

### 方式一：Docker 环境（推荐）

1. **启动所有服务**
```bash
cd sync_service
docker-compose -f deployments/docker/docker-compose.yml up -d
```

2. **查看服务状态**
```bash
docker-compose -f deployments/docker/docker-compose.yml ps
```

3. **查看日志**
```bash
docker-compose -f deployments/docker/docker-compose.yml logs -f sync_service
```

### 方式二：本地开发环境

1. **启动依赖服务（MySQL、Redis）**
```bash
docker-compose -f deployments/docker/docker-compose.yml up -d mysql redis
```

2. **修改配置文件**
```yaml
# configs/config.yaml
server:
  addr: ":5500"
mysql:
  dsn: "root:123456@tcp(127.0.0.1:5506)/video_service?charset=utf8mb4&parseTime=True&loc=Local"
redis:
  addr: "127.0.0.1:5509"
etcd:
  addr: ""  # 本地开发可以不使用etcd
```

3. **运行服务**
```bash
# 方式1：直接运行
go run cmd/server/main.go

# 方式2：编译后运行
go build -o bin/server cmd/server/main.go
./bin/server

# 方式3：使用Make
make run
```

## 📋 配置说明

### 基础配置 (configs/config.yaml)

```yaml
server:
  addr: ":5500"          # 服务监听端口

mysql:
  dsn: "root:123456@tcp(mysql:3306)/video_service?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  addr: "redis:6379"     # Redis地址
  
etcd:
  addr: "http://etcd:2379"  # etcd地址（可选，用于存储敏感配置）
  
prometheus:
  global:
    scrape_interval: 60s
  scrape_configs:
    - job_name: 'sync_service'
      metrics_path: /metrics
      static_configs:
        - targets: ['sync_service:5500']
```

### etcd敏感配置（可选）

如果使用etcd存储敏感信息：

```bash
docker exec -it Etcd /bin/sh
etcdctl put /video-service/secret '{
  "jwt_key": "your-secret-jwt-key-change-me",
  "mysql_dsn": "root:123456@tcp(mysql:3306)/video_service?charset=utf8mb4&parseTime=True&loc=Local"
}'
```

## 🧪 测试接口

### 健康检查
```bash
curl http://localhost:5500/ping
# 返回: {"code":0,"message":"pong","data":{"time":"ok"},...}
```

### 手动触发豆瓣同步
```bash
curl -X POST http://localhost:5500/api/sync/douban/movies
# 返回: {"code":0,"message":"同步任务已启动，正在后台执行",...}
```

### Prometheus指标
```bash
curl http://localhost:5500/metrics
```

## 📊 监控访问

- **后端服务**: http://localhost:5500
- **Prometheus**: http://localhost:5590
- **健康检查**: http://localhost:5500/ping
- **监控指标**: http://localhost:5500/metrics

## 📁 项目结构

```
sync_service/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── handler/                 # HTTP请求处理器
│   │   ├── health.go           # 健康检查
│   │   └── sync.go             # 同步任务接口
│   ├── middleware/              # 中间件
│   │   ├── auth.go             # JWT认证
│   │   ├── logger.go           # 日志记录
│   │   ├── recovery.go         # 错误恢复
│   │   └── trace.go            # 链路追踪
│   ├── model/                   # 数据模型
│   │   └── model.go            # 数据库模型定义
│   ├── repository/              # 数据访问层
│   │   └── video_repository.go # 视频数据操作
│   ├── service/                 # 业务逻辑层
│   │   └── douban_sync_service.go  # 豆瓣同步服务
│   ├── router/                  # 路由配置
│   │   └── router.go
│   └── pkg/                     # 内部工具包
│       ├── auth/               # JWT工具
│       ├── errors/             # 错误定义
│       ├── response/           # 统一响应格式
│       └── utils/              # 工具函数
├── pkg/infrastructure/          # 基础设施
│   ├── cache/                  # Redis缓存
│   ├── config/                 # 配置管理
│   ├── database/               # MySQL数据库
│   ├── logger/                 # 日志系统
│   ├── metrics/                # Prometheus指标
│   └── scheduler/              # 定时任务调度器
├── configs/                     # 配置文件
│   ├── config.yaml             # 主配置文件
│   └── prometheus.yml          # Prometheus配置
├── deployments/docker/          # Docker部署文件
│   ├── Dockerfile
│   └── docker-compose.yml
├── docs/                        # 文档
│   ├── DOUBAN_SYNC.md          # 豆瓣同步详细文档
│   └── QUICKSTART_DOUBAN_SYNC.md  # 快速开始指南
├── scripts/                     # 脚本工具
│   ├── build.sh                # 编译脚本
│   └── test_douban_sync.sh     # 同步功能测试脚本
├── migrations/                  # 数据库迁移文件
│   └── init.sql
└── logs/                        # 日志文件
    └── app.log
```

## 🔧 Docker Compose 管理

### 启动所有服务
```bash
docker-compose -f deployments/docker/docker-compose.yml up -d
```

### 停止所有服务
```bash
docker-compose -f deployments/docker/docker-compose.yml down
```

### 重新构建并启动
```bash
docker-compose -f deployments/docker/docker-compose.yml down
docker-compose -f deployments/docker/docker-compose.yml build --no-cache
docker-compose -f deployments/docker/docker-compose.yml up -d
```

### 查看日志
```bash
# 查看所有服务日志
docker-compose -f deployments/docker/docker-compose.yml logs -f

# 查看特定服务日志
docker-compose -f deployments/docker/docker-compose.yml logs -f sync_service
docker-compose -f deployments/docker/docker-compose.yml logs -f mysql
```

## 📖 详细文档

- **[豆瓣同步功能详解](docs/DOUBAN_SYNC.md)** - 完整的功能说明和技术实现
- **[快速开始指南](docs/QUICKSTART_DOUBAN_SYNC.md)** - 测试、调试和故障排查指南

## ⚙️ 数据库

### 自动迁移

服务启动时会自动执行 GORM 的 `AutoMigrate()`，创建以下数据表：

- `users` - 用户表
- `user_tokens` - 用户Token记录
- `videos` - 视频信息表
- `episodes` - 剧集信息表
- `danmakus` - 弹幕表
- `user_favorites` - 用户收藏表
- `user_watch_progress` - 观看进度表
- `app_versions` - 应用版本表

### 手动初始化

也可以使用SQL文件手动初始化：

```bash
docker exec -i Mysql mysql -uroot -p123456 video_service < migrations/init.sql
```

## 📝 日志系统

日志同时输出到：
- **控制台**：JSON格式，方便开发调试
- **文件**：`logs/app.log`，自动按日期和大小轮转

### 查看日志

```bash
# 实时查看所有日志
tail -f logs/app.log

# 只看豆瓣同步相关日志
tail -f logs/app.log | grep "豆瓣"

# 查看错误日志
tail -f logs/app.log | grep "error"
```

## 🔐 注意事项

1. **数据库初始化**
   - 首次启动会自动创建表结构
   - 确保MySQL容器健康检查通过后再启动应用

2. **敏感信息**
   - 生产环境请修改MySQL、Redis密码
   - 更新etcd中的JWT密钥

3. **定时任务**
   - 豆瓣同步任务默认每8小时执行一次
   - 可在 `pkg/infrastructure/scheduler/scheduler.go` 中修改Cron表达式

4. **请求频率**
   - 豆瓣详情页请求间隔2秒，避免被封禁
   - 如需调整，修改 `internal/service/douban_sync_service.go`

5. **数据去重**
   - 通过 `source_id` 字段确保不会重复保存相同电影
   - 已存在的电影会跳过，只保存新增的

## 🛠️ 开发建议

### 修改同步频率

编辑 `pkg/infrastructure/scheduler/scheduler.go`：

```go
// 当前：每8小时执行一次
_, err := cronScheduler.AddFunc("0 0 */8 * * *", func() {
    // ...
})

// 修改为每4小时：
_, err := cronScheduler.AddFunc("0 0 */4 * * *", func() {
    // ...
})

// 修改为每天凌晨2点：
_, err := cronScheduler.AddFunc("0 0 2 * * *", func() {
    // ...
})
```

### 添加新的数据源

参考 `internal/service/douban_sync_service.go`，创建新的同步服务：

1. 创建新的service文件（如 `tmdb_sync_service.go`）
2. 实现数据获取和解析逻辑
3. 在scheduler中注册定时任务
4. 在router中添加手动触发接口

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

---

**快速链接**：
- 🐛 [报告问题](https://github.com/enjoula/sync_service/issues)
- 📚 [查看文档](docs/)
- 🎬 [豆瓣同步详解](docs/DOUBAN_SYNC.md)
