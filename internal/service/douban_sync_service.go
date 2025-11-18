// service 包提供业务逻辑层
package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"video-service/internal/model"
	"video-service/internal/pkg/utils"
	"video-service/internal/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DoubanMovieItem 豆瓣电影列表项
type DoubanMovieItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Type   string `json:"type"`
	Rating struct {
		Value float64 `json:"value"`
	} `json:"rating"`
	Pic struct {
		Normal string `json:"normal"`
	} `json:"pic"`
}

// DoubanMovieListResponse 豆瓣电影列表响应
type DoubanMovieListResponse struct {
	Items []DoubanMovieItem `json:"items"`
}

// DoubanSyncService 豆瓣同步服务
type DoubanSyncService struct {
	videoRepo repository.VideoRepository
}

// NewDoubanSyncService 创建豆瓣同步服务实例
func NewDoubanSyncService() *DoubanSyncService {
	return &DoubanSyncService{
		videoRepo: repository.NewVideoRepository(),
	}
}

// SyncMovies 同步豆瓣电影数据
func (s *DoubanSyncService) SyncMovies() error {
	zap.L().Info("开始同步豆瓣电影数据")

	// 第一步：获取最新电影列表并保存基本信息
	if err := s.fetchAndSaveMovieList(); err != nil {
		zap.L().Error("获取电影列表失败", zap.Error(err))
		return err
	}

	// 第二步：补充电影详情信息
	if err := s.fetchAndUpdateMovieDetails(); err != nil {
		zap.L().Error("更新电影详情失败", zap.Error(err))
		return err
	}

	zap.L().Info("豆瓣电影数据同步完成")
	return nil
}

// fetchAndSaveMovieList 获取并保存电影列表
func (s *DoubanSyncService) fetchAndSaveMovieList() error {
	url := "https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie?start=0&limit=80&category=%E6%9C%80%E6%96%B0&type=%E5%85%A8%E9%83%A8"

	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("origin", "https://movie.douban.com")
	req.Header.Set("referer", "https://movie.douban.com/explore")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析JSON
	var movieList DoubanMovieListResponse
	if err := json.Unmarshal(body, &movieList); err != nil {
		return fmt.Errorf("解析JSON失败: %w", err)
	}

	zap.L().Info("获取到电影列表", zap.Int("count", len(movieList.Items)))

	// 遍历电影列表，保存不存在的电影
	savedCount := 0
	for _, item := range movieList.Items {
		// 将字符串ID转换为整数
		sourceID, err := strconv.Atoi(item.ID)
		if err != nil {
			zap.L().Warn("无效的电影ID", zap.String("id", item.ID))
			continue
		}

		// 检查是否已存在
		_, err = s.videoRepo.FindBySourceID(sourceID)
		if err == nil {
			// 已存在，跳过
			continue
		}
		if err != gorm.ErrRecordNotFound {
			zap.L().Error("查询数据库失败", zap.Error(err))
			continue
		}

		// 创建新视频记录
		video := &model.Video{
			ID:           utils.GenerateUserID(), // 使用雪花算法生成ID
			SourceID:     sourceID,
			Source:       "douban",
			Title:        item.Title,
			Type:         item.Type,
			CoverURL:     item.Pic.Normal,
			Rating:       fmt.Sprintf("%.1f", item.Rating.Value),
			EpisodeCount: 0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := s.videoRepo.Create(video); err != nil {
			zap.L().Error("保存视频失败", zap.Error(err), zap.String("title", item.Title))
			continue
		}

		savedCount++
		zap.L().Info("保存新电影", zap.String("title", item.Title), zap.Int("source_id", sourceID))
	}

	zap.L().Info("电影列表同步完成", zap.Int("saved_count", savedCount))
	return nil
}

// fetchAndUpdateMovieDetails 获取并更新电影详情
func (s *DoubanSyncService) fetchAndUpdateMovieDetails() error {
	// 查找需要补充详情的视频（每次处理10条）
	videos, err := s.videoRepo.FindNeedDetailVideos(10)
	if err != nil {
		return fmt.Errorf("查询需要更新的视频失败: %w", err)
	}

	if len(videos) == 0 {
		zap.L().Info("没有需要更新详情的电影")
		return nil
	}

	zap.L().Info("找到需要更新详情的电影", zap.Int("count", len(videos)))

	// 遍历每个视频，获取详情
	for _, video := range videos {
		if err := s.fetchAndUpdateSingleMovieDetail(video); err != nil {
			zap.L().Error("更新电影详情失败", zap.Error(err), zap.String("title", video.Title))
			continue
		}

		// 避免请求过快，休眠2秒
		time.Sleep(2 * time.Second)
	}

	return nil
}

// fetchAndUpdateSingleMovieDetail 获取并更新单个电影详情
func (s *DoubanSyncService) fetchAndUpdateSingleMovieDetail(video *model.Video) error {
	url := fmt.Sprintf("https://movie.douban.com/subject/%d/", video.SourceID)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("referer", "https://movie.douban.com/explore")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	html := string(body)

	// 解析HTML，提取信息
	video.Director = extractFieldWithAttrs(html, "导演")
	video.Actors = extractFieldWithAttrs(html, "主演")
	video.Tags = extractGenres(html)
	video.Country = extractField(html, `<span class="pl">制片国家/地区:</span>`, `<br`)

	// 提取上映日期（完整日期格式 YYYY-MM-DD）
	video.Year = extractReleaseDate(html)

	// 提取片长（只保留数字）
	runtimeStr := extractField(html, `<span class="pl">片长:</span>`, `<br`)
	if runtimeStr != "" {
		video.Runtime = extractNumber(runtimeStr)
	}

	// 提取IMDb ID
	video.IMDbID = extractIMDbID(html)

	// 提取简介
	video.Description = extractDescription(html)

	// 记录提取到的信息（调试用）
	zap.L().Debug("提取电影详情",
		zap.String("title", video.Title),
		zap.String("director", video.Director),
		zap.String("actors", video.Actors),
		zap.String("tags", video.Tags),
		zap.String("country", video.Country),
		zap.String("year", video.Year),
		zap.Int("runtime", video.Runtime),
		zap.String("imdb_id", video.IMDbID),
		zap.Int("desc_len", len(video.Description)))

	// 设置集数为0（电影）
	video.EpisodeCount = 0

	// 更新时间
	video.UpdatedAt = time.Now()

	// 保存到数据库
	if err := s.videoRepo.Update(video); err != nil {
		return fmt.Errorf("更新数据库失败: %w", err)
	}

	zap.L().Info("更新电影详情成功", zap.String("title", video.Title), zap.Int("source_id", video.SourceID))
	return nil
}

// extractField 从HTML中提取字段值
func extractField(html, startTag, endTag string) string {
	startIdx := strings.Index(html, startTag)
	if startIdx == -1 {
		return ""
	}

	startIdx += len(startTag)
	endIdx := strings.Index(html[startIdx:], endTag)
	if endIdx == -1 {
		return ""
	}

	content := html[startIdx : startIdx+endIdx]

	// 移除HTML标签
	content = removeHTMLTags(content)

	// 清理空白字符
	content = strings.TrimSpace(content)

	return content
}

// extractFieldWithAttrs 从HTML中提取包含attrs的字段值
func extractFieldWithAttrs(html, label string) string {
	// 查找标签，例如：<span class='pl'>导演</span>: <span class='attrs'>...</span>
	pattern := fmt.Sprintf(`<span class='pl'>%s</span>:\s*<span class='attrs'>(.*?)</span>`, regexp.QuoteMeta(label))
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		content := matches[1]
		content = removeHTMLTags(content)
		content = strings.TrimSpace(content)
		// 替换多个空格为一个
		content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
		// 替换 " / " 为 ", "
		content = strings.ReplaceAll(content, " / ", ", ")
		return content
	}
	return ""
}

// extractGenres 提取类型标签
func extractGenres(html string) string {
	// 查找所有 <span property="v:genre">...</span>
	re := regexp.MustCompile(`<span property="v:genre">([^<]+)</span>`)
	matches := re.FindAllStringSubmatch(html, -1)
	var genres []string
	for _, match := range matches {
		if len(match) > 1 {
			genres = append(genres, strings.TrimSpace(match[1]))
		}
	}
	return strings.Join(genres, ", ")
}

// extractIMDbID 提取IMDb ID
func extractIMDbID(html string) string {
	// 查找 <span class="pl">IMDb:</span> tt12368458
	re := regexp.MustCompile(`<span class="pl">IMDb:</span>\s*([a-zA-Z0-9]+)`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// extractReleaseDate 提取上映日期（完整日期格式 YYYY-MM-DD）
func extractReleaseDate(html string) string {
	// 优先从 content 属性提取：<span property="v:initialReleaseDate" content="2025-11-15(中国大陆)">
	re := regexp.MustCompile(`<span property="v:initialReleaseDate" content="([0-9]{4}-[0-9]{2}-[0-9]{2})`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}

	// 如果没有找到，尝试从显示文本中提取日期格式
	// 例如：<span class="pl">上映日期:</span> <span>2025-11-15(中国大陆)</span>
	re2 := regexp.MustCompile(`<span class="pl">上映日期:</span>.*?([0-9]{4}-[0-9]{2}-[0-9]{2})`)
	matches2 := re2.FindStringSubmatch(html)
	if len(matches2) > 1 {
		return matches2[1]
	}

	return ""
}

// extractDescription 提取电影简介
func extractDescription(html string) string {
	// 查找 <span property="v:summary" class="">...</span>
	re := regexp.MustCompile(`<span property="v:summary"[^>]*>([\s\S]*?)</span>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		content := matches[1]
		content = removeHTMLTags(content)
		content = strings.TrimSpace(content)
		return content
	}
	return ""
}

// extractNumber 从字符串中提取第一个数字
func extractNumber(str string) int {
	// 使用正则提取数字
	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindStringSubmatch(str)
	if len(matches) > 1 {
		num, _ := strconv.Atoi(matches[1])
		return num
	}
	return 0
}

// removeHTMLTags 移除HTML标签
func removeHTMLTags(html string) string {
	// 移除所有HTML标签
	re := regexp.MustCompile(`<[^>]*>`)
	content := re.ReplaceAllString(html, "")

	// 替换HTML实体
	content = strings.ReplaceAll(content, "&nbsp;", " ")
	content = strings.ReplaceAll(content, "&lt;", "<")
	content = strings.ReplaceAll(content, "&gt;", ">")
	content = strings.ReplaceAll(content, "&amp;", "&")
	content = strings.ReplaceAll(content, "&quot;", "\"")

	return content
}
