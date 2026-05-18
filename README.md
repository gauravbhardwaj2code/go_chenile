# Chenile-Go Framework

Chenile-Go is a modular microservices framework for Go, inspired by Chenile's service architecture. It demonstrates explicit service registration, a shared HTTP adapter, standardized responses, service generation, module packaging, and Godog-based BDD tests.

## Checkout And Prerequisites

```bash
git clone <repo-url>
cd go-ajapro
```

Requirements:

- Go 1.22 or later
- Make

This repository is a multi-module Go workspace. There is no root `go.mod`, so run service commands from the module directories shown below.

## Run All Tests

From the repository root:

```bash
make test
```

This runs framework tests and all example module tests:

```bash
(cd chenile-framework && go test ./base/... ./owiz/... ./core/... ./http/... ./test/... ./servicegen/... ./packager/...)
(cd chenile-examples/customer-service && go test ./...)
(cd chenile-examples/order-service && go test ./...)
(cd chenile-examples/mainweb-app && go test ./...)
```

Coverage:

```bash
make coverage
```

## Run Godog BDD Tests

Godog is wired through `chenile-framework/test/godog` and runs the Gherkin feature files in-process against the HTTP router. No deployed server is required for BDD tests.

Customer BDD:

```bash
(cd chenile-examples/customer-service && go test ./test/... -v)
```

Order BDD:

```bash
(cd chenile-examples/order-service && go test ./test/... -v)
```

Feature files:

- `chenile-examples/customer-service/test/features/customer.feature`
- `chenile-examples/order-service/test/features/order.feature`

The feature steps post JSON to `/customers` or `/orders`, assert HTTP 200, assert `success` is true, and inspect response payload fields.

## Run Independent Services

Run one service at a time because both examples listen on `:8080`.

Customer service:

```bash
(cd chenile-examples/customer-service && go run ./cmd/customer-service)
```

In another terminal:

```bash
curl -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice"}'
```

Order service:

```bash
(cd chenile-examples/order-service && go run ./cmd/order-service)
```

In another terminal:

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"name":"Order 1"}'
```

## Run The Combined Web App

The mainweb app packages customer and order modules into one HTTP server on `:8080`.

```bash
(cd chenile-examples/mainweb-app && go run ./cmd/mainweb-app)
```

In another terminal, test both routes:

```bash
curl -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice"}'

curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"name":"Order 1"}'
```

Swagger/OpenAPI is available on every running service and combined app:

```bash
curl http://localhost:8080/openapi.json
```

Open the Swagger UI in a browser:

```text
http://localhost:8080/swagger
```

The combined app is assembled in `chenile-examples/mainweb-app/app.go`:

```go
func NewApp() (*packager.App, error) {
    return packager.NewWebApp(
        packager.Module{Name: "customer", Register: customer.Register},
        packager.Module{Name: "order", Register: order.Register},
    )
}
```

## Generate A New Service

From the repository root:

```bash
(cd chenile-framework/servicegen && go run ./cmd/chenile-servicegen new --name inventory --out ../../chenile-examples)
cd chenile-examples/inventory-service
go test ./...
go run ./cmd/inventory-service
```

Test the generated endpoint:

```bash
curl -X POST http://localhost:8080/inventorys \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice"}'
```

The generator creates:

- `cmd/<service>/main.go`
- model, service, controller, and module files
- unit tests for service and controller
- Godog BDD runner and Gherkin feature file
- `go.mod` with local framework `replace` directives

There is no config folder. Services are wired explicitly in Go code.

## Architecture Map

```text
service modules
  -> packager
  -> http
  -> core
  -> base
  -> owiz
```

| Java / Chenile side | Go alternative used here | Where |
| --- | --- | --- |
| Maven parent POM | `go.work`, per-module `go.mod`, local `replace` directives | `chenile-framework/go.work`, module `go.mod` files |
| Spring Boot bootstrap | explicit Go `main` with packager setup | `chenile-examples/*/cmd/*/main.go` |
| Spring MVC controller | `net/http` router backed by operation metadata | `chenile-framework/http/router.go` |
| Spring annotations | explicit operation registration | `customer/controller.go`, `order/controller.go` |
| Spring DI | constructors and explicit module registration | `NewService()`, `Register(...)` |
| Cucumber JVM | Godog | `chenile-framework/test/godog`, `*.feature` files |
| Spring MockMvc | `net/http/httptest` in Godog utilities | `chenile-framework/test/godog/rest.go` |
| OpenAPI / Swagger | generated OpenAPI JSON and Swagger UI | `/openapi.json`, `/swagger` |

## More Documentation

- `MODULE_QUICKSTART.md`: service generation and module workflow
- `ARCHITECTURE_GUIDE.md`: framework layers and design
- `DEPENDENCY_MANAGEMENT.md`: Go module layout
- `VALIDATION_REPORT.md`: latest validation summary
