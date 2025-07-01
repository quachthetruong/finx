package test

import (
	"encoding/json"
	http2 "net/http"
	"testing"
	"time"

	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/financialproduct/repository"
	orderServiceRepository "financing-offer/internal/core/orderservice/repository"
	"financing-offer/internal/core/promotion_loan_package/transport/http"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/jwttoken"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestPromotionLoanPackageHandler(t *testing.T) {
	t.Parallel()
	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	do.OverrideValue(injector, &mock.EmptyCache{})
	financialProductRepoMock := mock.NewMockFinancialProductRepository(t)
	orderServiceMock := mock.NewMockOrderServiceRepository(t)
	do.OverrideValue[repository.FinancialProductRepository](injector, financialProductRepoMock)
	do.OverrideValue[orderServiceRepository.OrderServiceRepository](injector, orderServiceMock)
	promotionLoanPackageHandler := do.MustInvoke[*http.PromotionLoanPackageHandler](injector)
	cfg := do.MustInvoke[config.AppConfig](injector)

	t.Run(
		"test get cheapest promotion configuration ", func(t *testing.T) {
			ctx, _, recorder := gintest.GetTestContext()
			promotionLoanPackageHandler.GetBestLoanPackageIds(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(
				t,
				len(cfg.BestPromotions.LoanPackageIds),
				testhelper.GetArrayLength(body, "data"),
			)
		},
	)

	t.Run(
		"test set promotion configuration and success", func(t *testing.T) {
			defer truncateData()
			bytes, err := json.Marshal(
				entity.PromotionLoanPackage{
					LoanProducts: []entity.PromotionLoanProduct{
						{
							LoanPackageId:    1,
							RetailSymbols:    []string{"VNM", "NT2"},
							NonRetailSymbols: []string{"HPG", "FPT"},
							Symbols:          []string{"VNM", "NT2", "HPG", "FPT"},
						},
						{
							LoanPackageId:    2,
							RetailSymbols:    []string{"SSI"},
							NonRetailSymbols: []string{"VIC"},
							Symbols:          []string{"SSI", "VIC"},
						},
					},
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "promotionLoanPackage",
					Value:         string(bytes),
					LastUpdatedBy: "admin",
				},
			)
			ctx, _, recorder := gintest.GetTestContext()
			ctx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "admin",
				},
			)
			ctx.Request = gintest.MustMakeRequest(
				"POST", "/v1/configurations/promotion-loan-packages", entity.PromotionLoanPackage{
					LoanProducts: []entity.PromotionLoanProduct{
						{
							LoanPackageId:    1,
							RetailSymbols:    []string{"VNM", "NT2"},
							NonRetailSymbols: []string{"HPG", "FPT"},
						},
						{
							LoanPackageId:    2,
							RetailSymbols:    []string{"SSI"},
							NonRetailSymbols: []string{"VIC"},
						},
					},
				},
			)
			promotionLoanPackageHandler.SetPromotionLoanPackages(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(
				t, []string{"VNM", "NT2"},
				testhelper.GetArrayString(body, "data", "loanProducts", "[0]", "retailSymbols"),
			)
			assert.Equal(
				t, []string{"HPG", "FPT"},
				testhelper.GetArrayString(body, "data", "loanProducts", "[0]", "nonRetailSymbols"),
			)
			assert.Equal(
				t, []string{"VNM", "NT2", "HPG", "FPT"},
				testhelper.GetArrayString(body, "data", "loanProducts", "[0]", "symbols"),
			)
			assert.Equal(t, int64(1), testhelper.GetInt(body, "data", "loanProducts", "[0]", "loanPackageId"))

		},
	)

	t.Run(
		"test get promotion configuration", func(t *testing.T) {
			defer truncateData()
			bytes, err := json.Marshal(
				entity.PromotionLoanPackage{
					LoanProducts: []entity.PromotionLoanProduct{
						{
							LoanPackageId:    1,
							RetailSymbols:    []string{"VNM", "NT2"},
							NonRetailSymbols: []string{"HPG", "FPT"},
							Symbols:          []string{"VNM", "NT2", "HPG", "FPT"},
						},
						{
							LoanPackageId:    2,
							RetailSymbols:    []string{"SSI"},
							NonRetailSymbols: []string{"VIC"},
							Symbols:          []string{"SSI", "VIC"},
						},
					},
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "promotionLoanPackage",
					Value:         string(bytes),
					LastUpdatedBy: "admin",
				},
			)
			ctx, _, recorder := gintest.GetTestContext()
			ctx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "investor",
				},
			)
			promotionLoanPackageHandler.GetPromotionLoanPackages(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(
				t, []string{"SSI"}, testhelper.GetArrayString(body, "data", "loanProducts", "[1]", "retailSymbols"),
			)
			assert.Equal(
				t, []string{"VIC"}, testhelper.GetArrayString(body, "data", "loanProducts", "[1]", "nonRetailSymbols"),
			)
			assert.Equal(
				t, []string{"SSI", "VIC"}, testhelper.GetArrayString(body, "data", "loanProducts", "[1]", "symbols"),
			)
		},
	)

	t.Run(
		"get public promotion loan packages (v2)", func(t *testing.T) {
			defer truncateData()
			ctx, _, recorder := gintest.GetTestContext()
			bytes, err := json.Marshal(
				entity.PromotionLoanPackage{
					LoanProducts: []entity.PromotionLoanProduct{
						{
							LoanPackageId: 4916,
							RetailSymbols: []string{"DSN", "FPT", "HPG", "MBB", "PVT"},
						},
						{
							LoanPackageId: 7360,
							RetailSymbols: []string{"CTR", "SHS", "AAA", "VCB"},
						},
					},
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "promotionLoanPackage",
					Value:         string(bytes),
					LastUpdatedBy: "admin",
				},
			)
			financialProductRepoMock.EXPECT().GetLoanPackageDetails(testifyMock.Anything, []int64{4916}).Return(
				[]entity.FinancialProductLoanPackage{
					{
						Id:                       4916,
						Name:                     "RocketX KQ 50% - LS 5.99%",
						InitialRate:              decimal.NewFromFloat(0.5),
						InterestRate:             decimal.NewFromFloat(0.0599),
						Term:                     180,
						BuyingFeeRate:            decimal.NewFromFloat(0.00045),
						LoanBasketId:             4914,
						InitialRateForWithdraw:   decimal.NewFromFloat(0.51),
						MaintenanceRate:          decimal.NewFromFloat(0.4),
						LiquidRate:               decimal.NewFromFloat(0.3),
						PreferentialPeriod:       0,
						PreferentialInterestRate: decimal.NewFromFloat(0),
						AllowExtendLoanTerm:      false,
						AllowEarlyPayment:        false,
						LoanType:                 "M",
						BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.00045),
						TransferFee:              decimal.NewFromFloat(0.3),
						Description:              "",
					},
				}, nil,
			).Once()
			financialProductRepoMock.EXPECT().GetMarginBasketsByIds(testifyMock.Anything, []int64{4914}).Return(
				[]entity.MarginBasket{
					{
						Id:             4914,
						Name:           "Rổ 5.99% Test",
						LoanProductIds: nil,
						Symbols: []string{
							"DGC",
							"DSN",
							"FPT",
							"HPG",
							"MBB",
							"PVT",
							"SSI",
							"STB",
							"VIC",
							"VND",
							"VNM",
						},
					},
				}, nil,
			).Once()
			financialProductRepoMock.EXPECT().GetLoanProducts(
				testifyMock.Anything, entity.MarginProductFilter{Symbol: "FPT"},
			).Return(
				[]entity.MarginProduct{
					{
						Id:           12314,
						Name:         "product 1",
						Symbol:       "FPT",
						LoanRateId:   4232,
						LoanPolicies: nil,
					},
				}, nil,
			).Once()
			ctx.AddParam("symbol", "FPT")
			promotionLoanPackageHandler.GetPublicPromotionLoanPackageBySymbol(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "RocketX KQ 50% - LS 5.99%", testhelper.GetString(body, "name"))
			assert.Equal(t, int64(4916), testhelper.GetInt(body, "id"))
			assert.Equal(t, 0.0599, testhelper.GetFloat(body, "loanProducts", "[0]", "interestRate"))
		},
	)

	t.Run(
		"get public promotion loan packages (v3)", func(t *testing.T) {
			defer truncateData()
			ctx, _, recorder := gintest.GetTestContext()
			bytes, err := json.Marshal(
				entity.PromotionLoanPackage{
					LoanProducts: []entity.PromotionLoanProduct{
						{
							LoanPackageId: 4915,
							RetailSymbols: []string{"DSN", "HPG", "MBB", "PVT"},
						},
						{
							LoanPackageId: 7360,
							RetailSymbols: []string{"CTR", "SHS", "FPT", "AAA", "VCB"},
						},
					},
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "promotionLoanPackage",
					Value:         string(bytes),
					LastUpdatedBy: "admin",
				},
			)
			financialProductRepoMock.EXPECT().GetLoanPackageDetails(testifyMock.Anything, []int64{7360}).Return(
				[]entity.FinancialProductLoanPackage{
					{
						Id:                       7360,
						Name:                     "RocketX KQ 50% - LS 5.99%",
						InitialRate:              decimal.NewFromFloat(0.5),
						InterestRate:             decimal.NewFromFloat(0.0599),
						Term:                     180,
						BuyingFeeRate:            decimal.NewFromFloat(0.00045),
						LoanBasketId:             7359,
						InitialRateForWithdraw:   decimal.NewFromFloat(0.51),
						MaintenanceRate:          decimal.NewFromFloat(0.4),
						LiquidRate:               decimal.NewFromFloat(0.3),
						PreferentialPeriod:       0,
						PreferentialInterestRate: decimal.NewFromFloat(0),
						AllowExtendLoanTerm:      false,
						AllowEarlyPayment:        false,
						LoanType:                 "M",
						BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.00045),
						TransferFee:              decimal.NewFromFloat(0.3),
						Description:              "",
					},
				}, nil,
			).Once()
			financialProductRepoMock.EXPECT().GetMarginBasketsByIds(testifyMock.Anything, []int64{7359}).Return(
				[]entity.MarginBasket{
					{
						Id:   7359,
						Name: "Rổ 5.99% Test",
						LoanProducts: []entity.MarginProduct{
							{
								Id:         101871,
								Name:       "SHS 9.99%",
								Symbol:     "SHS",
								LoanRateId: 7332,
								LoanRate: entity.LoanRate{
									Id:                     7332,
									Name:                   "Margin KQ 20%",
									InitialRate:            decimal.NewFromFloat(0.2),
									InitialRateForWithdraw: decimal.NewFromFloat(0.21),
									MaintenanceRate:        decimal.NewFromFloat(0.16),
									LiquidRate:             decimal.NewFromFloat(0.14),
								},
								LoanPolicies: []entity.LoanProductPolicy{
									{
										Rate:            decimal.NewFromFloat(0.8),
										RateForWithdraw: decimal.NewFromFloat(0.79),
										LoanPolicyId:    5785,
										LoanPolicy: entity.FinancialProductLoanPolicy{
											Id:                       5785,
											Name:                     "[Test] CS Rocket X 5.99%",
											Source:                   "DNSE",
											InterestRate:             decimal.NewFromFloat(0.0299),
											InterestBasis:            365,
											Term:                     180,
											OverdueInterest:          decimal.NewFromFloat(0.08985),
											AllowExtendLoanTerm:      true,
											AllowEarlyPayment:        true,
											PreferentialPeriod:       0,
											PreferentialInterestRate: decimal.Zero,
										},
									},
								},
							},
						},
						LoanProductIds: []int64{
							101873,
							101872,
							101862,
							101817,
							101322,
							101869,
							101871,
							101870,
						},
						Symbols: nil,
					},
				}, nil,
			).Once()
			financialProductRepoMock.EXPECT().GetLoanProducts(
				testifyMock.Anything, entity.MarginProductFilter{Symbol: "SHS"},
			).Return(
				[]entity.MarginProduct{
					{
						Id:         101871,
						Name:       "SHS 9.99%",
						Symbol:     "SHS",
						LoanRateId: 7332,
						LoanPolicies: []entity.LoanProductPolicy{
							{
								Rate:            decimal.NewFromFloat(0.8),
								RateForWithdraw: decimal.NewFromFloat(0.79),
								LoanPolicyId:    5785,
							},
						},
					},
				}, nil,
			).Once()
			ctx.AddParam("symbol", "SHS")
			promotionLoanPackageHandler.GetPublicPromotionLoanPackageBySymbol(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "RocketX KQ 50% - LS 5.99%", testhelper.GetString(body, "name"))
			assert.Equal(t, int64(7360), testhelper.GetInt(body, "id"))
			assert.Equal(t, 0.0299, testhelper.GetFloat(body, "loanProducts", "[0]", "interestRate"))
		},
	)

	t.Run(
		"get promotion loan packages for authenticated user", func(t *testing.T) {
			defer truncateData()
			ctx, _, recorder := gintest.GetTestContext()
			custodyCode := "064C999911"
			accountNo := "00094321"
			bytes, err := json.Marshal(
				entity.PromotionLoanPackage{
					LoanProducts: []entity.PromotionLoanProduct{
						{
							LoanPackageId: 4915,
							RetailSymbols: []string{"DSN", "HPG", "MBB", "PVT"},
						},
						{
							LoanPackageId: 7360,
							RetailSymbols: []string{"CTR", "SHS", "FPT", "AAA", "VCB"},
						},
					},
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "promotionLoanPackage",
					Value:         string(bytes),
					LastUpdatedBy: "admin",
				},
			)
			financialProductRepoMock.EXPECT().GetAllAccountDetailByCustodyCode(
				testifyMock.Anything, custodyCode,
			).Return(
				[]entity.FinancialAccountDetail{
					{
						AccountNo: accountNo,
					},
				}, nil,
			).Once()
			orderServiceMock.EXPECT().GetAllAccountLoanPackages(testifyMock.Anything, accountNo).Return(
				[]entity.AccountLoanPackage{
					{
						Id:                       7360,
						Name:                     "Test 980 1020 (5)",
						Type:                     "M",
						BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.002),
						BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.002),
						TransferFee:              decimal.NewFromFloat(0.3),
						Description:              "",
						LoanProducts: []entity.LoanProduct{
							{
								Id:                       "4635-FPT",
								Name:                     "FinX Tran My Huong - VND KQ 20%",
								Symbol:                   "FPT",
								InitialRate:              decimal.NewFromFloat(0.2),
								InitialRateForWithdraw:   decimal.NewFromFloat(0.21),
								MaintenanceRate:          decimal.NewFromFloat(0.16),
								LiquidRate:               decimal.NewFromFloat(0.14),
								InterestRate:             decimal.NewFromFloat(0.16),
								PreferentialPeriod:       0,
								PreferentialInterestRate: decimal.Zero,
								Term:                     90,
								AllowExtendLoanTerm:      true,
								AllowEarlyPayment:        true,
							},
						},
						BasketId: 4838,
					},
				}, nil,
			).Once()
			ctx.AddParam("symbol", "FPT")
			ctx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub:         custodyCode,
					CustodyCode: custodyCode,
				},
			)
			ctx.Request = gintest.MustMakeRequest("POST", "/v1/promotion-loan-packages/FPT", nil)
			promotionLoanPackageHandler.GetPromotionLoanPackageBySymbol(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "FPT", testhelper.GetString(body, "symbol"))
			assert.Equal(t, int64(7360), testhelper.GetInt(body, "accountNos", "[0]", "id"))
			assert.Equal(t, "4635-FPT", testhelper.GetString(body, "accountNos", "[0]", "loanProducts", "[0]", "id"))
			assert.Equal(t, 0.16, testhelper.GetFloat(body, "accountNos", "[0]", "loanProducts", "[0]", "interestRate"))
		},
	)

	t.Run(
		"get promotion loan packages for authenticated user error", func(t *testing.T) {
			defer truncateData()
			ctx, _, recorder := gintest.GetTestContext()
			custodyCode := "064C999911"
			bytes, err := json.Marshal(
				entity.PromotionLoanPackage{
					LoanProducts: []entity.PromotionLoanProduct{
						{
							LoanPackageId: 4915,
							RetailSymbols: []string{"DSN", "HPG", "MBB", "PVT"},
						},
						{
							LoanPackageId: 7360,
							RetailSymbols: []string{"CTR", "SHS", "FPT", "AAA", "VCB"},
						},
					},
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "promotionLoanPackage",
					Value:         string(bytes),
					LastUpdatedBy: "admin",
				},
			)
			financialProductRepoMock.EXPECT().GetAllAccountDetailByCustodyCode(
				testifyMock.Anything, custodyCode,
			).Return(
				nil, assert.AnError,
			).Once()
			ctx.AddParam("symbol", "FPT")
			ctx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub:         custodyCode,
					CustodyCode: custodyCode,
				},
			)
			ctx.Request = gintest.MustMakeRequest("POST", "/v1/promotion-loan-packages/FPT", nil)
			promotionLoanPackageHandler.GetPromotionLoanPackageBySymbol(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			assert.Equal(t, http2.StatusInternalServerError, result.StatusCode)
		},
	)

	t.Run(
		"get promotion loan packages for authenticated user with no suitable loan package", func(t *testing.T) {
			defer truncateData()
			ctx, _, recorder := gintest.GetTestContext()
			custodyCode := "064C999911"
			bytes, err := json.Marshal(
				entity.PromotionLoanPackage{
					LoanProducts: []entity.PromotionLoanProduct{
						{
							LoanPackageId: 4915,
							RetailSymbols: []string{"DSN", "HPG", "MBB", "PVT"},
						},
						{
							LoanPackageId: 7360,
							RetailSymbols: []string{"CTR", "SHS", "FPT", "AAA", "VCB"},
						},
					},
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "promotionLoanPackage",
					Value:         string(bytes),
					LastUpdatedBy: "admin",
				},
			)
			ctx.AddParam("symbol", "BID")
			ctx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub:         custodyCode,
					CustodyCode: custodyCode,
				},
			)
			ctx.Request = gintest.MustMakeRequest("POST", "/v1/promotion-loan-packages/BID", nil)
			promotionLoanPackageHandler.GetPromotionLoanPackageBySymbol(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "BID", testhelper.GetString(body, "symbol"))
			assert.Equal(t, 0, testhelper.GetArrayLength(body, "accountNos"))
		},
	)

	t.Run(
		"get all symbol promotion loan packages (public)", func(t *testing.T) {
			defer truncateData()
			ctx, _, recorder := gintest.GetTestContext()
			bytes, err := json.Marshal(
				entity.PromotionLoanPackage{
					LoanProducts: []entity.PromotionLoanProduct{
						{
							LoanPackageId: 7368,
							RetailSymbols: []string{"SHS", "VIB"},
						},
					},
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "promotionLoanPackage",
					Value:         string(bytes),
					LastUpdatedBy: "admin",
				},
			)
			financialProductRepoMock.EXPECT().GetLoanPackageDetails(testifyMock.Anything, []int64{7368}).Return(
				[]entity.FinancialProductLoanPackage{
					{
						Id:                       7368,
						Name:                     "RocketX KQ 50% - LS 5.99%",
						InitialRate:              decimal.NewFromFloat(0.5),
						InterestRate:             decimal.NewFromFloat(0.0599),
						Term:                     180,
						BuyingFeeRate:            decimal.NewFromFloat(0.00045),
						LoanBasketId:             7357,
						InitialRateForWithdraw:   decimal.NewFromFloat(0.51),
						MaintenanceRate:          decimal.NewFromFloat(0.4),
						LiquidRate:               decimal.NewFromFloat(0.3),
						PreferentialPeriod:       0,
						PreferentialInterestRate: decimal.NewFromFloat(0),
						AllowExtendLoanTerm:      false,
						AllowEarlyPayment:        false,
						LoanType:                 "M",
						BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.00045),
						TransferFee:              decimal.NewFromFloat(0.3),
						Description:              "",
					},
				}, nil,
			)
			financialProductRepoMock.EXPECT().GetMarginBasketsByIds(testifyMock.Anything, []int64{7357}).Return(
				[]entity.MarginBasket{
					{
						Id:   7357,
						Name: "Rổ 5.99% Test",
						LoanProducts: []entity.MarginProduct{
							{
								Id:         101871,
								Name:       "SHS 9.99%",
								Symbol:     "SHS",
								LoanRateId: 7332,
								LoanRate: entity.LoanRate{
									Id:                     7332,
									Name:                   "Margin KQ 20%",
									InitialRate:            decimal.NewFromFloat(0.2),
									InitialRateForWithdraw: decimal.NewFromFloat(0.21),
									MaintenanceRate:        decimal.NewFromFloat(0.16),
									LiquidRate:             decimal.NewFromFloat(0.14),
								},
								LoanPolicies: []entity.LoanProductPolicy{
									{
										Rate:            decimal.NewFromFloat(0.8),
										RateForWithdraw: decimal.NewFromFloat(0.79),
										LoanPolicyId:    5785,
										LoanPolicy: entity.FinancialProductLoanPolicy{
											Id:                       5785,
											Name:                     "[Test] CS Rocket X 5.99%",
											Source:                   "DNSE",
											InterestRate:             decimal.NewFromFloat(0.0299),
											InterestBasis:            365,
											Term:                     180,
											OverdueInterest:          decimal.NewFromFloat(0.08985),
											AllowExtendLoanTerm:      true,
											AllowEarlyPayment:        true,
											PreferentialPeriod:       0,
											PreferentialInterestRate: decimal.Zero,
										},
									},
								},
							},
						},
						LoanProductIds: []int64{
							101873,
							101872,
							101862,
							101817,
							101322,
							101869,
							101871,
							101870,
						},
						Symbols: nil,
					},
				}, nil,
			)
			financialProductRepoMock.EXPECT().GetLoanProducts(
				testifyMock.Anything, entity.MarginProductFilter{Symbol: "SHS"},
			).Return(
				[]entity.MarginProduct{
					{
						Id:         101871,
						Name:       "SHS 9.99%",
						Symbol:     "SHS",
						LoanRateId: 7332,
						LoanPolicies: []entity.LoanProductPolicy{
							{
								Rate:            decimal.NewFromFloat(0.8),
								RateForWithdraw: decimal.NewFromFloat(0.79),
								LoanPolicyId:    5785,
							},
						},
					},
				}, nil,
			).Once()
			financialProductRepoMock.EXPECT().GetLoanProducts(
				testifyMock.Anything, entity.MarginProductFilter{Symbol: "VIB"},
			).Return([]entity.MarginProduct{}, nil).Once()
			ctx.Request = gintest.MustMakeRequest("POST", "public/v1/promotion-loan-packages", nil)
			promotionLoanPackageHandler.GetPublicPromotionLoanPackages(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "SHS", testhelper.GetString(body, "data", "[0]", "symbol"))
			assert.Equal(t, int64(7368), testhelper.GetInt(body, "data", "[0]", "id"))
			assert.Equal(t, 0.0299, testhelper.GetFloat(body, "data", "[0]", "loanProducts", "[0]", "interestRate"))
		},
	)

	t.Run(
		"get all promotion loan packages (authenticated user)", func(t *testing.T) {
			defer truncateData()
			ctx, _, recorder := gintest.GetTestContext()
			custodyCode := "064C999911"
			accountNo := "00094321"
			bytes, err := json.Marshal(
				entity.PromotionLoanPackage{
					LoanProducts: []entity.PromotionLoanProduct{
						{
							LoanPackageId:    7368,
							RetailSymbols:    []string{"TCH", "MBB"},
							NonRetailSymbols: []string{"VNM", "NT2"},
							Symbols:          []string{"TCH", "MBB", "VNM", "NT2"},
						},
					},
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "promotionLoanPackage",
					Value:         string(bytes),
					LastUpdatedBy: "admin",
				},
			)
			financialProductRepoMock.EXPECT().GetAllAccountDetailByCustodyCode(
				testifyMock.Anything, custodyCode,
			).Return(
				[]entity.FinancialAccountDetail{
					{
						AccountNo: accountNo,
					},
				}, nil,
			).Once()
			orderServiceMock.EXPECT().GetAllAccountLoanPackages(testifyMock.Anything, accountNo).Return(
				[]entity.AccountLoanPackage{
					{
						Id:                       7368,
						Name:                     "Test 980 1020 (5)",
						Type:                     "M",
						BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.002),
						BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.002),
						TransferFee:              decimal.NewFromFloat(0.3),
						Description:              "",
						LoanProducts: []entity.LoanProduct{
							{
								Id:                       "4635-TCH",
								Name:                     "FinX Tran My Huong - VND KQ 20%",
								Symbol:                   "TCH",
								InitialRate:              decimal.NewFromFloat(0.2),
								InitialRateForWithdraw:   decimal.NewFromFloat(0.21),
								MaintenanceRate:          decimal.NewFromFloat(0.16),
								LiquidRate:               decimal.NewFromFloat(0.14),
								InterestRate:             decimal.NewFromFloat(0.16),
								PreferentialPeriod:       0,
								PreferentialInterestRate: decimal.Zero,
								Term:                     90,
								AllowExtendLoanTerm:      true,
								AllowEarlyPayment:        true,
							},
						},
						BasketId: 4838,
					},
					{
						Id:                       7368,
						Name:                     "Test 980 1020 (5)",
						Type:                     "M",
						BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.002),
						BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.002),
						TransferFee:              decimal.NewFromFloat(0.3),
						Description:              "",
						LoanProducts: []entity.LoanProduct{
							{
								Id:                       "4636-MBB",
								Name:                     "FinX Tran My Huong - VND KQ 20%",
								Symbol:                   "MBB",
								InitialRate:              decimal.NewFromFloat(0.2),
								InitialRateForWithdraw:   decimal.NewFromFloat(0.21),
								MaintenanceRate:          decimal.NewFromFloat(0.16),
								LiquidRate:               decimal.NewFromFloat(0.14),
								InterestRate:             decimal.NewFromFloat(0.16),
								PreferentialPeriod:       0,
								PreferentialInterestRate: decimal.Zero,
								Term:                     90,
								AllowExtendLoanTerm:      true,
								AllowEarlyPayment:        true,
							},
						},
						BasketId: 4838,
					},
				}, nil,
			).Once()
			ctx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub:         custodyCode,
					CustodyCode: custodyCode,
				},
			)
			ctx.Request = gintest.MustMakeRequest("POST", "/v1/promotion-loan-packages", nil)
			promotionLoanPackageHandler.GetInvestorPromotionLoanPackage(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, accountNo, testhelper.GetString(body, "data", "[0]", "accountNo"))
			assert.Equal(t, 2, testhelper.GetArrayLength(body, "data", "[0]", "symbols"))
		},
	)

	t.Run(
		"get all promotion loan packages v2 (authenticated user)", func(t *testing.T) {
			defer truncateData()
			ctx, _, recorder := gintest.GetTestContext()
			custodyCode := "064C999911"
			accountNo := "00094321"
			bytes, err := json.Marshal(
				entity.PromotionCampaignMetadata{
					Products: []entity.PromotionCampaignProduct{
						{
							LoanPackageId: 7368,
							RetailSymbols: []string{"TCH", "MBB"},
							Symbols:       []string{"TCH", "MBB", "VNM", "NT2"},
						},
					},
				},
			)
			require.NoError(t, err)
			mock.SeedPromotionCampaign(
				t, db, model.PromotionCampaign{
					ID:          1,
					Name:        "Test campaign",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Status:      "ACTIVE",
					UpdatedBy:   "admin",
					Tag:         "5.99%",
					Description: "description",
					Metadata:    string(bytes),
				},
			)
			financialProductRepoMock.EXPECT().GetAllAccountDetailByCustodyCode(
				testifyMock.Anything, custodyCode,
			).Return(
				[]entity.FinancialAccountDetail{
					{
						AccountNo: accountNo,
					},
				}, nil,
			).Once()
			orderServiceMock.EXPECT().GetAllAccountLoanPackages(testifyMock.Anything, accountNo).Return(
				[]entity.AccountLoanPackage{
					{
						Id:                       7368,
						Name:                     "Test 980 1020 (5)",
						Type:                     "M",
						BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.002),
						BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.002),
						TransferFee:              decimal.NewFromFloat(0.3),
						Description:              "",
						LoanProducts: []entity.LoanProduct{
							{
								Id:                       "4635-TCH",
								Name:                     "FinX Tran My Huong - VND KQ 20%",
								Symbol:                   "TCH",
								InitialRate:              decimal.NewFromFloat(0.2),
								InitialRateForWithdraw:   decimal.NewFromFloat(0.21),
								MaintenanceRate:          decimal.NewFromFloat(0.16),
								LiquidRate:               decimal.NewFromFloat(0.14),
								InterestRate:             decimal.NewFromFloat(0.16),
								PreferentialPeriod:       0,
								PreferentialInterestRate: decimal.Zero,
								Term:                     90,
								AllowExtendLoanTerm:      true,
								AllowEarlyPayment:        true,
							},
						},
						BasketId: 4838,
					},
					{
						Id:                       7368,
						Name:                     "Test 980 1020 (5)",
						Type:                     "M",
						BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.002),
						BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.002),
						TransferFee:              decimal.NewFromFloat(0.3),
						Description:              "",
						LoanProducts: []entity.LoanProduct{
							{
								Id:                       "4636-MBB",
								Name:                     "FinX Tran My Huong - VND KQ 20%",
								Symbol:                   "MBB",
								InitialRate:              decimal.NewFromFloat(0.2),
								InitialRateForWithdraw:   decimal.NewFromFloat(0.21),
								MaintenanceRate:          decimal.NewFromFloat(0.16),
								LiquidRate:               decimal.NewFromFloat(0.14),
								InterestRate:             decimal.NewFromFloat(0.16),
								PreferentialPeriod:       0,
								PreferentialInterestRate: decimal.Zero,
								Term:                     90,
								AllowExtendLoanTerm:      true,
								AllowEarlyPayment:        true,
							},
						},
						BasketId: 4838,
					},
				}, nil,
			).Once()
			ctx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub:         custodyCode,
					CustodyCode: custodyCode,
				},
			)
			ctx.Request = gintest.MustMakeRequest("GET", "/v2/promotion-loan-packages", nil)
			promotionLoanPackageHandler.GetInvestorPromotionLoanPackageV2(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, accountNo, testhelper.GetString(body, "data", "[0]", "accountNo"))
			assert.Equal(t, 2, testhelper.GetArrayLength(body, "data", "[0]", "loanPackages"))
		},
	)

	t.Run(
		"get all symbol promotion loan packages v2 (public)", func(t *testing.T) {
			defer truncateData()
			ctx, _, recorder := gintest.GetTestContext()
			bytes, err := json.Marshal(
				entity.PromotionCampaignMetadata{
					Products: []entity.PromotionCampaignProduct{
						{
							LoanPackageId: 7368,
							RetailSymbols: []string{"SHS", "VIB"},
						},
					},
				},
			)
			require.NoError(t, err)
			mock.SeedPromotionCampaign(
				t, db, model.PromotionCampaign{
					ID:          1,
					Name:        "Test campaign",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Status:      "ACTIVE",
					UpdatedBy:   "admin",
					Tag:         "5.99%",
					Description: "description",
					Metadata:    string(bytes),
				},
			)
			financialProductRepoMock.EXPECT().GetLoanPackageDetails(testifyMock.Anything, []int64{7368}).Return(
				[]entity.FinancialProductLoanPackage{
					{
						Id:                       7368,
						Name:                     "RocketX KQ 50% - LS 5.99%",
						InitialRate:              decimal.NewFromFloat(0.5),
						InterestRate:             decimal.NewFromFloat(0.0599),
						Term:                     180,
						BuyingFeeRate:            decimal.NewFromFloat(0.00045),
						LoanBasketId:             7357,
						InitialRateForWithdraw:   decimal.NewFromFloat(0.51),
						MaintenanceRate:          decimal.NewFromFloat(0.4),
						LiquidRate:               decimal.NewFromFloat(0.3),
						PreferentialPeriod:       0,
						PreferentialInterestRate: decimal.NewFromFloat(0),
						AllowExtendLoanTerm:      false,
						AllowEarlyPayment:        false,
						LoanType:                 "M",
						BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.00045),
						TransferFee:              decimal.NewFromFloat(0.3),
						Description:              "",
					},
				}, nil,
			)
			financialProductRepoMock.EXPECT().GetMarginBasketsByIds(testifyMock.Anything, []int64{7357}).Return(
				[]entity.MarginBasket{
					{
						Id:   7357,
						Name: "Rổ 5.99% Test",
						LoanProducts: []entity.MarginProduct{
							{
								Id:         101871,
								Name:       "SHS 9.99%",
								Symbol:     "SHS",
								LoanRateId: 7332,
								LoanRate: entity.LoanRate{
									Id:                     7332,
									Name:                   "Margin KQ 20%",
									InitialRate:            decimal.NewFromFloat(0.2),
									InitialRateForWithdraw: decimal.NewFromFloat(0.21),
									MaintenanceRate:        decimal.NewFromFloat(0.16),
									LiquidRate:             decimal.NewFromFloat(0.14),
								},
								LoanPolicies: []entity.LoanProductPolicy{
									{
										Rate:            decimal.NewFromFloat(0.8),
										RateForWithdraw: decimal.NewFromFloat(0.79),
										LoanPolicyId:    5785,
										LoanPolicy: entity.FinancialProductLoanPolicy{
											Id:                       5785,
											Name:                     "[Test] CS Rocket X 5.99%",
											Source:                   "DNSE",
											InterestRate:             decimal.NewFromFloat(0.0299),
											InterestBasis:            365,
											Term:                     180,
											OverdueInterest:          decimal.NewFromFloat(0.08985),
											AllowExtendLoanTerm:      true,
											AllowEarlyPayment:        true,
											PreferentialPeriod:       0,
											PreferentialInterestRate: decimal.Zero,
										},
									},
								},
							},
						},
						LoanProductIds: []int64{
							101873,
							101872,
							101862,
							101817,
							101322,
							101869,
							101871,
							101870,
						},
						Symbols: nil,
					},
				}, nil,
			)
			financialProductRepoMock.EXPECT().GetLoanProducts(
				testifyMock.Anything, entity.MarginProductFilter{Symbol: "SHS"},
			).Return(
				[]entity.MarginProduct{
					{
						Id:         101871,
						Name:       "SHS 9.99%",
						Symbol:     "SHS",
						LoanRateId: 7332,
						LoanPolicies: []entity.LoanProductPolicy{
							{
								Rate:            decimal.NewFromFloat(0.8),
								RateForWithdraw: decimal.NewFromFloat(0.79),
								LoanPolicyId:    5785,
							},
						},
					},
				}, nil,
			).Once()
			financialProductRepoMock.EXPECT().GetLoanProducts(
				testifyMock.Anything, entity.MarginProductFilter{Symbol: "VIB"},
			).Return([]entity.MarginProduct{}, nil).Once()
			ctx.Request = gintest.MustMakeRequest("GET", "public/v2/promotion-loan-packages", nil)
			promotionLoanPackageHandler.GetPublicPromotionLoanPackagesV2(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, int64(7368), testhelper.GetInt(body, "data", "[0]", "id"))
			assert.Equal(t, 0.0299, testhelper.GetFloat(body, "data", "[0]", "campaignProducts", "[0]", "product", "interestRate"))
			assert.Equal(t, "SHS", testhelper.GetString(body, "data", "[0]", "campaignProducts", "[0]", "product", "symbol"))
			assert.Equal(t, "Test campaign", testhelper.GetString(body, "data", "[0]", "campaignProducts", "[0]", "campaign", "name"))
		},
	)
}
