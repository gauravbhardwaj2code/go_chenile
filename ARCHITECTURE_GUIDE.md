# Chenile-Go Framework Architecture Guide

## Overview

Chenile-Go is a modular service framework inspired by Java/Spring architecture, built for Go. This guide explains the framework layers, core components, and how to create service modules that work together in a single package.

## Framework Layers

The framework is organized into distinct layers, each with specific responsibilities:

```
┌─────────────────────────────────────────┐
│         Application Layer               │
│    (chenile-examples/customer-service)  │
│  - Service Modules (customer, order)    │
│  - Main Entry Point                     │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│         Packager Layer                  │
│    (packager) │
│  - App Assembly                         │
│  - Module Registration                  │
│  - HTTP Server Bootstrap                │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│         HTTP Layer                      │
│    (http)     │
│  - Router                               │
│  - Request/Response Handling            │
│  - Path Matching                        │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│         Core Layer                      │
│    (core)     │
│  - Registry                             │
│  - EntryPoint                           │
│  - Exchange                             │
│  - Interceptors                         │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│         Base Layer                      │
│    (base)     │
│  - Errors                               │
│  - Response Types                       │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│         O Wiz (Orchestration)           │
│    (owiz)     │
│  - Chain Execution                      │
│  - Command Pattern                      │
└─────────────────────────────────────────┘
```

## Module Descriptions

### 1. Base Layer (`base`)

**Purpose**: Foundational types and utilities used across all layers.

**Sub-packages**:
- `errors/`: Chenile-specific error types with status codes
- `response/`: Standardized response structures

**Key Types**:
```go
// GenericResponse - Standard API response format
type GenericResponse struct {
    Code         int
    Description  string
    Errors       []ResponseMessage
    Payload      any
    Success      bool
}

// ChenileError - Framework error type
type ChenileError struct {
    Status       int
    Description  string
    SubErrorCode int
}
```

**When to use**: Always imported by other layers. Never depends on other chenile-go modules.

---

### 2. O Wiz Layer (`owiz`)

**Purpose**: Orchestration utilities for command chain execution.

**Key Features**:
- Command pattern implementation
- Before/After interceptor chains
- Context-aware execution

**Key Types**:
```go
// Chain - Executable command chain
type Chain[T any] struct {
    commands []Command[T]
}

// Command - Single executable unit
type Command[T any] interface {
    Execute(ctx context.Context, input *T) error
}
```

**When to use**: When you need ordered execution of operations (interceptors, middleware).

---

### 3. Core Layer (`core`)

**Purpose**: Framework kernel - service registration, operation routing, execution engine.

**Key Components**:

#### Registry
Central registry for all services and operations.

```go
type Registry struct {
    services map[string]*RegisteredService
}

type ServiceDefinition struct {
    ID         string
    Name       string
    Operations []OperationDefinition
}

type OperationDefinition struct {
    Name      string
    Method    string  // HTTP method
    Path      string  // URL path
    NewInput  func() any
    Handler   func(ctx context.Context, exchange *Exchange) (any, error)
}
```

#### Exchange
Context carrier for request/response data flow.

```go
type Exchange struct {
    Context       context.Context
    ServiceID     string
    OperationName string
    Method        string
    Path          string
    Headers       map[string]string
    QueryParams   map[string]string
    PathParams    map[string]string
    Body          []byte
    Input         any       // Parsed request body
    Output        any       // Handler result
    Response      *response.GenericResponse
    Err           error
}
```

#### EntryPoint
Execution engine that processes exchanges through interceptors to handlers.

```go
type EntryPoint struct {
    registry     *Registry
    interceptors []Interceptor
}

func (e *EntryPoint) Execute(ctx context.Context, exchange *Exchange) response.GenericResponse
```

#### Interceptor
Middleware interface for cross-cutting concerns.

```go
type Interceptor interface {
    Before(context.Context, *Exchange) error
    After(context.Context, *Exchange) error
}
```

**When to use**: Core is used by packager and http layers. Service modules interact with Registry.

---

### 4. HTTP Layer (`http`)

**Purpose**: HTTP server implementation, routing, request/response translation.

**Key Components**:

#### Router
HTTP router that maps requests to registered operations.

```go
type Router struct {
    entryPoint *core.EntryPoint
    routes     map[string]route
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request)
```

**Features**:
- Mounts operations from Registry
- Parses HTTP requests into Exchange objects
- Writes GenericResponse as JSON

**When to use**: Automatically used by packager. Rarely used directly.

---

### 5. Packager Layer (`packager`)

**Purpose**: Application assembler - brings modules together into runnable apps.

**Key Components**:

#### Module
Unit of service registration.

```go
type Module struct {
    Name     string
    Register func(*core.Registry) error
}
```

#### App
Fully assembled application with all modules.

```go
type App struct {
    Registry *core.Registry
    Router   *chenilehttp.Router
    Modules  []Module
}

func NewWebApp(modules ...Module) (*App, error)
func (a *App) ListenAndServe(address string) error
```

**When to use**: Used in main.go to bootstrap the application.

---

## Creating a Service Module

### Step 1: Module Structure

Prefer the generator:

```bash
go run ./chenile-framework/servicegen/cmd/chenile-servicegen new --name myservice --out ./chenile-examples
```

It creates this structure:

```
myservice-service/
├── cmd/
│   └── myservice-service/
│       └── main.go
├── myservice/
│   ├── controller.go    # Registry and operations
│   ├── controller_test.go
│   ├── service.go       # Business logic
│   ├── service_test.go
│   ├── model.go         # Data structures
│   └── module.go        # Package declaration
├── go.mod
└── test/
    ├── myservice_service_test.go
    ├── features/
    └── fixtures/
```

There is no config directory; configuration is explicit Go code.

### Step 2: Define Models (`model.go`)

```go
package myservice

type CreateItemRequest struct {
    Name string `json:"name"`
}

type Item struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

### Step 3: Implement Service Logic (`service.go`)

```go
package myservice

import "context"

type Service interface {
    Create(context.Context, CreateItemRequest) (Item, error)
}

type service struct{}

func NewService() Service {
    return service{}
}

func (service) Create(ctx context.Context, req CreateItemRequest) (Item, error) {
    return Item{
        ID:   "item-1",
        Name: req.Name,
    }, nil
}
```

### Step 4: Register Operations (`controller.go`)

```go
package myservice

import (
    "context"
    "net/http"

    "core"
)

func Register(registry *core.Registry) error {
    service := NewService()
    
    return registry.RegisterService(core.ServiceDefinition{
        ID:   "myserviceService",
        Name: "myserviceService",
        Operations: []core.OperationDefinition{{
            Name:   "create",
            Method: http.MethodPost,
            Path:   "/myservices",
            NewInput: func() any {
                return &CreateItemRequest{}
            },
            Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
                request := exchange.Input.(*CreateItemRequest)
                return service.Create(ctx, *request)
            },
        }},
    })
}
```

The HTTP router currently matches exact method/path pairs. Template routes such as `/items/{id}` and automatic path parameter extraction are not implemented yet.

### Step 5: Create Main Entry Point (`cmd/my-service/main.go`)

```go
package main

import (
    "log"

    "packager"
    "myservice-service/myservice"
)

func main() {
    app, err := packager.NewWebApp(
        packager.Module{
            Name:     "myservice",
            Register: myservice.Register,
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("listening on :8080")
    log.Fatal(app.ListenAndServe(":8080"))
}
```

### Step 6: Configure Module Dependencies (`go.mod`)

```go
module myservice-service

go 1.26

toolchain go1.26.3

require (
    bdd-utils v0.0.0
    config v0.0.0
    core v0.0.0
    http v0.0.0
    packager v0.0.0
)

replace bdd-utils => ../../chenile-framework/bdd-utils
replace base => ../../chenile-framework/base
replace config => ../../chenile-framework/config
replace core => ../../chenile-framework/core
replace http => ../../chenile-framework/http
replace owiz => ../../chenile-framework/owiz
replace packager => ../../chenile-framework/packager
```

---

## How Service Modules Come Together

### Single Module Application

```go
// cmd/my-service/main.go
app, err := packager.NewWebApp(
    packager.Module{
        Name:     "customer",
        Register: customer.Register,
    },
)
```

### Multi-Module Application

Multiple service modules can be combined into a single application:

```go
// cmd/my-app/main.go
app, err := packager.NewWebApp(
    packager.Module{
        Name:     "customer",
        Register: customer.Register,
    },
    packager.Module{
        Name:     "order",
        Register: order.Register,
    },
    packager.Module{
        Name:     "inventory",
        Register: inventory.Register,
    },
)
```

**What happens internally**:

1. **Registry Creation**: `NewWebApp` creates a single `core.Registry`
2. **Module Registration**: Each module's `Register` function is called with the same registry
3. **Operation Aggregation**: All operations from all modules are stored in the registry
4. **Router Mounting**: Router reads all registered operations and creates routes
5. **Single Server**: One HTTP server handles requests for all modules

### Execution Flow

```
HTTP Request → Router → EntryPoint → Interceptors → Handler → Response
                    ↓
              Looks up operation in Registry
              (aggregated from all modules)
```

---

## Key Design Patterns

### 1. Registry Pattern
Central service discovery - all operations register themselves, enabling loose coupling.

### 2. Exchange Object
Context carrier pattern - all data flows through a single object, making it easy to add cross-cutting concerns.

### 3. Command Chain
Interceptor pattern using command chain - before/after hooks for logging, auth, transactions.

### 4. Module Isolation
Each module is self-contained with its own models, services, and controllers, but shares the same registry.

---

## Best Practices

### DO:
- ✅ Keep modules independent - no direct imports between service modules
- ✅ Use the Exchange object for all handler signatures
- ✅ Return standardized errors using base/errors
- ✅ Define clear operation boundaries in controller.go
- ✅ Use interceptors for cross-cutting concerns (logging, auth)

### DON'T:
- ❌ Import one service module into another (use events or shared kernel instead)
- ❌ Access HTTP request/response directly in handlers
- ❌ Bypass the registry for service discovery
- ❌ Mix business logic in controllers

---

## Example: Complete Working Module

See `chenile-examples/customer-service/` for a complete working example including:
- Models (`customer/model.go`)
- Service layer (`customer/service.go`)
- Controller/Registry (`customer/controller.go`)
- Tests (`customer/controller_test.go`, `test/features/`)
- Main entry point (`cmd/customer-service/main.go`)

---

## Development Workflow

1. **Generate module**: `go run ./chenile-framework/servicegen/cmd/chenile-servicegen new --name myservice --out ./chenile-examples`
2. **Run generated tests**: `cd chenile-examples/myservice-service && go test ./...`
3. **Adjust dependencies**: Edit go.mod if external packages are needed
4. **Define models**: Create request/response structs
5. **Implement service**: Business logic without framework dependencies
6. **Create controller**: Register operations with registry
7. **Write main.go**: Assemble modules with packager
8. **Run**: `go run ./cmd/myservice-service`
9. **Test**: Write godog BDD tests in `test/features/`

---

## Troubleshooting

### Issue: "module not found"
**Solution**: Ensure replace directives in go.mod point to correct relative paths.

### Issue: "operation not found"
**Solution**: Verify the service ID and operation name match exactly in registry and request.

### Issue: Port already in use
**Solution**: Stop the other process or edit the service `main.go` to listen on a different address.

### Issue: Path parameters not extracted
**Solution**: The current router uses exact path matching only. Register exact paths, or implement route-template matching before using `exchange.PathParams`.

---

## Future Enhancements

- [ ] Path parameter extraction from route templates
- [ ] Built-in validation framework
- [ ] Async operation support
- [ ] Distributed tracing integration
- [ ] Configuration-driven service registration
- [ ] Hot-reload for development

---

## Summary

The Chenile-Go framework provides a clean, modular architecture for building microservices:

- **Base**: Foundation types (errors, responses)
- **O Wiz**: Orchestration (command chains)
- **Core**: Kernel (registry, exchange, entrypoint)
- **HTTP**: Web layer (router, request handling)
- **Packager**: Application assembler
- **Service Modules**: Your business logic

Modules register their operations with a shared registry, and the packager assembles them into a single runnable application. This architecture enables clean separation of concerns while maintaining simplicity.
