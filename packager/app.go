package packager

import (
	"fmt"
	"net/http"

	"core"
	chenilehttp "http"
)

type Module struct {
	Name     string
	Register func(*core.Registry) error
}

type App struct {
	Registry *core.Registry
	Router   *chenilehttp.Router
	Modules  []Module
}

func NewWebApp(modules ...Module) (*App, error) {
	registry := core.NewRegistry()
	for _, module := range modules {
		if module.Name == "" {
			return nil, fmt.Errorf("module name is required")
		}
		if module.Register == nil {
			return nil, fmt.Errorf("module %q register function is required", module.Name)
		}
		if err := module.Register(registry); err != nil {
			return nil, fmt.Errorf("register module %q: %w", module.Name, err)
		}
	}
	entryPoint := core.NewEntryPoint(registry)
	router := chenilehttp.NewRouter(entryPoint)
	if err := router.MountRegistry(registry); err != nil {
		return nil, err
	}
	return &App{
		Registry: registry,
		Router:   router,
		Modules:  append([]Module{}, modules...),
	}, nil
}

func (a *App) ListenAndServe(address string) error {
	return http.ListenAndServe(address, a.Router)
}
