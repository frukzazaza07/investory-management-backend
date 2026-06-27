package service

import (
	"errors"

	"investory-management-backend/internal/models"
	"investory-management-backend/internal/repository"

	"gorm.io/gorm"
)

func ListProducts(page, limit int, search string) ([]models.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return repository.GetProducts(page, limit, search)
}

func GetProduct(id string) (*models.Product, error) {
	p, err := repository.GetProductByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("product not found")
	}
	return p, err
}

func CreateProduct(posProductID, name, sku string, isActive bool) (*models.Product, error) {
	p := &models.Product{
		POSProductID: posProductID,
		Name:         name,
		SKU:          sku,
		IsActive:     isActive,
	}
	if err := repository.CreateProduct(p); err != nil {
		return nil, err
	}
	return p, nil
}

func UpdateProduct(id, posProductID, name, sku string, isActive bool) (*models.Product, error) {
	p, err := repository.GetProductByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}
	p.POSProductID = posProductID
	p.Name = name
	p.SKU = sku
	p.IsActive = isActive
	p.BOM = nil // don't overwrite BOM via Save
	if err := repository.UpdateProduct(p); err != nil {
		return nil, err
	}
	return repository.GetProductByID(id)
}

func DeleteProduct(id string) error {
	if _, err := repository.GetProductByID(id); errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("product not found")
	}
	return repository.DeleteProduct(id)
}

type BOMItem struct {
	InventoryItemID  string  `json:"inventory_item_id"`
	QuantityRequired float64 `json:"quantity_required"`
}

func GetProductBOM(productID string) ([]models.ProductInventoryMapping, error) {
	if _, err := repository.GetProductByID(productID); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("product not found")
	}
	return repository.GetBOMByProductID(productID)
}

func UpdateProductBOM(productID string, items []BOMItem) ([]models.ProductInventoryMapping, error) {
	if _, err := repository.GetProductByID(productID); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("product not found")
	}

	mappings := make([]models.ProductInventoryMapping, len(items))
	for i, item := range items {
		if item.QuantityRequired <= 0 {
			return nil, errors.New("quantity_required must be greater than 0")
		}
		mappings[i] = models.ProductInventoryMapping{
			ProductID:        productID,
			InventoryItemID:  item.InventoryItemID,
			QuantityRequired: item.QuantityRequired,
		}
	}

	if err := repository.ReplaceBOM(productID, mappings); err != nil {
		return nil, err
	}
	return repository.GetBOMByProductID(productID)
}
