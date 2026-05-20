package service

import (
	"context"
	"testing"

	"inventory-service/inventory/domain"
	"inventory-service/inventory/repository"
)

func TestServiceCreatesInventory(t *testing.T) {
	service := New(repository.NewMemoryRepository())

	inventory, err := service.Create(context.Background(), domain.CreateInventoryCommand{Name: "Alice"})
	if err != nil {
		t.Fatal(err)
	}
	if inventory.ID == "" {
		t.Fatal("expected generated id")
	}
	if inventory.Name != "Alice" {
		t.Fatalf("expected Alice, got %q", inventory.Name)
	}
}

func TestServiceRejectsMissingName(t *testing.T) {
	service := New(repository.NewMemoryRepository())

	_, err := service.Create(context.Background(), domain.CreateInventoryCommand{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
