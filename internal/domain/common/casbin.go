package common

type CasbinRule struct {
    ID    uint   `gorm:"primaryKey;autoIncrement"`
    Ptype string `gorm:"size:100;index"` // jenis policy: p (policy), g (grouping)

    V0 string `gorm:"size:100;index"` // subject (user/role)
    V1 string `gorm:"size:100;index"` // object (resource)
    V2 string `gorm:"size:100;index"` // action (read/write/update/delete)
    V3 string `gorm:"size:100;index"` // domain (companyId → untuk multi-tenancy)
    V4 string `gorm:"size:100;index"`
    V5 string `gorm:"size:100;index"`
}

func (CasbinRule) TableName() string {
	return "casbin_rule"
}