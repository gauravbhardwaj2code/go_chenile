package module

import (
	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile"

	"customer-service/customer/contract"
	"customer-service/customer/repository"
	"customer-service/customer/service"
)

func New() chenile.Module {
	return chenile.NewModule("customer", func(builder *chenile.Builder) error {
		repo := repository.NewMemoryRepository()
		svc := service.New(repo)
		return builder.Service("customerService").Routes(contract.Routes(svc)...)
	})
}
