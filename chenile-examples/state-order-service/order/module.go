package order

import (
	"core"
	"stateentity"
	"stm"
)

func Register(registry *core.Registry) error {
	machine, err := stm.FromJSON([]byte(`{
		"created": {"initial": true, "events": {"confirm": {"newState": "confirmed"}, "cancel": {"newState": "cancelled"}}},
		"confirmed": {"events": {"fulfill": {"newState": "fulfilled"}}},
		"fulfilled": {},
		"cancelled": {}
	}`), stm.Actions{})
	if err != nil {
		return err
	}
	return stateentity.Register[*Order](registry, stateentity.Options[*Order]{
		ServiceID: "stateOrderService",
		Name:      "state-orders",
		BasePath:  "/state-orders",
		NewEntity: func() *Order { return &Order{} },
		Repo:      NewRepository(),
		Machine:   machine,
	})
}
