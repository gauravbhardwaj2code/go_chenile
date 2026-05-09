package chain

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestChainExecutesCommandsInOrder(t *testing.T) {
	type state struct {
		steps []string
	}
	value := state{}
	c := New[state](
		CommandFunc[state](func(ctx context.Context, value *state) error {
			value.steps = append(value.steps, "first")
			return nil
		}),
		CommandFunc[state](func(ctx context.Context, value *state) error {
			value.steps = append(value.steps, "second")
			return nil
		}),
	)

	if err := c.Execute(context.Background(), &value); err != nil {
		t.Fatal(err)
	}

	expected := []string{"first", "second"}
	if !reflect.DeepEqual(value.steps, expected) {
		t.Fatalf("expected %v, got %v", expected, value.steps)
	}
}

func TestChainStopsOnError(t *testing.T) {
	type state struct {
		steps []string
	}
	expectedErr := errors.New("stop")
	value := state{}
	c := New[state](
		CommandFunc[state](func(ctx context.Context, value *state) error {
			value.steps = append(value.steps, "first")
			return expectedErr
		}),
		CommandFunc[state](func(ctx context.Context, value *state) error {
			value.steps = append(value.steps, "second")
			return nil
		}),
	)

	err := c.Execute(context.Background(), &value)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
	if len(value.steps) != 1 {
		t.Fatalf("expected one executed command, got %d", len(value.steps))
	}
}
