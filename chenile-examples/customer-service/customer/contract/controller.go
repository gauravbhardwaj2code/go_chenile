package contract

import (
	"context"

	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile"

	"customer-service/customer/domain"
	servicepkg "customer-service/customer/service"
)

type Controller struct {
	service servicepkg.Service
}

func NewController(service servicepkg.Service) Controller {
	return Controller{service: service}
}

func Routes(service servicepkg.Service) []chenile.Route {
	return NewController(service).Routes()
}

func (c Controller) Routes() []chenile.Route {
	return []chenile.Route{
		chenile.POST("/customers", "create", func() *CreateCustomerRequest {
			return &CreateCustomerRequest{}
		}, c.Create),
	}
}

func (c Controller) Create(ctx context.Context, request CreateCustomerRequest) (Customer, error) {
	entity, err := c.service.Create(ctx, domain.CreateCustomerCommand{Name: request.Name})
	if err != nil {
		return Customer{}, err
	}
	return NewCustomerResponse(entity), nil
}
