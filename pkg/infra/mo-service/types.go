package mo_service

import "fmt"

type ListResponse[T any] struct {
	Data  []T `json:"data"`
	Total int `json:"total"`
	Start int `json:"start"`
	End   int `json:"end"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
}
