package repository

import (
	"context"

	"backend/internal/domain/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.UserModel, error)
	FindByID(ctx context.Context, id uuid.UUID) (*user.UserModel, error)
	Create(ctx context.Context, u *user.UserModel) error
	Update(ctx context.Context, u *user.UserModel) error
	UpdateVerified(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.UserModel, error) {
	var u user.UserModel
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		Preload("Employees.Role").
		Preload("Employees.Shift").
		First(&u).Error

	return &u, err
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.UserModel, error) {
	var u user.UserModel
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("Employees.Role").
		Preload("Employees.Shift").
		First(&u).Error

	return &u, err
}

func (r *userRepository) Create(ctx context.Context, u *user.UserModel) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepository) Update(ctx context.Context, u *user.UserModel) error {
	return r.db.WithContext(ctx).
		Model(&user.UserModel{}).
		Where("id = ?", u.ID).
		Updates(u).Error
}

func (r *userRepository) UpdateVerified(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&user.UserModel{}).
		Where("id = ?", id).
		Update("is_verified", true).Error
}