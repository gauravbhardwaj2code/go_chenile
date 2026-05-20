package module

import (
	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile"

	"slot-service/slot/contract"
	"slot-service/slot/repository"
	"slot-service/slot/service"
)

func New() chenile.Module {
	return chenile.NewModule("slot", func(builder *chenile.Builder) error {
		repo := repository.NewMemoryRepository()
		svc := service.New(repo)
		return builder.Service("slotService").Routes(contract.Routes(svc)...)
	})
}
