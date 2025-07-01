package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/financialproduct"
	"financing-offer/internal/handler"
)

type FinancialProductHandler struct {
	handler.BaseHandler
	useCase financialproduct.UseCase
}

func NewFinancialProductHandler(
	baseHandler handler.BaseHandler,
	useCase financialproduct.UseCase,
) *FinancialProductHandler {
	return &FinancialProductHandler{
		BaseHandler: baseHandler,
		useCase:     useCase,
	}
}

// GetLoanRates godoc
//
//	@Summary		Get loan rates from config
//	@Description	Get loan rates from config
//	@Tags			loan rates, admin
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	handler.BaseResponse[[]entity.LoanRate]
//	@Failure		500	{object}	handler.ErrorResponse
//	@Router			/v1/mo/applicable-loan-rates [get]
func (h *FinancialProductHandler) GetLoanRates(ctx *gin.Context) {
	res, err := h.useCase.GetLoanRates(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK,
		handler.BaseResponse[[]entity.LoanRate]{
			Data: res,
		},
	)
}

// GetMarginPools godoc
//
//	@Summary		Get margin pools from config
//	@Description	Get margin pools from config
//	@Tags			margin pools, admin
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	handler.BaseResponse[[]entity.MarginPool]
//	@Failure		500	{object}	handler.ErrorResponse
//	@Router			/v1/mo/applicable-margin-pools [get]
func (h *FinancialProductHandler) GetMarginPools(ctx *gin.Context) {
	res, err := h.useCase.GetMarginPools(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK,
		handler.BaseResponse[[]entity.MarginPool]{
			Data: res,
		},
	)
}
