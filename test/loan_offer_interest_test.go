package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/volatiletech/null/v9"
	"go.temporal.io/sdk/client"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/financialproduct/repository"
	loanOfferInterestHttp "financing-offer/internal/core/loanofferinterest/http"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/event"
	"financing-offer/internal/jwttoken"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestLoanOfferInterestHandler(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	investorId := "0001000115"
	accountNo := "0001000115"
	financialProductClientMock := mock.NewMockFinancialProductRepository(t)
	temporalClientMock := mock.NewMockTemporalClient(t)
	eventPublisher := mock.EventPublisher{}
	do.OverrideValue[repository.FinancialProductRepository](injector, financialProductClientMock)
	do.OverrideValue[event.Publisher](injector, eventPublisher)
	do.OverrideValue[client.Client](injector, temporalClientMock)
	h := do.MustInvoke[*loanOfferInterestHttp.LoanOfferInterestHandler](injector)

	t.Run(
		"user confirm loan package interest existed loan packages success", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					InvestorId: investorId,
					Sub:        investorId,
				},
			)
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
			offerInterest1 := mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           request.LoanRate,
					InterestRate:       decimal.Zero,
					Status:             "PENDING",
					LoanID:             986,
					AssetType:          "UNDERLYING",
				},
			)
			offerInterest2 := mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           request.LoanRate,
					InterestRate:       decimal.Zero,
					Status:             "PENDING",
					LoanID:             168,
					AssetType:          "UNDERLYING",
				},
			)
			offerInterest3 := mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           request.LoanRate,
					InterestRate:       decimal.Zero,
					Status:             "PENDING",
					LoanID:             164,
					AssetType:          "UNDERLYING",
				},
			)
			financialProductClientMock.EXPECT().GetAllAccountDetail(testifyMock.Anything, investorId).Return(
				[]entity.FinancialAccountDetail{}, nil,
			).Maybe()
			financialProductClientMock.EXPECT().AssignLoanPackageOrGetLoanPackageAccountId(
				testifyMock.AnythingOfType("*context.valueCtx"), accountNo, offerInterest1.LoanID,
				entity.AssetTypeUnderlying,
			).Return(456, nil)
			financialProductClientMock.EXPECT().AssignLoanPackageOrGetLoanPackageAccountId(
				testifyMock.AnythingOfType("*context.valueCtx"), accountNo, offerInterest3.LoanID,
				entity.AssetTypeUnderlying,
			).Return(876, nil)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/my-loan-offer-interests",
				loanOfferInterestHttp.InvestorConfirmLoanPackageOfferInterestRequest{LoanPackageOfferInterestIds: []int64{offerInterest1.ID, offerInterest3.ID}},
			)
			h.InvestorConfirmMultipleLoanPackageInterest(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "ok", testhelper.GetString(body, "data"))
			updatedOfferInterest := testhelper.GetLoanOfferInterest(t, db, offerInterest1.ID)
			assert.Equal(t, "PACKAGE_CREATED", updatedOfferInterest.Status)
			updatedOfferInterest2 := testhelper.GetLoanOfferInterest(t, db, offerInterest2.ID)
			assert.Equal(t, "CANCELLED", updatedOfferInterest2.Status)
			updatedOfferInterest3 := testhelper.GetLoanOfferInterest(t, db, offerInterest3.ID)
			assert.Equal(t, "PACKAGE_CREATED", updatedOfferInterest3.Status)
		},
	)

	t.Run(
		"fill loan package data success", func(t *testing.T) {
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
			offerInterest1 := mock.SeedLoanPackageOfferInterest(
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
			offerInterest2 := mock.SeedLoanPackageOfferInterest(
				t, db, model.LoanPackageOfferInterest{
					LoanPackageOfferID: offer.ID,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.Zero,
					InterestRate:       decimal.Zero,
					FeeRate:            decimal.Zero,
					LoanID:             187,
					Status:             "PACKAGE_CREATED",
					AssetType:          "UNDERLYING",
				},
			)
			financialProductClientMock.EXPECT().GetLoanPackageDetail(
				testifyMock.Anything, testifyMock.AnythingOfType("int64"),
			).Return(
				entity.FinancialProductLoanPackage{
					Id:            1,
					Name:          "FinX 40%",
					InitialRate:   decimal.NewFromFloat(0.4),
					InterestRate:  decimal.NewFromFloat(0.1),
					Term:          200,
					BuyingFeeRate: decimal.NewFromFloat(0.012),
				}, nil,
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("GET", "/loan-package-offer-interests/sync-loan-package-data", nil)
			h.FillWithLoanPackageData(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, int64(2), testhelper.GetInt(body, "data"))
			updatedOfferInterest1 := testhelper.GetLoanOfferInterest(t, db, offerInterest1.ID)
			assert.True(t, decimal.NewFromFloat(0.6).Equal(updatedOfferInterest1.LoanRate))
			assert.True(t, decimal.NewFromFloat(0.012).Equal(updatedOfferInterest1.FeeRate))
			updatedOfferInterest2 := testhelper.GetLoanOfferInterest(t, db, offerInterest2.ID)
			assert.True(t, decimal.NewFromFloat(0.6).Equal(updatedOfferInterest2.LoanRate))
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
					Sub:        "investor",
					InvestorId: investorId,
				},
			)
			ginCtx.Request = gintest.MustMakeRequest("POST", "/loan-package-offer-interests/1/admin-cancel", nil)
			ginCtx.AddParam("id", fmt.Sprintf("%d", offerInterest.ID))
			h.InvestorCancelLoanPackageOfferInterest(ginCtx)
			cancelledOfferInterest := testhelper.GetLoanOfferInterest(t, db, offerInterest.ID)
			assert.Equal(t, "CANCELLED", cancelledOfferInterest.Status)
			assert.Equal(t, investorId, cancelledOfferInterest.CancelledBy)
		},
	)

	t.Run(
		"Test create loan contract success", func(t *testing.T) {
			defer truncateData()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code: "HOSE",
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "VIB",
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
					LoanRate:           request.LoanRate,
					InterestRate:       decimal.Zero,
					Status:             "PENDING",
					LoanID:             986,
					AssetType:          "UNDERLYING",
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub:        "investor",
					InvestorId: investorId,
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", fmt.Sprintf("/loan-offer-interests/%d/assign-loan-contract", offerInterest.ID),
				loanOfferInterestHttp.CreateAssignedLoanOfferInterestLoanContractRequest{
					LoanPackageAccountId: 123,
					LoanProductIdRef:     45,
					LoanPackage: entity.FinancialProductLoanPackage{
						Id:            654,
						Name:          "Mo loan package 654",
						InitialRate:   decimal.NewFromFloat(0.43),
						InterestRate:  decimal.NewFromFloat(0.11),
						Term:          180,
						BuyingFeeRate: decimal.NewFromFloat(0.001),
					},
				},
			)
			ginCtx.AddParam("id", fmt.Sprintf("%d", offerInterest.ID))
			h.CreateAssignedLoanOfferInterestLoanContract(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, int64(45), testhelper.GetInt(body, "data", "loanProductIdRef"))
			cancelledOfferInterest := testhelper.GetLoanOfferInterest(t, db, offerInterest.ID)
			assert.Equal(t, "PACKAGE_CREATED", cancelledOfferInterest.Status)
		},
	)
}

func TestLoanOfferInterestHandler_Admin(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	investorId := "0001000115"
	accountNo := "0001000115"
	financialProductClientMock := mock.NewMockFinancialProductRepository(t)
	eventPublisher := mock.EventPublisher{}
	do.OverrideValue[repository.FinancialProductRepository](injector, financialProductClientMock)
	do.OverrideValue[event.Publisher](injector, eventPublisher)
	h := do.MustInvoke[*loanOfferInterestHttp.LoanOfferInterestHandler](injector)

	t.Run(
		"(online flow) investor confirm offer interest", func(t *testing.T) {
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
					LoanID:             loanId,
					LimitAmount:        request.LimitAmount,
					LoanRate:           decimal.NewFromFloat(0.6),
					InterestRate:       decimal.NewFromFloat(0.15),
					FeeRate:            decimal.NewFromFloat(0.012),
					Status:             "PENDING",
					Term:               200,
					AssetType:          "UNDERLYING",
				},
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
					InvestorId: investorId,
					Sub:        investorId,
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/my-loan-offer-interests/confirm",
				loanOfferInterestHttp.InvestorConfirmLoanPackageOfferInterestRequest{
					LoanPackageOfferInterestIds: []int64{offerInterest.ID},
				},
			)
			ginCtx.AddParam("id", fmt.Sprintf("%d", offerInterest.ID))
			h.InvestorConfirmMultipleLoanPackageInterest(ginCtx)
			cancelledOfferInterest := testhelper.GetLoanOfferInterest(t, db, offerInterest.ID)
			assert.Equal(t, "PACKAGE_CREATED", cancelledOfferInterest.Status)
			assert.Equal(t, loanId, cancelledOfferInterest.LoanID)
		},
	)
}
