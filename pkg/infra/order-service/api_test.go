package order_service

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
)

func TestClient_GetAllAccountLoanPackages(t *testing.T) {
	defer gock.Off()
	orderServiceConfig := config.OrderServiceConfig{
		Url:   "http://order-service",
		Token: "OrderServiceToken",
	}
	client := NewClient(orderServiceConfig)

	t.Run(
		"GetAllAccountLoanPackages success", func(t *testing.T) {
			gock.New(orderServiceConfig.Url).Get("/v2/accounts/123/loan-packages").HeaderPresent("Authorization").Reply(200).JSON(
				LoanPackagesResponse{
					LoanPackages: []entity.AccountLoanPackage{
						{
							Id:                       1,
							Name:                     "Loan Package 1",
							Type:                     "Type A",
							BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.02),
							BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.02),
							TransferFee:              decimal.NewFromFloat(2000),
							Description:              "This is a description for Loan Package 1",
							LoanProducts: []entity.LoanProduct{
								{
									Id:                       "1",
									Name:                     "Product 1",
									Symbol:                   "BID",
									InitialRate:              decimal.NewFromFloat(0.02),
									InitialRateForWithdraw:   decimal.NewFromFloat(0.02),
									MaintenanceRate:          decimal.NewFromFloat(0.02),
									LiquidRate:               decimal.NewFromFloat(0.02),
									InterestRate:             decimal.NewFromFloat(0.02),
									PreferentialPeriod:       0,
									PreferentialInterestRate: decimal.NewFromFloat(0.02),
									Term:                     0,
									AllowExtendLoanTerm:      false,
									AllowEarlyPayment:        false,
								},
							},
							BasketId: 1,
						},
					},
				},
			)
			res, err := client.GetAllAccountLoanPackages(context.Background(), "123")
			assert.Nil(t, err)
			assert.Equal(t, res[0].Name, "Loan Package 1")
		},
	)

	t.Run(
		"GetAllAccountLoanPackages verify investor_account not found", func(t *testing.T) {
			gock.New(orderServiceConfig.Url).Get("/v2/accounts/123/loan-packages").HeaderPresent("Authorization").Reply(404).JSON(
				ErrorResponse{
					Message: "Account not found 123",
				},
			)
			_, err := client.GetAllAccountLoanPackages(context.Background(), "123")
			assert.NotNil(t, err)
			assert.Equal(
				t, "GetAllAccountLoanPackages got error Status 404, Message: Account not found 123", err.Error(),
			)
		},
	)

	t.Run(
		"GetAllAccountLoanPackages client do error", func(t *testing.T) {
			gock.New(orderServiceConfig.Url).Get("/v2/accounts/123/loan-packages").HeaderPresent("Authorization").ReplyError(assert.AnError)
			_, err := client.GetAllAccountLoanPackages(context.Background(), "123")
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		},
	)

	t.Run(
		"GetAllAccountLoanPackages response decode error", func(t *testing.T) {
			gock.New(orderServiceConfig.Url).Get("/v2/accounts/123/loan-packages").HeaderPresent("Authorization").Reply(200).BodyString("invalid json")
			_, err := client.GetAllAccountLoanPackages(context.Background(), "123")
			assert.NotNil(t, err)
			assert.ErrorContains(t, err, "GetAllAccountLoanPackages invalid character")
		},
	)
}

func TestClient_GetAccountByAccountNoAndCustodyCode(t *testing.T) {
	defer gock.Off()
	orderServiceConfig := config.OrderServiceConfig{
		Url:   "http://order-service",
		Token: "OrderServiceToken",
	}
	client := NewClient(orderServiceConfig)

	t.Run(
		"test get account success", func(t *testing.T) {
			gock.New(orderServiceConfig.Url).Get("/internal/investors/064C000115/accounts/0001000115").HeaderPresent("Authorization").Reply(200).JSON(
				entity.OrderServiceAccount{
					CustodyCode:     "064C000115",
					AccountNo:       "0001000115",
					AccountTypeName: "test",
				},
			)
			res, err := client.GetAccountByAccountNoAndCustodyCode(context.Background(), "064C000115", "0001000115")
			assert.Nil(t, err)
			assert.Equal(t, res.AccountTypeName, "test")
			assert.Equal(t, res.CustodyCode, "064C000115")
			assert.Equal(t, res.AccountNo, "0001000115")
		},
	)

	t.Run(
		"test get account not found", func(t *testing.T) {
			gock.New(orderServiceConfig.Url).Get("/internal/investors/064C000115/accounts/0001000115").HeaderPresent("Authorization").Reply(404).JSON(
				ErrorResponse{
					Message: "Account not found 0001000115",
				},
			)
			_, err := client.GetAccountByAccountNoAndCustodyCode(context.Background(), "064C000115", "0001000115")
			assert.NotNil(t, err)
			assert.Equal(t, "invalid accountNo", err.Error())
		},
	)

	t.Run(
		"test get account 5xx error", func(t *testing.T) {
			gock.New(orderServiceConfig.Url).Get("/internal/investors/064C000115/accounts/0001000115").HeaderPresent("Authorization").Reply(500).JSON(
				ErrorResponse{
					Message: "Internal server error",
				},
			)
			_, err := client.GetAccountByAccountNoAndCustodyCode(context.Background(), "064C000115", "0001000115")
			assert.NotNil(t, err)
			assert.Equal(
				t, "GetAccountByAccountNoAndCustodyCode got error Status 500, Message: Internal server error",
				err.Error(),
			)
		},
	)

	t.Run(
		"test get account with malformed response format ", func(t *testing.T) {
			gock.New(orderServiceConfig.Url).Get("/internal/investors/064C000115/accounts/0001000115").HeaderPresent("Authorization").Reply(200).JSON(
				`{"accountNo": 1, "custodyCode":"064C000115"}`,
			)
			_, err := client.GetAccountByAccountNoAndCustodyCode(context.Background(), "064C000115", "0001000115")
			assert.NotNil(t, err)
			assert.Equal(
				t,
				"GetAccountByAccountNoAndCustodyCode json: cannot unmarshal number into Go struct field OrderServiceAccount.accountNo of type string",
				err.Error(),
			)
		},
	)
}
