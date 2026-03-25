package pkg

import (
	"backend/internal/domain/common"
	"time"

	"github.com/google/uuid"
)

type PackageHistoryModel struct {
	ID               uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CompanyID        uuid.UUID       `gorm:"type:uuid;not null;index" json:"company_id"`
	PackageID        uuid.UUID       `gorm:"type:uuid;not null;index" json:"package_id"`
	

	IsActive  bool       `gorm:"default:true" json:"is_active"`
	StartDate time.Time  `gorm:"not null" json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	IsTrial   bool       `gorm:"default:true" json:"is_trial"`

	common.BaseModel
}
