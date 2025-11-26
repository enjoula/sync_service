// repository 包提供数据访问层，封装数据库操作
package repository

import (
	"video-service/internal/model"
	"video-service/pkg/infrastructure/database"
)

// VideoRepository 视频仓库接口
type VideoRepository interface {
	// FindBySourceID 根据来源ID查找视频
	FindBySourceID(sourceID int64) (*model.Video, error)

	// Create 创建视频记录
	Create(video *model.Video) error

	// Update 更新视频记录
	Update(video *model.Video) error

	// FindNeedDetailVideos 查找需要补充详情的视频（source_id不为空且release_date和country_json都为空）
	FindNeedDetailVideos(limit int) ([]*model.Video, error)

	// FindNeedDetailVideosByType 根据类型查找需要补充详情的视频（source_id不为空且release_date和country_json都为空）
	FindNeedDetailVideosByType(videoType string, limit int) ([]*model.Video, error)

	// FindAllVideos 查找所有视频（仅返回 id 和 title）
	FindAllVideos() ([]*model.Video, error)
}

// videoRepository 视频仓库实现
type videoRepository struct{}

// NewVideoRepository 创建视频仓库实例
func NewVideoRepository() VideoRepository {
	return &videoRepository{}
}

// FindBySourceID 根据来源ID查找视频
func (r *videoRepository) FindBySourceID(sourceID int64) (*model.Video, error) {
	var video model.Video
	err := database.DB.Where("source_id = ?", sourceID).First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// Create 创建视频记录
func (r *videoRepository) Create(video *model.Video) error {
	return database.DB.Create(video).Error
}

// Update 更新视频记录
func (r *videoRepository) Update(video *model.Video) error {
	return database.DB.Save(video).Error
}

// FindNeedDetailVideos 查找需要补充详情的视频
// 条件：source_id不为空且country_json为空（country_json为空说明详情未获取）
func (r *videoRepository) FindNeedDetailVideos(limit int) ([]*model.Video, error) {
	var videos []*model.Video
	err := database.DB.Where("source_id IS NOT NULL AND source_id != 0 AND (country_json IS NULL OR country_json = '[]')").
		Limit(limit).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// FindNeedDetailVideosByType 根据类型查找需要补充详情的视频
// 条件：source_id不为空且release_date和country_json都为空
func (r *videoRepository) FindNeedDetailVideosByType(videoType string, limit int) ([]*model.Video, error) {
	var videos []*model.Video
	err := database.DB.Where("source_id IS NOT NULL AND source_id != 0 AND type = ? AND (release_date IS NULL) AND (country_json IS NULL OR country_json = '[]')", videoType).
		Limit(limit).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// FindAllVideos 查找所有视频（仅返回 id 和 title）
func (r *videoRepository) FindAllVideos() ([]*model.Video, error) {
	var videos []*model.Video
	err := database.DB.Select("id", "title").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}
