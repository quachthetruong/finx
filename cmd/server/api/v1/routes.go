package v1

import (
	submissionDefaultHttp "financing-offer/internal/core/submission_default/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"

	"financing-offer/cmd/server/middlewares"
	configHttp "financing-offer/internal/config/transport/http"
	"financing-offer/internal/core/awaiting_confirm_request/transport/http"
	blacklistSymbolHttp "financing-offer/internal/core/blacklistsymbol/transport/http"
	combinedRequestHttp "financing-offer/internal/core/combined_loan_request/transport/http"
	configurationHttp "financing-offer/internal/core/configuration/transport/http"
	financialProductHttp "financing-offer/internal/core/financialproduct/transport/http"
	investorAccountHttp "financing-offer/internal/core/investor_account/transport/http"
	loanOfferHttp "financing-offer/internal/core/loanoffer/transport/http"
	loanOfferInterestHttp "financing-offer/internal/core/loanofferinterest/http"
	loanPackageRequestHttp "financing-offer/internal/core/loanpackagerequest/transport/http"
	loanPolicyTemplateHttp "financing-offer/internal/core/loanpolicytemplate/transport/http"
	promotionCampaignHttp "financing-offer/internal/core/promotion_campaign/transport/http"
	promotionLoanPackageHttp "financing-offer/internal/core/promotion_loan_package/transport/http"
	schedulerHttp "financing-offer/internal/core/scheduler/transport/http"
	scoreGroupHttp "financing-offer/internal/core/scoregroup/transport/http"
	scoreGroupInterestHttp "financing-offer/internal/core/scoregroupinterest/transport/http"
	stockExchangeHttp "financing-offer/internal/core/stockexchange/transport/http"
	submissionSheetHttp "financing-offer/internal/core/submissionsheet/transport/http"
	suggestedOfferHttp "financing-offer/internal/core/suggested_offer/transport/http"
	suggestedOfferConfigHttp "financing-offer/internal/core/suggested_offer_config/transport/http"
	symbolHttp "financing-offer/internal/core/symbol/transport/http"
	symbolScoreHttp "financing-offer/internal/core/symbolscore/transport/http"
	featureHttp "financing-offer/internal/featureflag/transport/http"
)

func NewRoutes(engine *gin.RouterGroup, middleware middlewares.Middleware, injector *do.Injector) {
	symbolHandler := do.MustInvoke[*symbolHttp.SymbolHandler](injector)
	blacklistSymbolHandler := do.MustInvoke[*blacklistSymbolHttp.BlacklistSymbolHandler](injector)
	stockExchangeHandler := do.MustInvoke[*stockExchangeHttp.StockExchangeHandler](injector)
	symbolScoreHandler := do.MustInvoke[*symbolScoreHttp.SymbolScoreHandler](injector)
	loanPackageRequestHandler := do.MustInvoke[*loanPackageRequestHttp.LoanPackageRequestHandler](injector)
	scoreGroupHandler := do.MustInvoke[*scoreGroupHttp.ScoreGroupHandler](injector)
	scoreGroupInterestHandler := do.MustInvoke[*scoreGroupInterestHttp.ScoreGroupInterestHandler](injector)
	loanOfferHandler := do.MustInvoke[*loanOfferHttp.LoanPackageOfferHandler](injector)
	offerInterestHandler := do.MustInvoke[*loanOfferInterestHttp.LoanOfferInterestHandler](injector)
	featureHandler := do.MustInvoke[*featureHttp.FeatureHandler](injector)
	configHandler := do.MustInvoke[*configHttp.ConfigHandler](injector)
	loanRequestSchedulerConfigHandler := do.MustInvoke[*schedulerHttp.SchedulerHandler](injector)
	awaitingConfirmRequestHandler := do.MustInvoke[*http.AwaitingConfirmRequestHandler](injector)
	combinedRequestHandler := do.MustInvoke[*combinedRequestHttp.CombinedLoanRequestHandler](injector)
	investorAccountHandler := do.MustInvoke[*investorAccountHttp.InvestorAccountHandler](injector)
	loanPolicyTemplateHandler := do.MustInvoke[*loanPolicyTemplateHttp.LoanPolicyTemplateHandler](injector)
	financialProductHandler := do.MustInvoke[*financialProductHttp.FinancialProductHandler](injector)
	submissionSheetHandler := do.MustInvoke[*submissionSheetHttp.SubmissionSheetHandler](injector)
	suggestedOfferConfigHandler := do.MustInvoke[*suggestedOfferConfigHttp.SuggestedOfferConfigHandler](injector)
	suggestedOfferHandler := do.MustInvoke[*suggestedOfferHttp.SuggestedOfferHandler](injector)
	promotionLoanPackageHandler := do.MustInvoke[*promotionLoanPackageHttp.PromotionLoanPackageHandler](injector)
	configurationHandler := do.MustInvoke[*configurationHttp.ConfigurationHandler](injector)
	submissionDefaultHandler := do.MustInvoke[*submissionDefaultHttp.SubmissionDefaultHandler](injector)
	promotionCampaignHandler := do.MustInvoke[*promotionCampaignHttp.PromotionCampaignHandler](injector)

	v1Routes := engine.Group("/v1")
	v2Routes := engine.Group("/v2")

	groupSymbol := v1Routes.Group("/symbols", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"))
	groupSymbol.GET("", symbolHandler.GetAll)
	groupSymbol.GET("/:id", symbolHandler.GetById)
	groupSymbol.POST("", symbolHandler.Create)
	groupSymbol.PATCH("/:id", symbolHandler.Update)
	groupSymbol.POST("/:id/blacklist-symbols", blacklistSymbolHandler.Create)
	groupSymbol.POST("/:id/cancel-requests", loanPackageRequestHandler.CancelAllLoanPackageRequestBySymbolId)

	groupStockExchange := v1Routes.Group("/stock-exchanges", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"))
	groupStockExchange.GET("", stockExchangeHandler.GetAll)
	groupStockExchange.POST("", stockExchangeHandler.Create)
	groupStockExchange.PATCH("/:id", stockExchangeHandler.Update)
	groupStockExchange.DELETE("/:id", stockExchangeHandler.Delete)

	groupSymbolScore := v1Routes.Group("/symbol-scores", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"))
	groupSymbolScore.POST("", symbolScoreHandler.Create)
	groupSymbolScore.PATCH("/:id", symbolScoreHandler.Update)

	groupAdminLoanPackageRequest := v1Routes.Group(
		"/loan-package-requests", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"),
	)
	groupAdminLoanPackageRequest.GET("", loanPackageRequestHandler.GetAll)
	groupAdminLoanPackageRequest.GET("/:id", loanPackageRequestHandler.AdminGetById)
	groupAdminLoanPackageRequest.POST("/:id/admin-confirm", loanPackageRequestHandler.AdminConfirmUserRequest)
	groupAdminLoanPackageRequest.POST("/:id/cancel", loanPackageRequestHandler.AdminCancelLoanRequest)
	groupAdminLoanPackageRequest.GET("/:id/available-packages", loanPackageRequestHandler.GetAvailablePackages)
	groupAdminLoanPackageRequest.POST("/:id/submissions", loanPackageRequestHandler.AdminConfirmWithNewLoanPackage)
	groupAdminLoanPackageRequest.POST(
		"/:id/cancel-with-submission", loanPackageRequestHandler.AdminDeclineLoanRequestWithNewLoanPackage,
	)

	groupAdminLoanPackageRequest.GET("/:id/latest-submission", loanPackageRequestHandler.AdminGetLatestSubmissionSheet)

	groupAdminLoanPackageRequest.GET("/underlying", loanPackageRequestHandler.GetAllUnderlyingRequests)

	scoreGroup := v1Routes.Group("/score-groups", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"))
	scoreGroup.GET("", scoreGroupHandler.GetAll)
	scoreGroup.POST("", scoreGroupHandler.Create)
	scoreGroup.GET("/:id/available-packages", scoreGroupHandler.GetAvailablePackages)
	scoreGroup.PATCH("/:id", scoreGroupHandler.Update)
	scoreGroup.DELETE("/:id", scoreGroupHandler.Delete)

	groupScoreGroupInterest := v1Routes.Group(
		"/score-group-interests", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"),
	)
	groupScoreGroupInterest.GET("", scoreGroupInterestHandler.GetAll)
	groupScoreGroupInterest.GET("/:id", scoreGroupInterestHandler.GetById)
	groupScoreGroupInterest.POST("", scoreGroupInterestHandler.Create)
	groupScoreGroupInterest.PATCH("/:id", scoreGroupInterestHandler.Update)
	groupScoreGroupInterest.DELETE("/:id", scoreGroupInterestHandler.Delete)

	groupLoanOfferInterest := v1Routes.Group(
		"/loan-offer-interests", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"),
	)
	groupLoanOfferInterest.GET("", offerInterestHandler.GetAllWithFilter)
	groupLoanOfferInterest.POST(
		"/:id/assign-loan-contract", offerInterestHandler.CreateAssignedLoanOfferInterestLoanContract,
	)

	offlineUpdatesWithIdUri := "/:id/offline-updates"
	groupLoanOffer := v1Routes.Group("/loan-package-offers", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"))
	groupLoanOffer.GET(offlineUpdatesWithIdUri, loanOfferHandler.GetOfflineOfferUpdateHistory)
	groupLoanOffer.POST(offlineUpdatesWithIdUri, loanOfferHandler.CreateOfflineOfferUpdate)
	groupLoanOffer.POST("/:id/assign-loan", loanOfferHandler.AdminAssignLoanId)
	groupLoanOffer.POST("/:id/cancel", loanOfferHandler.AdminCancelLoanPackageOfferInterest)

	groupDerivativeLoanOffer := v1Routes.Group(
		"/derivative-loan-package-offers", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"),
	)
	groupDerivativeLoanOffer.GET(offlineUpdatesWithIdUri, loanOfferHandler.GetDerivativeOfflineOfferUpdateHistory)
	groupDerivativeLoanOffer.POST(offlineUpdatesWithIdUri, loanOfferHandler.CreateDerivativeOfflineOfferUpdate)
	groupDerivativeLoanOffer.POST("/:id/assign-loan", loanOfferHandler.AdminAssignLoanId)

	groupAwaitingConfirmRequest := v1Routes.Group(
		"/awaiting-confirm-requests", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"),
	)
	groupAwaitingConfirmRequest.GET("", awaitingConfirmRequestHandler.GetAll)

	groupDerivativeAwaitingConfirmRequest := v1Routes.Group(
		"/derivative-awaiting-confirm-requests", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"),
	)

	groupDerivativeAwaitingConfirmRequest.GET("", awaitingConfirmRequestHandler.GetAllDerivative)

	groupCombinedRequest := v1Routes.Group(
		"combined-requests", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"),
	)
	groupCombinedRequest.GET("", combinedRequestHandler.GetAll)

	groupAdminConfiguration := v1Routes.Group(
		"/configurations", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"),
	)
	groupAdminConfiguration.POST("/promotion-loan-packages", promotionLoanPackageHandler.SetPromotionLoanPackages)
	groupAdminConfiguration.POST("/loan-rate", configurationHandler.SetLoanRate)
	groupAdminConfiguration.GET("/loan-rate", configurationHandler.GetLoanRate)
	groupAdminConfiguration.POST("/margin-pool", configurationHandler.SetMarginPool)
	groupAdminConfiguration.GET("/margin-pool", configurationHandler.GetMarginPool)

	// investor routes

	groupInvestorSymbol := v1Routes.Group("/my-symbols", middleware.RequireAuthenticatedUser())
	groupInvestorSymbol.GET("/:symbol", symbolHandler.GetBySymbol)

	groupInvestorLoanPackageRequest := v1Routes.Group(
		"/my-loan-package-request", middleware.RequireAuthenticatedUser(),
		middleware.RequireFeatureEnable("loanRequest"),
	)
	groupInvestorLoanPackageRequest.GET("", loanPackageRequestHandler.InvestorGetAll)
	groupInvestorLoanPackageRequest.GET("/:id", loanPackageRequestHandler.InvestorGetById)
	groupInvestorLoanPackageRequest.POST("", loanPackageRequestHandler.InvestorRequest)

	groupInvestorDerivativeRequest := v1Routes.Group(
		"/my-derivative-requests", middleware.RequireAuthenticatedUser(),
	)
	groupInvestorDerivativeRequest.POST("", loanPackageRequestHandler.InvestorRequestDerivative)

	groupInvestorLoggedRequest := v1Routes.Group("/my-logged-requests", middleware.RequireAuthenticatedUser())
	groupInvestorLoggedRequest.POST("", loanPackageRequestHandler.SaveLoanRateExistedRequest)

	groupInvestorLoanPackageOffer := v1Routes.Group("/my-loan-package-offer", middleware.RequireAuthenticatedUser())
	groupInvestorLoanPackageOffer.GET("", loanOfferHandler.InvestorFindLoanOffers)
	groupInvestorLoanPackageOffer.POST("/:id/cancel", loanOfferHandler.InvestorCancelLoanPackageOffer)
	groupInvestorLoanPackageOffer.GET("/:id", loanOfferHandler.InvestorGetById)

	groupOfferInterest := v1Routes.Group("/my-loan-offer-interests", middleware.RequireAuthenticatedUser())
	groupOfferInterest.POST(
		"/confirm", middleware.RequireHOActive(), offerInterestHandler.InvestorConfirmMultipleLoanPackageInterest,
	)
	groupOfferInterest.POST(
		"/:id/confirm", middleware.RequireHOActive(), offerInterestHandler.InvestorConfirmLoanPackageInterest,
	)
	groupOfferInterest.POST("/:id/cancel", offerInterestHandler.InvestorCancelLoanPackageOfferInterest)

	groupFeature := v1Routes.Group("/features")
	groupFeature.GET("/:name/verify", featureHandler.CheckFeatureEnable)

	groupBlacklistSymbol := v1Routes.Group(
		"/blacklist-symbols",
		middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"),
	)
	groupBlacklistSymbol.PATCH("/:id", blacklistSymbolHandler.Update)

	groupUserConfig := v1Routes.Group("/my-configurations", middleware.RequireAuthenticatedUser())
	groupUserConfig.GET("", configHandler.GetConfiguration)
	groupUserConfig.GET("/best-loan-package-ids", promotionLoanPackageHandler.GetBestLoanPackageIds)
	groupUserConfig.GET("/promotion-loan-packages", promotionLoanPackageHandler.GetPromotionLoanPackages)

	groupPromotionLoanPackage := v1Routes.Group("/promotion-loan-packages", middleware.RequireAuthenticatedUser())
	groupPromotionLoanPackage.GET("", promotionLoanPackageHandler.GetInvestorPromotionLoanPackage)
	groupPromotionLoanPackage.GET(":symbol", promotionLoanPackageHandler.GetPromotionLoanPackageBySymbol)

	groupPromotionLoanPackageV2 := v2Routes.Group("/promotion-loan-packages", middleware.RequireAuthenticatedUser())
	groupPromotionLoanPackageV2.GET("", promotionLoanPackageHandler.GetInvestorPromotionLoanPackageV2)

	groupLoanOfferSymbol := v1Routes.Group("/loan-offerable-symbols", middleware.RequireAuthenticatedUser())
	groupLoanOfferSymbol.GET("/:symbol-code", symbolHandler.GetSymbolNotActiveBlacklist)

	loanPolicyTemplate := v1Routes.Group(
		"/loan-policy-template", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"),
	)
	loanPolicyTemplate.GET("", loanPolicyTemplateHandler.GetAll)
	loanPolicyTemplate.POST("", loanPolicyTemplateHandler.Create)
	loanPolicyTemplate.GET("/:id", loanPolicyTemplateHandler.GetById)
	loanPolicyTemplate.PUT("/:id", loanPolicyTemplateHandler.Update)
	loanPolicyTemplate.DELETE("/:id", loanPolicyTemplateHandler.Delete)

	marginOperationGroup := v1Routes.Group("mo")
	marginOperationGroup.GET("/applicable-loan-rates", financialProductHandler.GetLoanRates)
	marginOperationGroup.GET("/applicable-margin-pools", financialProductHandler.GetMarginPools)

	schedulerGroup := v1Routes.Group(
		"/schedulers",
		middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"),
	)
	schedulerGroup.GET(
		"/loan-request-scheduler-config", loanRequestSchedulerConfigHandler.GetAllLoanRequestSchedulerConfigs,
	)
	schedulerGroup.GET(
		"/loan-request-scheduler-config/current",
		loanRequestSchedulerConfigHandler.GetCurrentLoanRequestSchedulerConfig,
	)
	schedulerGroup.POST(
		"/loan-request-scheduler-config", loanRequestSchedulerConfigHandler.CreateLoanRequestSchedulerConfig,
	)

	groupInvestorAccount := v1Routes.Group(
		"/investor-accounts", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"),
	)
	groupInvestorAccount.PUT(
		"/:account-no/margin-status", investorAccountHandler.VerifyAndUpdateInvestorAccountMarginStatus,
	)

	submissionSheetGroup := v1Routes.Group(
		"/submission-sheets", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"),
	)
	submissionSheetGroup.POST("", submissionSheetHandler.Upsert)
	submissionSheetGroup.POST("/:id/approve", submissionSheetHandler.AdminApproveSubmissionSheet)
	submissionSheetGroup.POST("/:id/reject", submissionSheetHandler.AdminRejectSubmissionSheet)

	groupSuggestedOfferConfig := v1Routes.Group(
		"/suggested-offer-configs", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"),
	)
	groupSuggestedOfferConfig.POST("", suggestedOfferConfigHandler.Create)
	groupSuggestedOfferConfig.PATCH("/:id", suggestedOfferConfigHandler.Update)
	groupSuggestedOfferConfig.PATCH("/:id/status", suggestedOfferConfigHandler.UpdateStatus)
	groupSuggestedOfferConfig.GET("", suggestedOfferConfigHandler.GetAll)
	groupSuggestedOfferConfig.GET("/:id", suggestedOfferConfigHandler.GetById)

	groupUserSuggestedOfferConfig := v1Routes.Group(
		"/active-suggested-offer-config", middleware.RequireAuthenticatedUser(),
	)
	groupUserSuggestedOfferConfig.GET("", suggestedOfferConfigHandler.Get)

	groupSuggestedOffer := v1Routes.Group("/suggested-offers")
	groupSuggestedOffer.POST("", suggestedOfferHandler.CreateOffer)

	groupSubmissionDefault := v1Routes.Group("/submission-defaults", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"))
	groupSubmissionDefault.GET("", submissionDefaultHandler.GetSubmissionDefault)
	groupSubmissionDefault.POST("", submissionDefaultHandler.SetSubmissionDefault)

	promotionCampaign := v1Routes.Group("/promotion-campaigns", middleware.RequireOneOfRoles("ADMIN", "FINANCIAL_ADMIN"))
	promotionCampaign.GET("", promotionCampaignHandler.GetAll)
	promotionCampaign.POST("", promotionCampaignHandler.Create)
	promotionCampaign.PATCH("/:id", promotionCampaignHandler.Update)

	groupUserPromotionCampaignPackage := v1Routes.Group("/my-promotion-campaigns", middleware.RequireAuthenticatedUser())
	groupUserPromotionCampaignPackage.GET("", promotionCampaignHandler.GetAll)

}
