package stateentity

import (
	"context"
	"net/http"

	"core"
	"stm"
)

type Repository[E stm.Entity] interface {
	Create(context.Context, E) (E, error)
	Get(context.Context, string) (E, error)
	Save(context.Context, E) (E, error)
}

type NewEntity[E stm.Entity] func() E

type ApplyEventRequest struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	Param any    `json:"param,omitempty"`
}

type Options[E stm.Entity] struct {
	ServiceID string
	Name      string
	BasePath  string
	NewEntity NewEntity[E]
	Repo      Repository[E]
	Machine   *stm.Machine
}

func Register[E stm.Entity](registry *core.Registry, options Options[E]) error {
	serviceID := options.ServiceID
	if serviceID == "" {
		serviceID = options.Name
	}
	basePath := options.BasePath
	if basePath == "" {
		basePath = "/" + options.Name
	}
	return registry.RegisterService(core.ServiceDefinition{
		ID:   serviceID,
		Name: options.Name,
		Operations: []core.OperationDefinition{
			{
				Name:   "create",
				Method: http.MethodPost,
				Path:   basePath,
				NewInput: func() any {
					return options.NewEntity()
				},
				Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
					entity := exchange.Input.(E)
					if entity.GetState() == "" && options.Machine != nil {
						processed, err := options.Machine.Process(ctx, entity, stm.InitialEvent, nil)
						if err != nil {
							var zero E
							return zero, err
						}
						entity = processed.(E)
					}
					return options.Repo.Create(ctx, entity)
				},
			},
			{
				Name:   "applyEvent",
				Method: http.MethodPost,
				Path:   basePath + "/event",
				NewInput: func() any {
					return &ApplyEventRequest{}
				},
				Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
					request := exchange.Input.(*ApplyEventRequest)
					entity, err := options.Repo.Get(ctx, request.ID)
					if err != nil {
						var zero E
						return zero, err
					}
					processed, err := options.Machine.Process(ctx, entity, request.Event, request.Param)
					if err != nil {
						var zero E
						return zero, err
					}
					return options.Repo.Save(ctx, processed.(E))
				},
			},
		},
	})
}
