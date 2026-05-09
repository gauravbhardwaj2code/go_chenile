package godogtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cucumber/godog"

	chenilehttp "http"
)

type Suite struct {
	Name         string
	Router       *chenilehttp.Router
	FeaturePaths []string
	TestingT     *testing.T
	Output       io.Writer
	Format       string
}

func (s Suite) Run() int {
	harness := NewRESTHarness(s.Router)
	format := s.Format
	if format == "" {
		format = "pretty"
	}
	output := s.Output
	if output == nil {
		output = io.Discard
	}
	return godog.TestSuite{
		Name:                s.Name,
		ScenarioInitializer: harness.Register,
		Options: &godog.Options{
			Format:   format,
			Paths:    s.FeaturePaths,
			Strict:   true,
			TestingT: s.TestingT,
			Output:   output,
		},
	}.Run()
}

type RESTHarness struct {
	router  *chenilehttp.Router
	headers map[string]string
	code    int
	body    []byte
	json    map[string]any
}

func NewRESTHarness(router *chenilehttp.Router) *RESTHarness {
	return &RESTHarness{router: router}
}

func (h *RESTHarness) Register(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(sc *godog.Scenario) {
		h.reset()
	})
	ctx.Step(`^I construct a REST request with header "([^"]*)" and value "([^"]*)"$`, h.setHeader)
	ctx.Step(`^I GET a REST request to URL "([^"]*)"$`, h.get)
	ctx.Step(`^I POST a REST request to URL "([^"]*)" with payload$`, h.postWithPayload)
	ctx.Step(`^I PUT a REST request to URL "([^"]*)" with payload$`, h.putWithPayload)
	ctx.Step(`^I PATCH a REST request to URL "([^"]*)" with payload$`, h.patchWithPayload)
	ctx.Step(`^I DELETE a REST request to URL "([^"]*)" with payload$`, h.deleteWithPayload)
	ctx.Step(`^the http status code is (\d+)$`, h.statusIs)
	ctx.Step(`^success is true$`, h.successIsTrue)
	ctx.Step(`^success is false$`, h.successIsFalse)
	ctx.Step(`^the REST response key "([^"]*)" is "([^"]*)"$`, h.payloadStringIs)
	ctx.Step(`^the REST response contains key "([^"]*)"$`, h.payloadContainsKey)
	ctx.Step(`^the error array size is (\d+)$`, h.errorArraySizeIs)
}

func (h *RESTHarness) reset() {
	h.headers = map[string]string{}
	h.code = 0
	h.body = nil
	h.json = nil
}

func (h *RESTHarness) setHeader(name string, value string) error {
	h.headers[name] = value
	return nil
}

func (h *RESTHarness) get(url string) error {
	return h.do(http.MethodGet, url, "")
}

func (h *RESTHarness) postWithPayload(url string, doc *godog.DocString) error {
	return h.do(http.MethodPost, url, doc.Content)
}

func (h *RESTHarness) putWithPayload(url string, doc *godog.DocString) error {
	return h.do(http.MethodPut, url, doc.Content)
}

func (h *RESTHarness) patchWithPayload(url string, doc *godog.DocString) error {
	return h.do(http.MethodPatch, url, doc.Content)
}

func (h *RESTHarness) deleteWithPayload(url string, doc *godog.DocString) error {
	return h.do(http.MethodDelete, url, doc.Content)
}

func (h *RESTHarness) statusIs(code int) error {
	if h.code != code {
		return fmt.Errorf("expected status %d, got %d; body=%s", code, h.code, string(h.body))
	}
	return nil
}

func (h *RESTHarness) successIsTrue() error {
	return h.successIs(true)
}

func (h *RESTHarness) successIsFalse() error {
	return h.successIs(false)
}

func (h *RESTHarness) successIs(success bool) error {
	actual, _ := h.json["success"].(bool)
	if actual != success {
		return fmt.Errorf("expected success %v, got %v; body=%s", success, actual, string(h.body))
	}
	return nil
}

func (h *RESTHarness) payloadStringIs(path string, value string) error {
	actual := h.payloadValue(path)
	if fmt.Sprint(actual) != value {
		return fmt.Errorf("expected payload.%s to be %q, got %q; body=%s", path, value, fmt.Sprint(actual), string(h.body))
	}
	return nil
}

func (h *RESTHarness) payloadContainsKey(path string) error {
	if h.payloadValue(path) == nil {
		return fmt.Errorf("expected payload.%s to exist; body=%s", path, string(h.body))
	}
	return nil
}

func (h *RESTHarness) errorArraySizeIs(size int) error {
	errorsValue, ok := h.json["errors"].([]any)
	if !ok && size == 0 {
		return nil
	}
	if !ok {
		return fmt.Errorf("expected errors array size %d, but errors is missing; body=%s", size, string(h.body))
	}
	if len(errorsValue) != size {
		return fmt.Errorf("expected errors array size %d, got %d; body=%s", size, len(errorsValue), string(h.body))
	}
	return nil
}

func (h *RESTHarness) do(method string, url string, payload string) error {
	request := httptest.NewRequest(method, url, bytes.NewBufferString(payload))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	for key, value := range h.headers {
		request.Header.Set(key, value)
	}
	recorder := httptest.NewRecorder()
	h.router.ServeHTTP(recorder, request)
	h.code = recorder.Code
	h.body = recorder.Body.Bytes()
	if err := json.Unmarshal(h.body, &h.json); err != nil {
		return fmt.Errorf("response is not valid json: %w; body=%s", err, string(h.body))
	}
	return nil
}

func (h *RESTHarness) payloadValue(path string) any {
	value, ok := h.json["payload"]
	if !ok {
		return nil
	}
	current := value
	for _, part := range strings.Split(path, ".") {
		object, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current = object[part]
	}
	return current
}
