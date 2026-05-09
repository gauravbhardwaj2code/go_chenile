package errors

import "fmt"

type ChenileError struct {
	Status       int
	SubErrorCode int
	Description  string
}

func (e ChenileError) Error() string {
	if e.Description == "" {
		return fmt.Sprintf("chenile error status=%d subErrorCode=%d", e.Status, e.SubErrorCode)
	}
	return e.Description
}

func New(status int, subErrorCode int, description string) ChenileError {
	return ChenileError{
		Status:       status,
		SubErrorCode: subErrorCode,
		Description:  description,
	}
}
