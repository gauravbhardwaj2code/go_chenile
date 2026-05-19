package stm

import (
	"context"
	"testing"
)

type order struct {
	State string
	Calls []string
}

func (o *order) GetState() string {
	return o.State
}

func (o *order) SetState(state string) {
	o.State = state
}

func TestMachineProcessesTransitionsAndActions(t *testing.T) {
	machine, err := FromJSON([]byte(`{
		"created": {"initial": true, "events": {"confirm": {"newState": "confirmed"}}},
		"confirmed": {"events": {"fulfill": {"newState": "fulfilled"}}},
		"fulfilled": {}
	}`), Actions{
		Transition: map[string]Action{
			"confirm": ActionFunc(func(ctx context.Context, transition Transition) error {
				transition.Entity.(*order).Calls = append(transition.Entity.(*order).Calls, "confirm")
				return nil
			}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	entity := &order{}

	if _, err := machine.Process(context.Background(), entity, InitialEvent, nil); err != nil {
		t.Fatal(err)
	}
	if entity.State != "created" {
		t.Fatalf("expected created, got %q", entity.State)
	}
	if _, err := machine.Process(context.Background(), entity, "confirm", nil); err != nil {
		t.Fatal(err)
	}
	if entity.State != "confirmed" || len(entity.Calls) != 1 {
		t.Fatalf("unexpected entity %#v", entity)
	}
}

func TestMachineRejectsInvalidEvent(t *testing.T) {
	machine, err := FromJSON([]byte(`{
		"created": {"initial": true}
	}`), Actions{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = machine.Process(context.Background(), &order{State: "created"}, "confirm", nil)
	if err == nil {
		t.Fatal("expected invalid event error")
	}
}
