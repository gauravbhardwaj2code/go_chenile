package contract

import "order-service/order/domain"

type Order struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewOrderResponse(entity domain.Order) Order {
	return Order{
		ID:   entity.ID,
		Name: entity.Name,
	}
}
