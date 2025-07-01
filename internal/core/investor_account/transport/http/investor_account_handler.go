package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/investor_account"
	"financing-offer/internal/handler"
)

type InvestorAccountHandler struct {
	handler.BaseHandler
	logger  *slog.Logger
	useCase investor_account.UseCase
}

func NewInvestorAccountHandler(baseHandler handler.BaseHandler, logger *slog.Logger, useCase investor_account.UseCase) *InvestorAccountHandler {
	return &InvestorAccountHandler{
		BaseHandler: baseHandler,
		logger:      logger,
		useCase:     useCase,
	}
}

// VerifyAndUpdateInvestorAccountMarginStatus godoc
//
//	@Summary		Verify and update investor_account margin status
//	@Description	Verify and update investor_account margin status
//	@Tags			investor account margin status,admin
//	@Accept			json
//	@Produce		json
//	@Param			accountNo	path		string										true	"account-no"
//	@Param			request		body		VerifyInvestorAccountMarginStatusRequest	true	"body"
//	@Success		200			{object}	handler.BaseResponse[entity.InvestorAccount]
//	@Failure		400			{object}	handler.ErrorResponse
//	@Failure		404			{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/investor-accounts/{account-no}/margin-status [put]
func (h *InvestorAccountHandler) VerifyAndUpdateInvestorAccountMarginStatus(ctx *gin.Context) {
	accountNo, err := h.ParamNotEmpty(ctx, "account-no")
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	var req VerifyInvestorAccountMarginStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	res, err := h.useCase.VerifyAndUpdateInvestorAccountMarginStatus(ctx, req.toEntity(accountNo))
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.InvestorAccount]{Data: res})
}
