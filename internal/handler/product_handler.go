package handler

import (
	"strconv"

	"investory-management-backend/internal/service"
	"investory-management-backend/pkg/i18n"
	"investory-management-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

// ListProducts godoc
// @Summary      List products
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Param        page   query int    false "Page"
// @Param        limit  query int    false "Limit"
// @Param        search query string false "Search"
// @Success      200 {object} response.Response
// @Router       /api/v1/products [get]
func ListProducts(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	products, total, err := service.ListProducts(page, limit, c.Query("search"))
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "product.list"), fiber.Map{
		"items": products, "total": total, "page": page, "limit": limit,
	})
}

// GetProduct godoc
// @Summary      Get product by ID
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} response.Response
// @Router       /api/v1/products/{id} [get]
func GetProduct(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	product, err := service.GetProduct(c.Params("id"))
	if err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "product.get"), product)
}

// GetProductByBarcode godoc
// @Summary      Get product by barcode
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Param        barcode path string true "Barcode"
// @Success      200 {object} response.Response
// @Router       /api/v1/products/barcode/{barcode} [get]
func GetProductByBarcode(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	product, err := service.GetProductByBarcode(c.Params("barcode"))
	if err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "product.get"), product)
}

type productRequest struct {
	POSProductID string `json:"pos_product_id"`
	Name         string `json:"name"`
	SKU          string `json:"sku"`
	Barcode      string `json:"barcode"`
	IsActive     bool   `json:"is_active"`
}

// CreateProduct godoc
// @Summary      Create product
// @Tags         products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body productRequest true "Product"
// @Success      201 {object} response.Response
// @Router       /api/v1/products [post]
func CreateProduct(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	var req productRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, i18n.T(lang, "err.invalid_body"))
	}
	if req.POSProductID == "" || req.Name == "" {
		return response.BadRequest(c, i18n.T(lang, "err.pos_product_id_name_required"))
	}
	product, err := service.CreateProduct(req.POSProductID, req.Name, req.SKU, req.Barcode, req.IsActive)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Created(c, i18n.T(lang, "product.create"), product)
}

// UpdateProduct godoc
// @Summary      Update product
// @Tags         products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string        true "Product ID"
// @Param        body body productRequest true "Product"
// @Success      200 {object} response.Response
// @Router       /api/v1/products/{id} [put]
func UpdateProduct(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	var req productRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, i18n.T(lang, "err.invalid_body"))
	}
	product, err := service.UpdateProduct(c.Params("id"), req.POSProductID, req.Name, req.SKU, req.Barcode, req.IsActive)
	if err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "product.update"), product)
}

// DeleteProduct godoc
// @Summary      Delete product
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} response.Response
// @Router       /api/v1/products/{id} [delete]
func DeleteProduct(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	if err := service.DeleteProduct(c.Params("id")); err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "product.delete"), nil)
}

// GetProductBOM godoc
// @Summary      Get product BOM (Bill of Materials)
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} response.Response
// @Router       /api/v1/products/{id}/bom [get]
func GetProductBOM(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	bom, err := service.GetProductBOM(c.Params("id"))
	if err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "product.bom_get"), bom)
}

type updateBOMRequest struct {
	Items []service.BOMItem `json:"items"`
}

// UpdateProductBOM godoc
// @Summary      Update product BOM (full replace)
// @Tags         products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string          true "Product ID"
// @Param        body body updateBOMRequest true "BOM items"
// @Success      200 {object} response.Response
// @Router       /api/v1/products/{id}/bom [put]
func UpdateProductBOM(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	var req updateBOMRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, i18n.T(lang, "err.invalid_body"))
	}
	bom, err := service.UpdateProductBOM(c.Params("id"), req.Items)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "product.bom_update"), bom)
}
