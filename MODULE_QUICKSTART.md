# Chenile-Go Module Quick Start

## Create Your First Service Module in 5 Minutes

This guide shows you how to create a new service module from scratch and run it.

### Prerequisites

- Go 1.22+ installed
- This repository cloned locally
- Basic Go knowledge

---

## Step 1: Create Directory Structure

**Recommended:** Use the service generator instead of manual creation:

```bash
cd /workspace
go run ./chenile-framework/servicegen/cmd/chenile-servicegen new --name myservice --out ./chenile-examples
cd chenile-examples/myservice-service
```

If creating manually:

```bash
cd /workspace/chenile-examples
mkdir -p myservice/cmd/myservice
mkdir -p myservice/item
mkdir -p myservice/test/features
```

Your structure should look like:
```
myservice/
├── cmd/
│   └── myservice/
│       └── main.go
├── item/
│   ├── controller.go
│   ├── service.go
│   ├── model.go
│   └── module.go
├── test/
│   └── features/
└── go.mod
```

**Note:** No config folder is needed. Configuration is done explicitly in code for type safety and simpler deployment.

---

## Step 2: Initialize Go Module

If using the generator, this is done automatically. For manual creation:

```bash
cd /workspace/chenile-examples/myservice
go mod init myservice-service
```

Edit `go.mod` to add framework dependencies:

```go
module myservice-service

go 1.22

require (
    core v0.0.0
    http v0.0.0
    packager v0.0.0
    test v0.0.0
)

replace base => ../../chenile-framework/base
replace core => ../../chenile-framework/core
replace http => ../../chenile-framework/http
replace owiz => ../../chenile-framework/owiz
replace packager => ../../chenile-framework/packager
replace test => ../../chenile-framework/test
```

---

## Step 3: Define Models (`item/model.go`)

```go
package item

type CreateItemRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Price       float64 `json:"price"`
}

type ItemResponse struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}

type GetItemRequest struct {
    ID string `json:"id"`
}
```

---

## Step 4: Implement Service Logic (`item/service.go`)

```go
package item

import (
    "context"
    "fmt"
)

type Service struct {
    // Add repositories, clients, etc. here
}

func NewService() *Service {
    return &Service{}
}

func (s *Service) Create(ctx context.Context, req CreateItemRequest) (*ItemResponse, error) {
    // Business logic here
    // In real app: validate, save to DB, publish events, etc.
    
    return &ItemResponse{
        ID:          fmt.Sprintf("item-%d", 1), // Replace with actual ID generation
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
    }, nil
}

func (s *Service) Get(ctx context.Context, id string) (*ItemResponse, error) {
    // Business logic: fetch from database
    
    // For demo, return a mock item
    return &ItemResponse{
        ID:          id,
        Name:        "Sample Item",
        Description: "A sample item description",
        Price:       99.99,
    }, nil
}
```

---

## Step 5: Register Operations (`item/controller.go`)

```go
package item

import (
    "context"
    "net/http"

    "core"
)

func Register(registry *core.Registry) error {
    service := NewService()
    
    return registry.RegisterService(core.ServiceDefinition{
        ID:   "itemService",
        Name: "itemService",
        Operations: []core.OperationDefinition{
            {
                Name:   "create",
                Method: http.MethodPost,
                Path:   "/items",
                NewInput: func() any {
                    return &CreateItemRequest{}
                },
                Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
                    request := exchange.Input.(*CreateItemRequest)
                    return service.Create(ctx, *request)
                },
            },
            {
                Name:   "get",
                Method: http.MethodGet,
                Path:   "/items/{id}",
                NewInput: func() any {
                    return &GetItemRequest{}
                },
                Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
                    // Path params are extracted manually for now
                    // In production, add middleware to parse path templates
                    id := exchange.PathParams["id"]
                    if id == "" {
                        // Extract from URL path as fallback
                        // Simple split for demo - use proper router in production
                        return nil, fmt.Errorf("id parameter required")
                    }
                    return service.Get(ctx, id)
                },
            },
        },
    })
}
```

---

## Step 6: Create Main Entry Point (`cmd/myservice/main.go`)

```go
package main

import (
    "log"
    "os"

    "packager"
    "myservice-service/item"
)

func main() {
    // Allow port configuration via environment variable
    port := os.Getenv("PORT")
    if port == "" {
        port = ":8080"
    }
    
    app, err := packager.NewWebApp(
        packager.Module{
            Name:     "item",
            Register: item.Register,
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("listening on %s", port)
    log.Fatal(app.ListenAndServe(port))
}
```

---

## Step 7: Run Your Service

```bash
cd /workspace/chenile-examples/myservice-service
go run ./cmd/myservice
```

You should see:
```
2024/01/01 12:00:00 listening on :8080
```

**Note:** Configuration is done in code, not YAML files. To change settings:
- Port: Use `PORT` environment variable
- Other config: Add to your service code or use environment variables

---

## Step 8: Test Your API

Open another terminal and test the endpoints:

### Create an Item
```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 1299.99
  }'
```

Expected response:
```json
{
  "code": 200,
  "success": true,
  "payload": {
    "id": "item-1",
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 1299.99
  }
}
```

### Get an Item
```bash
curl http://localhost:8080/items/item-123
```

Expected response:
```json
{
  "code": 200,
  "success": true,
  "payload": {
    "id": "item-123",
    "name": "Sample Item",
    "description": "A sample item description",
    "price": 99.99
  }
}
```

---

## Step 9: Write Tests

Create `test/item_service_test.go`:

```go
package test

import (
    "io"
    "testing"

    "packager"
    godogtest "test/godog"
    "myservice-service/item"
)

func TestItemService(t *testing.T) {
    app, err := packager.NewWebApp(
        packager.Module{Name: "item", Register: item.Register},
    )
    if err != nil {
        t.Fatal(err)
    }

    status := godogtest.Suite{
        Name:         "item-service",
        Router:       app.Router,
        FeaturePaths: []string{"features/item.feature"},
        TestingT:     t,
        Output:       io.Discard,
    }.Run()
    if status != 0 {
        t.Fatalf("godog suite failed with status %d", status)
    }
}
```

Create `test/features/item.feature`:

```gherkin
Feature: Item Service
  As a user
  I want to manage items
  So that I can track inventory

  Scenario: Create a new item
    Given I have an item creation request
    When I send a POST request to /items
    Then I should receive a successful response
    And the response should contain the item details

  Scenario: Get an item by ID
    Given an item exists with ID "test-123"
    When I send a GET request to /items/test-123
    Then I should receive the item details
```

Run tests:
```bash
cd /workspace/chenile-examples/myservice-service
go test ./test/... -v
```

---

## Step 10: Combining Multiple Modules

To create an application with multiple service modules:

### Create `chenile-examples/myapp/cmd/myapp/main.go`

```go
package main

import (
    "log"

    "packager"
    "myservice-service/item"
    "customer-service/customer"  // Import other modules
)

func main() {
    app, err := packager.NewWebApp(
        packager.Module{
            Name:     "item",
            Register: item.Register,
        },
        packager.Module{
            Name:     "customer",
            Register: customer.Register,
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("listening on :8080")
    log.Fatal(app.ListenAndServe(":8080"))
}
```

**Important:** Update your `go.mod` to include all service modules:

```go
module myapp

go 1.22

require (
    myservice-service v0.0.0
    customer-service v0.0.0
    packager v0.0.0
)

replace myservice-service => ../myservice-service
replace customer-service => ../customer-service
replace base => ../../chenile-framework/base
replace core => ../../chenile-framework/core
replace http => ../../chenile-framework/http
replace owiz => ../../chenile-framework/owiz
replace packager => ../../chenile-framework/packager
replace test => ../../chenile-framework/test
```

Now both `/items` and `/customers` routes are served from the same application!

See `chenile-examples/mainweb-app` for a working example combining customer and order services.

---

## Common Patterns

### Adding Validation

```go
func (s *Service) Create(ctx context.Context, req CreateItemRequest) (*ItemResponse, error) {
    if req.Name == "" {
        return nil, fmt.Errorf("name is required")
    }
    if req.Price < 0 {
        return nil, fmt.Errorf("price must be positive")
    }
    // ... rest of logic
}
```

### Using Interceptors

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

### Error Handling

```go
import chenileerrors "base/errors"

func (s *Service) Get(ctx context.Context, id string) (*ItemResponse, error) {
    if id == "" {
        return nil, chenileerrors.NewNotFoundError("item not found")
    }
    // ...
}
```

---

## Next Steps

1. **Add Database Integration**: Use GORM, sqlx, or your preferred ORM
2. **Add Authentication**: Implement JWT or OAuth2 interceptors
3. **Add Metrics**: Integrate Prometheus or OpenTelemetry
4. **Add Documentation**: Generate OpenAPI specs from operation definitions
5. **Deploy**: Containerize with Docker and deploy to Kubernetes

---

## Troubleshooting

### "module not found"
Ensure replace directives point to correct paths relative to your service directory.

### "port already in use"
Set a different port: `PORT=:8081 go run ./cmd/myservice`

### "operation not found"
Check that the HTTP method and path exactly match what's registered in controller.go.

### JSON unmarshaling errors
Verify your request struct has correct JSON tags matching the incoming payload.

---

## See Also

- `ARCHITECTURE_GUIDE.md` - Detailed framework architecture
- `chenile-examples/customer-service/` - Complete working example
- `chenile-examples/mainweb-app/` - Multi-service application example
- `README.md` - Project overview
- `AGENT_GUIDE.md` - Guide for AI agents working with this framework
