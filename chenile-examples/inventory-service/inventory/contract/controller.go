package contract

import (
	"context"

	"chenile"

	"inventory-service/inventory/domain"
	servicepkg "inventory-service/inventory/service"
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
		chenile.POST("/inventorys", "create", func() *CreateInventoryRequest {
			return &CreateInventoryRequest{}
		}, c.Create),
	}
}

func (c Controller) Create(ctx context.Context, request CreateInventoryRequest) (Inventory, error) {
	entity, err := c.service.Create(ctx, domain.CreateInventoryCommand{Name: request.Name})
	if err != nil {
		return Inventory{}, err
	}
	return NewInventoryResponse(entity), nil
}
