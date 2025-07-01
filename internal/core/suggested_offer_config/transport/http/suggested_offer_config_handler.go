package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/suggested_offer_config"
	"financing-offer/internal/handler"
	"financing-offer/pkg/optional"
)

type SuggestedOfferConfigHandler struct {
	handler.BaseHandler
	logger  *slog.Logger
	useCase suggested_offer_config.UseCase
}

func NewSuggestedOfferConfigHandler(baseHandler handler.BaseHandler, logger *slog.Logger, useCase suggested_offer_config.UseCase) *SuggestedOfferConfigHandler {
	return &SuggestedOfferConfigHandler{BaseHandler: baseHandler, logger: logger, useCase: useCase}
}

// Create godoc
//
//	@Summary		Create suggested offer config
//	@Description	Create suggested offer config
//	@Tags			suggested offer config,admin
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateSuggestedOfferConfigRequest	true	"body"
//	@Success		201		{object}	handler.BaseResponse[entity.SuggestedOfferConfig]
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/suggested-offer-configs [post]
func (h *SuggestedOfferConfigHandler) Create(ctx *gin.Context) {
	errorMessage := "Suggested Offer Config Handler Create"
	req := &CreateSuggestedOfferConfigRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	suggestedOfferConfig := req.toEntity()
	suggestedOfferConfig.CreatedBy = appcontext.ContextGetCustomerInfo(ctx).Sub
	created, err := h.useCase.Create(ctx, suggestedOfferConfig)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusCreated, handler.BaseResponse[entity.SuggestedOfferConfig]{
			Data: created,
		},
	)
}

// Update godoc
//
//	@Summary		Update suggested offer config
//	@Description	Update suggested offer config
//	@Tags			suggested offer config,admin
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"id"
//	@Param			request	body		entity.SuggestedOfferConfig	true	"body"
//	@Success		200		{object}	handler.BaseResponse[entity.SuggestedOfferConfig]
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/suggested-offer-configs/{id} [patch]
func (h *SuggestedOfferConfigHandler) Update(ctx *gin.Context) {
	errorMessage := "Suggested Offer Config Handler Update"
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	oldSuggestedOfferConfig, err := h.useCase.GetById(ctx, id)
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderError(ctx, err)
		return
	}
	if err := ctx.ShouldBindJSON(&oldSuggestedOfferConfig); err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	oldSuggestedOfferConfig.LastUpdatedBy = appcontext.ContextGetCustomerInfo(ctx).Sub
	updated, err := h.useCase.Update(ctx, oldSuggestedOfferConfig)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.SuggestedOfferConfig]{
			Data: updated,
		},
	)
}

// Get godoc
//
//	@Summary		Get active suggested offer config
//	@Description	Get active suggested offer config
//	@Tags			suggested offer config,admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[optional.Optional[entity.SuggestedOfferConfig]]
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/active-suggested-offer-config [get]
func (h *SuggestedOfferConfigHandler) Get(ctx *gin.Context) {
	activeSuggestedOfferConfig, err := h.useCase.GetActiveSuggestedOfferConfig(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[optional.Optional[entity.SuggestedOfferConfig]]{
			Data: activeSuggestedOfferConfig,
		},
	)
}

// GetAll godoc
//
//	@Summary		Get all suggested offer config
//	@Description	Get all suggested offer config
//	@Tags			suggested offer config,admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[[]entity.SuggestedOfferConfig]
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/suggested-offer-configs [get]
func (h *SuggestedOfferConfigHandler) GetAll(ctx *gin.Context) {
	suggestedOfferConfigs, err := h.useCase.GetAll(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]entity.SuggestedOfferConfig]{
			Data: suggestedOfferConfigs,
		},
	)
}

// GetById godoc
//
//	@Summary		Get suggested offer config by id
//	@Description	Get suggested offer config by id
//	@Tags			suggested offer config,admin
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		200	{object}	handler.BaseResponse[entity.SuggestedOfferConfig]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/suggested-offer-configs/{id} [get]
func (h *SuggestedOfferConfigHandler) GetById(ctx *gin.Context) {
	errorMessage := "Suggested Offer Config Handler GetById"
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	res, err := h.useCase.GetById(ctx, id)
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.SuggestedOfferConfig]{
			Data: res,
		},
	)
}

// UpdateStatus godoc
//
//	@Summary		Update suggested offer config status
//	@Description	Update suggested offer config status
//	@Tags			suggested offer config,admin
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int										true	"id"
//	@Param			status	body		UpdateSuggestedOfferConfigStatusRequest	true	"status"
//	@Success		200		{object}	handler.BaseResponse[entity.SuggestedOfferConfig]
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/suggested-offer-configs/{id}/status [patch]
func (h *SuggestedOfferConfigHandler) UpdateStatus(ctx *gin.Context) {
	errorMessage := "Suggested Offer Config Handler UpdateStatus"
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	var status UpdateSuggestedOfferConfigStatusRequest
	if err = ctx.ShouldBindJSON(&status); err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	updated, err := h.useCase.UpdateStatus(ctx, id, status.Status, h.UserSubOrEmpty(ctx))
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.SuggestedOfferConfig]{
			Data: updated,
		},
	)
}
