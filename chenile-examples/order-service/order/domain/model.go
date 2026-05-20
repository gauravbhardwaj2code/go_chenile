package domain

type CreateOrderCommand struct {
	Name string
}

type Order struct {
	ID   string
	Name string
}
