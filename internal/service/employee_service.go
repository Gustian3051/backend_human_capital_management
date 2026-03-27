package service

import (
	"backend/internal/domain/common"
	"backend/internal/domain/company"
	"backend/internal/domain/employee"
	"backend/internal/domain/user"
	"backend/internal/repository"
	"backend/pkg/helpers"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeServiceInterface interface {
	CreateDefaultEmployee(ctx context.Context, user *user.UserModel, company *company.CompanyModel, fullName string) (*employee.EmployeeModel, error)
}

type EmployeeService struct {
	EmployeeRepo             repository.EmployeeRepository
	RoleAndPermissionService RolePermissionServiceInterface
}

func NewEmployeeService(
	employeeRepo repository.EmployeeRepository,
	roleAndPermissionService RolePermissionServiceInterface,
) EmployeeServiceInterface {
	return &EmployeeService{
		EmployeeRepo:             employeeRepo,
		RoleAndPermissionService: roleAndPermissionService,
	}
}

func (s *EmployeeService) CreateDefaultEmployee(ctx context.Context, user *user.UserModel, company *company.CompanyModel, fullName string) (*employee.EmployeeModel, error) {

	var result *employee.EmployeeModel

	err := s.EmployeeRepo.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// ===== 1. cek existing =====
		existing, err := s.EmployeeRepo.FindByUserAndCompanyTx(tx, user.ID, company.ID)
		if err == nil {
			result = existing
			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		role, err := s.RoleAndPermissionService.GetOrCreateAdminRole(tx, company.ID)
		if err != nil {
			return err
		}
		name := helpers.NormalizeName(fullName, user.Email)

		now := time.Now()

		// ===== 4. create employee =====
		emp := &employee.EmployeeModel{
			ID:        uuid.New(),
			UserID:    user.ID,
			CompanyID: company.ID,
			RoleID:    role.ID,
			FullName:  name,
			Status:    "active",
			BaseModel: common.BaseModel{
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		if err := s.EmployeeRepo.CreateTx(tx, emp); err != nil {
			return err
		}

		result = emp
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}