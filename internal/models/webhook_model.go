package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebhookSubscription struct {
	Base
	Name     string `gorm:"not null" json:"name"`
	URL      string `gorm:"not null" json:"url"`
	Secret   string `gorm:"not null" json:"-"`
	Events   string `gorm:"type:text;not null" json:"events"` // JSON array: ["STOCK_UPDATED","STOCK_LOW"]
	IsActive bool   `gorm:"not null;default:true" json:"is_active"`
}

type WebhookDeliveryLog struct {
	ID                    string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	WebhookSubscriptionID string    `gorm:"type:varchar(36);not null;index" json:"webhook_subscription_id"`
	EventType             string    `gorm:"not null" json:"event_type"`
	Payload               string    `gorm:"type:text" json:"payload"`
	ResponseStatus        int       `json:"response_status"`
	ResponseBody          string    `json:"response_body"`
	IsSuccess             bool      `json:"is_success"`
	CreatedAt             time.Time `json:"created_at"`
}

func (w *WebhookDeliveryLog) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return nil
}
