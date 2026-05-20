package service

import (
	"context"
	"testing"

	"customer-service/customer/domain"
	"customer-service/customer/repository"
)

func TestServiceCreatesCustomer(t *testing.T) {
	service := New(repository.NewMemoryRepository())

	customer, err := service.Create(context.Background(), domain.CreateCustomerCommand{Name: "Alice"})
	if err != nil {
		t.Fatal(err)
	}
	if customer.ID == "" {
		t.Fatal("expected generated id")
	}
	if customer.Name != "Alice" {
		t.Fatalf("expected Alice, got %q", customer.Name)
	}
}

func TestServiceRejectsMissingName(t *testing.T) {
	service := New(repository.NewMemoryRepository())

	_, err := service.Create(context.Background(), domain.CreateCustomerCommand{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
