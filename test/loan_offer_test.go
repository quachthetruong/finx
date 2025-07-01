package test

import (
	"fmt"
	gohttp "net/http"
	"strconv"
	"testing"
	"time"

	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/volatiletech/null/v9"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/financialproduct/repository"
	"financing-offer/internal/core/loanoffer/transport/http"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/event"
	"financing-offer/internal/jwttoken"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestLoanOfferHandler(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	financialProductClientMock := mock.NewMockFinancialProductRepository(t)
	eventPublisher := mock.EventPublisher{}
	do.OverrideValue[repository.FinancialProductRepository](injector, financialProductClientMock)
	do.OverrideValue[event.Publisher](injector, eventPublisher)
	h := do.MustInvoke[*http.LoanPackageOfferHandler](injector)
	investorId := "0001000115"
	accountNo := "0001000115"
	t.Run(
		"test cancel loan package offer", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
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
					SymbolID:   symbol.ID,
					InvestorID: investorId,
					AccountNo:  accountNo,
					LoanRate:   decimal.NewFromFloat(0.8),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			offer := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().Add(time.Hour * 24 * 30)),
					FlowType:             entity.FlowTypeDnseOnline.String(),
				},
			)
			offerInterest := mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           request.LoanRate,
					InterestRate:       decimal.Zero,
					Status:             "PENDING",
					AssetType:          "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/my-loan-package-offer", nil,
			)
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					InvestorId: investorId,
					Sub:        investorId,
				},
			)
			ginCtx.AddParam("id", strconv.FormatInt(offerInterest.ID, 10))
			h.InvestorCancelLoanPackageOffer(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "ok", testhelper.GetString(body, "data"))
			updatedOfferInterest := testhelper.GetLoanOfferInterest(t, db, offerInterest.ID)
			assert.Equal(t, "CANCELLED", updatedOfferInterest.Status)
			assert.Equal(t, investorId, updatedOfferInterest.CancelledBy)
		},
	)

	t.Run(
		"cancel expired loan offer success", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
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
					SymbolID:   symbol.ID,
					InvestorID: investorId,
					AccountNo:  accountNo,
					LoanRate:   decimal.NewFromFloat(0.8),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			offer := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().Add(time.Hour * -1)),
					FlowType:             entity.FlowTypeDnseOnline.String(),
				},
			)
			offer1 := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().Add(time.Hour * -1)),
					FlowType:             entity.FlowTypeDnseOnline.String(),
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
			offerInterest1 := mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.NewFromFloat(0.6),
					InterestRate:       decimal.NewFromFloat(0.12),
					Status:             "PENDING",
					AssetType:          "UNDERLYING",
				},
			)
			offerInterest2 := mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.NewFromFloat(0.7),
					InterestRate:       decimal.NewFromFloat(0.15),
					Status:             "PENDING",
					AssetType:          "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/internal/loan-package-offers/expire", nil,
			)
			h.ManualTriggerExpireLoanOffers(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "ok", testhelper.GetString(body, "data"))
			updatedOfferInterest1 := testhelper.GetLoanOfferInterest(t, db, offerInterest1.ID)
			assert.Equal(t, "CANCELLED", updatedOfferInterest1.Status)
			assert.Equal(t, "system", updatedOfferInterest1.CancelledBy)
			updatedOfferInterest2 := testhelper.GetLoanOfferInterest(t, db, offerInterest2.ID)
			assert.Equal(t, "CANCELLED", updatedOfferInterest2.Status)
			assert.Equal(t, "system", updatedOfferInterest2.CancelledBy)
			assert.Equal(
				t, entity.LoanPackageOfferCancelledReasonExpired.String(), updatedOfferInterest2.CancelledReason,
			)
		},
	)

	t.Run(
		"investor get loan offer", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
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
					SymbolID:   symbol.ID,
					InvestorID: investorId,
					AccountNo:  accountNo,
					LoanRate:   decimal.NewFromFloat(0.8),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			offer := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().UTC().Add(time.Hour * 1)),
					FlowType:             entity.FlowTypeDnseOnline.String(),
				},
			)
			mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().UTC().Add(time.Hour * 1)),
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
					AssetType:          "UNDERLYING",
				},
			)
			mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.NewFromFloat(0.7),
					InterestRate:       decimal.NewFromFloat(0.12),
					Status:             "PENDING",
					AssetType:          "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/my-loan-package-offer", nil,
			)
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					InvestorId: investorId,
					Sub:        investorId,
				},
			)
			h.InvestorFindLoanOffers(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 2, testhelper.GetArrayLength(body, "data"))
			assert.Equal(t, 2, testhelper.GetArrayLength(body, "data", "[0]", "loanPackageOfferInterests"))
		},
	)

	t.Run(
		"investor get loan offer with expired offers", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
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
					SymbolID:   symbol.ID,
					InvestorID: investorId,
					AccountNo:  accountNo,
					LoanRate:   decimal.NewFromFloat(0.8),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
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
					LoanPackageRequestID: request.ID,
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
					AssetType:          "UNDERLYING",
				},
			)
			mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer1.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.NewFromFloat(0.6),
					InterestRate:       decimal.NewFromFloat(0.12),
					Status:             "PACKAGE_CREATED",
					AssetType:          "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest("GET", "/api/v1/my-loan-package-offer", nil)
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					InvestorId: investorId,
					Sub:        investorId,
				},
			)
			h.InvestorFindLoanOffers(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 1, testhelper.GetArrayLength(body, "data"))
			assert.Equal(t, offer1.ID, testhelper.GetInt(body, "data", "[0]", "id"))
		},
	)

	t.Run(
		"investor get loan offer by id", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
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
					SymbolID:   symbol.ID,
					InvestorID: investorId,
					AccountNo:  accountNo,
					LoanRate:   decimal.NewFromFloat(0.8),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
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
			mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.NewFromFloat(0.6),
					InterestRate:       decimal.NewFromFloat(0.12),
					Status:             "CANCELLED",
					AssetType:          "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest("GET", "/api/v1/my-loan-package-offer", nil)
			ginCtx.AddParam("id", strconv.Itoa(int(offer.ID)))
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					InvestorId: investorId,
					Sub:        investorId,
				},
			)
			h.InvestorGetById(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, offer.ID, testhelper.GetInt(body, "data", "id"))
		},
	)

	t.Run(
		"create offline offer update success", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 0,
					MaxScore: 78,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "VND",
					AssetType:       "UNDERLYING",
				},
			)
			request := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: investorId,
					AccountNo:  accountNo,
					LoanRate:   decimal.NewFromFloat(0.8),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			offer := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().UTC()),
					FlowType:             entity.FlowTypeDnseOnline.String(),
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/loan-package-offers/:id/offline-updates", http.CreateOfferUpdateRequest{
					Status:   "APPROVED",
					Category: "Đã hoàn thiện HĐ",
					Note:     "note",
				},
			)
			ginCtx.AddParam("id", strconv.Itoa(int(offer.ID)))
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "customer success",
				},
			)
			h.CreateOfflineOfferUpdate(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, offer.ID, testhelper.GetInt(body, "data", "offerId"))
			assert.Equal(t, "customer success", testhelper.GetString(body, "data", "createdBy"))
		},
	)

	t.Run(
		"create offline offer update success with status processing", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 0,
					MaxScore: 78,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "VND",
					AssetType:       "UNDERLYING",
				},
			)
			request := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: investorId,
					AccountNo:  accountNo,
					LoanRate:   decimal.NewFromFloat(0.8),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			offer := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().UTC()),
					FlowType:             entity.FlowTypeDnseOnline.String(),
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/loan-package-offers/:id/offline-updates", http.CreateOfferUpdateRequest{
					Status:   "PROCESSING",
					Category: "Không liên hệ được",
					Note:     "note",
				},
			)
			ginCtx.AddParam("id", strconv.Itoa(int(offer.ID)))
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "customer success",
				},
			)
			h.CreateOfflineOfferUpdate(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, offer.ID, testhelper.GetInt(body, "data", "offerId"))
			assert.Equal(t, "PROCESSING", testhelper.GetString(body, "data", "status"))
		},
	)

	t.Run(
		"create offline offer update success with status cancelled", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 0,
					MaxScore: 78,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "VND",
					AssetType:       "UNDERLYING",
				},
			)
			request := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: investorId,
					AccountNo:  accountNo,
					LoanRate:   decimal.NewFromFloat(0.8),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			offer := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().UTC()),
					FlowType:             entity.FlowTypeDnseOnline.String(),
				},
			)

			mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					ID:                 1,
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           request.LoanRate,
					Status:             string(entity.LoanPackageOfferInterestStatusPending),
					AssetType:          "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/loan-package-offers/:id/offline-updates", http.CreateOfferUpdateRequest{
					Status:   "REJECTED",
					Category: "Lãi xuất vay cao",
					Note:     "note",
				},
			)
			ginCtx.AddParam("id", strconv.Itoa(int(offer.ID)))
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "customer success",
				},
			)
			h.CreateOfflineOfferUpdate(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, offer.ID, testhelper.GetInt(body, "data", "offerId"))
			assert.Equal(t, "REJECTED", testhelper.GetString(body, "data", "status"))

			ginCtx, _, recorder = gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/my-loan-package-offer/:id", nil,
			)
			ginCtx.AddParam("id", strconv.Itoa(int(offer.ID)))
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub:        investorId,
					InvestorId: investorId,
				},
			)
			h.InvestorGetById(ginCtx)
			result = recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body = gintest.ExtractBody(result.Body)
			assert.Equal(t, offer.ID, testhelper.GetInt(body, "data", "id"))
			assert.Equal(
				t, "CANCELLED", testhelper.GetString(body, "data", "loanPackageOfferInterests", "[0]", "status"),
			)
			assert.Equal(
				t, "INVESTOR",
				testhelper.GetString(body, "data", "loanPackageOfferInterests", "[0]", "cancelledReason"),
			)
			assert.Equal(
				t, "customer success",
				testhelper.GetString(body, "data", "loanPackageOfferInterests", "[0]", "cancelledBy"),
			)
		},
	)

	t.Run(
		"create offline offer update when offer not exist", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/loan-package-offers/:id/offline-updates", http.CreateOfferUpdateRequest{
					Status:   "APPROVED",
					Category: "Đã hoàn thiện HĐ",
					Note:     "note",
				},
			)
			ginCtx.AddParam("id", "1")
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "customer success",
				},
			)
			h.CreateOfflineOfferUpdate(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusNotFound, result.StatusCode)
		},
	)

	t.Run(
		"get offline offer history success", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 0,
					MaxScore: 78,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "VND",
					AssetType:       "UNDERLYING",
				},
			)
			request := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: investorId,
					AccountNo:  accountNo,
					LoanRate:   decimal.NewFromFloat(0.8),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			offer := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().UTC()),
					FlowType:             entity.FlowTypeDnseOnline.String(),
				},
			)
			mock.SeedOfflineOfferUpdate(
				t, db, model.OfflineOfferUpdate{
					OfferID:   offer.ID,
					Status:    "1",
					Category:  "2",
					Note:      "3",
					CreatedBy: "duc",
				},
			)

			mock.SeedOfflineOfferUpdate(
				t, db, model.OfflineOfferUpdate{
					OfferID:   offer.ID,
					Status:    "4",
					Category:  "5",
					Note:      "6",
					CreatedBy: "duc",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/loan-package-offer/:id/offline-updates", nil,
			)
			ginCtx.AddParam("id", strconv.Itoa(int(offer.ID)))
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "customer success",
				},
			)
			h.GetOfflineOfferUpdateHistory(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 2, testhelper.GetArrayLength(body, "data"))
		},
	)

	t.Run(
		"admin cancel loan offer interest", func(t *testing.T) {
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
					SymbolID:   symbol.ID,
					InvestorID: investorId,
					AccountNo:  accountNo,
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			offer := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().Add(time.Hour * 24 * 30)),
					FlowType:             entity.FlowTypeDnseOnline.String(),
				},
			)
			offerInterest := mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.Zero,
					InterestRate:       decimal.Zero,
					FeeRate:            decimal.Zero,
					LoanID:             145,
					Status:             "PENDING",
					AssetType:          "UNDERLYING",
				},
			)
			ginCtx, _, _ := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "admin",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest("POST", "/loan-package-offer-interests/1/admin-cancel", nil)
			ginCtx.AddParam("id", fmt.Sprintf("%d", offer.ID))
			h.AdminCancelLoanPackageOfferInterest(ginCtx)
			cancelledOfferInterest := testhelper.GetLoanOfferInterest(t, db, offerInterest.ID)
			assert.Equal(t, "CANCELLED", cancelledOfferInterest.Status)
			assert.Equal(t, "admin", cancelledOfferInterest.CancelledBy)
		},
	)

	t.Run(
		"(offline flow) admin assign loan id to offer interest", func(t *testing.T) {
			defer truncateData()
			loanId := int64(9978)
			loanPackageAccountId := int64(5454)
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
					SymbolID:   symbol.ID,
					InvestorID: investorId,
					AccountNo:  accountNo,
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			offer := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().Add(time.Hour * 24 * 30)),
					FlowType:             entity.FLowTypeDnseOffline.String(),
				},
			)
			offerInterest := mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.Zero,
					InterestRate:       decimal.Zero,
					FeeRate:            decimal.Zero,
					Status:             "PENDING",
					AssetType:          "UNDERLYING",
				},
			)
			loanPackage := entity.FinancialProductLoanPackage{
				Id:            loanId,
				Name:          "FinX 40%",
				InitialRate:   decimal.NewFromFloat(0.6),
				InterestRate:  decimal.NewFromFloat(0.15),
				Term:          100,
				BuyingFeeRate: decimal.NewFromFloat(0.012),
			}
			financialProductClientMock.EXPECT().GetLoanPackageDetail(testifyMock.Anything, loanId).Return(
				loanPackage, nil,
			)
			financialProductClientMock.EXPECT().AssignLoanPackageOrGetLoanPackageAccountId(
				testifyMock.Anything, accountNo, loanId, entity.AssetTypeUnderlying,
			).Return(loanPackageAccountId, nil)
			financialProductClientMock.EXPECT().GetAllAccountDetail(testifyMock.Anything, investorId).Return(
				[]entity.FinancialAccountDetail{}, nil,
			).Maybe()
			ginCtx, _, _ := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "admin",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/loan-offer-interests/:id/assign-loan", http.AdminAssignLoanIdRequest{
					LoanId: loanId,
				},
			)
			ginCtx.AddParam("id", fmt.Sprintf("%d", offer.ID))
			h.AdminAssignLoanId(ginCtx)
			cancelledOfferInterest := testhelper.GetLoanOfferInterest(t, db, offerInterest.ID)
			assert.Equal(t, "PACKAGE_CREATED", cancelledOfferInterest.Status)
			assert.Equal(t, loanId, cancelledOfferInterest.LoanID)
			assert.True(
				t, decimal.NewFromInt(1).Sub(loanPackage.InitialRate).Equal(cancelledOfferInterest.LoanRate),
			)
		},
	)

	t.Run(
		"(derivative) admin assign loan id to offer interest success", func(t *testing.T) {
			defer truncateData()
			loanId := int64(7886)
			loanPackageAccountId := int64(2344)
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code: "HOSE",
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "TCB",
					AssetType:       "UNDERLYING",
				},
			)
			request := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: investorId,
					AccountNo:  accountNo,
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  model.AssetType_Derivative,
				},
			)
			offer := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().Add(time.Hour * 24 * 30)),
					FlowType:             entity.FLowTypeDnseOffline.String(),
				},
			)
			offerInterest := mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.Zero,
					InterestRate:       decimal.Zero,
					FeeRate:            decimal.Zero,
					Status:             "PENDING",
					AssetType:          model.AssetType_Derivative,
				},
			)
			loanPackage := entity.FinancialProductLoanPackageDerivative{
				Id:          loanId,
				Name:        "Goi vay 02",
				InitialRate: decimal.NewFromFloat(0.178),
			}
			financialProductClientMock.EXPECT().GetLoanPackageDerivative(testifyMock.Anything, loanId).Return(
				loanPackage, nil,
			)
			financialProductClientMock.EXPECT().AssignLoanPackageOrGetLoanPackageAccountId(
				testifyMock.Anything, accountNo, loanId, entity.AssetTypeDerivative,
			).Return(loanPackageAccountId, nil)
			financialProductClientMock.EXPECT().GetAllAccountDetail(testifyMock.Anything, investorId).Return(
				[]entity.FinancialAccountDetail{}, nil,
			).Maybe()
			ginCtx, _, _ := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "admin",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/loan-offer-interests/:id/assign-loan", http.AdminAssignLoanIdRequest{
					LoanId: loanId,
				},
			)
			ginCtx.AddParam("id", fmt.Sprintf("%d", offer.ID))
			h.AdminAssignLoanId(ginCtx)
			cancelledOfferInterest := testhelper.GetLoanOfferInterest(t, db, offerInterest.ID)
			assert.Equal(t, "PACKAGE_CREATED", cancelledOfferInterest.Status)
			assert.Equal(t, loanId, cancelledOfferInterest.LoanID)
			assert.True(
				t, decimal.NewFromFloat(0.178).Equal(cancelledOfferInterest.InitialRate),
			)
		},
	)
}
