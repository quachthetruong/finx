package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/promotion_loan_package"
	"financing-offer/internal/handler"
	"financing-offer/pkg/cache"
)

type PromotionLoanPackageHandler struct {
	handler.BaseHandler
	cacheStore cache.Cache
	useCase    promotionloanpackage.UseCase
}

func NewPromotionLoanPackageHandler(baseHandler handler.BaseHandler, cache cache.Cache, useCase promotionloanpackage.UseCase) *PromotionLoanPackageHandler {
	return &PromotionLoanPackageHandler{
		BaseHandler: baseHandler,
		cacheStore:  cache,
		useCase:     useCase,
	}
}

// GetBestLoanPackageIds godoc
//
//	@Summary		Get best loan package ids
//	@Description	Get best loan package ids
//	@Tags			configuration,investor
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[[]int64]
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-configurations/best-loan-package-ids [get]
func (h *PromotionLoanPackageHandler) GetBestLoanPackageIds(ctx *gin.Context) {
	res, err := h.useCase.GetOngoingPromotionLoanPackageIds(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]int64]{
			Data: res,
		},
	)
}

// GetPromotionLoanPackages godoc
//
//	@Summary		Get promotion loan packages
//	@Description	Get promotion loan packages
//	@Tags			configuration,investor
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[entity.PromotionLoanPackage]
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-configurations/promotion-loan-packages [get]
func (h *PromotionLoanPackageHandler) GetPromotionLoanPackages(ctx *gin.Context) {
	res, err := h.useCase.GetOnGoingPromotionLoanPackages(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.PromotionLoanPackage]{
			Data: res,
		},
	)
}

// SetPromotionLoanPackages godoc
//
//	@Summary		Set promotion loan packages
//	@Description	Set promotion loan packages
//	@Tags			configuration,admin
//	@Accept			json
//	@Produce		json
//	@Param			promotionLoanPackages	body		SetPromotionLoanPackagesRequest	true	"promotion loan packages"
//	@Success		200						{object}	handler.BaseResponse[entity.PromotionLoanPackage]
//	@Failure		400						{object}	handler.ErrorResponse
//	@Failure		500						{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/configurations/promotion-loan-packages [post]
func (h *PromotionLoanPackageHandler) SetPromotionLoanPackages(ctx *gin.Context) {
	newPromotionLoanPackagesRequest := SetPromotionLoanPackagesRequest{}
	err := ctx.ShouldBindJSON(&newPromotionLoanPackagesRequest)
	if err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	newPromotionLoanPackages := newPromotionLoanPackagesRequest.toEntity()
	loanPackage, err := h.useCase.SetPromotionLoanPackage(ctx, newPromotionLoanPackages, h.UserSubOrEmpty(ctx))
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.PromotionLoanPackage]{
			Data: loanPackage,
		},
	)
}

// GetPromotionLoanPackageBySymbol godoc
//
//	@Summary		Get promotion loan package by symbol
//	@Description	Get promotion loan package by symbol
//	@Tags			promotion,investor
//	@Accept			json
//	@Produce		json
//	@Param			symbol	query	string	true	"symbol"
//	@Param			custodyCode	query	string	true	"custody code"
//	@Param			accountNo	query	string	false	"account no"
//	@Success		200	{object}	GetPromotionLoanPackageBySymbolResponse
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/promotion-loan-packages/{symbol} [get]
func (h *PromotionLoanPackageHandler) GetPromotionLoanPackageBySymbol(ctx *gin.Context) {
	symbol, err := h.ParamNotEmpty(ctx, "symbol")
	if err != nil {
		h.RenderBadRequest(ctx, "symbol is required")
		return
	}
	investor, err := h.Investor(ctx)
	if err != nil {
		h.RenderUnauthenticated(ctx, "investor is required")
		return

	}
	req := GetPromotionLoanPackageBySymbolRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	res, err := h.useCase.GetPromotionLoanPackageBySymbol(ctx, symbol, req.AccountNo, investor.CustodyCode)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, GetPromotionLoanPackageBySymbolResponse{
			Symbol:      symbol,
			CustodyCode: investor.CustodyCode,
			AccountNos:  res,
		},
	)
}

// GetPublicPromotionLoanPackageBySymbol godoc
//
//	@Summary		Get public promotion loan package by symbol
//	@Description	Get public promotion loan package by symbol
//	@Tags			promotion,public
//	@Accept			json
//	@Produce		json
//	@Param			symbol	path	string	true	"symbol"
//	@Success		200	{object}	entity.PromotionLoanPackage
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Router			/v1/public/promotion-loan-packages/{symbol} [get]
func (h *PromotionLoanPackageHandler) GetPublicPromotionLoanPackageBySymbol(ctx *gin.Context) {
	symbol, err := h.ParamNotEmpty(ctx, "symbol")
	if err != nil {
		h.RenderBadRequest(ctx, "symbol is required")
		return
	}
	res, err := cache.Do[*entity.AccountLoanPackageWithSymbol](
		h.cacheStore,
		fmt.Sprintf("public-promotion-loan-package-%s", symbol),
		30*time.Second,
		func() (*entity.AccountLoanPackageWithSymbol, error) {
			return h.useCase.GetPublicPromotionLoanPackageBySymbol(ctx, symbol)
		},
	)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, res,
	)
}

// GetPublicPromotionLoanPackages godoc
//
//	@Summary		Get public promotion loan packages
//	@Description	Get public promotion loan packages
//	@Tags			promotion,public
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[[]entity.AccountLoanPackageWithSymbol]
//	@Failure		500	{object}	handler.ErrorResponse
//	@Router			public/v1/promotion-loan-packages [get]
func (h *PromotionLoanPackageHandler) GetPublicPromotionLoanPackages(ctx *gin.Context) {
	res, err := h.useCase.GetPublicPromotionLoanPackages(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]entity.AccountLoanPackageWithSymbol]{
			Data: res,
		},
	)
}

// GetInvestorPromotionLoanPackage godoc
//
//	@Summary		Get investor promotion loan package
//	@Description	Get investor promotion loan package
//	@Tags			promotion,investor
//	@Accept			json
//	@Produce		json
//	@Param			accountNo	query	string	true	"account no"
//	@Success		200	{object}	handler.BaseResponse[[]PromotionPackagesWithAccountNo]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/promotion-loan-packages [get]
func (h *PromotionLoanPackageHandler) GetInvestorPromotionLoanPackage(ctx *gin.Context) {
	investor, err := h.Investor(ctx)
	if err != nil {
		h.RenderUnauthenticated(ctx, "investor is required")
		return
	}
	req := GetInvestorPromotionLoanPackagesRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	res, err := h.useCase.GetInvestorPromotionLoanPackages(ctx, req.AccountNo, investor.CustodyCode)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	dataResponse := make([]PromotionPackagesWithAccountNo, 0, len(res))
	for accountNo, promotionPackages := range res {
		dataResponse = append(
			dataResponse, PromotionPackagesWithAccountNo{
				AccountNo: accountNo,
				Symbols:   promotionPackages,
			},
		)
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]PromotionPackagesWithAccountNo]{
			Data: dataResponse,
		},
	)
}

// GetInvestorPromotionLoanPackageV2 godoc
//
//	@Summary		Get investor promotion loan package
//	@Description	Get investor promotion loan package
//	@Tags			promotion,investor
//	@Accept			json
//	@Produce		json
//	@Param			accountNo	query	string	true	"account no"
//	@Param			symbol	    query	string	true	"symbol"
//	@Success		200	{object}	handler.BaseResponse[[]PromotionPackagesWithAccountNo]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v2/promotion-loan-packages [get]
func (h *PromotionLoanPackageHandler) GetInvestorPromotionLoanPackageV2(ctx *gin.Context) {
	investor, err := h.Investor(ctx)
	if err != nil {
		h.RenderUnauthenticated(ctx, "investor is required")
		return
	}
	req := GetPromotionLoanPackagesRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	res, err := h.useCase.GetPromotionLoanPackages(ctx, req.AccountNo, investor.CustodyCode, req.Symbol)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	dataResponse := make([]PromotionLoanPackages, 0, len(res))
	for accountNo, promotionPackages := range res {
		if len(promotionPackages) != 0 {
			dataResponse = append(
				dataResponse, PromotionLoanPackages{
					AccountNo:    accountNo,
					LoanPackages: promotionPackages,
				},
			)
		}
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]PromotionLoanPackages]{
			Data: dataResponse,
		},
	)
}

// GetPublicPromotionLoanPackagesV2 godoc
//
//	@Summary		Get public promotion loan packages
//	@Description	Get public promotion loan packages
//	@Tags			promotion,public
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[[]entity.AccountLoanPackageWithSymbol]
//	@Failure		500	{object}	handler.ErrorResponse
//	@Router			public/v2/promotion-loan-packages [get]
func (h *PromotionLoanPackageHandler) GetPublicPromotionLoanPackagesV2(ctx *gin.Context) {
	req := GetPublicLoanPackagesRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	res, err := h.useCase.GetPublicPromotionLoanPackagesWithCampaigns(ctx, req.Symbol)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]entity.LoanPackageWithCampaignProduct]{
			Data: res,
		},
	)
}
