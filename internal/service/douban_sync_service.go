// service 包提供业务逻辑层
package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"video-service/internal/model"
	"video-service/internal/pkg/utils"
	"video-service/internal/repository"
	"video-service/pkg/infrastructure/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DoubanItem 豆瓣列表项
type DoubanItem struct {
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

// DoubanListResponse 豆瓣列表响应
type DoubanListResponse struct {
	Items []DoubanItem `json:"items"`
}

// SearchResponse 搜索播放地址响应
type SearchResponse struct {
	Results []SearchResult `json:"results"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Title      string   `json:"title"`
	SourceName string   `json:"source_name"`
	Episodes   []string `json:"episodes"`
}

// DoubanSyncService 豆瓣同步服务
type DoubanSyncService struct {
	videoRepo   repository.VideoRepository
	episodeRepo repository.EpisodeRepository
}

// NewDoubanSyncService 创建豆瓣同步服务实例
func NewDoubanSyncService() *DoubanSyncService {
	return &DoubanSyncService{
		videoRepo:   repository.NewVideoRepository(),
		episodeRepo: repository.NewEpisodeRepository(),
	}
}

// SyncAll 同步所有豆瓣数据
func (s *DoubanSyncService) SyncAll() error {
	zap.L().Info("开始同步豆瓣数据")

	// 第一步：获取最新列表并保存基本信息
	if err := s.fetchAndSaveAllLists(); err != nil {
		zap.L().Error("获取列表失败", zap.Error(err))
		return err
	}

	// 第二步：更新详细信息
	// a. 更新电影详细信息
	if err := s.fetchAndUpdateMovieDetails(); err != nil {
		zap.L().Error("更新电影详情失败", zap.Error(err))
	}

	// b. 更新电视详细信息
	if err := s.fetchAndUpdateTVDetails(); err != nil {
		zap.L().Error("更新电视详情失败", zap.Error(err))
	}

	// c. 更新动漫详细信息（b执行完才能执行）
	if err := s.fetchAndUpdateAnimeDetails(); err != nil {
		zap.L().Error("更新动漫详情失败", zap.Error(err))
	}

	// d. 更新综艺详细信息（b执行完才能执行）
	if err := s.fetchAndUpdateShowDetails(); err != nil {
		zap.L().Error("更新综艺详情失败", zap.Error(err))
	}

	// e. 更新纪录片详细信息（b执行完才能执行）
	if err := s.fetchAndUpdateDocDetails(); err != nil {
		zap.L().Error("更新纪录片详情失败", zap.Error(err))
	}

	// 第三步：搜索播放地址并插入episodes表
	if err := s.searchAndSavePlayURLs(); err != nil {
		zap.L().Error("搜索播放地址失败", zap.Error(err))
	}

	// 第四步：更新存在 episodes 记录的 videos 的 status 为 1
	if err := s.updateVideosStatusByEpisodes(); err != nil {
		zap.L().Error("更新视频状态失败", zap.Error(err))
	}

	zap.L().Info("豆瓣数据同步完成")
	return nil
}

// fetchAndSaveAllLists 获取并保存所有列表
func (s *DoubanSyncService) fetchAndSaveAllLists() error {
	// 1. 最新电影列表
	if err := s.fetchAndSaveList(
		"https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie?start=0&limit=100&category=%E6%9C%80%E6%96%B0&type=%E5%85%A8%E9%83%A8",
		"https://movie.douban.com/explore",
		"movie",
		"", // 使用items.type
	); err != nil {
		zap.L().Error("获取电影列表失败", zap.Error(err))
	}

	// 2. 最新电视列表
	if err := s.fetchAndSaveList(
		"https://m.douban.com/rexxar/api/v2/subject/recent_hot/tv?start=0&limit=100&category=tv&type=tv",
		"https://movie.douban.com/tv/",
		"tv",
		"", // 使用items.type
	); err != nil {
		zap.L().Error("获取电视列表失败", zap.Error(err))
	}

	// 3. 动画列表
	if err := s.fetchAndSaveList(
		"https://m.douban.com/rexxar/api/v2/subject/recent_hot/tv?start=0&limit=200&category=tv&type=tv_animation",
		"https://movie.douban.com/tv/",
		"anime",
		"anime", // 固定为anime
	); err != nil {
		zap.L().Error("获取动画列表失败", zap.Error(err))
	}

	// 4. 纪录片列表
	if err := s.fetchAndSaveList(
		"https://m.douban.com/rexxar/api/v2/subject/recent_hot/tv?start=0&limit=200&category=tv&type=tv_documentary",
		"https://movie.douban.com/tv/",
		"doc",
		"doc", // 固定为doc，便于第二步区分
	); err != nil {
		zap.L().Error("获取纪录片列表失败", zap.Error(err))
	}

	// 5. 综艺列表
	if err := s.fetchAndSaveList(
		"https://m.douban.com/rexxar/api/v2/subject/recent_hot/tv?start=0&limit=200&category=show&type=show",
		"https://movie.douban.com/tv/",
		"tvshow",
		"tvshow", // 固定为tvshow
	); err != nil {
		zap.L().Error("获取综艺列表失败", zap.Error(err))
	}

	return nil
}

// fetchAndSaveList 获取并保存单个列表
func (s *DoubanSyncService) fetchAndSaveList(url, referer, defaultType, fixedType string) error {
	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("dnt", "1")
	req.Header.Set("origin", "https://movie.douban.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", referer)
	req.Header.Set("sec-ch-ua", `"Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
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
	var listResponse DoubanListResponse
	if err := json.Unmarshal(body, &listResponse); err != nil {
		return fmt.Errorf("解析JSON失败: %w", err)
	}

	zap.L().Info("获取到列表", zap.String("type", defaultType), zap.Int("count", len(listResponse.Items)))

	// 遍历列表，保存不存在的项
	savedCount := 0
	for _, item := range listResponse.Items {
		// 将字符串ID转换为整数
		sourceIDInt, err := strconv.Atoi(item.ID)
		if err != nil {
			zap.L().Warn("无效的ID", zap.String("id", item.ID))
			continue
		}
		sourceID := int64(sourceIDInt)

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

		// 确定type值
		videoType := fixedType
		if videoType == "" {
			videoType = item.Type
		}

		// 创建新视频记录
		score := item.Rating.Value
		video := &model.Video{
			ID:        utils.GenerateUserID(), // 使用雪花算法生成ID
			SourceID:  sourceID,
			Source:    "douban",
			Title:     item.Title,
			Type:      videoType,
			CoverURL:  item.Pic.Normal,
			Score:     &score,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.videoRepo.Create(video); err != nil {
			zap.L().Error("保存视频失败", zap.Error(err), zap.String("title", item.Title))
			continue
		}

		savedCount++
		zap.L().Info("保存新视频", zap.String("title", item.Title), zap.Int64("source_id", sourceID), zap.String("type", videoType))
	}

	zap.L().Info("列表同步完成", zap.String("type", defaultType), zap.Int("saved_count", savedCount))
	return nil
}

// fetchAndUpdateMovieDetails 获取并更新电影详情
func (s *DoubanSyncService) fetchAndUpdateMovieDetails() error {
	// 查找需要补充详情的电影（每次处理100条）
	videos, err := s.videoRepo.FindNeedDetailVideosByType("movie", 100)
	if err != nil {
		return fmt.Errorf("查询需要更新的电影失败: %w", err)
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

		// 避免请求过快，休眠4秒
		time.Sleep(4 * time.Second)
	}

	return nil
}

// fetchAndUpdateTVDetails 获取并更新电视详情
func (s *DoubanSyncService) fetchAndUpdateTVDetails() error {
	// 查找需要补充详情的电视（每次处理100条）
	videos, err := s.videoRepo.FindNeedDetailVideosByType("tv", 100)
	if err != nil {
		return fmt.Errorf("查询需要更新的电视失败: %w", err)
	}

	if len(videos) == 0 {
		zap.L().Info("没有需要更新详情的电视")
		return nil
	}

	zap.L().Info("找到需要更新详情的电视", zap.Int("count", len(videos)))

	// 遍历每个视频，获取详情
	for _, video := range videos {
		if err := s.fetchAndUpdateSingleTVDetail(video); err != nil {
			zap.L().Error("更新电视详情失败", zap.Error(err), zap.String("title", video.Title))
			continue
		}

		// 避免请求过快，休眠4秒
		time.Sleep(4 * time.Second)
	}

	return nil
}

// fetchAndUpdateAnimeDetails 获取并更新动漫详情
func (s *DoubanSyncService) fetchAndUpdateAnimeDetails() error {
	// 查找需要补充详情的动漫（type=anime且release_date和country_json都为空，每次处理100条）
	// 注意：第一步调用3时已经保存为type=anime，所以这里从anime查找
	videos, err := s.videoRepo.FindNeedDetailVideosByType("anime", 100)
	if err != nil {
		return fmt.Errorf("查询需要更新的动漫失败: %w", err)
	}

	if len(videos) == 0 {
		zap.L().Info("没有需要更新详情的动漫")
		return nil
	}

	zap.L().Info("找到需要更新详情的动漫", zap.Int("count", len(videos)))

	// 遍历每个视频，获取详情
	for _, video := range videos {
		if err := s.fetchAndUpdateSingleTVDetail(video); err != nil {
			zap.L().Error("更新动漫详情失败", zap.Error(err), zap.String("title", video.Title))
			continue
		}

		// type已经是anime，不需要更新
		video.CreatedAt = time.Now()

		if err := s.videoRepo.Update(video); err != nil {
			zap.L().Error("更新动漫详情失败", zap.Error(err), zap.String("title", video.Title))
			continue
		}

		// 避免请求过快，休眠4秒
		time.Sleep(4 * time.Second)
	}

	return nil
}

// fetchAndUpdateShowDetails 获取并更新综艺详情
func (s *DoubanSyncService) fetchAndUpdateShowDetails() error {
	// 查找需要补充详情的综艺（type=tvshow且release_date和country_json都为空，每次处理100条）
	// 注意：第一步调用5时已经保存为type=tvshow，所以这里从tvshow查找
	videos, err := s.videoRepo.FindNeedDetailVideosByType("tvshow", 100)
	if err != nil {
		return fmt.Errorf("查询需要更新的综艺失败: %w", err)
	}

	if len(videos) == 0 {
		zap.L().Info("没有需要更新详情的综艺")
		return nil
	}

	zap.L().Info("找到需要更新详情的综艺", zap.Int("count", len(videos)))

	// 遍历每个视频，获取详情
	for _, video := range videos {
		if err := s.fetchAndUpdateSingleShowDetail(video); err != nil {
			zap.L().Error("更新综艺详情失败", zap.Error(err), zap.String("title", video.Title))
			continue
		}

		// type已经是tvshow，不需要更新
		video.CreatedAt = time.Now()

		if err := s.videoRepo.Update(video); err != nil {
			zap.L().Error("更新综艺详情失败", zap.Error(err), zap.String("title", video.Title))
			continue
		}

		// 避免请求过快，休眠4秒
		time.Sleep(4 * time.Second)
	}

	return nil
}

// fetchAndUpdateDocDetails 获取并更新纪录片详情
func (s *DoubanSyncService) fetchAndUpdateDocDetails() error {
	// 查找需要补充详情的纪录片（type=doc且release_date和country_json都为空，每次处理100条）
	// 注意：第一步调用4时已经保存为type=doc，所以这里从doc查找
	videos, err := s.videoRepo.FindNeedDetailVideosByType("doc", 100)
	if err != nil {
		return fmt.Errorf("查询需要更新的纪录片失败: %w", err)
	}

	if len(videos) == 0 {
		zap.L().Info("没有需要更新详情的纪录片")
		return nil
	}

	zap.L().Info("找到需要更新详情的纪录片", zap.Int("count", len(videos)))

	// 遍历每个视频，获取详情
	for _, video := range videos {
		if err := s.fetchAndUpdateSingleDocDetail(video); err != nil {
			zap.L().Error("更新纪录片详情失败", zap.Error(err), zap.String("title", video.Title))
			continue
		}

		// type已经是doc，不需要更新
		video.CreatedAt = time.Now()

		if err := s.videoRepo.Update(video); err != nil {
			zap.L().Error("更新纪录片详情失败", zap.Error(err), zap.String("title", video.Title))
			continue
		}

		// 避免请求过快，休眠4秒
		time.Sleep(4 * time.Second)
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
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("dnt", "1")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("referer", "https://movie.douban.com/explore?support_type=movie&is_all=false&category=%E8%B1%86%E7%93%A3%E9%AB%98%E5%88%86&type=%E5%85%A8%E9%83%A8")
	req.Header.Set("sec-ch-ua", `"Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
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
	directorStr := extractFieldWithAttrs(html, "导演")
	actorsStr := extractFieldWithAttrs(html, "主演")
	tagsStr := extractGenres(html)
	countryStr := extractField(html, `<span class="pl">制片国家/地区:</span>`, `<br`)

	// 转换为JSON数组并截断到512字节以内
	if directorJSON, err := stringToJSONArray(directorStr); err == nil {
		video.DirectorJSON = truncateJSONArray(directorJSON)
	}
	if actorsJSON, err := stringToJSONArray(actorsStr); err == nil {
		video.ActorsJSON = truncateJSONArray(actorsJSON)
	}
	if tagsJSON, err := stringToJSONArray(tagsStr); err == nil {
		video.TagsJSON = truncateJSONArray(tagsJSON)
	}
	if countryJSON, err := stringToCountryJSONArray(countryStr); err == nil {
		video.CountryJSON = truncateJSONArray(countryJSON)
	}

	// 提取评分
	if score := extractScore(html); score != nil {
		video.Score = score
	}

	// 提取上映日期（完整日期）
	dateStr := extractField(html, `<span class="pl">上映日期:</span>`, `<br`)
	video.ReleaseDate = parseDateString(dateStr)

	// 提取片长（只保留数字）
	runtimeStr := extractField(html, `<span class="pl">片长:</span>`, `<br`)
	if runtimeStr != "" {
		runtime := int64(extractNumber(runtimeStr))
		video.Runtime = &runtime
	}

	// 提取IMDb ID
	video.IMDbID = extractIMDbID(html)

	// 提取简介
	video.Description = extractDescription(html)

	// 设置集数为0（电影）
	episodeCount := int64(0)
	video.EpisodeCount = &episodeCount

	// 设置创建时间
	video.CreatedAt = time.Now()

	// 更新时间
	video.UpdatedAt = time.Now()

	// 保存到数据库
	if err := s.videoRepo.Update(video); err != nil {
		return fmt.Errorf("更新数据库失败: %w", err)
	}

	zap.L().Info("更新电影详情成功", zap.String("title", video.Title), zap.Int64("source_id", video.SourceID))
	return nil
}

// fetchAndUpdateSingleTVDetail 获取并更新单个电视详情
func (s *DoubanSyncService) fetchAndUpdateSingleTVDetail(video *model.Video) error {
	url := fmt.Sprintf("https://movie.douban.com/subject/%d/", video.SourceID)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头（与电影详情相同）
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("dnt", "1")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("referer", "https://movie.douban.com/tv/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
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
	directorStr := extractFieldWithAttrs(html, "导演")
	actorsStr := extractFieldWithAttrs(html, "主演")
	tagsStr := extractGenres(html)
	countryStr := extractField(html, `<span class="pl">制片国家/地区:</span>`, `<br`)

	// 转换为JSON数组并截断到512字节以内
	if directorJSON, err := stringToJSONArray(directorStr); err == nil {
		video.DirectorJSON = truncateJSONArray(directorJSON)
	}
	if actorsJSON, err := stringToJSONArray(actorsStr); err == nil {
		video.ActorsJSON = truncateJSONArray(actorsJSON)
	}
	if tagsJSON, err := stringToJSONArray(tagsStr); err == nil {
		video.TagsJSON = truncateJSONArray(tagsJSON)
	}
	if countryJSON, err := stringToCountryJSONArray(countryStr); err == nil {
		video.CountryJSON = truncateJSONArray(countryJSON)
	}

	// 提取评分
	if score := extractScore(html); score != nil {
		video.Score = score
	}

	// 提取首播日期（完整日期）
	dateStr := extractField(html, `<span class="pl">首播:</span>`, `<br`)
	video.ReleaseDate = parseDateString(dateStr)

	// 提取集数（只保留数字）
	episodeStr := extractField(html, `<span class="pl">集数:</span>`, `<br`)
	if episodeStr != "" {
		episodeCount := int64(extractNumber(episodeStr))
		video.EpisodeCount = &episodeCount
	}

	// 提取IMDb ID
	video.IMDbID = extractIMDbID(html)

	// 提取简介
	video.Description = extractDescription(html)

	// 设置创建时间
	video.CreatedAt = time.Now()

	// 更新时间
	video.UpdatedAt = time.Now()

	// 保存到数据库
	if err := s.videoRepo.Update(video); err != nil {
		return fmt.Errorf("更新数据库失败: %w", err)
	}

	zap.L().Info("更新电视详情成功", zap.String("title", video.Title), zap.Int64("source_id", video.SourceID))
	return nil
}

// fetchAndUpdateSingleShowDetail 获取并更新单个综艺详情
func (s *DoubanSyncService) fetchAndUpdateSingleShowDetail(video *model.Video) error {
	url := fmt.Sprintf("https://movie.douban.com/subject/%d/", video.SourceID)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头（与电视详情相同）
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("dnt", "1")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("referer", "https://movie.douban.com/tv/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
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

	// 解析HTML，提取信息（综艺没有导演）
	actorsStr := extractFieldWithAttrs(html, "主演")
	tagsStr := extractGenres(html)
	countryStr := extractField(html, `<span class="pl">制片国家/地区:</span>`, `<br`)

	// 转换为JSON数组并截断到512字节以内
	if actorsJSON, err := stringToJSONArray(actorsStr); err == nil {
		video.ActorsJSON = truncateJSONArray(actorsJSON)
	}
	if tagsJSON, err := stringToJSONArray(tagsStr); err == nil {
		video.TagsJSON = truncateJSONArray(tagsJSON)
	}
	if countryJSON, err := stringToCountryJSONArray(countryStr); err == nil {
		video.CountryJSON = truncateJSONArray(countryJSON)
	}

	// 提取评分
	if score := extractScore(html); score != nil {
		video.Score = score
	}

	// 提取首播日期（完整日期）
	dateStr := extractField(html, `<span class="pl">首播:</span>`, `<br`)
	video.ReleaseDate = parseDateString(dateStr)

	// 提取集数（只保留数字）
	episodeStr := extractField(html, `<span class="pl">集数:</span>`, `<br`)
	if episodeStr != "" {
		episodeCount := int64(extractNumber(episodeStr))
		video.EpisodeCount = &episodeCount
	}

	// 提取简介
	video.Description = extractDescription(html)

	// 设置创建时间
	video.CreatedAt = time.Now()

	// 更新时间
	video.UpdatedAt = time.Now()

	// 保存到数据库
	if err := s.videoRepo.Update(video); err != nil {
		return fmt.Errorf("更新数据库失败: %w", err)
	}

	zap.L().Info("更新综艺详情成功", zap.String("title", video.Title), zap.Int64("source_id", video.SourceID))
	return nil
}

// fetchAndUpdateSingleDocDetail 获取并更新单个纪录片详情
func (s *DoubanSyncService) fetchAndUpdateSingleDocDetail(video *model.Video) error {
	url := fmt.Sprintf("https://movie.douban.com/subject/%d/", video.SourceID)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头（与电视详情相同）
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("dnt", "1")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("referer", "https://movie.douban.com/tv/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
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

	// 解析HTML，提取信息（纪录片没有导演和主演）
	tagsStr := extractGenres(html)
	countryStr := extractField(html, `<span class="pl">制片国家/地区:</span>`, `<br`)

	// 转换为JSON数组并截断到512字节以内
	if tagsJSON, err := stringToJSONArray(tagsStr); err == nil {
		video.TagsJSON = truncateJSONArray(tagsJSON)
	}
	if countryJSON, err := stringToCountryJSONArray(countryStr); err == nil {
		video.CountryJSON = truncateJSONArray(countryJSON)
	}

	// 提取评分
	if score := extractScore(html); score != nil {
		video.Score = score
	}

	// 提取首播日期（完整日期）
	dateStr := extractField(html, `<span class="pl">首播:</span>`, `<br`)
	video.ReleaseDate = parseDateString(dateStr)

	// 提取集数（只保留数字）
	episodeStr := extractField(html, `<span class="pl">集数:</span>`, `<br`)
	if episodeStr != "" {
		episodeCount := int64(extractNumber(episodeStr))
		video.EpisodeCount = &episodeCount
	}

	// 提取简介
	video.Description = extractDescription(html)

	// 设置创建时间
	video.CreatedAt = time.Now()

	// 更新时间
	video.UpdatedAt = time.Now()

	// 保存到数据库
	if err := s.videoRepo.Update(video); err != nil {
		return fmt.Errorf("更新数据库失败: %w", err)
	}

	zap.L().Info("更新纪录片详情成功", zap.String("title", video.Title), zap.Int64("source_id", video.SourceID))
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

// extractScore 提取评分
func extractScore(html string) *float64 {
	// 查找 <strong class="ll rating_num" property="v:average">9.7</strong>
	// 或者 <span class="rating_num">9.7</span>
	patterns := []string{
		`<strong[^>]*class="[^"]*rating_num[^"]*"[^>]*property="v:average"[^>]*>([0-9.]+)</strong>`,
		`<span[^>]*class="[^"]*rating_num[^"]*"[^>]*>([0-9.]+)</span>`,
		`property="v:average"[^>]*>([0-9.]+)</strong>`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) > 1 {
			scoreStr := strings.TrimSpace(matches[1])
			if score, err := strconv.ParseFloat(scoreStr, 64); err == nil && score >= 0 && score <= 10 {
				return &score
			}
		}
	}
	return nil
}

// parseDateString 解析日期字符串，支持多种格式
// 支持的格式：
// - "2025-01-07"
// - "2025-01-07(中国大陆)"
// - "2025年1月7日"
// - "2025-1-7"
// - "2025年01月07日"
// - "2025" (仅年份，转换为该年1月1日)
func parseDateString(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}

	// 清理字符串，移除HTML标签和多余空白
	dateStr = removeHTMLTags(dateStr)
	dateStr = strings.TrimSpace(dateStr)

	// 如果包含括号，提取括号前的内容（例如："2025-01-07(中国大陆)" -> "2025-01-07"）
	if idx := strings.Index(dateStr, "("); idx != -1 {
		dateStr = strings.TrimSpace(dateStr[:idx])
	}

	// 尝试解析各种日期格式
	dateFormats := []string{
		"2006-01-02",  // 标准格式：2025-01-07
		"2006-1-2",    // 无前导零：2025-1-7
		"2006年01月02日", // 中文格式：2025年01月07日
		"2006年1月2日",   // 中文格式无前导零：2025年1月7日
		"2006年01月2日",  // 混合格式
		"2006年1月02日",  // 混合格式
		"2006",        // 仅年份
	}

	for _, format := range dateFormats {
		if t, err := time.Parse(format, dateStr); err == nil {
			// 如果只解析到年份，设置为该年1月1日
			if format == "2006" {
				t = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
			}
			// 验证日期范围合理
			if t.Year() >= 1900 && t.Year() <= 2100 {
				return &t
			}
		}
	}

	// 如果所有格式都失败，尝试提取年份
	re := regexp.MustCompile(`(\d{4})`)
	matches := re.FindStringSubmatch(dateStr)
	if len(matches) > 1 {
		if year, err := strconv.Atoi(matches[1]); err == nil {
			if year >= 1900 && year <= 2100 {
				date := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
				return &date
			}
		}
	}

	return nil
}

// extractDescription 提取简介
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

// stringToJSONArray 将逗号分隔的字符串转换为JSON数组
func stringToJSONArray(str string) ([]byte, error) {
	if str == "" {
		return []byte("[]"), nil
	}
	// 分割字符串并清理空白
	parts := strings.Split(str, ",")
	var items []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return json.Marshal(items)
}

// stringToCountryJSONArray 将国家/地区字符串转换为JSON数组
// 支持多种分隔符：", "、","、" / "、"/"
// 每个国家作为一个独立的JSON值存储，如：["中国大陆","美国"]
func stringToCountryJSONArray(str string) ([]byte, error) {
	if str == "" {
		return []byte("[]"), nil
	}

	// 先替换 " / " 为 ", "，统一分隔符
	str = strings.ReplaceAll(str, " / ", ", ")
	str = strings.ReplaceAll(str, "/", ",")

	// 分割字符串并清理空白
	parts := strings.Split(str, ",")
	var items []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return json.Marshal(items)
}

// truncateJSONArray 截断JSON数组，确保不超过512字节
// 如果超过512字节，循环删除最后一个元素，直到长度低于512字节
func truncateJSONArray(jsonData []byte) []byte {
	const maxSize = 512

	// 如果已经小于等于512字节，直接返回
	if len(jsonData) <= maxSize {
		return jsonData
	}

	// 解析JSON数组
	var items []string
	if err := json.Unmarshal(jsonData, &items); err != nil {
		// 如果解析失败，返回空数组
		return []byte("[]")
	}

	// 循环删除最后一个元素，直到长度低于512字节
	for len(items) > 0 {
		// 重新序列化
		truncated, err := json.Marshal(items)
		if err != nil {
			// 如果序列化失败，返回空数组
			return []byte("[]")
		}

		// 如果长度符合要求，返回
		if len(truncated) <= maxSize {
			return truncated
		}

		// 删除最后一个元素
		items = items[:len(items)-1]
	}

	// 如果所有元素都被删除，返回空数组
	return []byte("[]")
}

// searchAndSavePlayURLs 搜索播放地址并保存到episodes表
func (s *DoubanSyncService) searchAndSavePlayURLs() error {
	zap.L().Info("开始搜索播放地址")

	// 查询 status 不等于 0 和 1 的视频的 id、type、title（用于更新episodes）
	videos, err := s.videoRepo.FindVideosNeedUpdateEpisodes()
	if err != nil {
		return fmt.Errorf("查询视频列表失败: %w", err)
	}

	if len(videos) == 0 {
		zap.L().Info("没有需要搜索播放地址的视频")
		return nil
	}

	zap.L().Info("找到需要搜索播放地址的视频", zap.Int("count", len(videos)))

	// 遍历每个视频，搜索播放地址
	for _, video := range videos {
		if video.Title == "" {
			continue
		}

		if err := s.searchAndSavePlayURLsForVideo(video); err != nil {
			zap.L().Error("搜索播放地址失败", zap.Error(err), zap.String("title", video.Title), zap.Int64("id", video.ID))
			continue
		}

		// 避免请求过快，休眠1秒
		time.Sleep(1 * time.Second)
	}

	zap.L().Info("播放地址搜索完成")
	return nil
}

// searchAndSavePlayURLsForVideo 为单个视频搜索播放地址并保存
func (s *DoubanSyncService) searchAndSavePlayURLsForVideo(video *model.Video) error {
	// 构建搜索URL，使用title替换q参数
	searchURL := fmt.Sprintf("http://124.222.196.128:3000/api/search?q=%s", url.QueryEscape(video.Title))

	// 创建HTTP请求
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")

	// 设置Cookie
	req.Header.Set("Cookie", "auth=%257B%2522role%2522%253A%2522user%2522%252C%2522password%2522%253A%252212345%2522%257D")

	// 创建HTTP客户端（跳过SSL验证）
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// 发送请求
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

	// 解析JSON响应
	var searchResponse SearchResponse
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return fmt.Errorf("解析JSON失败: %w", err)
	}

	// 遍历搜索结果，只处理第一个匹配的 result
	for _, result := range searchResponse.Results {
		// 判断 results.title = videos.title
		if result.Title != video.Title {
			continue
		}

		// 根据 type 区分处理逻辑
		if video.Type == "movie" {
			// movie 类型：优先获取包含 "vip" 的项，如果没有则获取包含 "ryplay7" 的项，如果都没有则按顺序取第一个，只取第一行
			var selectedEpisode string
			found := false

			// 1. 优先查找包含 "vip" 的项
			for _, episode := range result.Episodes {
				if strings.Contains(strings.ToLower(episode), "vip") {
					selectedEpisode = episode
					found = true
					break
				}
			}

			// 2. 如果没有找到 vip，查找包含 "ryplay7" 的项
			if !found {
				for _, episode := range result.Episodes {
					if strings.Contains(strings.ToLower(episode), "ryplay7") {
						selectedEpisode = episode
						found = true
						break
					}
				}
			}

			// 3. 如果都没有，按顺序取第一个
			if !found && len(result.Episodes) > 0 {
				selectedEpisode = result.Episodes[0]
				found = true
			}

			// 如果没有找到任何项，跳过
			if !found {
				continue
			}

			// 将episode值按行分割（支持\n和\r\n）
			lines := strings.Split(strings.ReplaceAll(selectedEpisode, "\r\n", "\n"), "\n")
			// movie 类型：只取第一行
			var firstLine string
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					firstLine = line
					break
				}
			}

			if firstLine == "" {
				continue
			}

			// 限制播放地址长度不超过255字符
			playURL := firstLine
			if len(playURL) > 255 {
				playURL = playURL[:255]
				zap.L().Warn("播放地址长度超过255字符，已截断", zap.String("original", firstLine), zap.String("truncated", playURL))
			}

			// 创建episode记录
			episodeNumber := int64(1)
			episode := &model.Episode{
				Channel:         result.SourceName,
				ChannelID:       nil, // channel_id 为 null
				VideoID:         video.ID,
				EpisodeNumber:   &episodeNumber,
				Name:            result.Title,
				PlayURLs:        playURL,
				DurationSeconds: nil, // duration_seconds 为 null
				SubtitleURLs:    nil, // subtitle_urls 为 null
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}

			// 插入到数据库
			if err := s.episodeRepo.Create(episode); err != nil {
				zap.L().Error("插入episode失败", zap.Error(err), zap.String("title", result.Title), zap.Int64("video_id", video.ID))
				return fmt.Errorf("插入episode失败: %w", err)
			}

			zap.L().Info("插入episode成功", zap.String("title", result.Title), zap.Int64("video_id", video.ID), zap.Int64("episode_number", episodeNumber), zap.String("play_url", firstLine))
		} else {
			// 非 movie 类型：episodes 本身就是一个字符串数组
			// 1. 找到匹配的 episode（优先包含 "vip"，其次 "ryplay7"，否则取第一个）
			// 2. 实现增量更新逻辑
			var selectedEpisodes []string
			found := false

			// 1. 优先查找包含 "vip" 的项
			for _, episode := range result.Episodes {
				if strings.Contains(strings.ToLower(episode), "vip") {
					// episodes 本身就是数组，直接使用整个数组
					selectedEpisodes = result.Episodes
					found = true
					break
				}
			}

			// 2. 如果没有找到 vip，查找包含 "ryplay7" 的项
			if !found {
				for _, episode := range result.Episodes {
					if strings.Contains(strings.ToLower(episode), "ryplay7") {
						// episodes 本身就是数组，直接使用整个数组
						selectedEpisodes = result.Episodes
						found = true
						break
					}
				}
			}

			// 3. 如果都没有，按顺序取第一个（使用整个数组）
			if !found && len(result.Episodes) > 0 {
				selectedEpisodes = result.Episodes
				found = true
			}

			// 如果没有找到任何项，跳过
			if !found || len(selectedEpisodes) == 0 {
				continue
			}

			// 检查videos.id在episodes表的video_id是否存在
			exists, err := s.episodeRepo.ExistsByVideoID(video.ID)
			if err != nil {
				zap.L().Error("检查episode是否存在失败", zap.Error(err), zap.Int64("video_id", video.ID))
			} else if exists {
				// 如果存在，更新videos表status的值为1
				if err := s.videoRepo.UpdateVideoStatus(video.ID, "1"); err != nil {
					zap.L().Error("更新视频status失败", zap.Error(err), zap.Int64("video_id", video.ID))
				} else {
					zap.L().Info("更新视频status为1", zap.Int64("video_id", video.ID))
				}
			}

			// 获取当前已存在的episode数量
			existingCount, err := s.episodeRepo.CountByVideoID(video.ID)
			if err != nil {
				zap.L().Error("统计episode数量失败", zap.Error(err), zap.Int64("video_id", video.ID))
				existingCount = 0
			}

			// 源数据总条数
			sourceCount := int64(len(selectedEpisodes))

			// 如果源数据条数大于已存在的episodes数量，则增量更新
			if sourceCount > existingCount {
				// 从existingCount+1开始，增量插入新的episodes
				startIndex := int(existingCount)
				newEpisodesCount := 0

				for i := startIndex; i < len(selectedEpisodes); i++ {
					episodeValue := strings.TrimSpace(selectedEpisodes[i])
					if episodeValue == "" {
						continue
					}

					// 限制播放地址长度不超过255字符
					playURL := episodeValue
					if len(playURL) > 255 {
						playURL = playURL[:255]
						zap.L().Warn("播放地址长度超过255字符，已截断", zap.String("original", episodeValue), zap.String("truncated", playURL))
					}

					// 创建episode记录，episode_number从existingCount+1开始
					episodeNumber := int64(i + 1)
					episode := &model.Episode{
						Channel:         result.SourceName,
						ChannelID:       nil, // channel_id 为 null
						VideoID:         video.ID,
						EpisodeNumber:   &episodeNumber,
						Name:            result.Title,
						PlayURLs:        playURL,
						DurationSeconds: nil, // duration_seconds 为 null
						SubtitleURLs:    nil, // subtitle_urls 为 null
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					}

					// 插入到数据库
					if err := s.episodeRepo.Create(episode); err != nil {
						zap.L().Error("插入episode失败", zap.Error(err), zap.String("title", result.Title), zap.Int64("video_id", video.ID), zap.Int64("episode_number", episodeNumber))
						continue
					}

					newEpisodesCount++
					zap.L().Info("插入episode成功", zap.String("title", result.Title), zap.Int64("video_id", video.ID), zap.Int64("episode_number", episodeNumber), zap.String("play_url", episodeValue))
				}

				// 如果有新增episodes，更新videos表的updated_at为当前时间
				if newEpisodesCount > 0 {
					if err := database.DB.Model(&model.Video{}).
						Where("id = ?", video.ID).
						Update("updated_at", time.Now()).Error; err != nil {
						zap.L().Error("更新视频updated_at失败", zap.Error(err), zap.Int64("video_id", video.ID))
					}
				}
			}

			// 检查最后一条episode记录的created_at是否为最近三天
			lastEpisode, err := s.episodeRepo.FindLastByVideoID(video.ID)
			var isUpdate bool
			if err == nil && lastEpisode != nil {
				threeDaysAgo := time.Now().AddDate(0, 0, -3)
				isUpdate = lastEpisode.CreatedAt.After(threeDaysAgo)
			} else {
				// 如果没有episode记录，is_update设为false
				isUpdate = false
			}
			// 更新is_update字段
			if err := s.videoRepo.UpdateVideoIsUpdate(video.ID, isUpdate); err != nil {
				zap.L().Error("更新视频is_update失败", zap.Error(err), zap.Int64("video_id", video.ID))
			} else {
				zap.L().Info("更新视频is_update", zap.Int64("video_id", video.ID), zap.Bool("is_update", isUpdate))
			}

			// 获取视频完整信息，检查episode_count
			videoInfo, err := s.videoRepo.FindByID(video.ID)
			if err == nil && videoInfo != nil {
				// 获取当前episodes总数
				currentCount, err := s.episodeRepo.CountByVideoID(video.ID)
				if err == nil && videoInfo.EpisodeCount != nil {
					// 如果episodes总数等于episode_count，则is_completed为1
					isCompleted := currentCount == *videoInfo.EpisodeCount
					if err := s.videoRepo.UpdateVideoIsCompleted(video.ID, isCompleted); err != nil {
						zap.L().Error("更新视频is_completed失败", zap.Error(err), zap.Int64("video_id", video.ID))
					} else {
						zap.L().Info("更新视频is_completed", zap.Int64("video_id", video.ID), zap.Bool("is_completed", isCompleted), zap.Int64("current_count", currentCount), zap.Int64("episode_count", *videoInfo.EpisodeCount))
					}
				}
			}
		}

		// 只处理第一个匹配的 result，处理完就退出
		break
	}

	return nil
}

// updateVideosStatusByEpisodes 更新存在 episodes 记录的 videos 的 status 为 1
func (s *DoubanSyncService) updateVideosStatusByEpisodes() error {
	zap.L().Info("开始更新视频状态")

	if err := s.videoRepo.UpdateVideosStatusByEpisodes("1"); err != nil {
		return fmt.Errorf("更新视频状态失败: %w", err)
	}

	zap.L().Info("视频状态更新完成")
	return nil
}
