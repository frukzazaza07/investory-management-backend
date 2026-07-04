package repository

import (
	"investory-management-backend/internal/database"
	"investory-management-backend/internal/models"
)

func GetInventoryUnits() ([]models.InventoryUnit, error) {
	var units []models.InventoryUnit
	err := database.DB.Order("code ASC").Find(&units).Error
	return units, err
}

func GetInventoryUnitByID(id string) (*models.InventoryUnit, error) {
	var unit models.InventoryUnit
	err := database.DB.First(&unit, "id = ?", id).Error
	return &unit, err
}

func GetInventoryUnitByCode(code string) (*models.InventoryUnit, error) {
	var unit models.InventoryUnit
	err := database.DB.First(&unit, "code = ?", code).Error
	return &unit, err
}

func CreateInventoryUnit(unit *models.InventoryUnit) error {
	return database.DB.Create(unit).Error
}

func UpdateInventoryUnit(unit *models.InventoryUnit) error {
	return database.DB.Save(unit).Error
}

func DeleteInventoryUnit(id string) error {
	return database.DB.Delete(&models.InventoryUnit{}, "id = ?", id).Error
}

func CountInventoryItemsByUnit(code string) (int64, error) {
	var count int64
	err := database.DB.Model(&models.InventoryItem{}).Where("unit = ?", code).Count(&count).Error
	return count, err
}
