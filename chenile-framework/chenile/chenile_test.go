package chenile

import (
	"context"
	"testing"

	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/core"
)

type createRequest struct {
	Name string
}

type createResponse struct {
	Name string
}

func TestBuilderRegistersTypedPostRoute(t *testing.T) {
	registry := core.NewRegistry()
	builder := NewBuilder(registry)

	err := builder.Service("inventoryService").Routes(POST("/inventorys", "create", func() *createRequest {
		return &createRequest{}
	}, func(ctx context.Context, request createRequest) (createResponse, error) {
		return createResponse{Name: request.Name}, nil
	}))
	if err != nil {
		t.Fatal(err)
	}

	operation, ok := registry.Operation("inventoryService", "create")
	if !ok {
		t.Fatal("expected registered operation")
	}
	if operation.Path != "/inventorys" {
		t.Fatalf("unexpected path %q", operation.Path)
	}
	payload, err := operation.Handler(context.Background(), &core.Exchange{Input: &createRequest{Name: "Alice"}})
	if err != nil {
		t.Fatal(err)
	}
	if payload.(createResponse).Name != "Alice" {
		t.Fatalf("unexpected payload %#v", payload)
	}
}
