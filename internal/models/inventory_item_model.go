package models

type InventoryItem struct {
	Base
	SKU             string  `gorm:"uniqueIndex;not null" json:"sku"`
	Name            string  `gorm:"not null" json:"name"`
	Description     string  `json:"description"`
	Unit            string  `gorm:"not null" json:"unit"`
	QuantityInStock float64 `gorm:"not null;default:0" json:"quantity_in_stock"`
	MinQuantity     float64 `gorm:"not null;default:0" json:"min_quantity"`
	CostPerUnit     float64 `gorm:"not null;default:0" json:"cost_per_unit"`
}
