package repository

import (
	"investory-management-backend/internal/database"
	"investory-management-backend/internal/models"

	"gorm.io/gorm"
)

func GetProducts(page, limit int, search string) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := database.DB.Model(&models.Product{})
	if search != "" {
		query = query.Where("name ILIKE ? OR pos_product_id ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	err := query.Offset((page - 1) * limit).Limit(limit).Order("created_at DESC").Find(&products).Error
	return products, total, err
}

func GetProductByID(id string) (*models.Product, error) {
	var product models.Product
	err := database.DB.Preload("BOM.InventoryItem").First(&product, "id = ?", id).Error
	return &product, err
}

func GetProductByPOSID(posProductID string) (*models.Product, error) {
	var product models.Product
	err := database.DB.Preload("BOM.InventoryItem").First(&product, "pos_product_id = ?", posProductID).Error
	return &product, err
}

func CreateProduct(product *models.Product) error {
	return database.DB.Create(product).Error
}

func UpdateProduct(product *models.Product) error {
	return database.DB.Save(product).Error
}

func DeleteProduct(id string) error {
	return database.DB.Delete(&models.Product{}, "id = ?", id).Error
}

func GetBOMByProductID(productID string) ([]models.ProductInventoryMapping, error) {
	var bom []models.ProductInventoryMapping
	err := database.DB.Preload("InventoryItem").Where("product_id = ?", productID).Find(&bom).Error
	return bom, err
}

func ReplaceBOM(productID string, items []models.ProductInventoryMapping) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", productID).Delete(&models.ProductInventoryMapping{}).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		return tx.Create(&items).Error
	})
}

func GetProductsByInventoryItemID(inventoryItemID string) ([]models.Product, error) {
	var products []models.Product
	err := database.DB.
		Joins("JOIN product_inventory_mappings pm ON pm.product_id = products.id").
		Where("pm.inventory_item_id = ?", inventoryItemID).
		Find(&products).Error
	return products, err
}
