package company

import (
	"backend/internal/domain/common"
	"backend/internal/domain/subscription"

	"github.com/google/uuid"
)

type CompanyModel struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	IDCompany     *string `gorm:"type:varchar(255);uniqueIndex" json:"id_company"`
	Name          string  `gorm:"type:varchar(255);not null" json:"name"`
	EmployeeCount int     `json:"employee_count"`
	BussinessType string  `json:"bussiness_type"`
	Picture       string  `json:"picture"`
	CoverPicture  string  `json:"cover_picture"`

	CurrentSubscriptionID      *uuid.UUID                `gorm:"type:uuid" json:"current_subscription_id"`
	CompanySubscriptionHistory []subscription.SubscriptionHistoryModel `gorm:"foreignKey:CompanyID;references:ID" json:"company_subscription_history"`

	EmailAddress       *string `gorm:"type:varchar(100)" json:"email_address"`
	PhoneNumberCompany *string `gorm:"type:varchar(50)" json:"phone_number_company"`
	Website            *string `gorm:"type:varchar(255)" json:"website"`
	NPWP               *string `gorm:"type:varchar(20);uniqueIndex" json:"npwp"`
	FoundedYear        *string `gorm:"type:string" json:"founded_year"`
	Instagram          *string `gorm:"type:varchar(100)" json:"instagram"`
	Facebook           *string `gorm:"type:varchar(100)" json:"facebook"`
	Youtube            *string `gorm:"type:varchar(100)" json:"youtube"`
	RegionCode         *string `gorm:"type:varchar(50);index" json:"region_code"`
	Address            *string `gorm:"type:text" json:"address"`

	OwnerName        *string `gorm:"type:varchar(100)" json:"owner_name"`
	OwnerPosition    *string `gorm:"type:varchar(100)" json:"owner_position"`
	OwnerEmail       *string `gorm:"type:varchar(100)" json:"owner_email"`
	OwnerPhoneNumber *string `gorm:"type:varchar(50)" json:"owner_phone_number"`

	common.BaseModel
}
