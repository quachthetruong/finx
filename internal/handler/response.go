package handler

import "financing-offer/internal/core"

type BaseResponse[T any] struct {
	Data T `json:"data"`
}

type ResponseWithPaging[T any] struct {
	Data     T                   `json:"data"`
	MetaData core.PagingMetaData `json:"metaData"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
