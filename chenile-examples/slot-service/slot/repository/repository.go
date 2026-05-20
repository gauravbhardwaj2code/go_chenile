package repository

import (
	"context"

	"slot-service/slot/domain"
)

type Repository interface {
	AddRunner(context.Context, domain.Runner) (domain.Runner, error)
	ListRunners(context.Context) ([]domain.Runner, error)
	TrySaveAllocation(context.Context, domain.Allocation) (domain.Allocation, bool, error)
	ListAllocations(context.Context) ([]domain.Allocation, error)
}
