package service

import (
	"errors"
	"fmt"
	"time"

	"investory-management-backend/internal/database"
	"investory-management-backend/internal/models"
	"investory-management-backend/internal/repository"

	"gorm.io/gorm"
)

func ListPurchaseOrders(page, limit int, status string) ([]models.PurchaseOrder, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return repository.GetPurchaseOrders(page, limit, status)
}

func GetPurchaseOrder(id string) (*models.PurchaseOrder, error) {
	po, err := repository.GetPurchaseOrderByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("purchase order not found")
	}
	return po, err
}

type CreatePOInput struct {
	SupplierID string
	Notes      string
	ExpectedAt *time.Time
	Items      []CreatePOItemInput
	UserID     uint
}

type CreatePOItemInput struct {
	InventoryItemID string
	QuantityOrdered float64
	CostPerUnit     float64
}

func CreatePurchaseOrder(input CreatePOInput) (*models.PurchaseOrder, error) {
	if len(input.Items) == 0 {
		return nil, errors.New("purchase order must have at least one item")
	}

	poNumber := generatePONumber()
	now := time.Now()
	po := &models.PurchaseOrder{
		PONumber:   poNumber,
		SupplierID: input.SupplierID,
		Status:     models.POStatusDraft,
		OrderedAt:  &now,
		ExpectedAt: input.ExpectedAt,
		Notes:      input.Notes,
		CreatedBy:  &input.UserID,
	}

	for _, i := range input.Items {
		po.Items = append(po.Items, models.PurchaseOrderItem{
			InventoryItemID: i.InventoryItemID,
			QuantityOrdered: i.QuantityOrdered,
			CostPerUnit:     i.CostPerUnit,
		})
	}

	if err := repository.CreatePurchaseOrder(po); err != nil {
		return nil, err
	}
	return repository.GetPurchaseOrderByID(po.ID)
}

func UpdatePurchaseOrder(id, notes string, status models.PurchaseOrderStatus, expectedAt *time.Time) (*models.PurchaseOrder, error) {
	po, err := repository.GetPurchaseOrderByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("purchase order not found")
	}
	if err != nil {
		return nil, err
	}
	if po.Status == models.POStatusReceived || po.Status == models.POStatusCancelled {
		return nil, fmt.Errorf("cannot update a %s purchase order", po.Status)
	}
	if notes != "" {
		po.Notes = notes
	}
	if status != "" {
		po.Status = status
	}
	if expectedAt != nil {
		po.ExpectedAt = expectedAt
	}
	po.Items = nil
	po.Supplier = nil
	if err := repository.UpdatePurchaseOrder(po); err != nil {
		return nil, err
	}
	return repository.GetPurchaseOrderByID(id)
}

type ReceivePOItemInput struct {
	PurchaseOrderItemID string
	QuantityReceived    float64
}

func ReceivePurchaseOrder(id, note string, items []ReceivePOItemInput, userID uint) (*models.PurchaseOrder, error) {
	po, err := repository.GetPurchaseOrderByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("purchase order not found")
	}
	if err != nil {
		return nil, err
	}
	if po.Status == models.POStatusReceived || po.Status == models.POStatusCancelled {
		return nil, fmt.Errorf("cannot receive a %s purchase order", po.Status)
	}

	var updatedItems []*models.InventoryItem
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		for _, input := range items {
			poItem, err := repository.GetPurchaseOrderItemByID(input.PurchaseOrderItemID)
			if err != nil {
				return fmt.Errorf("purchase order item %s not found", input.PurchaseOrderItemID)
			}
			if poItem.PurchaseOrderID != po.ID {
				return errors.New("item does not belong to this purchase order")
			}

			remaining := poItem.QuantityOrdered - poItem.QuantityReceived
			if input.QuantityReceived > remaining {
				return fmt.Errorf("receive quantity exceeds ordered quantity for item %s", input.PurchaseOrderItemID)
			}

			item, err := AdjustStockTx(tx, StockAdjustment{
				InventoryItemID: poItem.InventoryItemID,
				Quantity:        input.QuantityReceived,
				TxType:          models.TxTypeIn,
				ReferenceType:   "PURCHASE_ORDER",
				ReferenceID:     po.ID,
				Note:            note,
				CreatedBy:       &userID,
			})
			if err != nil {
				return err
			}
			updatedItems = append(updatedItems, item)

			poItem.QuantityReceived += input.QuantityReceived
			if err := repository.UpdatePurchaseOrderItem(poItem); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Refresh PO and determine new status
	po, _ = repository.GetPurchaseOrderByID(id)
	allReceived := true
	for _, item := range po.Items {
		if item.QuantityReceived < item.QuantityOrdered {
			allReceived = false
			break
		}
	}
	now := time.Now()
	po.Status = models.POStatusPartiallyReceived
	if allReceived {
		po.Status = models.POStatusReceived
		po.ReceivedAt = &now
	}
	po.Items = nil
	po.Supplier = nil
	repository.UpdatePurchaseOrder(po)

	// Fire webhooks for each updated item
	for _, item := range updatedItems {
		FireWebhookEvent("STOCK_UPDATED", item)
	}

	return repository.GetPurchaseOrderByID(id)
}

func CancelPurchaseOrder(id string) (*models.PurchaseOrder, error) {
	po, err := repository.GetPurchaseOrderByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("purchase order not found")
	}
	if err != nil {
		return nil, err
	}
	if po.Status == models.POStatusReceived {
		return nil, errors.New("cannot cancel a received purchase order")
	}
	if po.Status == models.POStatusCancelled {
		return nil, errors.New("purchase order is already cancelled")
	}
	po.Status = models.POStatusCancelled
	po.Items = nil
	po.Supplier = nil
	if err := repository.UpdatePurchaseOrder(po); err != nil {
		return nil, err
	}
	return repository.GetPurchaseOrderByID(id)
}

func generatePONumber() string {
	today := time.Now().Format("20060102")
	prefix := "PO-" + today + "-"
	var count int64
	database.DB.Model(&models.PurchaseOrder{}).Where("po_number LIKE ?", prefix+"%").Count(&count)
	return fmt.Sprintf("%s%04d", prefix, count+1)
}
