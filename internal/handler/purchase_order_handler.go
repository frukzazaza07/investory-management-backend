package handler

import (
	"strconv"
	"time"

	"investory-management-backend/internal/models"
	"investory-management-backend/internal/service"
	"investory-management-backend/pkg/i18n"
	"investory-management-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

// ListPurchaseOrders godoc
// @Summary      List purchase orders
// @Tags         purchase-orders
// @Security     BearerAuth
// @Produce      json
// @Param        page   query int    false "Page"
// @Param        limit  query int    false "Limit"
// @Param        status query string false "Status filter"
// @Success      200 {object} response.Response
// @Router       /api/v1/purchase-orders [get]
func ListPurchaseOrders(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	orders, total, err := service.ListPurchaseOrders(page, limit, c.Query("status"))
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "po.list"), fiber.Map{
		"items": orders, "total": total, "page": page, "limit": limit,
	})
}

// GetPurchaseOrder godoc
// @Summary      Get purchase order by ID
// @Tags         purchase-orders
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "PO ID"
// @Success      200 {object} response.Response
// @Router       /api/v1/purchase-orders/{id} [get]
func GetPurchaseOrder(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	po, err := service.GetPurchaseOrder(c.Params("id"))
	if err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "po.get"), po)
}

type poItemRequest struct {
	InventoryItemID string  `json:"inventory_item_id"`
	QuantityOrdered float64 `json:"quantity_ordered"`
	CostPerUnit     float64 `json:"cost_per_unit"`
}

type createPORequest struct {
	SupplierID string          `json:"supplier_id"`
	Notes      string          `json:"notes"`
	ExpectedAt string          `json:"expected_at"` // RFC3339
	Items      []poItemRequest `json:"items"`
}

// CreatePurchaseOrder godoc
// @Summary      Create purchase order
// @Tags         purchase-orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body createPORequest true "Purchase Order"
// @Success      201 {object} response.Response
// @Router       /api/v1/purchase-orders [post]
func CreatePurchaseOrder(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	var req createPORequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, i18n.T(lang, "err.invalid_body"))
	}
	if req.SupplierID == "" {
		return response.BadRequest(c, i18n.T(lang, "err.supplier_id_required"))
	}

	input := service.CreatePOInput{
		SupplierID: req.SupplierID,
		Notes:      req.Notes,
		UserID:     c.Locals("userID").(uint),
	}
	if req.ExpectedAt != "" {
		t, err := time.Parse(time.RFC3339, req.ExpectedAt)
		if err == nil {
			input.ExpectedAt = &t
		}
	}
	for _, item := range req.Items {
		input.Items = append(input.Items, service.CreatePOItemInput{
			InventoryItemID: item.InventoryItemID,
			QuantityOrdered: item.QuantityOrdered,
			CostPerUnit:     item.CostPerUnit,
		})
	}

	po, err := service.CreatePurchaseOrder(input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, i18n.T(lang, "po.create"), po)
}

type updatePORequest struct {
	Notes      string                     `json:"notes"`
	Status     models.PurchaseOrderStatus `json:"status"`
	ExpectedAt string                     `json:"expected_at"`
}

// UpdatePurchaseOrder godoc
// @Summary      Update purchase order
// @Tags         purchase-orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string          true "PO ID"
// @Param        body body updatePORequest true "Update"
// @Success      200 {object} response.Response
// @Router       /api/v1/purchase-orders/{id} [put]
func UpdatePurchaseOrder(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	var req updatePORequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, i18n.T(lang, "err.invalid_body"))
	}
	var expectedAt *time.Time
	if req.ExpectedAt != "" {
		t, err := time.Parse(time.RFC3339, req.ExpectedAt)
		if err == nil {
			expectedAt = &t
		}
	}
	po, err := service.UpdatePurchaseOrder(c.Params("id"), req.Notes, req.Status, expectedAt)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "po.update"), po)
}

type receiveItemRequest struct {
	PurchaseOrderItemID string  `json:"purchase_order_item_id"`
	QuantityReceived    float64 `json:"quantity_received"`
}

type receivePORequest struct {
	Items []receiveItemRequest `json:"items"`
	Note  string               `json:"note"`
}

// ReceivePurchaseOrder godoc
// @Summary      Receive goods for purchase order
// @Tags         purchase-orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string           true "PO ID"
// @Param        body body receivePORequest true "Receive"
// @Success      200 {object} response.Response
// @Router       /api/v1/purchase-orders/{id}/receive [post]
func ReceivePurchaseOrder(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	var req receivePORequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, i18n.T(lang, "err.invalid_body"))
	}
	if len(req.Items) == 0 {
		return response.BadRequest(c, i18n.T(lang, "err.items_required"))
	}

	receiveItems := make([]service.ReceivePOItemInput, len(req.Items))
	for i, item := range req.Items {
		receiveItems[i] = service.ReceivePOItemInput{
			PurchaseOrderItemID: item.PurchaseOrderItemID,
			QuantityReceived:    item.QuantityReceived,
		}
	}

	po, err := service.ReceivePurchaseOrder(c.Params("id"), req.Note, receiveItems, c.Locals("userID").(uint))
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "po.receive"), po)
}

// CancelPurchaseOrder godoc
// @Summary      Cancel purchase order
// @Tags         purchase-orders
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "PO ID"
// @Success      200 {object} response.Response
// @Router       /api/v1/purchase-orders/{id}/cancel [post]
func CancelPurchaseOrder(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	po, err := service.CancelPurchaseOrder(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "po.cancel"), po)
}
