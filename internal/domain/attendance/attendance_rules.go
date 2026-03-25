package attendance

import (
	"backend/internal/domain/common"

	"github.com/google/uuid"
)

type AttendanceRulesModel struct {
	ID               uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CompanyID uuid.UUID       `gorm:"type:uuid;not null;uniqueIndex" json:"company_id"`

	OfficeLatitude  *float64 `json:"office_latitude,omitempty"`
	OfficeLongitude *float64 `json:"office_longitude,omitempty"`
	RadiusMeters    int      `gorm:"default:100" json:"radius_meters"`

	MaxLateMinutes  int `gorm:"default:30" json:"max_late_minutes"`
	MaxLateCheckIn  int `gorm:"default:30" json:"max_late_check_in"`
	MaxLateCheckOut int `gorm:"default:30" json:"max_late_check_out"`

	common.BaseModel
}
