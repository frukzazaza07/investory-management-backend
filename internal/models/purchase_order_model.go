package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseOrderStatus string

const (
	POStatusDraft             PurchaseOrderStatus = "DRAFT"
	POStatusOrdered           PurchaseOrderStatus = "ORDERED"
	POStatusPartiallyReceived PurchaseOrderStatus = "PARTIALLY_RECEIVED"
	POStatusReceived          PurchaseOrderStatus = "RECEIVED"
	POStatusCancelled         PurchaseOrderStatus = "CANCELLED"
)

type PurchaseOrder struct {
	Base
	PONumber   string              `gorm:"uniqueIndex;not null" json:"po_number"`
	SupplierID string              `gorm:"type:varchar(36);not null" json:"supplier_id"`
	Status     PurchaseOrderStatus `gorm:"not null;default:'DRAFT'" json:"status"`
	OrderedAt  *time.Time          `json:"ordered_at"`
	ExpectedAt *time.Time          `json:"expected_at"`
	ReceivedAt *time.Time          `json:"received_at"`
	Notes      string              `json:"notes"`
	CreatedBy  *uint               `json:"created_by"`
	Supplier   *Supplier           `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Items      []PurchaseOrderItem `gorm:"foreignKey:PurchaseOrderID" json:"items,omitempty"`
}

type PurchaseOrderItem struct {
	ID               string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	PurchaseOrderID  string         `gorm:"type:varchar(36);not null;index" json:"purchase_order_id"`
	InventoryItemID  string         `gorm:"type:varchar(36);not null" json:"inventory_item_id"`
	QuantityOrdered  float64        `gorm:"not null" json:"quantity_ordered"`
	QuantityReceived float64        `gorm:"not null;default:0" json:"quantity_received"`
	CostPerUnit      float64        `gorm:"not null;default:0" json:"cost_per_unit"`
	InventoryItem    *InventoryItem `gorm:"foreignKey:InventoryItemID" json:"inventory_item,omitempty"`
}

func (p *PurchaseOrderItem) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}
