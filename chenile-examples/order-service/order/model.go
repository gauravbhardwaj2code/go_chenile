package order

type CreateOrderRequest struct {
	Name string `json:"name"`
}

type Order struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
