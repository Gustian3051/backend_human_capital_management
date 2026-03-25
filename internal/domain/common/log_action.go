package common

import (
	"time"

	"github.com/google/uuid"
)

type LogActionModel struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	UserID           *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	CompanyID *uuid.UUID `gorm:"type:uuid;index" json:"company_id"`

	Action    string    `gorm:"type:varchar(255);not null" json:"action"`
	IPAddress string    `gorm:"type:varchar(50)" json:"ip_address"`
	DateTime  time.Time `gorm:"autoCreateTime" json:"date_time"`
}
