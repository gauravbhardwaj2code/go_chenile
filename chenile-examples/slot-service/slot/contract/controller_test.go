package contract

import (
	"context"
	"testing"

	"slot-service/slot/domain"
)

type fakeService struct{}

func (fakeService) AddRunner(ctx context.Context, command domain.AddRunnerCommand) (domain.Runner, error) {
	return domain.Runner{ID: "runner-1", Name: command.Name, Skills: command.Skills, Attributes: command.Attributes, Slots: command.Slots}, nil
}

func (fakeService) Allocate(ctx context.Context, command domain.AllocateCommand) (domain.Allocation, error) {
	return domain.Allocation{ID: "allocation-1", RequestID: command.RequestID, RunnerID: "runner-1", RunnerName: "Asha", Skill: command.Skill, Slot: command.Slot, SoftScore: 1}, nil
}

func TestRoutesDeclareRunnerAndAllocationOperations(t *testing.T) {
	routes := Routes(fakeService{})

	if len(routes) != 2 {
		t.Fatalf("expected two routes, got %d", len(routes))
	}
	if routes[0].Name != "addRunner" || routes[0].Path != "/runners" {
		t.Fatalf("unexpected runner route %#v", routes[0])
	}
	if routes[1].Name != "allocate" || routes[1].Path != "/allocations" {
		t.Fatalf("unexpected allocation route %#v", routes[1])
	}
	if routes[0].NewInput == nil || routes[1].NewInput == nil {
		t.Fatal("expected input factories")
	}
}

func TestAddRunnerInvokesService(t *testing.T) {
	controller := NewController(fakeService{})

	response, err := controller.AddRunner(context.Background(), AddRunnerRequest{
		Name:   "Asha",
		Skills: []domain.Skill{domain.SkillCook},
		Slots:  []domain.TimeSlot{{Date: "2026-06-01", Start: "09:00", End: "11:00"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if response.Name != "Asha" {
		t.Fatalf("expected Asha, got %#v", response)
	}
}

func TestAllocateInvokesService(t *testing.T) {
	controller := NewController(fakeService{})

	response, err := controller.Allocate(context.Background(), AllocateRequest{
		RequestID: "req-1",
		Skill:     domain.SkillCook,
		Slot:      domain.TimeSlot{Date: "2026-06-01", Start: "09:00", End: "11:00"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if response.RunnerName != "Asha" {
		t.Fatalf("expected Asha, got %#v", response)
	}
}
