package test

import (
	"errors"
	"testing"

	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/financialproduct/repository"
	"financing-offer/internal/core/investor/transport/http"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func Test_InvestorHandler(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	t.Run(
		"test fill investor ids from requests success", func(t *testing.T) {
			defer truncateData()
			injector := testhelper.NewInjector(testhelper.WithDb(db))
			financialProductClientMock := mock.NewMockFinancialProductRepository(t)
			do.OverrideValue[repository.FinancialProductRepository](injector, financialProductClientMock)
			h := do.MustInvoke[*http.InvestorHandler](injector)
			investor1 := "00032421"
			custodyCode1 := "064C56453"
			investor2 := "00032422"
			custodyCode2 := "064C56454"
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
					InvestorID: investor1,
					LoanRate:   decimal.NewFromFloat(0.8),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: investor2,
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			financialProductClientMock.EXPECT().GetAllAccountDetail(mock2.Anything, investor1).Return(
				[]entity.FinancialAccountDetail{
					{
						Id:      "1",
						Custody: custodyCode1,
					},
				}, nil,
			)
			financialProductClientMock.EXPECT().GetAllAccountDetail(mock2.Anything, investor2).Return(
				[]entity.FinancialAccountDetail{
					{
						Id:      "2",
						Custody: custodyCode2,
					},
				}, nil,
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("POST", "/internal/investors/sync-investor-data", nil)
			h.FillInvestorIdsFromRequests(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 200, result.StatusCode)
			assert.Equal(t, int64(2), testhelper.GetInt(body, "data"))
		},
	)

	t.Run(
		"test fill investor ids from requests fail when get external user", func(t *testing.T) {
			defer truncateData()
			injector := testhelper.NewInjector(testhelper.WithDb(db))
			financialProductClientMock := mock.NewMockFinancialProductRepository(t)
			do.OverrideValue[repository.FinancialProductRepository](injector, financialProductClientMock)
			h := do.MustInvoke[*http.InvestorHandler](injector)
			investor1 := "00032421"
			investor2 := "00032422"
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
					InvestorID: investor1,
					LoanRate:   decimal.NewFromFloat(0.8),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: investor2,
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "CONFIRMED",
					AssetType:  "UNDERLYING",
				},
			)
			financialProductClientMock.EXPECT().GetAllAccountDetail(mock2.Anything, mock2.Anything).Return(
				nil, errors.New("error"),
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("POST", "/internal/investors/sync-investor-data", nil)
			h.FillInvestorIdsFromRequests(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, 500, result.StatusCode)
		},
	)
}
