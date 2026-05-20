package contract

import "customer-service/customer/domain"

type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewCustomerResponse(entity domain.Customer) Customer {
	return Customer{
		ID:   entity.ID,
		Name: entity.Name,
	}
}
