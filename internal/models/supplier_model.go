package models

type Supplier struct {
	Base
	Name        string `gorm:"not null" json:"name"`
	ContactName string `json:"contact_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
}
