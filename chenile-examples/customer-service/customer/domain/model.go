package domain

type CreateCustomerCommand struct {
	Name string
}

type Customer struct {
	ID   string
	Name string
}
