# Go Ajapro

This repository is a Go prototype of the Chenile service framework shape.

It keeps the MVC runtime small:

```text
app service
  -> chenile-go-http
      -> chenile-go-core
          -> chenile-go-base
          -> chenile-go-owiz
```

Developer tooling is separate:

```text
chenile-go-test       # Cucumber-style service tests
chenile-go-servicegen # ready-to-use service skeleton generator
chenile-go-packager   # packaging/manifest validation
```

Framework modules use normal Go unit tests. Generated services use real Godog tests through `chenile-go-test/godog`, so service tests execute `.feature` files while still running with `go test`.

The `.feature` files belong to Godog. They are not Java/Cucumber placeholders. They are executable Gherkin inputs read by `github.com/cucumber/godog` from tests such as `examples/customer-service/test/customer_service_test.go`.

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

## Standalone Service vs Packager

A generated service can run standalone. For example, `customer-service` starts its own web app with only the customer module:

```bash
go run ./examples/customer-service/cmd/customer-service
```

The packager is for the Spring Boot packager style: one main web application imports many service modules and exposes all their routes from one server.

The example main web app combines `customer-service` and `order-service`:

```bash
go run ./examples/mainweb-app/cmd/mainweb-app
```

In code, the mainweb app is just explicit module composition:

```go
app, err := packager.NewWebApp(
    packager.Module{Name: "customer", Register: customer.Register},
    packager.Module{Name: "order", Register: order.Register},
)
```

This is intentionally not runtime classpath scanning. Go requires imported packages at compile time, so the packager/mainweb app imports the service modules it wants to include.

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

```bash
go run ./servicegen/cmd/chenile-servicegen new --name customer --out ./examples
cd examples/customer-service
go test ./...
go run ./cmd/customer-service
```
