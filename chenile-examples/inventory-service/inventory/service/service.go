package service

import (
	"context"

	"inventory-service/inventory/domain"
	"inventory-service/inventory/repository"
)

type Service interface {
	Create(context.Context, domain.CreateInventoryCommand) (domain.Inventory, error)
}

type service struct {
	repository repository.Repository
}

func New(repository repository.Repository) Service {
	return service{repository: repository}
}

func (s service) Create(ctx context.Context, command domain.CreateInventoryCommand) (domain.Inventory, error) {
	if command.Name == "" {
		return domain.Inventory{}, domain.NameRequired()
	}
	return s.repository.Create(ctx, domain.Inventory{
		Name: command.Name,
	})
}
