package repository

import (
	"context"

	"backend/internal/domain/role_permission"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PermissionRepository interface {
	UpsertPermission(tx *gorm.DB, perm *role_permission.PermissionModel) error
	GetRoleWithPermissions(ctx context.Context, roleID uuid.UUID) (*role_permission.RoleModel, error)
	FindByNamePermission(tx *gorm.DB, name string) (*role_permission.PermissionModel, error)
	BatchUpsertPermission(tx *gorm.DB, perms []role_permission.PermissionModel) error	
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) UpsertPermission(tx *gorm.DB, perm *role_permission.PermissionModel) error {
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "name"},
		},
		DoNothing: true,
	}).Create(perm).Error
}

func (r *permissionRepository) GetRoleWithPermissions(ctx context.Context, roleID uuid.UUID) (*role_permission.RoleModel, error) {
	var role role_permission.RoleModel

	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("id = ?", roleID).
		First(&role).Error

	return &role, err
}

func (r *permissionRepository) FindByNamePermission(tx *gorm.DB, name string) (*role_permission.PermissionModel, error) {
	var p role_permission.PermissionModel
	err := tx.Where("name = ?", name).
		First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *permissionRepository) BatchUpsertPermission(tx *gorm.DB, perms []role_permission.PermissionModel) error {
	for _, p := range perms {
		if err := r.UpsertPermission(tx, &p); err != nil {
			return err
		}
	}
	return nil
}