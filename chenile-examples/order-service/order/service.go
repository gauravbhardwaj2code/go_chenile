package order

import "context"

type Service interface {
	Create(context.Context, CreateOrderRequest) (Order, error)
}

type service struct{}

func NewService() Service {
	return service{}
}

func (service) Create(ctx context.Context, request CreateOrderRequest) (Order, error) {
	return Order{
		ID:   "order-1",
		Name: request.Name,
	}, nil
}
