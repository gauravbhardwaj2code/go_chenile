package customer

import "context"

type Service interface {
	Create(context.Context, CreateCustomerRequest) (Customer, error)
}

type service struct{}

func NewService() Service {
	return service{}
}

func (service) Create(ctx context.Context, request CreateCustomerRequest) (Customer, error) {
	return Customer{
		ID:   "customer-1",
		Name: request.Name,
	}, nil
}
