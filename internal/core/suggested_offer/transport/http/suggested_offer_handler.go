package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/entity"
	suggestedOffer "financing-offer/internal/core/suggested_offer"
	suggestedOfferConfig "financing-offer/internal/core/suggested_offer_config"
	"financing-offer/internal/handler"
)

type SuggestedOfferHandler struct {
	handler.BaseHandler
	logger                      *slog.Logger
	suggestedOfferUseCase       suggestedOffer.UseCase
	suggestedOfferConfigUseCase suggestedOfferConfig.UseCase
}

func NewSuggestedOfferHandler(
	baseHandler handler.BaseHandler,
	logger *slog.Logger,
	useCase suggestedOffer.UseCase,
	configUseCase suggestedOfferConfig.UseCase,
) *SuggestedOfferHandler {
	return &SuggestedOfferHandler{
		BaseHandler:                 baseHandler,
		logger:                      logger,
		suggestedOfferUseCase:       useCase,
		suggestedOfferConfigUseCase: configUseCase,
	}
}

// CreateOffer godoc
//
//	@Summary		Create suggested offer
//	@Description	Create suggested offer
//	@Tags			suggested offer
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SuggestedOfferRequest	true	"body"
//	@Success		201		{object}	handler.BaseResponse[entity.SuggestedOffer]
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/suggested-offers [post]
func (h *SuggestedOfferHandler) CreateOffer(ctx *gin.Context) {
	errorMessage := "Suggested Offer CreateOffer"
	request := &SuggestedOfferRequest{}
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	investor, err := h.Investor(ctx)
	if err != nil {
		h.RenderUnauthenticated(ctx, err.Error())
		return
	}
	offer, err := h.suggestedOfferUseCase.CreateSuggestedOffer(
		ctx,
		investor.InvestorId,
		investor.CustodyCode,
		entity.SuggestedOffer{ConfigId: request.ConfigId, AccountNo: request.AccountNo, Symbols: request.Symbols},
	)
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.SuggestedOffer]{
			Data: offer,
		},
	)
}
