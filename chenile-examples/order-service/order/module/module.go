package module

import (
	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile"

	"order-service/order/contract"
	"order-service/order/repository"
	"order-service/order/service"
)

func New() chenile.Module {
	return chenile.NewModule("order", func(builder *chenile.Builder) error {
		repo := repository.NewMemoryRepository()
		svc := service.New(repo)
		return builder.Service("orderService").Routes(contract.Routes(svc)...)
	})
}
