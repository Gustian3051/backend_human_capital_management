package role_permission

import "github.com/google/uuid"

type RolePermissionModel struct {
	RoleID       uuid.UUID `gorm:"type:uuid;not null" json:"role_id"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null" json:"permission_id"`
}
