package stateentity

import (
	"context"
	"testing"

	"core"
	"stm"
)

type entity struct {
	ID    string `json:"id"`
	State string `json:"state"`
}

func (e *entity) GetState() string      { return e.State }
func (e *entity) SetState(state string) { e.State = state }

type repo struct {
	value *entity
}

func (r *repo) Create(ctx context.Context, e *entity) (*entity, error) {
	r.value = e
	return e, nil
}

func (r *repo) Get(ctx context.Context, id string) (*entity, error) {
	return r.value, nil
}

func (r *repo) Save(ctx context.Context, e *entity) (*entity, error) {
	r.value = e
	return e, nil
}

func TestRegisterAddsCreateAndApplyEventOperations(t *testing.T) {
	machine, err := stm.FromJSON([]byte(`{
		"created": {"initial": true, "events": {"confirm": {"newState": "confirmed"}}},
		"confirmed": {}
	}`), stm.Actions{})
	if err != nil {
		t.Fatal(err)
	}
	registry := core.NewRegistry()
	storage := &repo{}

	if err := Register[*entity](registry, Options[*entity]{
		ServiceID: "orderService",
		Name:      "orders",
		NewEntity: func() *entity { return &entity{} },
		Repo:      storage,
		Machine:   machine,
	}); err != nil {
		t.Fatal(err)
	}

	create, ok := registry.Operation("orderService", "create")
	if !ok {
		t.Fatal("expected create operation")
	}
	created, err := create.Handler(context.Background(), &core.Exchange{Input: &entity{ID: "1"}})
	if err != nil {
		t.Fatal(err)
	}
	if created.(*entity).State != "created" {
		t.Fatalf("expected initial state, got %#v", created)
	}

	apply, ok := registry.Operation("orderService", "applyEvent")
	if !ok {
		t.Fatal("expected applyEvent operation")
	}
	updated, err := apply.Handler(context.Background(), &core.Exchange{Input: &ApplyEventRequest{ID: "1", Event: "confirm"}})
	if err != nil {
		t.Fatal(err)
	}
	if updated.(*entity).State != "confirmed" {
		t.Fatalf("expected confirmed state, got %#v", updated)
	}
}
