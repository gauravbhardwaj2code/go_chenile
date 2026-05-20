package repository

import (
	"context"
	"fmt"
	"sync"

	"inventory-service/inventory/domain"
)

type MemoryRepository struct {
	mu     sync.Mutex
	nextID int
	values map[string]domain.Inventory
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextID: 1,
		values: map[string]domain.Inventory{},
	}
}

func (r *MemoryRepository) Create(ctx context.Context, entity domain.Inventory) (domain.Inventory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if entity.ID == "" {
		entity.ID = fmt.Sprintf("inventory-%d", r.nextID)
		r.nextID++
	}
	r.values[entity.ID] = entity
	return entity, nil
}
