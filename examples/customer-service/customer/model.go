package customer

type CreateCustomerRequest struct {
	Name string `json:"name"`
}

type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
