# Chenile-Go Module Quick Start

This guide shows the fastest path for creating and running a Chenile-Go service module.

## Prerequisites

- Go 1.22 or later
- This repository checked out locally

## Generate a Service

From the repository root:

```bash
go run ./chenile-framework/servicegen/cmd/chenile-servicegen new --name inventory --out ./chenile-examples
cd chenile-examples/inventory-service
```

The generator creates a runnable service with framework `replace` directives that point back to `../../chenile-framework`.

Generated structure:

```text
inventory-service/
├── cmd/
│   └── inventory-service/
│       └── main.go
├── inventory/
│   ├── contract/
│   │   ├── controller.go
│   │   ├── controller_test.go
│   │   ├── request.go
│   │   └── response.go
│   ├── domain/
│   │   ├── errors.go
│   │   └── model.go
│   ├── module/
│   │   └── module.go
│   ├── repository/
│   │   ├── memory_repository.go
│   │   └── repository.go
│   └── service/
│       ├── service.go
│       └── service_test.go
├── config/
│   └── application.yaml
├── test/
│   ├── features/
│   │   └── inventory.feature
│   ├── fixtures/
│   │   ├── create_inventory.json
│   │   └── create_inventory_missing_name.json
│   └── inventory_service_test.go
└── go.mod
```

Generated services use layered packages: `contract` for HTTP request/response adapters, `domain` for entities and domain errors, `service` for business logic, `repository` for persistence abstraction, and `module` for framework wiring.

## Run Tests

```bash
go test ./...
```

From the repository root, run every framework and example test:

```bash
make test
```

## Run the Service

```bash
go run ./cmd/inventory-service
```

The generated service listens on `:8080`.

Test the generated endpoint from another terminal:

```bash
curl -X POST http://localhost:8080/inventorys \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice"}'
```

Expected response shape:

```json
{
  "code": 200,
  "payload": {
    "id": "inventory-1",
    "name": "Alice"
  },
  "success": true
}
```

The current generator pluralizes routes by appending `s` to the kebab-case service name. For example, `order-item` becomes `/order-items`.

## Generated Code Pattern

The service layer returns a value and an error:

```go
type Service interface {
    Create(context.Context, CreateInventoryRequest) (Inventory, error)
}
```

The controller registers one exact HTTP route:

```go
func Register(registry *core.Registry) error {
    service := NewService()
    return registry.RegisterService(core.ServiceDefinition{
        ID:   "inventoryService",
        Name: "inventoryService",
        Operations: []core.OperationDefinition{{
            Name:   "create",
            Method: http.MethodPost,
            Path:   "/inventorys",
            NewInput: func() any {
                return &CreateInventoryRequest{}
            },
            Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
                request := exchange.Input.(*CreateInventoryRequest)
                return service.Create(ctx, *request)
            },
        }},
    })
}
```

The HTTP router currently matches exact method/path pairs. Template routes such as `/items/{id}` and automatic path parameter extraction are not implemented yet.

## Custom Locations

If you generate a service somewhere other than `chenile-examples`, pass the framework path relative to the generated service directory:

```bash
go run ./chenile-framework/servicegen/cmd/chenile-servicegen new \
  --name payment \
  --out ./tmp/services \
  --framework-root ../../chenile-framework
```

Use the same rule when editing generated `go.mod` files manually: every `replace` path is relative to the service module directory.

## Combine Services

The packager combines independent service modules into one HTTP server:

```go
func NewApp() (*packager.App, error) {
    return packager.NewWebApp(
        packager.Module{Name: "customer", Register: customer.Register},
        packager.Module{Name: "order", Register: order.Register},
    )
}
```

See `chenile-examples/mainweb-app` for the working multi-service example.

## Troubleshooting

- `module not found`: check that `replace` directives point to `../../chenile-framework/<module>` from the service directory.
- `route not found`: the router requires the exact method and path registered in `controller.go`.
- `listen tcp :8080: bind: address already in use`: stop the other process or edit the service `main.go` to listen on a different address.
