package service

import (
	"fmt"

	"investory-management-backend/internal/models"
	"investory-management-backend/internal/repository"

	"gorm.io/gorm"
)

type StockAdjustment struct {
	InventoryItemID string
	Quantity        float64
	TxType          models.TransactionType
	ReferenceType   string
	ReferenceID     string
	Note            string
	CreatedBy       *uint
}

// AdjustStockTx performs a stock movement within an existing transaction.
// Returns the updated item so callers can inspect post-transaction quantity.
func AdjustStockTx(tx *gorm.DB, adj StockAdjustment) (*models.InventoryItem, error) {
	item, err := repository.GetInventoryItemForUpdateTx(tx, adj.InventoryItemID)
	if err != nil {
		return nil, fmt.Errorf("inventory item not found: %s", adj.InventoryItemID)
	}

	before := item.QuantityInStock
	var after float64

	switch adj.TxType {
	case models.TxTypeIn, models.TxTypeAdjustmentAdd:
		after = before + adj.Quantity
	case models.TxTypeOut, models.TxTypeAdjustmentRemove:
		if before < adj.Quantity {
			return nil, fmt.Errorf("insufficient stock for %s: have %.4g, need %.4g", item.SKU, before, adj.Quantity)
		}
		after = before - adj.Quantity
	default:
		return nil, fmt.Errorf("unknown transaction type: %s", adj.TxType)
	}

	if err := repository.UpdateInventoryStockTx(tx, item.ID, after); err != nil {
		return nil, err
	}
	item.QuantityInStock = after

	record := &models.StockTransaction{
		InventoryItemID: item.ID,
		TransactionType: adj.TxType,
		Quantity:        adj.Quantity,
		QuantityBefore:  before,
		QuantityAfter:   after,
		ReferenceType:   adj.ReferenceType,
		ReferenceID:     adj.ReferenceID,
		Note:            adj.Note,
		CreatedBy:       adj.CreatedBy,
	}
	if err := repository.CreateStockTransactionTx(tx, record); err != nil {
		return nil, err
	}

	return item, nil
}
