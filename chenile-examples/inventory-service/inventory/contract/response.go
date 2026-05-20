package contract

import "inventory-service/inventory/domain"

type Inventory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewInventoryResponse(entity domain.Inventory) Inventory {
	return Inventory{
		ID:   entity.ID,
		Name: entity.Name,
	}
}
