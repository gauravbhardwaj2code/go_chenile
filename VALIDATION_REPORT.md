# Chenile-Go Framework Validation Report

## Summary

Status: validated on 2026-05-19.

The framework modules and example services build and pass the repository test command:

```bash
make test
make coverage
```

## Commands Run

```bash
make test
make coverage
```

The command runs:

```bash
cd chenile-framework && go test ./base/... ./owiz/... ./core/... ./http/... ./bdd-utils/... ./config/... ./middleware/... ./stm/... ./stateentity/... ./servicegen/... ./packager/...
cd chenile-examples/customer-service && go test ./...
cd chenile-examples/order-service && go test ./...
cd chenile-examples/state-order-service && go test ./...
cd chenile-examples/mainweb-app && go test ./...
```

The coverage command runs the same module groups with `go test -cover`.

## Results

Framework modules:

- `base/errors`: pass
- `base/response`: pass
- `owiz/chain`: pass
- `core`: pass
- `http`: pass
- `bdd-utils/godog`: pass
- `config`: pass
- `middleware`: pass
- `stm`: pass
- `stateentity`: pass
- `servicegen/cmd/chenile-servicegen`: pass
- `packager`: pass
- `packager/cmd/chenile-packager`: pass

Example modules:

- `chenile-examples/customer-service`: pass
- `chenile-examples/order-service`: pass
- `chenile-examples/state-order-service`: pass
- `chenile-examples/mainweb-app`: pass

## Verified Architecture

Dependency direction remains layered:

```text
application services
  -> packager
  -> http
  -> core
  -> base
  -> owiz
```

The examples demonstrate both standalone services and a combined app:

- `chenile-examples/customer-service`
- `chenile-examples/order-service`
- `chenile-examples/mainweb-app`

## Swagger/OpenAPI

Every router exposes:

- `GET /openapi.json`: OpenAPI 3.0.3 document generated from registered operations
- `GET /swagger`: Swagger UI pointing at `/openapi.json`

The OpenAPI document includes operation IDs, tags, JSON request body schemas derived from `NewInput()`, and the standardized generic JSON response shape.

## Current Constraints

- The HTTP router matches exact method/path pairs. Route templates such as `/items/{id}` and automatic `PathParams` extraction are not implemented.
- Configuration is explicit Go code. The generator does not create YAML config files or config directories.
- Generated services use local `replace` directives that point to `../../chenile-framework/<module>` when generated under `chenile-examples`.

## Fixes Applied During Validation

- Updated `chenile-framework/go.work` to `go 1.22` so it matches the framework modules.
- Updated `Makefile` so example modules are tested from their own module directories.
- Updated `Makefile` coverage so it also runs from each module directory.
- Updated service generator defaults so services generated under `chenile-examples` receive correct framework `replace` paths.
- Updated documentation to remove stale `/workspace`, YAML config, path-template, and error-helper references.
