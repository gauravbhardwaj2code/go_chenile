package chenile

import (
	"context"
	"net/http"

	"core"
)

type Module interface {
	Name() string
	Register(*Builder) error
}

type module struct {
	name     string
	register func(*Builder) error
}

func NewModule(name string, register func(*Builder) error) Module {
	return module{name: name, register: register}
}

func (m module) Name() string {
	return m.name
}

func (m module) Register(builder *Builder) error {
	return m.register(builder)
}

type Builder struct {
	registry *core.Registry
}

func NewBuilder(registry *core.Registry) *Builder {
	return &Builder{registry: registry}
}

func (b *Builder) Service(id string) *ServiceBuilder {
	return &ServiceBuilder{
		registry: b.registry,
		service: core.ServiceDefinition{
			ID:   id,
			Name: id,
		},
	}
}

type ServiceBuilder struct {
	registry *core.Registry
	service  core.ServiceDefinition
}

func (b *ServiceBuilder) Name(name string) *ServiceBuilder {
	b.service.Name = name
	return b
}

func (b *ServiceBuilder) Routes(routes ...Route) error {
	for _, route := range routes {
		b.service.Operations = append(b.service.Operations, core.OperationDefinition{
			Name:     route.Name,
			Method:   route.Method,
			Path:     route.Path,
			NewInput: route.NewInput,
			Handler:  route.Handler,
		})
	}
	return b.registry.RegisterService(b.service)
}

type Route struct {
	Name     string
	Method   string
	Path     string
	NewInput func() any
	Handler  core.HandlerFunc
}

func GET[Resp any](path string, name string, handler func(context.Context) (Resp, error)) Route {
	return Route{
		Name:   name,
		Method: http.MethodGet,
		Path:   path,
		Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
			return handler(ctx)
		},
	}
}

func POST[Req any, Resp any](path string, name string, newRequest func() *Req, handler func(context.Context, Req) (Resp, error)) Route {
	return jsonRoute(http.MethodPost, path, name, newRequest, handler)
}

func PUT[Req any, Resp any](path string, name string, newRequest func() *Req, handler func(context.Context, Req) (Resp, error)) Route {
	return jsonRoute(http.MethodPut, path, name, newRequest, handler)
}

func PATCH[Req any, Resp any](path string, name string, newRequest func() *Req, handler func(context.Context, Req) (Resp, error)) Route {
	return jsonRoute(http.MethodPatch, path, name, newRequest, handler)
}

func DELETE[Req any, Resp any](path string, name string, newRequest func() *Req, handler func(context.Context, Req) (Resp, error)) Route {
	return jsonRoute(http.MethodDelete, path, name, newRequest, handler)
}

func jsonRoute[Req any, Resp any](method string, path string, name string, newRequest func() *Req, handler func(context.Context, Req) (Resp, error)) Route {
	return Route{
		Name:   name,
		Method: method,
		Path:   path,
		NewInput: func() any {
			return newRequest()
		},
		Handler: func(ctx context.Context, exchange *core.Exchange) (any, error) {
			request := exchange.Input.(*Req)
			return handler(ctx, *request)
		},
	}
}
