package department

import (
	"backend/internal/domain/common"

	"github.com/google/uuid"
)

type DepartmentModel struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	CompanyID uuid.UUID `gorm:"type:uuid;not null;index" json:"company_id"`

	Name string `gorm:"type:varchar(255);not null;uniqueIndex:uniq_department_company" json:"name"`

	common.BaseModel
}

