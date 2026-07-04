package service

import (
	"errors"

	"investory-management-backend/internal/models"
	"investory-management-backend/internal/repository"

	"gorm.io/gorm"
)

func ListInventoryUnits() ([]models.InventoryUnit, error) {
	return repository.GetInventoryUnits()
}

func GetInventoryUnit(id string) (*models.InventoryUnit, error) {
	unit, err := repository.GetInventoryUnitByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("inventory unit not found")
	}
	return unit, err
}

func CreateInventoryUnit(code, name string) (*models.InventoryUnit, error) {
	if _, err := repository.GetInventoryUnitByCode(code); err == nil {
		return nil, errors.New("unit code already exists")
	}
	unit := &models.InventoryUnit{Code: code, Name: name}
	if err := repository.CreateInventoryUnit(unit); err != nil {
		return nil, err
	}
	return unit, nil
}

func UpdateInventoryUnit(id, code, name string) (*models.InventoryUnit, error) {
	unit, err := repository.GetInventoryUnitByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("inventory unit not found")
	}
	if err != nil {
		return nil, err
	}
	if code != unit.Code {
		if existing, err := repository.GetInventoryUnitByCode(code); err == nil && existing.ID != unit.ID {
			return nil, errors.New("unit code already exists")
		}
	}
	unit.Code = code
	unit.Name = name
	if err := repository.UpdateInventoryUnit(unit); err != nil {
		return nil, err
	}
	return unit, nil
}

func DeleteInventoryUnit(id string) error {
	unit, err := repository.GetInventoryUnitByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("inventory unit not found")
	}
	if err != nil {
		return err
	}
	count, err := repository.CountInventoryItemsByUnit(unit.Code)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("unit is in use by inventory items")
	}
	return repository.DeleteInventoryUnit(id)
}

func IsValidUnitCode(code string) bool {
	_, err := repository.GetInventoryUnitByCode(code)
	return err == nil
}
