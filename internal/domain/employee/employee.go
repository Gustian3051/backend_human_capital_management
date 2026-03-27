package employee

import (
	"backend/internal/domain/common"
	"backend/internal/domain/company"
	"backend/internal/domain/department"
	"backend/internal/domain/position"
	"backend/internal/domain/role_permission"
	"backend/internal/domain/shift"

	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type EmployeeModel struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	UserID       uuid.UUID  `gorm:"type:uuid;index" json:"user_id,omitempty"`
	
	CompanyID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"company_id"`
	Company      *company.CompanyModel `gorm:"foreignKey:CompanyID;references:ID" json:"company"`
	
	DepartmentID *uuid.UUID `gorm:"type:uuid" json:"department_id,omitempty"`
	Department   *department.DepartmentModel `gorm:"foreignKey:DepartmentID;references:ID" json:"department"`
	
	PositionID   *uuid.UUID `gorm:"type:uuid" json:"position_id,omitempty"`
	Position     *position.PositionModel `gorm:"foreignKey:PositionID;references:ID" json:"position"`
	
	RoleID       uuid.UUID  `gorm:"type:uuid;index" json:"role_id,omitempty"`
	Role         *role_permission.RoleModel `gorm:"foreignKey:RoleID;references:ID" json:"role"`
	
	ShiftID      *uuid.UUID `gorm:"type:uuid;index" json:"shift_id,omitempty"`
	Shift        *shift.ShiftModel `gorm:"foreignKey:ShiftID;references:ID" json:"shift"`

	// personal information
	Picture         *string    `gorm:"size:255" json:"picture"`
	FullName        string     `gorm:"size:150" json:"full_name"`
	BirthPlace      *string    `gorm:"size:255" json:"birth_place"`
	BirthDate       *time.Time `json:"birth_date"`
	Gender          *string    `gorm:"size:50" json:"gender"`
	KTPNumber       *string    `gorm:"uniqueIndex:uni_ktp_company;size:50" json:"ktp_number"`
	NPWPNumber      *string    `gorm:"uniqueIndex:uni_npwp_company;size:50" json:"npwp_number"`
	MaritalStatus   *string    `gorm:"size:50" json:"marital_status"`
	Citizenship     *string    `gorm:"size:50" json:"citizenship"`
	Religion        *string    `gorm:"size:50" json:"religion"`
	Address         *string    `gorm:"size:255" json:"address"`
	DomicileAddress *string    `gorm:"size:255" json:"domicile_address"`
	BloodType       *string    `gorm:"size:50" json:"blood_type"`

	// work data
	IDCardNumber     *string     `gorm:"uniqueIndex:uni_idcard_company;size:50" json:"id_card_number"`
	JoinedDate       *time.Time  `json:"joined_date"`
	EmploymentStatus string      `gorm:"size:50" json:"employment_status"`
	ResignDate       *time.Time  `json:"resign_date"`
	Status           string      `gorm:"size:50" json:"status"`
	
	// information data
	PhoneNumber    *string `gorm:"size:50" json:"phone_number"`
	EmergencyPhone *string `gorm:"size:50" json:"emergency_phone"`
	
	// education data
	LastEducation      *string         `gorm:"size:255" json:"last_education"`
	EducationInstitute *string         `gorm:"size:255" json:"education_institute"`
	Major              *string         `gorm:"size:255" json:"major"`
	GraduationYear     *string         `gorm:"size:50" json:"graduation_year"`
	Certification      *datatypes.JSON `gorm:"type:json" json:"certification"`

	// family data
	SpouseName *string         `gorm:"size:255" json:"spouse_name"`
	Children   *datatypes.JSON `gorm:"type:json" json:"children"`

	// health data
	BPJSForHealth    *string    `gorm:"size:50" json:"bpjs_for_health"`
	BPJSForWork      *string    `gorm:"size:50" json:"bpjs_for_work"`
	Height           *string    `gorm:"size:50" json:"height"`
	Weight           *string    `gorm:"size:50" json:"weight"`
	DiseaseHistory   *string    `gorm:"size:255" json:"disease_history"`
	LastMedicalCheck *time.Time `json:"last_medical_check"`

	// administrative data
	BankAccountNumber *string  `gorm:"size:50" json:"bank_account_number"`
	BankName          *string  `gorm:"size:255" json:"bank_name"`
	BankAccountName   *string  `gorm:"size:255" json:"bank_account_name"`
	TaxStatus         *string  `gorm:"size:50" json:"tax_status"`
	BasicSalary       *float64 `json:"basic_salary"`
	Allowances        *float64 `json:"allowances"`
	TotalIncome       *float64 `json:"totalIncome"`

	common.BaseModel
}
