package service

import (
	"errors"

	"investory-management-backend/internal/database"
	"investory-management-backend/internal/models"
	"investory-management-backend/internal/repository"

	"gorm.io/gorm"
)

func ListInventoryItems(page, limit int, search string) ([]models.InventoryItem, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return repository.GetInventoryItems(page, limit, search)
}

func GetInventoryItem(id string) (*models.InventoryItem, error) {
	item, err := repository.GetInventoryItemByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("inventory item not found")
	}
	return item, err
}

func CreateInventoryItem(sku, name, description, unit string, minQty, costPerUnit float64) (*models.InventoryItem, error) {
	item := &models.InventoryItem{
		SKU:         sku,
		Name:        name,
		Description: description,
		Unit:        unit,
		MinQuantity: minQty,
		CostPerUnit: costPerUnit,
	}
	if err := repository.CreateInventoryItem(item); err != nil {
		return nil, err
	}
	return item, nil
}

func UpdateInventoryItem(id, sku, name, description, unit string, minQty, costPerUnit float64) (*models.InventoryItem, error) {
	item, err := repository.GetInventoryItemByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("inventory item not found")
	}
	if err != nil {
		return nil, err
	}
	item.SKU = sku
	item.Name = name
	item.Description = description
	item.Unit = unit
	item.MinQuantity = minQty
	item.CostPerUnit = costPerUnit
	if err := repository.UpdateInventoryItem(item); err != nil {
		return nil, err
	}
	return item, nil
}

func DeleteInventoryItem(id string) error {
	if _, err := repository.GetInventoryItemByID(id); errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("inventory item not found")
	}
	return repository.DeleteInventoryItem(id)
}

type AdjustStockInput struct {
	InventoryItemID string
	Quantity        float64
	IsAdd           bool
	Note            string
	UserID          uint
}

func AdjustInventoryStock(input AdjustStockInput) (*models.InventoryItem, error) {
	txType := models.TxTypeAdjustmentAdd
	if !input.IsAdd {
		txType = models.TxTypeAdjustmentRemove
	}

	var resultItem *models.InventoryItem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		item, err := AdjustStockTx(tx, StockAdjustment{
			InventoryItemID: input.InventoryItemID,
			Quantity:        input.Quantity,
			TxType:          txType,
			ReferenceType:   "MANUAL",
			Note:            input.Note,
			CreatedBy:       &input.UserID,
		})
		if err != nil {
			return err
		}
		resultItem = item
		return nil
	})
	if err != nil {
		return nil, err
	}

	go FireWebhookEvent("STOCK_UPDATED", resultItem)
	return resultItem, nil
}

func GetInventoryItemTransactions(itemID string, page, limit int) ([]models.StockTransaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	_, err := repository.GetInventoryItemByID(itemID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, errors.New("inventory item not found")
	}
	return repository.GetStockTransactions(itemID, page, limit)
}
