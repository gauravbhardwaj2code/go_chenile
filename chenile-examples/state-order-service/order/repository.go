package order

import (
	"context"
	"fmt"
)

type Repository struct {
	values map[string]*Order
}

func NewRepository() *Repository {
	return &Repository{values: map[string]*Order{}}
}

func (r *Repository) Create(ctx context.Context, order *Order) (*Order, error) {
	if order.ID == "" {
		order.ID = "order-1"
	}
	r.values[order.ID] = order
	return order, nil
}

func (r *Repository) Get(ctx context.Context, id string) (*Order, error) {
	order := r.values[id]
	if order == nil {
		return nil, fmt.Errorf("order %q not found", id)
	}
	return order, nil
}

func (r *Repository) Save(ctx context.Context, order *Order) (*Order, error) {
	r.values[order.ID] = order
	return order, nil
}
