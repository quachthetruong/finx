package test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/investor_account/repository"
	"financing-offer/internal/core/investor_account/transport/http"
	orderServiceRepo "financing-offer/internal/core/orderservice/repository"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestInvestorAccountHandler(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()

	t.Run(
		"VerifyAndUpdateInvestorAccountMarginStatus verify investor_account version 3 success", func(t *testing.T) {
			defer truncateData()
			injector := testhelper.NewInjector(testhelper.WithDb(db))
			h := do.MustInvoke[*http.InvestorAccountHandler](injector)
			account := mock.SeedInvestorAccount(
				t, db, model.InvestorAccount{
					AccountNo:    "123",
					InvestorID:   "456",
					MarginStatus: "v3",
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"PUT", "/v1/investor-accounts/123/margin-status", http.VerifyInvestorAccountMarginStatusRequest{
					InvestorId: account.InvestorID,
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "account-no",
				Value: account.AccountNo,
			}}
			h.VerifyAndUpdateInvestorAccountMarginStatus(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 200, result.StatusCode)
			assert.Equal(t, "v3", testhelper.GetString(body, "data", "marginStatus"))
		},
	)

	t.Run(
		"VerifyAndUpdateInvestorAccountMarginStatus bind json error", func(t *testing.T) {
			defer truncateData()
			injector := testhelper.NewInjector(testhelper.WithDb(db))
			h := do.MustInvoke[*http.InvestorAccountHandler](injector)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"PUT", "/v1/investor-accounts/123/margin-status", struct {
					AccountNo bool `json:"accountNo"`
				}{AccountNo: true},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "account-no",
				Value: "123",
			}}
			h.VerifyAndUpdateInvestorAccountMarginStatus(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 400, result.StatusCode)
			assert.Equal(t, "parse body", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"VerifyAndUpdateInvestorAccountMarginStatus verify then create investor_account version", func(t *testing.T) {
			defer truncateData()
			injector := testhelper.NewInjector(testhelper.WithDb(db))
			orderServiceRepositoryMock := mock.NewMockOrderServiceRepository(t)
			do.OverrideValue[orderServiceRepo.OrderServiceRepository](injector, orderServiceRepositoryMock)
			h := do.MustInvoke[*http.InvestorAccountHandler](injector)
			account := entity.InvestorAccount{
				AccountNo:    "123",
				InvestorId:   "456",
				MarginStatus: entity.MarginStatusVersion3,
			}
			loanPackage := entity.AccountLoanPackage{
				Id:                       1,
				Name:                     "Loan Package 1",
				Type:                     "M",
				BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.02),
				BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.02),
				TransferFee:              decimal.NewFromFloat(0.03),
				Description:              "This is a description for Loan Package 1",
				LoanProducts: []entity.LoanProduct{
					{
						Id:   "1-1204",
						Name: "Loan Product 1",
					},
				},
				BasketId: 1,
			}
			orderServiceRepositoryMock.EXPECT().GetAllAccountLoanPackages(mock2.Anything, account.AccountNo).Return(
				[]entity.AccountLoanPackage{loanPackage}, nil,
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"PUT", "/v1/investor-accounts/123/margin-status", http.VerifyInvestorAccountMarginStatusRequest{
					InvestorId: account.InvestorId,
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "account-no",
				Value: account.AccountNo,
			}}
			h.VerifyAndUpdateInvestorAccountMarginStatus(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 200, result.StatusCode)
			assert.Equal(t, "v3", testhelper.GetString(body, "data", "marginStatus"))
		},
	)

	t.Run(
		"VerifyAndUpdateInvestorAccountMarginStatus verify and update error", func(t *testing.T) {
			defer truncateData()
			injector := testhelper.NewInjector(testhelper.WithDb(db))
			accountRepositoryMock := mock.NewMockInvestorAccountRepository(t)
			do.OverrideValue[repository.InvestorAccountRepository](injector, accountRepositoryMock)
			h := do.MustInvoke[*http.InvestorAccountHandler](injector)
			accountRepositoryMock.EXPECT().GetByAccountNo(mock2.Anything, "123").Return(
				entity.InvestorAccount{}, assert.AnError,
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"PUT", "/v1/investor-accounts/123/margin-status", http.VerifyInvestorAccountMarginStatusRequest{
					InvestorId: "1",
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "account-no",
				Value: "123",
			}}
			h.VerifyAndUpdateInvestorAccountMarginStatus(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 500, result.StatusCode)
			assert.Equal(t, "an error happened, please try again later", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"VerifyAndUpdateInvestorAccountMarginStatus account no not found error", func(t *testing.T) {
			defer truncateData()
			injector := testhelper.NewInjector(testhelper.WithDb(db))
			h := do.MustInvoke[*http.InvestorAccountHandler](injector)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"PUT", "/v1/investor-accounts/123/margin-status", struct {
					AccountNo bool `json:"accountNo"`
				}{AccountNo: true},
			)
			h.VerifyAndUpdateInvestorAccountMarginStatus(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 400, result.StatusCode)
			assert.Equal(t, "invalid param: account-no", testhelper.GetString(body, "error"))
		},
	)
}
