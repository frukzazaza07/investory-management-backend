package models

type InventoryUnit struct {
	Base
	Code string `gorm:"uniqueIndex;not null" json:"code"`
	Name string `gorm:"not null" json:"name"`
}
