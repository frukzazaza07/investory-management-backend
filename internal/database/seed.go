package database

import (
	"log"

	"investory-management-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func Seed() {
	seedUsers()
	seedSuppliers()
	seedInventoryUnits()
	seedInventoryItems()
	seedProducts()
	log.Println("Seeding completed")
}

func seedUsers() {
	users := []struct {
		Email    string
		Password string
		Status   string
		Role     string
	}{
		{"admin@example.com", "admin123", "active", models.RoleAdmin},
		{"user@example.com", "user123", "active", models.RoleStaff},
	}

	for _, u := range users {
		var existing models.User
		if err := DB.Where("email = ?", u.Email).First(&existing).Error; err == nil {
			if existing.Role != u.Role {
				DB.Model(&existing).Update("role", u.Role)
			}
			continue
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", u.Email, err)
			continue
		}
		user := models.User{Email: u.Email, Password: string(hashed), Status: u.Status, Role: u.Role}
		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", u.Email, err)
		}
	}
}

func seedInventoryUnits() {
	units := []models.InventoryUnit{
		{Base: models.Base{ID: "unit-g"}, Code: "g", Name: "Gram"},
		{Base: models.Base{ID: "unit-kg"}, Code: "kg", Name: "Kilogram"},
		{Base: models.Base{ID: "unit-ml"}, Code: "ml", Name: "Milliliter"},
		{Base: models.Base{ID: "unit-l"}, Code: "l", Name: "Liter"},
		{Base: models.Base{ID: "unit-piece"}, Code: "piece", Name: "Piece"},
		{Base: models.Base{ID: "unit-pack"}, Code: "pack", Name: "Pack"},
		{Base: models.Base{ID: "unit-box"}, Code: "box", Name: "Box"},
		{Base: models.Base{ID: "unit-bottle"}, Code: "bottle", Name: "Bottle"},
	}
	for _, u := range units {
		var existing models.InventoryUnit
		if DB.Where("id = ?", u.ID).First(&existing).Error == nil {
			continue
		}
		if err := DB.Create(&u).Error; err != nil {
			log.Printf("Failed to seed inventory unit %s: %v", u.Code, err)
		}
	}
}

func seedSuppliers() {
	suppliers := []models.Supplier{
		{Base: models.Base{ID: "sup-coffee-world-0001"}, Name: "Coffee World Co.", ContactName: "Somchai", Phone: "0812345678", Email: "order@coffeeworld.th", Address: "123 Silom Rd, Bangkok"},
		{Base: models.Base{ID: "sup-dairy-farm-0001"}, Name: "Fresh Dairy Farm", ContactName: "Malee", Phone: "0898765432", Email: "supply@dairyfarm.th", Address: "456 Rama 9, Bangkok"},
	}
	for _, s := range suppliers {
		var existing models.Supplier
		if DB.Where("id = ?", s.ID).First(&existing).Error == nil {
			continue
		}
		if err := DB.Create(&s).Error; err != nil {
			log.Printf("Failed to seed supplier %s: %v", s.Name, err)
		}
	}
}

func seedInventoryItems() {
	items := []models.InventoryItem{
		{Base: models.Base{ID: "inv-coffee-beans-001"}, SKU: "RAW-COFFEE-BEANS", Name: "Coffee Beans (Arabica)", Description: "Premium Arabica coffee beans", Unit: "g", QuantityInStock: 5000, MinQuantity: 500, CostPerUnit: 0.5},
		{Base: models.Base{ID: "inv-milk-001"}, SKU: "RAW-MILK-FRESH", Name: "Fresh Milk", Description: "Full-fat fresh milk", Unit: "ml", QuantityInStock: 10000, MinQuantity: 1000, CostPerUnit: 0.04},
		{Base: models.Base{ID: "inv-sugar-001"}, SKU: "RAW-SUGAR-WHITE", Name: "White Sugar", Description: "Refined white sugar", Unit: "g", QuantityInStock: 3000, MinQuantity: 300, CostPerUnit: 0.02},
		{Base: models.Base{ID: "inv-cup-hot-001"}, SKU: "PKG-CUP-HOT-12OZ", Name: "Hot Cup 12oz", Description: "Paper hot cup 12oz", Unit: "piece", QuantityInStock: 500, MinQuantity: 50, CostPerUnit: 2.5},
		{Base: models.Base{ID: "inv-water-001"}, SKU: "RAW-WATER-FILTER", Name: "Filtered Water", Description: "RO filtered water", Unit: "ml", QuantityInStock: 20000, MinQuantity: 2000, CostPerUnit: 0.001},
	}
	for _, item := range items {
		var existing models.InventoryItem
		if DB.Where("id = ?", item.ID).First(&existing).Error == nil {
			continue
		}
		if err := DB.Create(&item).Error; err != nil {
			log.Printf("Failed to seed inventory item %s: %v", item.Name, err)
		}
	}
}

func seedProducts() {
	products := []struct {
		product models.Product
		bom     []models.ProductInventoryMapping
	}{
		{
			product: models.Product{Base: models.Base{ID: "prod-americano-001"}, POSProductID: "pos-americano", Name: "Americano", SKU: "BEV-AMERICANO", Barcode: "8850001000016", IsActive: true},
			bom: []models.ProductInventoryMapping{
				{ID: "bom-ame-coffee-001", ProductID: "prod-americano-001", InventoryItemID: "inv-coffee-beans-001", QuantityRequired: 18},
				{ID: "bom-ame-water-001", ProductID: "prod-americano-001", InventoryItemID: "inv-water-001", QuantityRequired: 200},
				{ID: "bom-ame-cup-001", ProductID: "prod-americano-001", InventoryItemID: "inv-cup-hot-001", QuantityRequired: 1},
			},
		},
		{
			product: models.Product{Base: models.Base{ID: "prod-latte-001"}, POSProductID: "pos-latte", Name: "Cafe Latte", SKU: "BEV-LATTE", Barcode: "8850001000023", IsActive: true},
			bom: []models.ProductInventoryMapping{
				{ID: "bom-lat-coffee-001", ProductID: "prod-latte-001", InventoryItemID: "inv-coffee-beans-001", QuantityRequired: 18},
				{ID: "bom-lat-milk-001", ProductID: "prod-latte-001", InventoryItemID: "inv-milk-001", QuantityRequired: 200},
				{ID: "bom-lat-sugar-001", ProductID: "prod-latte-001", InventoryItemID: "inv-sugar-001", QuantityRequired: 5},
				{ID: "bom-lat-cup-001", ProductID: "prod-latte-001", InventoryItemID: "inv-cup-hot-001", QuantityRequired: 1},
			},
		},
		{
			product: models.Product{Base: models.Base{ID: "prod-cappuccino-001"}, POSProductID: "pos-cappuccino", Name: "Cappuccino", SKU: "BEV-CAPPUCCINO", Barcode: "8850001000030", IsActive: true},
			bom: []models.ProductInventoryMapping{
				{ID: "bom-cap-coffee-001", ProductID: "prod-cappuccino-001", InventoryItemID: "inv-coffee-beans-001", QuantityRequired: 18},
				{ID: "bom-cap-milk-001", ProductID: "prod-cappuccino-001", InventoryItemID: "inv-milk-001", QuantityRequired: 120},
				{ID: "bom-cap-water-001", ProductID: "prod-cappuccino-001", InventoryItemID: "inv-water-001", QuantityRequired: 60},
				{ID: "bom-cap-cup-001", ProductID: "prod-cappuccino-001", InventoryItemID: "inv-cup-hot-001", QuantityRequired: 1},
			},
		},
		{
			product: models.Product{Base: models.Base{ID: "prod-espresso-001"}, POSProductID: "pos-espresso", Name: "Espresso", SKU: "BEV-ESPRESSO", Barcode: "8850001000047", IsActive: true},
			bom: []models.ProductInventoryMapping{
				{ID: "bom-esp-coffee-001", ProductID: "prod-espresso-001", InventoryItemID: "inv-coffee-beans-001", QuantityRequired: 18},
				{ID: "bom-esp-water-001", ProductID: "prod-espresso-001", InventoryItemID: "inv-water-001", QuantityRequired: 30},
				{ID: "bom-esp-cup-001", ProductID: "prod-espresso-001", InventoryItemID: "inv-cup-hot-001", QuantityRequired: 1},
			},
		},
	}

	for _, p := range products {
		var existing models.Product
		if DB.Where("id = ?", p.product.ID).First(&existing).Error == nil {
			continue
		}
		if err := DB.Create(&p.product).Error; err != nil {
			log.Printf("Failed to seed product %s: %v", p.product.Name, err)
			continue
		}
		for _, bom := range p.bom {
			var existingBOM models.ProductInventoryMapping
			if DB.Where("id = ?", bom.ID).First(&existingBOM).Error == nil {
				continue
			}
			if err := DB.Create(&bom).Error; err != nil {
				log.Printf("Failed to seed BOM %s: %v", bom.ID, err)
			}
		}
	}
}
