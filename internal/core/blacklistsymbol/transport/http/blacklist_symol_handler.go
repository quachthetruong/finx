package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/blacklistsymbol"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/handler"
)

type BlacklistSymbolHandler struct {
	handler.BaseHandler
	logger                 *slog.Logger
	blacklistSymbolUseCase blacklistsymbol.UseCase
}

func NewBlacklistSymbolHandler(baseHandler handler.BaseHandler, logger *slog.Logger,
	blacklistSymbolUseCase blacklistsymbol.UseCase,
) *BlacklistSymbolHandler {
	return &BlacklistSymbolHandler{
		BaseHandler:            baseHandler,
		logger:                 logger,
		blacklistSymbolUseCase: blacklistSymbolUseCase,
	}
}

func (h *BlacklistSymbolHandler) Create(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.logger.Error("symbol id", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	blacklistSymbol := entity.BlacklistSymbol{}
	if err := ctx.ShouldBindJSON(&blacklistSymbol); err != nil {
		h.logger.Error("create black list symbol", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	blacklistSymbol.SymbolId = id
	created, err := h.blacklistSymbolUseCase.Create(ctx, blacklistSymbol)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusCreated, handler.BaseResponse[entity.BlacklistSymbol]{
			Data: created,
		},
	)
}

func (h *BlacklistSymbolHandler) Update(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.logger.Error("symbol id", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	currentBlacklistSymbol, err := h.blacklistSymbolUseCase.GetById(ctx, id)
	if err != nil {
		h.logger.Error("blacklist symbol id not found", slog.String("err", err.Error()))
		h.RenderNotFound(ctx, fmt.Sprintf("blacklist symbol id: %v not found", id))
		return
	}
	if err := ctx.ShouldBindJSON(&currentBlacklistSymbol); err != nil {
		h.logger.Error("update blacklist symbol", slog.String("err", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	updated, err := h.blacklistSymbolUseCase.Update(ctx, currentBlacklistSymbol)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusCreated, handler.BaseResponse[entity.BlacklistSymbol]{
			Data: updated,
		},
	)
}
