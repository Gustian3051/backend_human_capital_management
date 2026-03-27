package service

import (
	"backend/internal/domain/common"
	"backend/internal/domain/role_permission"
	"backend/internal/domain/user"
	"backend/internal/dto"
	"backend/internal/infrastructure/database"
	"backend/internal/repository"
	"backend/pkg/log"
	"fmt"
	"time"

	"context"
	"errors"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateRoleWithPermission(userID, companyProfileID uuid.UUID, req dto.RolePermissionRequest, ipAddress string) (*dto.RolePermissionPayload, error)
// UpdateRoleWithPermission(userID, companyProfileID, roleID uuid.UUID, req dto.RolePermissionRequest, ipAddress string) (*dto.RolePermissionPayload, error)
// DeleteRoleWithPermission(userID, companyProfileID, roleID uuid.UUID, ipAddress string) error

// DetailRoleAndPermission(ctx context.Context, companyProfileID, roleID uuid.UUID) (*dto.RolePermissionPayload, error)
// DataRoleAndPermission(ctx context.Context, companyProfileID uuid.UUID) (*dto.DataRoleAndPermission, error)
// DataPermission(ctx context.Context, companyProfileID uuid.UUID) ([]dto.PermissionGroup, error)

type RolePermissionServiceInterface interface {
	GetPermissionsForUser(ctx context.Context, user *user.UserModel) ([]dto.PermissionInfo, error)
	GetOrCreateAdminRole(ctx *gorm.DB, companyID uuid.UUID) (*role_permission.RoleModel, error)
}

type RoleAndPermissionService struct {
	Repo     repository.PermissionRepository
	Seeder   database.RolePermissionSeederInterface
	Enforcer *casbin.Enforcer
}

func NewRoleAndPermissionService(repo repository.PermissionRepository, seeder database.RolePermissionSeederInterface, enforcer *casbin.Enforcer) RolePermissionServiceInterface {
	return &RoleAndPermissionService{
		Repo:     repo,
		Seeder:   seeder,
		Enforcer: enforcer,
	}
}

func (s *RoleAndPermissionService) GetPermissionsForUser(ctx context.Context, user *user.UserModel) ([]dto.PermissionInfo, error) {

	if user == nil || len(user.Employees) == 0 {
		return []dto.PermissionInfo{}, nil
	}

	emp := user.Employees[0]
	if emp.RoleID == uuid.Nil {
		return []dto.PermissionInfo{}, nil
	}

	role, err := s.Repo.GetRoleWithPermissions(ctx, emp.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Fatal("Role not found",
				zap.Error(errors.New("role not found")),
			)
		}
		logger.Log.Fatal("Failed to get role",
			zap.Error(err),
		)
	}

	perms := make([]dto.PermissionInfo, 0, len(role.Permissions))

	for _, p := range role.Permissions {
		perms = append(perms, dto.PermissionInfo{
			Name: p.Name,
		})
	}

	return perms, nil
}

func (s *RoleAndPermissionService) GetOrCreateAdminRole(tx *gorm.DB, companyID uuid.UUID) (*role_permission.RoleModel, error) {

	var role role_permission.RoleModel

	// ===== 1. Find existing =====
	err := tx.
		Where("name = ? AND company_id = ?", "Admin", companyID).
		Preload("Permissions").
		First(&role).Error

	if err != nil {

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to get role: %w", err)
		}

		// ===== 2. Create role =====
		now := time.Now()

		role = role_permission.RoleModel{
			ID:          uuid.New(),
			CompanyID:   companyID,
			Name:        "Admin",
			Description: "Administrator role with full permissions",
			BaseModel: common.BaseModel{
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		if err := tx.Create(&role).Error; err != nil {
			return nil, fmt.Errorf("failed to create role: %w", err)
		}
	}

	// ===== 3. VALIDASI (WAJIB) =====
	if role.ID == uuid.Nil {
		return nil, fmt.Errorf("invalid role created")
	}

	// ===== 4. Seed permission =====
	allPermission, err := s.Seeder.SeedAdminPermissionTx(tx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to seed permission: %w", err)
	}

	// ===== 5. Assign permission (SAFE) =====
	var count int64

	err = tx.
		Table("role_permission_models").
		Where("role_id = ? AND permission_id = ?", role.ID, allPermission.ID).
		Count(&count).Error

	if err != nil {
		return nil, err
	}

	if count == 0 {
		if err := tx.Model(&role).
			Association("Permissions").
			Append(allPermission); err != nil {
			return nil, fmt.Errorf("failed to assign permission: %w", err)
		}
	}

	return &role, nil
}