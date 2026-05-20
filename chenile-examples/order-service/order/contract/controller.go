package contract

import (
	"context"

	"chenile"

	"order-service/order/domain"
	servicepkg "order-service/order/service"
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
		chenile.POST("/orders", "create", func() *CreateOrderRequest {
			return &CreateOrderRequest{}
		}, c.Create),
	}
}

func (c Controller) Create(ctx context.Context, request CreateOrderRequest) (Order, error) {
	entity, err := c.service.Create(ctx, domain.CreateOrderCommand{Name: request.Name})
	if err != nil {
		return Order{}, err
	}
	return NewOrderResponse(entity), nil
}
