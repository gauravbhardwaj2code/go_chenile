package response

import "net/http"

type ErrorType string

const (
	Error ErrorType = "ERROR"
	Warn  ErrorType = "WARN"
)

type ResponseMessage struct {
	Description  string    `json:"description,omitempty"`
	Code         int       `json:"code,omitempty"`
	SubErrorCode int       `json:"subErrorCode,omitempty"`
	MessageKey   string    `json:"messageKey,omitempty"`
	Severity     ErrorType `json:"severity,omitempty"`
	Field        string    `json:"field,omitempty"`
	Params       []any     `json:"params,omitempty"`
}

type GenericResponse struct {
	Code         int               `json:"code"`
	Description  string            `json:"description,omitempty"`
	Errors       []ResponseMessage `json:"errors,omitempty"`
	Payload      any               `json:"payload,omitempty"`
	Severity     ErrorType         `json:"severity,omitempty"`
	SubErrorCode int               `json:"subErrorCode,omitempty"`
	Success      bool              `json:"success"`
}

func Success(payload any) GenericResponse {
	return GenericResponse{
		Code:    http.StatusOK,
		Payload: payload,
		Success: true,
	}
}

func Failure(status int, description string, subErrorCode int) GenericResponse {
	if status == 0 {
		status = http.StatusInternalServerError
	}
	message := ResponseMessage{
		Description:  description,
		Code:         status,
		SubErrorCode: subErrorCode,
		Severity:     Error,
	}
	return GenericResponse{
		Code:         status,
		Description:  description,
		Errors:       []ResponseMessage{message},
		Severity:     Error,
		SubErrorCode: subErrorCode,
		Success:      false,
	}
}
