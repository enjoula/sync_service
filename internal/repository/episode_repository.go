// repository 包提供数据访问层，封装数据库操作
package repository

import (
	"time"
	"video-service/internal/model"
	"video-service/pkg/infrastructure/database"

	"gorm.io/gorm"
)

// EpisodeRepository 剧集仓库接口
type EpisodeRepository interface {
	// Create 创建剧集记录
	Create(episode *model.Episode) error

	// FindByVideoID 根据视频ID查找所有剧集
	FindByVideoID(videoID int64) ([]*model.Episode, error)

	// CountByVideoID 根据视频ID统计episode数量
	CountByVideoID(videoID int64) (int64, error)

	// FindLastByVideoID 根据视频ID查找最后一条episode记录（按created_at降序）
	FindLastByVideoID(videoID int64) (*model.Episode, error)

	// ExistsByVideoID 检查视频ID是否存在episode记录
	ExistsByVideoID(videoID int64) (bool, error)
}

// episodeRepository 剧集仓库实现
type episodeRepository struct{}

// NewEpisodeRepository 创建剧集仓库实例
func NewEpisodeRepository() EpisodeRepository {
	return &episodeRepository{}
}

// Create 创建剧集记录，同时更新对应视频的updated_at字段
func (r *episodeRepository) Create(episode *model.Episode) error {
	// 使用事务确保原子性
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 插入episode记录
		if err := tx.Create(episode).Error; err != nil {
			return err
		}

		// 2. 更新对应video的updated_at字段为当前时间
		if err := tx.Model(&model.Video{}).
			Where("id = ?", episode.VideoID).
			Update("updated_at", time.Now()).Error; err != nil {
			return err
		}

		return nil
	})
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

// CountByVideoID 根据视频ID统计episode数量
func (r *episodeRepository) CountByVideoID(videoID int64) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Episode{}).
		Where("video_id = ?", videoID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// FindLastByVideoID 根据视频ID查找最后一条episode记录（按created_at降序）
func (r *episodeRepository) FindLastByVideoID(videoID int64) (*model.Episode, error) {
	var episode model.Episode
	err := database.DB.Where("video_id = ?", videoID).
		Order("created_at DESC").
		First(&episode).Error
	if err != nil {
		return nil, err
	}
	return &episode, nil
}

// ExistsByVideoID 检查视频ID是否存在episode记录
func (r *episodeRepository) ExistsByVideoID(videoID int64) (bool, error) {
	var count int64
	err := database.DB.Model(&model.Episode{}).
		Where("video_id = ?", videoID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
