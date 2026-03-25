package shift

import (
	"backend/internal/domain/common"
	"time"

	"github.com/google/uuid"
)

type ShiftModel struct {
	ID               uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CompanyID uuid.UUID       `gorm:"type:uuid;not null;index" json:"company_id"`
	
	ShiftName string `gorm:"type:varchar(255);not null" json:"shift_name"`

	
	DateStart time.Time `gorm:"not null" json:"date_start"`
	DateEnd   time.Time `gorm:"not null" json:"date_end"`

	ShiftStartTime time.Time `gorm:"not null" json:"shift_start_time"`
	ShiftEndTime   time.Time `gorm:"not null" json:"shift_end_time"`

	IsNightShift bool `gorm:"default:false" json:"is_night_shift"`
	IsActive     bool `gorm:"default:true" json:"is_active"`

	common.BaseModel

}

type WorkDayModel struct {
	ID               uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CompanyID uuid.UUID       `gorm:"type:uuid;not null;index" json:"company_id"`
	
	Name string `gorm:"type:varchar(100);not null" json:"name"`

	common.BaseModel

}

type ShiftWorkDayModel struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ShiftID   uuid.UUID `gorm:"type:uuid;not null;index"`
	WorkDayID uuid.UUID `gorm:"type:uuid;not null;index"`
}
