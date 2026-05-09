package order

import (
	"context"
	"net/http"

	"core"
)

func Register(registry *core.Registry) error {
	service := NewService()
	return registry.RegisterService(core.ServiceDefinition{
		ID:   "orderService",
		Name: "orderService",
		Operations: []core.OperationDefinition{
			{
				Name:   "create",
				Method: http.MethodPost,
				Path:   "/orders",
				NewInput: func() any {
					return &CreateOrderRequest{}
				},
				Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
					request := exchange.Input.(*CreateOrderRequest)
					return service.Create(ctx, *request)
				},
			},
		},
	})
}
