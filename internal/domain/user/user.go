package user

import (
	"backend/internal/domain/common"
	"backend/internal/domain/employee"

	"github.com/google/uuid"
)

type UserModel struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	CompanyID *uuid.UUID `gorm:"type:uuid;index" json:"company_id"`
	
	Email           string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password        string `gorm:"type:varchar(255);not null" json:"password"`
	PasswordDefault string `gorm:"type:varchar(255)" json:"password_default"`
	IsVerified      bool   `json:"is_verified"`
	NeedsProfile    bool   `json:"needs_profile"`
	common.BaseModel

	Employees []employee.EmployeeModel `gorm:"foreignKey:UserID;references:ID" json:"employees"`
}
