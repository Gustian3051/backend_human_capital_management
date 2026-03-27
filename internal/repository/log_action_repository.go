package repository

import (
	"context"

	"backend/internal/domain/common"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LogRepository interface {
	Create(ctx context.Context, log *common.LogActionModel) error
	GetEmployeeNameByUserID(ctx context.Context, userID uuid.UUID) (string, error)
}

type logRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

func (r *logRepository) Create(ctx context.Context, log *common.LogActionModel) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *logRepository) GetEmployeeNameByUserID(ctx context.Context, userID uuid.UUID) (string, error) {
	type result struct {
		FullName string
	}

	var res result

	err := r.db.WithContext(ctx).
		Table("employee_models").
		Select("full_name").
		Where("user_id = ?", userID).
		First(&res).Error

	if err != nil {
		return "", err
	}

	return res.FullName, nil
}