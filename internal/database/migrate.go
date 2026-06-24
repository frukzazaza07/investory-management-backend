package database

import (
	"log"

	"investory-management-backend/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migrate() {
	m := gormigrate.New(DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20260624000002",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.User{})
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migrations ran successfully")
}
