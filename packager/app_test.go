package packager

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"core"
)

func TestNewWebAppCombinesMultipleModules(t *testing.T) {
	app, err := NewWebApp(
		testModule("customer", "customerService", "/customers"),
		testModule("order", "orderService", "/orders"),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertRoute(t, app, "/customers", `{"name":"Alice"}`, "customerService")
	assertRoute(t, app, "/orders", `{"name":"Order 1"}`, "orderService")
}

func TestNewWebAppRejectsInvalidModules(t *testing.T) {
	if _, err := NewWebApp(Module{}); err == nil {
		t.Fatal("expected missing module name error")
	}
	if _, err := NewWebApp(Module{Name: "customer"}); err == nil {
		t.Fatal("expected missing register function error")
	}
}

func testModule(name string, serviceID string, path string) Module {
	return Module{
		Name: name,
		Register: func(registry *core.Registry) error {
			return registry.RegisterService(core.ServiceDefinition{
				ID: serviceID,
				Operations: []core.OperationDefinition{
					{
						Name:   "create",
						Method: http.MethodPost,
						Path:   path,
						NewInput: func() any {
							return &struct {
								Name string `json:"name"`
							}{}
						},
						Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
							return map[string]string{"service": serviceID}, nil
						},
					},
				},
			})
		},
	}
}

func assertRoute(t *testing.T, app *App, path string, body string, serviceID string) {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	recorder := httptest.NewRecorder()

	app.Router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 for %s, got %d; body=%s", path, recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), serviceID) {
		t.Fatalf("expected response to contain %s, got %s", serviceID, recorder.Body.String())
	}
}
