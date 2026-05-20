package module

import (
	"chenile"

	"inventory-service/inventory/contract"
	"inventory-service/inventory/repository"
	"inventory-service/inventory/service"
)

func New() chenile.Module {
	return chenile.NewModule("inventory", func(builder *chenile.Builder) error {
		repo := repository.NewMemoryRepository()
		svc := service.New(repo)
		return builder.Service("inventoryService").Routes(contract.Routes(svc)...)
	})
}
