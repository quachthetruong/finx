package dependencies

import (
	"log/slog"

	"financing-offer/internal/apperrors/repository"
	"financing-offer/internal/atomicity"
	configHttp "financing-offer/internal/config/transport/http"
	blacklistSymbolHttp "financing-offer/internal/core/blacklistsymbol/transport/http"
	financialProductRepo "financing-offer/internal/core/financialproduct/repository"
	loanContractHttp "financing-offer/internal/core/loancontract/transport/http"
	loanPackageOfferHttp "financing-offer/internal/core/loanoffer/transport/http"
	loanOfferInterestHttp "financing-offer/internal/core/loanofferinterest/http"
	loanPackageRequestHttp "financing-offer/internal/core/loanpackagerequest/transport/http"
	scoreGroupHttp "financing-offer/internal/core/scoregroup/transport/http"
	scoreGroupInterestHttp "financing-offer/internal/core/scoregroupinterest/transport/http"
	stockExchangeHttp "financing-offer/internal/core/stockexchange/transport/http"
	symbolHttp "financing-offer/internal/core/symbol/transport/http"
	symbolScoreHttp "financing-offer/internal/core/symbolscore/transport/http"
	"financing-offer/internal/database"
	"financing-offer/internal/event"
	"financing-offer/internal/featureflag"
	featureHandler "financing-offer/internal/featureflag/transport/http"
)

type Provider struct {
	GetDbFunc                 database.GetDbFunc
	Logger                    *slog.Logger
	AtomicExecutor            atomicity.AtomicExecutor
	SymbolHandler             *symbolHttp.SymbolHandler
	StockExchangeHandler      *stockExchangeHttp.StockExchangeHandler
	SymbolScoreHandler        *symbolScoreHttp.SymbolScoreHandler
	ScoreGroupHandler         *scoreGroupHttp.ScoreGroupHandler
	LoanPackageRequestHandler *loanPackageRequestHttp.LoanPackageRequestHandler
	ScoreGroupInterestHandler *scoreGroupInterestHttp.ScoreGroupInterestHandler
	LoanPackageOfferHandler   *loanPackageOfferHttp.LoanPackageOfferHandler
	LoanContractHandler       *loanContractHttp.LoanContractHandler
	OfferInterestHandler      *loanOfferInterestHttp.LoanOfferInterestHandler
	FeatureHandler            *featureHandler.FeatureHandler
	ConfigHandler             *configHttp.ConfigHandler
	BlacklistSymbolHandler    *blacklistSymbolHttp.BlacklistSymbolHandler

	FeatureFlagUseCase featureflag.UseCase
}

type ExternalProvider struct {
	GetDbFunc              database.GetDbFunc
	AtomicExecutor         atomicity.AtomicExecutor
	Publisher              event.Publisher
	NotifyWebhookClient    repository.NotifyWebhookRepository
	FinancialProductClient financialProductRepo.FinancialProductRepository
}
