package ragflow

import "fmt"

type APIError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("RAGFlow API error (code: %d, status: %d): %s", e.Code, e.StatusCode, e.Message)
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	ErrorCodeBadRequest           = 400
	ErrorCodeUnauthorized         = 401
	ErrorCodeForbidden            = 403
	ErrorCodeNotFound             = 404
	ErrorCodeInternalServerError  = 500
	ErrorCodeDuplicatedName       = 1001
	ErrorCodeFileTypeNotSupported = 1002
	ErrorCodeSuccess              = 200
	ErrorCodeGenericSuccess       = 0
)

func IsErrorCode(err error, code int) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Code == code
	}
	return false
}