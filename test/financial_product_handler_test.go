package test

import (
	configRepo "financing-offer/internal/config/repository"
	http2 "net/http"
	"testing"

	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/financialproduct/repository"
	"financing-offer/internal/core/financialproduct/transport/http"
	marginOperationRepo "financing-offer/internal/core/marginoperation/repository"
	"financing-offer/internal/jwttoken"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestFinancialProductHandler(t *testing.T) {
	t.Parallel()
	injector := testhelper.NewInjector()
	financialProductClientMock := mock.NewMockFinancialProductRepository(t)
	configurationRepoMock := mock.NewMockConfigurationPersistenceRepository(t)
	marginOperationClientMock := mock.NewMockMarginOperationRepository(t)
	cfg := config.AppConfig{}
	do.OverrideValue[repository.FinancialProductRepository](injector, financialProductClientMock)
	do.OverrideValue[marginOperationRepo.MarginOperationRepository](injector, marginOperationClientMock)
	do.OverrideValue[configRepo.ConfigurationPersistenceRepository](injector, configurationRepoMock)
	do.OverrideValue[config.AppConfig](injector, cfg)
	h := do.MustInvoke[*http.FinancialProductHandler](injector)
	t.Run(
		"GetLoanRates success", func(t *testing.T) {
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "admin",
				},
			)
			configurationRepoMock.EXPECT().GetLoanRateConfiguration(testifyMock.Anything).Return(entity.LoanRateConfiguration{
				Ids: []int64{45, 56, 32},
			}, nil)
			financialProductClientMock.EXPECT().GetLoanRatesByIds(
				testifyMock.Anything, []int64{45, 56, 32},
			).Return(
				[]entity.LoanRate{
					{
						Id:          45,
						Name:        "Loan rate 1",
						InitialRate: decimal.NewFromFloat(0.1),
					},
					{
						Id:          56,
						Name:        "Loan rate 2",
						InitialRate: decimal.NewFromFloat(0.2),
					},
					{
						Id:          32,
						Name:        "Loan rate 3",
						InitialRate: decimal.NewFromFloat(0.3),
					},
				}, nil,
			).Once()
			h.GetLoanRates(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 3, testhelper.GetArrayLength(body, "data"))
			assert.Equal(t, "Loan rate 1", testhelper.GetString(body, "data", "[0]", "name"))
		},
	)

	t.Run(
		"Get loanRates failed", func(t *testing.T) {
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "admin",
				},
			)
			configurationRepoMock.EXPECT().GetLoanRateConfiguration(testifyMock.Anything).Return(entity.LoanRateConfiguration{
				Ids: []int64{45, 56, 32},
			}, nil)
			financialProductClientMock.EXPECT().GetLoanRatesByIds(
				testifyMock.Anything, []int64{45, 56, 32},
			).Return(nil, assert.AnError).Once()
			h.GetLoanRates(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, http2.StatusInternalServerError, result.StatusCode)
			assert.Equal(t, "an error happened, please try again later", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"GetMarginPools success", func(t *testing.T) {
			defer financialProductClientMock.AssertExpectations(t)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "admin",
				},
			)
			configurationRepoMock.EXPECT().GetMarginPoolConfiguration(testifyMock.Anything).Return(entity.MarginPoolConfiguration{
				Ids: []int64{78, 79, 80},
			}, nil)
			marginOperationClientMock.EXPECT().GetMarginPoolsByIds(
				testifyMock.Anything, []int64{78, 79, 80},
			).Return(
				[]entity.MarginPool{
					{
						Id:          78,
						Name:        "Margin pool 1",
						PoolGroupId: 245,
					},
					{
						Id:          79,
						Name:        "Margin pool 2",
						PoolGroupId: 246,
					},
					{
						Id:          80,
						Name:        "Margin pool 3",
						PoolGroupId: 247,
					},
				}, nil,
			).Once()
			marginOperationClientMock.EXPECT().GetMarginPoolGroupsByIds(
				testifyMock.Anything, []int64{245, 246, 247},
			).Return(
				[]entity.MarginPoolGroup{
					{
						Id:   245,
						Name: "Group 1",
					},
					{
						Id:   246,
						Name: "Group 2",
					},
					{
						Id:   247,
						Name: "Group 3",
					},
				}, nil,
			).Once()
			h.GetMarginPools(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 3, testhelper.GetArrayLength(body, "data"))
			assert.Equal(t, "Margin pool 1", testhelper.GetString(body, "data", "[0]", "name"))
			assert.Equal(t, int64(247), testhelper.GetInt(body, "data", "[2]", "group", "id"))
		},
	)

	t.Run(
		"GetMarginPools get group failed", func(t *testing.T) {
			defer financialProductClientMock.AssertExpectations(t)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "admin",
				},
			)
			configurationRepoMock.EXPECT().GetMarginPoolConfiguration(testifyMock.Anything).Return(entity.MarginPoolConfiguration{
				Ids: []int64{78, 79, 80},
			}, nil)
			marginOperationClientMock.EXPECT().GetMarginPoolsByIds(
				testifyMock.Anything, []int64{78, 79, 80},
			).Return(
				[]entity.MarginPool{
					{
						Id:          78,
						Name:        "Margin pool 1",
						PoolGroupId: 245,
					},
					{
						Id:          79,
						Name:        "Margin pool 2",
						PoolGroupId: 246,
					},
					{
						Id:          80,
						Name:        "Margin pool 3",
						PoolGroupId: 247,
					},
				}, nil,
			).Once()
			marginOperationClientMock.EXPECT().GetMarginPoolGroupsByIds(
				testifyMock.Anything, []int64{245, 246, 247},
			).Return(
				nil, assert.AnError,
			).Once()
			h.GetMarginPools(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, http2.StatusInternalServerError, result.StatusCode)
		},
	)
}
