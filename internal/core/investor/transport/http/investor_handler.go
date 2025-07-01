package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/investor"
	"financing-offer/internal/handler"
)

type InvestorHandler struct {
	handler.BaseHandler
	logger          *slog.Logger
	investorUseCase investor.UseCase
}

// FillInvestorIdsFromRequests godoc
//
//	@Summary		Fill investor ids from requests
//	@Description	Fill investor ids from requests
//	@Tags			investor,admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[int]
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/internal/investors/sync-investor-data [post]
func (h *InvestorHandler) FillInvestorIdsFromRequests(ctx *gin.Context) {
	investorIdsCount, err := h.investorUseCase.FillInvestorIdsFromRequests(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[int]{
			Data: investorIdsCount,
		},
	)
}

func NewInvestorHandler(baseHandler handler.BaseHandler, logger *slog.Logger, investorUseCase investor.UseCase) *InvestorHandler {
	return &InvestorHandler{BaseHandler: baseHandler, logger: logger, investorUseCase: investorUseCase}
}
