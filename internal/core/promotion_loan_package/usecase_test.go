package promotionloanpackage

import (
	"context"
	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"financing-offer/test/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	testify "github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestPromotionLoanPackageUseCase_GetPromotionLoanPackages(t *testing.T) {
	t.Parallel()

	t.Run("get promotion loan packages success", func(t *testing.T) {
		promotionCampaignRepo := mock.NewMockPromotionCampaignRepository(t)
		appConfig := config.AppConfig{}
		configurationPersistenceRepo := mock.NewMockConfigurationPersistenceRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		financialProductRepo := mock.NewMockFinancialProductRepository(t)
		campaigns := []entity.PromotionCampaign{
			{
				Id:          1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				UpdatedBy:   "kiennt",
				Name:        "name",
				Tag:         "5.99*",
				Status:      entity.Active,
				Description: "description",
				Metadata: entity.PromotionCampaignMetadata{
					Products: []entity.PromotionCampaignProduct{
						{
							Symbols:       []string{"HPG", "DGW"},
							LoanPackageId: 1,
							RetailSymbols: []string{"HPG", "DGW"},
						},
					},
				},
			},
		}
		promotionCampaignRepo.EXPECT().GetAll(testify.Anything, testify.Anything).Return(campaigns, nil)

		financialProductRepo.EXPECT().GetAllAccountDetailByCustodyCode(testify.Anything, testify.Anything).Return([]entity.FinancialAccountDetail{
			{
				AccountNo: "accountNo1",
			},
			{
				AccountNo: "accountNo2",
			},
		}, nil)

		orderServiceRepo.EXPECT().GetAllAccountLoanPackages(testify.Anything, "accountNo1").Return([]entity.AccountLoanPackage{
			{
				Id:                       1,
				Name:                     "name",
				Type:                     "M",
				BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.1),
				BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.1),
				TransferFee:              decimal.NewFromFloat(0.1),
				Description:              "description",
				LoanProducts: []entity.LoanProduct{
					{
						Id:                       "1",
						Name:                     "name",
						Symbol:                   "HPG",
						InterestRate:             decimal.NewFromFloat(0.1),
						Term:                     1,
						PreferentialPeriod:       1,
						PreferentialInterestRate: decimal.NewFromFloat(0.1),
						AllowExtendLoanTerm:      true,
						AllowEarlyPayment:        true,
					},
				},
				BasketId: 1,
			},
		}, nil)

		orderServiceRepo.EXPECT().GetAllAccountLoanPackages(testify.Anything, "accountNo2").Return([]entity.AccountLoanPackage{
			{
				Id:                       2,
				Name:                     "name",
				Type:                     "M",
				BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.1),
				BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.1),
				TransferFee:              decimal.NewFromFloat(0.1),
				Description:              "description",
				LoanProducts: []entity.LoanProduct{
					{
						Id:                       "2",
						Name:                     "name",
						Symbol:                   "DGW",
						InterestRate:             decimal.NewFromFloat(0.1),
						Term:                     1,
						PreferentialPeriod:       1,
						PreferentialInterestRate: decimal.NewFromFloat(0.1),
						AllowExtendLoanTerm:      true,
						AllowEarlyPayment:        true,
					},
				},
			},
		}, nil)

		useCase := NewUseCase(appConfig, configurationPersistenceRepo, orderServiceRepo, financialProductRepo, promotionCampaignRepo)
		res, err := useCase.GetPromotionLoanPackages(context.Background(), "", "custodyCode", "")
		assert.Nil(t, err)
		assert.Equal(t, 2, len(res))
		assert.Equal(t, int64(1), res["accountNo1"][0].Id)
		assert.Equal(t, "name", res["accountNo1"][0].CampaignProducts[0].Campaign.Name)
		assert.Equal(t, "1", res["accountNo1"][0].CampaignProducts[0].Product.Id)
	})

	t.Run("get promotion loan packages error when get campaign error", func(t *testing.T) {
		promotionCampaignRepo := mock.NewMockPromotionCampaignRepository(t)
		appConfig := config.AppConfig{}
		configurationPersistenceRepo := mock.NewMockConfigurationPersistenceRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		financialProductRepo := mock.NewMockFinancialProductRepository(t)
		promotionCampaignRepo.EXPECT().GetAll(testify.Anything, testify.Anything).Return(nil, assert.AnError)
		useCase := NewUseCase(appConfig, configurationPersistenceRepo, orderServiceRepo, financialProductRepo, promotionCampaignRepo)
		_, err := useCase.GetPromotionLoanPackages(context.Background(), "", "custodyCode", "")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("get promotion loan packages error when get custody code error", func(t *testing.T) {
		promotionCampaignRepo := mock.NewMockPromotionCampaignRepository(t)
		appConfig := config.AppConfig{}
		configurationPersistenceRepo := mock.NewMockConfigurationPersistenceRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		financialProductRepo := mock.NewMockFinancialProductRepository(t)
		campaigns := []entity.PromotionCampaign{
			{
				Id:          1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				UpdatedBy:   "kiennt",
				Name:        "name",
				Tag:         "5.99*",
				Status:      entity.Active,
				Description: "description",
				Metadata: entity.PromotionCampaignMetadata{
					Products: []entity.PromotionCampaignProduct{
						{
							Symbols:       []string{"HPG", "DGW"},
							LoanPackageId: 1,
							RetailSymbols: []string{"HPG", "DGW"},
						},
					},
				},
			},
		}
		promotionCampaignRepo.EXPECT().GetAll(testify.Anything, testify.Anything).Return(campaigns, nil)

		financialProductRepo.EXPECT().GetAllAccountDetailByCustodyCode(testify.Anything, testify.Anything).Return(nil, assert.AnError)
		useCase := NewUseCase(appConfig, configurationPersistenceRepo, orderServiceRepo, financialProductRepo, promotionCampaignRepo)
		_, err := useCase.GetPromotionLoanPackages(context.Background(), "", "custodyCode", "")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("get promotion loan packages error when get user loan package error", func(t *testing.T) {
		promotionCampaignRepo := mock.NewMockPromotionCampaignRepository(t)
		appConfig := config.AppConfig{}
		configurationPersistenceRepo := mock.NewMockConfigurationPersistenceRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		financialProductRepo := mock.NewMockFinancialProductRepository(t)
		campaigns := []entity.PromotionCampaign{
			{
				Id:          1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				UpdatedBy:   "kiennt",
				Name:        "name",
				Tag:         "5.99*",
				Status:      entity.Active,
				Description: "description",
				Metadata: entity.PromotionCampaignMetadata{
					Products: []entity.PromotionCampaignProduct{
						{
							Symbols:       []string{"HPG", "DGW"},
							LoanPackageId: 1,
							RetailSymbols: []string{"HPG", "DGW"},
						},
					},
				},
			},
		}
		promotionCampaignRepo.EXPECT().GetAll(testify.Anything, testify.Anything).Return(campaigns, nil)

		financialProductRepo.EXPECT().GetAllAccountDetailByCustodyCode(testify.Anything, testify.Anything).Return([]entity.FinancialAccountDetail{
			{
				AccountNo: "accountNo1",
			},
			{
				AccountNo: "accountNo2",
			},
		}, nil)

		orderServiceRepo.EXPECT().GetAllAccountLoanPackages(testify.Anything, "accountNo1").Return(nil, assert.AnError)

		orderServiceRepo.EXPECT().GetAllAccountLoanPackages(testify.Anything, "accountNo2").Return([]entity.AccountLoanPackage{
			{
				Id:                       2,
				Name:                     "name",
				Type:                     "M",
				BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.1),
				BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.1),
				TransferFee:              decimal.NewFromFloat(0.1),
				Description:              "description",
				LoanProducts: []entity.LoanProduct{
					{
						Id:                       "2",
						Name:                     "name",
						Symbol:                   "DGW",
						InterestRate:             decimal.NewFromFloat(0.1),
						Term:                     1,
						PreferentialPeriod:       1,
						PreferentialInterestRate: decimal.NewFromFloat(0.1),
						AllowExtendLoanTerm:      true,
						AllowEarlyPayment:        true,
					},
				},
			},
		}, nil)

		useCase := NewUseCase(appConfig, configurationPersistenceRepo, orderServiceRepo, financialProductRepo, promotionCampaignRepo)
		_, err := useCase.GetPromotionLoanPackages(context.Background(), "", "custodyCode", "")
		assert.ErrorIs(t, err, assert.AnError)
	})

}

func TestPromotionLoanPackageUseCase_GetPublicPromotionLoanPackagesWithCampaigns(t *testing.T) {
	t.Parallel()

	t.Run("get public promotion loan packages success", func(t *testing.T) {
		promotionCampaignRepo := mock.NewMockPromotionCampaignRepository(t)
		appConfig := config.AppConfig{}
		configurationPersistenceRepo := mock.NewMockConfigurationPersistenceRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		financialProductRepo := mock.NewMockFinancialProductRepository(t)
		campaigns := []entity.PromotionCampaign{
			{
				Id:          1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				UpdatedBy:   "kiennt",
				Name:        "name",
				Tag:         "5.99*",
				Status:      entity.Active,
				Description: "description",
				Metadata: entity.PromotionCampaignMetadata{
					Products: []entity.PromotionCampaignProduct{
						{
							Symbols:       []string{"HPG", "DGW"},
							LoanPackageId: 1,
							RetailSymbols: []string{"HPG", "DGW"},
						},
					},
				},
			},
		}
		promotionCampaignRepo.EXPECT().GetAll(testify.Anything, testify.Anything).Return(campaigns, nil)

		financialProductRepo.EXPECT().GetLoanPackageDetails(testify.Anything, []int64{1}).Return([]entity.FinancialProductLoanPackage{
			{
				Id:                       1,
				Name:                     "name",
				InterestRate:             decimal.NewFromFloat(0.1),
				Term:                     1,
				BuyingFeeRate:            decimal.NewFromFloat(0.1),
				LoanBasketId:             1,
				Description:              "description",
				InitialRateForWithdraw:   decimal.NewFromFloat(0.1),
				MaintenanceRate:          decimal.NewFromFloat(0.1),
				LiquidRate:               decimal.NewFromFloat(0.1),
				PreferentialPeriod:       1,
				PreferentialInterestRate: decimal.NewFromFloat(0.1),
				AllowEarlyPayment:        true,
				AllowExtendLoanTerm:      true,
				LoanType:                 "M",
				BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.1),
				TransferFee:              decimal.NewFromFloat(0.1),
			},
		}, nil)

		financialProductRepo.EXPECT().GetMarginBasketsByIds(testify.Anything, []int64{1}).
			Return([]entity.MarginBasket{
				{
					Id:   1,
					Name: "name",
					Symbols: []string{
						"HPG",
						"DGW",
					},
					LoanProductIds: []int64{1},
					LoanProducts: []entity.MarginProduct{
						{
							Id:         1,
							Name:       "name",
							Symbol:     "HPG",
							LoanRateId: 1,
							LoanRate: entity.LoanRate{
								Id:                     1,
								Name:                   "name",
								InitialRate:            decimal.NewFromFloat(0.1),
								InitialRateForWithdraw: decimal.NewFromFloat(0.1),
								MaintenanceRate:        decimal.NewFromFloat(0.1),
								LiquidRate:             decimal.NewFromFloat(0.1),
							},
							LoanPolicies: []entity.LoanProductPolicy{
								{
									Rate:            decimal.NewFromFloat(0.1),
									RateForWithdraw: decimal.NewFromFloat(0.1),
									LoanPolicyId:    1,
									LoanPolicy: entity.FinancialProductLoanPolicy{
										Id:                       1,
										Name:                     "name",
										Source:                   "source",
										InterestRate:             decimal.NewFromFloat(0.1),
										InterestBasis:            1,
										Term:                     1,
										OverdueInterest:          decimal.NewFromFloat(0.1),
										AllowExtendLoanTerm:      true,
										AllowEarlyPayment:        true,
										PreferentialPeriod:       1,
										PreferentialInterestRate: decimal.NewFromFloat(0.1),
										CreatedDate:              time.Now(),
										ModifiedDate:             time.Now(),
									},
								},
							},
						},
						{
							Id:         2,
							Name:       "name",
							Symbol:     "DGW",
							LoanRateId: 2,
							LoanRate: entity.LoanRate{
								Id:                     2,
								Name:                   "name",
								InitialRate:            decimal.NewFromFloat(0.1),
								InitialRateForWithdraw: decimal.NewFromFloat(0.1),
								MaintenanceRate:        decimal.NewFromFloat(0.1),
								LiquidRate:             decimal.NewFromFloat(0.1),
							},
							LoanPolicies: []entity.LoanProductPolicy{
								{
									Rate:            decimal.NewFromFloat(0.1),
									RateForWithdraw: decimal.NewFromFloat(0.1),
									LoanPolicyId:    2,
									LoanPolicy: entity.FinancialProductLoanPolicy{
										Id:                       2,
										Name:                     "name",
										Source:                   "source",
										InterestRate:             decimal.NewFromFloat(0.1),
										InterestBasis:            1,
										Term:                     1,
										OverdueInterest:          decimal.NewFromFloat(0.1),
										AllowExtendLoanTerm:      true,
										AllowEarlyPayment:        true,
										PreferentialPeriod:       1,
										PreferentialInterestRate: decimal.NewFromFloat(0.1),
										CreatedDate:              time.Now(),
										ModifiedDate:             time.Now(),
									},
								},
							},
						},
					},
				},
			}, nil)

		financialProductRepo.EXPECT().GetLoanProducts(testify.Anything, entity.MarginProductFilter{Symbol: "HPG"}).Return([]entity.MarginProduct{
			{
				Id:         1,
				Name:       "name",
				Symbol:     "HPG",
				LoanRateId: 1,
				LoanRate: entity.LoanRate{
					Id:                     1,
					Name:                   "name",
					InitialRate:            decimal.NewFromFloat(0.1),
					InitialRateForWithdraw: decimal.NewFromFloat(0.1),
					MaintenanceRate:        decimal.NewFromFloat(0.1),
					LiquidRate:             decimal.NewFromFloat(0.1),
				},
				LoanPolicies: []entity.LoanProductPolicy{
					{
						Rate:            decimal.NewFromFloat(0.1),
						RateForWithdraw: decimal.NewFromFloat(0.1),
						LoanPolicyId:    1,
						LoanPolicy: entity.FinancialProductLoanPolicy{
							Id:                       1,
							Name:                     "name",
							Source:                   "source",
							InterestRate:             decimal.NewFromFloat(0.1),
							InterestBasis:            1,
							Term:                     1,
							OverdueInterest:          decimal.NewFromFloat(0.1),
							AllowExtendLoanTerm:      true,
							AllowEarlyPayment:        true,
							PreferentialPeriod:       1,
							PreferentialInterestRate: decimal.NewFromFloat(0.1),
							CreatedDate:              time.Now(),
							ModifiedDate:             time.Now(),
						},
					},
				},
			},
		}, nil)

		financialProductRepo.EXPECT().GetLoanProducts(testify.Anything, entity.MarginProductFilter{Symbol: "DGW"}).Return([]entity.MarginProduct{
			{
				Id:         2,
				Name:       "name",
				Symbol:     "DGW",
				LoanRateId: 2,
				LoanRate: entity.LoanRate{
					Id:                     2,
					Name:                   "name",
					InitialRate:            decimal.NewFromFloat(0.1),
					InitialRateForWithdraw: decimal.NewFromFloat(0.1),
					MaintenanceRate:        decimal.NewFromFloat(0.1),
					LiquidRate:             decimal.NewFromFloat(0.1),
				},
				LoanPolicies: []entity.LoanProductPolicy{
					{
						Rate:            decimal.NewFromFloat(0.1),
						RateForWithdraw: decimal.NewFromFloat(0.1),
						LoanPolicyId:    2,
						LoanPolicy: entity.FinancialProductLoanPolicy{
							Id:                       2,
							Name:                     "name",
							Source:                   "source",
							InterestRate:             decimal.NewFromFloat(0.1),
							InterestBasis:            1,
							Term:                     1,
							OverdueInterest:          decimal.NewFromFloat(0.1),
							AllowExtendLoanTerm:      true,
							AllowEarlyPayment:        true,
							PreferentialPeriod:       1,
							PreferentialInterestRate: decimal.NewFromFloat(0.1),
							CreatedDate:              time.Now(),
							ModifiedDate:             time.Now(),
						},
					},
				},
			},
		}, nil)

		useCase := NewUseCase(appConfig, configurationPersistenceRepo, orderServiceRepo, financialProductRepo, promotionCampaignRepo)
		res, err := useCase.GetPublicPromotionLoanPackagesWithCampaigns(context.Background(), "")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, 2, len(res[0].CampaignProducts))
	})

	t.Run("get public promotion loan packages error when get campaign error", func(t *testing.T) {
		promotionCampaignRepo := mock.NewMockPromotionCampaignRepository(t)
		appConfig := config.AppConfig{}
		configurationPersistenceRepo := mock.NewMockConfigurationPersistenceRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		financialProductRepo := mock.NewMockFinancialProductRepository(t)
		promotionCampaignRepo.EXPECT().GetAll(testify.Anything, testify.Anything).Return(nil, assert.AnError)
		useCase := NewUseCase(appConfig, configurationPersistenceRepo, orderServiceRepo, financialProductRepo, promotionCampaignRepo)
		_, err := useCase.GetPublicPromotionLoanPackagesWithCampaigns(context.Background(), "")
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("get public promotion loan packages error when get loan package details error", func(t *testing.T) {
		promotionCampaignRepo := mock.NewMockPromotionCampaignRepository(t)
		appConfig := config.AppConfig{}
		configurationPersistenceRepo := mock.NewMockConfigurationPersistenceRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		financialProductRepo := mock.NewMockFinancialProductRepository(t)
		campaigns := []entity.PromotionCampaign{
			{
				Id:          1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				UpdatedBy:   "kiennt",
				Name:        "name",
				Tag:         "5.99*",
				Status:      entity.Active,
				Description: "description",
				Metadata: entity.PromotionCampaignMetadata{
					Products: []entity.PromotionCampaignProduct{
						{
							Symbols:       []string{"HPG", "DGW"},
							LoanPackageId: 1,
							RetailSymbols: []string{"HPG", "DGW"},
						},
					},
				},
			},
		}
		promotionCampaignRepo.EXPECT().GetAll(testify.Anything, testify.Anything).Return(campaigns, nil)

		financialProductRepo.EXPECT().GetLoanPackageDetails(testify.Anything, []int64{1}).Return(nil, assert.AnError)
		useCase := NewUseCase(appConfig, configurationPersistenceRepo, orderServiceRepo, financialProductRepo, promotionCampaignRepo)
		_, err := useCase.GetPublicPromotionLoanPackagesWithCampaigns(context.Background(), "")
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("get public promotion loan packages error when get margin basket ids error", func(t *testing.T) {
		promotionCampaignRepo := mock.NewMockPromotionCampaignRepository(t)
		appConfig := config.AppConfig{}
		configurationPersistenceRepo := mock.NewMockConfigurationPersistenceRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		financialProductRepo := mock.NewMockFinancialProductRepository(t)
		campaigns := []entity.PromotionCampaign{
			{
				Id:          1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				UpdatedBy:   "kiennt",
				Name:        "name",
				Tag:         "5.99*",
				Status:      entity.Active,
				Description: "description",
				Metadata: entity.PromotionCampaignMetadata{
					Products: []entity.PromotionCampaignProduct{
						{
							Symbols:       []string{"HPG", "DGW"},
							LoanPackageId: 1,
							RetailSymbols: []string{"HPG", "DGW"},
						},
					},
				},
			},
		}
		promotionCampaignRepo.EXPECT().GetAll(testify.Anything, testify.Anything).Return(campaigns, nil)

		financialProductRepo.EXPECT().GetLoanPackageDetails(testify.Anything, []int64{1}).Return([]entity.FinancialProductLoanPackage{
			{
				Id:                       1,
				Name:                     "name",
				InterestRate:             decimal.NewFromFloat(0.1),
				Term:                     1,
				BuyingFeeRate:            decimal.NewFromFloat(0.1),
				LoanBasketId:             1,
				Description:              "description",
				InitialRateForWithdraw:   decimal.NewFromFloat(0.1),
				MaintenanceRate:          decimal.NewFromFloat(0.1),
				LiquidRate:               decimal.NewFromFloat(0.1),
				PreferentialPeriod:       1,
				PreferentialInterestRate: decimal.NewFromFloat(0.1),
				AllowEarlyPayment:        true,
				AllowExtendLoanTerm:      true,
				LoanType:                 "M",
				BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.1),
				TransferFee:              decimal.NewFromFloat(0.1),
			},
		}, nil)

		financialProductRepo.EXPECT().GetMarginBasketsByIds(testify.Anything, []int64{1}).
			Return(nil, assert.AnError)

		useCase := NewUseCase(appConfig, configurationPersistenceRepo, orderServiceRepo, financialProductRepo, promotionCampaignRepo)
		_, err := useCase.GetPublicPromotionLoanPackagesWithCampaigns(context.Background(), "")
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("get public promotion loan packages error when get loan product error", func(t *testing.T) {
		promotionCampaignRepo := mock.NewMockPromotionCampaignRepository(t)
		appConfig := config.AppConfig{}
		configurationPersistenceRepo := mock.NewMockConfigurationPersistenceRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		financialProductRepo := mock.NewMockFinancialProductRepository(t)
		campaigns := []entity.PromotionCampaign{
			{
				Id:          1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				UpdatedBy:   "kiennt",
				Name:        "name",
				Tag:         "5.99*",
				Status:      entity.Active,
				Description: "description",
				Metadata: entity.PromotionCampaignMetadata{
					Products: []entity.PromotionCampaignProduct{
						{
							Symbols:       []string{"HPG", "DGW"},
							LoanPackageId: 1,
							RetailSymbols: []string{"HPG", "DGW"},
						},
					},
				},
			},
		}
		promotionCampaignRepo.EXPECT().GetAll(testify.Anything, testify.Anything).Return(campaigns, nil)

		financialProductRepo.EXPECT().GetLoanPackageDetails(testify.Anything, []int64{1}).Return([]entity.FinancialProductLoanPackage{
			{
				Id:                       1,
				Name:                     "name",
				InterestRate:             decimal.NewFromFloat(0.1),
				Term:                     1,
				BuyingFeeRate:            decimal.NewFromFloat(0.1),
				LoanBasketId:             1,
				Description:              "description",
				InitialRateForWithdraw:   decimal.NewFromFloat(0.1),
				MaintenanceRate:          decimal.NewFromFloat(0.1),
				LiquidRate:               decimal.NewFromFloat(0.1),
				PreferentialPeriod:       1,
				PreferentialInterestRate: decimal.NewFromFloat(0.1),
				AllowEarlyPayment:        true,
				AllowExtendLoanTerm:      true,
				LoanType:                 "M",
				BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.1),
				TransferFee:              decimal.NewFromFloat(0.1),
			},
		}, nil)

		financialProductRepo.EXPECT().GetMarginBasketsByIds(testify.Anything, []int64{1}).
			Return([]entity.MarginBasket{
				{
					Id:   1,
					Name: "name",
					Symbols: []string{
						"HPG",
						"DGW",
					},
					LoanProductIds: []int64{1},
					LoanProducts: []entity.MarginProduct{
						{
							Id:         1,
							Name:       "name",
							Symbol:     "HPG",
							LoanRateId: 1,
							LoanRate: entity.LoanRate{
								Id:                     1,
								Name:                   "name",
								InitialRate:            decimal.NewFromFloat(0.1),
								InitialRateForWithdraw: decimal.NewFromFloat(0.1),
								MaintenanceRate:        decimal.NewFromFloat(0.1),
								LiquidRate:             decimal.NewFromFloat(0.1),
							},
							LoanPolicies: []entity.LoanProductPolicy{
								{
									Rate:            decimal.NewFromFloat(0.1),
									RateForWithdraw: decimal.NewFromFloat(0.1),
									LoanPolicyId:    1,
									LoanPolicy: entity.FinancialProductLoanPolicy{
										Id:                       1,
										Name:                     "name",
										Source:                   "source",
										InterestRate:             decimal.NewFromFloat(0.1),
										InterestBasis:            1,
										Term:                     1,
										OverdueInterest:          decimal.NewFromFloat(0.1),
										AllowExtendLoanTerm:      true,
										AllowEarlyPayment:        true,
										PreferentialPeriod:       1,
										PreferentialInterestRate: decimal.NewFromFloat(0.1),
										CreatedDate:              time.Now(),
										ModifiedDate:             time.Now(),
									},
								},
							},
						},
						{
							Id:         2,
							Name:       "name",
							Symbol:     "DGW",
							LoanRateId: 2,
							LoanRate: entity.LoanRate{
								Id:                     2,
								Name:                   "name",
								InitialRate:            decimal.NewFromFloat(0.1),
								InitialRateForWithdraw: decimal.NewFromFloat(0.1),
								MaintenanceRate:        decimal.NewFromFloat(0.1),
								LiquidRate:             decimal.NewFromFloat(0.1),
							},
							LoanPolicies: []entity.LoanProductPolicy{
								{
									Rate:            decimal.NewFromFloat(0.1),
									RateForWithdraw: decimal.NewFromFloat(0.1),
									LoanPolicyId:    2,
									LoanPolicy: entity.FinancialProductLoanPolicy{
										Id:                       2,
										Name:                     "name",
										Source:                   "source",
										InterestRate:             decimal.NewFromFloat(0.1),
										InterestBasis:            1,
										Term:                     1,
										OverdueInterest:          decimal.NewFromFloat(0.1),
										AllowExtendLoanTerm:      true,
										AllowEarlyPayment:        true,
										PreferentialPeriod:       1,
										PreferentialInterestRate: decimal.NewFromFloat(0.1),
										CreatedDate:              time.Now(),
										ModifiedDate:             time.Now(),
									},
								},
							},
						},
					},
				},
			}, nil)

		financialProductRepo.EXPECT().GetLoanProducts(testify.Anything, entity.MarginProductFilter{Symbol: "DGW"}).Return(nil, assert.AnError)
		financialProductRepo.EXPECT().GetLoanProducts(testify.Anything, entity.MarginProductFilter{Symbol: "HPG"}).Return(nil, assert.AnError)

		useCase := NewUseCase(appConfig, configurationPersistenceRepo, orderServiceRepo, financialProductRepo, promotionCampaignRepo)
		_, err := useCase.GetPublicPromotionLoanPackagesWithCampaigns(context.Background(), "")
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("get public promotion loan packages success when symbol not null", func(t *testing.T) {
		promotionCampaignRepo := mock.NewMockPromotionCampaignRepository(t)
		appConfig := config.AppConfig{}
		configurationPersistenceRepo := mock.NewMockConfigurationPersistenceRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		financialProductRepo := mock.NewMockFinancialProductRepository(t)
		campaigns := []entity.PromotionCampaign{
			{
				Id:          1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				UpdatedBy:   "kiennt",
				Name:        "name",
				Tag:         "5.99*",
				Status:      entity.Active,
				Description: "description",
				Metadata: entity.PromotionCampaignMetadata{
					Products: []entity.PromotionCampaignProduct{
						{
							Symbols:       []string{"HPG", "DGW"},
							LoanPackageId: 1,
							RetailSymbols: []string{"HPG", "DGW"},
						},
					},
				},
			},
		}
		promotionCampaignRepo.EXPECT().GetAll(testify.Anything, testify.Anything).Return(campaigns, nil)

		financialProductRepo.EXPECT().GetLoanPackageDetails(testify.Anything, []int64{1}).Return([]entity.FinancialProductLoanPackage{
			{
				Id:                       1,
				Name:                     "name",
				InterestRate:             decimal.NewFromFloat(0.1),
				Term:                     1,
				BuyingFeeRate:            decimal.NewFromFloat(0.1),
				LoanBasketId:             1,
				Description:              "description",
				InitialRateForWithdraw:   decimal.NewFromFloat(0.1),
				MaintenanceRate:          decimal.NewFromFloat(0.1),
				LiquidRate:               decimal.NewFromFloat(0.1),
				PreferentialPeriod:       1,
				PreferentialInterestRate: decimal.NewFromFloat(0.1),
				AllowEarlyPayment:        true,
				AllowExtendLoanTerm:      true,
				LoanType:                 "M",
				BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.1),
				TransferFee:              decimal.NewFromFloat(0.1),
			},
		}, nil)

		financialProductRepo.EXPECT().GetMarginBasketsByIds(testify.Anything, []int64{1}).
			Return([]entity.MarginBasket{
				{
					Id:   1,
					Name: "name",
					Symbols: []string{
						"HPG",
						"DGW",
					},
					LoanProductIds: []int64{1},
					LoanProducts: []entity.MarginProduct{
						{
							Id:         1,
							Name:       "name",
							Symbol:     "HPG",
							LoanRateId: 1,
							LoanRate: entity.LoanRate{
								Id:                     1,
								Name:                   "name",
								InitialRate:            decimal.NewFromFloat(0.1),
								InitialRateForWithdraw: decimal.NewFromFloat(0.1),
								MaintenanceRate:        decimal.NewFromFloat(0.1),
								LiquidRate:             decimal.NewFromFloat(0.1),
							},
							LoanPolicies: []entity.LoanProductPolicy{
								{
									Rate:            decimal.NewFromFloat(0.1),
									RateForWithdraw: decimal.NewFromFloat(0.1),
									LoanPolicyId:    1,
									LoanPolicy: entity.FinancialProductLoanPolicy{
										Id:                       1,
										Name:                     "name",
										Source:                   "source",
										InterestRate:             decimal.NewFromFloat(0.1),
										InterestBasis:            1,
										Term:                     1,
										OverdueInterest:          decimal.NewFromFloat(0.1),
										AllowExtendLoanTerm:      true,
										AllowEarlyPayment:        true,
										PreferentialPeriod:       1,
										PreferentialInterestRate: decimal.NewFromFloat(0.1),
										CreatedDate:              time.Now(),
										ModifiedDate:             time.Now(),
									},
								},
							},
						},
						{
							Id:         2,
							Name:       "name",
							Symbol:     "DGW",
							LoanRateId: 2,
							LoanRate: entity.LoanRate{
								Id:                     2,
								Name:                   "name",
								InitialRate:            decimal.NewFromFloat(0.1),
								InitialRateForWithdraw: decimal.NewFromFloat(0.1),
								MaintenanceRate:        decimal.NewFromFloat(0.1),
								LiquidRate:             decimal.NewFromFloat(0.1),
							},
							LoanPolicies: []entity.LoanProductPolicy{
								{
									Rate:            decimal.NewFromFloat(0.1),
									RateForWithdraw: decimal.NewFromFloat(0.1),
									LoanPolicyId:    2,
									LoanPolicy: entity.FinancialProductLoanPolicy{
										Id:                       2,
										Name:                     "name",
										Source:                   "source",
										InterestRate:             decimal.NewFromFloat(0.1),
										InterestBasis:            1,
										Term:                     1,
										OverdueInterest:          decimal.NewFromFloat(0.1),
										AllowExtendLoanTerm:      true,
										AllowEarlyPayment:        true,
										PreferentialPeriod:       1,
										PreferentialInterestRate: decimal.NewFromFloat(0.1),
										CreatedDate:              time.Now(),
										ModifiedDate:             time.Now(),
									},
								},
							},
						},
					},
				},
			}, nil)

		financialProductRepo.EXPECT().GetLoanProducts(testify.Anything, entity.MarginProductFilter{Symbol: "HPG"}).Return([]entity.MarginProduct{
			{
				Id:         1,
				Name:       "name",
				Symbol:     "HPG",
				LoanRateId: 1,
				LoanRate: entity.LoanRate{
					Id:                     1,
					Name:                   "name",
					InitialRate:            decimal.NewFromFloat(0.1),
					InitialRateForWithdraw: decimal.NewFromFloat(0.1),
					MaintenanceRate:        decimal.NewFromFloat(0.1),
					LiquidRate:             decimal.NewFromFloat(0.1),
				},
				LoanPolicies: []entity.LoanProductPolicy{
					{
						Rate:            decimal.NewFromFloat(0.1),
						RateForWithdraw: decimal.NewFromFloat(0.1),
						LoanPolicyId:    1,
						LoanPolicy: entity.FinancialProductLoanPolicy{
							Id:                       1,
							Name:                     "name",
							Source:                   "source",
							InterestRate:             decimal.NewFromFloat(0.1),
							InterestBasis:            1,
							Term:                     1,
							OverdueInterest:          decimal.NewFromFloat(0.1),
							AllowExtendLoanTerm:      true,
							AllowEarlyPayment:        true,
							PreferentialPeriod:       1,
							PreferentialInterestRate: decimal.NewFromFloat(0.1),
							CreatedDate:              time.Now(),
							ModifiedDate:             time.Now(),
						},
					},
				},
			},
		}, nil)

		useCase := NewUseCase(appConfig, configurationPersistenceRepo, orderServiceRepo, financialProductRepo, promotionCampaignRepo)
		res, err := useCase.GetPublicPromotionLoanPackagesWithCampaigns(context.Background(), "HPG")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, 1, len(res[0].CampaignProducts))
		assert.Equal(t, "HPG", res[0].CampaignProducts[0].Product.Symbol)
	})
}
