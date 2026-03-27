package repository

import (
	"backend/internal/domain/employee"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	GetEmployeeByID(ctx context.Context, id uuid.UUID) (*employee.EmployeeModel, error)
	FindByUserAndCompanyTx(tx *gorm.DB, userID, companyID uuid.UUID) (*employee.EmployeeModel, error)
	GetEmployeeByEmail(ctx context.Context, email string) (*employee.EmployeeModel, error)
	CreateTx(tx *gorm.DB, employee *employee.EmployeeModel) error
	UpdateEmployee(ctx context.Context, employee *employee.EmployeeModel) error
	DeleteEmployee(ctx context.Context, id uuid.UUID) error
	WithContext(ctx context.Context) EmployeeRepository
	Transaction(f func(tx *gorm.DB) error) error
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) WithContext(ctx context.Context) EmployeeRepository {
	return &employeeRepository{db: r.db.WithContext(ctx)}
}

func (r *employeeRepository) Transaction(f func(tx *gorm.DB) error) error {
	return r.db.Transaction(f)
}


func (r *employeeRepository) GetEmployeeByID(ctx context.Context, id uuid.UUID) (*employee.EmployeeModel, error) {
	var e employee.EmployeeModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&e).Error
	return &e, err
}

func (r *employeeRepository) GetEmployeeByEmail(ctx context.Context, email string) (*employee.EmployeeModel, error) {
	var e employee.EmployeeModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&e).Error
	return &e, err
}

func (r *employeeRepository) CreateTx(tx *gorm.DB, employee *employee.EmployeeModel) error {
	return tx.Create(employee).Error
}

func (r *employeeRepository) FindByUserAndCompanyTx(tx *gorm.DB, userID, companyID uuid.UUID) (*employee.EmployeeModel, error) {
	var e employee.EmployeeModel
	err := tx.
		Where("user_id = ? AND company_id = ?", userID, companyID).
		First(&e).Error
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *employeeRepository) UpdateEmployee(ctx context.Context, employee *employee.EmployeeModel) error {
	return r.db.WithContext(ctx).Save(employee).Error
}

func (r *employeeRepository) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&employee.EmployeeModel{}, "id = ?", id).Error
}