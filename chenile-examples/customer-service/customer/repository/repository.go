package repository

import (
	"context"

	"customer-service/customer/domain"
)

type Repository interface {
	Create(context.Context, domain.Customer) (domain.Customer, error)
}
