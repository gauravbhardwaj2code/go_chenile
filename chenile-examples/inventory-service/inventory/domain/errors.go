package domain

import (
	"net/http"

	chenileerrors "github.com/gauravbhardwaj2code/go_chenile/chenile-framework/base/errors"
)

const (
	ErrorInventoryNameRequired = 1001
)

func NameRequired() error {
	return chenileerrors.Builder().
		Status(http.StatusBadRequest).
		Code(ErrorInventoryNameRequired).
		MessageKey("inventory.name.required").
		Description("name is required").
		Field("name", "required", "name is required").
		Build()
}
