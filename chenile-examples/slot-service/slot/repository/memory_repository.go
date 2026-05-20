package repository

import (
	"context"
	"fmt"
	"sync"

	"slot-service/slot/domain"
)

type MemoryRepository struct {
	mu               sync.Mutex
	nextRunnerID     int
	nextAllocationID int
	runners          map[string]domain.Runner
	allocations      map[string]domain.Allocation
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextRunnerID:     1,
		nextAllocationID: 1,
		runners:          map[string]domain.Runner{},
		allocations:      map[string]domain.Allocation{},
	}
}

func (r *MemoryRepository) AddRunner(ctx context.Context, runner domain.Runner) (domain.Runner, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if runner.ID == "" {
		runner.ID = fmt.Sprintf("runner-%d", r.nextRunnerID)
		r.nextRunnerID++
	}
	stored := cloneRunner(runner)
	r.runners[stored.ID] = stored
	return cloneRunner(stored), nil
}

func (r *MemoryRepository) ListRunners(ctx context.Context) ([]domain.Runner, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	runners := make([]domain.Runner, 0, len(r.runners))
	for _, runner := range r.runners {
		runners = append(runners, cloneRunner(runner))
	}
	return runners, nil
}

func (r *MemoryRepository) TrySaveAllocation(ctx context.Context, allocation domain.Allocation) (domain.Allocation, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.allocations {
		if existing.RunnerID == allocation.RunnerID && existing.Slot == allocation.Slot {
			return domain.Allocation{}, false, nil
		}
	}
	if allocation.ID == "" {
		allocation.ID = fmt.Sprintf("allocation-%d", r.nextAllocationID)
		r.nextAllocationID++
	}
	stored := cloneAllocation(allocation)
	r.allocations[stored.ID] = stored
	return cloneAllocation(stored), true, nil
}

func (r *MemoryRepository) ListAllocations(ctx context.Context) ([]domain.Allocation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	allocations := make([]domain.Allocation, 0, len(r.allocations))
	for _, allocation := range r.allocations {
		allocations = append(allocations, cloneAllocation(allocation))
	}
	return allocations, nil
}

func cloneRunner(runner domain.Runner) domain.Runner {
	runner.Skills = append([]domain.Skill(nil), runner.Skills...)
	runner.Slots = append([]domain.TimeSlot(nil), runner.Slots...)
	if runner.Attributes != nil {
		attributes := make(map[string]string, len(runner.Attributes))
		for key, value := range runner.Attributes {
			attributes[key] = value
		}
		runner.Attributes = attributes
	}
	return runner
}

func cloneAllocation(allocation domain.Allocation) domain.Allocation {
	allocation.MatchedHard = append([]string(nil), allocation.MatchedHard...)
	allocation.MatchedSoft = append([]string(nil), allocation.MatchedSoft...)
	allocation.UnmatchedSoft = append([]string(nil), allocation.UnmatchedSoft...)
	return allocation
}
