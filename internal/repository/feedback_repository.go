package repository

import (
	"video-service/internal/model"
	"video-service/pkg/infrastructure/database"
)

// FeedbackListFilter 管理端反馈列表筛选
type FeedbackListFilter struct {
	Source      string
	Status      string
	ProblemType string
	VideoID     *int64
	UserID      *int64
	Keyword     string
	Page        int
	PageSize    int
}

// FeedbackRepository 用户反馈数据访问
type FeedbackRepository interface {
	FindByID(id int64) (*model.VideoFeedback, error)
	List(filter FeedbackListFilter) ([]model.VideoFeedback, int64, error)
	Update(feedback *model.VideoFeedback) error
}

type feedbackRepository struct{}

// NewFeedbackRepository 创建反馈仓库
func NewFeedbackRepository() FeedbackRepository {
	return &feedbackRepository{}
}

func (r *feedbackRepository) FindByID(id int64) (*model.VideoFeedback, error) {
	var feedback model.VideoFeedback
	if err := database.DB.Where("id = ?", id).First(&feedback).Error; err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *feedbackRepository) List(filter FeedbackListFilter) ([]model.VideoFeedback, int64, error) {
	page := filter.Page
	pageSize := filter.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query := database.DB.Model(&model.VideoFeedback{})
	if filter.Source != "" {
		query = query.Where("source = ?", filter.Source)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ProblemType != "" {
		query = query.Where("problem_type = ?", filter.ProblemType)
	}
	if filter.VideoID != nil {
		query = query.Where("video_id = ?", *filter.VideoID)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("video_title LIKE ? OR content LIKE ? OR search_keyword LIKE ?", like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []model.VideoFeedback
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *feedbackRepository) Update(feedback *model.VideoFeedback) error {
	return database.DB.Save(feedback).Error
}
