package financialproduct

import (
	"context"
	"strconv"
	"testing"

	"github.com/h2non/gock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
)

func TestClient_AssignLoanPackage(t *testing.T) {
	defer gock.Off()
	financingProductConfig := config.FinancialProductConfig{
		Url:   "http://financing-product",
		Token: "financing-product-token",
	}
	client := NewClient(financingProductConfig)

	t.Run(
		"assign loan package success", func(t *testing.T) {
			gock.New(financingProductConfig.Url).Post("/v2/loan-package-accounts").HeaderPresent("Authorization").Reply(200).JSON(
				AssignLoanPackageResponse{
					Id:            145,
					AccountNo:     "12312",
					LoanPackageId: 878,
				},
			)
			loanAccountId, err := client.AssignLoanPackage(context.Background(), "12312", 878, "")
			assert.Nil(t, err)
			assert.Equal(t, int64(145), loanAccountId)
		},
	)

	t.Run(
		"assign loan package already assigned", func(t *testing.T) {
			gock.New(financingProductConfig.Url).Post("/v2/loan-package-accounts").HeaderPresent("Authorization").Reply(400).JSON(
				ErrorResponse{
					Status:  400,
					Code:    "INVALID_INPUT",
					Message: "already existed",
				},
			)
			_, err := client.AssignLoanPackage(context.Background(), "12312", 878, "")
			assert.ErrorIs(t, err, apperrors.ErrLoanPackageAccountAlreadyExisted("12312", 878))
		},
	)

	t.Run(
		"assign loan package or get existed loan package investor_account", func(t *testing.T) {
			accountNo := "an89"
			packageId := int64(1541)
			gock.New(financingProductConfig.Url).Post("/v2/loan-package-accounts").HeaderPresent("Authorization").Reply(400).JSON(
				ErrorResponse{
					Status:  400,
					Code:    "INVALID_INPUT",
					Message: "already existed",
				},
			)
			gock.New(financingProductConfig.Url).Get("/loan-package-accounts").HeaderPresent("Authorization").MatchParam(
				"accountNo", accountNo,
			).MatchParam("loanPackageId", strconv.FormatInt(packageId, 10)).Reply(200).BodyString(
				`
						{
							"data": [
								{
									"id": 9680,
									"accountNo": "0001175983",
									"loanPackageId": 1541,
									"createdDate": "2023-10-09T03:31:35.97992Z",
									"updatedDate": "2023-10-09T03:31:35.97992Z"
								}
							],
							"total": 1,
							"start": 0,
							"end": 1
						}`,
			)
			loanAccountId, err := client.AssignLoanPackageOrGetLoanPackageAccountId(
				context.Background(), accountNo, packageId, entity.AssetTypeUnderlying,
			)
			assert.Nil(t, err)
			assert.Equal(t, int64(9680), loanAccountId)
		},
	)

	t.Run(
		"assign derivative loan package or get existed loan package investor_account", func(t *testing.T) {
			accountNo := "an89"
			packageId := int64(1541)
			gock.New(financingProductConfig.Url).Post("/derivatives/package-accounts").HeaderPresent("Authorization").Reply(400).JSON(
				ErrorResponse{
					Status:  409,
					Code:    "RESOURCE_ALREADY_EXISTS",
					Message: "Pair loanPackageId and accountNo already exists",
				},
			)
			gock.New(financingProductConfig.Url).Get("/derivatives/package-accounts").HeaderPresent("Authorization").MatchParam(
				"accountNo", accountNo,
			).MatchParam("loanPackageId", strconv.FormatInt(packageId, 10)).Reply(200).BodyString(
				`
						{
							"data": [
								{
									"id": 876,
									"accountNo": "0001175983",
									"loanPackageId": 1541
								}
							],
							"total": 1,
							"start": 0,
							"end": 1
						}`,
			)
			loanAccountId, err := client.AssignLoanPackageOrGetLoanPackageAccountId(
				context.Background(), accountNo, packageId, entity.AssetTypeDerivative,
			)
			assert.Nil(t, err)
			assert.Equal(t, int64(876), loanAccountId)
		},
	)

	t.Run(
		"assign loan package not exist", func(t *testing.T) {
			gock.New(financingProductConfig.Url).Post("/v2/loan-package-accounts").HeaderPresent("Authorization").Reply(400).JSON(
				ErrorResponse{
					Status:  400,
					Code:    "RESOURCE_NOT_FOUND",
					Message: "loanPackage is not existed",
				},
			)
			_, err := client.AssignLoanPackage(context.Background(), "12312", 878, "")
			assert.ErrorIs(t, err, apperrors.ErrLoanPackageAccountNotExisted(878))
		},
	)
}

func TestClient_GetLoanPackageDetail(t *testing.T) {
	defer gock.Off()
	financingProductConfig := config.FinancialProductConfig{
		Url:   "http://financing-product",
		Token: "financing-product-token",
	}
	client := NewClient(financingProductConfig)

	t.Run(
		"test get loan package underlying detail", func(t *testing.T) {
			gock.New(financingProductConfig.Url).Get("/v2/loan-packages/4634").HeaderPresent("Authorization").Reply(200).BodyString(
				`
						{
							"id": 4634,
							"name": "FinX Tran My Huong - VND KQ 30%",
							"loanType": "M",
							"source": "PRIVATE",
							"initialRate": 0.3,
							"initialRateForWithdraw": 0.31,
							"maintenanceRate": 0.24,
							"liquidRate": 0.16,
							"preferentialPeriod": 0,
							"preferentialInterestRate": 0.0,
							"interestRate": 0.155,
							"term": 90,
							"overdueInterest": 0.2325,
							"allowExtendLoanTerm": true,
							"allowEarlyPayment": true,
							"brokerFirmBuyingFeeRate": 0.0012,
							"brokerFirmSellingFeeRate": 0.0012,
							"transferFee": 300,
							"loanBasketId": 4594,
							"interestBasis": 365,
							"modifiedDate": "2023-09-07T04:28:21.133653Z"
						}`,
			)
			loanPackageDetail, err := client.GetLoanPackageDetail(context.Background(), 4634)
			assert.Nil(t, err)
			assert.Equal(t, int64(4634), loanPackageDetail.Id)
			assert.Equal(t, decimal.NewFromFloat(0.3), loanPackageDetail.InitialRate)
		},
	)

	t.Run(
		"test get loan package detail with error code", func(t *testing.T) {
			gock.New(financingProductConfig.Url).Post("/v2/loan-packages/4634").HeaderPresent("Authorization").Reply(400)
			_, err := client.GetLoanPackageDetail(context.Background(), 4634)
			assert.Error(t, err)
		},
	)
}

func TestClient_GetLoanPackageDerivative(t *testing.T) {
	defer gock.Off()
	financingProductConfig := config.FinancialProductConfig{
		Url:   "http://financing-product",
		Token: "financing-product-token",
	}
	client := NewClient(financingProductConfig)
	t.Run(
		"test get loan package derivative success", func(t *testing.T) {
			gock.New(financingProductConfig.Url).Get("/derivatives/packages/3688").HeaderPresent("Authorization").Reply(200).BodyString(
				`
						{
							"id": 3688,
							"name": "Goi giao dich 1",
							"initialRate": 0.1848
						}`,
			)
			loanPackageDetail, err := client.GetLoanPackageDerivative(context.Background(), 3688)
			assert.Nil(t, err)
			assert.Equal(t, int64(3688), loanPackageDetail.Id)
			assert.Equal(t, decimal.NewFromFloat(0.1848), loanPackageDetail.InitialRate)
		},
	)

	t.Run(
		"test get loan package derivative with error code", func(t *testing.T) {
			gock.New(financingProductConfig.Url).Post("/derivatives/packages/3688").HeaderPresent("Authorization").Reply(400)
			_, err := client.GetLoanPackageDetail(context.Background(), 3688)
			assert.Error(t, err)
		},
	)
}

func TestClient_GetAllAccountDetail(t *testing.T) {
	defer gock.Off()
	financingProductConfig := config.FinancialProductConfig{
		Url:   "http://financing-product",
		Token: "financing-product-token",
	}
	client := NewClient(financingProductConfig)
	t.Run(
		"test get all investor_account detail", func(t *testing.T) {
			gock.New(financingProductConfig.Url).Get("/accounts").HeaderPresent("Authorization").MatchParam(
				"investorId", "0001000115",
			).Reply(200).BodyString(
				`
						{
							"accounts": [
								{
									"investorId": "0001000115",
									"accountNo": "0001007031",
									"custody": "064C000115",
									"fullName": "Nguyễn Văn Mâm",
									"foreigner": false,
									"status": "ACTIVE",
									"accountType": "0402",
									"accountTypeName": "RocketX 0.075%",
									"accountTypeBriefName": "Loại hình: Giao dịch ký quỹ; Lãi suất: 12%; Phí giao dịch: 0.075%",
									"marginAccount": true,
									"dealAccount": true,
									"carebyId": "0085",
									"customerId": "0001000115",
									"id": "0001007031"
								},
								{
									"investorId": "0001000115",
									"accountNo": "0001000115",
									"custody": "064C000115",
									"fullName": "Nguyễn Văn Mâm",
									"foreigner": false,
									"status": "ACTIVE",
									"accountType": "0394",
									"accountTypeName": "SpaceX",
									"accountTypeBriefName": "Loại hình: Giao dịch tiền mặt; Phí giao dịch: 0.03%",
									"marginAccount": false,
									"dealAccount": true,
									"carebyId": "0076",
									"customerId": "0001000115",
									"id": "0001000115"
								},
								{
									"investorId": "0001000115",
									"accountNo": "0001176887",
									"custody": "064C000115",
									"fullName": "Nguyễn Văn Mâm",
									"foreigner": false,
									"status": "ACTIVE",
									"accountType": "0395",
									"accountTypeName": "RocketX 0.03%",
									"accountTypeBriefName": "Loại hình: Giao dịch ký quỹ; Lãi suất:12%; Phí giao dịch: 0.03%",
									"marginAccount": true,
									"dealAccount": false,
									"carebyId": "0085",
									"customerId": "0001000115",
									"id": "0001176887"
								},
								{
									"investorId": "0001000115",
									"accountNo": "0011008845",
									"custody": "064C000115",
									"fullName": "Nguyễn Văn Mâm",
									"foreigner": false,
									"status": "ACTIVE",
									"accountType": "0394",
									"accountTypeName": "SpaceX",
									"accountTypeBriefName": "Loại hình: Giao dịch tiền mặt; Phí giao dịch: 0.03%",
									"marginAccount": false,
									"dealAccount": true,
									"carebyId": "0101",
									"customerId": "0001000115",
									"id": "0011008845"
								}
							]
						}`,
			)
			accounts, err := client.GetAllAccountDetail(context.Background(), "0001000115")
			assert.Nil(t, err)
			assert.Equal(t, 4, len(accounts))
		},
	)

	t.Run(
		"test get all investors account detail with by custodyCode", func(t *testing.T) {
			custodyCode := "064C000115"
			gock.New(financingProductConfig.Url).Get("/accounts").HeaderPresent("Authorization").MatchParam(
				"custodyCode", custodyCode,
			).Reply(200).BodyString(
				`
						{
							"accounts": [
								{
									"investorId": "0001000115",
									"accountNo": "0001007031",
									"custody": "064C000115",
									"fullName": "Nguyễn Văn Mâm",
									"foreigner": false,
									"status": "ACTIVE",
									"accountType": "0402",
									"accountTypeName": "RocketX 0.075%",
									"accountTypeBriefName": "Loại hình: Giao dịch ký quỹ; Lãi suất: 12%; Phí giao dịch: 0.075%",
									"marginAccount": true,
									"dealAccount": true,
									"carebyId": "0085",
									"customerId": "0001000115",
									"id": "0001007031"
								},
								{
									"investorId": "0001000115",
									"accountNo": "0001000115",
									"custody": "064C000115",
									"fullName": "Nguyễn Văn Mâm",
									"foreigner": false,
									"status": "ACTIVE",
									"accountType": "0394",
									"accountTypeName": "SpaceX",
									"accountTypeBriefName": "Loại hình: Giao dịch tiền mặt; Phí giao dịch: 0.03%",
									"marginAccount": false,
									"dealAccount": true,
									"carebyId": "0076",
									"customerId": "0001000115",
									"id": "0001000115"
								},
								{
									"investorId": "0001000115",
									"accountNo": "0001176887",
									"custody": "064C000115",
									"fullName": "Nguyễn Văn Mâm",
									"foreigner": false,
									"status": "ACTIVE",
									"accountType": "0395",
									"accountTypeName": "RocketX 0.03%",
									"accountTypeBriefName": "Loại hình: Giao dịch ký quỹ; Lãi suất:12%; Phí giao dịch: 0.03%",
									"marginAccount": true,
									"dealAccount": false,
									"carebyId": "0085",
									"customerId": "0001000115",
									"id": "0001176887"
								}
							]
						}`,
			)
			accounts, err := client.GetAllAccountDetailByCustodyCode(context.Background(), custodyCode)
			assert.Nil(t, err)
			assert.Equal(t, 3, len(accounts))
		},
	)

	t.Run(
		"test get all investor_account detail with error code", func(t *testing.T) {
			gock.New(financingProductConfig.Url).Get("/v2/loan-packages/4634").HeaderPresent("Authorization").Reply(400)
			_, err := client.GetAllAccountDetail(context.Background(), "0001000115")
			assert.Error(t, err)
		},
	)
}

func TestClient_GetLoanRatesByIds(t *testing.T) {
	defer gock.Off()
	financingProductConfig := config.FinancialProductConfig{
		Url:   "http://financing-product",
		Token: "financing-product-token",
	}
	client := NewClient(financingProductConfig)
	t.Run(
		"get loan rates by ids success", func(t *testing.T) {
			gock.New(financingProductConfig.Url).Get("/loan-rates").HeaderPresent("Authorization").MatchParam(
				"ids", "6443,6441",
			).Reply(200).BodyString(
				`
						{
							"data": [
								{
									"id": 6443,
									"name": "ky quy 40%",
									"initialRate": 0.4,
									"initialRateForWithdraw": 0.41,
									"maintenanceRate": 0.35,
									"liquidRate": 0.25,
									"createdDate": "2024-01-26T04:40:42.69872Z",
									"modifiedDate": "2024-01-26T04:40:42.698722Z"
								},
								{
									"id": 6441,
									"name": "ky quy mana 60%",
									"initialRate": 1.0,
									"initialRateForWithdraw": 0.6,
									"maintenanceRate": 0.5,
									"liquidRate": 0.4,
									"createdDate": "2024-01-26T04:40:40.854676Z",
									"modifiedDate": "2024-01-26T04:40:40.854678Z"
								}
							],
							"total": 2,
							"start": 0,
							"end": 2
						}`,
			)
			rates, err := client.GetLoanRatesByIds(context.Background(), []int64{6443, 6441})
			assert.Nil(t, err)
			assert.Equal(t, 2, len(rates))
			assert.Equal(t, int64(6443), rates[0].Id)
			assert.True(t, decimal.NewFromFloat(1).Equal(rates[1].InitialRate))
		},
	)
	t.Run(
		"get loan rates by ids failed", func(t *testing.T) {
			gock.New(financingProductConfig.Url).Get("/loan-rates").HeaderPresent("Authorization").MatchParam(
				"ids", "6443,6441",
			).Reply(400).BodyString("{}")
			_, err := client.GetLoanRatesByIds(context.Background(), []int64{6443, 6441})
			assert.Error(t, err)
			dest := apperrors.AppError{}
			assert.ErrorAs(t, err, &dest)
			assert.Equal(t, 400, dest.Code)
		},
	)
}

func TestClient_GetLoanPackageAccountIdByAccountNoAndLoanPackageId(t *testing.T) {
	defer gock.Off()
	financingProductConfig := config.FinancialProductConfig{
		Url:   "http://financing-product",
		Token: "financing-product-token",
	}
	client := NewClient(financingProductConfig)
	t.Run(
		"get loan package investor_account id by investor_account no and loan package id success", func(t *testing.T) {
			gock.New(financingProductConfig.Url).Get("/loan-package-accounts").HeaderPresent("Authorization").MatchParam(
				"accountNo", "0001175983",
			).MatchParam("loanPackageId", "1541").Reply(200).BodyString(
				`
						{
							"data": [
								{
									"id": 9680,
									"accountNo": "0001175983",
									"loanPackageId": 1541,
									"createdDate": "2023-10-09T03:31:35.97992Z",
									"updatedDate": "2023-10-09T03:31:35.97992Z"
								}
							],
							"total": 1,
							"start": 0,
							"end": 1
						}`,
			)
			loanPackageAccountId, err := client.GetLoanPackageAccountIdByAccountNoAndLoanPackageId(
				context.Background(), "0001175983", 1541, "",
			)
			assert.Nil(t, err)
			assert.Equal(t, int64(9680), loanPackageAccountId)
		},
	)
}

func TestClient_GetLoanProducts(t *testing.T) {
	defer gock.Off()
	financingProductConfig := config.FinancialProductConfig{
		Url:   "http://financing-product",
		Token: "financing-product-token",
	}
	client := NewClient(financingProductConfig)
	t.Run(
		"get loan products success", func(t *testing.T) {
			gock.New(financingProductConfig.Url).Get("/loan-products").HeaderPresent("Authorization").MatchParam(
				"symbol", "BID",
			).Reply(200).BodyString(
				`
					{
						"data": [
							{
								"id": 1,
								"name": "BID-1",
								"symbol": "BID"	
							},
							{
								"id": 2,
								"name": "BID-2",
								"symbol": "BID"	
							}
						],
						"total": 2,
						"start": 0,
						"end": 2
					}`,
			)
			loanProducts, err := client.GetLoanProducts(
				context.Background(), entity.MarginProductFilter{
					Symbol: "BID",
				},
			)
			assert.Nil(t, err)
			assert.Equal(t, 2, len(loanProducts))
			assert.Equal(t, int64(1), loanProducts[0].Id)
		},
	)
}

func TestClient_GetMarginBasketsByIds(t *testing.T) {
	defer gock.Off()
	financingProductConfig := config.FinancialProductConfig{
		Url:   "http://financing-product",
		Token: "financing-product-token",
	}
	client := NewClient(financingProductConfig)
	t.Run(
		"get margin baskets by ids success", func(t *testing.T) {
			gock.New(financingProductConfig.Url).Get("v2/margin-baskets/1/detail").HeaderPresent("Authorization").Reply(200).BodyString(
				`
							{
								"id": 1,
								"symbol": "BID"
							}`,
			)
			marginBaskets, err := client.GetMarginBasketsByIds(context.Background(), []int64{1})
			assert.Nil(t, err)
			assert.Equal(t, 1, len(marginBaskets))
			assert.Equal(t, int64(1), marginBaskets[0].Id)
		},
	)
}
