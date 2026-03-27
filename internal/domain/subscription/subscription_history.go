package subscription

import (
	"backend/internal/domain/common"
	"time"

	"github.com/google/uuid"
)

type SubscriptionHistoryModel struct {
	ID               uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CompanyID        uuid.UUID       `gorm:"type:uuid;not null;index" json:"company_id"`
	
	SubscriptionID        uuid.UUID       `gorm:"type:uuid;not null;index" json:"subscription_id"`
	Subscription          *SubscriptionModel   `gorm:"foreignKey:SubscriptionID;references:ID" json:"subscription"`

	IsActive  bool       `gorm:"default:true" json:"is_active"`
	StartDate time.Time  `gorm:"not null" json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	IsTrial   bool       `gorm:"default:true" json:"is_trial"`

	common.BaseModel
}
