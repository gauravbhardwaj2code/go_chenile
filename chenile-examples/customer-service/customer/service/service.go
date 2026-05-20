package service

import (
	"context"

	"customer-service/customer/domain"
	"customer-service/customer/repository"
)

type Service interface {
	Create(context.Context, domain.CreateCustomerCommand) (domain.Customer, error)
}

type service struct {
	repository repository.Repository
}

func New(repository repository.Repository) Service {
	return service{repository: repository}
}

func (s service) Create(ctx context.Context, command domain.CreateCustomerCommand) (domain.Customer, error) {
	if command.Name == "" {
		return domain.Customer{}, domain.NameRequired()
	}
	return s.repository.Create(ctx, domain.Customer{
		Name: command.Name,
	})
}
