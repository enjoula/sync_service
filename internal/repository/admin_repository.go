package repository

import (
	"strings"
	"time"

	"video-service/internal/model"
	"video-service/pkg/infrastructure/database"

	"gorm.io/gorm"
)

// AdminRepository 管理员仓库
type AdminRepository interface {
	FindByUsername(username string) (*model.Admin, error)
	FindByID(id int64) (*model.Admin, error)
	UpdateLastLogin(id int64, at time.Time) error
	List(keyword string, limit, offset int) ([]model.Admin, int64, error)
	Create(admin *model.Admin) error
	Update(id int64, fields map[string]interface{}) error
	CountEnabledByRole(role string, excludeID int64) (int64, error)
}

type adminRepository struct{}

// NewAdminRepository 创建管理员仓库
func NewAdminRepository() AdminRepository {
	return &adminRepository{}
}

func (r *adminRepository) FindByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	err := database.DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) FindByID(id int64) (*model.Admin, error) {
	var admin model.Admin
	err := database.DB.First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) UpdateLastLogin(id int64, at time.Time) error {
	return database.DB.Model(&model.Admin{}).Where("id = ?", id).Update("last_login_at", at).Error
}

func (r *adminRepository) List(keyword string, limit, offset int) ([]model.Admin, int64, error) {
	query := database.DB.Model(&model.Admin{}).
		Select("id", "username", "nickname", "role", "status", "last_login_at", "created_at", "updated_at")

	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR nickname LIKE ?", like, like)
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var admins []model.Admin
	err := query.Session(&gorm.Session{}).Order("id DESC").Limit(limit).Offset(offset).Find(&admins).Error
	if err != nil {
		return nil, 0, err
	}
	return admins, total, nil
}

func (r *adminRepository) Create(admin *model.Admin) error {
	return database.DB.Create(admin).Error
}

func (r *adminRepository) Update(id int64, fields map[string]interface{}) error {
	return database.DB.Model(&model.Admin{}).Where("id = ?", id).Updates(fields).Error
}

func (r *adminRepository) CountEnabledByRole(role string, excludeID int64) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Admin{}).
		Where("role = ? AND status = 1 AND id <> ?", role, excludeID).
		Count(&count).Error
	return count, err
}
