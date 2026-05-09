package customer

import (
	"context"
	"net/http"

	"github.com/ajapro/chenile-go/core"
)

func Register(registry *core.Registry) error {
	service := NewService()
	return registry.RegisterService(core.ServiceDefinition{
		ID:   "customerService",
		Name: "customerService",
		Operations: []core.OperationDefinition{
			{
				Name:   "create",
				Method: http.MethodPost,
				Path:   "/customers",
				NewInput: func() any {
					return &CreateCustomerRequest{}
				},
				Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
					request := exchange.Input.(*CreateCustomerRequest)
					return service.Create(ctx, *request)
				},
			},
		},
	})
}
