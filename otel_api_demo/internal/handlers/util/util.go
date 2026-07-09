package util

import (
	"errors"
	"net/http"
)

type ErrorResponse struct {
	HttpCode int    `json:"http_code"`
	Message  string `json:"message"`
	Code     string `json:"code,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

func BuildError(httpCode int, message string) error {
	return &ErrorResponse{
		HttpCode: httpCode,
		Message:  message,
	}
}

func BuildErrorWithCode(httpCode int, message, code string) error {
	return &ErrorResponse{
		HttpCode: httpCode,
		Message:  message,
		Code:     code,
	}
}

func HandleError(w http.ResponseWriter, err error) {
	var e *ErrorResponse
	if errors.As(err, &e) {
		WriteJSON(w, e.HttpCode, e)
		return
	}

	WriteJSON(w, http.StatusInternalServerError, &ErrorResponse{
		HttpCode: http.StatusInternalServerError,
		Message:  err.Error(),
		Code:     "INTERNAL_ERROR",
	})
}

type SuccessResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}
