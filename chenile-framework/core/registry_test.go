package core

import (
	"context"
	"net/http"
	"testing"
)

func TestRegistryRegistersAndFindsOperations(t *testing.T) {
	registry := NewRegistry()
	err := registry.RegisterService(ServiceDefinition{
		ID:   "customerService",
		Name: "customerService",
		Operations: []OperationDefinition{
			{
				Name:   "create",
				Method: http.MethodPost,
				Path:   "/customers",
				Handler: func(ctx context.Context, exchange *Exchange) (any, error) {
					return "ok", nil
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	operation, ok := registry.Operation("customerService", "create")
	if !ok {
		t.Fatal("expected operation")
	}
	if operation.Path != "/customers" {
		t.Fatalf("unexpected path %q", operation.Path)
	}
}

func TestRegistryRejectsDuplicateServices(t *testing.T) {
	registry := NewRegistry()
	service := ServiceDefinition{ID: "customerService"}

	if err := registry.RegisterService(service); err != nil {
		t.Fatal(err)
	}
	if err := registry.RegisterService(service); err == nil {
		t.Fatal("expected duplicate service error")
	}
}

func TestRegistryRejectsMissingServiceID(t *testing.T) {
	registry := NewRegistry()

	if err := registry.RegisterService(ServiceDefinition{}); err == nil {
		t.Fatal("expected missing service id error")
	}
}
