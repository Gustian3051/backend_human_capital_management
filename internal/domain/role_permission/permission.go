package role_permission

import (
	"backend/internal/domain/common"

	"github.com/google/uuid"
)

type PermissionModel struct {
    ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    CompanyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uni_company_permission" json:"company_id"`

	Name     string `gorm:"type:varchar(255);not null;uniqueIndex:uni_company_permission" json:"name"`
    Action   string `gorm:"type:varchar(100);not null" json:"action"`
    Resource string `gorm:"type:varchar(100);not null" json:"resource"`
	common.BaseModel
}
