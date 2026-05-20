# Chenile-Go Framework

Chenile-Go is a modular microservices framework for Go, inspired by Chenile's service architecture. It now combines the cleaner explicit service model from the original Go prototype with the important WeGO framework capabilities: structured errors, config loading, middleware/interceptors, state machines, state-entity services, service generation, module packaging, and Godog-based BDD tests.

## Checkout And Prerequisites

```bash
git clone <repo-url>
cd go-ajapro
```

Requirements:

- Go 1.26.3 or later
- Make

This repository is a multi-module Go workspace. There is no root `go.mod`, so run service commands from the module directories shown below.

## Run All Tests

From the repository root:

```bash
make test
```

This runs framework tests and all example module tests:

```bash
(cd chenile-framework && go test ./base/... ./owiz/... ./core/... ./http/... ./bdd-utils/... ./config/... ./middleware/... ./stm/... ./stateentity/... ./servicegen/... ./packager/...)
(cd chenile-examples/customer-service && go test ./...)
(cd chenile-examples/order-service && go test ./...)
(cd chenile-examples/slot-service && GOWORK=$(cd ../../chenile-framework && pwd)/go.work go test ./...)
(cd chenile-examples/state-order-service && go test ./...)
(cd chenile-examples/mainweb-app && go test ./...)
```

Coverage:

```bash
make coverage
```

## Run Godog BDD Tests

Godog is wired through `chenile-framework/bdd-utils/godog` and runs the Gherkin feature files in-process against the HTTP router. No deployed server is required for BDD tests.

Customer BDD:

```bash
(cd chenile-examples/customer-service && go test ./test/... -v)
```

Order BDD:

```bash
(cd chenile-examples/order-service && go test ./test/... -v)
```

Slot BDD:

```bash
(cd chenile-examples/slot-service && GOWORK=$(cd ../../chenile-framework && pwd)/go.work go test ./test/... -v)
```

Feature files:

- `chenile-examples/customer-service/test/features/customer.feature`
- `chenile-examples/order-service/test/features/order.feature`
- `chenile-examples/slot-service/test/features/slot.feature`

The feature steps post JSON to `/customers`, `/orders`, `/runners`, or `/allocations`, assert HTTP status and success fields, and inspect response payload fields.

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

Slot service:

```bash
(cd chenile-examples/slot-service && GOWORK=$(cd ../../chenile-framework && pwd)/go.work go run ./cmd/slot-service)
```

In another terminal:

```bash
curl -X POST http://localhost:8080/runners \
  -H "Content-Type: application/json" \
  -d '{"name":"Asha","skills":["cook"],"attributes":{"diet":"veg"},"slots":[{"date":"2026-06-01","start":"09:00","end":"11:00"}]}'

curl -X POST http://localhost:8080/allocations \
  -H "Content-Type: application/json" \
  -d '{"requestId":"req-1","skill":"cook","slot":{"date":"2026-06-01","start":"09:00","end":"11:00"},"constraints":[{"key":"diet","value":"veg","type":"hard"}]}'
```

## Framework Modules

Core modules:

- `base`: structured errors and standard response envelopes.
- `core`: service registry, operation registry, exchange, entrypoint, and interceptor execution.
- `http`: HTTP router, JSON binding, OpenAPI JSON, and Swagger UI.
- `middleware`: reusable interceptors such as request ID, validation, and panic recovery helpers.
- `config`: lightweight environment and file-backed configuration.
- `stm`: JSON-defined state transition machine with transition actions and automatic states.
- `stateentity`: lifecycle-driven service registration on top of `stm`.
- `bdd-utils`: reusable Godog REST test harness for generated services.
- `packager`: module assembly into runnable apps.
- `servicegen`: BDD-first service generator.

## Run The State Service Example

The state-order service demonstrates `stm`, `stateentity`, and BDD tests:

```bash
(cd chenile-examples/state-order-service && go test ./...)
(cd chenile-examples/state-order-service && go run ./cmd/state-order-service)
```

Create and transition an order:

```bash
curl -X POST http://localhost:8080/state-orders \
  -H "Content-Type: application/json" \
  -d '{"id":"order-1","name":"Order 1"}'

curl -X POST http://localhost:8080/state-orders/event \
  -H "Content-Type: application/json" \
  -d '{"id":"order-1","event":"confirm"}'
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

Use `--public-deps` when the generated service should depend directly on published Chenile module versions without local `replace` directives:

```bash
(cd chenile-framework/servicegen && go run ./cmd/chenile-servicegen new --name slot --out ../../chenile-examples --public-deps)
```

Generated services include a config folder for the port and service name. Service wiring is explicit in Go code.

Each generated service also includes its own `README.md` with run commands. If you copy a generated service before the public Chenile module tags exist, run the included helper first:

```bash
bash scripts/use-local-chenile.sh /absolute/path/to/go_chenile/chenile-framework
go test ./...
go run ./cmd/<service-name>
```

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
| Spring annotations | explicit operation registration | generated `*/contract/controller.go` and `*/module/module.go` |
| Spring DI | constructors and explicit module registration | `NewService()`, `Register(...)` |
| Cucumber JVM | Godog | `chenile-framework/bdd-utils/godog`, `*.feature` files |
| Spring MockMvc | `net/http/httptest` in Godog utilities | `chenile-framework/bdd-utils/godog/rest.go` |
| OpenAPI / Swagger | generated OpenAPI JSON and Swagger UI | `/openapi.json`, `/swagger` |

## More Documentation

- `MODULE_QUICKSTART.md`: service generation and module workflow
- `ARCHITECTURE_GUIDE.md`: framework layers and design
- `DEPENDENCY_MANAGEMENT.md`: Go module layout
- `VALIDATION_REPORT.md`: latest validation summary
