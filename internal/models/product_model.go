package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	Base
	POSProductID string                    `gorm:"uniqueIndex;not null" json:"pos_product_id"`
	Name         string                    `gorm:"not null" json:"name"`
	SKU          string                    `json:"sku"`
	IsActive     bool                      `gorm:"not null;default:true" json:"is_active"`
	BOM          []ProductInventoryMapping `gorm:"foreignKey:ProductID" json:"bom,omitempty"`
}

type ProductInventoryMapping struct {
	ID               string        `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProductID        string        `gorm:"type:varchar(36);not null;index" json:"product_id"`
	InventoryItemID  string        `gorm:"type:varchar(36);not null" json:"inventory_item_id"`
	QuantityRequired float64       `gorm:"not null" json:"quantity_required"`
	InventoryItem    *InventoryItem `gorm:"foreignKey:InventoryItemID" json:"inventory_item,omitempty"`
}

func (p *ProductInventoryMapping) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}
