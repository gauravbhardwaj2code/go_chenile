package domain

import (
	"net/http"

	chenileerrors "base/errors"
)

const (
	ErrorCustomerNameRequired = 1001
)

func NameRequired() error {
	return chenileerrors.Builder().
		Status(http.StatusBadRequest).
		Code(ErrorCustomerNameRequired).
		MessageKey("customer.name.required").
		Description("name is required").
		Field("name", "required", "name is required").
		Build()
}
