package core

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	chenileerrors "github.com/gauravbhardwaj2code/go_chenile/chenile-framework/base/errors"
)

type recordingInterceptor struct {
	name  string
	steps *[]string
}

func (r recordingInterceptor) Before(ctx context.Context, exchange *Exchange) error {
	*r.steps = append(*r.steps, "before:"+r.name)
	return nil
}

func (r recordingInterceptor) After(ctx context.Context, exchange *Exchange) error {
	*r.steps = append(*r.steps, "after:"+r.name)
	return nil
}

func TestEntryPointInvokesHandlerAndInterceptors(t *testing.T) {
	steps := []string{}
	registry := NewRegistry()
	err := registry.RegisterService(ServiceDefinition{
		ID: "customerService",
		Operations: []OperationDefinition{
			{
				Name: "create",
				Handler: func(ctx context.Context, exchange *Exchange) (any, error) {
					steps = append(steps, "handler")
					return map[string]string{"name": "Alice"}, nil
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	entryPoint := NewEntryPoint(
		registry,
		recordingInterceptor{name: "one", steps: &steps},
		recordingInterceptor{name: "two", steps: &steps},
	)

	result := entryPoint.Execute(context.Background(), &Exchange{
		ServiceID:     "customerService",
		OperationName: "create",
	})

	if !result.Success {
		t.Fatalf("expected success, got %+v", result)
	}
	expectedSteps := []string{"before:one", "before:two", "handler", "after:two", "after:one"}
	if !reflect.DeepEqual(steps, expectedSteps) {
		t.Fatalf("expected steps %v, got %v", expectedSteps, steps)
	}
}

func TestEntryPointMapsChenileErrors(t *testing.T) {
	registry := NewRegistry()
	err := registry.RegisterService(ServiceDefinition{
		ID: "customerService",
		Operations: []OperationDefinition{
			{
				Name: "create",
				Handler: func(ctx context.Context, exchange *Exchange) (any, error) {
					return nil, chenileerrors.New(http.StatusBadRequest, 42, "invalid customer")
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	result := NewEntryPoint(registry).Execute(context.Background(), &Exchange{
		ServiceID:     "customerService",
		OperationName: "create",
	})

	if result.Success {
		t.Fatal("expected failure")
	}
	if result.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", result.Code)
	}
	if result.SubErrorCode != 42 {
		t.Fatalf("expected sub error 42, got %d", result.SubErrorCode)
	}
}

func TestEntryPointMapsUnknownErrorsToInternalServerError(t *testing.T) {
	registry := NewRegistry()
	err := registry.RegisterService(ServiceDefinition{
		ID: "customerService",
		Operations: []OperationDefinition{
			{
				Name: "create",
				Handler: func(ctx context.Context, exchange *Exchange) (any, error) {
					return nil, errors.New("boom")
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	result := NewEntryPoint(registry).Execute(context.Background(), &Exchange{
		ServiceID:     "customerService",
		OperationName: "create",
	})

	if result.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", result.Code)
	}
}
