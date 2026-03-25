package shift

import (
	"backend/internal/domain/common"

	"github.com/google/uuid"
)

type ShiftAssignmentModel struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	EmployeeID       uuid.UUID `gorm:"type:uuid;index" json:"employee_id"`
	ShiftID          uuid.UUID `gorm:"type:uuid;index" json:"shift_id"`
	CompanyID uuid.UUID `gorm:"type:uuid;index" json:"company_id"`

	IsActive bool `gorm:"default:false" json:"is_active"`

	common.BaseModel
}
