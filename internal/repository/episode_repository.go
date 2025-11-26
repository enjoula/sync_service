// repository 包提供数据访问层，封装数据库操作
package repository

import (
	"video-service/internal/model"
	"video-service/pkg/infrastructure/database"
)

// EpisodeRepository 剧集仓库接口
type EpisodeRepository interface {
	// Create 创建剧集记录
	Create(episode *model.Episode) error

	// FindByVideoID 根据视频ID查找所有剧集
	FindByVideoID(videoID int64) ([]*model.Episode, error)
}

// episodeRepository 剧集仓库实现
type episodeRepository struct{}

// NewEpisodeRepository 创建剧集仓库实例
func NewEpisodeRepository() EpisodeRepository {
	return &episodeRepository{}
}

// Create 创建剧集记录
func (r *episodeRepository) Create(episode *model.Episode) error {
	return database.DB.Create(episode).Error
}

// FindByVideoID 根据视频ID查找所有剧集
func (r *episodeRepository) FindByVideoID(videoID int64) ([]*model.Episode, error) {
	var episodes []*model.Episode
	err := database.DB.Where("video_id = ?", videoID).Find(&episodes).Error
	if err != nil {
		return nil, err
	}
	return episodes, nil
}

