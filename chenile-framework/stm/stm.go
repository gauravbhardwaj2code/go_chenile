package stm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	chenileerrors "base/errors"
)

const InitialEvent = "InitialEvent"

type Entity interface {
	GetState() string
	SetState(string)
}

type Event struct {
	NewState string `json:"newState"`
}

type State struct {
	Initial   bool             `json:"initial,omitempty"`
	Automatic bool             `json:"automatic,omitempty"`
	Events    map[string]Event `json:"events,omitempty"`
}

type Action interface {
	Process(context.Context, Transition) error
}

type ActionFunc func(context.Context, Transition) error

func (f ActionFunc) Process(ctx context.Context, transition Transition) error {
	return f(ctx, transition)
}

type AutomaticAction interface {
	Process(context.Context, Entity) (string, error)
}

type AutomaticActionFunc func(context.Context, Entity) (string, error)

func (f AutomaticActionFunc) Process(ctx context.Context, entity Entity) (string, error) {
	return f(ctx, entity)
}

type Transition struct {
	OldState string
	NewState string
	Event    string
	Param    any
	Entity   Entity
}

type Actions struct {
	PreState   map[string]Action
	PostState  map[string]Action
	Transition map[string]Action
	Automatic  map[string]AutomaticAction
}

type Machine struct {
	states       map[string]State
	initialState string
	actions      Actions
}

func FromFile(path string, actions Actions) (*Machine, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read state machine %q: %w", path, err)
	}
	return FromJSON(bytes, actions)
}

func FromJSON(bytes []byte, actions Actions) (*Machine, error) {
	states := map[string]State{}
	if err := json.Unmarshal(bytes, &states); err != nil {
		return nil, fmt.Errorf("parse state machine: %w", err)
	}
	machine := &Machine{states: states, actions: actions}
	if err := machine.Validate(); err != nil {
		return nil, err
	}
	return machine, nil
}

func (m *Machine) Validate() error {
	initialCount := 0
	for name, state := range m.states {
		if state.Initial {
			initialCount++
			m.initialState = name
		}
		for event, target := range state.Events {
			if target.NewState == "" {
				return Error(http.StatusBadRequest, "stm.event.target.required", "event target state is required", "event", event)
			}
			if _, ok := m.states[target.NewState]; !ok {
				return Error(http.StatusBadRequest, "stm.state.unknown", "event targets unknown state", "state", target.NewState)
			}
		}
	}
	if initialCount != 1 {
		return Error(http.StatusBadRequest, "stm.initial.invalid", "state machine must define exactly one initial state", "count", initialCount)
	}
	return nil
}

func (m *Machine) Process(ctx context.Context, entity Entity, event string, param any) (Entity, error) {
	transition := Transition{
		OldState: entity.GetState(),
		Event:    event,
		Param:    param,
		Entity:   entity,
	}
	if transition.OldState == "" || event == InitialEvent {
		transition.Event = InitialEvent
		transition.NewState = m.initialState
		return m.apply(ctx, transition)
	}
	state, ok := m.states[transition.OldState]
	if !ok {
		return entity, Error(http.StatusBadRequest, "stm.state.invalid", "invalid current state", "state", transition.OldState)
	}
	next, ok := state.Events[event]
	if !ok {
		return entity, Error(http.StatusBadRequest, "stm.event.invalid", "invalid event for state", "event", event)
	}
	transition.NewState = next.NewState
	return m.apply(ctx, transition)
}

func (m *Machine) apply(ctx context.Context, transition Transition) (Entity, error) {
	if action := m.actions.PostState[transition.OldState]; action != nil {
		if err := action.Process(ctx, transition); err != nil {
			return transition.Entity, err
		}
	}
	if action := m.actions.Transition[transition.Event]; action != nil {
		if err := action.Process(ctx, transition); err != nil {
			return transition.Entity, err
		}
	}
	transition.Entity.SetState(transition.NewState)
	if action := m.actions.PreState[transition.NewState]; action != nil {
		if err := action.Process(ctx, transition); err != nil {
			return transition.Entity, err
		}
	}
	return m.processAutomatic(ctx, transition.Entity, transition.Param)
}

func (m *Machine) processAutomatic(ctx context.Context, entity Entity, param any) (Entity, error) {
	state := m.states[entity.GetState()]
	if !state.Automatic {
		return entity, nil
	}
	action := m.actions.Automatic[entity.GetState()]
	if action == nil {
		return entity, Error(http.StatusBadRequest, "stm.automatic.missing", "automatic state action is not configured", "state", entity.GetState())
	}
	event, err := action.Process(ctx, entity)
	if err != nil {
		return entity, err
	}
	return m.Process(ctx, entity, event, param)
}

func (m *Machine) InitialState() string {
	return m.initialState
}

func Error(status int, key string, description string, paramName string, paramValue any) error {
	return chenileerrors.Builder().
		Status(status).
		MessageKey(key).
		Description(description).
		Param(paramName, paramValue).
		Build()
}
