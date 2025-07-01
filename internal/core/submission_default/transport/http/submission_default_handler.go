package http

import (
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/submission_default"
	"financing-offer/internal/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SubmissionDefaultHandler struct {
	handler.BaseHandler
	useCase submission_default.UseCase
}

func NewSubmissionDefaultHandler(baseHandler handler.BaseHandler, useCase submission_default.UseCase) *SubmissionDefaultHandler {
	return &SubmissionDefaultHandler{
		BaseHandler: baseHandler,
		useCase:     useCase,
	}
}

func (h *SubmissionDefaultHandler) GetSubmissionDefault(ctx *gin.Context) {
	res, err := h.useCase.GetSubmissionDefault(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.SubmissionDefault]{
			Data: res,
		},
	)
}

func (h *SubmissionDefaultHandler) SetSubmissionDefault(ctx *gin.Context) {
	request := SetSubmissionDefaultRequest{}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	req := entity.SubmissionDefault{
		FirmSellingFeeRate:       request.FirmSellingFeeRate,
		FirmBuyingFeeRate:        request.FirmBuyingFeeRate,
		TransferFee:              request.TransferFee,
		AllowedOverdueLoanInDays: request.AllowedOverdueLoanInDays,
	}
	result, err := h.useCase.SetSubmissionDefault(ctx, req, h.UserSubOrEmpty(ctx))
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.SubmissionDefault]{
			Data: result,
		},
	)
}
