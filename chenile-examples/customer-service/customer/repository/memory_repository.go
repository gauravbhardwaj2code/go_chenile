package repository

import (
	"context"
	"fmt"
	"sync"

	"customer-service/customer/domain"
)

type MemoryRepository struct {
	mu     sync.Mutex
	nextID int
	values map[string]domain.Customer
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextID: 1,
		values: map[string]domain.Customer{},
	}
}

func (r *MemoryRepository) Create(ctx context.Context, entity domain.Customer) (domain.Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if entity.ID == "" {
		entity.ID = fmt.Sprintf("customer-%d", r.nextID)
		r.nextID++
	}
	r.values[entity.ID] = entity
	return entity, nil
}
