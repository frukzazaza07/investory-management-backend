package service

import (
	"errors"

	"investory-management-backend/internal/models"
	"investory-management-backend/internal/repository"

	"gorm.io/gorm"
)

func ListSuppliers(page, limit int, search string) ([]models.Supplier, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return repository.GetSuppliers(page, limit, search)
}

func GetSupplier(id string) (*models.Supplier, error) {
	s, err := repository.GetSupplierByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("supplier not found")
	}
	return s, err
}

func CreateSupplier(name, contactName, phone, email, address string) (*models.Supplier, error) {
	s := &models.Supplier{
		Name:        name,
		ContactName: contactName,
		Phone:       phone,
		Email:       email,
		Address:     address,
	}
	if err := repository.CreateSupplier(s); err != nil {
		return nil, err
	}
	return s, nil
}

func UpdateSupplier(id, name, contactName, phone, email, address string) (*models.Supplier, error) {
	s, err := repository.GetSupplierByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("supplier not found")
	}
	if err != nil {
		return nil, err
	}
	s.Name = name
	s.ContactName = contactName
	s.Phone = phone
	s.Email = email
	s.Address = address
	if err := repository.UpdateSupplier(s); err != nil {
		return nil, err
	}
	return s, nil
}

func DeleteSupplier(id string) error {
	if _, err := repository.GetSupplierByID(id); errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("supplier not found")
	}
	return repository.DeleteSupplier(id)
}
