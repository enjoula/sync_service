package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"video-service/internal/model"
	"video-service/internal/pkg/utils"
	"video-service/pkg/infrastructure/config"
	"video-service/pkg/infrastructure/database"
	"video-service/pkg/infrastructure/logger"

	"gorm.io/gorm"
)

// DoubanTVShowItem 豆瓣剧集列表项
type DoubanTVShowItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Rating struct {
		Value float64 `json:"value"`
	} `json:"rating"`
	Pic struct {
		Normal string `json:"normal"`
	} `json:"pic"`
}

// DoubanTVShowResponse 豆瓣剧集列表响应
type DoubanTVShowResponse struct {
	Items []DoubanTVShowItem `json:"items"`
}

func main() {
	log.Println("=== 开始同步豆瓣热门剧集 ===")

	// 初始化日志系统
	logger.InitLogger()

	// 初始化配置
	config.InitConfig()

	// 如果设置了环境变量 MYSQL_DSN，则覆盖配置文件中的数据库连接
	// 例如: export MYSQL_DSN="root:123456@tcp(localhost:3306)/video_service?charset=utf8mb4&parseTime=True&loc=Local"
	if mysqlDsn := os.Getenv("MYSQL_DSN"); mysqlDsn != "" {
		log.Printf("使用环境变量中的数据库配置")
		config.Cfg.Set("mysql.dsn", mysqlDsn)
	}

	// 初始化数据库
	database.InitMySQL()

	// 检查数据库是否连接成功
	if database.DB == nil {
		log.Fatalf("数据库连接失败，请检查配置或设置环境变量 MYSQL_DSN")
	}

	// 调用豆瓣API获取热门剧集
	url := "https://m.douban.com/rexxar/api/v2/subject/recent_hot/tv?start=0&limit=100&category=show&type=show"

	log.Printf("正在请求URL: %s", url)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("dnt", "1")
	req.Header.Set("origin", "https://movie.douban.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://movie.douban.com/tv/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", `ll="118172"; bid=YpzrGzsc_RI; __utmc=30149280; _vwo_uuid_v2=D10E35F36802D61AC27D96415ADCCEA16|ff5433783190543ff19fd731261905f5; __utmz=30149280.1763582035.22.4.utmcsr=cn.bing.com|utmccn=(referral)|utmcmd=referral|utmcct=/; ap_v=0,6.0; __utma=30149280.1857371696.1762958131.1763582035.1763629051.23; __utmb=30149280.0.10.1763629051`)

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("响应状态码: %d", resp.StatusCode)

	if resp.StatusCode != 200 {
		log.Fatalf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应失败: %v", err)
	}

	// 解析JSON
	var tvShowResponse DoubanTVShowResponse
	if err := json.Unmarshal(body, &tvShowResponse); err != nil {
		log.Fatalf("解析JSON失败: %v", err)
	}

	log.Printf("成功获取 %d 个剧集", len(tvShowResponse.Items))

	// 统计信息
	totalCount := len(tvShowResponse.Items)
	existCount := 0
	insertCount := 0
	errorCount := 0

	// 遍历所有items，检查并插入数据库
	for i, item := range tvShowResponse.Items {
		log.Printf("[%d/%d] 处理剧集: %s (ID: %s)", i+1, totalCount, item.Title, item.ID)

		// 将字符串ID转换为整数
		sourceIDInt, err := strconv.Atoi(item.ID)
		if err != nil {
			log.Printf("  ✗ 无效的剧集ID: %s, 错误: %v", item.ID, err)
			errorCount++
			continue
		}
		sourceID := int64(sourceIDInt)

		// 检查数据库中是否已存在该source_id
		var existingVideo model.Video
		err = database.DB.Where("source_id = ?", sourceID).First(&existingVideo).Error

		if err == nil {
			// 记录已存在
			log.Printf("  ○ 剧集已存在，跳过 (source_id: %d)", sourceID)
			existCount++
			continue
		}

		if err != gorm.ErrRecordNotFound {
			// 数据库查询错误
			log.Printf("  ✗ 查询数据库失败: %v", err)
			errorCount++
			continue
		}

		// 创建新的视频记录
		score := item.Rating.Value
		video := &model.Video{
			ID:        utils.GenerateUserID(), // 使用雪花算法生成ID
			SourceID:  sourceID,
			Source:    "douban",
			Title:     item.Title,
			Type:      "tvshow",
			CoverURL:  item.Pic.Normal,
			Score:     &score,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// 插入数据库
		if err := database.DB.Create(video).Error; err != nil {
			log.Printf("  ✗ 插入数据库失败: %v", err)
			errorCount++
			continue
		}

		scoreStr := "N/A"
		if video.Score != nil {
			scoreStr = fmt.Sprintf("%.1f", *video.Score)
		}
		log.Printf("  ✓ 成功插入新剧集 (ID: %d, source_id: %d, score: %s)", video.ID, video.SourceID, scoreStr)
		insertCount++
	}

	// 输出统计信息
	log.Println("\n=== 同步完成 ===")
	log.Printf("总数: %d", totalCount)
	log.Printf("已存在: %d", existCount)
	log.Printf("新插入: %d", insertCount)
	log.Printf("错误: %d", errorCount)
}
