package di

import (
	"context"
	"financing-offer/internal/core/configuration"
	"financing-offer/internal/core/marginoperation"
	"financing-offer/internal/core/promotion_campaign"
	"financing-offer/internal/core/submission_default"
	"log/slog"

	"github.com/samber/do"
	"go.temporal.io/sdk/client"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/apperrors/repository"
	"financing-offer/internal/atomicity"
	"financing-offer/internal/config"
	configRepo "financing-offer/internal/config/repository"
	configPostgres "financing-offer/internal/config/repository/postgres"
	configHttp "financing-offer/internal/config/transport/http"
	awaitingconfirmrequest "financing-offer/internal/core/awaiting_confirm_request"
	awaitingConfirmRequestRepo "financing-offer/internal/core/awaiting_confirm_request/repository"
	awaitingConfirmRequestPostgres "financing-offer/internal/core/awaiting_confirm_request/repository/postgres"
	awaitingConfirmRequestHttp "financing-offer/internal/core/awaiting_confirm_request/transport/http"
	"financing-offer/internal/core/blacklistsymbol"
	blSymbolPostgres "financing-offer/internal/core/blacklistsymbol/repository/postgres"
	blSymbolHttp "financing-offer/internal/core/blacklistsymbol/transport/http"
	combinedloanrequest "financing-offer/internal/core/combined_loan_request"
	combinedRequestRepo "financing-offer/internal/core/combined_loan_request/repository"
	combinedRequestPostgres "financing-offer/internal/core/combined_loan_request/repository/postgres"
	combinedRequestHttp "financing-offer/internal/core/combined_loan_request/transport/http"
	configurationHttp "financing-offer/internal/core/configuration/transport/http"
	financialProductDomain "financing-offer/internal/core/financialproduct"
	financialProductRepo "financing-offer/internal/core/financialproduct/repository"
	financialProductHttp "financing-offer/internal/core/financialproduct/transport/http"
	financingApiRepository "financing-offer/internal/core/financing/repository"
	flexOpenApiRepo "financing-offer/internal/core/flex/repository"
	"financing-offer/internal/core/investor"
	investorRepo "financing-offer/internal/core/investor/repository"
	investorPostgres "financing-offer/internal/core/investor/repository/postgres"
	investorHttp "financing-offer/internal/core/investor/transport/http"
	"financing-offer/internal/core/investor_account"
	investorAccountRepo "financing-offer/internal/core/investor_account/repository"
	investorAccountPostgres "financing-offer/internal/core/investor_account/repository/postgres"
	investorAccountHttp "financing-offer/internal/core/investor_account/transport/http"
	"financing-offer/internal/core/loancontract"
	loanContractPostgres "financing-offer/internal/core/loancontract/repository/postgres"
	loanContractHttp "financing-offer/internal/core/loancontract/transport/http"
	"financing-offer/internal/core/loanoffer"
	loanPackageOfferPostgres "financing-offer/internal/core/loanoffer/repository/postgres"
	loanOfferHttp "financing-offer/internal/core/loanoffer/transport/http"
	loanOfferScheduler "financing-offer/internal/core/loanoffer/transport/scheduler"
	"financing-offer/internal/core/loanofferinterest"
	loanOfferInterestHttp "financing-offer/internal/core/loanofferinterest/http"
	loanOfferInterestRepo "financing-offer/internal/core/loanofferinterest/repository"
	loanOfferInterestKafka "financing-offer/internal/core/loanofferinterest/repository/kafka"
	loanPackageOfferInterestPostgres "financing-offer/internal/core/loanofferinterest/repository/postgres"
	"financing-offer/internal/core/loanpackagerequest"
	loanPackageRequestRepo "financing-offer/internal/core/loanpackagerequest/repository"
	loanRequestKafka "financing-offer/internal/core/loanpackagerequest/repository/kafka"
	loanPackageRequestPostgres "financing-offer/internal/core/loanpackagerequest/repository/postgres"
	loanPackageRequestHttp "financing-offer/internal/core/loanpackagerequest/transport/http"
	loanPackageScheduler "financing-offer/internal/core/loanpackagerequest/transport/scheduler"
	"financing-offer/internal/core/loanpolicytemplate"
	loanPolicyTemplatePostgres "financing-offer/internal/core/loanpolicytemplate/repository/postgres"
	loanPolicyTemplateHttp "financing-offer/internal/core/loanpolicytemplate/transport/http"
	marginOperationRepo "financing-offer/internal/core/marginoperation/repository"
	odooServiceRepo "financing-offer/internal/core/odoo_service/repository"
	offlineofferupdate "financing-offer/internal/core/offline_offer_update"
	offlineOfferRepo "financing-offer/internal/core/offline_offer_update/repository"
	offlineOfferPosgres "financing-offer/internal/core/offline_offer_update/repository/postgres"
	orderServiceRepo "financing-offer/internal/core/orderservice/repository"
	promotionCampaignPostgres "financing-offer/internal/core/promotion_campaign/repository/postgres"
	promotionCampaignHttp "financing-offer/internal/core/promotion_campaign/transport/http"
	promotionloanpackage "financing-offer/internal/core/promotion_loan_package"
	promotionLoanPackageHttp "financing-offer/internal/core/promotion_loan_package/transport/http"
	"financing-offer/internal/core/scheduler"
	schedulerRepo "financing-offer/internal/core/scheduler/repository"
	schedulerRepoPostgres "financing-offer/internal/core/scheduler/repository/postgres"
	schedulerHttp "financing-offer/internal/core/scheduler/transport/http"
	"financing-offer/internal/core/scoregroup"
	scoreGroupPostgres "financing-offer/internal/core/scoregroup/repository/postgres"
	scoreGroupHttp "financing-offer/internal/core/scoregroup/transport/http"
	"financing-offer/internal/core/scoregroupinterest"
	scoreGroupInterestPostgres "financing-offer/internal/core/scoregroupinterest/repository/postgres"
	scoreGroupInterestHttp "financing-offer/internal/core/scoregroupinterest/transport/http"
	"financing-offer/internal/core/stockexchange"
	stockExchangePostgres "financing-offer/internal/core/stockexchange/repository/postgres"
	"financing-offer/internal/core/stockexchange/transport/http"
	submissionDefaultHttp "financing-offer/internal/core/submission_default/transport/http"
	"financing-offer/internal/core/submissionsheet"
	submissionSheetPostgres "financing-offer/internal/core/submissionsheet/repository/postgres"
	submissionSheetHttp "financing-offer/internal/core/submissionsheet/transport/http"
	suggestedOffer "financing-offer/internal/core/suggested_offer"
	suggestedOfferRepo "financing-offer/internal/core/suggested_offer/repository"
	suggestedOfferKafka "financing-offer/internal/core/suggested_offer/repository/kafka"
	suggestedOfferPostgres "financing-offer/internal/core/suggested_offer/repository/postgres"
	suggestedOfferHttp "financing-offer/internal/core/suggested_offer/transport/http"
	suggestedOfferConfig "financing-offer/internal/core/suggested_offer_config"
	suggestedOfferConfigRepo "financing-offer/internal/core/suggested_offer_config/repository"
	suggestedOfferConfigPostgres "financing-offer/internal/core/suggested_offer_config/repository/postgres"
	suggestedOfferConfigHttp "financing-offer/internal/core/suggested_offer_config/transport/http"
	"financing-offer/internal/core/symbol"
	symbolPostgres "financing-offer/internal/core/symbol/repository/postgres"
	symbolHttp "financing-offer/internal/core/symbol/transport/http"
	"financing-offer/internal/core/symbolscore"
	symbolScorePostgres "financing-offer/internal/core/symbolscore/repository/postgres"
	symbolScoreHttp "financing-offer/internal/core/symbolscore/transport/http"
	"financing-offer/internal/database"
	"financing-offer/internal/event"
	"financing-offer/internal/featureflag"
	http2 "financing-offer/internal/featureflag/transport/http"
	"financing-offer/internal/handler"
	"financing-offer/pkg/cache"
	"financing-offer/pkg/environment"
	"financing-offer/pkg/infra/financialproduct"
	financingApi "financing-offer/pkg/infra/financing-api"
	flexApi "financing-offer/pkg/infra/flex-api"
	mo_service "financing-offer/pkg/infra/mo-service"
	odooService "financing-offer/pkg/infra/odoo_service"
	orderService "financing-offer/pkg/infra/order-service"
	"financing-offer/pkg/mattermost"
	"financing-offer/pkg/shutdown"
)

func NewInjector(logger *slog.Logger) *do.Injector {
	injector := do.New()
	do.ProvideValue(injector, logger)
	do.Provide(injector, NewPublisher)
	do.Provide(injector, NewMattermostClient)
	do.Provide(injector, NewErrorService)
	do.Provide(injector, NewFinancialProductClient)
	do.Provide(injector, NewMoServiceClient)
	do.Provide(injector, NewFinancingApiClient)
	do.Provide(injector, NewOrderServiceClient)
	do.Provide(injector, NewTemporalClient)
	do.Provide(injector, NewFlexOpenApiClient)
	do.Provide(injector, NewOdooServiceClient)
	do.Provide(injector, NewCache)

	do.Provide(injector, NewBlackListRepository)
	do.Provide(injector, NewStockExchangeRepository)
	do.Provide(injector, NewSymbolRepository)
	do.Provide(injector, NewSymbolScoreRepository)
	do.Provide(injector, NewScoreGroupRepository)
	do.Provide(injector, NewScoreGroupInterestRepository)
	do.Provide(injector, NewLoanPackageRequestRepository)
	do.Provide(injector, NewLoanPackageOfferRepository)
	do.Provide(injector, NewLoanPackageOfferInterestRepository)
	do.Provide(injector, NewLoanContractRepository)
	do.Provide(injector, NewLoanRequestSchedulerConfigRepository)
	do.Provide(injector, NewSchedulerJobRepository)
	do.Provide(injector, NewOfflineOfferUpdateRepository)
	do.Provide(injector, NewPromotionCampaignRepository)
	do.Provide(injector, NewAwaitingConfirmRequestRepository)
	do.Provide(injector, NewCombinedRequestRepository)
	do.Provide(injector, NewInvestorRepository)
	do.Provide(injector, NewInvestorAccountRepository)
	do.Provide(injector, NewLoanPolicyTemplateRepository)
	do.Provide(injector, NewSubmissionSheetRepository)
	do.Provide(injector, NewSuggestedOfferConfigRepository)
	do.Provide(injector, NewSuggestedOfferRepository)

	do.Provide(injector, NewLoanOfferRequestEventPublisher)
	do.Provide(injector, NewLoanOfferInterestEventPublisher)
	do.Provide(injector, NewConfigurationPersistenceRepository)
	do.Provide(injector, NewSuggestedOfferEventPublisher)

	do.Provide(injector, NewBlackListUseCase)
	do.Provide(injector, NewStockExchangeUseCase)
	do.Provide(injector, NewSymbolUseCase)
	do.Provide(injector, NewSymbolScoreUseCase)
	do.Provide(injector, NewScoreGroupUseCase)
	do.Provide(injector, NewScoreGroupInterestUseCase)
	do.Provide(injector, NewLoanPackageRequestUseCase)
	do.Provide(injector, NewLoanPackageOfferUseCase)
	do.Provide(injector, NewLoanOfferInterestUseCase)
	do.Provide(injector, NewLoanContractUseCase)
	do.Provide(injector, NewFeatureUseCase)
	do.Provide(injector, NewConfigUseCase)
	do.Provide(injector, NewSchedulerUseCase)
	do.Provide(injector, NewOfflineOfferUpdateUseCase)
	do.Provide(injector, NewAwaitingConfirmRequestUseCase)
	do.Provide(injector, NewCombinedRequestUseCase)
	do.Provide(injector, NewInvestorUseCase)
	do.Provide(injector, NewInvestorAccountUseCase)
	do.Provide(injector, NewLoanPolicyTemplateUseCase)
	do.Provide(injector, NewFinancialProductUseCase)
	do.Provide(injector, NewSubmissionSheetUseCase)
	do.Provide(injector, NewSuggestedOfferConfigUseCase)
	do.Provide(injector, NewSuggestedOfferUseCase)
	do.Provide(injector, NewPromotionLoanPackageUseCase)
	do.Provide(injector, NewConfigurationUseCase)
	do.Provide(injector, NewMarginOperationUseCase)
	do.Provide(injector, NewSubmissionDefaultUseCase)
	do.Provide(injector, NewPromotionCampaignUseCase)

	do.Provide(injector, NewBaseHandler)
	do.Provide(injector, NewBlackListHandler)
	do.Provide(injector, NewStockExchangeHandler)
	do.Provide(injector, NewSymbolScoreHandler)
	do.Provide(injector, NewSymbolHandler)
	do.Provide(injector, NewScoreGroupHandler)
	do.Provide(injector, NewScoreGroupInterestHandler)
	do.Provide(injector, NewLoanPackageRequestHandler)
	do.Provide(injector, NewLoanPackageOfferHandler)
	do.Provide(injector, NewLoanContractHandler)
	do.Provide(injector, NewLoanPackageOfferInterestHandler)
	do.Provide(injector, NewFeatureHandler)
	do.Provide(injector, NewConfigHandler)
	do.Provide(injector, NewSchedulerHandler)
	do.Provide(injector, NewAwaitingConfirmRequestHandler)
	do.Provide(injector, NewCombinedRequestHandler)
	do.Provide(injector, NewInvestorHandler)
	do.Provide(injector, NewInvestorAccountHandler)
	do.Provide(injector, NewLoanPolicyTemplateHandler)
	do.Provide(injector, NewFinancialProductHandler)
	do.Provide(injector, NewSuggestedOfferConfigHandler)
	do.Provide(injector, NewSuggestedOfferHandler)

	do.Provide(injector, NewLoanOfferScheduler)
	do.Provide(injector, NewLoanPackageRequestScheduler)
	do.Provide(injector, NewSubmissionSheetHandler)
	do.Provide(injector, NewPromotionLoanPackageHandler)
	do.Provide(injector, NewConfigurationHandler)
	do.Provide(injector, NewSubmissionDefaultHandler)
	do.Provide(injector, NewPromotionCampaignHandler)
	return injector
}

func NewMattermostClient(i *do.Injector) (repository.NotifyWebhookRepository, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	env := do.MustInvoke[environment.Environment](i)
	return mattermost.NewClient(env, cfg.Mattermost.WebhookUrl), nil
}

func NewPublisher(i *do.Injector) (event.Publisher, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	tasks := do.MustInvoke[*shutdown.Tasks](i)
	return event.NewPublisher(cfg, tasks), nil
}

func NewErrorService(i *do.Injector) (apperrors.Service, error) {
	logger := do.MustInvoke[*slog.Logger](i)
	mattermostClient := do.MustInvoke[repository.NotifyWebhookRepository](i)
	return apperrors.NewService(mattermostClient, logger), nil
}

func NewFinancialProductClient(i *do.Injector) (financialProductRepo.FinancialProductRepository, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	baseClient := financialproduct.NewClient(cfg.FinancialProduct)
	cacheStore := do.MustInvoke[cache.Cache](i)
	return financialproduct.NewCachedFinancialProductRepository(baseClient, cacheStore), nil
}

func NewMoServiceClient(i *do.Injector) (marginOperationRepo.MarginOperationRepository, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	return mo_service.NewClient(cfg.MoService), nil
}

func NewFinancingApiClient(i *do.Injector) (financingApiRepository.FinancingRepository, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	return financingApi.NewClient(cfg.FinancingApi), nil
}

func NewOrderServiceClient(i *do.Injector) (orderServiceRepo.OrderServiceRepository, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	return orderService.NewClient(cfg.OrderService), nil
}

func NewOdooServiceClient(i *do.Injector) (odooServiceRepo.OdooServiceRepository, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	if cfg.Env == "prod" {
		return odooService.NewClient(cfg.OdooService, nil)
	}
	return odooService.NewInMemoryClient(cfg.OdooService)
}

func NewFlexOpenApiClient(i *do.Injector) (flexOpenApiRepo.FlexOpenApiRepository, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	return flexApi.NewClient(cfg.FlexOpenApi), nil
}

func NewSchedulerJobRepository(i *do.Injector) (schedulerRepo.SchedulerJobRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return schedulerRepoPostgres.NewSchedulerJobRepository(getDbFunc), nil
}

func NewBaseHandler(i *do.Injector) (handler.BaseHandler, error) {
	logger := do.MustInvoke[*slog.Logger](i)
	errorService := do.MustInvoke[apperrors.Service](i)
	return handler.NewBaseHandler(logger, errorService), nil
}

func NewBlackListRepository(i *do.Injector) (*blSymbolPostgres.BlackListSymbolRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return blSymbolPostgres.NewBlackListSymbolRepository(getDbFunc), nil
}

func NewStockExchangeRepository(i *do.Injector) (*stockExchangePostgres.StockExchangeRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return stockExchangePostgres.NewStockExchangeRepository(getDbFunc), nil
}

func NewSymbolScoreRepository(i *do.Injector) (*symbolScorePostgres.SymbolScoreRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return symbolScorePostgres.NewSymbolScoreRepository(getDbFunc), nil
}

func NewSymbolRepository(i *do.Injector) (*symbolPostgres.SymbolRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return symbolPostgres.NewSymbolRepository(getDbFunc), nil
}

func NewScoreGroupRepository(i *do.Injector) (*scoreGroupPostgres.ScoreGroupRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return scoreGroupPostgres.NewScoreGroupRepository(getDbFunc), nil
}

func NewScoreGroupInterestRepository(i *do.Injector) (*scoreGroupInterestPostgres.ScoreGroupInterestSqlRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return scoreGroupInterestPostgres.NewScoreGroupInterestSqlRepository(getDbFunc), nil
}

func NewLoanPackageRequestRepository(i *do.Injector) (*loanPackageRequestPostgres.LoanPackageRequestPostgresRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return loanPackageRequestPostgres.NewLoanPackageRequestPostgresRepository(getDbFunc), nil
}

func NewLoanPackageOfferRepository(i *do.Injector) (*loanPackageOfferPostgres.LoanPackageOfferPostgresRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return loanPackageOfferPostgres.NewLoanPackageOfferPostgresRepository(getDbFunc), nil
}

func NewLoanPackageOfferInterestRepository(i *do.Injector) (*loanPackageOfferInterestPostgres.LoanPackageOfferInterestPostgresRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return loanPackageOfferInterestPostgres.NewLoanPackageOfferInterestPostgresRepository(getDbFunc), nil
}

func NewLoanContractRepository(i *do.Injector) (*loanContractPostgres.LoanContractRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return loanContractPostgres.NewLoanContractRepository(getDbFunc), nil
}

func NewOfflineOfferUpdateRepository(i *do.Injector) (offlineOfferRepo.OfflineOfferUpdatePersistenceRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return offlineOfferPosgres.NewOfflineOfferUpdatePostgresRepository(getDbFunc), nil
}

func NewPromotionCampaignRepository(i *do.Injector) (*promotionCampaignPostgres.PromotionCampaignPostgresRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return promotionCampaignPostgres.NewPromotionCampaignRepository(getDbFunc), nil
}

func NewLoanOfferInterestEventPublisher(i *do.Injector) (loanOfferInterestRepo.LoanPackageOfferInterestEventRepository, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	publisher := do.MustInvoke[event.Publisher](i)
	temporalClient := do.MustInvoke[client.Client](i)
	return loanOfferInterestKafka.NewLoanOfferInterestEventPublisher(cfg.Kafka, publisher, temporalClient), nil
}

func NewLoanOfferRequestEventPublisher(i *do.Injector) (loanPackageRequestRepo.LoanPackageRequestEventRepository, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	publisher := do.MustInvoke[event.Publisher](i)
	return loanRequestKafka.NewLoanPackageRequestEventPublisher(cfg.Kafka, publisher), nil
}

func NewSuggestedOfferEventPublisher(i *do.Injector) (suggestedOfferRepo.SuggestedOfferEventRepository, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	publisher := do.MustInvoke[event.Publisher](i)
	return suggestedOfferKafka.NewSuggestedOfferEventPublisher(cfg.Kafka, publisher), nil
}

func NewLoanRequestSchedulerConfigRepository(i *do.Injector) (schedulerRepo.LoanRequestSchedulerConfigRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return schedulerRepoPostgres.NewLoanRequestSchedulerConfigRepo(getDbFunc), nil
}

func NewAwaitingConfirmRequestRepository(i *do.Injector) (awaitingConfirmRequestRepo.AwaitingConfirmRequestPersistenceRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return awaitingConfirmRequestPostgres.NewAwaitingConfirmRequestPostgresRepository(getDbFunc), nil
}

func NewCombinedRequestRepository(i *do.Injector) (combinedRequestRepo.CombinedLoanPackageRequestPersistenceRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return combinedRequestPostgres.NewCombinedLoanPackageRequestPostgresRepository(getDbFunc), nil
}

func NewInvestorRepository(i *do.Injector) (investorRepo.InvestorPersistenceRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return investorPostgres.NewInvestorPostgresRepository(getDbFunc), nil
}

func NewInvestorAccountRepository(i *do.Injector) (investorAccountRepo.InvestorAccountRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return investorAccountPostgres.NewInvestorAccountPostgresRepository(getDbFunc), nil
}

func NewSuggestedOfferConfigRepository(i *do.Injector) (suggestedOfferConfigRepo.SuggestedOfferConfigRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return suggestedOfferConfigPostgres.NewSuggestedOfferConfigRepository(getDbFunc), nil
}

func NewSuggestedOfferRepository(i *do.Injector) (suggestedOfferRepo.SuggestedOfferRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return suggestedOfferPostgres.NewSuggestedOfferRepository(getDbFunc), nil
}

func NewBlackListUseCase(i *do.Injector) (blacklistsymbol.UseCase, error) {
	blackListSymbolRepo := do.MustInvoke[*blSymbolPostgres.BlackListSymbolRepository](i)
	return blacklistsymbol.NewUseCase(blackListSymbolRepo), nil
}

func NewStockExchangeUseCase(i *do.Injector) (stockexchange.UseCase, error) {
	stockExchangeRepo := do.MustInvoke[*stockExchangePostgres.StockExchangeRepository](i)
	return stockexchange.NewUseCase(stockExchangeRepo), nil
}

func NewSymbolScoreUseCase(i *do.Injector) (symbolscore.UseCase, error) {
	symbolScoreRepo := do.MustInvoke[*symbolScorePostgres.SymbolScoreRepository](i)
	stockExchangeRepo := do.MustInvoke[*stockExchangePostgres.StockExchangeRepository](i)
	return symbolscore.NewUseCase(symbolScoreRepo, stockExchangeRepo), nil
}

func NewSymbolUseCase(i *do.Injector) (symbol.UseCase, error) {
	symbolScoreRepo := do.MustInvoke[*symbolScorePostgres.SymbolScoreRepository](i)
	symbolRepo := do.MustInvoke[*symbolPostgres.SymbolRepository](i)
	atomicExecutor := do.MustInvoke[*atomicity.DbAtomicExecutor](i)
	return symbol.NewUseCase(symbolRepo, symbolScoreRepo, atomicExecutor), nil
}

func NewScoreGroupUseCase(i *do.Injector) (scoregroup.UseCase, error) {
	scoreGroupRepo := do.MustInvoke[*scoreGroupPostgres.ScoreGroupRepository](i)
	scoreGroupInterestRepo := do.MustInvoke[*scoreGroupInterestPostgres.ScoreGroupInterestSqlRepository](i)
	return scoregroup.NewUseCase(scoreGroupRepo, scoreGroupInterestRepo), nil
}

func NewLoanPolicyTemplateUseCase(i *do.Injector) (loanpolicytemplate.UseCase, error) {
	repo := do.MustInvoke[*loanPolicyTemplatePostgres.LoanPolicyTemplateRepository](i)
	return loanpolicytemplate.NewUseCase(repo), nil
}

func NewScoreGroupInterestUseCase(i *do.Injector) (scoregroupinterest.UseCase, error) {
	scoreGroupInterestRepo := do.MustInvoke[*scoreGroupInterestPostgres.ScoreGroupInterestSqlRepository](i)
	loanRequestRepo := do.MustInvoke[*loanPackageRequestPostgres.LoanPackageRequestPostgresRepository](i)
	symbolScoreRepo := do.MustInvoke[*symbolScorePostgres.SymbolScoreRepository](i)
	return scoregroupinterest.NewUseCase(
		scoreGroupInterestRepo, loanRequestRepo, symbolScoreRepo,
	), nil
}

func NewSchedulerUseCase(i *do.Injector) (scheduler.UseCase, error) {
	logger := do.MustInvoke[*slog.Logger](i)
	repo := do.MustInvoke[schedulerRepo.LoanRequestSchedulerConfigRepository](i)
	return scheduler.NewSchedulerUseCase(
		repo, logger,
	), nil
}

func NewFinancialProductUseCase(i *do.Injector) (financialProductDomain.UseCase, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	financialProductRepository := do.MustInvoke[financialProductRepo.FinancialProductRepository](i)
	marginOperationRepository := do.MustInvoke[marginOperationRepo.MarginOperationRepository](i)
	configurationRepository := do.MustInvoke[configRepo.ConfigurationPersistenceRepository](i)
	return financialProductDomain.NewUseCase(cfg, financialProductRepository, marginOperationRepository, configurationRepository), nil
}

func NewSuggestedOfferConfigUseCase(i *do.Injector) (suggestedOfferConfig.UseCase, error) {
	repo := do.MustInvoke[suggestedOfferConfigRepo.SuggestedOfferConfigRepository](i)
	return suggestedOfferConfig.NewUseCase(repo), nil
}

func NewSuggestedOfferUseCase(i *do.Injector) (suggestedOffer.UseCase, error) {
	repo := do.MustInvoke[suggestedOfferRepo.SuggestedOfferRepository](i)
	configRepo := do.MustInvoke[suggestedOfferConfigRepo.SuggestedOfferConfigRepository](i)
	eventRepo := do.MustInvoke[suggestedOfferRepo.SuggestedOfferEventRepository](i)
	orderServiceRepo := do.MustInvoke[orderServiceRepo.OrderServiceRepository](i)
	return suggestedOffer.NewUseCase(repo, configRepo, eventRepo, orderServiceRepo), nil
}

func NewPromotionCampaignUseCase(i *do.Injector) (promotion_campaign.UseCase, error) {
	repo := do.MustInvoke[*promotionCampaignPostgres.PromotionCampaignPostgresRepository](i)
	return promotion_campaign.NewUseCase(repo), nil
}

func NewLoanPolicyTemplateRepository(i *do.Injector) (*loanPolicyTemplatePostgres.LoanPolicyTemplateRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return loanPolicyTemplatePostgres.NewLoanPolicyTemplateRepository(getDbFunc), nil
}

func NewSubmissionSheetRepository(i *do.Injector) (*submissionSheetPostgres.SubmissionSheetPostgresRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return submissionSheetPostgres.NewSubmissionSheetPostgresRepository(getDbFunc), nil
}

func NewConfigurationPersistenceRepository(i *do.Injector) (configRepo.ConfigurationPersistenceRepository, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return configPostgres.NewConfigurationPostgresRepository(getDbFunc), nil
}

func NewLoanPackageRequestUseCase(i *do.Injector) (loanpackagerequest.UseCase, error) {
	loanRequestRepo := do.MustInvoke[*loanPackageRequestPostgres.LoanPackageRequestPostgresRepository](i)
	atomicExecutor := do.MustInvoke[*atomicity.DbAtomicExecutor](i)
	scoreGroupInterestRepo := do.MustInvoke[*scoreGroupInterestPostgres.ScoreGroupInterestSqlRepository](i)
	loanPackageOfferRepo := do.MustInvoke[*loanPackageOfferPostgres.LoanPackageOfferPostgresRepository](i)
	loanPackageOfferInterestRepo := do.MustInvoke[*loanPackageOfferInterestPostgres.LoanPackageOfferInterestPostgresRepository](i)
	loanPackageRequestEventRepo := do.MustInvoke[loanPackageRequestRepo.LoanPackageRequestEventRepository](i)
	symbolRepo := do.MustInvoke[*symbolPostgres.SymbolRepository](i)
	loanContractRepo := do.MustInvoke[*loanContractPostgres.LoanContractRepository](i)
	financialProductClient := do.MustInvoke[financialProductRepo.FinancialProductRepository](i)
	financingApiClient := do.MustInvoke[financingApiRepository.FinancingRepository](i)
	schedulerJobRepo := do.MustInvoke[schedulerRepo.SchedulerJobRepository](i)
	loanPolicyRepository := do.MustInvoke[*loanPolicyTemplatePostgres.LoanPolicyTemplateRepository](i)
	errorService := do.MustInvoke[apperrors.Service](i)
	cfg := do.MustInvoke[config.AppConfig](i)
	logger := do.MustInvoke[*slog.Logger](i)
	investorRepository := do.MustInvoke[investorRepo.InvestorPersistenceRepository](i)
	submissionSheetRepository := do.MustInvoke[*submissionSheetPostgres.SubmissionSheetPostgresRepository](i)
	marginOperationRepository := do.MustInvoke[marginOperationRepo.MarginOperationRepository](i)
	configurationRepository := do.MustInvoke[configRepo.ConfigurationPersistenceRepository](i)
	odooServiceRepository := do.MustInvoke[odooServiceRepo.OdooServiceRepository](i)
	return loanpackagerequest.NewUseCase(
		loanRequestRepo,
		atomicExecutor,
		scoreGroupInterestRepo,
		loanPackageOfferRepo,
		loanPackageOfferInterestRepo,
		loanPackageRequestEventRepo,
		symbolRepo,
		loanContractRepo,
		financialProductClient,
		cfg,
		loanPolicyRepository,
		logger,
		financingApiClient,
		schedulerJobRepo,
		errorService,
		investorRepository,
		submissionSheetRepository,
		marginOperationRepository,
		configurationRepository,
		odooServiceRepository,
	), nil
}

func NewLoanPackageOfferUseCase(i *do.Injector) (loanoffer.UseCase, error) {
	loanPackageOfferRepo := do.MustInvoke[*loanPackageOfferPostgres.LoanPackageOfferPostgresRepository](i)
	cfg := do.MustInvoke[config.AppConfig](i)
	atomicExecutor := do.MustInvoke[*atomicity.DbAtomicExecutor](i)
	loanPackageOfferInterestRepo := do.MustInvoke[*loanPackageOfferInterestPostgres.LoanPackageOfferInterestPostgresRepository](i)
	return loanoffer.NewUseCase(
		loanPackageOfferRepo, cfg.LoanRequest, loanPackageOfferInterestRepo, atomicExecutor,
	), nil
}

func NewLoanOfferInterestUseCase(i *do.Injector) (loanofferinterest.UseCase, error) {
	atomicExecutor := do.MustInvoke[*atomicity.DbAtomicExecutor](i)
	loanPackageOfferInterestRepo := do.MustInvoke[*loanPackageOfferInterestPostgres.LoanPackageOfferInterestPostgresRepository](i)
	loanPackageOfferRepo := do.MustInvoke[*loanPackageOfferPostgres.LoanPackageOfferPostgresRepository](i)
	loanContractRepo := do.MustInvoke[*loanContractPostgres.LoanContractRepository](i)
	financialProductClient := do.MustInvoke[financialProductRepo.FinancialProductRepository](i)
	loanOfferInterestEventRepo := do.MustInvoke[loanOfferInterestRepo.LoanPackageOfferInterestEventRepository](i)
	errorService := do.MustInvoke[apperrors.Service](i)
	symbolRepo := do.MustInvoke[*symbolPostgres.SymbolRepository](i)
	submissionSheetRepo := do.MustInvoke[*submissionSheetPostgres.SubmissionSheetPostgresRepository](i)
	loanTemplateRepo := do.MustInvoke[*loanPolicyTemplatePostgres.LoanPolicyTemplateRepository](i)
	appConfig := do.MustInvoke[config.AppConfig](i)
	return loanofferinterest.NewUseCase(
		loanPackageOfferInterestRepo,
		atomicExecutor,
		loanPackageOfferRepo,
		loanContractRepo,
		financialProductClient,
		loanOfferInterestEventRepo,
		errorService,
		symbolRepo,
		submissionSheetRepo,
		loanTemplateRepo,
		appConfig,
	), nil
}

func NewLoanContractUseCase(i *do.Injector) (loancontract.UseCase, error) {
	atomicExecutor := do.MustInvoke[*atomicity.DbAtomicExecutor](i)
	loanContractRepo := do.MustInvoke[*loanContractPostgres.LoanContractRepository](i)
	loanRequestRepo := do.MustInvoke[*loanPackageRequestPostgres.LoanPackageRequestPostgresRepository](i)
	loanPackageOfferInterestRepo := do.MustInvoke[*loanPackageOfferInterestPostgres.LoanPackageOfferInterestPostgresRepository](i)
	loanPackageOfferRepo := do.MustInvoke[*loanPackageOfferPostgres.LoanPackageOfferPostgresRepository](i)
	return loancontract.NewUseCase(
		atomicExecutor,
		loanContractRepo,
		loanRequestRepo,
		loanPackageOfferInterestRepo,
		loanPackageOfferRepo,
	), nil
}

func NewFeatureUseCase(i *do.Injector) (featureflag.UseCase, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	return featureflag.NewUseCase(cfg.Features), nil
}

func NewConfigUseCase(i *do.Injector) (config.UseCase, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	configurationRepository := do.MustInvoke[configRepo.ConfigurationPersistenceRepository](i)
	return config.NewUseCase(cfg, configurationRepository), nil
}

func NewOfflineOfferUpdateUseCase(i *do.Injector) (offlineofferupdate.UseCase, error) {
	offlineOfferUpdateRepo := do.MustInvoke[offlineOfferRepo.OfflineOfferUpdatePersistenceRepository](i)
	offerRepo := do.MustInvoke[*loanPackageOfferPostgres.LoanPackageOfferPostgresRepository](i)
	offerInterestRepo := do.MustInvoke[*loanPackageOfferInterestPostgres.LoanPackageOfferInterestPostgresRepository](i)
	atomicExecutor := do.MustInvoke[*atomicity.DbAtomicExecutor](i)
	return offlineofferupdate.NewUseCase(offlineOfferUpdateRepo, offerRepo, offerInterestRepo, atomicExecutor), nil
}

func NewAwaitingConfirmRequestUseCase(i *do.Injector) (awaitingconfirmrequest.UseCase, error) {
	repo := do.MustInvoke[awaitingConfirmRequestRepo.AwaitingConfirmRequestPersistenceRepository](i)
	return awaitingconfirmrequest.NewUseCase(repo), nil
}

func NewCombinedRequestUseCase(i *do.Injector) (combinedloanrequest.UseCase, error) {
	repo := do.MustInvoke[combinedRequestRepo.CombinedLoanPackageRequestPersistenceRepository](i)
	return combinedloanrequest.NewUseCase(repo), nil
}

func NewInvestorUseCase(i *do.Injector) (investor.UseCase, error) {
	repo := do.MustInvoke[investorRepo.InvestorPersistenceRepository](i)
	financialProductClient := do.MustInvoke[financialProductRepo.FinancialProductRepository](i)
	return investor.NewUseCase(repo, financialProductClient), nil
}

func NewInvestorAccountUseCase(i *do.Injector) (investor_account.UseCase, error) {
	repo := do.MustInvoke[investorAccountRepo.InvestorAccountRepository](i)
	orderServiceClient := do.MustInvoke[orderServiceRepo.OrderServiceRepository](i)
	return investor_account.NewUseCase(repo, orderServiceClient), nil
}

func NewSubmissionSheetUseCase(i *do.Injector) (submissionsheet.UseCase, error) {
	repo := do.MustInvoke[*submissionSheetPostgres.SubmissionSheetPostgresRepository](i)
	atomicExecutor := do.MustInvoke[*atomicity.DbAtomicExecutor](i)
	loanPolicyRepository := do.MustInvoke[*loanPolicyTemplatePostgres.LoanPolicyTemplateRepository](i)
	marginOperationRepository := do.MustInvoke[marginOperationRepo.MarginOperationRepository](i)
	financialProductRepository := do.MustInvoke[financialProductRepo.FinancialProductRepository](i)
	loanPackageRequestRepository := do.MustInvoke[*loanPackageRequestPostgres.LoanPackageRequestPostgresRepository](i)
	loanPackageOfferRepository := do.MustInvoke[*loanPackageOfferPostgres.LoanPackageOfferPostgresRepository](i)
	loanPackageOfferInterestRepository := do.MustInvoke[*loanPackageOfferInterestPostgres.LoanPackageOfferInterestPostgresRepository](i)
	financingApiClient := do.MustInvoke[financingApiRepository.FinancingRepository](i)
	appConfig := do.MustInvoke[config.AppConfig](i)
	errorService := do.MustInvoke[apperrors.Service](i)
	loanPackageRequestEventRepository := do.MustInvoke[loanPackageRequestRepo.LoanPackageRequestEventRepository](i)
	symbolRepository := do.MustInvoke[*symbolPostgres.SymbolRepository](i)
	return submissionsheet.NewUseCase(
		repo,
		atomicExecutor,
		loanPolicyRepository,
		marginOperationRepository,
		financialProductRepository,
		loanPackageRequestRepository,
		loanPackageOfferRepository,
		loanPackageOfferInterestRepository,
		financingApiClient,
		appConfig,
		errorService,
		loanPackageRequestEventRepository,
		symbolRepository,
	), nil
}

func NewSubmissionDefaultUseCase(i *do.Injector) (submission_default.UseCase, error) {
	configRepository := do.MustInvoke[configRepo.ConfigurationPersistenceRepository](i)
	return submission_default.NewUseCase(configRepository), nil
}

func NewPromotionLoanPackageUseCase(i *do.Injector) (promotionloanpackage.UseCase, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	configRepository := do.MustInvoke[configRepo.ConfigurationPersistenceRepository](i)
	orderServiceRepository := do.MustInvoke[orderServiceRepo.OrderServiceRepository](i)
	financialProductRepository := do.MustInvoke[financialProductRepo.FinancialProductRepository](i)
	promotionCampaignRepo := do.MustInvoke[*promotionCampaignPostgres.PromotionCampaignPostgresRepository](i)
	return promotionloanpackage.NewUseCase(
		cfg, configRepository, orderServiceRepository, financialProductRepository, promotionCampaignRepo,
	), nil
}

func NewMarginOperationUseCase(i *do.Injector) (marginoperation.UseCase, error) {
	marginOperationRepository := do.MustInvoke[marginOperationRepo.MarginOperationRepository](i)
	return marginoperation.NewUseCase(marginOperationRepository), nil
}

func NewConfigurationUseCase(i *do.Injector) (configuration.UseCase, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	configRepository := do.MustInvoke[configRepo.ConfigurationPersistenceRepository](i)
	return configuration.NewUseCase(cfg, configRepository), nil
}

func NewBlackListHandler(i *do.Injector) (*blSymbolHttp.BlacklistSymbolHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	blackListUseCase := do.MustInvoke[blacklistsymbol.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	return blSymbolHttp.NewBlacklistSymbolHandler(baseHandler, logger, blackListUseCase), nil
}

func NewStockExchangeHandler(i *do.Injector) (*http.StockExchangeHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	stockExchangeUseCase := do.MustInvoke[stockexchange.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	return http.NewStockExchangeHandler(baseHandler, logger, stockExchangeUseCase), nil
}

func NewSymbolHandler(i *do.Injector) (*symbolHttp.SymbolHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	symbolUseCase := do.MustInvoke[symbol.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	return symbolHttp.NewSymbolHandler(baseHandler, logger, symbolUseCase), nil
}

func NewSymbolScoreHandler(i *do.Injector) (*symbolScoreHttp.SymbolScoreHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	symbolScoreUseCase := do.MustInvoke[symbolscore.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	return symbolScoreHttp.NewSymbolScoreHandler(baseHandler, logger, symbolScoreUseCase), nil
}

func NewScoreGroupHandler(i *do.Injector) (*scoreGroupHttp.ScoreGroupHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	scoreGroupUseCase := do.MustInvoke[scoregroup.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	return scoreGroupHttp.NewScoreGroupHandler(baseHandler, logger, scoreGroupUseCase), nil
}

func NewScoreGroupInterestHandler(i *do.Injector) (*scoreGroupInterestHttp.ScoreGroupInterestHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	scoreGroupInterestUseCase := do.MustInvoke[scoregroupinterest.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	return scoreGroupInterestHttp.NewScoreGroupInterestHandler(baseHandler, logger, scoreGroupInterestUseCase), nil
}

func NewLoanPackageRequestHandler(i *do.Injector) (*loanPackageRequestHttp.LoanPackageRequestHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	loanPackageRequestUseCase := do.MustInvoke[loanpackagerequest.UseCase](i)
	scoreGroupInterestUseCase := do.MustInvoke[scoregroupinterest.UseCase](i)
	submissionSheetUseCase := do.MustInvoke[submissionsheet.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	return loanPackageRequestHttp.NewLoanPackageRequestHandler(
		baseHandler, logger, loanPackageRequestUseCase, scoreGroupInterestUseCase, submissionSheetUseCase,
	), nil
}

func NewLoanPackageOfferHandler(i *do.Injector) (*loanOfferHttp.LoanPackageOfferHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	loanPackageOfferUseCase := do.MustInvoke[loanoffer.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	offlineOfferUpdateUseCase := do.MustInvoke[offlineofferupdate.UseCase](i)
	offerInterestUseCase := do.MustInvoke[loanofferinterest.UseCase](i)
	return loanOfferHttp.NewLoanPackageOfferHandler(
		baseHandler, logger, loanPackageOfferUseCase, offlineOfferUpdateUseCase, offerInterestUseCase,
	), nil
}

func NewLoanPackageOfferInterestHandler(i *do.Injector) (*loanOfferInterestHttp.LoanOfferInterestHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	loanPackageOfferInterestUseCase := do.MustInvoke[loanofferinterest.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	return loanOfferInterestHttp.NewLoanOfferInterestHandler(
		baseHandler, logger, loanPackageOfferInterestUseCase,
	), nil
}

func NewLoanContractHandler(i *do.Injector) (*loanContractHttp.LoanContractHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	loanContractUseCase := do.MustInvoke[loancontract.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	return loanContractHttp.NewLoanContractHandler(baseHandler, logger, loanContractUseCase), nil
}

func NewFeatureHandler(i *do.Injector) (*http2.FeatureHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	featureUseCase := do.MustInvoke[featureflag.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	return http2.NewFeatureHandler(baseHandler, logger, featureUseCase), nil
}

func NewConfigHandler(i *do.Injector) (*configHttp.ConfigHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	configUseCase := do.MustInvoke[config.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	return configHttp.NewConfigHandler(baseHandler, logger, configUseCase), nil
}

func NewSchedulerHandler(i *do.Injector) (*schedulerHttp.SchedulerHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	logger := do.MustInvoke[*slog.Logger](i)
	useCase := do.MustInvoke[scheduler.UseCase](i)

	return schedulerHttp.NewSchedulerHandler(baseHandler, logger, useCase), nil
}

func NewAwaitingConfirmRequestHandler(i *do.Injector) (*awaitingConfirmRequestHttp.AwaitingConfirmRequestHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	logger := do.MustInvoke[*slog.Logger](i)
	useCase := do.MustInvoke[awaitingconfirmrequest.UseCase](i)
	return awaitingConfirmRequestHttp.NewAwaitingConfirmRequestHandler(baseHandler, logger, useCase), nil
}

func NewCombinedRequestHandler(i *do.Injector) (*combinedRequestHttp.CombinedLoanRequestHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	logger := do.MustInvoke[*slog.Logger](i)
	useCase := do.MustInvoke[combinedloanrequest.UseCase](i)
	return combinedRequestHttp.NewCombinedLoanRequestHandler(baseHandler, logger, useCase), nil
}

func NewInvestorHandler(i *do.Injector) (*investorHttp.InvestorHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	logger := do.MustInvoke[*slog.Logger](i)
	useCase := do.MustInvoke[investor.UseCase](i)
	return investorHttp.NewInvestorHandler(baseHandler, logger, useCase), nil
}

func NewInvestorAccountHandler(i *do.Injector) (*investorAccountHttp.InvestorAccountHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	logger := do.MustInvoke[*slog.Logger](i)
	useCase := do.MustInvoke[investor_account.UseCase](i)
	return investorAccountHttp.NewInvestorAccountHandler(baseHandler, logger, useCase), nil
}

func NewLoanPolicyTemplateHandler(i *do.Injector) (*loanPolicyTemplateHttp.LoanPolicyTemplateHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	useCase := do.MustInvoke[loanpolicytemplate.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	marginOperationUseCase := do.MustInvoke[marginoperation.UseCase](i)
	return loanPolicyTemplateHttp.NewLoanPolicyTemplateHandler(baseHandler, logger, useCase, marginOperationUseCase), nil
}

func NewFinancialProductHandler(i *do.Injector) (*financialProductHttp.FinancialProductHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	useCase := do.MustInvoke[financialProductDomain.UseCase](i)
	return financialProductHttp.NewFinancialProductHandler(baseHandler, useCase), nil
}

func NewSuggestedOfferConfigHandler(i *do.Injector) (*suggestedOfferConfigHttp.SuggestedOfferConfigHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	useCase := do.MustInvoke[suggestedOfferConfig.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	return suggestedOfferConfigHttp.NewSuggestedOfferConfigHandler(baseHandler, logger, useCase), nil
}

func NewSuggestedOfferHandler(i *do.Injector) (*suggestedOfferHttp.SuggestedOfferHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	useCase := do.MustInvoke[suggestedOffer.UseCase](i)
	configUseCase := do.MustInvoke[suggestedOfferConfig.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	return suggestedOfferHttp.NewSuggestedOfferHandler(baseHandler, logger, useCase, configUseCase), nil
}

func NewPromotionLoanPackageHandler(i *do.Injector) (*promotionLoanPackageHttp.PromotionLoanPackageHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	useCase := do.MustInvoke[promotionloanpackage.UseCase](i)
	cacheStore := do.MustInvoke[cache.Cache](i)
	return promotionLoanPackageHttp.NewPromotionLoanPackageHandler(baseHandler, cacheStore, useCase), nil
}

func NewConfigurationHandler(i *do.Injector) (*configurationHttp.ConfigurationHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	useCase := do.MustInvoke[configuration.UseCase](i)
	cacheStore := do.MustInvoke[cache.Cache](i)
	return configurationHttp.NewConfigurationHandler(baseHandler, cacheStore, useCase), nil
}

func NewLoanOfferScheduler(i *do.Injector) (*loanOfferScheduler.LoanOfferScheduler, error) {
	logger := do.MustInvoke[*slog.Logger](i)
	useCase := do.MustInvoke[loanoffer.UseCase](i)
	errorService := do.MustInvoke[apperrors.Service](i)
	return loanOfferScheduler.NewLoanOfferScheduler(logger, useCase, errorService), nil
}

func NewLoanPackageRequestScheduler(i *do.Injector) (*loanPackageScheduler.LoanRequestScheduler, error) {
	logger := do.MustInvoke[*slog.Logger](i)
	schedulerUseCase := do.MustInvoke[scheduler.UseCase](i)
	useCase := do.MustInvoke[loanpackagerequest.UseCase](i)
	errorService := do.MustInvoke[apperrors.Service](i)
	return loanPackageScheduler.NewLoanRequestScheduler(logger, schedulerUseCase, useCase, errorService), nil
}

func NewSubmissionSheetHandler(i *do.Injector) (*submissionSheetHttp.SubmissionSheetHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	logger := do.MustInvoke[*slog.Logger](i)
	useCase := do.MustInvoke[submissionsheet.UseCase](i)
	return submissionSheetHttp.NewSubmissionSheetHandler(baseHandler, logger, useCase), nil
}

func NewSubmissionDefaultHandler(i *do.Injector) (*submissionDefaultHttp.SubmissionDefaultHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	useCase := do.MustInvoke[submission_default.UseCase](i)
	return submissionDefaultHttp.NewSubmissionDefaultHandler(baseHandler, useCase), nil
}

func NewPromotionCampaignHandler(i *do.Injector) (*promotionCampaignHttp.PromotionCampaignHandler, error) {
	baseHandler := do.MustInvoke[handler.BaseHandler](i)
	useCase := do.MustInvoke[promotion_campaign.UseCase](i)
	logger := do.MustInvoke[*slog.Logger](i)
	return promotionCampaignHttp.NewPromotionCampaignHandler(baseHandler, logger, useCase), nil
}

func NewCache(i *do.Injector) (cache.Cache, error) {
	tasks := do.MustInvoke[*shutdown.Tasks](i)
	return cache.NewInProcessCache(tasks)
}

func NewTemporalClient(i *do.Injector) (client.Client, error) {
	logger := do.MustInvoke[*slog.Logger](i)
	tasks := do.MustInvoke[*shutdown.Tasks](i)
	cfg := do.MustInvoke[config.AppConfig](i)
	c, err := client.Dial(
		client.Options{
			HostPort:  cfg.Temporal.Host,
			Namespace: cfg.Temporal.Namespace,
			Logger:    logger,
		},
	)
	if err != nil {
		if cfg.Env == "local" {
			logger.Error("failed to connect to temporal client", slog.String("error", err.Error()))
			return nil, nil
		}
		return nil, err
	}
	tasks.AddShutdownTask(
		func(_ context.Context) error {
			c.Close()
			return nil
		},
	)
	return c, nil
}
