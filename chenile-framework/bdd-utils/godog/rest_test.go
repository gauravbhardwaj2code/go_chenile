package godogtest

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/cucumber/godog"

	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/core"
	chenilehttp "github.com/gauravbhardwaj2code/go_chenile/chenile-framework/http"
)

func TestSuiteRunsGodogFeature(t *testing.T) {
	router := testRouter(t)
	status := Suite{
		Name:     "godog-utils",
		Router:   router,
		TestingT: t,
		Output:   io.Discard,
		FeaturePaths: []string{
			"features/customer.feature",
		},
	}.Run()

	if status != 0 {
		t.Fatalf("expected godog suite status 0, got %d", status)
	}
}

func TestHarnessReportsAssertionErrors(t *testing.T) {
	harness := NewRESTHarness(testRouter(t))
	harness.reset()

	if err := harness.get("/missing"); err != nil {
		t.Fatal(err)
	}
	if err := harness.statusIs(http.StatusOK); err == nil {
		t.Fatal("expected status assertion error")
	}
}

func TestHarnessRunsPayloadStepDirectly(t *testing.T) {
	harness := NewRESTHarness(testRouter(t))
	harness.reset()

	if err := harness.postWithPayload("/customers", &godog.DocString{Content: `{"name":"Alice"}`}); err != nil {
		t.Fatal(err)
	}
	if err := harness.statusIs(http.StatusOK); err != nil {
		t.Fatal(err)
	}
	if err := harness.successIsTrue(); err != nil {
		t.Fatal(err)
	}
	if err := harness.payloadStringIs("name", "Alice"); err != nil {
		t.Fatal(err)
	}
}

func testRouter(t *testing.T) *chenilehttp.Router {
	t.Helper()
	type createRequest struct {
		Name string `json:"name"`
	}
	registry := core.NewRegistry()
	err := registry.RegisterService(core.ServiceDefinition{
		ID: "customerService",
		Operations: []core.OperationDefinition{
			{
				Name:   "create",
				Method: http.MethodPost,
				Path:   "/customers",
				NewInput: func() any {
					return &createRequest{}
				},
				Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
					input := exchange.Input.(*createRequest)
					return map[string]string{"id": "customer-1", "name": input.Name}, nil
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	entryPoint := core.NewEntryPoint(registry)
	router := chenilehttp.NewRouter(entryPoint)
	if err := router.MountRegistry(registry); err != nil {
		t.Fatal(err)
	}
	return router
}
