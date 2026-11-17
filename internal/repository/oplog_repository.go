package repository

import (
	"video-service/internal/model"
	"video-service/pkg/infrastructure/database"
)

// OpLogListFilter 操作日志列表筛选
type OpLogListFilter struct {
	AdminID    *int64
	Action     string
	ObjectType string
	Result     string
	Keyword    string
	Page       int
	PageSize   int
}

// OpLogRepository 操作日志数据访问
type OpLogRepository interface {
	List(filter OpLogListFilter) ([]model.AdminOperationLog, int64, error)
	FindByID(id int64) (*model.AdminOperationLog, error)
}

type opLogRepository struct{}

// NewOpLogRepository 创建操作日志仓库
func NewOpLogRepository() OpLogRepository {
	return &opLogRepository{}
}

func (r *opLogRepository) List(filter OpLogListFilter) ([]model.AdminOperationLog, int64, error) {
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

	query := database.DB.Model(&model.AdminOperationLog{})
	if filter.AdminID != nil {
		query = query.Where("admin_id = ?", *filter.AdminID)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.ObjectType != "" {
		query = query.Where("object_type = ?", filter.ObjectType)
	}
	if filter.Result != "" {
		query = query.Where("result = ?", filter.Result)
	}
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where(
			"admin_username LIKE ? OR object_name LIKE ? OR object_id LIKE ? OR error_msg LIKE ?",
			like, like, like, like,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []model.AdminOperationLog
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC, id DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *opLogRepository) FindByID(id int64) (*model.AdminOperationLog, error) {
	var row model.AdminOperationLog
	if err := database.DB.Where("id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}
