package core

import "context"

type HandlerFunc func(context.Context, *Exchange) (any, error)

type OperationDefinition struct {
	Name     string
	Method   string
	Path     string
	NewInput func() any
	Handler  HandlerFunc
}

type ServiceDefinition struct {
	ID         string
	Name       string
	Operations []OperationDefinition
}

type RegisteredOperation struct {
	Service   ServiceDefinition
	Operation OperationDefinition
}
