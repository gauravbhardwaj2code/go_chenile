package contract

import (
	"context"

	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile"

	"slot-service/slot/domain"
	servicepkg "slot-service/slot/service"
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
		chenile.POST("/runners", "addRunner", func() *AddRunnerRequest {
			return &AddRunnerRequest{}
		}, c.AddRunner),
		chenile.POST("/allocations", "allocate", func() *AllocateRequest {
			return &AllocateRequest{}
		}, c.Allocate),
	}
}

func (c Controller) AddRunner(ctx context.Context, request AddRunnerRequest) (RunnerResponse, error) {
	runner, err := c.service.AddRunner(ctx, domain.AddRunnerCommand{
		Name:       request.Name,
		Skills:     request.Skills,
		Attributes: request.Attributes,
		Slots:      request.Slots,
	})
	if err != nil {
		return RunnerResponse{}, err
	}
	return NewRunnerResponse(runner), nil
}

func (c Controller) Allocate(ctx context.Context, request AllocateRequest) (AllocationResponse, error) {
	allocation, err := c.service.Allocate(ctx, domain.AllocateCommand{
		RequestID:   request.RequestID,
		Skill:       request.Skill,
		Slot:        request.Slot,
		Constraints: request.Constraints,
	})
	if err != nil {
		return AllocationResponse{}, err
	}
	return NewAllocationResponse(allocation), nil
}
