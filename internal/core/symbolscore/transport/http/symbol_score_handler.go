package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/symbolscore"
	"financing-offer/internal/handler"
)

type SymbolScoreHandler struct {
	handler.BaseHandler
	logger  *slog.Logger
	useCase symbolscore.UseCase
}

func NewSymbolScoreHandler(baseHandler handler.BaseHandler, logger *slog.Logger, useCase symbolscore.UseCase) *SymbolScoreHandler {
	return &SymbolScoreHandler{BaseHandler: baseHandler, logger: logger, useCase: useCase}
}

// Create godoc
//
//	@Summary		Create symbol score
//	@Description	Create symbol score
//	@Tags			symbol score,admin
//	@Accept			json
//	@Produce		json
//	@Param			symbolScore	body		CreateSymbolScoreRequest	true	"symbol score"
//	@Success		201			{object}	handler.BaseResponse[entity.SymbolScore]
//	@Failure		400			{object}	handler.ErrorResponse
//	@Failure		500			{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/symbol-scores [post]
func (h *SymbolScoreHandler) Create(ctx *gin.Context) {
	var req CreateSymbolScoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("create symbol score", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	e := req.toEntity(0)
	e.Creator = h.UserSubOrEmpty(ctx)
	created, err := h.useCase.Create(ctx, e)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.SymbolScore]{
			Data: created,
		},
	)
}

// Update godoc
//
//	@Summary		Update symbol score
//	@Description	Update symbol score
//	@Tags			symbol score,admin
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int					true	"id"
//	@Param			symbolScore	body		entity.SymbolScore	true	"symbol score"
//	@Success		200			{object}	handler.BaseResponse[entity.SymbolScore]
//	@Failure		400			{object}	handler.ErrorResponse
//	@Failure		500			{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/symbol-scores/{id} [patch]
func (h *SymbolScoreHandler) Update(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.logger.Error("update symbol score", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	oldSymbolScore, err := h.useCase.GetById(ctx, id)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	if err := ctx.ShouldBindJSON(&oldSymbolScore); err != nil {
		h.logger.Error("update symbol score", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	updated, err := h.useCase.Update(ctx, oldSymbolScore)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.SymbolScore]{
			Data: updated,
		},
	)
}
