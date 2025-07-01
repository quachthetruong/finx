package test

import (
	configRepo "financing-offer/internal/config/repository"
	"github.com/go-jet/jet/v2/qrm"
	gohttp "net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/apperrors"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/financialproduct/repository"
	financingApiRepository "financing-offer/internal/core/financing/repository"
	loanPackageRequestRepo "financing-offer/internal/core/loanpackagerequest/repository"
	"financing-offer/internal/core/loanpackagerequest/transport/http"
	marginOperationRepo "financing-offer/internal/core/marginoperation/repository"
	odooServiceRepo "financing-offer/internal/core/odoo_service/repository"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/internal/jwttoken"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestLoggedLoanPackageRequest(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	do.OverrideValue[loanPackageRequestRepo.LoanPackageRequestEventRepository](
		injector, &mock.LoanPackageRequestEventRepository{},
	)
	h := do.MustInvoke[*http.LoanPackageRequestHandler](injector)
	investorId := "0954589234"
	t.Run(
		"save logged request", func(t *testing.T) {
			defer truncateData()
			stockExchange := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code: "HNX",
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					Symbol:          "VIB",
					StockExchangeID: stockExchange.ID,
					AssetType:       "UNDERLYING",
				},
			)
			req := http.LoggedRequestRequest{
				Request: http.CreateLoanPackageRequestUnderlyingRequest{
					SymbolId:    symbol.ID,
					LoanRate:    decimal.NewFromFloat(0.5),
					LimitAmount: decimal.NewFromFloat(1000000),
					AccountNo:   "accountNo",
					Type:        entity.LoanPackageRequestTypeFlexible,
				},
			}
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					InvestorId: investorId,
					Sub:        investorId,
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/my-logged_requests", req,
			)
			h.SaveLoanRateExistedRequest(ginCtx)
			result := recorder.Result()
			assert.Equal(t, 200, result.StatusCode)
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, investorId, testhelper.GetString(body, "data", "investorId"))
		},
	)
}

func TestLoanPackageRequest(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	do.OverrideValue[loanPackageRequestRepo.LoanPackageRequestEventRepository](
		injector, &mock.LoanPackageRequestEventRepository{},
	)
	financialProductRepository := &mock.MockFinancialProductRepository{}
	do.OverrideValue[repository.FinancialProductRepository](
		injector, financialProductRepository,
	)
	do.OverrideValue[financingApiRepository.FinancingRepository](injector, mock.FinancingApiMock{})
	loanPackageRequestHandler := do.MustInvoke[*http.LoanPackageRequestHandler](injector)

	t.Run(
		"get available package of loan package with symbol score greater than stock exchange score",
		func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			groupA := mock.SeedScoreGroup(
				t, db, model.ScoreGroup{
					Code:     "A",
					MinScore: 10,
					MaxScore: 20,
				},
			)
			groupB := mock.SeedScoreGroup(
				t, db, model.ScoreGroup{
					Code:     "B",
					MinScore: 21,
					MaxScore: 30,
				},
			)
			stockExchange := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:         "B",
					ScoreGroupID: &groupA.ID,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					Symbol:          "FPT",
					StockExchangeID: stockExchange.ID,
					AssetType:       "UNDERLYING",
				},
			)
			mock.SeedSymbolScore(
				t, db, model.SymbolScore{
					SymbolID:     symbol.ID,
					Score:        25,
					Type:         string(entity.SymbolScoreTypeManual),
					Status:       string(entity.SymbolScoreStatusActive),
					AffectedFrom: time.Now().Add(-1 * time.Second),
				},
			)
			mock.SeedScoreGroupInterest(
				t, db, model.ScoreGroupInterest{
					LimitAmount:  decimal.NewFromInt(1000000),
					LoanRate:     decimal.NewFromFloat(3.0),
					InterestRate: decimal.NewFromFloat(4.0),
					ScoreGroupID: groupA.ID,
				},
			)
			mock.SeedScoreGroupInterest(
				t, db, model.ScoreGroupInterest{
					LimitAmount:  decimal.NewFromInt(2000000),
					LoanRate:     decimal.NewFromFloat(3.0),
					InterestRate: decimal.NewFromFloat(4.0),
					ScoreGroupID: groupB.ID,
				},
			)
			loanRequest := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:    symbol.ID,
					LimitAmount: decimal.NewFromInt(1000000),
					LoanRate:    decimal.NewFromFloat(3.0),
					Status:      entity.LoanPackageRequestStatusPending.String(),
					AssetType:   "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "", nil,
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(loanRequest.ID, 10),
			}}
			loanPackageRequestHandler.GetAvailablePackages(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 1, testhelper.GetArrayLength(body, "data"))
			assert.Equal(t, groupA.ID, testhelper.GetInt(body, "data", "[0]", "scoreGroupId"))
		},
	)

	t.Run(
		"get available package of score group that return empty", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			groupA := mock.SeedScoreGroup(
				t, db, model.ScoreGroup{
					Code:     "A",
					MinScore: 10,
					MaxScore: 20,
				},
			)
			groupB := mock.SeedScoreGroup(
				t, db, model.ScoreGroup{
					Code:     "B",
					MinScore: 21,
					MaxScore: 30,
				},
			)
			stockExchange := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:         "B",
					ScoreGroupID: &groupB.ID,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					Symbol:          "FPT",
					StockExchangeID: stockExchange.ID,
					AssetType:       "UNDERLYING",
				},
			)
			mock.SeedSymbolScore(
				t, db, model.SymbolScore{
					SymbolID:     symbol.ID,
					Score:        15,
					Type:         string(entity.SymbolScoreTypeManual),
					Status:       string(entity.SymbolScoreStatusActive),
					AffectedFrom: time.Now().Add(-1 * time.Second),
				},
			)
			mock.SeedScoreGroupInterest(
				t, db, model.ScoreGroupInterest{
					LimitAmount:  decimal.NewFromInt(1000000),
					LoanRate:     decimal.NewFromFloat(3.0),
					InterestRate: decimal.NewFromFloat(4.0),
					ScoreGroupID: groupA.ID,
				},
			)
			mock.SeedScoreGroupInterest(
				t, db, model.ScoreGroupInterest{
					LimitAmount:  decimal.NewFromInt(2000000),
					LoanRate:     decimal.NewFromFloat(3.0),
					InterestRate: decimal.NewFromFloat(4.0),
					ScoreGroupID: groupB.ID,
				},
			)
			loanRequest := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:    symbol.ID,
					LimitAmount: decimal.NewFromInt(1000000),
					LoanRate:    decimal.NewFromFloat(3.0),
					Status:      entity.LoanPackageRequestStatusPending.String(),
					AssetType:   "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "", nil,
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(loanRequest.ID, 10),
			}}
			loanPackageRequestHandler.GetAvailablePackages(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 1, testhelper.GetArrayLength(body, "data"))
			assert.Equal(t, groupA.ID, testhelper.GetInt(body, "data", "[0]", "scoreGroupId"))
		},
	)

	t.Run(
		"user request success", func(t *testing.T) {
			defer truncateData()
			accountNo := "accountNo198"
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					InvestorId:  "inv1437",
					Sub:         "investor1",
					CustodyCode: "064C859439",
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
			financialProductRepository.EXPECT().GetAllAccountDetail(ginCtx, "inv1437").Return(
				[]entity.FinancialAccountDetail{
					{
						AccountNo: accountNo,
					},
				}, nil,
			).Maybe()
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/my-loan-package-request", http.CreateLoanPackageRequestUnderlyingRequest{
					SymbolId:    symbol.ID,
					LoanRate:    decimal.NewFromFloat(0.8),
					LimitAmount: decimal.NewFromFloat(2000000),
					AccountNo:   accountNo,
					Type:        entity.LoanPackageRequestTypeFlexible,
				},
			)
			loanPackageRequestHandler.InvestorRequest(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 200, result.StatusCode)
			assert.Equal(t, "PENDING", testhelper.GetString(body, "data", "status"))
		},
	)

	t.Run(
		"user request when accountNo do not belong to investor", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					InvestorId: "investorId143",
					Sub:        "investor1",
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
			financialProductRepository.EXPECT().GetAllAccountDetail(ginCtx, "investorId143").Return(
				[]entity.FinancialAccountDetail{
					{
						AccountNo: "123",
					},
				}, nil,
			).Maybe()
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/my-loan-package-request", http.CreateLoanPackageRequestUnderlyingRequest{
					SymbolId:    symbol.ID,
					LoanRate:    decimal.NewFromFloat(0.8),
					LimitAmount: decimal.NewFromFloat(2000000),
					AccountNo:   "456",
					Type:        entity.LoanPackageRequestTypeFlexible,
				},
			)
			loanPackageRequestHandler.InvestorRequest(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 400, result.StatusCode)
			assert.Equal(t, int64(apperrors.ErrAccountNoInvalid.Code), testhelper.GetInt(body, "code"))
		},
	)

	t.Run(
		"user create derivative request when symbol is not derivative", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					InvestorId: "investorId143",
					Sub:        "investor1",
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
			financialProductRepository.EXPECT().GetAllAccountDetail(ginCtx, "investorId143").Return(
				[]entity.FinancialAccountDetail{
					{
						AccountNo: "876876",
					},
				}, nil,
			).Maybe()
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/my-loan-package-request", http.CreateLoanPackageRequestDerivativeRequest{
					SymbolId:     symbol.ID,
					InitialRate:  decimal.NewFromFloat(0.18),
					ContractSize: 10,
					AccountNo:    "876876",
				},
			)
			loanPackageRequestHandler.InvestorRequestDerivative(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 400, result.StatusCode)
			assert.Equal(t, int64(apperrors.ErrMismatchAssetType.Code), testhelper.GetInt(body, "code"))
		},
	)
}

func TestLoanPackageRequest_Admin(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))

	do.OverrideValue[loanPackageRequestRepo.LoanPackageRequestEventRepository](
		injector, &mock.LoanPackageRequestEventRepository{},
	)
	financialProductRepository := &mock.MockFinancialProductRepository{}
	do.OverrideValue[repository.FinancialProductRepository](
		injector, financialProductRepository,
	)
	configurationRepoMock := mock.NewMockConfigurationPersistenceRepository(t)
	marginOperationRepository := &mock.MockMarginOperationRepository{}
	odooServiceRepository := &mock.MockOdooServiceRepository{}
	do.OverrideValue[marginOperationRepo.MarginOperationRepository](
		injector, marginOperationRepository,
	)
	do.OverrideValue[financingApiRepository.FinancingRepository](injector, mock.FinancingApiMock{})
	do.OverrideValue[configRepo.ConfigurationPersistenceRepository](injector, configurationRepoMock)
	do.OverrideValue[odooServiceRepo.OdooServiceRepository](injector, odooServiceRepository)
	loanPackageRequestHandler := do.MustInvoke[*http.LoanPackageRequestHandler](injector)

	t.Run(
		"admin confirm loan package request with loan package id success", func(t *testing.T) {
			defer truncateData()
			loanId := int64(11990)
			ginCtx, _, recorder := gintest.GetTestContext()
			investorId := "0001000115"
			accountNo := "0001000115"
			financialProductRepository.EXPECT().GetLoanPackageDetail(ginCtx, loanId).Return(
				entity.FinancialProductLoanPackage{
					Id:            loanId,
					Name:          "FinX",
					InitialRate:   decimal.NewFromFloat(0.4),
					InterestRate:  decimal.NewFromFloat(0.15),
					Term:          80,
					BuyingFeeRate: decimal.NewFromFloat(0.0012),
				}, nil,
			)
			financialProductRepository.EXPECT().GetAllAccountDetail(mock2.Anything, investorId).Return(
				[]entity.FinancialAccountDetail{
					{
						AccountNo: accountNo,
					},
				}, nil,
			).Maybe()
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
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/:id/admin-confirm", http.ConfirmLoanPackageRequestRequest{
					OfferedBy: "admin",
					LoanId:    loanId,
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(request.ID, 10),
			}}
			loanPackageRequestHandler.AdminConfirmUserRequest(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "CONFIRMED", testhelper.GetString(body, "data", "status"))
			offer := &model.LoanPackageOffer{}
			err := postgres.SELECT(table.LoanPackageOffer.AllColumns).FROM(
				table.LoanPackageOffer.INNER_JOIN(
					table.LoanPackageRequest,
					table.LoanPackageRequest.ID.EQ(table.LoanPackageOffer.LoanPackageRequestID),
				),
			).Query(db, offer)
			assert.Nil(t, err)
			assert.Equal(t, "admin", offer.OfferedBy)
			assert.Equal(t, entity.FlowTypeDnseOnline.String(), offer.FlowType)
			offerLine := &model.LoanPackageOfferInterest{}
			err = table.LoanPackageOfferInterest.SELECT(table.LoanPackageOfferInterest.AllColumns).Query(db, offerLine)
			assert.Nil(t, err)
			assert.Equal(t, loanId, offerLine.LoanID)
		},
	)

	t.Run(
		"admin confirm loan package request not in pending state", func(t *testing.T) {
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
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/:id/admin-confirm", http.ConfirmLoanPackageRequestRequest{
					OfferedBy: "admin",
					LoanId:    0,
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(request.ID, 10),
			}}
			loanPackageRequestHandler.AdminConfirmUserRequest(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, 400, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, int64(apperrors.ErrInvalidRequestStatus.Code), testhelper.GetInt(body, "code"))
		},
	)

	t.Run(
		"admin confirm derivative loan package with loan id", func(t *testing.T) {
			defer truncateData()
			loanId := int64(1533)
			ginCtx, _, recorder := gintest.GetTestContext()
			financialProductRepository.EXPECT().GetLoanPackageDetail(ginCtx, loanId).Return(
				entity.FinancialProductLoanPackage{
					Id:            loanId,
					Name:          "FinX",
					InitialRate:   decimal.NewFromFloat(0.4),
					InterestRate:  decimal.NewFromFloat(0.15),
					Term:          80,
					BuyingFeeRate: decimal.NewFromFloat(0.0012),
				}, nil,
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
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  model.AssetType_Derivative,
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/:id/admin-confirm", http.ConfirmLoanPackageRequestRequest{
					OfferedBy: "admin",
					LoanId:    loanId,
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(request.ID, 10),
			}}
			loanPackageRequestHandler.AdminConfirmUserRequest(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, 400, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, int64(apperrors.ErrInvalidInput("any").Code), testhelper.GetInt(body, "code"))
		},
	)

	t.Run(
		"admin confirm loan package request not in pending state", func(t *testing.T) {
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
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/:id/admin-confirm", http.ConfirmLoanPackageRequestRequest{
					OfferedBy: "admin",
					LoanId:    0,
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(request.ID, 10),
			}}
			loanPackageRequestHandler.AdminConfirmUserRequest(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, 400, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, int64(apperrors.ErrInvalidRequestStatus.Code), testhelper.GetInt(body, "code"))
		},
	)

	t.Run(
		"admin confirm derivative loan package with loan id", func(t *testing.T) {
			defer truncateData()
			loanId := int64(1533)
			ginCtx, _, recorder := gintest.GetTestContext()
			financialProductRepository.EXPECT().GetLoanPackageDetail(ginCtx, loanId).Return(
				entity.FinancialProductLoanPackage{
					Id:            loanId,
					Name:          "FinX",
					InitialRate:   decimal.NewFromFloat(0.4),
					InterestRate:  decimal.NewFromFloat(0.15),
					Term:          80,
					BuyingFeeRate: decimal.NewFromFloat(0.0012),
				}, nil,
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
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  model.AssetType_Derivative,
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/:id/admin-confirm", http.ConfirmLoanPackageRequestRequest{
					OfferedBy: "admin",
					LoanId:    loanId,
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(request.ID, 10),
			}}
			loanPackageRequestHandler.AdminConfirmUserRequest(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, 400, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, int64(apperrors.ErrInvalidInput("any").Code), testhelper.GetInt(body, "code"))
		},
	)

	t.Run(
		"admin confirm loan package request without loan package id success", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			investorId := "0001000115"
			accountNo := "0001000115"
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "user1",
				},
			)
			financialProductRepository.EXPECT().GetAllAccountDetail(ginCtx, investorId).Return(
				[]entity.FinancialAccountDetail{
					{
						AccountNo: accountNo,
					},
				}, nil,
			).Maybe()
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
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/:id/admin-confirm", http.ConfirmLoanPackageRequestRequest{},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(request.ID, 10),
			}}
			loanPackageRequestHandler.AdminConfirmUserRequest(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "CONFIRMED", testhelper.GetString(body, "data", "status"))
			offer := &model.LoanPackageOffer{}
			err := postgres.SELECT(table.LoanPackageOffer.AllColumns).FROM(
				table.LoanPackageOffer.INNER_JOIN(
					table.LoanPackageRequest,
					table.LoanPackageRequest.ID.EQ(table.LoanPackageOffer.LoanPackageRequestID),
				),
			).Query(db, offer)
			assert.Nil(t, err)
			assert.Equal(t, "user1", offer.OfferedBy)
			assert.Equal(t, entity.FLowTypeDnseOffline.String(), offer.FlowType)
		},
	)

	t.Run(
		"admin cancel loan package request with alternative options", func(t *testing.T) {
			defer truncateData()
			investorId := "0001000115"
			accountNo := "0001000115"
			loanPackage1 := entity.FinancialProductLoanPackage{
				Id:            1132,
				Name:          "FinX",
				InitialRate:   decimal.NewFromFloat(0.4),
				InterestRate:  decimal.NewFromFloat(0.15),
				Term:          80,
				BuyingFeeRate: decimal.NewFromFloat(0.0012),
			}

			loanPackage2 := entity.FinancialProductLoanPackage{
				Id:            6578,
				Name:          "FinX 80%",
				InitialRate:   decimal.NewFromFloat(0.2),
				InterestRate:  decimal.NewFromFloat(0.2),
				Term:          30,
				BuyingFeeRate: decimal.NewFromFloat(0.0012),
			}
			ginCtx, _, recorder := gintest.GetTestContext()
			financialProductRepository.EXPECT().GetLoanPackageDetails(
				mock2.Anything, []int64{loanPackage1.Id, loanPackage2.Id},
			).Return([]entity.FinancialProductLoanPackage{loanPackage1, loanPackage2}, nil)
			financialProductRepository.EXPECT().GetAllAccountDetail(ginCtx, investorId).Return(
				[]entity.FinancialAccountDetail{
					{
						AccountNo: accountNo,
					},
				}, nil,
			).Maybe()
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
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/:id/cancel", http.CancelLoanPackageRequestRequest{
					LoanIds:   []int64{loanPackage1.Id, loanPackage2.Id},
					OfferedBy: "admin",
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(request.ID, 10),
			}}
			loanPackageRequestHandler.AdminCancelLoanRequest(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(
				t, entity.LoanPackageRequestStatusConfirmed.String(), testhelper.GetString(body, "data", "status"),
			)
			createdOfferInterest := make([]model.LoanPackageOfferInterest, 0, 3)
			if err := postgres.SELECT(table.LoanPackageOfferInterest.AllColumns).FROM(
				table.LoanPackageOfferInterest.
					INNER_JOIN(
						table.LoanPackageOffer,
						table.LoanPackageOffer.ID.EQ(table.LoanPackageOfferInterest.LoanPackageOfferID),
					).
					INNER_JOIN(
						table.LoanPackageRequest,
						table.LoanPackageRequest.ID.EQ(table.LoanPackageOffer.LoanPackageRequestID),
					),
			).WHERE(table.LoanPackageOffer.LoanPackageRequestID.EQ(postgres.Int64(request.ID))).Query(
				db, &createdOfferInterest,
			); err != nil {
				t.Error(err)
			}
			assert.Equal(t, 3, len(createdOfferInterest))
			for _, offerInterest := range createdOfferInterest {
				if offerInterest.LoanID == 0 {
					assert.Equal(t, "CANCELLED", offerInterest.Status)
				}
			}
		},
	)

	t.Run(
		"cancel loan request with no alternative option success", func(t *testing.T) {
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
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/:id/cancel", http.CancelLoanPackageRequestRequest{
					OfferedBy: "admin",
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(request.ID, 10),
			}}
			financialProductRepository.EXPECT().GetAllAccountDetail(ginCtx, request.InvestorID).Return(
				nil, nil,
			).Maybe()
			loanPackageRequestHandler.AdminCancelLoanRequest(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(
				t, entity.LoanPackageRequestStatusConfirmed.String(),
				testhelper.GetString(body, "data", "status"),
			)
		},
	)

	t.Run(
		"get loan package requests with paging success", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HNX",
					MinScore: 0,
					MaxScore: 90,
				},
			)
			symbol1 := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "HSB",
					AssetType:       "UNDERLYING",
				},
			)
			symbol2 := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "VIB",
					AssetType:       "UNDERLYING",
				},
			)
			mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					ID:         1,
					SymbolID:   symbol2.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)

			request2 := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					ID:         2,
					SymbolID:   symbol1.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			request3 := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					ID:         3,
					SymbolID:   symbol1.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.9),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					ID:         4,
					SymbolID:   symbol1.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.5),
					Type:       "GUARANTEED",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET",
				"?page[size]=2&page[number]=1&sort=-id,status&symbols=HSB&types=FLEXIBLE&types=GUARANTEED&loanPercentFrom=70&statuses=PENDING",
				nil,
			)
			loanPackageRequestHandler.GetAll(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 2, testhelper.GetArrayLength(body, "data"))
			assert.Equal(t, request3.ID, testhelper.GetInt(body, "data", "[0]", "id"))
			assert.Equal(t, request2.ID, testhelper.GetInt(body, "data", "[1]", "id"))
			assert.Equal(t, int64(2), testhelper.GetInt(body, "metaData", "total"))
			assert.Equal(t, int64(1), testhelper.GetInt(body, "metaData", "totalPages"))
		},
	)

	t.Run(
		"admin get request by id success", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HNX",
					MinScore: 0,
					MaxScore: 90,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "HSB",
					AssetType:       "UNDERLYING",
				},
			)
			request := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/id", nil,
			)
			ginCtx.AddParam("id", strconv.FormatInt(request.ID, 10))
			loanPackageRequestHandler.AdminGetById(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, request.ID, testhelper.GetInt(body, "data", "id"))
			assert.Equal(t, request.Status, testhelper.GetString(body, "data", "status"))
		},
	)

	t.Run(
		"test admin get by id when not found", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/id", nil,
			)
			ginCtx.AddParam("id", "1")
			loanPackageRequestHandler.AdminGetById(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, 404, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "not found resources", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"investor get by id", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HNX",
					MinScore: 0,
					MaxScore: 90,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "HSB",
					AssetType:       "UNDERLYING",
				},
			)
			request := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/id", nil,
			)
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					InvestorId: "0001000115",
				},
			)
			ginCtx.AddParam("id", strconv.FormatInt(request.ID, 10))
			loanPackageRequestHandler.InvestorGetById(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, request.ID, testhelper.GetInt(body, "data", "id"))
			assert.Equal(t, request.Status, testhelper.GetString(body, "data", "status"))
		},
	)

	t.Run(
		"investor get by id when not found", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/id", nil,
			)
			ginCtx.AddParam("id", "1")
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					InvestorId: "0001000115",
				},
			)
			loanPackageRequestHandler.InvestorGetById(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, 404, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "not found resources", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"investor get all requests", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HNX",
					MinScore: 0,
					MaxScore: 90,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "VIB",
					AssetType:       "UNDERLYING",
				},
			)
			request1 := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			request2 := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "", nil,
			)
			loanPackageRequestHandler.GetAll(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 2, testhelper.GetArrayLength(body, "data"))
			assert.Equal(t, request1.Status, testhelper.GetString(body, "data", "[0]", "status"))
			assert.Equal(t, request2.Status, testhelper.GetString(body, "data", "[1]", "status"))
		},
	)

	t.Run(
		"cancel all loan package with symbol id", func(t *testing.T) {
			defer truncateData()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HNX",
					MinScore: 0,
					MaxScore: 90,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "VIB",
					AssetType:       "UNDERLYING",
				},
			)
			mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			request := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "test@dnse.com",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/symbol/1", nil,
			)
			ginCtx.AddParam("id", strconv.FormatInt(symbol.ID, 10))
			financialProductRepository.EXPECT().GetAllAccountDetail(ginCtx, request.InvestorID).Return(
				nil, nil,
			).Maybe()
			loanPackageRequestHandler.CancelAllLoanPackageRequestBySymbolId(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(
				t, entity.LoanPackageRequestStatusConfirmed.String(),
				testhelper.GetString(body, "data", "[0]", "status"),
			)
			assert.Equal(
				t, entity.LoanPackageRequestStatusConfirmed.String(),
				testhelper.GetString(body, "data", "[1]", "status"),
			)
		},
	)
	t.Run(
		"cancel all loan package invalid id", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "test@dnse.com",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/symbol/", nil,
			)
			ginCtx.AddParam("id", "invalid")
			loanPackageRequestHandler.CancelAllLoanPackageRequestBySymbolId(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "id invalid", testhelper.GetString(body, "error"))
		},
	)
	//t.Run(
	//	"cancel all loan package error", func(t *testing.T) {
	//		se := mock.SeedStockExchange(
	//			t, db, model.StockExchange{
	//				Code:     "HNX",
	//				MinScore: 0,
	//				MaxScore: 90,
	//			},
	//		)
	//		symbol := mock.SeedSymbol(
	//			t, db, model.Symbol{
	//				StockExchangeID: se.ID,
	//				Symbol:          "VIB",
	//				AssetType:       "UNDERLYING",
	//			},
	//		)
	//		mock.SeedLoanPackageRequest(
	//			t, db, model.LoanPackageRequest{
	//				SymbolID:   symbol.ID,
	//				InvestorID: "0001000115",
	//				AccountNo:  "0001000115",
	//				LoanRate:   decimal.NewFromFloat(0.7),
	//				Type:       "FLEXIBLE",
	//				Status:     "PENDING",
	//				AssetType:  "UNDERLYING",
	//			},
	//		)
	//		request := mock.SeedLoanPackageRequest(
	//			t, db, model.LoanPackageRequest{
	//				SymbolID:   symbol.ID,
	//				InvestorID: "0001000115",
	//				AccountNo:  "0001000115",
	//				LoanRate:   decimal.NewFromFloat(0.7),
	//				Type:       "FLEXIBLE",
	//				Status:     "PENDING",
	//				AssetType:  "UNDERLYING",
	//			},
	//		)
	//		defer func() {
	//			truncateData()
	//			db, tearDownDb, truncateData = dbtest.NewDb(t)
	//			//injector = testhelper.NewInjector(testhelper.WithDb(db))
	//			do.OverrideValue[*sql.DB](injector, db)
	//			//do.OverrideValue[loanPackageRequestRepo.LoanPackageRequestEventRepository](
	//			//	injector, &mock.LoanPackageRequestEventRepository{},
	//			//)
	//			//financialProductRepository = &mock.MockFinancialProductRepository{}
	//			//do.OverrideValue[repository.FinancialProductRepository](
	//			//	injector, financialProductRepository,
	//			//)
	//			//do.OverrideValue[financingApiRepository.FinancingRepository](injector, mock.FinancingApiMock{})
	//			//loanPackageRequestHandler = do.MustInvoke[*http.LoanPackageRequestHandler](injector)
	//		}()
	//		err := db.Close()
	//		assert.Nil(t, err)
	//		ginCtx, _, recorder := gintest.GetTestContext()
	//		ginCtx.Set(
	//			appcontext.UserInformation, &jwttoken.AdminClaims{
	//				Sub: "test@dnse.com",
	//			},
	//		)
	//		ginCtx.Request = gintest.MustMakeRequest(
	//			"POST", "/symbol/1", nil,
	//		)
	//		ginCtx.AddParam("id", strconv.FormatInt(symbol.ID, 10))
	//		financialProductRepository.EXPECT().GetAllAccountDetail(ginCtx, request.InvestorID).Return(
	//			nil, nil,
	//		).Maybe()
	//		loanPackageRequestHandler.CancelAllLoanPackageRequestBySymbolId(ginCtx)
	//		result := recorder.Result()
	//		defer assert.Nil(t, result.Body.Close())
	//		body := gintest.ExtractBody(result.Body)
	//		assert.Equal(t, "an error happened, please try again later", testhelper.GetString(body, "error"))
	//	},
	//)

	t.Run(
		"Admin confirm with new loan package", func(t *testing.T) {
			defer truncateData()
			loanRateId := int64(7071)
			investorId := "0001000115"
			accountNo := "0001000115"

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
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       1,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test1",
					InterestRate:             decimal.NewFromFloat(0.12),
					InterestBasis:            31,
					Term:                     1,
					PoolIDRef:                1005,
					OverdueInterest:          decimal.NewFromFloat(0.66),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.43),
				},
			)
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       5,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test2",
					InterestRate:             decimal.NewFromFloat(0.53),
					InterestBasis:            30,
					Term:                     2,
					PoolIDRef:                1004,
					OverdueInterest:          decimal.NewFromFloat(0.12),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.442),
				},
			)
			odooServiceRepository.EXPECT().SendLoanApprovalRequest(mock2.Anything).Return(nil)

			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/loan-package-requests/:id/submissions", entity.SubmissionSheetShorten{
					Metadata: entity.SubmissionSheetMetadata{
						ActionType:           "APPROVE",
						FlowType:             "002",
						ProposeType:          "NEW_LOAN_PACKAGE",
						LoanPackageRequestId: request.ID,
						Status:               "SUBMITTED",
					},
					Detail: entity.SubmissionSheetDetailShorten{
						LoanPackageRateId: loanRateId,
						FirmBuyingFee:     decimal.NewFromFloat(0.234),
						FirmSellingFee:    decimal.NewFromFloat(0.236),
						TransferFee:       decimal.NewFromFloat(234.5),
						LoanPolicies: []entity.LoanPolicyShorten{
							{
								LoanPolicyTemplateId:   1,
								InitialRateForWithdraw: decimal.NewFromFloat(0.4606),
								InitialRate:            decimal.NewFromFloat(0.47),
							},
							{
								LoanPolicyTemplateId:   5,
								InitialRateForWithdraw: decimal.NewFromFloat(0.0294),
								InitialRate:            decimal.NewFromFloat(0.03),
							},
						},
						Comment: "description",
					},
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(request.ID, 10),
			}}
			configurationRepoMock.EXPECT().GetLoanRateConfiguration(mock2.Anything).Return(entity.LoanRateConfiguration{
				Ids: []int64{loanRateId},
			}, nil)
			financialProductRepository.EXPECT().GetLoanRateDetail(mock2.Anything, mock2.Anything).Return(
				entity.LoanRate{
					Id:                     loanRateId,
					Name:                   "ky quy 50%",
					InitialRate:            decimal.NewFromFloat(0.5),
					InitialRateForWithdraw: decimal.NewFromFloat(0.51),
					MaintenanceRate:        decimal.NewFromFloat(0.4),
					LiquidRate:             decimal.NewFromFloat(0.3),
				}, nil,
			)
			marginOperationRepository.EXPECT().GetMarginPoolsByIds(ginCtx, []int64{1005, 1004}).Return(

				[]entity.MarginPool{
					{
						Id:          1004,
						PoolGroupId: 1000,
						Name:        "Pool mặc định",
						Type:        "DEFAULT",
						Status:      "ACTIVE",
					},
					{
						Id:          1005,
						PoolGroupId: 1000,
						Name:        "Pool high margin",
						Type:        "RESERVE",
						Status:      "ACTIVE",
					},
				}, nil,
			)

			marginOperationRepository.EXPECT().GetMarginPoolGroupsByIds(ginCtx, []int64{1000, 1000}).Return(
				[]entity.MarginPoolGroup{
					{
						Id:     1000,
						Name:   "POOL TỔNG",
						Status: "ACTIVE",
						Type:   "MARGIN",
						Source: "DNSE",
					},
				}, nil,
			)
			financialProductRepository.EXPECT().GetAllAccountDetail(mock2.Anything, mock2.Anything).Return(
				[]entity.FinancialAccountDetail{
					{
						AccountNo: accountNo,
					},
				}, nil,
			).Maybe()
			loanPackageRequestHandler.AdminConfirmWithNewLoanPackage(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "PENDING", testhelper.GetString(body, "data", "status"))

			submissionSheetMetadata := &model.SubmissionSheetMetadata{}
			err := postgres.SELECT(table.SubmissionSheetMetadata.AllColumns).FROM(table.SubmissionSheetMetadata).WHERE(table.SubmissionSheetMetadata.LoanPackageRequestID.EQ(postgres.Int64(request.ID))).Query(
				db, submissionSheetMetadata,
			)
			assert.Nil(t, err)

			assert.Equal(t, submissionSheetMetadata.Status, entity.SubmissionSheetStatusSubmitted.String())
		},
	)
	t.Run(
		"Admin send other propose with new loan package with existed draft", func(t *testing.T) {
			defer truncateData()
			loanRateId := int64(7071)
			investorId := "0001000115"
			accountNo := "0001000115"

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
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       1,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test1",
					InterestRate:             decimal.NewFromFloat(0.12),
					InterestBasis:            31,
					Term:                     1,
					PoolIDRef:                1005,
					OverdueInterest:          decimal.NewFromFloat(0.66),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.43),
				},
			)
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       5,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test2",
					InterestRate:             decimal.NewFromFloat(0.53),
					InterestBasis:            30,
					Term:                     2,
					PoolIDRef:                1004,
					OverdueInterest:          decimal.NewFromFloat(0.12),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.442),
				},
			)
			submissionSheetMetadata := mock.SeedSubmissionSheetMetadata(
				t, db, model.SubmissionSheetMetadata{
					ID:                   1,
					LoanPackageRequestID: request.ID,
					Status:               "DRAFT",
					ActionType:           "REJECT_AND_SEND_OTHER_PROPOSAL",
					FlowType:             "002",
					ProposeType:          "NEW_LOAN_PACKAGE",
				},
			)
			submissionSheetDetail := mock.SeedSubmissionSheetDetail(
				t, db, model.SubmissionSheetDetail{
					ID:                1,
					SubmissionSheetID: submissionSheetMetadata.ID,
					LoanRate:          "{\"id\": 7071, \"name\": \"ky quy 30%\", \"liquidRate\": 0.2, \"initialRate\": 0.3, \"interestRate\": 0, \"maintenanceRate\": 0.25, \"initialRateForWithdraw\": 0.31}",
					FirmBuyingFee:     decimal.NewFromFloat(0.234),
					FirmSellingFee:    decimal.NewFromFloat(0.236),
					TransferFee:       decimal.NewFromFloat(234.5),
					LoanPolicies:      "[\n  {\n    \"loanPolicyTemplateName\": \"truong.quach\",\n    \"term\": 3,\n    \"source\": \"C\",\n    \"createdAt\": \"2024-03-16T16:26:23.875309Z\",\n    \"poolIdRef\": 1001,\n    \"updatedAt\": \"2024-03-16T16:26:23.875309Z\",\n    \"updatedBy\": \"dnse.admin@dnse.com.vn\",\n    \"initialRate\": 0.7,\n    \"interestRate\": 0.245,\n    \"interestBasis\": 24,\n    \"overdueInterest\": 0.435,\n    \"allowEarlyPayment\": true,\n    \"preferentialPeriod\": 365,\n    \"allowExtendLoanTerm\": true,\n    \"loanPolicyTemplateId\": 1,\n    \"initialRateForWithdraw\": 0.69,\n    \"preferentialInterestRate\": 0.325\n  }\n]",
					Comment:           "test lan 1",
				},
			)

			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/loan-package-requests/:id/submissions", entity.SubmissionSheetShorten{
					Metadata: entity.SubmissionSheetMetadata{
						Id:                   submissionSheetMetadata.ID,
						ActionType:           "REJECT_AND_SEND_OTHER_PROPOSAL",
						FlowType:             "002",
						ProposeType:          "NEW_LOAN_PACKAGE",
						LoanPackageRequestId: request.ID,
						Status:               "SUBMITTED",
					},
					Detail: entity.SubmissionSheetDetailShorten{
						Id:                submissionSheetDetail.ID,
						SubmissionSheetId: submissionSheetMetadata.ID,
						LoanPackageRateId: loanRateId,
						FirmBuyingFee:     decimal.NewFromFloat(0.234),
						FirmSellingFee:    decimal.NewFromFloat(0.236),
						TransferFee:       decimal.NewFromFloat(234.5),
						LoanPolicies: []entity.LoanPolicyShorten{
							{
								LoanPolicyTemplateId:   1,
								InitialRateForWithdraw: decimal.NewFromFloat(0.4606),
								InitialRate:            decimal.NewFromFloat(0.47),
							},
							{
								LoanPolicyTemplateId:   5,
								InitialRateForWithdraw: decimal.NewFromFloat(0.0294),
								InitialRate:            decimal.NewFromFloat(0.03),
							},
						},
						Comment: "description",
					},
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(request.ID, 10),
			}}

			financialProductRepository.EXPECT().GetLoanRateDetail(mock2.Anything, mock2.Anything).Return(
				entity.LoanRate{
					Id:                     loanRateId,
					Name:                   "ky quy 50%",
					InitialRate:            decimal.NewFromFloat(0.5),
					InitialRateForWithdraw: decimal.NewFromFloat(0.51),
					MaintenanceRate:        decimal.NewFromFloat(0.4),
					LiquidRate:             decimal.NewFromFloat(0.3),
				}, nil,
			)

			configurationRepoMock.EXPECT().GetLoanRateConfiguration(mock2.Anything).Return(entity.LoanRateConfiguration{
				Ids: []int64{loanRateId},
			}, nil)

			marginOperationRepository.EXPECT().GetMarginPoolsByIds(ginCtx, []int64{1005, 1004}).Return(

				[]entity.MarginPool{
					{
						Id:          1004,
						PoolGroupId: 1000,
						Name:        "Pool mặc định",
						Type:        "DEFAULT",
						Status:      "ACTIVE",
					},
					{
						Id:          1005,
						PoolGroupId: 1000,
						Name:        "Pool high margin",
						Type:        "RESERVE",
						Status:      "ACTIVE",
					},
				}, nil,
			)

			marginOperationRepository.EXPECT().GetMarginPoolGroupsByIds(ginCtx, []int64{1000, 1000}).Return(
				[]entity.MarginPoolGroup{
					{
						Id:     1000,
						Name:   "POOL TỔNG",
						Status: "ACTIVE",
						Type:   "MARGIN",
						Source: "DNSE",
					},
				}, nil,
			)

			odooServiceRepository.EXPECT().SendLoanApprovalRequest(mock2.Anything).Return(nil)

			loanPackageRequestHandler.AdminDeclineLoanRequestWithNewLoanPackage(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, entity.LoanPackageRequestStatusPending.String(), testhelper.GetString(body, "data", "status"))
			offer := &model.LoanPackageOffer{}
			err := postgres.SELECT(table.LoanPackageOffer.AllColumns).FROM(
				table.LoanPackageOffer.INNER_JOIN(
					table.LoanPackageRequest,
					table.LoanPackageRequest.ID.EQ(table.LoanPackageOffer.LoanPackageRequestID),
				),
			).Query(db, offer)
			assert.ErrorIs(t, err, qrm.ErrNoRows)
			offerLine := &model.LoanPackageOfferInterest{}
			err = table.LoanPackageOfferInterest.SELECT(table.LoanPackageOfferInterest.AllColumns).Query(db, offerLine)
			assert.ErrorIs(t, err, qrm.ErrNoRows)

			newSubmissionSheetMetadata := &model.SubmissionSheetMetadata{}
			err = postgres.SELECT(table.SubmissionSheetMetadata.AllColumns).FROM(table.SubmissionSheetMetadata).WHERE(table.SubmissionSheetMetadata.LoanPackageRequestID.EQ(postgres.Int64(request.ID))).Query(
				db, newSubmissionSheetMetadata,
			)
			assert.Nil(t, err)

			newSubmissionSheetDetail := &model.SubmissionSheetDetail{}
			err = postgres.SELECT(table.SubmissionSheetDetail.AllColumns).FROM(table.SubmissionSheetDetail).WHERE(table.SubmissionSheetDetail.SubmissionSheetID.EQ(postgres.Int64(submissionSheetMetadata.ID))).Query(
				db, newSubmissionSheetDetail,
			)
			assert.Nil(t, err)
			assert.Equal(t, newSubmissionSheetMetadata.LoanPackageRequestID, request.ID)
			assert.Equal(t, newSubmissionSheetMetadata.Status, entity.SubmissionSheetStatusSubmitted.String())
		},
	)

	t.Run(
		"Admin rejected with existed draft", func(t *testing.T) {
			defer truncateData()
			loanRateId := int64(7071)
			investorId := "0001000115"
			accountNo := "0001000115"

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
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       1,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test1",
					InterestRate:             decimal.NewFromFloat(0.12),
					InterestBasis:            31,
					Term:                     1,
					PoolIDRef:                1005,
					OverdueInterest:          decimal.NewFromFloat(0.66),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.43),
				},
			)
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       5,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test2",
					InterestRate:             decimal.NewFromFloat(0.53),
					InterestBasis:            30,
					Term:                     2,
					PoolIDRef:                1004,
					OverdueInterest:          decimal.NewFromFloat(0.12),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.442),
				},
			)
			submissionSheetMetadata := mock.SeedSubmissionSheetMetadata(
				t, db, model.SubmissionSheetMetadata{
					ID:                   1,
					LoanPackageRequestID: request.ID,
					Status:               "DRAFT",
					ActionType:           "REJECTED",
					FlowType:             "002",
					ProposeType:          "NEW_LOAN_PACKAGE",
				},
			)
			mock.SeedSubmissionSheetDetail(
				t, db, model.SubmissionSheetDetail{
					ID:                1,
					SubmissionSheetID: submissionSheetMetadata.ID,
					LoanRate:          "{\"id\": 7071, \"name\": \"ky quy 30%\", \"liquidRate\": 0.2, \"initialRate\": 0.3, \"interestRate\": 0, \"maintenanceRate\": 0.25, \"initialRateForWithdraw\": 0.31}",
					FirmBuyingFee:     decimal.NewFromFloat(0.234),
					FirmSellingFee:    decimal.NewFromFloat(0.236),
					TransferFee:       decimal.NewFromFloat(234.5),
					LoanPolicies:      "[\n  {\n    \"loanPolicyTemplateName\": \"truong.quach\",\n    \"term\": 3,\n    \"source\": \"C\",\n    \"createdAt\": \"2024-03-16T16:26:23.875309Z\",\n    \"poolIdRef\": 1001,\n    \"updatedAt\": \"2024-03-16T16:26:23.875309Z\",\n    \"updatedBy\": \"dnse.admin@dnse.com.vn\",\n    \"initialRate\": 0.7,\n    \"interestRate\": 0.245,\n    \"interestBasis\": 24,\n    \"overdueInterest\": 0.435,\n    \"allowEarlyPayment\": true,\n    \"preferentialPeriod\": 365,\n    \"allowExtendLoanTerm\": true,\n    \"loanPolicyTemplateId\": 1,\n    \"initialRateForWithdraw\": 0.69,\n    \"preferentialInterestRate\": 0.325\n  }\n]",
					Comment:           "test lan 1",
				},
			)
			loanPackage1 := entity.FinancialProductLoanPackage{
				Id:            1132,
				Name:          "FinX",
				InitialRate:   decimal.NewFromFloat(0.4),
				InterestRate:  decimal.NewFromFloat(0.15),
				Term:          80,
				BuyingFeeRate: decimal.NewFromFloat(0.0012),
			}

			loanPackage2 := entity.FinancialProductLoanPackage{
				Id:            6578,
				Name:          "FinX 80%",
				InitialRate:   decimal.NewFromFloat(0.2),
				InterestRate:  decimal.NewFromFloat(0.2),
				Term:          30,
				BuyingFeeRate: decimal.NewFromFloat(0.0012),
			}
			ginCtx, _, recorder := gintest.GetTestContext()

			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/loan-package-requests/:id/cancel", http.CancelLoanPackageRequestRequest{
					LoanIds:   []int64{loanPackage1.Id, loanPackage2.Id},
					OfferedBy: "admin",
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(request.ID, 10),
			}}

			financialProductRepository.EXPECT().GetLoanPackageDetails(
				mock2.Anything, []int64{loanPackage1.Id, loanPackage2.Id},
			).Return([]entity.FinancialProductLoanPackage{loanPackage1, loanPackage2}, nil)
			financialProductRepository.EXPECT().GetLoanRateDetail(ginCtx, loanRateId).Return(
				entity.LoanRate{
					Id:                     loanRateId,
					Name:                   "ky quy 50%",
					InitialRate:            decimal.NewFromFloat(0.5),
					InitialRateForWithdraw: decimal.NewFromFloat(0.51),
					MaintenanceRate:        decimal.NewFromFloat(0.4),
					LiquidRate:             decimal.NewFromFloat(0.3),
				}, nil,
			)
			marginOperationRepository.EXPECT().GetMarginPoolsByIds(ginCtx, []int64{1005, 1004}).Return(

				[]entity.MarginPool{
					{
						Id:          1005,
						PoolGroupId: 1000,
						Name:        "Pool mặc định",
						Type:        "DEFAULT",
						Status:      "ACTIVE",
					},
					{
						Id:          1005,
						PoolGroupId: 1000,
						Name:        "Pool high margin",
						Type:        "RESERVE",
						Status:      "ACTIVE",
					},
				}, nil,
			)

			marginOperationRepository.EXPECT().GetMarginPoolGroupsByIds(ginCtx, []int64{1000, 1000}).Return(
				[]entity.MarginPoolGroup{
					{
						Id:     1000,
						Name:   "POOL TỔNG",
						Status: "ACTIVE",
						Type:   "MARGIN",
						Source: "DNSE",
					},
				}, nil,
			)
			financialProductRepository.EXPECT().GetAllAccountDetail(mock2.Anything, investorId).Return(
				[]entity.FinancialAccountDetail{
					{
						AccountNo: accountNo,
					},
				}, nil,
			).Maybe()
			loanPackageRequestHandler.AdminCancelLoanRequest(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "CONFIRMED", testhelper.GetString(body, "data", "status"))
			offer := &model.LoanPackageOffer{}
			err := postgres.SELECT(table.LoanPackageOffer.AllColumns).FROM(
				table.LoanPackageOffer.INNER_JOIN(
					table.LoanPackageRequest,
					table.LoanPackageRequest.ID.EQ(table.LoanPackageOffer.LoanPackageRequestID),
				),
			).Query(db, offer)
			assert.Nil(t, err)
			assert.Equal(t, entity.FlowTypeDnseOnline.String(), offer.FlowType)
			offerLines := make([]model.LoanPackageOfferInterest, 0)
			err = table.LoanPackageOfferInterest.SELECT(table.LoanPackageOfferInterest.AllColumns).WHERE(
				table.LoanPackageOfferInterest.LoanPackageOfferID.EQ(postgres.Int64(offer.ID)),
			).Query(db, &offerLines)
			assert.Nil(t, err)

			assert.Equal(t, 3, len(offerLines))
			assert.Equal(t, "CANCELLED", offerLines[len(offerLines)-1].Status)
			assert.Equal(t, int64(1132), offerLines[0].LoanID)
			assert.Equal(t, "PENDING", offerLines[1].Status)

			newSubmissionSheetMetadata := &model.SubmissionSheetMetadata{}
			err = postgres.SELECT(table.SubmissionSheetMetadata.AllColumns).FROM(table.SubmissionSheetMetadata).WHERE(table.SubmissionSheetMetadata.LoanPackageRequestID.EQ(postgres.Int64(request.ID))).Query(
				db, newSubmissionSheetMetadata,
			)
			assert.Nil(t, err)
			assert.Equal(t, "DRAFT", newSubmissionSheetMetadata.Status)
		},
	)

	t.Run(
		"Get submission sheet by request id", func(t *testing.T) {
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
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       1,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test1",
					InterestRate:             decimal.NewFromFloat(0.12),
					InterestBasis:            31,
					Term:                     1,
					PoolIDRef:                1001,
					OverdueInterest:          decimal.NewFromFloat(0.66),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.43),
				},
			)
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       5,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test2",
					InterestRate:             decimal.NewFromFloat(0.53),
					InterestBasis:            30,
					Term:                     2,
					PoolIDRef:                1000,
					OverdueInterest:          decimal.NewFromFloat(0.12),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.442),
				},
			)

			submissionSheetMetadata := mock.SeedSubmissionSheetMetadata(t, db, model.SubmissionSheetMetadata{
				ID:                   1,
				LoanPackageRequestID: request.ID,
				Status:               "DRAFT",
				ActionType:           "REJECT_AND_SEND_OTHER_PROPOSAL",
				FlowType:             "002",
				ProposeType:          "NEW_LOAN_PACKAGE",
			})
			mock.SeedSubmissionSheetDetail(t, db, model.SubmissionSheetDetail{
				ID:                1,
				SubmissionSheetID: submissionSheetMetadata.ID,
				LoanRate:          "{\"id\": 7071, \"name\": \"ky quy 30%\", \"liquidRate\": 0.2, \"initialRate\": 0.3, \"interestRate\": 0, \"maintenanceRate\": 0.25, \"initialRateForWithdraw\": 0.31}",
				FirmBuyingFee:     decimal.NewFromFloat(0.234),
				FirmSellingFee:    decimal.NewFromFloat(0.236),
				TransferFee:       decimal.NewFromFloat(234),
				LoanPolicies:      "[\n  {\n    \"loanPolicyTemplateName\": \"truong.quach\",\n    \"term\": 3,\n    \"source\": \"C\",\n    \"createdAt\": \"2024-03-16T16:26:23.875309Z\",\n    \"poolIdRef\": 1001,\n    \"updatedAt\": \"2024-03-16T16:26:23.875309Z\",\n    \"updatedBy\": \"dnse.admin@dnse.com.vn\",\n    \"initialRate\": 0.7,\n    \"interestRate\": 0.245,\n    \"interestBasis\": 24,\n    \"overdueInterest\": 0.435,\n    \"allowEarlyPayment\": true,\n    \"preferentialPeriod\": 365,\n    \"allowExtendLoanTerm\": true,\n    \"loanPolicyTemplateId\": 1,\n    \"initialRateForWithdraw\": 0.69,\n    \"preferentialInterestRate\": 0.325\n  }\n]",
				Comment:           "test lan 1",
			})
			ginCtx.Request = gintest.MustMakeRequest("GET", "/api/v1/submission-sheets", nil)
			ginCtx.AddParam("id", strconv.Itoa(int(request.ID)))

			loanPackageRequestHandler.AdminGetLatestSubmissionSheet(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "DRAFT", testhelper.GetString(body, "data", "metadata", "status"))
			assert.Equal(t, "C", testhelper.GetString(body, "data", "detail", "loanPolicies", "[0]", "source"))
		},
	)
}

func TestLoanPackageRequestHandler_GetAllUnderlyingRequests(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))

	do.OverrideValue[loanPackageRequestRepo.LoanPackageRequestEventRepository](
		injector, &mock.LoanPackageRequestEventRepository{},
	)
	financialProductRepository := &mock.MockFinancialProductRepository{}
	do.OverrideValue[repository.FinancialProductRepository](
		injector, financialProductRepository,
	)
	configurationRepoMock := mock.NewMockConfigurationPersistenceRepository(t)
	marginOperationRepository := &mock.MockMarginOperationRepository{}
	odooServiceRepository := &mock.MockOdooServiceRepository{}
	do.OverrideValue[marginOperationRepo.MarginOperationRepository](
		injector, marginOperationRepository,
	)
	do.OverrideValue[financingApiRepository.FinancingRepository](injector, mock.FinancingApiMock{})
	do.OverrideValue[configRepo.ConfigurationPersistenceRepository](injector, configurationRepoMock)
	do.OverrideValue[odooServiceRepo.OdooServiceRepository](injector, odooServiceRepository)
	loanPackageRequestHandler := do.MustInvoke[*http.LoanPackageRequestHandler](injector)

	t.Run(
		"admin confirm loan package request with loan package id success", func(t *testing.T) {
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
			mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			request2 := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)

			request3 := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)

			mock.SeedSubmissionSheetMetadata(
				t, db, model.SubmissionSheetMetadata{
					ID:                   1,
					LoanPackageRequestID: request3.ID,
					Status:               "REJECTED",
					ActionType:           "REJECT_AND_SEND_OTHER_PROPOSAL",
					FlowType:             "002",
					ProposeType:          "NEW_LOAN_PACKAGE",
					CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
				},
			)

			submissionSheetMetadata2 := mock.SeedSubmissionSheetMetadata(
				t, db, model.SubmissionSheetMetadata{
					ID:                   2,
					LoanPackageRequestID: request3.ID,
					Status:               "DRAFT",
					ActionType:           "REJECT_AND_SEND_OTHER_PROPOSAL",
					FlowType:             "002",
					ProposeType:          "NEW_LOAN_PACKAGE",
					CreatedAt:            time.Date(2021, 2, 1, 0, 0, 0, 0, time.Local),
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET",
				"?statuses=PENDING",
				nil,
			)
			loanPackageRequestHandler.GetAllUnderlyingRequests(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 2, testhelper.GetArrayLength(body, "data"))
			assert.Equal(t, request2.ID, testhelper.GetInt(body, "data", "[1]", "id"))
			assert.Equal(t, request3.ID, testhelper.GetInt(body, "data", "[0]", "id"))
			assert.Equal(t, submissionSheetMetadata2.ID, testhelper.GetInt(body, "data", "[0]", "submissionId"))
			assert.Equal(t, "DRAFT", testhelper.GetString(body, "data", "[0]", "submissionStatus"))
		})
}
