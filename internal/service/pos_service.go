package service

import (
	"encoding/json"
	"errors"
	"time"

	"investory-management-backend/internal/database"
	"investory-management-backend/internal/models"
	"investory-management-backend/internal/repository"

	"gorm.io/gorm"
)

type POSSaleItem struct {
	POSProductID string  `json:"pos_product_id"`
	Quantity     float64 `json:"quantity"`
}

type DeductResult struct {
	POSOrderID string            `json:"pos_order_id"`
	Status     string            `json:"status"`
	Deductions []StockDeduction  `json:"deductions"`
}

type StockDeduction struct {
	InventoryItemID   string  `json:"inventory_item_id"`
	SKU               string  `json:"sku"`
	Name              string  `json:"name"`
	QuantityDeducted  float64 `json:"quantity_deducted"`
	QuantityRemaining float64 `json:"quantity_remaining"`
}

func DeductStockForSale(posOrderID string, items []POSSaleItem) (*DeductResult, error) {
	payloadBytes, _ := json.Marshal(items)

	// Idempotency check
	existing, err := repository.GetPosSaleLogByOrderID(posOrderID)
	if err == nil {
		if existing.Status == models.SaleStatusProcessed {
			return &DeductResult{POSOrderID: posOrderID, Status: "already_processed"}, nil
		}
		if existing.Status == models.SaleStatusFailed {
			return nil, errors.New("previous attempt failed: " + existing.ErrorMessage)
		}
	}

	saleLog := &models.PosSaleLog{
		PosOrderID: posOrderID,
		Status:     models.SaleStatusPending,
		Payload:    string(payloadBytes),
	}
	if err := repository.CreatePosSaleLog(saleLog); err != nil {
		// Already exists and is PENDING (race condition) - treat as processing
		existing, _ = repository.GetPosSaleLogByOrderID(posOrderID)
		if existing != nil {
			saleLog = existing
		}
	}

	// Collect all inventory deductions grouped by item ID
	type deductionPlan struct {
		inventoryItemID string
		quantity        float64
	}
	deductionMap := map[string]float64{}

	for _, saleItem := range items {
		product, err := repository.GetProductByPOSID(saleItem.POSProductID)
		if err != nil {
			errMsg := "product not found: " + saleItem.POSProductID
			saleLog.Status = models.SaleStatusFailed
			saleLog.ErrorMessage = errMsg
			repository.UpdatePosSaleLog(saleLog)
			return nil, errors.New(errMsg)
		}

		for _, bom := range product.BOM {
			needed := bom.QuantityRequired * saleItem.Quantity
			deductionMap[bom.InventoryItemID] += needed
		}
	}

	var deductions []StockDeduction
	var updatedItems []*models.InventoryItem

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		for itemID, qty := range deductionMap {
			item, err := AdjustStockTx(tx, StockAdjustment{
				InventoryItemID: itemID,
				Quantity:        qty,
				TxType:          models.TxTypeOut,
				ReferenceType:   "POS_SALE",
				ReferenceID:     saleLog.ID,
				Note:            "Sale: " + posOrderID,
			})
			if err != nil {
				return err
			}
			updatedItems = append(updatedItems, item)
			deductions = append(deductions, StockDeduction{
				InventoryItemID:   item.ID,
				SKU:               item.SKU,
				Name:              item.Name,
				QuantityDeducted:  qty,
				QuantityRemaining: item.QuantityInStock,
			})
		}
		return nil
	})

	if err != nil {
		saleLog.Status = models.SaleStatusFailed
		saleLog.ErrorMessage = err.Error()
		repository.UpdatePosSaleLog(saleLog)
		return nil, err
	}

	now := time.Now()
	saleLog.Status = models.SaleStatusProcessed
	saleLog.ProcessedAt = &now
	repository.UpdatePosSaleLog(saleLog)

	// Fire webhooks & check stock alerts
	for _, item := range updatedItems {
		event := "STOCK_UPDATED"
		if item.QuantityInStock <= 0 {
			event = "STOCK_OUT"
		} else if item.MinQuantity > 0 && item.QuantityInStock <= item.MinQuantity {
			event = "STOCK_LOW"
		}
		FireWebhookEvent(event, item)
	}

	return &DeductResult{
		POSOrderID: posOrderID,
		Status:     "processed",
		Deductions: deductions,
	}, nil
}

type StockLevel struct {
	InventoryItemID string  `json:"inventory_item_id"`
	SKU             string  `json:"sku"`
	Name            string  `json:"name"`
	Unit            string  `json:"unit"`
	QuantityInStock float64 `json:"quantity_in_stock"`
	MinQuantity     float64 `json:"min_quantity"`
	IsLow           bool    `json:"is_low"`
	IsOut           bool    `json:"is_out"`
}

func GetStockLevels() ([]StockLevel, error) {
	items, err := repository.GetAllInventoryItems()
	if err != nil {
		return nil, err
	}

	levels := make([]StockLevel, len(items))
	for i, item := range items {
		levels[i] = StockLevel{
			InventoryItemID: item.ID,
			SKU:             item.SKU,
			Name:            item.Name,
			Unit:            item.Unit,
			QuantityInStock: item.QuantityInStock,
			MinQuantity:     item.MinQuantity,
			IsLow:           item.MinQuantity > 0 && item.QuantityInStock <= item.MinQuantity,
			IsOut:           item.QuantityInStock <= 0,
		}
	}
	return levels, nil
}

type ProductAvailability struct {
	POSProductID string              `json:"pos_product_id"`
	Name         string              `json:"name"`
	IsAvailable  bool                `json:"is_available"`
	Details      []AvailabilityDetail `json:"details"`
}

type AvailabilityDetail struct {
	InventoryItemID  string  `json:"inventory_item_id"`
	SKU              string  `json:"sku"`
	Name             string  `json:"name"`
	Required         float64 `json:"required"`
	Available        float64 `json:"available"`
	IsSufficient     bool    `json:"is_sufficient"`
}

func CheckProductAvailability(posProductID string, quantity float64) (*ProductAvailability, error) {
	if quantity <= 0 {
		quantity = 1
	}
	product, err := repository.GetProductByPOSID(posProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	result := &ProductAvailability{
		POSProductID: product.POSProductID,
		Name:         product.Name,
		IsAvailable:  true,
	}

	for _, bom := range product.BOM {
		required := bom.QuantityRequired * quantity
		item, err := repository.GetInventoryItemByID(bom.InventoryItemID)
		if err != nil {
			return nil, err
		}
		sufficient := item.QuantityInStock >= required
		if !sufficient {
			result.IsAvailable = false
		}
		result.Details = append(result.Details, AvailabilityDetail{
			InventoryItemID: item.ID,
			SKU:             item.SKU,
			Name:            item.Name,
			Required:        required,
			Available:       item.QuantityInStock,
			IsSufficient:    sufficient,
		})
	}

	return result, nil
}
