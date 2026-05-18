# Dependency Management

Chenile Java modules inherit dependency versions and plugin behavior from `chenile-parent`.

The Go prototype uses a different mechanism:

1. `go.work` is the workspace-level module list.
2. Each framework component has its own `go.mod`.
3. Internal framework dependencies use stable module paths such as `core`, `base`, `http`, etc.
4. During local development, generated services use `replace` directives to point those module paths back to local workspace modules.
5. `Makefile` centralizes the test command across all modules.

## Runtime Dependency Graph

```text
http
  -> core
      -> base
      -> owiz
```

## Developer Tooling Graph

```text
test
  -> http
  -> github.com/cucumber/godog

packager
  -> core
  -> http

servicegen
  -> standard library only
```

## Generated Service Dependency Shape

Generated services receive a `go.mod` like this:

```text
require (
  core v0.0.0
  http v0.0.0
  packager v0.0.0
  test v0.0.0
)

replace core => ../../chenile-framework/core
replace base => ../../chenile-framework/base
replace http => ../../chenile-framework/http
replace owiz => ../../chenile-framework/owiz
replace packager => ../../chenile-framework/packager
replace test => ../../chenile-framework/test
```

That gives new developers a service that can be tested locally before framework modules are published.

When these modules are published, generated services can drop the local `replace` directives and use real released versions.

Generated service packages are intentionally not placed under Go `internal/` directories. A mainweb packager app must be able to import service modules, for example `customer-service/customer` and `order-service/order`.
