package repository

import (
	"investory-management-backend/internal/database"
	"investory-management-backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetInventoryItems(page, limit int, search string) ([]models.InventoryItem, int64, error) {
	var items []models.InventoryItem
	var total int64

	query := database.DB.Model(&models.InventoryItem{})
	if search != "" {
		query = query.Where("name ILIKE ? OR sku ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	err := query.Offset((page - 1) * limit).Limit(limit).Order("created_at DESC").Find(&items).Error
	return items, total, err
}

func GetInventoryItemByID(id string) (*models.InventoryItem, error) {
	var item models.InventoryItem
	err := database.DB.First(&item, "id = ?", id).Error
	return &item, err
}

func GetInventoryItemForUpdateTx(tx *gorm.DB, id string) (*models.InventoryItem, error) {
	var item models.InventoryItem
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item, "id = ?", id).Error
	return &item, err
}

func CreateInventoryItem(item *models.InventoryItem) error {
	return database.DB.Create(item).Error
}

func UpdateInventoryItem(item *models.InventoryItem) error {
	return database.DB.Save(item).Error
}

func UpdateInventoryStockTx(tx *gorm.DB, id string, qty float64) error {
	return tx.Model(&models.InventoryItem{}).Where("id = ?", id).Update("quantity_in_stock", qty).Error
}

func DeleteInventoryItem(id string) error {
	return database.DB.Delete(&models.InventoryItem{}, "id = ?", id).Error
}

func GetAllInventoryItems() ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	err := database.DB.Find(&items).Error
	return items, err
}
