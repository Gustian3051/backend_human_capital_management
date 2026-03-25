package attendance

import (
	"backend/internal/domain/common"
	"time"

	"github.com/google/uuid"
)

type AttendanceModel struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	EmployeeID uuid.UUID `gorm:"type:uuid;not null;index" json:"employee_id"`
	

	ShiftID uuid.UUID `gorm:"type:uuid;not null;index" json:"shift_id"`

	CompanyID uuid.UUID       `gorm:"type:uuid;not null;index" json:"company_id"`

	AttendanceRuleID uuid.UUID       `gorm:"type:uuid;index" json:"attendance_rule_id"`
	Date             time.Time       `gorm:"type:date;not null;index" json:"date"`

	// CheckInTime
	CheckInTime   *time.Time `json:"check_in_time,omitempty"`
	StatusCheckIn string     `gorm:"type:varchar(100)" json:"status_check_in,omitempty"`
	CheckInNote   *string    `gorm:"type:text" json:"check_in_note,omitempty"`

	// CheckOutTime
	CheckOutTime   *time.Time `json:"check_out_time,omitempty"`
	StatusCheckOut string     `gorm:"type:varchar(100)" json:"status_check_out,omitempty"`
	CheckOutNote   *string    `gorm:"type:text" json:"check_out_note,omitempty"`

	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`

	common.BaseModel
}
