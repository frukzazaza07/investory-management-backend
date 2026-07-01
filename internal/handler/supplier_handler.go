package handler

import (
	"strconv"

	"investory-management-backend/internal/service"
	"investory-management-backend/pkg/i18n"
	"investory-management-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

// ListSuppliers godoc
// @Summary      List suppliers
// @Tags         suppliers
// @Security     BearerAuth
// @Produce      json
// @Param        page   query int    false "Page"
// @Param        limit  query int    false "Limit"
// @Param        search query string false "Search"
// @Success      200 {object} response.Response
// @Router       /api/v1/suppliers [get]
func ListSuppliers(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	suppliers, total, err := service.ListSuppliers(page, limit, c.Query("search"))
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "supplier.list"), fiber.Map{
		"items": suppliers, "total": total, "page": page, "limit": limit,
	})
}

// GetSupplier godoc
// @Summary      Get supplier by ID
// @Tags         suppliers
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Supplier ID"
// @Success      200 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /api/v1/suppliers/{id} [get]
func GetSupplier(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	supplier, err := service.GetSupplier(c.Params("id"))
	if err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "supplier.get"), supplier)
}

type supplierRequest struct {
	Name        string `json:"name"`
	ContactName string `json:"contact_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
}

// CreateSupplier godoc
// @Summary      Create supplier
// @Tags         suppliers
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body supplierRequest true "Supplier"
// @Success      201 {object} response.Response
// @Router       /api/v1/suppliers [post]
func CreateSupplier(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	var req supplierRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, i18n.T(lang, "err.invalid_body"))
	}
	if req.Name == "" {
		return response.BadRequest(c, i18n.T(lang, "err.name_required"))
	}
	supplier, err := service.CreateSupplier(req.Name, req.ContactName, req.Phone, req.Email, req.Address)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Created(c, i18n.T(lang, "supplier.create"), supplier)
}

// UpdateSupplier godoc
// @Summary      Update supplier
// @Tags         suppliers
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string         true "Supplier ID"
// @Param        body body supplierRequest true "Supplier"
// @Success      200 {object} response.Response
// @Router       /api/v1/suppliers/{id} [put]
func UpdateSupplier(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	var req supplierRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, i18n.T(lang, "err.invalid_body"))
	}
	supplier, err := service.UpdateSupplier(c.Params("id"), req.Name, req.ContactName, req.Phone, req.Email, req.Address)
	if err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "supplier.update"), supplier)
}

// DeleteSupplier godoc
// @Summary      Delete supplier
// @Tags         suppliers
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Supplier ID"
// @Success      200 {object} response.Response
// @Router       /api/v1/suppliers/{id} [delete]
func DeleteSupplier(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	if err := service.DeleteSupplier(c.Params("id")); err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "supplier.delete"), nil)
}
