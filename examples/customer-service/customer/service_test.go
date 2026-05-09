package customer

import (
	"context"
	"testing"
)

func TestServiceCreatesCustomer(t *testing.T) {
	service := NewService()

	customer, err := service.Create(context.Background(), CreateCustomerRequest{Name: "Alice"})
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
