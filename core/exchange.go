package core

import (
	"context"

	"github.com/ajapro/chenile-go/base/response"
)

type Exchange struct {
	Context       context.Context
	ServiceID     string
	OperationName string
	Method        string
	Path          string
	Headers       map[string]string
	QueryParams   map[string]string
	PathParams    map[string]string
	Body          []byte
	Input         any
	Output        any
	Response      *response.GenericResponse
	Err           error
}
