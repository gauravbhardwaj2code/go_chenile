package module

import (
	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile"

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
