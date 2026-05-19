package errors

import (
	"fmt"
	"net/http"
)

type Severity string

const (
	Error Severity = "ERROR"
	Warn  Severity = "WARN"
)

type FieldError struct {
	Field       string `json:"field,omitempty"`
	MessageKey  string `json:"messageKey,omitempty"`
	Description string `json:"description,omitempty"`
}

type ChenileError struct {
	Status       int            `json:"status"`
	Code         int            `json:"code,omitempty"`
	SubErrorCode int            `json:"subErrorCode,omitempty"`
	MessageKey   string         `json:"messageKey,omitempty"`
	Description  string         `json:"description,omitempty"`
	Severity     Severity       `json:"severity,omitempty"`
	Params       map[string]any `json:"params,omitempty"`
	Fields       []FieldError   `json:"fields,omitempty"`
}

func (e ChenileError) Error() string {
	if e.Description == "" {
		return fmt.Sprintf("chenile error status=%d code=%d subErrorCode=%d", e.Status, e.Code, e.SubErrorCode)
	}
	return e.Description
}

func New(status int, subErrorCode int, description string) ChenileError {
	return Builder().Status(status).SubErrorCode(subErrorCode).Description(description).Build()
}

type ErrorBuilder struct {
	err ChenileError
}

func Builder() ErrorBuilder {
	return ErrorBuilder{err: ChenileError{
		Status:   http.StatusInternalServerError,
		Severity: Error,
	}}
}

func (b ErrorBuilder) Status(status int) ErrorBuilder {
	b.err.Status = status
	return b
}

func (b ErrorBuilder) Code(code int) ErrorBuilder {
	b.err.Code = code
	return b
}

func (b ErrorBuilder) SubErrorCode(code int) ErrorBuilder {
	b.err.SubErrorCode = code
	return b
}

func (b ErrorBuilder) MessageKey(key string) ErrorBuilder {
	b.err.MessageKey = key
	return b
}

func (b ErrorBuilder) Description(description string) ErrorBuilder {
	b.err.Description = description
	return b
}

func (b ErrorBuilder) Severity(severity Severity) ErrorBuilder {
	b.err.Severity = severity
	return b
}

func (b ErrorBuilder) Param(name string, value any) ErrorBuilder {
	if b.err.Params == nil {
		b.err.Params = map[string]any{}
	}
	b.err.Params[name] = value
	return b
}

func (b ErrorBuilder) Field(field string, messageKey string, description string) ErrorBuilder {
	b.err.Fields = append(b.err.Fields, FieldError{
		Field:       field,
		MessageKey:  messageKey,
		Description: description,
	})
	return b
}

func (b ErrorBuilder) Build() ChenileError {
	if b.err.Status == 0 {
		b.err.Status = http.StatusInternalServerError
	}
	if b.err.Severity == "" {
		b.err.Severity = Error
	}
	return b.err
}
