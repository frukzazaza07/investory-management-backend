package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string

const (
	TxTypeIn               TransactionType = "IN"
	TxTypeOut              TransactionType = "OUT"
	TxTypeAdjustmentAdd    TransactionType = "ADJUSTMENT_ADD"
	TxTypeAdjustmentRemove TransactionType = "ADJUSTMENT_REMOVE"
)

type StockTransaction struct {
	ID              string          `gorm:"type:varchar(36);primaryKey" json:"id"`
	InventoryItemID string          `gorm:"type:varchar(36);not null;index" json:"inventory_item_id"`
	TransactionType TransactionType `gorm:"not null" json:"transaction_type"`
	Quantity        float64         `gorm:"not null" json:"quantity"`
	QuantityBefore  float64         `gorm:"not null" json:"quantity_before"`
	QuantityAfter   float64         `gorm:"not null" json:"quantity_after"`
	ReferenceType   string          `gorm:"type:varchar(50)" json:"reference_type"`
	ReferenceID     string          `gorm:"type:varchar(36)" json:"reference_id"`
	Note            string          `json:"note"`
	CreatedBy       *uint           `json:"created_by"`
	CreatedAt       time.Time       `json:"created_at"`
	InventoryItem   *InventoryItem  `gorm:"foreignKey:InventoryItemID" json:"inventory_item,omitempty"`
}

func (s *StockTransaction) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}
