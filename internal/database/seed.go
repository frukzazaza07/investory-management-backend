package database

import (
	"log"

	"investory-management-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func Seed() {
	seedUsers()
	log.Println("Seeding completed")
}

func seedUsers() {
	users := []struct {
		Email    string
		Password string
		Status   string
	}{
		{"admin@example.com", "admin123", "active"},
		{"user@example.com", "user123", "active"},
	}

	for _, u := range users {
		var existing models.User
		if err := DB.Where("email = ?", u.Email).First(&existing).Error; err == nil {
			continue
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", u.Email, err)
			continue
		}

		user := models.User{
			Email:    u.Email,
			Password: string(hashed),
			Status:   u.Status,
		}
		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", u.Email, err)
		}
	}
}
