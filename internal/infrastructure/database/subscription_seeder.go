package database

import (
	"backend/internal/domain/common"
	"backend/internal/domain/subscription"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedSubscription(db *gorm.DB) error {
	time := time.Now()
	packages := []subscription.SubscriptionModel{
		{
			ID:           uuid.New(),
			Name:         "Basic",
			Description:  ptr("Default basic package with trial 30 days"),
			Price:        0,
			DurationDays: 30, // trial 30 hari
			BaseModel: common.BaseModel{
				CreatedAt: time,
				UpdatedAt: time,
			},
		},
		{
			ID:           uuid.New(),
			Name:         "Premium",
			Description:  ptr("Premium package with more features"),
			Price:        1499000, // contoh harga
			DurationDays: 30,
			BaseModel: common.BaseModel{
				CreatedAt: time,
				UpdatedAt: time,
			},
		},
		{
			ID:           uuid.New(),
			Name:         "Enterprise",
			Description:  ptr("Enterprise package with full features"),
			Price:        2999000, // contoh harga
			DurationDays: 30,
			BaseModel: common.BaseModel{
				CreatedAt: time,
				UpdatedAt: time,
			},
		},
	}

	for _, p := range packages {
		var existing subscription.SubscriptionModel
		if err := db.Where("name = ?", p.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&p).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func ptr(s string) *string {
	return &s
}
