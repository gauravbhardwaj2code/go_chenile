package mainweb

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ajapro/chenile-go/packager"
)

func TestMainwebAppCombinesCustomerAndOrderServices(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatal(err)
	}

	assertRoute(t, app, "/customers", `{"name":"Alice"}`, "Alice")
	assertRoute(t, app, "/orders", `{"name":"Order 1"}`, "Order 1")
}

func assertRoute(t *testing.T, app *packager.App, path string, body string, expected string) {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	recorder := httptest.NewRecorder()

	app.Router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200 for %s, got %d; body=%s", path, recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), expected) {
		t.Fatalf("expected response to contain %q, got %s", expected, recorder.Body.String())
	}
}
