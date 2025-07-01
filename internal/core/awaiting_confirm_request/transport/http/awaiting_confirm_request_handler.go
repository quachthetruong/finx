package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	awaitingconfirmrequest "financing-offer/internal/core/awaiting_confirm_request"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/handler"
)

type AwaitingConfirmRequestHandler struct {
	handler.BaseHandler
	useCase awaitingconfirmrequest.UseCase
	logger  *slog.Logger
}

// GetAll godoc
//
//	@Summary		Get all awaiting confirm request
//	@Description	Get all awaiting confirm request
//	@Tags			awaiting confirm request,admin
//	@Accept			json
//	@Produce		json
//
//	@Param			page[number]			query		int			false	"pageNumber"
//	@Param			page[size]				query		int			false	"pageSize"
//	@Param			symbols					query		[]string	false	"symbols"
//	@Param			flowTypes				query		[]string	false	"flowTypes"
//	@Param			accountNumbers			query		[]string	false	"accountNumbers"
//	@Param			investorId				query		string		false	"investorId"
//	@Param			ids						query		[]int64		false	"ids"
//	@Param			startDate				query		string		false	"startDate"
//	@Param			endDate					query		string		false	"endDate"
//	@Param			latestUpdateCategories	query		[]string	false	"latestUpdateCategories"
//	@Param			assetType				query		string		false	"assetType"
//	@Param			custodyCode				query		string		false	"custodyCode"
//	@Param			custodyCodes			query		[]string	false	"custodyCodes"
//
//	@Success		200						{object}	handler.ResponseWithPaging[[]entity.AwaitingConfirmRequest]
//	@Failure		500						{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/awaiting-confirm-requests [get]
func (h *AwaitingConfirmRequestHandler) GetAll(ctx *gin.Context) {
	req := GetAllAwaitingConfirmRequestRequest{}
	if err := h.ParseQueryWithPagination(ctx, &req.Paging, &req); err != nil {
		h.logger.Error("AwaitingConfirmRequestHandler GetAll", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	res, pagingMetaData, err := h.useCase.GetAll(ctx, req.toFilter())
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	ctx.JSON(
		http.StatusOK, handler.ResponseWithPaging[[]entity.AwaitingConfirmRequest]{
			Data:     res,
			MetaData: pagingMetaData,
		},
	)
}

// GetAllDerivative godoc
//
//	@Summary		Get all awaiting confirm request derivative
//	@Description	Get all awaiting confirm request derivative
//	@Tags			awaiting confirm request,admin
//	@Accept			json
//	@Produce		json
//
//	@Param			page[number]			query		int			false	"pageNumber"
//	@Param			page[size]				query		int			false	"pageSize"
//	@Param			symbols					query		[]string	false	"symbols"
//	@Param			flowTypes				query		[]string	false	"flowTypes"
//	@Param			accountNumbers			query		[]string	false	"accountNumbers"
//	@Param			investorId				query		string		false	"investorId"
//	@Param			ids						query		[]int64		false	"ids"
//	@Param			startDate				query		string		false	"startDate"
//	@Param			endDate					query		string		false	"endDate"
//	@Param			latestUpdateCategories	query		[]string	false	"latestUpdateCategories"
//	@Param			assetType				query		string		false	"assetType"
//	@Param			custodyCode				query		string		false	"custodyCode"
//	@Param			custodyCodes			query		[]string	false	"custodyCodes"
//
//	@Success		200						{object}	handler.ResponseWithPaging[[]entity.AwaitingConfirmRequest]
//	@Failure		500						{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/derivative-awaiting-confirm-requests [get]
func (h *AwaitingConfirmRequestHandler) GetAllDerivative(ctx *gin.Context) {
	req := GetAllAwaitingConfirmRequestRequest{}
	if err := h.ParseQueryWithPagination(ctx, &req.Paging, &req); err != nil {
		h.logger.Error("AwaitingConfirmRequestHandler GetAllDerivative", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	req.AssetType = entity.AssetTypeDerivative.String()
	res, pagingMetaData, err := h.useCase.GetAll(ctx, req.toFilter())
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	ctx.JSON(
		http.StatusOK, handler.ResponseWithPaging[[]entity.AwaitingConfirmRequest]{
			Data:     res,
			MetaData: pagingMetaData,
		},
	)
}

func NewAwaitingConfirmRequestHandler(bh handler.BaseHandler, logger *slog.Logger, useCase awaitingconfirmrequest.UseCase) *AwaitingConfirmRequestHandler {
	return &AwaitingConfirmRequestHandler{
		BaseHandler: bh,
		useCase:     useCase,
		logger:      logger,
	}
}
