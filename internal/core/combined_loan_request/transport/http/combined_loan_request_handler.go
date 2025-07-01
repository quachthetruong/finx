package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	combinedloanrequest "financing-offer/internal/core/combined_loan_request"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/handler"
)

type CombinedLoanRequestHandler struct {
	handler.BaseHandler
	logger  *slog.Logger
	useCase combinedloanrequest.UseCase
}

// GetAll godoc
//
//	@Summary		Get all combined loan request
//	@Description	Get all combined loan request
//	@Tags			combined loan request,admin
//	@Accept			json
//	@Produce		json
//
//	@Param			page[number]	query		int			false	"pageNumber"
//	@Param			page[size]		query		int			false	"pageSize"
//	@Param			symbols			query		[]string	false	"symbols"
//	@Param			startDate		query		string		false	"startDate"
//	@Param			endDate			query		string		false	"endDate"
//	@Param			offerDateFrom	query		string		false	"offerDateFrom"
//	@Param			offerDateTo		query		string		false	"offerDateTo"
//	@Param			flowTypes		query		[]string	false	"flowTypes"
//	@Param			accountNumbers	query		[]string	false	"accountNumbers"
//	@Param			investorId		query		string		false	"investorId"
//	@Param			status			query		string		false	"status"
//	@Param			assignedLoanId	query		int64		false	"assignedLoanId"
//	@Param			activatedLoanId	query		int64		false	"activatedLoanId"
//	@Param			ids				query		[]int64		false	"ids"
//	@Param			assetType		query		string		false	"assetType"
//	@Param			custodyCode		query		string		false	"custodyCode"
//
//	@Success		200				{object}	handler.ResponseWithPaging[[]entity.CombinedLoanRequest]
//	@Failure		500				{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/combined-requests [get]
func (h *CombinedLoanRequestHandler) GetAll(ctx *gin.Context) {
	req := GetAllCombinedRequestsRequest{}
	if err := h.ParseQueryWithPagination(ctx, &req.Paging, &req); err != nil {
		h.logger.Error("CombinedLoanRequestHandler GetAllWithFilter", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	res, pagingMetaData, err := h.useCase.GetAll(ctx, req.toFilter())
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	ctx.JSON(
		http.StatusOK, handler.ResponseWithPaging[[]entity.CombinedLoanRequest]{
			Data:     res,
			MetaData: pagingMetaData,
		},
	)
}

func NewCombinedLoanRequestHandler(bh handler.BaseHandler, logger *slog.Logger, useCase combinedloanrequest.UseCase) *CombinedLoanRequestHandler {
	return &CombinedLoanRequestHandler{
		BaseHandler: bh,
		useCase:     useCase,
		logger:      logger,
	}
}
