package service

import (
	"context"

	"order-service/order/domain"
	"order-service/order/repository"
)

type Service interface {
	Create(context.Context, domain.CreateOrderCommand) (domain.Order, error)
}

type service struct {
	repository repository.Repository
}

func New(repository repository.Repository) Service {
	return service{repository: repository}
}

func (s service) Create(ctx context.Context, command domain.CreateOrderCommand) (domain.Order, error) {
	if command.Name == "" {
		return domain.Order{}, domain.NameRequired()
	}
	return s.repository.Create(ctx, domain.Order{
		Name: command.Name,
	})
}
