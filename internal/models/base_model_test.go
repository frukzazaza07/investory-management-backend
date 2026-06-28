package models_test

import (
	"testing"

	"investory-management-backend/internal/models"
)

func TestBase_BeforeCreate_GeneratesUUID(t *testing.T) {
	b := &models.Base{}
	if err := b.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate error: %v", err)
	}
	if b.ID == "" {
		t.Error("expected ID to be set, got empty string")
	}
	if len(b.ID) != 36 {
		t.Errorf("expected UUID length 36, got %d", len(b.ID))
	}
}

func TestBase_BeforeCreate_DoesNotOverwriteExistingID(t *testing.T) {
	b := &models.Base{ID: "existing-id"}
	if err := b.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate error: %v", err)
	}
	if b.ID != "existing-id" {
		t.Errorf("expected ID to remain existing-id, got %q", b.ID)
	}
}

func TestBase_BeforeCreate_UniquePerCall(t *testing.T) {
	b1 := &models.Base{}
	b2 := &models.Base{}
	_ = b1.BeforeCreate(nil)
	_ = b2.BeforeCreate(nil)
	if b1.ID == b2.ID {
		t.Error("expected different UUIDs for different Base instances")
	}
}

func TestProductInventoryMapping_BeforeCreate_GeneratesUUID(t *testing.T) {
	m := &models.ProductInventoryMapping{}
	if err := m.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate error: %v", err)
	}
	if m.ID == "" {
		t.Error("expected ID to be set, got empty string")
	}
	if len(m.ID) != 36 {
		t.Errorf("expected UUID length 36, got %d", len(m.ID))
	}
}

func TestProductInventoryMapping_BeforeCreate_DoesNotOverwriteExistingID(t *testing.T) {
	m := &models.ProductInventoryMapping{ID: "my-bom-id"}
	if err := m.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate error: %v", err)
	}
	if m.ID != "my-bom-id" {
		t.Errorf("expected ID to remain my-bom-id, got %q", m.ID)
	}
}
