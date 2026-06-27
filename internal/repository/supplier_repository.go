package repository

import (
	"investory-management-backend/internal/database"
	"investory-management-backend/internal/models"
)

func GetSuppliers(page, limit int, search string) ([]models.Supplier, int64, error) {
	var suppliers []models.Supplier
	var total int64

	query := database.DB.Model(&models.Supplier{})
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	err := query.Offset((page - 1) * limit).Limit(limit).Order("created_at DESC").Find(&suppliers).Error
	return suppliers, total, err
}

func GetSupplierByID(id string) (*models.Supplier, error) {
	var supplier models.Supplier
	err := database.DB.First(&supplier, "id = ?", id).Error
	return &supplier, err
}

func CreateSupplier(supplier *models.Supplier) error {
	return database.DB.Create(supplier).Error
}

func UpdateSupplier(supplier *models.Supplier) error {
	return database.DB.Save(supplier).Error
}

func DeleteSupplier(id string) error {
	return database.DB.Delete(&models.Supplier{}, "id = ?", id).Error
}
