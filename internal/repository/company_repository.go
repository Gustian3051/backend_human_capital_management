package repository

import (
	"context"
	"errors"

	"backend/internal/domain/company"
	"backend/internal/domain/employee"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyRepository interface {
	GetCompanyWithPackage(ctx context.Context, companyID uuid.UUID) (*company.CompanyModel, error)
	FindByName(tx *gorm.DB, name string) (*company.CompanyModel, error)
	Create(tx *gorm.DB, name string) (*company.CompanyModel, error)
	GetOrCreate(tx *gorm.DB, name string) (*company.CompanyModel, error)
	CountEmployees(tx *gorm.DB, companyID uuid.UUID) (int64, error)

	WithContext(ctx context.Context) CompanyRepository
	Transaction(f func(tx *gorm.DB) error) error
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) WithContext(ctx context.Context) CompanyRepository {
	return &companyRepository{db: r.db.WithContext(ctx)}
}

func (r *companyRepository) Transaction(f func(tx *gorm.DB) error) error {
	return r.db.Transaction(f)
}

func (r *companyRepository) GetCompanyWithPackage(ctx context.Context, companyID uuid.UUID) (*company.CompanyModel, error) {
	var c company.CompanyModel

	err := r.db.WithContext(ctx).
		Preload("CompanySubscriptionHistory", "is_active = ?", true).
		Preload("CompanySubscriptionHistory.Subscription").
		Where("id = ?", companyID).
		First(&c).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &c, err
}

func (r *companyRepository) FindByName(tx *gorm.DB, name string) (*company.CompanyModel, error) {
	var c company.CompanyModel

	err := tx.Where("name = ?", name).First(&c).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &c, err
}

func (r *companyRepository) Create(tx *gorm.DB, name string) (*company.CompanyModel, error) {
	company := company.CompanyModel{
		ID:   uuid.New(),
		Name: name,
	}

	if err := tx.Create(&company).Error; err != nil {
		return nil, err
	}

	return &company, nil
}

func (r *companyRepository) GetOrCreate(tx *gorm.DB, name string) (*company.CompanyModel, error) {

	// 🔍 cek dulu
	existing, err := r.FindByName(tx, name)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return existing, nil
	}

	// 🏗 create baru
	return r.Create(tx, name)
}

func (r *companyRepository) CountEmployees(tx *gorm.DB, companyID uuid.UUID) (int64, error) {
	var count int64
	err := tx.Model(&employee.EmployeeModel{}).
		Where("company_id = ?", companyID).
		Count(&count).Error
	return count, err
}