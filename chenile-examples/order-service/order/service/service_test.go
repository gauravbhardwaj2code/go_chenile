package service

import (
	"context"
	"testing"

	"order-service/order/domain"
	"order-service/order/repository"
)

func TestServiceCreatesOrder(t *testing.T) {
	service := New(repository.NewMemoryRepository())

	order, err := service.Create(context.Background(), domain.CreateOrderCommand{Name: "Alice"})
	if err != nil {
		t.Fatal(err)
	}
	if order.ID == "" {
		t.Fatal("expected generated id")
	}
	if order.Name != "Alice" {
		t.Fatalf("expected Alice, got %q", order.Name)
	}
}

func TestServiceRejectsMissingName(t *testing.T) {
	service := New(repository.NewMemoryRepository())

	_, err := service.Create(context.Background(), domain.CreateOrderCommand{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
