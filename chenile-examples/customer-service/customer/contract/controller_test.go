package contract

import (
	"context"
	"testing"

	"customer-service/customer/domain"
)

type fakeService struct{}

func (fakeService) Create(ctx context.Context, command domain.CreateCustomerCommand) (domain.Customer, error) {
	return domain.Customer{ID: "customer-1", Name: command.Name}, nil
}

func TestRoutesDeclareCreateOperation(t *testing.T) {
	routes := Routes(fakeService{})

	if len(routes) != 1 {
		t.Fatalf("expected one route, got %d", len(routes))
	}
	if routes[0].Name != "create" {
		t.Fatalf("expected create operation, got %q", routes[0].Name)
	}
	if routes[0].Path != "/customers" {
		t.Fatalf("expected /customers, got %q", routes[0].Path)
	}
	if routes[0].NewInput == nil {
		t.Fatal("expected input factory")
	}
}

func TestCreateInvokesService(t *testing.T) {
	controller := NewController(fakeService{})

	payload, err := controller.Create(context.Background(), CreateCustomerRequest{Name: "Alice"})
	if err != nil {
		t.Fatal(err)
	}
	if payload.Name != "Alice" {
		t.Fatalf("expected Alice, got %q", payload.Name)
	}
}
