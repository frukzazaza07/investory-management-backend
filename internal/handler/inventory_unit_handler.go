package handler

import (
	"investory-management-backend/internal/service"
	"investory-management-backend/pkg/i18n"
	"investory-management-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

// ListInventoryUnits godoc
// @Summary      List inventory units
// @Tags         inventory-units
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} response.Response
// @Router       /api/v1/inventory/units [get]
func ListInventoryUnits(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	units, err := service.ListInventoryUnits()
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "inventory.unit.list"), units)
}

// GetInventoryUnit godoc
// @Summary      Get inventory unit
// @Tags         inventory-units
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Unit ID"
// @Success      200 {object} response.Response
// @Router       /api/v1/inventory/units/{id} [get]
func GetInventoryUnit(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	unit, err := service.GetInventoryUnit(c.Params("id"))
	if err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "inventory.unit.get"), unit)
}

type inventoryUnitRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// CreateInventoryUnit godoc
// @Summary      Create inventory unit (admin only)
// @Tags         inventory-units
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body inventoryUnitRequest true "Inventory Unit"
// @Success      201 {object} response.Response
// @Router       /api/v1/inventory/units [post]
func CreateInventoryUnit(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	var req inventoryUnitRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, i18n.T(lang, "err.invalid_body"))
	}
	if req.Code == "" || req.Name == "" {
		return response.BadRequest(c, i18n.T(lang, "err.unit_code_name_required"))
	}
	unit, err := service.CreateInventoryUnit(req.Code, req.Name)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, i18n.T(lang, "inventory.unit.create"), unit)
}

// UpdateInventoryUnit godoc
// @Summary      Update inventory unit (admin only)
// @Tags         inventory-units
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string               true "Unit ID"
// @Param        body body inventoryUnitRequest true "Inventory Unit"
// @Success      200 {object} response.Response
// @Router       /api/v1/inventory/units/{id} [put]
func UpdateInventoryUnit(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	var req inventoryUnitRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, i18n.T(lang, "err.invalid_body"))
	}
	if req.Code == "" || req.Name == "" {
		return response.BadRequest(c, i18n.T(lang, "err.unit_code_name_required"))
	}
	unit, err := service.UpdateInventoryUnit(c.Params("id"), req.Code, req.Name)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "inventory.unit.update"), unit)
}

// DeleteInventoryUnit godoc
// @Summary      Delete inventory unit (admin only)
// @Tags         inventory-units
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Unit ID"
// @Success      200 {object} response.Response
// @Router       /api/v1/inventory/units/{id} [delete]
func DeleteInventoryUnit(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	if err := service.DeleteInventoryUnit(c.Params("id")); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "inventory.unit.delete"), nil)
}
