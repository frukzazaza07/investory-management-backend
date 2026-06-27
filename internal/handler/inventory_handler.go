package handler

import (
	"strconv"

	"investory-management-backend/internal/service"
	"investory-management-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

// ListInventoryItems godoc
// @Summary      List inventory items
// @Tags         inventory
// @Security     BearerAuth
// @Produce      json
// @Param        page   query int    false "Page"
// @Param        limit  query int    false "Limit"
// @Param        search query string false "Search"
// @Success      200 {object} response.Response
// @Router       /api/v1/inventory/items [get]
func ListInventoryItems(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	items, total, err := service.ListInventoryItems(page, limit, c.Query("search"))
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, "inventory items retrieved", fiber.Map{
		"items": items, "total": total, "page": page, "limit": limit,
	})
}

// GetInventoryItem godoc
// @Summary      Get inventory item
// @Tags         inventory
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Item ID"
// @Success      200 {object} response.Response
// @Router       /api/v1/inventory/items/{id} [get]
func GetInventoryItem(c fiber.Ctx) error {
	item, err := service.GetInventoryItem(c.Params("id"))
	if err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, "inventory item retrieved", item)
}

type inventoryItemRequest struct {
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	MinQuantity float64 `json:"min_quantity"`
	CostPerUnit float64 `json:"cost_per_unit"`
}

// CreateInventoryItem godoc
// @Summary      Create inventory item
// @Tags         inventory
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body inventoryItemRequest true "Inventory Item"
// @Success      201 {object} response.Response
// @Router       /api/v1/inventory/items [post]
func CreateInventoryItem(c fiber.Ctx) error {
	var req inventoryItemRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.SKU == "" || req.Name == "" || req.Unit == "" {
		return response.BadRequest(c, "sku, name and unit are required")
	}
	item, err := service.CreateInventoryItem(req.SKU, req.Name, req.Description, req.Unit, req.MinQuantity, req.CostPerUnit)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Created(c, "inventory item created", item)
}

// UpdateInventoryItem godoc
// @Summary      Update inventory item
// @Tags         inventory
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string               true "Item ID"
// @Param        body body inventoryItemRequest true "Inventory Item"
// @Success      200 {object} response.Response
// @Router       /api/v1/inventory/items/{id} [put]
func UpdateInventoryItem(c fiber.Ctx) error {
	var req inventoryItemRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	item, err := service.UpdateInventoryItem(c.Params("id"), req.SKU, req.Name, req.Description, req.Unit, req.MinQuantity, req.CostPerUnit)
	if err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, "inventory item updated", item)
}

// DeleteInventoryItem godoc
// @Summary      Delete inventory item
// @Tags         inventory
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Item ID"
// @Success      200 {object} response.Response
// @Router       /api/v1/inventory/items/{id} [delete]
func DeleteInventoryItem(c fiber.Ctx) error {
	if err := service.DeleteInventoryItem(c.Params("id")); err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, "inventory item deleted", nil)
}

type adjustStockRequest struct {
	Quantity float64 `json:"quantity"`
	IsAdd    bool    `json:"is_add"`
	Note     string  `json:"note"`
}

// AdjustInventoryStock godoc
// @Summary      Adjust inventory stock manually
// @Tags         inventory
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string             true "Item ID"
// @Param        body body adjustStockRequest true "Adjustment"
// @Success      200 {object} response.Response
// @Router       /api/v1/inventory/items/{id}/adjust [post]
func AdjustInventoryStock(c fiber.Ctx) error {
	var req adjustStockRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Quantity <= 0 {
		return response.BadRequest(c, "quantity must be greater than 0")
	}
	userID := c.Locals("userID").(uint)
	item, err := service.AdjustInventoryStock(service.AdjustStockInput{
		InventoryItemID: c.Params("id"),
		Quantity:        req.Quantity,
		IsAdd:           req.IsAdd,
		Note:            req.Note,
		UserID:          userID,
	})
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, "stock adjusted", item)
}

// GetInventoryItemTransactions godoc
// @Summary      Get stock transaction history
// @Tags         inventory
// @Security     BearerAuth
// @Produce      json
// @Param        id    path  string true  "Item ID"
// @Param        page  query int    false "Page"
// @Param        limit query int    false "Limit"
// @Success      200 {object} response.Response
// @Router       /api/v1/inventory/items/{id}/transactions [get]
func GetInventoryItemTransactions(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	txs, total, err := service.GetInventoryItemTransactions(c.Params("id"), page, limit)
	if err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, "transactions retrieved", fiber.Map{
		"items": txs, "total": total, "page": page, "limit": limit,
	})
}
