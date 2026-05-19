package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	chenileerrors "base/errors"
	"core"
)

type RequestIDGenerator func() string

type RequestID struct {
	Header    string
	Generate  RequestIDGenerator
	Generated string
}

func (m *RequestID) Before(ctx context.Context, exchange *core.Exchange) error {
	header := m.Header
	if header == "" {
		header = "X-Request-Id"
	}
	if exchange.Headers == nil {
		exchange.Headers = map[string]string{}
	}
	if exchange.Headers[header] == "" && m.Generate != nil {
		exchange.Headers[header] = m.Generate()
		m.Generated = exchange.Headers[header]
	}
	return nil
}

func (m *RequestID) After(context.Context, *core.Exchange) error {
	return nil
}

type Recovery struct{}

func (Recovery) Before(context.Context, *core.Exchange) error {
	return nil
}

func (Recovery) After(context.Context, *core.Exchange) error {
	return nil
}

func Recover(fn func() error) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = chenileerrors.Builder().
				Status(http.StatusInternalServerError).
				MessageKey("chenile.panic").
				Description(fmt.Sprintf("panic recovered: %v", recovered)).
				Param("stack", string(debug.Stack())).
				Build()
		}
	}()
	return fn()
}

type Validator interface {
	Validate() error
}

type ValidateInput struct{}

func (ValidateInput) Before(ctx context.Context, exchange *core.Exchange) error {
	input, ok := exchange.Input.(Validator)
	if !ok || input == nil {
		return nil
	}
	return input.Validate()
}

func (ValidateInput) After(context.Context, *core.Exchange) error {
	return nil
}
