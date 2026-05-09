package chenilehttp

import (
	"encoding/json"
	"io"
	"net/http"

	"base/response"
	"core"
)

type route struct {
	serviceID     string
	operationName string
	operation     core.OperationDefinition
}

type Router struct {
	entryPoint *core.EntryPoint
	routes     map[string]route
}

func NewRouter(entryPoint *core.EntryPoint) *Router {
	return &Router{
		entryPoint: entryPoint,
		routes:     map[string]route{},
	}
}

func (r *Router) MountRegistry(registry *core.Registry) error {
	for _, registered := range registry.Operations() {
		key := routeKey(registered.Operation.Method, registered.Operation.Path)
		r.routes[key] = route{
			serviceID:     registered.Service.ID,
			operationName: registered.Operation.Name,
			operation:     registered.Operation,
		}
	}
	return nil
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	route, ok := r.routes[routeKey(request.Method, request.URL.Path)]
	if !ok {
		writeResponse(writer, response.Failure(http.StatusNotFound, "route not found", 0))
		return
	}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		writeResponse(writer, response.Failure(http.StatusBadRequest, err.Error(), 0))
		return
	}
	exchange := &core.Exchange{
		Context:       request.Context(),
		ServiceID:     route.serviceID,
		OperationName: route.operationName,
		Method:        request.Method,
		Path:          request.URL.Path,
		Headers:       headers(request),
		QueryParams:   queryParams(request),
		Body:          body,
	}
	if route.operation.NewInput != nil && len(body) > 0 {
		input := route.operation.NewInput()
		if err := json.Unmarshal(body, input); err != nil {
			writeResponse(writer, response.Failure(http.StatusBadRequest, err.Error(), 0))
			return
		}
		exchange.Input = input
	}
	result := r.entryPoint.Execute(request.Context(), exchange)
	writeResponse(writer, result)
}

func routeKey(method string, path string) string {
	return method + " " + path
}

func headers(request *http.Request) map[string]string {
	values := map[string]string{}
	for key, headerValues := range request.Header {
		if len(headerValues) > 0 {
			values[key] = headerValues[0]
		}
	}
	return values
}

func queryParams(request *http.Request) map[string]string {
	values := map[string]string{}
	for key, queryValues := range request.URL.Query() {
		if len(queryValues) > 0 {
			values[key] = queryValues[0]
		}
	}
	return values
}

func writeResponse(writer http.ResponseWriter, result response.GenericResponse) {
	status := result.Code
	if status == 0 {
		status = http.StatusOK
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(result)
}
