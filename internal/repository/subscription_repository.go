package repository

import (
	"backend/internal/domain/subscription"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	FindByNameSubscription(tx *gorm.DB, name string) (*subscription.SubscriptionModel, error)
	CreateHistorySubscription(tx *gorm.DB, history *subscription.SubscriptionHistoryModel) error
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) FindByNameSubscription(tx *gorm.DB, name string) (*subscription.SubscriptionModel, error) {
	var p subscription.SubscriptionModel
	err := tx.Where("name = ?", name).First(&p).Error
	return &p, err
}

func (r *subscriptionRepository) CreateHistorySubscription(tx *gorm.DB, history *subscription.SubscriptionHistoryModel) error {
	return tx.Create(history).Error
}