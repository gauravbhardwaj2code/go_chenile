package domain

import (
	"net/http"

	chenileerrors "base/errors"
)

const (
	ErrorOrderNameRequired = 1001
)

func NameRequired() error {
	return chenileerrors.Builder().
		Status(http.StatusBadRequest).
		Code(ErrorOrderNameRequired).
		MessageKey("order.name.required").
		Description("name is required").
		Field("name", "required", "name is required").
		Build()
}
