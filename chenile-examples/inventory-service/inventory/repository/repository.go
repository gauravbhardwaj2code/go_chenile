package repository

import (
	"context"

	"inventory-service/inventory/domain"
)

type Repository interface {
	Create(context.Context, domain.Inventory) (domain.Inventory, error)
}
