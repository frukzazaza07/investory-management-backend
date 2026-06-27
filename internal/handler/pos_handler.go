package handler

import (
	"strconv"

	"investory-management-backend/internal/service"
	"investory-management-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type deductStockRequest struct {
	POSOrderID string                `json:"pos_order_id"`
	Items      []service.POSSaleItem `json:"items"`
}

// DeductStock godoc
// @Summary      Deduct stock after POS sale (idempotent)
// @Tags         pos
// @Accept       json
// @Produce      json
// @Param        X-API-Key header string            true "POS API Key"
// @Param        body      body   deductStockRequest true "Sale info"
// @Success      200 {object} response.Response
// @Router       /api/v1/pos/stock/deduct [post]
func DeductStock(c fiber.Ctx) error {
	var req deductStockRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.POSOrderID == "" {
		return response.BadRequest(c, "pos_order_id is required")
	}
	if len(req.Items) == 0 {
		return response.BadRequest(c, "items are required")
	}

	result, err := service.DeductStockForSale(req.POSOrderID, req.Items)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, "stock deducted", result)
}

// GetStockLevels godoc
// @Summary      Get current stock levels
// @Tags         pos
// @Produce      json
// @Param        X-API-Key header string true "POS API Key"
// @Success      200 {object} response.Response
// @Router       /api/v1/pos/stock/levels [get]
func GetStockLevels(c fiber.Ctx) error {
	levels, err := service.GetStockLevels()
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, "stock levels retrieved", levels)
}

// CheckProductAvailability godoc
// @Summary      Check if a POS product can be sold
// @Tags         pos
// @Produce      json
// @Param        X-API-Key      header string true  "POS API Key"
// @Param        pos_product_id path   string true  "POS Product ID"
// @Param        quantity       query  number false "Quantity (default 1)"
// @Success      200 {object} response.Response
// @Router       /api/v1/pos/products/{pos_product_id}/availability [get]
func CheckProductAvailability(c fiber.Ctx) error {
	posProductID := c.Params("pos_product_id")
	qty := 1.0
	if q := c.Query("quantity"); q != "" {
		if parsed, err := strconv.ParseFloat(q, 64); err == nil {
			qty = parsed
		}
	}

	avail, err := service.CheckProductAvailability(posProductID, qty)
	if err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, "availability checked", avail)
}
