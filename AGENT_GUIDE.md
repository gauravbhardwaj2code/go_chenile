# Chenile-Go Framework - Agent Guide

This guide is for AI agents working with the Chenile-Go framework. It provides essential information about the framework structure, conventions, and how to interact with the codebase.

## Quick Reference

### Repository Structure

```
<repo-root>/
├── chenile-framework/          # Core framework modules
│   ├── base/                   # Foundation types (errors, responses)
│   ├── core/                   # Registry, Exchange, EntryPoint
│   ├── bdd-utils/              # Godog BDD testing utilities
│   ├── config/                 # Config loading
│   ├── http/                   # HTTP router and handlers
│   ├── middleware/             # Reusable interceptors
│   ├── owiz/                   # Command chain orchestration
│   ├── packager/               # Application assembler
│   ├── servicegen/             # Service skeleton generator
│   ├── stateentity/            # State entity service registration
│   └── stm/                    # State transition machine
├── chenile-examples/           # Example services
│   ├── customer-service/       # Standalone customer service
│   ├── order-service/          # Standalone order service
│   └── mainweb-app/            # Combined multi-service app
├── README.md                   # Project overview
├── ARCHITECTURE_GUIDE.md       # Detailed architecture
├── MODULE_QUICKSTART.md        # Step-by-step module creation
├── DEPENDENCY_MANAGEMENT.md    # Go module management
└── AGENT_GUIDE.md              # This file
```

### Key Commands

```bash
# Run all tests
make test

# Generate a new service
go run ./chenile-framework/servicegen/cmd/chenile-servicegen new --name <service-name> --out ./chenile-examples

# Run a standalone service
go run ./chenile-examples/<service-name>-service/cmd/<service-name>-service

# Run the combined application
go run ./chenile-examples/mainweb-app/cmd/mainweb-app

# Test a specific service
cd chenile-examples/<service-name>-service && go test ./...
```

## Framework Architecture

### Dependency Flow

```
Application (your service)
    ↓
packager (app assembly)
    ↓
http (routing)
    ↓
core (registry, exchange, entrypoint)
    ↓
base (errors, responses)
    ↓
owiz (command chains)
```

### Module Responsibilities

| Module | Purpose | Key Types |
|--------|---------|-----------|
| `base` | Foundation types | `GenericResponse`, `ChenileError` |
| `owiz` | Orchestration | `Chain[T]`, `Command[T]` |
| `core` | Framework kernel | `Registry`, `Exchange`, `EntryPoint`, `Interceptor` |
| `http` | Web layer | `Router` |
| `packager` | App assembler | `App`, `Module` |
| `servicegen` | Code generation | CLI tool only |
| `test` | BDD testing | `godog.Suite` |

## Service Generation

### Generator Templates

The service generator (`servicegen/cmd/chenile-servicegen/main.go`) creates these files:

1. **go.mod** - Module dependencies with local replace directives
2. **cmd/<service>/main.go** - Application entry point using packager
3. **<package>/model.go** - Request/response DTOs
4. **<package>/service.go** - Business logic interface and implementation
5. **<package>/controller.go** - Operation registration with registry
6. **<package>/module.go** - Package declaration comment
7. **<package>/service_test.go** - Unit tests for service layer
8. **<package>/controller_test.go** - Unit tests for controller
9. **test/<package>_service_test.go** - Godog BDD test runner
10. **test/features/<package>.feature** - Gherkin feature file
11. **test/fixtures/create_<package>.json** - Test fixtures

**Important:** The generator does NOT create a config folder. Configuration is done in code.

### Generated Service Pattern

```go
// 1. Model layer (<package>/model.go)
type Create<Entity>Request struct {
    Name string `json:"name"`
}

type <Entity> struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

// 2. Service layer (<package>/service.go)
type Service interface {
    Create(context.Context, Create<Entity>Request) (<Entity>, error)
}

func NewService() Service { return service{} }

// 3. Controller layer (<package>/controller.go)
func Register(registry *core.Registry) error {
    service := NewService()
    return registry.RegisterService(core.ServiceDefinition{
        ID:   "<package>Service",
        Name: "<package>Service",
        Operations: []core.OperationDefinition{{
            Name:   "create",
            Method: http.MethodPost,
            Path:   "/<entities>",
            NewInput: func() any { return &Create<Entity>Request{} },
            Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
                request := exchange.Input.(*Create<Entity>Request)
                return service.Create(ctx, *request)
            },
        }},
    })
}

// 4. Main entry point (cmd/<service>/main.go)
func main() {
    app, err := packager.NewWebApp(
        packager.Module{Name: "<package>", Register: <package>.Register},
    )
    // handle error and start server
}
```

## Testing Patterns

### Unit Tests

```go
// Service unit test
func TestServiceCreatesEntity(t *testing.T) {
    service := NewService()
    result, err := service.Create(context.Background(), CreateRequest{Name: "Test"})
    if err != nil { t.Fatal(err) }
    if result.ID == "" { t.Fatal("expected id") }
}

// Controller unit test
func TestRegisterAddsOperation(t *testing.T) {
    registry := core.NewRegistry()
    if err := Register(registry); err != nil { t.Fatal(err) }
    operation, ok := registry.Operation("entityService", "create")
    if !ok { t.Fatal("expected create operation") }
}
```

### BDD Tests (Godog)

```go
// test/<package>_service_test.go
func Test<Entity>Service(t *testing.T) {
    app, err := packager.NewWebApp(
        packager.Module{Name: "<package>", Register: <package>.Register},
    )
    
    status := godogtest.Suite{
        Name:         "<package>-service",
        Router:       app.Router,
        FeaturePaths: []string{"features/<package>.feature"},
        TestingT:     t,
        Output:       io.Discard,
    }.Run()
    
    if status != 0 {
        t.Fatalf("godog suite failed with status %d", status)
    }
}
```

```gherkin
# test/features/<package>.feature
Feature: <Entity> service

  Scenario: Create <entity>
    When I POST a REST request to URL "/<entities>" with payload
      """
      {
        "name": "Alice"
      }
      """
    Then the http status code is 200
    And success is true
    And the REST response key "name" is "Alice"
```

## Configuration Philosophy

**No YAML configuration files are used.** The framework uses:

1. **Explicit code configuration** - All settings in Go code
2. **Constructor parameters** - For dependency injection
3. **Normal Go code in `main.go`** - For runtime choices such as listen address

Example:
```go
log.Fatal(app.ListenAndServe(":8080"))
```

This approach provides:
- Type safety at compile time
- No parsing errors at runtime
- Simpler deployment (no config files to manage)
- Clear documentation through code

**Important:** Do NOT create config folders or YAML configuration files in your services. This is by design.

## Multi-Service Applications

Services can be combined into a single application:

```go
// chenile-examples/mainweb-app/app.go
func NewApp() (*packager.App, error) {
    return packager.NewWebApp(
        packager.Module{Name: "customer", Register: customer.Register},
        packager.Module{Name: "order", Register: order.Register},
    )
}
```

**Key Points:**
- Each service remains independently testable
- The packager aggregates all operations into one router
- One HTTP server handles all routes
- Services must not import each other directly (use events or shared kernel)

## Common Tasks for Agents

### Adding a New Endpoint

1. Add request/response models to `<package>/model.go`
2. Implement business logic in `<package>/service.go`
3. Register operation in `<package>/controller.go`:
   ```go
   {
       Name:   "create",
       Method: http.MethodPost,
       Path:   "/entities",
       NewInput: func() any { return &CreateRequest{} },
       Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
           request := exchange.Input.(*CreateRequest)
           return service.Create(ctx, *request)
       },
   }
   ```

The current router matches exact method/path pairs. Do not document or rely on `/entities/{id}` path templates unless route-template matching is implemented first.

### Adding Validation

```go
func (s *service) Create(ctx context.Context, req CreateRequest) (Response, error) {
    if req.Name == "" {
        return Response{}, chenileerrors.New(http.StatusBadRequest, 0, "name is required")
    }
    // ... rest of logic
}
```

### Adding Interceptors

```go
type LoggingInterceptor struct{}

func (l *LoggingInterceptor) Before(ctx context.Context, exchange *core.Exchange) error {
    log.Printf("Request: %s %s", exchange.Method, exchange.Path)
    return nil
}

func (l *LoggingInterceptor) After(ctx context.Context, exchange *core.Exchange) error {
    log.Printf("Response: %+v", exchange.Output)
    return nil
}

// In main.go:
entryPoint := core.NewEntryPoint(registry, &LoggingInterceptor{})
```

## Error Handling

Use framework error types from `base/errors`:

```go
import chenileerrors "base/errors"

return Response{}, chenileerrors.New(http.StatusNotFound, 0, "entity not found")
return Response{}, chenileerrors.New(http.StatusBadRequest, 0, "invalid input")
return Response{}, chenileerrors.New(http.StatusInternalServerError, 0, "database error")
```

## Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Service name | kebab-case | `customer-service` |
| Package name | lowercase, no hyphens | `customer` |
| Service ID | package name + "Service" | `customerService` |
| Route base | plural, kebab-case | `/customers` |
| Model structs | PascalCase | `CreateCustomerRequest`, `Customer` |
| Files | snake_case | `model.go`, `service.go`, `controller.go` |

## File Locations Reference

| File Type | Location Pattern |
|-----------|-----------------|
| Service main | `chenile-examples/<name>-service/cmd/<name>-service/main.go` |
| Models | `chenile-examples/<name>-service/<package>/model.go` |
| Service logic | `chenile-examples/<name>-service/<package>/service.go` |
| Controller | `chenile-examples/<name>-service/<package>/controller.go` |
| Unit tests | `chenile-examples/<name>-service/<package>/*_test.go` |
| BDD tests | `chenile-examples/<name>-service/test/<package>_service_test.go` |
| Features | `chenile-examples/<name>-service/test/features/<package>.feature` |
| Fixtures | `chenile-examples/<name>-service/test/fixtures/*.json` |

## Troubleshooting Checklist

When something doesn't work, check:

1. **Module not found**: Verify `replace` directives in `go.mod` point to correct paths
2. **Operation not found**: Check service ID and operation name match exactly
3. **Port in use**: Stop the other process or edit `main.go` to use a different address
4. **Path params not extracted**: The router currently matches exact paths only
5. **JSON unmarshaling**: Verify JSON tags on request struct match payload
6. **Import cycle**: Services should not import each other directly

## Agent Best Practices

### DO:
- ✅ Use the service generator for new services
- ✅ Follow the model → service → controller layering
- ✅ Write both unit tests and BDD tests
- ✅ Use framework error types
- ✅ Keep services independent (no cross-imports)
- ✅ Use environment variables for configuration
- ✅ Update go.work when adding new modules

### DON'T:
- ❌ Create config folders or YAML files
- ❌ Import one service module into another
- ❌ Access HTTP request/response directly in handlers
- ❌ Mix business logic in controllers
- ❌ Skip writing tests
- ❌ Hardcode configuration values

## Examples Reference

Study these existing examples:

1. **customer-service**: Basic CRUD service example
2. **order-service**: Another service showing consistency
3. **mainweb-app**: Multi-service composition example
4. **servicegen templates**: Source of truth for generated code

## Version Information

- **Go version**: 1.22+
- **Framework modules**: Local development with replace directives
- **Published modules**: Not yet published (use local replace directives)
- **Testing**: Godog for BDD, standard testing package for unit tests

## Contact Points

For questions about:
- **Architecture**: See `ARCHITECTURE_GUIDE.md`
- **Quick start**: See `MODULE_QUICKSTART.md`
- **Dependencies**: See `DEPENDENCY_MANAGEMENT.md`
- **General overview**: See `README.md`
