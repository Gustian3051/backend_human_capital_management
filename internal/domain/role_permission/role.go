package role_permission

import (
	"backend/internal/domain/common"
	"backend/internal/domain/employee"

	"github.com/google/uuid"
)

type RoleModel struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CompanyID uuid.UUID `gorm:"type:uuid;not null;index" json:"company_id"`

	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	common.BaseModel
	employee.EmployeeModel
}
