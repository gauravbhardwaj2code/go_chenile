# Chenile-Go Framework

A modular microservices framework for Go, inspired by Chenile's architecture. Build scalable services with clean separation of concerns, BDD testing, and flexible deployment options.

## Quick Start

### Prerequisites

- Go 1.22 or later
- Make (optional, for running tests via Makefile)

### Run All Tests

```bash
# From the repository root
make test

# Or run tests manually for each module
cd chenile-framework && for dir in base owiz core http test servicegen packager; do cd $dir && go test ./... && cd ..; done
cd chenile-examples/customer-service && go test ./...
cd chenile-examples/order-service && go test ./...
cd chenile-examples/mainweb-app && go test ./...
```

### Run a Standalone Service

```bash
# Customer service (serves /customers endpoint on port 8080)
go run ./chenile-examples/customer-service/cmd/customer-service

# Order service (serves /orders endpoint on port 8080)
go run ./chenile-examples/order-service/cmd/order-service
```

Test the API:
```bash
curl -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice"}'
```

### Run the Combined Application

```bash
# Mainweb app (serves both /customers and /orders on port 8080)
go run ./chenile-examples/mainweb-app/cmd/mainweb-app
```

### Create a New Service

```bash
# Generate a new service named "inventory" in the examples folder
go run ./chenile-framework/servicegen/cmd/chenile-servicegen new --name inventory --out ./chenile-examples

# Navigate to the generated service
cd chenile-examples/inventory-service

# Run tests
go test ./...

# Run the service
go run ./cmd/inventory-service
```

The generator creates:
- Complete service skeleton with model, service, and controller layers
- Unit tests for service and controller
- BDD feature files with Godog integration
- Properly configured go.mod with framework dependencies
- **No config folder** - configuration is done in code for type safety

## Java-to-Go Framework Choices

This project preserves Chenile's service-execution ideas, but it does not try to copy Spring Boot or Java framework mechanics directly. The Go side uses Go-native libraries and explicit registration.

| Java / Chenile side | Go alternative used here | Where | Notes |
| --- | --- | --- | --- |
| Maven parent POM / dependency management | `go.work`, per-module `go.mod`, local `replace` directives | `go.work`, module `go.mod` files | This is the parent-POM equivalent for local development. |
| Maven build lifecycle | `Makefile` plus `go test` | `Makefile` | Keeps build/test commands centralized without recreating Maven. |
| Spring Boot app bootstrap | explicit Go `main` with registry/router setup | `examples/customer-service/cmd/customer-service/main.go` | No auto-scanning or Spring container mimicry. |
| Spring MVC / `@RestController` | `net/http` handler through `chenile-go-http` | `http/router.go` | Routes are mounted from registered operation metadata. |
| Spring annotations such as `@PostMapping` | code/config registration | `customer/controller.go`, service YAML | Go has no annotations; registration is explicit. |
| Spring DI / bean lookup | constructor functions and explicit registration | `NewService()`, `Register(...)` | No DI container used in this prototype. |
| Jackson JSON binding | standard library `encoding/json` | `http/router.go` | Request body binding uses Go structs and JSON tags. |
| OWIZ command chain | generic Go command chain | `owiz/chain` | Preserves Chenile's interceptor-chain idea with Go generics. |
| Chenile `GenericResponse` / errors | Go response and error packages | `base/response`, `base/errors` | Same response contract shape, Go implementation. |
| Cucumber JVM | Godog | `test/godog`, `*.feature` files | Real Go BDD runner for service-level tests. |
| Spring MockMvc | `net/http/httptest` inside Godog utilities | `test/godog/rest.go` | Service tests run in-process without a deployed server. |
| JUnit module tests | Go `testing` package | `*_test.go` | Framework modules use normal Go unit tests. |
| Java service generator concept | Go CLI using `text/template` | `servicegen/cmd/chenile-servicegen` | Generates Go-native skeletons, not Java-style controllers. |

The main Java idea intentionally carried over is Chenile's architecture: service metadata, exchange, interceptor chain, HTTP adapter, generator, and packager. The implementation mechanisms are Go-native.

## Standalone Service vs Multi-Service Application

### Running a Single Service

Each generated service can run independently:

```bash
# Customer service (serves /customers endpoint)
go run ./chenile-examples/customer-service/cmd/customer-service

# Order service (serves /orders endpoint)  
go run ./chenile-examples/order-service/cmd/order-service
```

### Combining Multiple Services

The packager allows you to combine multiple service modules into a single application:

```bash
# Mainweb app (serves both /customers and /orders)
go run ./chenile-examples/mainweb-app/cmd/mainweb-app
```

The mainweb app explicitly composes modules:

```go
// chenile-examples/mainweb-app/app.go
package mainweb

import (
    "packager"
    "customer-service/customer"
    "order-service/order"
)

func NewApp() (*packager.App, error) {
    return packager.NewWebApp(
        packager.Module{Name: "customer", Register: customer.Register},
        packager.Module{Name: "order", Register: order.Register},
    )
}
```

**Key Points:**
- No runtime classpath scanning - Go requires explicit imports at compile time
- Each service module remains independent and testable in isolation
- The packager aggregates all registered operations into a single router
- One HTTP server handles requests for all combined services

## Dependency Management

Java Chenile uses a parent POM for common dependency and module management. This Go prototype uses:

- `go.work` for workspace-wide module management
- one `go.mod` per framework module
- local `replace` directives while the modules live in this workspace

See `DEPENDENCY_MANAGEMENT.md` for the parent-POM equivalent model.

Run this from the repository root:

```bash
make test
```

## Generate a Service

Use the service generator to create a new service skeleton:

```bash
# Basic usage - generates in current directory
go run ./chenile-framework/servicegen/cmd/chenile-servicegen new --name inventory

# Specify output directory
go run ./chenile-framework/servicegen/cmd/chenile-servicegen new --name payment --out ./chenile-examples

# Full example with framework root (for custom locations)
go run ./chenile-framework/servicegen/cmd/chenile-servicegen new \
  --name order \
  --out ./chenile-examples \
  --framework-root ../..
```

### Generated Structure

The generator creates a complete service with the following structure:

```
payment-service/
├── cmd/
│   └── payment-service/
│       └── main.go              # Application entry point
├── payment/
│   ├── model.go                 # Request/response DTOs
│   ├── service.go               # Business logic
│   ├── controller.go            # Operation registration
│   ├── module.go                # Package declaration
│   ├── service_test.go          # Unit tests for service
│   └── controller_test.go       # Unit tests for controller
├── test/
│   ├── payment_service_test.go  # Godog BDD test runner
│   ├── features/
│   │   └── payment.feature      # Gherkin feature file
│   └── fixtures/
│       └── create_payment.json  # Test fixtures
└── go.mod                       # Module dependencies
```

**Important:** The generator does NOT create a config folder. Configuration is done explicitly in Go code for type safety and simpler deployment.

### After Generation

```bash
# Navigate to your service
cd chenile-examples/payment-service

# Run unit tests
go test ./payment/...

# Run BDD tests
go test ./test/...

# Run the service
go run ./cmd/payment-service

# Access the API
curl -X POST http://localhost:8080/payments \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Payment"}'
```

### Customizing Your Service

1. **Edit `payment/model.go`**: Define your request/response structures
2. **Edit `payment/service.go`**: Implement business logic
3. **Edit `payment/controller.go`**: Add more operations/endpoints
4. **Edit `test/features/payment.feature`**: Add BDD scenarios
5. **Update `go.mod`**: Add external dependencies as needed
