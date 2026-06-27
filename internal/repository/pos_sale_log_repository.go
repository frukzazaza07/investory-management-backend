package repository

import (
	"investory-management-backend/internal/database"
	"investory-management-backend/internal/models"
)

func GetPosSaleLogByOrderID(posOrderID string) (*models.PosSaleLog, error) {
	var log models.PosSaleLog
	err := database.DB.Where("pos_order_id = ?", posOrderID).First(&log).Error
	return &log, err
}

func CreatePosSaleLog(log *models.PosSaleLog) error {
	return database.DB.Create(log).Error
}

func UpdatePosSaleLog(log *models.PosSaleLog) error {
	return database.DB.Save(log).Error
}
