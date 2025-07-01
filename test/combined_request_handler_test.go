package test

import (
	"testing"
	"time"

	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v9"

	"financing-offer/internal/core/combined_loan_request/transport/http"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestCombinedRequestHandler(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	handler := do.MustInvoke[*http.CombinedLoanRequestHandler](injector)
	t.Run(
		"Get all combined requests success", func(t *testing.T) {
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
					SymbolID:    symbol.ID,
					InvestorID:  "investor",
					AccountNo:   "accountNo",
					LoanRate:    decimal.NewFromFloat(0.7),
					LimitAmount: decimal.NewFromInt(200000000),
					Type:        "FLEXIBLE",
					Status:      "CONFIRMED",
					AssetType:   "UNDERLYING",
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
					FlowType:             entity.FlowTypeDnseOnline.String(),
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
					LoanID:             420,
					InterestRate:       decimal.NewFromFloat(0.12),
					Status:             "PACKAGE_CREATED",
					AssetType:          "UNDERLYING",
				},
			)
			mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer1.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.NewFromFloat(0.6),
					InterestRate:       decimal.NewFromFloat(0.12),
					Status:             "CANCELLED",
					AssetType:          "UNDERLYING",
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("GET", "/api/v1/combined-requests", nil)
			handler.GetAll(ginCtx)
			result := recorder.Result()
			body := gintest.ExtractBody(result.Body)
			assert.Nil(t, result.Body.Close())
			assert.Equal(t, 2, testhelper.GetArrayLength(body, "data"))
			assert.Equal(t, int64(2), testhelper.GetInt(body, "metaData", "total"))
			assert.Equal(t, "AWAITING_CONFIRM", testhelper.GetString(body, "data", "[1]", "status"))
			assert.Equal(t, "PACKAGE_CREATED", testhelper.GetString(body, "data", "[0]", "status"))
			assert.Equal(t, "5, 78", testhelper.GetString(body, "data", "[1]", "adminAssignedLoanPackageIds"))
			assert.Equal(t, "420", testhelper.GetString(body, "data", "[0]", "activatedLoanPackageIds"))
			assert.Equal(t, int64(2), testhelper.GetInt(body, "metaData", "total"))

			ginCtx, _, recorder = gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/combined-requests?activatedLoanId=420", nil,
			)
			handler.GetAll(ginCtx)
			result = recorder.Result()
			body = gintest.ExtractBody(result.Body)
			assert.Nil(t, result.Body.Close())
			assert.Equal(t, 1, testhelper.GetArrayLength(body, "data"))
			assert.Equal(t, "PACKAGE_CREATED", testhelper.GetString(body, "data", "[0]", "status"))
			assert.Equal(t, int64(1), testhelper.GetInt(body, "metaData", "total"))
		},
	)
}
