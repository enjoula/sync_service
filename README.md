# 🧩 Video Service — 启动与开发说明

一个基于 Golang + Gin + GORM + Redis + etcd + Prometheus 的完整在线视频服务后端框架，支持 JWT 登录鉴权、TraceID 链路追踪、日志系统、自动迁移与定时任务。


## 🚀 快速启动（Docker 环境）
1️⃣ 解压项目
```bash
cd video_service
```

2️⃣ 配置检查
可在 configs/config.yaml 中修改服务配置（如 MySQL、Redis、etcd 地址）：
```yaml
server:
  addr: ":5501"

etcd:
  addr: "http://etcd:2379"

mysql:
  dsn: "root:123456@tcp(mysql:3306)/video_service?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  addr: "redis:6379"
  pass: ""
  db: 0

MoonTV:
  addr: ""
  pass: ""


```
3️⃣ 启动服务
```bash
docker-compose up --build

```
4️⃣ 初始化 etcd 配置（敏感信息）

服务启动后，在 etcd 中写入 JWT 密钥与数据库连接信息：

```bash
docker exec -it etcd /bin/sh
etcdctl put /video-service/secret '{
  "jwt_key": "super-secret-key-change-me",
  "mysql_dsn": "root:123456@tcp(mysql:3306)/video_service?charset=utf8mb4&parseTime=True&loc=Local"
}'
```
🔄 etcd 支持实时热加载，无需重启应用。

5️⃣ 测试接口
1、 注册用户

curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pwd"}'

2、登录
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pwd"}'

3、获取用户信息
curl http://localhost:8080/user/me \
  -H "Authorization: <token>"



返回示例(返回头部中也包含 X-Trace-Id。)：
```json
{
  "code": 0,
  "msg": "login success",
  "data": {
    "token": "xxxx.yyyy.zzzz"
  },
  "trace_id": "2a3c1b7f..."
}

```

6️⃣ 访问监控

- 后端 API: http://localhost:8080
- Prometheus: http://localhost:9090



# ⚙️ 注意事项

- 启动时系统会自动执行 GORM 的 AutoMigrate() 建表逻辑；

- 项目中同时包含 migrations/init.sql，可用于手动或自动初始化数据库结构；

- 日志同时输出到：

- 控制台（JSON 格式）

- 日志文件 ./logs/app.log 文件；

- Trace ID 自动附加到所有响应（JSON 字段 + HTTP Header）；

- etcd 支持实时更新 JWT 密钥与数据库配置；

- 定时任务示例位于 internal/scheduler/cron.go。


# 下一步建议

## 迁移管理优化

- 推荐集成 golang-migrate，在容器启动时自动执行 migrations/*.sql。

- 可在 docker-compose 的 app 服务中添加 migration 启动命令。

- SQL 初始化

- 当前包中包含 migrations/init.sql（空模板）。

- 可将你提供的完整表结构 SQL 文件替换进去，用于生产环境初始化。

- 日志与监控

- 已内置 Prometheus 采集接口 /metrics；

- 可在未来对接 Grafana 仪表盘。

- 安全增强

- 请尽快在 etcd 中更新默认 JWT 密钥；

- 生产环境建议关闭匿名 etcd 访问。
