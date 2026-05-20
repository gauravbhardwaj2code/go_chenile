package contract

import "slot-service/slot/domain"

type AddRunnerRequest struct {
	Name       string            `json:"name"`
	Skills     []domain.Skill    `json:"skills"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Slots      []domain.TimeSlot `json:"slots"`
}

type AllocateRequest struct {
	RequestID   string              `json:"requestId"`
	Skill       domain.Skill        `json:"skill"`
	Slot        domain.TimeSlot     `json:"slot"`
	Constraints []domain.Constraint `json:"constraints,omitempty"`
}
