package repository

import (
	"investory-management-backend/internal/database"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	result := database.DB.Where("email = ?", email).First(&user)
	return &user, result.Error
}

func CreateUser(user *User) error {
	return database.DB.Create(user).Error
}
