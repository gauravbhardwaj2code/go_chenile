package service

import (
	"context"
	"sort"

	"slot-service/slot/domain"
	"slot-service/slot/repository"
)

type Service interface {
	AddRunner(context.Context, domain.AddRunnerCommand) (domain.Runner, error)
	Allocate(context.Context, domain.AllocateCommand) (domain.Allocation, error)
}

type service struct {
	repository repository.Repository
}

func New(repository repository.Repository) Service {
	return service{repository: repository}
}

func (s service) AddRunner(ctx context.Context, command domain.AddRunnerCommand) (domain.Runner, error) {
	if command.Name == "" {
		return domain.Runner{}, domain.RunnerNameRequired()
	}
	if len(command.Skills) == 0 {
		return domain.Runner{}, domain.SkillRequired()
	}
	if len(command.Slots) == 0 {
		return domain.Runner{}, domain.SlotRequired()
	}
	for _, slot := range command.Slots {
		if !validSlot(slot) {
			return domain.Runner{}, domain.SlotRequired()
		}
	}
	return s.repository.AddRunner(ctx, domain.Runner{
		Name:       command.Name,
		Skills:     command.Skills,
		Attributes: command.Attributes,
		Slots:      command.Slots,
	})
}

func (s service) Allocate(ctx context.Context, command domain.AllocateCommand) (domain.Allocation, error) {
	if command.Skill == "" {
		return domain.Allocation{}, domain.SkillRequired()
	}
	if !validSlot(command.Slot) {
		return domain.Allocation{}, domain.SlotRequired()
	}
	runners, err := s.repository.ListRunners(ctx)
	if err != nil {
		return domain.Allocation{}, err
	}
	allocations, err := s.repository.ListAllocations(ctx)
	if err != nil {
		return domain.Allocation{}, err
	}
	candidates := candidatesFor(command, runners, bookedSlots(allocations))
	if len(candidates) == 0 {
		return domain.Allocation{}, domain.NoRunnerMatched()
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].allocation.SoftScore == candidates[j].allocation.SoftScore {
			return candidates[i].allocation.RunnerID < candidates[j].allocation.RunnerID
		}
		return candidates[i].allocation.SoftScore > candidates[j].allocation.SoftScore
	})
	for _, candidate := range candidates {
		allocation, saved, err := s.repository.TrySaveAllocation(ctx, candidate.allocation)
		if err != nil {
			return domain.Allocation{}, err
		}
		if saved {
			return allocation, nil
		}
	}
	return domain.Allocation{}, domain.NoRunnerMatched()
}

type candidate struct {
	allocation domain.Allocation
}

type bookedSlotKey struct {
	runnerID string
	slot     domain.TimeSlot
}

func candidatesFor(command domain.AllocateCommand, runners []domain.Runner, booked map[bookedSlotKey]bool) []candidate {
	candidates := []candidate{}
	for _, runner := range runners {
		if !supportsSkill(runner, command.Skill) || !hasSlot(runner, command.Slot) || booked[bookedSlotKey{runnerID: runner.ID, slot: command.Slot}] {
			continue
		}
		score, matchedHard, matchedSoft, unmatchedSoft, ok := evaluateConstraints(runner, command.Constraints)
		if !ok {
			continue
		}
		candidates = append(candidates, candidate{allocation: domain.Allocation{
			RequestID:         command.RequestID,
			RunnerID:          runner.ID,
			RunnerName:        runner.Name,
			Skill:             command.Skill,
			Slot:              command.Slot,
			SoftScore:         score,
			MatchedHard:       matchedHard,
			MatchedSoft:       matchedSoft,
			UnmatchedSoft:     unmatchedSoft,
			ConsideredRunners: len(runners),
		}})
	}
	return candidates
}

func evaluateConstraints(runner domain.Runner, constraints []domain.Constraint) (int, []string, []string, []string, bool) {
	score := 0
	matchedHard := []string{}
	matchedSoft := []string{}
	unmatchedSoft := []string{}
	for _, constraint := range constraints {
		key := constraint.Key + "=" + constraint.Value
		matches := runner.Attributes[constraint.Key] == constraint.Value
		switch constraint.Type {
		case domain.Soft:
			if matches {
				score++
				matchedSoft = append(matchedSoft, key)
			} else {
				unmatchedSoft = append(unmatchedSoft, key)
			}
		default:
			if !matches {
				return 0, nil, nil, nil, false
			}
			matchedHard = append(matchedHard, key)
		}
	}
	return score, matchedHard, matchedSoft, unmatchedSoft, true
}

func supportsSkill(runner domain.Runner, skill domain.Skill) bool {
	for _, candidate := range runner.Skills {
		if candidate == skill {
			return true
		}
	}
	return false
}

func hasSlot(runner domain.Runner, slot domain.TimeSlot) bool {
	for _, candidate := range runner.Slots {
		if candidate == slot {
			return true
		}
	}
	return false
}

func bookedSlots(allocations []domain.Allocation) map[bookedSlotKey]bool {
	booked := map[bookedSlotKey]bool{}
	for _, allocation := range allocations {
		booked[bookedSlotKey{runnerID: allocation.RunnerID, slot: allocation.Slot}] = true
	}
	return booked
}

func validSlot(slot domain.TimeSlot) bool {
	return slot.Date != "" && slot.Start != "" && slot.End != ""
}
