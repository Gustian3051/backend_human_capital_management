package position

import (
	"backend/internal/domain/common"

	"github.com/google/uuid"
)

type PositionModel struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	CompanyID uuid.UUID       `gorm:"type:uuid;not null;index" json:"company_id"`

	DepartmentID uuid.UUID `gorm:"type:uuid;not null;index" json:"department_id"`
	
	Name string `gorm:"type:varchar(255);not null" json:"name"`

	// Hierarchy & approval
	Level        int  `gorm:"not null;default:1;index" json:"level"` // 1=staff, 2=lead, 3=manager, 4=director
	IsManagerial bool `gorm:"not null;default:false" json:"is_managerial"`

	common.BaseModel
}
