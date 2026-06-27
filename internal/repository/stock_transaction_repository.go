package repository

import (
	"investory-management-backend/internal/database"
	"investory-management-backend/internal/models"

	"gorm.io/gorm"
)

func CreateStockTransactionTx(tx *gorm.DB, record *models.StockTransaction) error {
	return tx.Create(record).Error
}

func GetStockTransactions(inventoryItemID string, page, limit int) ([]models.StockTransaction, int64, error) {
	var records []models.StockTransaction
	var total int64

	query := database.DB.Model(&models.StockTransaction{}).Where("inventory_item_id = ?", inventoryItemID)
	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&records).Error
	return records, total, err
}
