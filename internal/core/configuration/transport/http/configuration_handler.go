package http

import (
	"financing-offer/internal/core/configuration"
	"github.com/gin-gonic/gin"
	"net/http"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/handler"
	"financing-offer/pkg/cache"
)

type ConfigurationHandler struct {
	handler.BaseHandler
	cacheStore cache.Cache
	useCase    configuration.UseCase
}

func NewConfigurationHandler(baseHandler handler.BaseHandler, cache cache.Cache, useCase configuration.UseCase) *ConfigurationHandler {
	return &ConfigurationHandler{
		BaseHandler: baseHandler,
		cacheStore:  cache,
		useCase:     useCase,
	}
}

// SetLoanRate godoc
//
//	@Summary		Set loan rate
//	@Description	Set loan rate
//	@Tags			configuration,admin
//	@Accept			json
//	@Produce		json
//	@Param			loanRate	body		SetLoanRateRequest	true	"loan rate"
//	@Success		200						{object}	handler.BaseResponse[entity.LoanRateConfiguration]
//	@Failure		400						{object}	handler.ErrorResponse
//	@Failure		500						{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/configurations/loan-rate [post]
func (h *ConfigurationHandler) SetLoanRate(ctx *gin.Context) {
	request := SetLoanRateRequest{}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	req := entity.LoanRateConfiguration{
		Ids: request.Ids,
	}
	result, err := h.useCase.SetLoanRate(ctx, req, h.UserSubOrEmpty(ctx))
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.LoanRateConfiguration]{
			Data: result,
		},
	)
}

// GetLoanRate godoc
//
//	@Summary		Get loan rate
//	@Description	Get loan rate
//	@Tags			promotion,investor
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[entity.LoanRateConfiguration]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/configurations/loan-rate [get]
func (h *ConfigurationHandler) GetLoanRate(ctx *gin.Context) {
	result, err := h.useCase.GetLoanRate(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.LoanRateConfiguration]{
			Data: result,
		},
	)
}

// SetMarginPool godoc
//
//	@Summary		Set margin pool
//	@Description	Set margin pool
//	@Tags			configuration,admin
//	@Accept			json
//	@Produce		json
//	@Param			marginPool	body		SetMarginPoolRequest	true	"margin pool
//	@Success		200						{object}	handler.BaseResponse[entity.MarginPoolConfiguration]
//	@Failure		400						{object}	handler.ErrorResponse
//	@Failure		500						{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/configurations/margin-pool [post]

func (h *ConfigurationHandler) SetMarginPool(ctx *gin.Context) {
	request := SetMarginPoolRequest{}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	req := entity.MarginPoolConfiguration{
		Ids: request.Ids,
	}
	result, err := h.useCase.SetMarginPool(ctx, req, h.UserSubOrEmpty(ctx))
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.MarginPoolConfiguration]{
			Data: result,
		},
	)
}

// GetMarginPool godoc
//
//	@Summary		Get margin pool
//	@Description	Get margin pool
//	@Tags			promotion,investor
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[entity.MarginPoolConfiguration]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/configurations/margin-pool [get]
func (h *ConfigurationHandler) GetMarginPool(ctx *gin.Context) {
	result, err := h.useCase.GetMarginPool(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.MarginPoolConfiguration]{
			Data: result,
		},
	)
}
