package apperrors

import (
	"net/http"
)

type AppError struct {
	Err     error
	Code    int
	Message string
}

type AppErrorOption func(appError *AppError)

func New(err error, opts ...AppErrorOption) AppError {
	appErr := AppError{
		Err:  err,
		Code: http.StatusInternalServerError,
	}
	for _, opt := range opts {
		opt(&appErr)
	}
	return appErr
}

func (e AppError) Error() string {
	if e.Err != nil {
		return e.Message + e.Err.Error()
	}
	return e.Message
}

func (e AppError) Unwrap() error {
	return e.Err
}

func WithCode(code int) AppErrorOption {
	return func(appError *AppError) {
		appError.Code = code
	}
}

func WithMessage(message string) AppErrorOption {
	return func(appError *AppError) {
		appError.Message = message
	}
}
