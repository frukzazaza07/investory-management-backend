package repository

import (
	"investory-management-backend/internal/database"
	"investory-management-backend/internal/models"
)

func GetPurchaseOrders(page, limit int, status string) ([]models.PurchaseOrder, int64, error) {
	var orders []models.PurchaseOrder
	var total int64

	query := database.DB.Model(&models.PurchaseOrder{}).Preload("Supplier")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Offset((page - 1) * limit).Limit(limit).Order("created_at DESC").Find(&orders).Error
	return orders, total, err
}

func GetPurchaseOrderByID(id string) (*models.PurchaseOrder, error) {
	var order models.PurchaseOrder
	err := database.DB.
		Preload("Supplier").
		Preload("Items.InventoryItem").
		First(&order, "id = ?", id).Error
	return &order, err
}

func CreatePurchaseOrder(order *models.PurchaseOrder) error {
	return database.DB.Create(order).Error
}

func UpdatePurchaseOrder(order *models.PurchaseOrder) error {
	return database.DB.Save(order).Error
}

func GetPurchaseOrderItemByID(id string) (*models.PurchaseOrderItem, error) {
	var item models.PurchaseOrderItem
	err := database.DB.First(&item, "id = ?", id).Error
	return &item, err
}

func UpdatePurchaseOrderItem(item *models.PurchaseOrderItem) error {
	return database.DB.Save(item).Error
}
