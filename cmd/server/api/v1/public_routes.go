package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"

	"financing-offer/cmd/server/middlewares"
	promotionLoanPackageHttp "financing-offer/internal/core/promotion_loan_package/transport/http"
)

func NewPublicRoutes(engine *gin.RouterGroup, _ middlewares.Middleware, injector *do.Injector) {
	promotionLoanPackageHandler := do.MustInvoke[*promotionLoanPackageHttp.PromotionLoanPackageHandler](injector)
	v1Public := engine.Group("/v1")
	groupPublicPromotionLoanPackage := v1Public.Group("/promotion-loan-packages")
	groupPublicPromotionLoanPackage.GET("", promotionLoanPackageHandler.GetPublicPromotionLoanPackages)
	groupPublicPromotionLoanPackage.GET(":symbol", promotionLoanPackageHandler.GetPublicPromotionLoanPackageBySymbol)

	v2Public := engine.Group("/v2")
	groupPublicPromotionLoanPackageV2 := v2Public.Group("/promotion-loan-packages")
	groupPublicPromotionLoanPackageV2.GET("", promotionLoanPackageHandler.GetPublicPromotionLoanPackagesV2)
}
