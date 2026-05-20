package repository

import (
	"context"
	"fmt"
	"sync"

	"order-service/order/domain"
)

type MemoryRepository struct {
	mu     sync.Mutex
	nextID int
	values map[string]domain.Order
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextID: 1,
		values: map[string]domain.Order{},
	}
}

func (r *MemoryRepository) Create(ctx context.Context, entity domain.Order) (domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if entity.ID == "" {
		entity.ID = fmt.Sprintf("order-%d", r.nextID)
		r.nextID++
	}
	r.values[entity.ID] = entity
	return entity, nil
}
