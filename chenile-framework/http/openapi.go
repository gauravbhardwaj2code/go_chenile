package chenilehttp

import (
	"encoding/json"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"core"
)

type openAPIDocument struct {
	OpenAPI string                     `json:"openapi"`
	Info    openAPIInfo                `json:"info"`
	Paths   map[string]openAPIPathItem `json:"paths"`
}

type openAPIInfo struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type openAPIPathItem map[string]openAPIOperation

type openAPIOperation struct {
	OperationID string                     `json:"operationId"`
	Summary     string                     `json:"summary"`
	Tags        []string                   `json:"tags,omitempty"`
	RequestBody *openAPIRequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]openAPIResponse `json:"responses"`
}

type openAPIRequestBody struct {
	Required bool                              `json:"required"`
	Content  map[string]openAPIMediaTypeObject `json:"content"`
}

type openAPIMediaTypeObject struct {
	Schema openAPISchema `json:"schema"`
}

type openAPIResponse struct {
	Description string                            `json:"description"`
	Content     map[string]openAPIMediaTypeObject `json:"content,omitempty"`
}

type openAPISchema struct {
	Type       string                   `json:"type,omitempty"`
	Format     string                   `json:"format,omitempty"`
	Properties map[string]openAPISchema `json:"properties,omitempty"`
	Items      *openAPISchema           `json:"items,omitempty"`
	Required   []string                 `json:"required,omitempty"`
	Additional *openAPISchema           `json:"additionalProperties,omitempty"`
	Ref        string                   `json:"$ref,omitempty"`
}

func buildOpenAPIDocument(registry *core.Registry) openAPIDocument {
	document := openAPIDocument{
		OpenAPI: "3.0.3",
		Info: openAPIInfo{
			Title:   "Chenile-Go API",
			Version: "0.0.1",
		},
		Paths: map[string]openAPIPathItem{},
	}
	operations := registry.Operations()
	sort.Slice(operations, func(i, j int) bool {
		left := operations[i].Operation.Path + " " + operations[i].Operation.Method
		right := operations[j].Operation.Path + " " + operations[j].Operation.Method
		return left < right
	})
	for _, registered := range operations {
		path := registered.Operation.Path
		method := strings.ToLower(registered.Operation.Method)
		if document.Paths[path] == nil {
			document.Paths[path] = openAPIPathItem{}
		}
		operation := openAPIOperation{
			OperationID: registered.Service.ID + "." + registered.Operation.Name,
			Summary:     registered.Operation.Name,
			Tags:        []string{registered.Service.ID},
			Responses: map[string]openAPIResponse{
				"200": {
					Description: "Successful response",
					Content: map[string]openAPIMediaTypeObject{
						"application/json": {
							Schema: genericResponseSchema(),
						},
					},
				},
			},
		}
		if registered.Operation.NewInput != nil {
			operation.RequestBody = &openAPIRequestBody{
				Required: true,
				Content: map[string]openAPIMediaTypeObject{
					"application/json": {
						Schema: schemaFor(reflect.TypeOf(registered.Operation.NewInput())),
					},
				},
			}
		}
		document.Paths[path][method] = operation
	}
	return document
}

func genericResponseSchema() openAPISchema {
	return openAPISchema{
		Type: "object",
		Properties: map[string]openAPISchema{
			"code":        {Type: "integer", Format: "int32"},
			"description": {Type: "string"},
			"errors": {
				Type: "array",
				Items: &openAPISchema{
					Type: "object",
					Properties: map[string]openAPISchema{
						"code":        {Type: "integer", Format: "int32"},
						"description": {Type: "string"},
					},
				},
			},
			"payload": {Type: "object"},
			"success": {Type: "boolean"},
		},
	}
}

func schemaFor(t reflect.Type) openAPISchema {
	if t == nil {
		return openAPISchema{Type: "object"}
	}
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Struct:
		return structSchema(t)
	case reflect.String:
		return openAPISchema{Type: "string"}
	case reflect.Bool:
		return openAPISchema{Type: "boolean"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return openAPISchema{Type: "integer", Format: "int32"}
	case reflect.Int64:
		return openAPISchema{Type: "integer", Format: "int64"}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return openAPISchema{Type: "integer", Format: "int32"}
	case reflect.Uint64:
		return openAPISchema{Type: "integer", Format: "int64"}
	case reflect.Float32:
		return openAPISchema{Type: "number", Format: "float"}
	case reflect.Float64:
		return openAPISchema{Type: "number", Format: "double"}
	case reflect.Slice, reflect.Array:
		item := schemaFor(t.Elem())
		return openAPISchema{Type: "array", Items: &item}
	case reflect.Map:
		value := schemaFor(t.Elem())
		return openAPISchema{Type: "object", Additional: &value}
	default:
		return openAPISchema{Type: "object"}
	}
}

func structSchema(t reflect.Type) openAPISchema {
	properties := map[string]openAPISchema{}
	required := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		name, omitEmpty, ok := jsonFieldName(field)
		if !ok {
			continue
		}
		properties[name] = schemaFor(field.Type)
		if !omitEmpty {
			required = append(required, name)
		}
	}
	sort.Strings(required)
	return openAPISchema{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}
}

func jsonFieldName(field reflect.StructField) (string, bool, bool) {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return "", false, false
	}
	if tag == "" {
		return field.Name, false, true
	}
	parts := strings.Split(tag, ",")
	name := parts[0]
	if name == "" {
		name = field.Name
	}
	for _, option := range parts[1:] {
		if option == "omitempty" {
			return name, true, true
		}
	}
	return name, false, true
}

func writeJSON(writer http.ResponseWriter, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(value)
}

func writeSwaggerUI(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(swaggerHTML))
}

const swaggerHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Chenile-Go Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/openapi.json",
        dom_id: "#swagger-ui"
      });
    };
  </script>
</body>
</html>`
