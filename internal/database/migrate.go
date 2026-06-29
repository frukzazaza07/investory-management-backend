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
		{
			ID: "20260627000001",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(
					&models.Supplier{},
					&models.InventoryItem{},
					&models.Product{},
					&models.ProductInventoryMapping{},
					&models.PurchaseOrder{},
					&models.PurchaseOrderItem{},
					&models.StockTransaction{},
					&models.PosSaleLog{},
					&models.WebhookSubscription{},
					&models.WebhookDeliveryLog{},
				)
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					&models.WebhookDeliveryLog{},
					&models.WebhookSubscription{},
					&models.PosSaleLog{},
					&models.StockTransaction{},
					&models.PurchaseOrderItem{},
					&models.PurchaseOrder{},
					&models.ProductInventoryMapping{},
					&models.Product{},
					&models.InventoryItem{},
					&models.Supplier{},
				)
			},
		},
		{
			ID: "20260629000001",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Product{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropColumn(&models.Product{}, "barcode")
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migrations ran successfully")
}
