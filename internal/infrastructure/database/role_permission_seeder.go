package database

import (
	"backend/internal/domain/common"
	"backend/internal/domain/role_permission"
	"backend/internal/domain/shift"
	"backend/internal/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RolePermissionSeederInterface interface {
	SeedDefaultRolesTx(tx *gorm.DB, companyID uuid.UUID) (*role_permission.RoleModel, error)
	AssignPermissionsToAdminTx(tx *gorm.DB, roleID uuid.UUID) error
	SeedAdminPermissionTx(tx *gorm.DB, companyID uuid.UUID) (*role_permission.PermissionModel, error)
	SeedDefaultPermissionsTx(tx *gorm.DB) error
}

type RolePermissionSeeder struct {
	PermissionRepo repository.PermissionRepository
}

func NewRolePermissionSeeder(permissionRepo repository.PermissionRepository) RolePermissionSeederInterface {
	return &RolePermissionSeeder{PermissionRepo: permissionRepo}
}

func (s *RolePermissionSeeder) SeedDefaultRolesTx(tx *gorm.DB, companyID uuid.UUID) (*role_permission.RoleModel, error) {

	adminRole := &role_permission.RoleModel{
		ID:        uuid.New(),
		Name:      "admin",		
		CompanyID: companyID,
	}

	if err := tx.
		Where("name = ? AND company_id = ?", "admin", companyID).
		FirstOrCreate(adminRole).Error; err != nil {
		return nil, err
	}

	return adminRole, nil
}

func (s *RolePermissionSeeder) AssignPermissionsToAdminTx(tx *gorm.DB, roleID uuid.UUID) error {

	var permissions []role_permission.PermissionModel
	if err := tx.Find(&permissions).Error; err != nil {
		return fmt.Errorf("fetch permissions: %w", err)
	}

	for _, perm := range permissions {
		rp := role_permission.RolePermissionModel{
			RoleID:       roleID,
			PermissionID: perm.ID,
		}

		if err := tx.
			Where("role_id = ? AND permission_id = ?", roleID, perm.ID).
			FirstOrCreate(&rp).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *RolePermissionSeeder) SeedAdminPermissionTx(tx *gorm.DB, companyID uuid.UUID) (*role_permission.PermissionModel, error) {

	perm := &role_permission.PermissionModel{
		ID:        uuid.New(),
		Name:      "all:all",
		Action:    "*",
		Resource:  "*",
	}

	if err := s.PermissionRepo.UpsertPermission(tx, perm); err != nil {
		return nil, fmt.Errorf("seed admin permission: %w", err)
	}

	existing, err := s.PermissionRepo.FindByNamePermission(tx, "all:all")
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *RolePermissionSeeder) SeedDefaultPermissionsTx(tx *gorm.DB) error {
	perms := []role_permission.PermissionModel{
		{ID: uuid.New(), Name: "dashboard:read", Action: "read", Resource: "dashboard"},
		// Employee Management
		{ID: uuid.New(), Name: "employee-management:create", Action: "create", Resource: "employee-management"},
		{ID: uuid.New(), Name: "employee-management:read", Action: "read", Resource: "employee-management"},
		{ID: uuid.New(), Name: "employee-management:update", Action: "update", Resource: "employee-management"},
		{ID: uuid.New(), Name: "employee-management:delete", Action: "delete", Resource: "employee-management"},
		// Department Management
		{ID: uuid.New(), Name: "department-management:create", Action: "create", Resource: "department-management"},
		{ID: uuid.New(), Name: "department-management:read", Action: "read", Resource: "department-management"},
		{ID: uuid.New(), Name: "department-management:update", Action: "update", Resource: "department-management"},
		{ID: uuid.New(), Name: "department-management:delete", Action: "delete", Resource: "department-management"},
		// Position Management
		{ID: uuid.New(), Name: "position-management:create", Action: "create", Resource: "position-management"},
		{ID: uuid.New(), Name: "position-management:read", Action: "read", Resource: "position-management"},
		{ID: uuid.New(), Name: "position-management:update", Action: "update", Resource: "position-management"},
		{ID: uuid.New(), Name: "position-management:delete", Action: "delete", Resource: "position-management"},
		// Role Permission Management
		{ID: uuid.New(), Name: "role-permission-management:create", Action: "create", Resource: "role-permission-management"},
		{ID: uuid.New(), Name: "role-permission-management:read", Action: "read", Resource: "role-permission-management"},
		{ID: uuid.New(), Name: "role-permission-management:update", Action: "update", Resource: "role-permission-management"},
		{ID: uuid.New(), Name: "role-permission-management:delete", Action: "delete", Resource: "role-permission-management"},
		// Company Management
		{ID: uuid.New(), Name: "company-management:read", Action: "read", Resource: "company-management"},
		{ID: uuid.New(), Name: "company-management:create", Action: "create", Resource: "company-management"},
		{ID: uuid.New(), Name: "company-management:update", Action: "update", Resource: "company-management"},
		{ID: uuid.New(), Name: "company-management:delete", Action: "delete", Resource: "company-management"},
		// Shift Management
		{ID: uuid.New(), Name: "shift-management:read", Action: "read", Resource: "shift-management"},
		{ID: uuid.New(), Name: "shift-management:create", Action: "create", Resource: "shift-management"},
		{ID: uuid.New(), Name: "shift-management:update", Action: "update", Resource: "shift-management"},
		{ID: uuid.New(), Name: "shift-management:delete", Action: "delete", Resource: "shift-management"},
		// Attendance
		{ID: uuid.New(), Name: "attendance:read", Action: "read", Resource: "attendance"},
		{ID: uuid.New(), Name: "attendance:create", Action: "create", Resource: "attendance"},
		{ID: uuid.New(), Name: "attendance:delete", Action: "delete", Resource: "attendance"},
		{ID: uuid.New(), Name: "attendance:update", Action: "update", Resource: "attendance"},
		// Attendance Management
		{ID: uuid.New(), Name: "attendance-management:read", Action: "read", Resource: "attendance-management"},
		{ID: uuid.New(), Name: "attendance-management:create", Action: "create", Resource: "attendance-management"},
		{ID: uuid.New(), Name: "attendance-management:delete", Action: "delete", Resource: "attendance-management"},
		{ID: uuid.New(), Name: "attendance-management:update", Action: "update", Resource: "attendance-management"},
		// Work Leave
		{ID: uuid.New(), Name: "work-leave:read", Action: "read", Resource: "work-leave"},
		{ID: uuid.New(), Name: "work-leave:create", Action: "create", Resource: "work-leave"},
		{ID: uuid.New(), Name: "work-leave:delete", Action: "delete", Resource: "work-leave"},
		{ID: uuid.New(), Name: "work-leave:update", Action: "update", Resource: "work-leave"},
		// Work Leave Management
		{ID: uuid.New(), Name: "work-leave-management:read", Action: "read", Resource: "work-leave-management"},
		{ID: uuid.New(), Name: "work-leave-management:create", Action: "create", Resource: "work-leave-management"},
		{ID: uuid.New(), Name: "work-leave-management:delete", Action: "delete", Resource: "work-leave-management"},
		{ID: uuid.New(), Name: "work-leave-management:update", Action: "update", Resource: "work-leave-management"},
		// Work Plan
		{ID: uuid.New(), Name: "work-plan:read", Action: "read", Resource: "work-plan"},
		{ID: uuid.New(), Name: "work-plan:create", Action: "create", Resource: "work-plan"},
		{ID: uuid.New(), Name: "work-plan:delete", Action: "delete", Resource: "work-plan"},
		{ID: uuid.New(), Name: "work-plan:update", Action: "update", Resource: "work-plan"},
		// Job Vacancy Management
		{ID: uuid.New(), Name: "job-vacancy-management:read", Action: "read", Resource: "job-vacancy-management"},
		{ID: uuid.New(), Name: "job-vacancy-management:create", Action: "create", Resource: "job-vacancy-management"},
		{ID: uuid.New(), Name: "job-vacancy-management:delete", Action: "delete", Resource: "job-vacancy-management"},
		{ID: uuid.New(), Name: "job-vacancy-management:update", Action: "update", Resource: "job-vacancy-management"},
		// Notification
		// {ID: uuid.New(), Name: "notification:read", Action: "read", Resource: "notification", CompanyProfileID: companyProfileID},
		// {ID: uuid.New(), Name: "notification:create", Action: "create", Resource: "notification", CompanyProfileID: companyProfileID},
		// {ID: uuid.New(), Name: "notification:delete", Action: "delete", Resource: "notification", CompanyProfileID: companyProfileID},
		// {ID: uuid.New(), Name: "notification:update", Action: "update", Resource: "notification", CompanyProfileID: companyProfileID},
	}

	return s.PermissionRepo.BatchUpsertPermission(tx, perms)
}

// SeedDefaultPermissions is a convenience function for seeding from main.go.
// It creates the repo internally and runs the seed in a transaction.
func SeedDefaultPermissions(db *gorm.DB) error {
	permRepo := repository.NewPermissionRepository(db)
	seeder := &RolePermissionSeeder{PermissionRepo: permRepo}
	return db.Transaction(func(tx *gorm.DB) error {
		return seeder.SeedDefaultPermissionsTx(tx)
	})
}

func SeedWorkDays(db *gorm.DB, companyID uuid.UUID) error {

	// Cek apakah data sudah ada
	var count int64
	if err := db.Model(&shift.WorkDayModel{}).
		Where("company_id = ?", companyID).
		Count(&count).Error; err != nil {

		return err
	}

	// Kalau sudah ada, skip
	if count > 0 {
		return nil
	}

	// List hari
	days := []string{
		"Senin",
		"Selasa",
		"Rabu",
		"Kamis",
		"Jumat",
		"Sabtu",
		"Minggu",
	}

	now := time.Now()

	var seeds []shift.WorkDayModel
	for _, day := range days {
		seeds = append(seeds, shift.WorkDayModel{
			ID:        uuid.New(),
			CompanyID: companyID,
			Name:      day,
			BaseModel: common.BaseModel{
				CreatedAt: now,
				UpdatedAt: now,
			},
		})
	}

	// Insert batch
	if err := db.Create(&seeds).Error; err != nil {
		return err
	}

	return nil
}