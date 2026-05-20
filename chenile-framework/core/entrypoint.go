package core

import (
	"context"
	"errors"
	"net/http"

	chenileerrors "github.com/gauravbhardwaj2code/go_chenile/chenile-framework/base/errors"
	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/base/response"
	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/owiz/chain"
)

type Interceptor interface {
	Before(context.Context, *Exchange) error
	After(context.Context, *Exchange) error
}

type EntryPoint struct {
	registry     *Registry
	interceptors []Interceptor
}

func NewEntryPoint(registry *Registry, interceptors ...Interceptor) *EntryPoint {
	return &EntryPoint{
		registry:     registry,
		interceptors: append([]Interceptor{}, interceptors...),
	}
}

func (e *EntryPoint) Execute(ctx context.Context, exchange *Exchange) (result response.GenericResponse) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = errorResponse(chenileerrors.Builder().
				Status(http.StatusInternalServerError).
				MessageKey("chenile.panic").
				Description("panic recovered during service execution").
				Param("panic", recovered).
				Build())
		}
	}()
	if exchange.Context == nil {
		exchange.Context = ctx
	}
	operation, ok := e.registry.Operation(exchange.ServiceID, exchange.OperationName)
	if !ok {
		return response.Failure(http.StatusNotFound, "operation not found", 0)
	}
	beforeChain := chain.New[Exchange]()
	for _, interceptor := range e.interceptors {
		current := interceptor
		beforeChain.Add(chain.CommandFunc[Exchange](func(ctx context.Context, exchange *Exchange) error {
			return current.Before(ctx, exchange)
		}))
	}
	if err := beforeChain.Execute(ctx, exchange); err != nil {
		return errorResponse(err)
	}
	payload, err := operation.Handler(ctx, exchange)
	if err != nil {
		return errorResponse(err)
	}
	exchange.Output = payload
	afterChain := chain.New[Exchange]()
	for i := len(e.interceptors) - 1; i >= 0; i-- {
		current := e.interceptors[i]
		afterChain.Add(chain.CommandFunc[Exchange](func(ctx context.Context, exchange *Exchange) error {
			return current.After(ctx, exchange)
		}))
	}
	if err := afterChain.Execute(ctx, exchange); err != nil {
		return errorResponse(err)
	}
	if exchange.Err != nil {
		return errorResponse(exchange.Err)
	}
	return response.Success(payload)
}

func errorResponse(err error) response.GenericResponse {
	var chenileErr chenileerrors.ChenileError
	if errors.As(err, &chenileErr) {
		result := response.Failure(chenileErr.Status, chenileErr.Description, chenileErr.SubErrorCode)
		result.Errors[0].MessageKey = chenileErr.MessageKey
		result.Errors[0].Code = chenileErr.Code
		if chenileErr.Code != 0 {
			result.Code = chenileErr.Status
		}
		for _, field := range chenileErr.Fields {
			result.Errors = append(result.Errors, response.ResponseMessage{
				Description: field.Description,
				MessageKey:  field.MessageKey,
				Severity:    response.Error,
				Field:       field.Field,
			})
		}
		return result
	}
	return response.Failure(http.StatusInternalServerError, err.Error(), 0)
}
