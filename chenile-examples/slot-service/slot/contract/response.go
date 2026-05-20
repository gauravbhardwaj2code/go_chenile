package contract

import "slot-service/slot/domain"

type RunnerResponse struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Skills     []domain.Skill    `json:"skills"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Slots      []domain.TimeSlot `json:"slots,omitempty"`
}

type AllocationResponse struct {
	ID                string          `json:"id"`
	RequestID         string          `json:"requestId"`
	RunnerID          string          `json:"runnerId"`
	RunnerName        string          `json:"runnerName"`
	Skill             domain.Skill    `json:"skill"`
	Slot              domain.TimeSlot `json:"slot"`
	SoftScore         int             `json:"softScore"`
	MatchedSoft       []string        `json:"matchedSoft,omitempty"`
	UnmatchedSoft     []string        `json:"unmatchedSoft,omitempty"`
	MatchedHard       []string        `json:"matchedHard,omitempty"`
	ConsideredRunners int             `json:"consideredRunners"`
}

func NewRunnerResponse(runner domain.Runner) RunnerResponse {
	return RunnerResponse{
		ID:         runner.ID,
		Name:       runner.Name,
		Skills:     runner.Skills,
		Attributes: runner.Attributes,
		Slots:      runner.Slots,
	}
}

func NewAllocationResponse(allocation domain.Allocation) AllocationResponse {
	return AllocationResponse{
		ID:                allocation.ID,
		RequestID:         allocation.RequestID,
		RunnerID:          allocation.RunnerID,
		RunnerName:        allocation.RunnerName,
		Skill:             allocation.Skill,
		Slot:              allocation.Slot,
		SoftScore:         allocation.SoftScore,
		MatchedSoft:       allocation.MatchedSoft,
		UnmatchedSoft:     allocation.UnmatchedSoft,
		MatchedHard:       allocation.MatchedHard,
		ConsideredRunners: allocation.ConsideredRunners,
	}
}
