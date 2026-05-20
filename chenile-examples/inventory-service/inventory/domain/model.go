package domain

type CreateInventoryCommand struct {
	Name string
}

type Inventory struct {
	ID   string
	Name string
}
