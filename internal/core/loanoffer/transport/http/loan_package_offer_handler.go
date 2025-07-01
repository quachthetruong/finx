package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/loanoffer"
	"financing-offer/internal/core/loanofferinterest"
	offlineofferupdate "financing-offer/internal/core/offline_offer_update"
	"financing-offer/internal/handler"
)

const invalidInvestorId = "invalid investorId"

type LoanPackageOfferHandler struct {
	handler.BaseHandler
	logger                    *slog.Logger
	useCase                   loanoffer.UseCase
	offlineOfferUpdateUseCase offlineofferupdate.UseCase
	offerInterestUseCase      loanofferinterest.UseCase
}

func NewLoanPackageOfferHandler(
	baseHandler handler.BaseHandler,
	logger *slog.Logger,
	useCase loanoffer.UseCase,
	offlineOfferUpdateUseCase offlineofferupdate.UseCase,
	offerInterestUseCase loanofferinterest.UseCase,
) *LoanPackageOfferHandler {
	return &LoanPackageOfferHandler{
		BaseHandler:               baseHandler,
		logger:                    logger,
		useCase:                   useCase,
		offlineOfferUpdateUseCase: offlineOfferUpdateUseCase,
		offerInterestUseCase:      offerInterestUseCase,
	}
}

// InvestorFindLoanOffers godoc
//
//	@Summary		Investor find loan offers
//	@Description	Investor find loan offers
//	@Tags			loan package offer,investor
//	@Accept			json
//	@Produce		json
//	@Param			page[number]	query		int		false	"pageNumber"
//	@Param			page[size]		query		int		false	"pageSize"
//	@Param			status			query		string	false	"status"
//	@Param			symbol			query		string	false	"symbol"
//	@Param			assetType		query		string	false	"assetType"
//	@Param			sort			query		string	false	"sort"
//	@Success		200				{object}	handler.BaseResponse[[]entity.LoanPackageOffer]
//	@Failure		400				{object}	handler.ErrorResponse
//	@Failure		500				{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-loan-package-offer [get]
func (h *LoanPackageOfferHandler) InvestorFindLoanOffers(ctx *gin.Context) {
	investorId, err := h.InvestorId(ctx)
	if err != nil {
		h.RenderUnauthenticated(ctx, invalidInvestorId)
		return
	}
	req := GetLoanPackageOfferRequest{}
	if err := h.ParseQueryWithPagination(ctx, &req.Paging, &req); err != nil {
		h.logger.Error("find loan offers", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "parse query")
		return
	}
	filter := req.toFilter()
	filter.InvestorId = investorId
	res, err := h.useCase.FindAllForInvestor(ctx, filter)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]entity.LoanPackageOffer]{
			Data: res,
		},
	)
}

// InvestorCancelLoanPackageOffer godoc
//
//	@Summary		Investor cancel loan package offer
//	@Description	Investor cancel loan package offer
//	@Tags			loan package offer,investor
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		200	{object}	handler.BaseResponse[string]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-loan-package-offer/{id}/cancel [post]
func (h *LoanPackageOfferHandler) InvestorCancelLoanPackageOffer(ctx *gin.Context) {
	investorId, err := h.InvestorId(ctx)
	if err != nil {
		h.RenderUnauthenticated(ctx, invalidInvestorId)
		return
	}
	loanOfferId, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	if err := h.useCase.InvestorCancel(ctx, investorId, loanOfferId); err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[string]{Data: "ok"})
}

// InvestorGetById godoc
//
//	@Summary		Investor get loan package offer by id
//	@Description	Investor get loan package offer by id
//	@Tags			loan package offer,investor
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		200	{object}	handler.BaseResponse[entity.LoanPackageOffer]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-loan-package-offer/{id} [get]
func (h *LoanPackageOfferHandler) InvestorGetById(ctx *gin.Context) {
	investorId, err := h.InvestorId(ctx)
	if err != nil {
		h.RenderUnauthenticated(ctx, invalidInvestorId)
		return
	}
	loanOfferId, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	res, err := h.useCase.InvestorGetById(ctx, loanOfferId, investorId)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.LoanPackageOffer]{
			Data: res,
		},
	)
}

// ManualTriggerExpireLoanOffers godoc
//
//	@Summary		Manual trigger expire loan offers
//	@Description	Manual trigger expire loan offers
//	@Tags			loan package offer,admin,internal
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[string]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/internal/loan-package-offers/expire [get]
func (h *LoanPackageOfferHandler) ManualTriggerExpireLoanOffers(ctx *gin.Context) {
	if err := h.useCase.ExpireLoanOffers(ctx); err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[string]{Data: "ok"})
}

// GetOfflineOfferUpdateHistory godoc
//
//	@Summary		Get offline offer update history
//	@Description	Get offline offer update history
//	@Tags			loan package offer,admin
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		200	{object}	handler.BaseResponse[[]entity.OfflineOfferUpdate]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/loan-package-offers/{id}/offline-updates [get]
func (h *LoanPackageOfferHandler) GetOfflineOfferUpdateHistory(ctx *gin.Context) {
	loanOfferId, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	res, err := h.offlineOfferUpdateUseCase.GetByOfferId(ctx, loanOfferId, entity.AssetTypeUnderlying)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]entity.OfflineOfferUpdate]{
			Data: res,
		},
	)
}

func (h *LoanPackageOfferHandler) GetDerivativeOfflineOfferUpdateHistory(ctx *gin.Context) {
	loanOfferId, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	res, err := h.offlineOfferUpdateUseCase.GetByOfferId(ctx, loanOfferId, entity.AssetTypeDerivative)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]entity.OfflineOfferUpdate]{
			Data: res,
		},
	)
}

// CreateOfflineOfferUpdate godoc
//
//	@Summary		Create offline offer update
//	@Description	Create offline offer update
//	@Tags			loan package offer,admin
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int							true	"id"
//	@Param			offerUpdate	body		CreateOfferUpdateRequest	true	"body"
//	@Success		200			{object}	handler.BaseResponse[entity.OfflineOfferUpdate]
//	@Failure		400			{object}	handler.ErrorResponse
//	@Failure		500			{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/loan-package-offers/{id}/offline-updates [post]
func (h *LoanPackageOfferHandler) CreateOfflineOfferUpdate(ctx *gin.Context) {
	loanOfferId, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	req := CreateOfferUpdateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	res, err := h.offlineOfferUpdateUseCase.Create(
		ctx, entity.OfflineOfferUpdate{
			OfferId:   loanOfferId,
			Status:    entity.OfflineOfferUpdateStatus(req.Status),
			Category:  req.Category,
			Note:      req.Note,
			CreatedBy: h.UserSubOrEmpty(ctx),
		},
		entity.AssetTypeUnderlying,
	)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.OfflineOfferUpdate]{Data: res})
}

func (h *LoanPackageOfferHandler) CreateDerivativeOfflineOfferUpdate(ctx *gin.Context) {
	loanOfferId, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	req := CreateOfferUpdateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	res, err := h.offlineOfferUpdateUseCase.Create(
		ctx, entity.OfflineOfferUpdate{
			OfferId:   loanOfferId,
			Status:    entity.OfflineOfferUpdateStatus(req.Status),
			Category:  req.Category,
			Note:      req.Note,
			CreatedBy: h.UserSubOrEmpty(ctx),
		},
		entity.AssetTypeDerivative,
	)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.OfflineOfferUpdate]{Data: res})
}

// AdminCancelLoanPackageOfferInterest godoc
//
//	@Summary		Admin cancel loan package offer interest
//	@Description	Admin cancel loan package offer interest
//	@Tags			loan package offer,admin
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		200	{object}	handler.BaseResponse[string]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/loan-package-offers/{id}/cancel [post]
func (h *LoanPackageOfferHandler) AdminCancelLoanPackageOfferInterest(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}

	if err := h.offerInterestUseCase.AdminCancelLoanPackageInterestByOfferId(
		ctx, id, h.UserSubOrEmpty(ctx),
	); err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[string]{Data: "ok"})
}

// AdminAssignLoanId godoc
//
//	@Summary		Admin assign loan id
//	@Description	Admin assign loan id
//	@Tags			loan package offer,admin
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"id"
//	@Param			body	body		AdminAssignLoanIdRequest	true	"body"
//	@Success		200		{object}	handler.BaseResponse[string]
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/loan-package-offers/{id}/assign-loan [post]
func (h *LoanPackageOfferHandler) AdminAssignLoanId(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	req := AdminAssignLoanIdRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.RenderBadRequest(ctx, err.Error())
		return
	}
	if err := h.offerInterestUseCase.AdminAssignLoanIdByOfferId(ctx, id, req.LoanId); err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[string]{Data: "ok"})
}
