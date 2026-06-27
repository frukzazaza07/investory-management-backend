package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SaleStatus string

const (
	SaleStatusPending   SaleStatus = "PENDING"
	SaleStatusProcessed SaleStatus = "PROCESSED"
	SaleStatusFailed    SaleStatus = "FAILED"
)

type PosSaleLog struct {
	ID           string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	PosOrderID   string     `gorm:"uniqueIndex;not null" json:"pos_order_id"`
	Status       SaleStatus `gorm:"not null;default:'PENDING'" json:"status"`
	Payload      string     `gorm:"type:text" json:"payload"`
	ProcessedAt  *time.Time `json:"processed_at"`
	ErrorMessage string     `json:"error_message"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (p *PosSaleLog) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}
