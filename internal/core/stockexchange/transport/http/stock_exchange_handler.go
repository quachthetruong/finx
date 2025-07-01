package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/stockexchange"
	"financing-offer/internal/handler"
)

type StockExchangeHandler struct {
	handler.BaseHandler
	logger  *slog.Logger
	useCase stockexchange.UseCase
}

func NewStockExchangeHandler(baseHandler handler.BaseHandler, logger *slog.Logger, useCase stockexchange.UseCase) *StockExchangeHandler {
	return &StockExchangeHandler{BaseHandler: baseHandler, logger: logger, useCase: useCase}
}

// Create godoc
//
//	@Summary		Create stock exchange
//	@Description	Create stock exchange
//	@Tags			stock exchange,admin
//	@Accept			json
//	@Produce		json
//	@Param			stockExchange	body		StockExchangeRequest	true	"stock exchange"
//	@Success		201				{object}	handler.BaseResponse[entity.StockExchange]
//	@Failure		400				{object}	handler.ErrorResponse
//	@Failure		500				{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/stock-exchanges [post]
func (h *StockExchangeHandler) Create(ctx *gin.Context) {
	req := StockExchangeRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("create stock exchange", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	created, err := h.useCase.Create(ctx, req.toEntity(0))
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusCreated, handler.BaseResponse[entity.StockExchange]{
			Data: created,
		},
	)
}

// GetAll godoc
//
//	@Summary		Get all stock exchange
//	@Description	Get all stock exchange
//	@Tags			stock exchange,admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[[]entity.StockExchange]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/stock-exchanges [get]
func (h *StockExchangeHandler) GetAll(ctx *gin.Context) {
	stockExchanges, err := h.useCase.GetAll(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]entity.StockExchange]{
			Data: stockExchanges,
		},
	)
}

// Update godoc
//
//	@Summary		Update stock exchange
//	@Description	Update stock exchange
//	@Tags			stock exchange,admin
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int						true	"id"
//	@Param			stockExchange	body		entity.StockExchange	true	"stock exchange"
//	@Success		200				{object}	handler.BaseResponse[entity.StockExchange]
//	@Failure		400				{object}	handler.ErrorResponse
//	@Failure		500				{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/stock-exchanges/{id} [patch]
func (h *StockExchangeHandler) Update(ctx *gin.Context) {
	errorMessage := "update stock exchange"
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	oldStockExchange, err := h.useCase.GetById(ctx, id)
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderError(ctx, err)
		return
	}
	if err := ctx.ShouldBindJSON(&oldStockExchange); err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	updated, err := h.useCase.Update(ctx, oldStockExchange)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.StockExchange]{
			Data: updated,
		},
	)
}

// Delete godoc
//
//	@Summary		Delete stock exchange
//	@Description	Delete stock exchange
//	@Tags			stock exchange,admin
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		204	{object}	handler.BaseResponse[string]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/stock-exchanges/{id} [delete]
func (h *StockExchangeHandler) Delete(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.logger.Error("delete stock exchange", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	if err := h.useCase.Delete(ctx, id); err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusNoContent, handler.BaseResponse[string]{
			Data: "ok",
		},
	)
}
