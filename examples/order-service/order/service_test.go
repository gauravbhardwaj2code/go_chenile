package order

import (
	"context"
	"testing"
)

func TestServiceCreatesOrder(t *testing.T) {
	service := NewService()

	order, err := service.Create(context.Background(), CreateOrderRequest{Name: "Order 1"})
	if err != nil {
		t.Fatal(err)
	}
	if order.ID == "" {
		t.Fatal("expected generated id")
	}
	if order.Name != "Order 1" {
		t.Fatalf("expected Order 1, got %q", order.Name)
	}
}
