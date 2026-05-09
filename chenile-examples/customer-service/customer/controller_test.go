package customer

import (
	"context"
	"testing"

	"core"
)

func TestRegisterAddsCreateOperation(t *testing.T) {
	registry := core.NewRegistry()

	if err := Register(registry); err != nil {
		t.Fatal(err)
	}

	operation, ok := registry.Operation("customerService", "create")
	if !ok {
		t.Fatal("expected create operation")
	}
	if operation.Path != "/customers" {
		t.Fatalf("expected /customers, got %q", operation.Path)
	}
	if operation.NewInput == nil {
		t.Fatal("expected input factory")
	}
}

func TestRegisteredCreateHandlerInvokesService(t *testing.T) {
	registry := core.NewRegistry()
	if err := Register(registry); err != nil {
		t.Fatal(err)
	}
	operation, ok := registry.Operation("customerService", "create")
	if !ok {
		t.Fatal("expected create operation")
	}

	payload, err := operation.Handler(context.Background(), &core.Exchange{
		Input: &CreateCustomerRequest{Name: "Alice"},
	})
	if err != nil {
		t.Fatal(err)
	}
	customer := payload.(Customer)
	if customer.Name != "Alice" {
		t.Fatalf("expected Alice, got %q", customer.Name)
	}
}
