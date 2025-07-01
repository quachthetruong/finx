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

func TestDerivativeLoanOfferHandler(t *testing.T) {
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
					AssetType:       model.AssetType_Derivative,
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
					AssetType:  model.AssetType_Derivative,
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
				"POST", "/api/v1/derivative-loan-package-offers/:id/offline-updates", http.CreateOfferUpdateRequest{
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
			h.CreateDerivativeOfflineOfferUpdate(ginCtx)
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
					AssetType:       model.AssetType_Derivative,
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
					AssetType:  model.AssetType_Derivative,
				},
			)
			offer := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().UTC()),
					FlowType:             entity.FLowTypeDnseOffline.String(),
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/derivative-loan-package-offers/:id/offline-updates", http.CreateOfferUpdateRequest{
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
			h.CreateDerivativeOfflineOfferUpdate(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, offer.ID, testhelper.GetInt(body, "data", "offerId"))
			assert.Equal(t, "PROCESSING", testhelper.GetString(body, "data", "status"))
		},
	)

	t.Run(
		"create offline offer update when offer not exist", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/derivative-loan-package-offers/:id/offline-updates", http.CreateOfferUpdateRequest{
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
			h.CreateDerivativeOfflineOfferUpdate(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusNotFound, result.StatusCode)
		},
	)

	t.Run(
		"create offline offer update when offer assetType does not match", func(t *testing.T) {
			defer truncateData()
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
					AssetType:       model.AssetType_Underlying,
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
					AssetType:  model.AssetType_Underlying,
				},
			)
			offer := mock.SeedLoanPackageOffer(
				t, db, model.LoanPackageOffer{
					LoanPackageRequestID: request.ID,
					OfferedBy:            "admin",
					ExpiredAt:            null.TimeFrom(time.Now().UTC()),
					FlowType:             entity.FLowTypeDnseOffline.String(),
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/derivative-loan-package-offers/:id/offline-updates", http.CreateOfferUpdateRequest{
					Status:   "APPROVED",
					Category: "Đã hoàn thiện HĐ",
					Note:     "note",
				},
			)
			ginCtx.AddParam("id", strconv.FormatInt(offer.ID, 10))
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "customer success",
				},
			)
			h.CreateDerivativeOfflineOfferUpdate(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusBadRequest, result.StatusCode)
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
					AssetType:       model.AssetType_Derivative,
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
					AssetType:  model.AssetType_Derivative,
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
					CreatedBy: "admin",
				},
			)

			mock.SeedOfflineOfferUpdate(
				t, db, model.OfflineOfferUpdate{
					OfferID:   offer.ID,
					Status:    "4",
					Category:  "5",
					Note:      "6",
					CreatedBy: "admin",
				},
			)
			mock.SeedOfflineOfferUpdate(
				t, db, model.OfflineOfferUpdate{
					OfferID:   offer.ID,
					Status:    "4",
					Category:  "5",
					Note:      "6",
					CreatedBy: "admin",
				},
			)

			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/derivative-loan-package-offers/:id/offline-updates", nil,
			)
			ginCtx.AddParam("id", strconv.Itoa(int(offer.ID)))
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "customer success",
				},
			)
			h.GetDerivativeOfflineOfferUpdateHistory(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 3, testhelper.GetArrayLength(body, "data"))
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
					AssetType:       model.AssetType_Derivative,
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
				Name:        "FinX 40%",
				InitialRate: decimal.NewFromFloat(0.6),
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
				t, loanPackage.InitialRate.Equal(cancelledOfferInterest.InitialRate),
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
					AssetType:       model.AssetType_Derivative,
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
