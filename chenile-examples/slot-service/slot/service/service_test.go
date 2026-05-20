package service

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"slot-service/slot/domain"
	"slot-service/slot/repository"
)

func TestAddRunnerCreatesAvailableRunner(t *testing.T) {
	service := New(repository.NewMemoryRepository())
	slot := domain.TimeSlot{Date: "2026-06-01", Start: "09:00", End: "11:00"}

	runner, err := service.AddRunner(context.Background(), domain.AddRunnerCommand{
		Name:   "Asha",
		Skills: []domain.Skill{domain.SkillCook},
		Slots:  []domain.TimeSlot{slot},
	})
	if err != nil {
		t.Fatal(err)
	}
	if runner.ID == "" {
		t.Fatal("expected generated id")
	}
	if runner.Name != "Asha" {
		t.Fatalf("expected Asha, got %q", runner.Name)
	}
}

func TestAllocateSatisfiesHardAndOptimizesSoftConstraints(t *testing.T) {
	service := New(repository.NewMemoryRepository())
	slot := domain.TimeSlot{Date: "2026-06-01", Start: "09:00", End: "11:00"}
	addRunner(t, service, "Asha", []domain.Skill{domain.SkillCook}, map[string]string{"diet": "veg", "gender": "female"}, slot)
	addRunner(t, service, "Ravi", []domain.Skill{domain.SkillCook}, map[string]string{"diet": "veg", "gender": "male"}, slot)

	allocation, err := service.Allocate(context.Background(), domain.AllocateCommand{
		RequestID: "req-1",
		Skill:     domain.SkillCook,
		Slot:      slot,
		Constraints: []domain.Constraint{
			{Key: "diet", Value: "veg", Type: domain.Hard},
			{Key: "gender", Value: "female", Type: domain.Soft},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if allocation.RunnerName != "Asha" {
		t.Fatalf("expected Asha, got %#v", allocation)
	}
	if allocation.SoftScore != 1 {
		t.Fatalf("expected soft score 1, got %d", allocation.SoftScore)
	}
}

func TestAllocateRejectsWhenHardConstraintCannotBeSatisfied(t *testing.T) {
	service := New(repository.NewMemoryRepository())
	slot := domain.TimeSlot{Date: "2026-06-01", Start: "09:00", End: "11:00"}
	addRunner(t, service, "Ravi", []domain.Skill{domain.SkillCook}, map[string]string{"diet": "non-veg"}, slot)

	_, err := service.Allocate(context.Background(), domain.AllocateCommand{
		RequestID: "req-1",
		Skill:     domain.SkillCook,
		Slot:      slot,
		Constraints: []domain.Constraint{
			{Key: "diet", Value: "veg", Type: domain.Hard},
		},
	})
	if err == nil {
		t.Fatal("expected no runner matched error")
	}
}

func TestAllocatePrefersHigherSoftScore(t *testing.T) {
	service := New(repository.NewMemoryRepository())
	slot := domain.TimeSlot{Date: "2026-06-01", Start: "09:00", End: "11:00"}
	addRunner(t, service, "Ravi", []domain.Skill{domain.SkillCook}, map[string]string{"diet": "veg", "gender": "male"}, slot)
	addRunner(t, service, "Asha", []domain.Skill{domain.SkillCook}, map[string]string{"diet": "veg", "gender": "female"}, slot)

	allocation, err := service.Allocate(context.Background(), domain.AllocateCommand{
		RequestID: "req-1",
		Skill:     domain.SkillCook,
		Slot:      slot,
		Constraints: []domain.Constraint{
			{Key: "diet", Value: "veg", Type: domain.Soft},
			{Key: "gender", Value: "female", Type: domain.Soft},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if allocation.RunnerName != "Asha" {
		t.Fatalf("expected Asha, got %#v", allocation)
	}
	if allocation.SoftScore != 2 {
		t.Fatalf("expected soft score 2, got %d", allocation.SoftScore)
	}
}

func TestAllocateDoesNotDoubleBookRunnerSlot(t *testing.T) {
	service := New(repository.NewMemoryRepository())
	slot := domain.TimeSlot{Date: "2026-06-01", Start: "09:00", End: "11:00"}
	addRunner(t, service, "Asha", []domain.Skill{domain.SkillCook}, map[string]string{"diet": "veg"}, slot)

	command := domain.AllocateCommand{RequestID: "req-1", Skill: domain.SkillCook, Slot: slot}
	if _, err := service.Allocate(context.Background(), command); err != nil {
		t.Fatal(err)
	}
	command.RequestID = "req-2"
	_, err := service.Allocate(context.Background(), command)
	if err == nil {
		t.Fatal("expected second allocation to fail")
	}
}

func TestAllocateTreatsDelimitedSlotFieldsAsDistinct(t *testing.T) {
	service := New(repository.NewMemoryRepository())
	slotOne := domain.TimeSlot{Date: "2026|06|01", Start: "09", End: "00|11:00"}
	slotTwo := domain.TimeSlot{Date: "2026|06", Start: "01|09", End: "00|11:00"}
	_, err := service.AddRunner(context.Background(), domain.AddRunnerCommand{
		Name:   "Asha",
		Skills: []domain.Skill{domain.SkillCook},
		Slots:  []domain.TimeSlot{slotOne, slotTwo},
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := service.Allocate(context.Background(), domain.AllocateCommand{RequestID: "req-1", Skill: domain.SkillCook, Slot: slotOne}); err != nil {
		t.Fatal(err)
	}
	allocation, err := service.Allocate(context.Background(), domain.AllocateCommand{RequestID: "req-2", Skill: domain.SkillCook, Slot: slotTwo})
	if err != nil {
		t.Fatal(err)
	}
	if allocation.Slot != slotTwo {
		t.Fatalf("expected second slot allocation, got %#v", allocation.Slot)
	}
}

func TestAllocateConcurrentlyDoesNotDoubleBookRunnerSlot(t *testing.T) {
	service := New(repository.NewMemoryRepository())
	slot := domain.TimeSlot{Date: "2026-06-01", Start: "09:00", End: "11:00"}
	addRunner(t, service, "Asha", []domain.Skill{domain.SkillCook}, map[string]string{"diet": "veg"}, slot)

	const attempts = 20
	results := make(chan error, attempts)
	var wg sync.WaitGroup
	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			_, err := service.Allocate(context.Background(), domain.AllocateCommand{
				RequestID: fmt.Sprintf("req-%d", index),
				Skill:     domain.SkillCook,
				Slot:      slot,
			})
			results <- err
		}(i)
	}
	wg.Wait()
	close(results)

	successes := 0
	for err := range results {
		if err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("expected exactly one successful allocation, got %d", successes)
	}
}

func addRunner(t *testing.T, service Service, name string, skills []domain.Skill, attributes map[string]string, slot domain.TimeSlot) {
	t.Helper()
	_, err := service.AddRunner(context.Background(), domain.AddRunnerCommand{
		Name:       name,
		Skills:     skills,
		Attributes: attributes,
		Slots:      []domain.TimeSlot{slot},
	})
	if err != nil {
		t.Fatal(err)
	}
}
