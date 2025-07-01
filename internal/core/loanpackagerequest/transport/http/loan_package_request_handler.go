package http

import (
	"financing-offer/internal/core/submissionsheet"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/loanpackagerequest"
	"financing-offer/internal/core/scoregroupinterest"
	"financing-offer/internal/funcs"
	"financing-offer/internal/handler"
	"financing-offer/pkg/optional"
)

type LoanPackageRequestHandler struct {
	handler.BaseHandler
	logger                    *slog.Logger
	useCase                   loanpackagerequest.UseCase
	scoreGroupInterestUseCase scoregroupinterest.UseCase
	submissionSheetUseCase    submissionsheet.UseCase
}

func NewLoanPackageRequestHandler(
	baseHandler handler.BaseHandler,
	logger *slog.Logger,
	useCase loanpackagerequest.UseCase,
	scoreGroupInterestUseCase scoregroupinterest.UseCase,
	submissionSheetUseCase submissionsheet.UseCase,
) *LoanPackageRequestHandler {
	return &LoanPackageRequestHandler{
		BaseHandler:               baseHandler,
		logger:                    logger,
		useCase:                   useCase,
		scoreGroupInterestUseCase: scoreGroupInterestUseCase,
		submissionSheetUseCase:    submissionSheetUseCase,
	}
}

// admin handlers

// GetAll godoc
//
//	@Summary		Get all loan package requests
//	@Description	Get all loan package requests (admin)
//	@Tags			loan package request,admin
//	@Accept			json
//	@Produce		json
//	@Param			page[size]		query		int64			false	"pageSize"
//	@Param			page[number]	query		int64			false	"pageNumber"
//	@Param			sort			query		string			false	"sort"
//	@Param			symbols			query		[]string		false	"symbols"
//	@Param			types			query		[]string		false	"types"
//	@Param			investorId		query		string			false	"investorId"
//	@Param			statuses		query		[]string		false	"statuses"
//	@Param			ids				query		[]int64			false	"ids"
//	@Param			startDate		query		string			false	"startDate"
//	@Param			endDate			query		string			false	"endDate"
//	@Param			loanPercentFrom	query		decimal.Decimal	false	"loanPercentFrom"
//	@Param			loanPercentTo	query		decimal.Decimal	false	"loanPercentTo"
//	@Param			limitAmountFrom	query		int64			false	"limitAmountFrom"
//	@Param			limitAmountTo	query		int64			false	"limitAmountTo"
//	@Param			accountNumbers	query		[]string		false	"accountNumbers"
//	@Param			assetType		query		string			false	"assetType"
//	@Param			custodyCode		query		string			false	"custodyCode"
//	@Param			custodyCodes	query		[]string		false	"custodyCodes"
//	@Success		200				{object}	handler.ResponseWithPaging[[]entity.LoanPackageRequest]
//	@Success		400				{object}	handler.ErrorResponse
//	@Failure		500				{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/loan-package-requests [get]
func (h *LoanPackageRequestHandler) GetAll(ctx *gin.Context) {
	req := GetLoanPackageRequest{}
	if err := h.ParseQueryWithPagination(ctx, &req.Paging, &req); err != nil {
		h.RenderBadRequest(ctx, "parse query")
		return
	}
	res, pagingMeta, err := h.useCase.GetAll(ctx, req.toEntity())
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.ResponseWithPaging[[]entity.LoanPackageRequest]{
			Data:     res,
			MetaData: pagingMeta,
		},
	)
}

func (h *LoanPackageRequestHandler) GetAllUnderlyingRequests(ctx *gin.Context) {
	req := GetUnderlyingLoanPackageRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	res, err := h.useCase.GetAllUnderlyingRequests(ctx, req.toEntity())
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]entity.UnderlyingLoanPackageRequest]{
			Data: res,
		},
	)
}

// AdminGetById godoc
//
//	@Summary		Get loan package request by id
//	@Description	Get loan package request by id (admin)
//	@Tags			loan package request,admin
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	"id"
//	@Success		200	{object}	handler.BaseResponse[entity.LoanPackageRequest]
//	@Success		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/loan-package-requests/{id} [get]
func (h *LoanPackageRequestHandler) AdminGetById(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderIdInvalid(ctx)
		return
	}
	res, err := h.useCase.GetById(ctx, id, entity.LoanPackageFilter{})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.LoanPackageRequest]{Data: res})
}

// AdminConfirmUserRequest godoc
//
//	@Summary		Admin confirm user request
//	@Description	Admin confirm user request
//	@Tags			loan package request,admin
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int64								true	"id"
//	@Param			request	body		ConfirmLoanPackageRequestRequest	true	"body"
//	@Success		200		{object}	handler.BaseResponse[entity.LoanPackageRequest]
//	@Success		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/loan-package-requests/{id}/admin-confirm [post]
func (h *LoanPackageRequestHandler) AdminConfirmUserRequest(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	req := ConfirmLoanPackageRequestRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	confirmUser := req.OfferedBy
	if confirmUser == "" {
		confirmUser = h.UserSubOrEmpty(ctx)
	}
	res, err := h.useCase.AdminConfirmLoanRequest(ctx, id, confirmUser, req.LoanId)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.LoanPackageRequest]{Data: res})
}

// AdminCancelLoanRequest godoc
//
//	@Summary		Admin cancel loan request
//	@Description	Admin cancel loan request (admin)
//	@Tags			loan package request,admin
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int64							true	"id"
//	@Param			request	body		CancelLoanPackageRequestRequest	true	"body"
//	@Success		200		{object}	handler.BaseResponse[entity.LoanPackageRequest]
//	@Success		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/loan-package-requests/{id}/cancel [post]
func (h *LoanPackageRequestHandler) AdminCancelLoanRequest(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	req := CancelLoanPackageRequestRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	confirmUser := req.OfferedBy
	if confirmUser == "" {
		confirmUser = h.UserSubOrEmpty(ctx)
	}
	res, err := h.useCase.AdminCancelLoanRequest(ctx, id, confirmUser, funcs.UniqueElements(req.LoanIds))
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.LoanPackageRequest]{Data: res})
}

func (h *LoanPackageRequestHandler) GetAvailablePackages(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderIdInvalid(ctx)
		return
	}
	res, err := h.scoreGroupInterestUseCase.GetForLoanPackageRequest(ctx, id)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[[]entity.ScoreGroupInterest]{Data: res})
}

// investor handlers

// InvestorGetById godoc
//
//	@Summary		Get loan package request by id
//	@Description	Get loan package request by id (investor)
//	@Tags			loan package request,investor
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	"id"
//	@Success		200	{object}	handler.BaseResponse[entity.LoanPackageRequest]
//	@Success		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-loan-package-request/{id} [get]
func (h *LoanPackageRequestHandler) InvestorGetById(ctx *gin.Context) {
	investorId, err := h.InvestorId(ctx)
	if err != nil {
		h.RenderUnauthenticated(ctx, err.Error())
		return
	}
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderIdInvalid(ctx)
		return
	}
	res, err := h.useCase.GetById(
		ctx, id, entity.LoanPackageFilter{
			InvestorId: optional.Some(investorId),
		},
	)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.LoanPackageRequest]{Data: res})
}

// InvestorGetAll godoc
//
//	@Summary		Get all loan package requests
//	@Description	Get all loan package requests (investor)
//	@Tags			loan package request,investor
//	@Accept			json
//	@Produce		json
//	@Param			page[size]		query		int64			false	"pageSize"
//	@Param			page[number]	query		int64			false	"pageNumber"
//	@Param			symbols			query		[]string		false	"symbols"
//	@Param			types			query		[]string		false	"types"
//	@Param			statuses		query		[]string		false	"statuses"
//	@Param			ids				query		[]int64			false	"ids"
//	@Param			startDate		query		string			false	"startDate"
//	@Param			endDate			query		string			false	"endDate"
//	@Param			loanPercentFrom	query		decimal.Decimal	false	"loanPercentFrom"
//	@Param			loanPercentTo	query		decimal.Decimal	false	"loanPercentTo"
//	@Param			limitAmountFrom	query		int64			false	"limitAmountFrom"
//	@Param			limitAmountTo	query		int64			false	"limitAmountTo"
//	@Param			accountNumbers	query		[]string		false	"accountNumbers"
//	@Param			assetType		query		string			false	"assetType"
//	@Success		200				{object}	handler.BaseResponse[[]entity.LoanPackageRequest]
//	@Success		400				{object}	handler.ErrorResponse
//	@Failure		500				{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-loan-package-request [get]
func (h *LoanPackageRequestHandler) InvestorGetAll(ctx *gin.Context) {
	req := GetLoanPackageRequest{}
	if err := h.ParseQueryWithPagination(ctx, &req.Paging, &req); err != nil {
		h.RenderBadRequest(ctx, "parse query")
		return
	}
	investorId, err := h.InvestorId(ctx)
	if err != nil {
		h.RenderUnauthenticated(ctx, err.Error())
		return
	}
	req.InvestorId = investorId
	res, err := h.useCase.InvestorGetAll(ctx, req.toEntity())
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]entity.LoanPackageRequest]{
			Data: res,
		},
	)
}

// InvestorRequest godoc
//
//	@Summary		Investor request loan package
//	@Description	Investor request loan package (investor)
//	@Tags			loan package request,investor
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateLoanPackageRequestUnderlyingRequest	true	"body"
//	@Success		200		{object}	handler.BaseResponse[entity.LoanPackageRequest]
//	@Success		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-loan-package-request [post]
func (h *LoanPackageRequestHandler) InvestorRequest(ctx *gin.Context) {
	req := CreateLoanPackageRequestUnderlyingRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	investor, err := h.Investor(ctx)
	if err != nil {
		h.RenderUnauthenticated(ctx, err.Error())
		return
	}
	res, err := h.useCase.InvestorRequest(ctx, req.toEntity(investor.InvestorId), investor)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.LoanPackageRequest]{Data: res})
}

// InvestorRequestDerivative godoc
//
//	@Summary		Investor request loan package derivative
//	@Description	Investor request loan package derivative (investor)
//	@Tags			loan package request,investor
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateLoanPackageRequestDerivativeRequest	true	"body"
//	@Success		200		{object}	handler.BaseResponse[entity.LoanPackageRequest]
//	@Success		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-derivative-requests [post]
func (h *LoanPackageRequestHandler) InvestorRequestDerivative(ctx *gin.Context) {
	req := CreateLoanPackageRequestDerivativeRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	investorId, err := h.InvestorId(ctx)
	if err != nil {
		h.RenderUnauthenticated(ctx, err.Error())
		return
	}
	requestEntity, err := req.toEntity(investorId)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	res, err := h.useCase.InvestorRequestDerivative(ctx, requestEntity)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.LoanPackageRequest]{Data: res})
}

// SaveLoanRateExistedRequest godoc
//
//	@Summary		Save loan rate existed request
//	@Description	Save loan rate existed request (investor)
//	@Tags			loan package request,investor
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoggedRequestRequest	true	"body"
//	@Success		200		{object}	handler.BaseResponse[entity.LoggedRequest]
//	@Success		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-logged-requests [post]
func (h *LoanPackageRequestHandler) SaveLoanRateExistedRequest(ctx *gin.Context) {
	req := LoggedRequestRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	investorId, err := h.InvestorId(ctx)
	if err != nil {
		h.RenderUnauthenticated(ctx, err.Error())
		return
	}
	res, err := h.useCase.SaveExistedLoanRateRequest(ctx, investorId, req.Request.toEntity(investorId))
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.LoggedRequest]{Data: res})
}

// CancelAllLoanPackageRequestBySymbolId godoc
//
//	@Summary		Cancel all loan package request by symbol id
//	@Description	Cancel all loan package request by symbol id
//	@Tags			loan package request,admin
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	"symbolId"
//	@Success		200	{object}	handler.BaseResponse[[]entity.LoanPackageRequest]
//	@Success		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/symbols/{id}/cancel-requests [post]
func (h *LoanPackageRequestHandler) CancelAllLoanPackageRequestBySymbolId(ctx *gin.Context) {
	symbolId, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderIdInvalid(ctx)
		return
	}
	res, err := h.useCase.CancelAllLoanPackageRequestBySymbolId(
		ctx, symbolId, appcontext.ContextGetCustomerInfo(ctx).Sub,
	)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[[]entity.LoanPackageRequest]{Data: res})
}

func (h *LoanPackageRequestHandler) AdminConfirmWithNewLoanPackage(ctx *gin.Context) {
	req := entity.SubmissionSheetShorten{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	res, err := h.useCase.AdminSubmitSubmission(ctx, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.LoanPackageRequest]{Data: res})
}

func (h *LoanPackageRequestHandler) AdminDeclineLoanRequestWithNewLoanPackage(ctx *gin.Context) {
	req := entity.SubmissionSheetShorten{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	res, err := h.useCase.AdminSubmitSubmission(ctx, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.LoanPackageRequest]{Data: res})
}

func (h *LoanPackageRequestHandler) AdminGetLatestSubmissionSheet(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderIdInvalid(ctx)
		return
	}
	res, err := h.submissionSheetUseCase.GetLatestByRequestId(ctx, id)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.SubmissionSheet]{Data: res})
}
