package pkg

import (
	"backend/internal/domain/common"

	"github.com/google/uuid"
)

type PackageModel struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(100);unique;not null" json:"name"` // Basic, Premium, Enterprise
	Price        float64   `gorm:"type:numeric" json:"price"`
	Description  *string   `gorm:"type:text" json:"description"`
	DurationDays int       `gorm:"type:int" json:"duration_days"`

	common.BaseModel
}