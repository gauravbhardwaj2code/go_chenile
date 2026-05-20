package domain

import (
	"net/http"

	chenileerrors "github.com/gauravbhardwaj2code/go_chenile/chenile-framework/base/errors"
)

const (
	ErrorRunnerNameRequired = 2001
	ErrorSkillRequired      = 2002
	ErrorSlotRequired       = 2003
	ErrorNoRunnerMatched    = 2004
	ErrorSlotAlreadyBooked  = 2005
)

func RunnerNameRequired() error {
	return validationError(ErrorRunnerNameRequired, "slot.runner.name.required", "runner name is required", "name")
}

func SkillRequired() error {
	return validationError(ErrorSkillRequired, "slot.skill.required", "skill is required", "skill")
}

func SlotRequired() error {
	return validationError(ErrorSlotRequired, "slot.required", "slot date, start, and end are required", "slot")
}

func NoRunnerMatched() error {
	return chenileerrors.Builder().
		Status(http.StatusNotFound).
		Code(ErrorNoRunnerMatched).
		MessageKey("slot.allocation.no_runner_matched").
		Description("no available runner satisfies the hard constraints").
		Build()
}

func SlotAlreadyBooked() error {
	return chenileerrors.Builder().
		Status(http.StatusConflict).
		Code(ErrorSlotAlreadyBooked).
		MessageKey("slot.allocation.slot_already_booked").
		Description("runner slot is already booked").
		Build()
}

func validationError(code int, key string, description string, field string) error {
	return chenileerrors.Builder().
		Status(http.StatusBadRequest).
		Code(code).
		MessageKey(key).
		Description(description).
		Field(field, "required", description).
		Build()
}
