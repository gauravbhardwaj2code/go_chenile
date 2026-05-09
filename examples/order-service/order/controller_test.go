package order

import (
	"context"
	"testing"

	"github.com/ajapro/chenile-go/core"
)

func TestRegisterAddsCreateOperation(t *testing.T) {
	registry := core.NewRegistry()

	if err := Register(registry); err != nil {
		t.Fatal(err)
	}

	operation, ok := registry.Operation("orderService", "create")
	if !ok {
		t.Fatal("expected create operation")
	}
	if operation.Path != "/orders" {
		t.Fatalf("expected /orders, got %q", operation.Path)
	}
}

func TestRegisteredCreateHandlerInvokesService(t *testing.T) {
	registry := core.NewRegistry()
	if err := Register(registry); err != nil {
		t.Fatal(err)
	}
	operation, ok := registry.Operation("orderService", "create")
	if !ok {
		t.Fatal("expected create operation")
	}

	payload, err := operation.Handler(context.Background(), &core.Exchange{
		Input: &CreateOrderRequest{Name: "Order 1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	order := payload.(Order)
	if order.Name != "Order 1" {
		t.Fatalf("expected Order 1, got %q", order.Name)
	}
}
