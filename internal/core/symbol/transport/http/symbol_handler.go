package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/symbol"
	"financing-offer/internal/handler"
)

type SymbolHandler struct {
	handler.BaseHandler
	logger        *slog.Logger
	symbolUseCase symbol.UseCase
}

func NewSymbolHandler(baseHandler handler.BaseHandler, logger *slog.Logger, symbolUseCase symbol.UseCase) *SymbolHandler {
	return &SymbolHandler{BaseHandler: baseHandler, logger: logger, symbolUseCase: symbolUseCase}
}

// Create godoc
//
//	@Summary		Create symbol
//	@Description	Create symbol
//	@Tags			symbol,admin
//	@Accept			json
//	@Produce		json
//	@Param			symbol	body		SymbolRequest	true	"symbol"
//	@Success		201		{object}	handler.BaseResponse[entity.Symbol]
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/symbols [post]
func (h *SymbolHandler) Create(ctx *gin.Context) {
	req := &SymbolRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.logger.Error("create symbol", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	created, err := h.symbolUseCase.Create(ctx, req.toEntity(0))
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusCreated, handler.BaseResponse[entity.Symbol]{
			Data: created,
		},
	)
}

// Update godoc
//
//	@Summary		Update symbol
//	@Description	Update symbol status
//	@Tags			symbol,admin
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"id"
//	@Param			status	body		UpdateSymbolStatusRequest	true	"status"
//	@Success		200		{object}	handler.BaseResponse[entity.Symbol]
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/symbols/{id} [patch]
func (h *SymbolHandler) Update(ctx *gin.Context) {
	errorMessage := "update symbol"
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	oldSymbol, err := h.symbolUseCase.GetById(ctx, id)
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderError(ctx, err)
		return
	}
	var updateSymbolRequest UpdateSymbolStatusRequest
	if err := ctx.ShouldBindJSON(&updateSymbolRequest); err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	oldSymbol.Status = updateSymbolRequest.Status
	oldSymbol.LastUpdatedBy = appcontext.ContextGetCustomerInfo(ctx).Sub
	updated, err := h.symbolUseCase.UpdateStatus(ctx, oldSymbol)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.Symbol]{
			Data: updated,
		},
	)
}

// GetBySymbol godoc
//
//	@Summary		Get symbol by symbol
//	@Description	Get symbol by symbol
//	@Tags			symbol,investor
//	@Accept			json
//	@Produce		json
//	@Param			symbol	path		string	true	"symbol"
//	@Success		200		{object}	handler.BaseResponse[entity.Symbol]
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		404		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-symbols/{symbol} [get]
func (h *SymbolHandler) GetBySymbol(ctx *gin.Context) {
	s := ctx.Param("symbol")
	res, err := h.symbolUseCase.GetBySymbol(ctx, s)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.Symbol]{
			Data: res,
		},
	)
}

// GetById godoc
//
//	@Summary		Get symbol by id
//	@Description	Get symbol by id
//	@Tags			symbol,admin
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		200	{object}	handler.BaseResponse[entity.Symbol]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		404	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/symbols/{id} [get]
func (h *SymbolHandler) GetById(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.logger.Error("get by id", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	res, err := h.symbolUseCase.GetById(ctx, id)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.Symbol]{
			Data: res,
		},
	)
}

// GetAll godoc
//
//	@Summary		Get all symbols
//	@Description	Get all symbols
//	@Tags			symbol,admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.ResponseWithPaging[[]entity.Symbol]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/symbols [get]
func (h *SymbolHandler) GetAll(ctx *gin.Context) {
	req := GetSymbolsRequest{}
	if err := h.ParseQueryWithPagination(ctx, &req.Paging, &req); err != nil {
		h.logger.Error("get all symbols", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "parse query")
		return
	}
	res, meta, err := h.symbolUseCase.GetAll(ctx, req.toFilter())
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.ResponseWithPaging[[]entity.Symbol]{
			Data:     res,
			MetaData: meta,
		},
	)
}

// GetSymbolNotActiveBlacklist godoc
//
//	@Summary		Get symbol not active blacklist
//	@Description	Get symbol not active blacklist
//	@Tags			symbol,admin
//	@Accept			json
//	@Produce		json
//	@Param			symbol-code	path		string	true	"symbol-code"
//	@Success		200			{object}	handler.BaseResponse[entity.Symbol]
//	@Failure		400			{object}	handler.ErrorResponse
//	@Failure		500			{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/loan-offerable-symbols/{symbol-code} [get]
func (h *SymbolHandler) GetSymbolNotActiveBlacklist(ctx *gin.Context) {
	symbolCode := ctx.Param("symbol-code")
	res, err := h.symbolUseCase.GetSymbolNotInActiveBlacklist(ctx, symbolCode)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.Symbol]{
			Data: res,
		},
	)
}
