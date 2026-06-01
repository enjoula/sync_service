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

	// FindByID 根据ID查找剧集
	FindByID(episodeID int64) (*model.Episode, error)

	// Update 更新剧集记录
	Update(episode *model.Episode) error

	// DeleteByID 根据ID删除剧集
	DeleteByID(episodeID int64) error

	// DeleteByVideoID 根据视频ID删除所有剧集
	DeleteByVideoID(videoID int64) error
}

// episodeRepository 剧集仓库实现
type episodeRepository struct{}

// NewEpisodeRepository 创建剧集仓库实例
func NewEpisodeRepository() EpisodeRepository {
	return &episodeRepository{}
}

// Create 创建剧集记录，同时更新对应视频的updated_at字段
func (r *episodeRepository) Create(episode *model.Episode) error {
	// 根据零值收集需要省略的字段，避免写入空数据
	omitFields := make([]string, 0, 5)
	if episode.Channel == "" {
		omitFields = append(omitFields, "Channel")
	}
	if episode.ChannelID == 0 {
		omitFields = append(omitFields, "ChannelID")
	}
	if episode.Name == "" {
		omitFields = append(omitFields, "Name")
	}
	if episode.DurationSeconds == 0 {
		omitFields = append(omitFields, "DurationSeconds")
	}
	if len(episode.SubtitleURLs) == 0 || string(episode.SubtitleURLs) == "[]" {
		omitFields = append(omitFields, "SubtitleURLs")
	}

	// 使用事务确保原子性
	return database.DB.Transaction(func(tx *gorm.DB) error {
		txCreate := tx
		if len(omitFields) > 0 {
			txCreate = txCreate.Omit(omitFields...)
		}
		// 1. 插入episode记录
		if err := txCreate.Create(episode).Error; err != nil {
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
	err := database.DB.Where("video_id = ?", videoID).
		Find(&episodes).Error
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

// FindByID 根据ID查找剧集
func (r *episodeRepository) FindByID(episodeID int64) (*model.Episode, error) {
	var episode model.Episode
	err := database.DB.Where("id = ?", episodeID).First(&episode).Error
	if err != nil {
		return nil, err
	}
	return &episode, nil
}

// Update 更新剧集记录
func (r *episodeRepository) Update(episode *model.Episode) error {
	return database.DB.Save(episode).Error
}

// DeleteByID 根据ID删除剧集
func (r *episodeRepository) DeleteByID(episodeID int64) error {
	return database.DB.Delete(&model.Episode{}, episodeID).Error
}

// DeleteByVideoID 根据视频ID删除所有剧集
func (r *episodeRepository) DeleteByVideoID(videoID int64) error {
	return database.DB.Where("video_id = ?", videoID).Delete(&model.Episode{}).Error
}
