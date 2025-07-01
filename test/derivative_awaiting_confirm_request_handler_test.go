package test

import (
	"testing"
	"time"

	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v9"

	http2 "financing-offer/internal/core/awaiting_confirm_request/transport/http"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestAwaitingConfirmRequestDerivativeHandler(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	handler := do.MustInvoke[*http2.AwaitingConfirmRequestHandler](injector)
	t.Run(
		"Get all awaiting confirm requests success", func(t *testing.T) {
			defer truncateData()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 0,
					MaxScore: 100,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "BID",
					AssetType:       "UNDERLYING",
				},
			)
			derivativeSymbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "VN30F1M",
					AssetType:       "DERIVATIVE",
				},
			)
			request := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:    symbol.ID,
					InvestorID:  "investor",
					AccountNo:   "accountNo",
					LoanRate:    decimal.NewFromFloat(0.8),
					LimitAmount: decimal.NewFromInt(1000000000),
					Type:        "FLEXIBLE",
					Status:      "CONFIRMED",
					AssetType:   "UNDERLYING",
				},
			)
			request1 := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:    derivativeSymbol.ID,
					InvestorID:  "investor",
					AccountNo:   "accountNo",
					LoanRate:    decimal.NewFromFloat(0.7),
					LimitAmount: decimal.NewFromInt(200000000),
					Type:        "FLEXIBLE",
					Status:      "CONFIRMED",
					AssetType:   "DERIVATIVE",
				},
			)
			request2 := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:    derivativeSymbol.ID,
					InvestorID:  "investor",
					AccountNo:   "accountNo",
					LoanRate:    decimal.NewFromFloat(0.7),
					LimitAmount: decimal.NewFromInt(300000000),
					Type:        "FLEXIBLE",
					Status:      "CONFIRMED",
					AssetType:   "DERIVATIVE",
				},
			)
			offer := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().UTC().Add(time.Hour * -1)),
					FlowType:             entity.FlowTypeDnseOnline.String(),
				},
			)
			offer1 := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request1.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().UTC().Add(time.Hour * -1)),
					FlowType:             entity.FLowTypeDnseOffline.String(),
				},
			)
			offer2 := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request2.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().UTC().Add(time.Hour * -1)),
					FlowType:             entity.FLowTypeDnseOffline.String(),
				},
			)
			mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.NewFromFloat(0.6),
					InterestRate:       decimal.NewFromFloat(0.12),
					Status:             "PENDING",
					LoanID:             5,
					AssetType:          "UNDERLYING",
				},
			)
			mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.NewFromFloat(0.6),
					InterestRate:       decimal.NewFromFloat(0.12),
					Status:             "PENDING",
					LoanID:             78,
					AssetType:          "UNDERLYING",
				},
			)
			mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer1.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.NewFromFloat(0.6),
					InterestRate:       decimal.NewFromFloat(0.12),
					Status:             "PENDING",
					LoanID:             6,
					AssetType:          "DERIVATIVE",
				},
			)
			mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer1.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.NewFromFloat(0.6),
					InterestRate:       decimal.NewFromFloat(0.12),
					Status:             "PENDING",
					LoanID:             7,
					AssetType:          "DERIVATIVE",
				},
			)
			mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer2.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.NewFromFloat(0.6),
					InterestRate:       decimal.NewFromFloat(0.12),
					Status:             "PENDING",
					LoanID:             8,
					AssetType:          "DERIVATIVE",
				},
			)
			mock.SeedOfflineOfferUpdate(
				t, db, model.OfflineOfferUpdate{
					OfferID:   offer1.ID,
					Status:    "đồng ý",
					Category:  "CA1",
					CreatedBy: "admin 1",
				},
			)
			mock.SeedOfflineOfferUpdate(
				t, db, model.OfflineOfferUpdate{
					OfferID:   offer2.ID,
					Status:    "chưa đồng ý",
					Category:  "CA2",
					CreatedBy: "admin 2",
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("GET", "/api/v1/derivative-awaiting-confirm-requests?sort=id", nil)
			handler.GetAllDerivative(ginCtx)
			result := recorder.Result()
			body := gintest.ExtractBody(result.Body)
			assert.Nil(t, result.Body.Close())
			assert.Equal(t, 2, testhelper.GetArrayLength(body, "data"))
			assert.Equal(t, int64(2), testhelper.GetInt(body, "metaData", "total"))
			assert.Equal(t, "6, 7", testhelper.GetString(body, "data", "[0]", "loanPackageIds"))
			assert.Equal(t, "CA1", testhelper.GetString(body, "data", "[0]", "latestUpdate", "category"))
			assert.Equal(t, "CA2", testhelper.GetString(body, "data", "[1]", "latestUpdate", "category"))
			assert.Equal(t, "8", testhelper.GetString(body, "data", "[1]", "loanPackageIds"))

			ginCtx, _, recorder = gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/derivative-awaiting-confirm-requests?sort=id&latestUpdateCategories=CA2", nil,
			)
			handler.GetAllDerivative(ginCtx)
			result = recorder.Result()
			body = gintest.ExtractBody(result.Body)
			assert.Nil(t, result.Body.Close())
			assert.Equal(t, 1, testhelper.GetArrayLength(body, "data"))
			assert.Equal(t, "CA2", testhelper.GetString(body, "data", "[0]", "latestUpdate", "category"))
		},
	)
}
