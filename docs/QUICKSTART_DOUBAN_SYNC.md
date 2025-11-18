# 豆瓣电影同步功能 - 快速开始

## 快速测试

### 1. 启动服务

```bash
# 方式1: 使用 make
make run

# 方式2: 直接运行编译后的程序
./bin/server

# 方式3: 使用 go run
go run cmd/server/main.go
```

### 2. 验证定时任务已注册

启动服务后，查看日志确认定时任务已注册：

```bash
tail -f logs/app.log | grep "豆瓣"
```

你应该看到类似的日志：

```
INFO    豆瓣电影同步定时任务已添加    {"schedule": "每8小时执行一次"}
INFO    定时任务调度器已启动
```

### 3. 手动触发同步（测试）

不想等待8小时？可以手动触发同步：

```bash
# 方式1: 使用测试脚本（推荐）
./scripts/test_douban_sync.sh

# 方式2: 使用 curl
curl -X POST http://localhost:8080/api/sync/douban/movies

# 方式3: 使用 httpie (如果已安装)
http POST http://localhost:8080/api/sync/douban/movies
```

### 4. 查看同步结果

#### 查看日志

```bash
# 实时查看同步日志
tail -f logs/app.log | grep "豆瓣\|电影\|同步"

# 查看最近的同步记录
tail -50 logs/app.log | grep "同步"
```

#### 查询数据库

```bash
# 连接到MySQL
mysql -h localhost -u your_user -p your_database

# 查看同步的电影数量
SELECT COUNT(*) FROM videos WHERE source = 'douban';

# 查看最近同步的电影
SELECT id, title, source_id, rating, year, country 
FROM videos 
WHERE source = 'douban' 
ORDER BY created_at DESC 
LIMIT 10;

# 查看需要补充详情的电影
SELECT COUNT(*) 
FROM videos 
WHERE source_id IS NOT NULL 
  AND source_id != 0 
  AND (year IS NULL OR year = 0) 
  AND (country IS NULL OR country = '');
```

## 预期结果

### 第一次同步

第一次运行同步任务时：

1. **第一阶段**：获取最新40部电影的基本信息
   - 如果数据库为空，会保存全部40部
   - 如果已有数据，只保存新增的电影
   - 时间：约5-10秒

2. **第二阶段**：补充前10部电影的详情
   - 每部电影需要2秒（包含延迟）
   - 总时间：约20-30秒

### 日志示例

```
INFO    开始执行豆瓣电影同步任务
INFO    开始同步豆瓣电影数据
INFO    获取到电影列表    {"count": 40}
INFO    保存新电影    {"title": "流浪地球2", "source_id": 35267208}
INFO    保存新电影    {"title": "满江红", "source_id": 35218627}
...
INFO    电影列表同步完成    {"saved_count": 40}
INFO    找到需要更新详情的电影    {"count": 10}
INFO    更新电影详情成功    {"title": "流浪地球2", "source_id": 35267208}
INFO    更新电影详情成功    {"title": "满江红", "source_id": 35218627}
...
INFO    豆瓣电影数据同步完成
INFO    豆瓣电影同步任务执行成功
```

## 常见问题

### 1. 同步任务没有执行？

**检查清单：**

- [ ] 服务是否正常启动？
- [ ] 数据库是否连接成功？
- [ ] 日志中是否有错误信息？

```bash
# 查看服务状态
ps aux | grep server

# 检查端口占用
lsof -i :8080

# 查看完整日志
cat logs/app.log
```

### 2. 获取电影列表失败？

**可能原因：**

- 网络连接问题
- 豆瓣API不可访问
- 请求头配置需要更新

**解决方案：**

```bash
# 测试网络连接
curl -I https://m.douban.com

# 查看详细错误日志
tail -100 logs/app.log | grep "错误\|失败\|error"
```

### 3. 电影详情更新失败？

**可能原因：**

- HTML结构变化（豆瓣页面更新）
- 请求频率过快被限制
- Cookie过期

**解决方案：**

1. 检查日志中的具体错误信息
2. 手动访问豆瓣电影详情页，确认可访问性
3. 如需要，更新 `internal/service/douban_sync_service.go` 中的请求头

### 4. 数据库连接失败？

**检查配置：**

```bash
# 查看数据库配置
cat configs/config.yaml | grep mysql

# 测试数据库连接
mysql -h localhost -u your_user -p
```

## 调试技巧

### 1. 启用详细日志

修改日志级别以获取更详细的信息：

```yaml
# configs/config.yaml
logger:
  level: debug  # 从 info 改为 debug
```

### 2. 单步测试

如果想测试特定步骤：

```go
// 在 internal/service/douban_sync_service.go 中临时修改
func (s *DoubanSyncService) SyncMovies() error {
    // 只测试第一阶段
    return s.fetchAndSaveMovieList()
    
    // 或只测试第二阶段
    // return s.fetchAndUpdateMovieDetails()
}
```

### 3. 检查网络请求

使用代理工具（如 Charles、Fiddler）查看实际的HTTP请求和响应。

### 4. 验证数据

```sql
-- 检查数据完整性
SELECT 
    COUNT(*) as total,
    COUNT(CASE WHEN year IS NOT NULL AND year != 0 THEN 1 END) as with_year,
    COUNT(CASE WHEN country IS NOT NULL AND country != '' THEN 1 END) as with_country
FROM videos 
WHERE source = 'douban';
```

## 性能优化建议

### 1. 调整批处理数量

如果同步速度太慢或太快，可以调整每次处理的数量：

```go
// internal/repository/video_repository.go
// 找到 FindNeedDetailVideos 方法的调用
videos, err := s.videoRepo.FindNeedDetailVideos(10) // 改为 20 或 5
```

### 2. 调整请求延迟

如果不担心被限制，可以减少延迟：

```go
// internal/service/douban_sync_service.go
// 找到 sleep 调用
time.Sleep(2 * time.Second) // 改为 1 * time.Second
```

### 3. 并发处理

如果需要更高的同步速度，可以考虑使用 goroutine 并发处理（需要注意并发控制）。

## 监控和告警

### Prometheus 指标

查看同步相关的指标：

```bash
# 查看所有指标
curl http://localhost:8080/metrics

# 查看特定指标
curl http://localhost:8080/metrics | grep video_service
```

### 日志告警

可以配置日志监控工具（如 ELK、Grafana Loki）来监控同步失败的日志。

## 下一步

1. 查看完整文档：[docs/DOUBAN_SYNC.md](DOUBAN_SYNC.md)
2. 了解定时任务调度：[pkg/infrastructure/scheduler/scheduler.go](../pkg/infrastructure/scheduler/scheduler.go)
3. 自定义同步逻辑：[internal/service/douban_sync_service.go](../internal/service/douban_sync_service.go)

## 需要帮助？

如果遇到问题：

1. 检查日志文件：`logs/app.log`
2. 查看数据库状态
3. 验证网络连接
4. 提交 Issue 或联系开发团队

