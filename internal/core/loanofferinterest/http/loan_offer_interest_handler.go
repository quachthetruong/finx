package loanOfferInterestHttp

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/loancontract"
	"financing-offer/internal/core/loanofferinterest"
	"financing-offer/internal/handler"
)

type LoanOfferInterestHandler struct {
	handler.BaseHandler
	logger              *slog.Logger
	useCase             loanofferinterest.UseCase
	loanContractUseCase loancontract.UseCase
}

// GetAllWithFilter godoc
//
//	@Summary		Get all loan offer interest with filter
//	@Description	Get all loan offer interest with filter
//	@Tags			loan package offer line,admin
//	@Accept			json
//	@Produce		json
//	@Param			page[number]	query		int		false	"pageNumber"
//	@Param			page[size]		query		int		false	"pageSize"
//	@Param			status			query		string	false	"status"
//	@Param			sort			query		string	false	"sort"
//	@Success		200				{object}	handler.ResponseWithPaging[[]entity.LoanPackageOfferInterest]
//	@Failure		400				{object}	handler.ErrorResponse
//	@Failure		500				{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/loan-offer-interests [get]
func (h *LoanOfferInterestHandler) GetAllWithFilter(c *gin.Context) {
	req := GetAllLoanOfferInterestRequest{}
	if err := h.ParseQueryWithPagination(c, &req.Paging, &req); err != nil {
		h.logger.Error("LoanOfferInterestHandler GetAllWithFilter", slog.String("error", err.Error()))
		h.RenderBadRequest(c, err.Error())
		return
	}
	res, pagingMetaData, err := h.useCase.GetAll(c, req.toFilter())
	if err != nil {
		h.RenderBadRequest(c, err.Error())
		return
	}
	c.JSON(
		http.StatusOK, handler.ResponseWithPaging[[]entity.LoanPackageOfferInterest]{
			Data:     res,
			MetaData: pagingMetaData,
		},
	)
}

// InvestorCancelLoanPackageOfferInterest godoc
//
//	@Summary		Investor cancel loan package offer interest
//	@Description	Investor cancel loan package offer interest
//	@Tags			loan package offer line,investor
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		200	{object}	handler.BaseResponse[string]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-loan-offer-interests/{id}/cancel [post]
func (h *LoanOfferInterestHandler) InvestorCancelLoanPackageOfferInterest(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	investorId, err := h.InvestorId(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	if err := h.useCase.InvestorCancelLoanPackageInterest(ctx, id, investorId); err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[string]{Data: "ok"})
}

// InvestorConfirmLoanPackageInterest godoc
//
//	@Summary		Investor confirm loan package offer interest
//	@Description	Investor confirm loan package offer interest
//	@Tags			loan package offer line,investor
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		200	{object}	handler.BaseResponse[string]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-loan-offer-interests/{id}/confirm [post]
func (h *LoanOfferInterestHandler) InvestorConfirmLoanPackageInterest(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	investorId, err := h.InvestorId(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	if err := h.useCase.InvestorConfirmLoanPackageInterest(ctx, []int64{id}, investorId); err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[string]{Data: "ok"})
}

// InvestorConfirmMultipleLoanPackageInterest godoc
//
//	@Summary		Investor confirm multiple loan package offer interest
//	@Description	Investor confirm multiple loan package offer interest
//	@Tags			loan package offer line,investor
//	@Accept			json
//	@Produce		json
//	@Param			loanPackageOfferInterestIds	body		InvestorConfirmLoanPackageOfferInterestRequest	true	"body"
//	@Success		200							{object}	handler.BaseResponse[string]
//	@Failure		400							{object}	handler.ErrorResponse
//	@Failure		500							{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-loan-offer-interests/confirm [post]
func (h *LoanOfferInterestHandler) InvestorConfirmMultipleLoanPackageInterest(ctx *gin.Context) {
	req := InvestorConfirmLoanPackageOfferInterestRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	investorId, err := h.InvestorId(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	if err := h.useCase.InvestorConfirmLoanPackageInterest(
		ctx, req.LoanPackageOfferInterestIds, investorId,
	); err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[string]{Data: "ok"})
}

func (h *LoanOfferInterestHandler) FillWithLoanPackageData(ctx *gin.Context) {
	total, err := h.useCase.SyncLoanPackageData(ctx)
	if err != nil {
		h.ReportError(ctx, err)
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[int]{Data: total})
}

// CreateAssignedLoanOfferInterestLoanContract godoc
//
//	@Summary		Create assigned loan offer interest loan contract
//	@Description	Create assigned loan offer interest loan contract
//	@Tags			loan package offer line,admin
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int													true	"id"
//	@Param			body	body		CreateAssignedLoanOfferInterestLoanContractRequest	true	"body"
//	@Success		200		{object}	handler.BaseResponse[entity.LoanContract]
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/loan-offer-interests/{id}/assign-loan-contract [post]
func (h *LoanOfferInterestHandler) CreateAssignedLoanOfferInterestLoanContract(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	req := CreateAssignedLoanOfferInterestLoanContractRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	res, err := h.useCase.CreateAssignedLoanOfferInterestLoanContract(
		ctx, id, req.LoanPackageAccountId, req.LoanProductIdRef, req.LoanPackage,
	)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.LoanContract]{Data: res})
}

func NewLoanOfferInterestHandler(
	baseHandler handler.BaseHandler,
	logger *slog.Logger,
	useCase loanofferinterest.UseCase,
) *LoanOfferInterestHandler {
	return &LoanOfferInterestHandler{
		BaseHandler: baseHandler,
		logger:      logger,
		useCase:     useCase,
	}
}
