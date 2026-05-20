package repository

import (
	"context"

	"order-service/order/domain"
)

type Repository interface {
	Create(context.Context, domain.Order) (domain.Order, error)
}
