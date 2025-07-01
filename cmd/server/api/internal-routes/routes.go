package internal_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"

	"financing-offer/cmd/server/middlewares"
	investorHttp "financing-offer/internal/core/investor/transport/http"
	loanOfferHttp "financing-offer/internal/core/loanoffer/transport/http"
	loanOfferInterestHttp "financing-offer/internal/core/loanofferinterest/http"
)

func NewRoutes(group *gin.RouterGroup, middleware middlewares.Middleware, injector *do.Injector) {
	internalRoutes := group.Group("/internal")
	internalRoutes.Use(middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"))
	loanOfferHandler := do.MustInvoke[*loanOfferHttp.LoanPackageOfferHandler](injector)
	loanOfferInterestHandler := do.MustInvoke[*loanOfferInterestHttp.LoanOfferInterestHandler](injector)
	investorHandler := do.MustInvoke[*investorHttp.InvestorHandler](injector)

	groupLoanPackageOffer := internalRoutes.Group("/loan-package-offers")
	groupLoanPackageOffer.GET("/expire", loanOfferHandler.ManualTriggerExpireLoanOffers)
	groupLoanPackageOfferInterest := internalRoutes.Group("/loan-package-offer-interests")
	groupLoanPackageOfferInterest.GET("/sync-loan-package-data", loanOfferInterestHandler.FillWithLoanPackageData)

	groupInvestor := internalRoutes.Group("/investors")
	groupInvestor.POST("/sync-investor-data", investorHandler.FillInvestorIdsFromRequests)
}
