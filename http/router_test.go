package chenilehttp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"core"
)

type createRequest struct {
	Name string `json:"name"`
}

func TestRouterBindsJSONAndWritesGenericResponse(t *testing.T) {
	router := testRouter(t)
	request := httptest.NewRequest(http.MethodPost, "/customers?source=test", strings.NewReader(`{"name":"Alice"}`))
	request.Header.Set("X-Tenant", "tenant-a")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body=%s", recorder.Code, recorder.Body.String())
	}
	var result map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["success"] != true {
		t.Fatalf("expected success, got %v", result["success"])
	}
	payload := result["payload"].(map[string]any)
	if payload["name"] != "Alice" {
		t.Fatalf("expected payload name Alice, got %v", payload["name"])
	}
	if payload["tenant"] != "tenant-a" {
		t.Fatalf("expected tenant header to reach handler, got %v", payload["tenant"])
	}
	if payload["source"] != "test" {
		t.Fatalf("expected query param to reach handler, got %v", payload["source"])
	}
}

func TestRouterReturnsNotFoundForUnknownRoute(t *testing.T) {
	router := testRouter(t)
	request := httptest.NewRequest(http.MethodGet, "/missing", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

func TestRouterReturnsBadRequestForInvalidJSON(t *testing.T) {
	router := testRouter(t)
	request := httptest.NewRequest(http.MethodPost, "/customers", strings.NewReader(`{`))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func testRouter(t *testing.T) *Router {
	t.Helper()
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
					return map[string]string{
						"name":   input.Name,
						"tenant": exchange.Headers["X-Tenant"],
						"source": exchange.QueryParams["source"],
					}, nil
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	entryPoint := core.NewEntryPoint(registry)
	router := NewRouter(entryPoint)
	if err := router.MountRegistry(registry); err != nil {
		t.Fatal(err)
	}
	return router
}
