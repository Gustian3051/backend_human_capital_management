package work_leave

import (
	"backend/internal/domain/common"

	"github.com/google/uuid"
)

type LeaveBalanceModel struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	CompanyID uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	EmployeeID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uniq_employee_year" json:"employee_id"`

	Year int `gorm:"not null;uniqueIndex:uniq_employee_year" json:"year"`

	MaxQuantity int `gorm:"not null;" json:"max_quantity"`
	Used        int `gorm:"not null;default:0" json:"used"`
	
	common.BaseModel
}
