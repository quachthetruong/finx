package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/promotion_campaign"
	"financing-offer/internal/handler"
)

type PromotionCampaignHandler struct {
	handler.BaseHandler
	logger  *slog.Logger
	useCase promotion_campaign.UseCase
}

func NewPromotionCampaignHandler(baseHandler handler.BaseHandler, logger *slog.Logger, useCase promotion_campaign.UseCase) *PromotionCampaignHandler {
	return &PromotionCampaignHandler{BaseHandler: baseHandler, logger: logger, useCase: useCase}
}

// Create godoc
//
//	@Summary		Create promotion campaign
//	@Description	Create promotion campaign
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreatePromotionCampaignRequest	true	"request"
//	@Success		201		{object}	handler.BaseResponse[entity.PromotionCampaign]
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/promotion-campaigns [post]
func (h *PromotionCampaignHandler) Create(ctx *gin.Context) {
	req := &CreatePromotionCampaignRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.logger.Error("create campaign", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	campaign := req.toEntity()
	campaign.UpdatedBy = appcontext.ContextGetCustomerInfo(ctx).Sub
	created, err := h.useCase.Create(ctx, campaign)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusCreated, handler.BaseResponse[entity.PromotionCampaign]{
			Data: created,
		},
	)
}

// Update godoc
//
//	@Summary		Update promotion campaign
//	@Description	Update promotion campaign
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"id"
//	@Param			request	body		PromotionCampaign	        true	"request"
//	@Success		200		{object}	handler.BaseResponse[entity.PromotionCampaign]
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/promotion-campaigns/{id} [patch]
func (h *PromotionCampaignHandler) Update(ctx *gin.Context) {
	errorMessage := "update promotion campaign"
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	current, err := h.useCase.GetById(ctx, id)
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderError(ctx, err)
		return
	}
	if err := ctx.ShouldBindJSON(&current); err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	current.UpdatedBy = appcontext.ContextGetCustomerInfo(ctx).Sub
	current.UpdatedAt = time.Now()
	updated, err := h.useCase.Update(ctx, current)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.PromotionCampaign]{
			Data: updated,
		},
	)
}

// GetAll godoc
//
//	@Summary		Get all promotion campaigns
//	@Description	Get all promotion campaigns
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[[]entity.PromotionCampaign]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/promotion-campaigns [get]
func (h *PromotionCampaignHandler) GetAll(ctx *gin.Context) {
	req := entity.GetPromotionCampaignsRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		h.logger.Error("get all promotion campaigns", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid query")
	}
	res, err := h.useCase.GetAll(ctx, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]entity.PromotionCampaign]{
			Data: res,
		},
	)
}
