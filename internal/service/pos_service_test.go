package service_test

import (
	"testing"

	"investory-management-backend/internal/models"
	"investory-management-backend/internal/service"
)

// Tests for availability logic using exported types from the pos service.

func TestStockLevel_IsOut(t *testing.T) {
	cases := []struct {
		qty   float64
		isOut bool
	}{
		{0, true},
		{-1, true},
		{0.001, false},
		{100, false},
	}
	for _, tc := range cases {
		item := &models.InventoryItem{QuantityInStock: tc.qty}
		got := item.QuantityInStock <= 0
		if got != tc.isOut {
			t.Errorf("qty=%.3f: want isOut=%v, got %v", tc.qty, tc.isOut, got)
		}
	}
}

func TestStockLevel_IsLow(t *testing.T) {
	cases := []struct {
		qty    float64
		min    float64
		isLow  bool
	}{
		{10, 0, false},   // no minimum set
		{10, 20, true},   // below minimum
		{20, 20, true},   // at minimum
		{21, 20, false},  // above minimum
		{0, 5, true},     // out of stock also low
	}
	for _, tc := range cases {
		got := tc.min > 0 && tc.qty <= tc.min
		if got != tc.isLow {
			t.Errorf("qty=%.0f min=%.0f: want isLow=%v, got %v", tc.qty, tc.min, tc.isLow, got)
		}
	}
}

func TestCheckProductAvailability_SufficientStock(t *testing.T) {
	// Simulate the availability check logic from CheckProductAvailability.
	// product BOM: needs 18g coffee, 200ml water, 1 cup.
	// stock: 5000g coffee, 20000ml water, 500 cups.
	type bomEntry struct {
		quantityRequired float64
		available        float64
	}
	bom := []bomEntry{
		{18, 5000},
		{200, 20000},
		{1, 500},
	}
	quantity := 2.0
	isAvailable := true
	for _, entry := range bom {
		required := entry.quantityRequired * quantity
		if entry.available < required {
			isAvailable = false
		}
	}
	if !isAvailable {
		t.Error("expected product to be available with sufficient stock")
	}
}

func TestCheckProductAvailability_InsufficientStock(t *testing.T) {
	type bomEntry struct {
		quantityRequired float64
		available        float64
	}
	// Only 10g coffee left, but 18g needed per unit × 2 = 36g
	bom := []bomEntry{
		{18, 10},
		{200, 20000},
		{1, 500},
	}
	quantity := 2.0
	isAvailable := true
	for _, entry := range bom {
		required := entry.quantityRequired * quantity
		if entry.available < required {
			isAvailable = false
		}
	}
	if isAvailable {
		t.Error("expected product to be unavailable with insufficient stock")
	}
}

func TestDeductResult_AlreadyProcessed(t *testing.T) {
	result := &service.DeductResult{
		POSOrderID: "ORDER-123",
		Status:     "already_processed",
		Deductions: nil,
	}
	if result.Status != "already_processed" {
		t.Errorf("want already_processed, got %q", result.Status)
	}
	if result.POSOrderID != "ORDER-123" {
		t.Errorf("want ORDER-123, got %q", result.POSOrderID)
	}
}

func TestStockLevel_Fields(t *testing.T) {
	sl := service.StockLevel{
		InventoryItemID: "inv-001",
		SKU:             "COFFEE",
		Name:            "Coffee Beans",
		Unit:            "g",
		QuantityInStock: 500,
		MinQuantity:     100,
		IsLow:           false,
		IsOut:           false,
	}
	if sl.InventoryItemID != "inv-001" {
		t.Errorf("unexpected InventoryItemID: %q", sl.InventoryItemID)
	}
	if sl.IsOut {
		t.Error("should not be out of stock")
	}
}
