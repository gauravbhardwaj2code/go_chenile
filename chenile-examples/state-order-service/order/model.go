package order

type Order struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

func (o *Order) GetState() string {
	return o.State
}

func (o *Order) SetState(state string) {
	o.State = state
}
