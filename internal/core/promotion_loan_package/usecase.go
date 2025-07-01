package promotionloanpackage

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"sync"

	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/config"
	"financing-offer/internal/config/repository"
	"financing-offer/internal/core/entity"
	financialProductRepo "financing-offer/internal/core/financialproduct/repository"
	orderServiceRepo "financing-offer/internal/core/orderservice/repository"
	promotionCampaignRepo "financing-offer/internal/core/promotion_campaign/repository"
	"financing-offer/internal/funcs"
)

type UseCase interface {
	GetOngoingPromotionLoanPackageIds(ctx context.Context) ([]int64, error)
	GetOnGoingPromotionLoanPackages(ctx context.Context) (entity.PromotionLoanPackage, error)
	SetPromotionLoanPackage(ctx context.Context, promotionLoanPackage entity.PromotionLoanPackage, updater string) (entity.PromotionLoanPackage, error)
	GetPromotionLoanPackageBySymbol(ctx context.Context, symbol, accountNo, custodyCode string) ([]entity.AccountLoanPackageWithAccountNo, error)
	GetPublicPromotionLoanPackageBySymbol(ctx context.Context, symbol string) (*entity.AccountLoanPackageWithSymbol, error)
	GetPublicPromotionLoanPackages(ctx context.Context) ([]entity.AccountLoanPackageWithSymbol, error)
	GetInvestorPromotionLoanPackages(ctx context.Context, accountNo, custodyCode string) (map[string][]entity.AccountLoanPackageWithSymbol, error)
	GetPromotionLoanPackages(ctx context.Context, accountNo, custodyCode string, symbol string) (map[string][]entity.LoanPackageWithCampaignProduct, error)
	GetPublicPromotionLoanPackagesWithCampaigns(ctx context.Context, symbol string) ([]entity.LoanPackageWithCampaignProduct, error)
}

type useCase struct {
	cfg                          config.AppConfig
	configurationPersistenceRepo repository.ConfigurationPersistenceRepository
	orderServiceRepo             orderServiceRepo.OrderServiceRepository
	financialProductRepo         financialProductRepo.FinancialProductRepository
	promotionCampaignRepo        promotionCampaignRepo.PromotionCampaignRepository
}

func NewUseCase(
	cfg config.AppConfig,
	configurationPersistenceRepo repository.ConfigurationPersistenceRepository,
	orderServiceRepo orderServiceRepo.OrderServiceRepository,
	financialProductRepo financialProductRepo.FinancialProductRepository,
	promotionCampaignRepo promotionCampaignRepo.PromotionCampaignRepository,
) UseCase {
	return &useCase{
		cfg:                          cfg,
		configurationPersistenceRepo: configurationPersistenceRepo,
		orderServiceRepo:             orderServiceRepo,
		financialProductRepo:         financialProductRepo,
		promotionCampaignRepo:        promotionCampaignRepo,
	}
}

func (u *useCase) GetOnGoingPromotionLoanPackages(ctx context.Context) (entity.PromotionLoanPackage, error) {
	promotionLoanPackage, err := u.configurationPersistenceRepo.GetPromotionConfiguration(ctx)
	if err != nil {
		return entity.PromotionLoanPackage{}, fmt.Errorf(
			"promotionLoanPackageUseCase GetOnGoingPromotionLoanProductIds %w", err,
		)
	}
	return promotionLoanPackage, nil
}

func (u *useCase) SetPromotionLoanPackage(ctx context.Context, promotionLoanPackage entity.PromotionLoanPackage, updater string) (entity.PromotionLoanPackage, error) {
	err := u.configurationPersistenceRepo.SetPromotionConfiguration(ctx, promotionLoanPackage, updater)
	if err != nil {
		return promotionLoanPackage, fmt.Errorf("promotionLoanPackageUseCase SetPromotionConfiguration %w", err)
	}
	return promotionLoanPackage, nil
}

func (u *useCase) GetOngoingPromotionLoanPackageIds(ctx context.Context) ([]int64, error) {
	loanPackageIdsFromEnv := u.cfg.BestPromotions.LoanPackageIds
	promotionLoanPackage, err := u.configurationPersistenceRepo.GetPromotionConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf("promotionLoanPackageUseCase GetOngoingPromotionLoanPackageIds %w", err)
	}
	loanPackageIds := funcs.Map(
		promotionLoanPackage.LoanProducts, func(loanProduct entity.PromotionLoanProduct) int64 {
			return loanProduct.LoanPackageId
		},
	)
	return funcs.UniqueElements(append(loanPackageIds, loanPackageIdsFromEnv...)), nil
}

func (u *useCase) GetPromotionLoanPackageBySymbol(ctx context.Context, symbol, accountNo, custodyCode string) ([]entity.AccountLoanPackageWithAccountNo, error) {
	errorTemplate := "promotionLoanPackageUseCase GetPromotionLoanPackageBySymbol %w"
	promotionPackage, err := u.configurationPersistenceRepo.GetPromotionConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	// get the configured promotion loan packages that have the symbol
	promotionSymbolPackagesSet := make(map[int64]bool)
	for _, promotion := range promotionPackage.LoanProducts {
		if allSymbols := promotion.AllSymbols(); slices.Contains(allSymbols, symbol) {
			promotionSymbolPackagesSet[promotion.LoanPackageId] = true
		}
	}
	// if no promotion loan package has the symbol, return nil
	if len(promotionSymbolPackagesSet) == 0 {
		return nil, nil
	}
	// get the user's assigned loan packages, and check whether the user is a margin user
	accountLoanPackagesByAccountNo, err := u.getUserAccountLoanPackages(ctx, custodyCode, accountNo)
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	result := make([]entity.AccountLoanPackageWithAccountNo, 0, len(accountLoanPackagesByAccountNo))
	for accountNo, accountLoanPackages := range accountLoanPackagesByAccountNo {

		if minimumInterestRateProduct := getMinimumInterestRateProductForAccountNumber(
			accountLoanPackages, promotionSymbolPackagesSet, symbol,
		); minimumInterestRateProduct != nil {
			result = append(
				result, entity.AccountLoanPackageWithAccountNo{
					AccountNo:          accountNo,
					AccountLoanPackage: *minimumInterestRateProduct,
				},
			)
		}
	}
	return result, nil
}

func (u *useCase) GetInvestorPromotionLoanPackages(ctx context.Context, accountNo, custodyCode string) (map[string][]entity.AccountLoanPackageWithSymbol, error) {
	errorTemplate := "promotionLoanPackageUseCase GetInvestorPromotionLoanPackages %w"
	promotionPackage, err := u.configurationPersistenceRepo.GetPromotionConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	groupedPromotionPackages := make(map[string]map[int64]bool) // map[symbol]map[loanPackageId]bool
	for _, promotion := range promotionPackage.LoanProducts {
		for _, symbol := range promotion.AllSymbols() {
			if _, ok := groupedPromotionPackages[symbol]; !ok {
				groupedPromotionPackages[symbol] = make(map[int64]bool)
			}
			groupedPromotionPackages[symbol][promotion.LoanPackageId] = true
		}
	}
	// get the user's assigned loan packages, and check whether the user is a margin user
	accountLoanPackagesByAccountNo, err := u.getUserAccountLoanPackages(ctx, custodyCode, accountNo)
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	result := make(map[string][]entity.AccountLoanPackageWithSymbol, len(groupedPromotionPackages))
	for accountNo, accountLoanPackages := range accountLoanPackagesByAccountNo {
		result[accountNo] = make([]entity.AccountLoanPackageWithSymbol, 0)

		for symbol, promotionSymbolPackagesSet := range groupedPromotionPackages {
			if minimumInterestRateProduct := getMinimumInterestRateProductForAccountNumber(
				accountLoanPackages, promotionSymbolPackagesSet, symbol,
			); minimumInterestRateProduct != nil {
				promotionPackages := result[accountNo]
				promotionPackages = append(
					promotionPackages, entity.AccountLoanPackageWithSymbol{
						AccountLoanPackage: *minimumInterestRateProduct,
						Symbol:             symbol,
					},
				)
				result[accountNo] = promotionPackages
			}
		}
	}
	return result, nil
}

// getUserAccountLoanPackages if accountNo is present, make sure it belongs to the custody code, otherwise return all investor's accountNos
func (u *useCase) getUserAccountLoanPackages(ctx context.Context, custodyCode, accountNo string) (map[string][]entity.AccountLoanPackage, error) {
	errorTemplate := "promotionLoanPackageUseCase getUserAccountLoanPackages %w"
	destAccountNos := make([]string, 0, 1)
	accounts, err := u.financialProductRepo.GetAllAccountDetailByCustodyCode(ctx, custodyCode)
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	accountNos := funcs.Map(
		accounts, func(account entity.FinancialAccountDetail) string {
			return account.AccountNo
		},
	)
	if accountNo != "" {
		if !slices.Contains(accountNos, accountNo) {
			return nil, apperrors.ErrAccountNoInvalid
		}
		destAccountNos = append(destAccountNos, accountNo)
	} else {
		destAccountNos = accountNos
	}
	investorAccountLoanPackages := make(map[string][]entity.AccountLoanPackage, len(destAccountNos))
	var (
		errGroup errgroup.Group
		mu       sync.Mutex
	)
	for _, accountNo := range destAccountNos {
		errGroup.Go(
			func() error {
				accountLoanPackages, err := u.orderServiceRepo.GetAllAccountLoanPackages(ctx, accountNo)
				if err != nil {
					return err
				}
				mu.Lock()
				investorAccountLoanPackages[accountNo] = accountLoanPackages
				mu.Unlock()
				return nil
			},
		)
	}
	if err := errGroup.Wait(); err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	return investorAccountLoanPackages, nil
}

func getMinimumInterestRateProductForAccountNumber(
	accountLoanPackages []entity.AccountLoanPackage,
	promotionSymbolPackagesSet map[int64]bool,
	symbol string,
) *entity.AccountLoanPackage {
	if !isMarginUser(accountLoanPackages) {
		return nil
	}
	var (
		resultLoanPackage          *entity.AccountLoanPackage
		minimumInterestRateProduct entity.LoanProduct
	)
	// find the loan package and margin product with minimum interest rate
	for _, accountLoanPackage := range accountLoanPackages {
		if _, ok := promotionSymbolPackagesSet[accountLoanPackage.Id]; !ok {
			continue
		}
		for _, product := range accountLoanPackage.LoanProducts {
			if product.Symbol != symbol {
				continue
			}
			if resultLoanPackage == nil {
				minimumInterestRateProduct = product
				resultLoanPackage = &accountLoanPackage
				resultLoanPackage.LoanProducts = []entity.LoanProduct{product}
				continue
			}
			if product.InterestRate.LessThan(minimumInterestRateProduct.InterestRate) {
				minimumInterestRateProduct = product
				resultLoanPackage = &accountLoanPackage
				resultLoanPackage.LoanProducts = []entity.LoanProduct{product}
			}
		}
	}
	return resultLoanPackage
}

func getMinimumInterestRateProduct(
	accountLoanPackages []entity.AccountLoanPackage,
	promotionSymbolPackagesSet map[int64]map[string]bool,
	symbol string,
) *entity.LoanProduct {
	if !isMarginUser(accountLoanPackages) {
		return nil
	}
	var result *entity.LoanProduct

	for _, accountLoanPackage := range accountLoanPackages {
		if promotionSymbolPackagesSet[accountLoanPackage.Id] == nil {
			continue
		}
		for _, product := range accountLoanPackage.LoanProducts {
			if product.Symbol != symbol {
				continue
			}
			if result == nil {
				result = &product
				continue
			}
			if product.InterestRate.LessThan(result.InterestRate) {
				result = &product
			}
		}
	}
	return result
}

type promotionLoanPackageContainer struct {
	minimumInterestRateProduct     entity.LoanProduct
	minimumInterestRateLoanPackage *entity.AccountLoanPackageWithSymbol
}

func (u *useCase) GetPublicPromotionLoanPackageBySymbol(ctx context.Context, symbol string) (*entity.AccountLoanPackageWithSymbol, error) {
	errorTemplate := "promotionLoanPackageUseCase GetPublicPromotionLoanPackageBySymbol %w"
	promotionPackage, err := u.configurationPersistenceRepo.GetPromotionConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	loanPackagesByBasketId, loanBaskets, err := u.prepareLoanBasketData(ctx, symbol, promotionPackage)
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	if len(loanBaskets) == 0 {
		return nil, nil
	}
	minimumInterestRateLoanPackage, err := u.calculateMinimumInterestRateLoanPackage(
		ctx, symbol, loanPackagesByBasketId, loanBaskets,
	)
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	return minimumInterestRateLoanPackage, nil
}

func (u *useCase) GetPublicPromotionLoanPackages(ctx context.Context) ([]entity.AccountLoanPackageWithSymbol, error) {
	errorTemplate := "promotionLoanPackageUseCase GetPublicPromotionLoanPackages %w"
	promotionPackage, err := u.configurationPersistenceRepo.GetPromotionConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	symbols := make(map[string]bool, len(promotionPackage.LoanProducts))
	for _, promotion := range promotionPackage.LoanProducts {
		for _, symbol := range promotion.RetailSymbols {
			symbols[symbol] = true
		}
	}
	var (
		result   []entity.AccountLoanPackageWithSymbol
		mu       sync.Mutex
		errGroup errgroup.Group
	)
	for symbol := range symbols {
		errGroup.Go(
			func() error {
				loanPackagesByBasketId, loanBaskets, err := u.prepareLoanBasketData(ctx, symbol, promotionPackage)
				if err != nil {
					return fmt.Errorf(errorTemplate, err)
				}
				minimumInterestRateLoanPackage, err := u.calculateMinimumInterestRateLoanPackage(
					ctx, symbol, loanPackagesByBasketId, loanBaskets,
				)
				if err != nil {
					return err
				}
				if minimumInterestRateLoanPackage != nil {
					mu.Lock()
					result = append(result, *minimumInterestRateLoanPackage)
					mu.Unlock()
				}
				return nil
			},
		)
	}
	if err := errGroup.Wait(); err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	return result, nil
}

func (u *useCase) GetPublicPromotionLoanPackagesWithCampaigns(ctx context.Context, symbol string) ([]entity.LoanPackageWithCampaignProduct, error) {
	errorTemplate := "promotionLoanPackageUseCase GetPublicPromotionLoanPackagesWithCampaigns %w"
	campaigns, err := u.promotionCampaignRepo.GetAll(ctx, entity.GetPromotionCampaignsRequest{
		Status: string(entity.Active),
	})
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	var products []entity.PromotionCampaignProduct
	for _, campaign := range campaigns {
		products = append(products, campaign.Metadata.Products...)
	}
	symbols := make(map[string]bool, len(products))
	for _, promotion := range products {
		for _, symbol := range promotion.RetailSymbols {
			symbols[symbol] = true
		}
	}
	campaignMapByLoanPackageAndSymbol := make(map[int64]map[string]*entity.PromotionCampaign)
	for _, campaign := range campaigns {
		for _, product := range campaign.Metadata.Products {
			if _, ok := campaignMapByLoanPackageAndSymbol[product.LoanPackageId]; !ok {
				campaignMapByLoanPackageAndSymbol[product.LoanPackageId] = make(map[string]*entity.PromotionCampaign)
			}
			for _, symbol := range product.RetailSymbols {
				campaignMapByLoanPackageAndSymbol[product.LoanPackageId][symbol] = &campaign
			}
		}
	}
	var (
		mu       sync.Mutex
		errGroup errgroup.Group
	)
	result := make([]entity.LoanPackageWithCampaignProduct, 0)
	if symbol != "" {
		loanPackagesByBasketId, loanBaskets, err := u.prepareBasketData(ctx, symbol, products)
		if err != nil {
			return nil, fmt.Errorf(errorTemplate, err)
		}
		if len(loanBaskets) == 0 {
			return result, nil
		}
		loanPackage, err := u.calculateMinimumInterestRateLoanPackage(
			ctx, symbol, loanPackagesByBasketId, loanBaskets,
		)
		products := make([]entity.CampaignWithProduct, 0)
		for _, product := range loanPackage.LoanProducts {
			if campaign := campaignMapByLoanPackageAndSymbol[loanPackage.Id][product.Symbol]; campaign != nil {
				products = append(
					products, entity.CampaignWithProduct{
						Product:  product,
						Campaign: entity.Campaign{Name: campaign.Name, Tag: campaign.Tag, Description: campaign.Description},
					},
				)
			}
		}
		if err != nil {
			return nil, err
		}
		if loanPackage != nil {
			result = append(result, entity.LoanPackageWithCampaignProduct{
				Id:                       loanPackage.Id,
				Name:                     loanPackage.Name,
				Type:                     loanPackage.Type,
				BrokerFirmBuyingFeeRate:  loanPackage.BrokerFirmBuyingFeeRate,
				BrokerFirmSellingFeeRate: loanPackage.BrokerFirmSellingFeeRate,
				TransferFee:              loanPackage.TransferFee,
				Description:              loanPackage.Description,
				BasketId:                 loanPackage.BasketId,
				CampaignProducts:         products,
			})
		}
		return result, nil
	}
	for symbol := range symbols {
		errGroup.Go(
			func() error {
				loanPackagesByBasketId, loanBaskets, err := u.prepareBasketData(ctx, symbol, products)
				if err != nil {
					return fmt.Errorf(errorTemplate, err)
				}
				if len(loanBaskets) == 0 {
					return nil
				}
				loanPackage, err := u.calculateMinimumInterestRateLoanPackage(
					ctx, symbol, loanPackagesByBasketId, loanBaskets,
				)
				if err != nil {
					return err
				}
				if loanPackage != nil {
					products := make([]entity.CampaignWithProduct, 0)
					for _, product := range loanPackage.LoanProducts {
						if campaign := campaignMapByLoanPackageAndSymbol[loanPackage.Id][product.Symbol]; campaign != nil {
							products = append(
								products, entity.CampaignWithProduct{
									Product:  product,
									Campaign: entity.Campaign{Name: campaign.Name, Tag: campaign.Tag, Description: campaign.Description},
								},
							)
						}
					}
					mu.Lock()
					result = append(result, entity.LoanPackageWithCampaignProduct{
						Id:                       loanPackage.Id,
						Name:                     loanPackage.Name,
						Type:                     loanPackage.Type,
						BrokerFirmBuyingFeeRate:  loanPackage.BrokerFirmBuyingFeeRate,
						BrokerFirmSellingFeeRate: loanPackage.BrokerFirmSellingFeeRate,
						TransferFee:              loanPackage.TransferFee,
						Description:              loanPackage.Description,
						BasketId:                 loanPackage.BasketId,
						CampaignProducts:         products,
					})
					mu.Unlock()
				}
				return nil
			},
		)
	}
	if err := errGroup.Wait(); err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	loanPackageMap := funcs.AssociateBy[entity.LoanPackageWithCampaignProduct, int64](
		result, func(group entity.LoanPackageWithCampaignProduct) int64 {
			return group.Id
		},
	)
	loanPackageGroup := make(map[int64][]entity.LoanPackageWithCampaignProduct, len(loanPackageMap))
	for _, loanPackage := range result {
		if _, ok := loanPackageGroup[loanPackage.Id]; !ok {
			loanPackageGroup[loanPackage.Id] = make([]entity.LoanPackageWithCampaignProduct, 0)
		}
		loanPackageGroup[loanPackage.Id] = append(loanPackageGroup[loanPackage.Id], loanPackage)
	}
	result = make([]entity.LoanPackageWithCampaignProduct, 0, len(loanPackageGroup))
	for id, loanPackage := range loanPackageMap {
		products := loanPackageGroup[id]
		campaignWithProducts := make([]entity.CampaignWithProduct, 0, len(products))
		for _, product := range products {
			campaignWithProducts = append(campaignWithProducts, product.CampaignProducts...)
		}
		loanPackage.CampaignProducts = campaignWithProducts
		result = append(result, loanPackage)
	}
	return result, nil
}

func (u *useCase) GetPromotionLoanPackages(ctx context.Context, accountNo, custodyCode string, symbol string) (map[string][]entity.LoanPackageWithCampaignProduct, error) {
	errorTemplate := "promotionLoanPackageUseCase GetPromotionLoanPackages %w"
	campaigns, err := u.promotionCampaignRepo.GetAll(ctx, entity.GetPromotionCampaignsRequest{
		Status: string(entity.Active),
	})
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	var products []entity.PromotionCampaignProduct
	for _, campaign := range campaigns {
		products = append(products, campaign.Metadata.Products...)
	}
	groupedPromotionPackages := make(map[int64]map[string]bool)
	for _, product := range products {
		for _, symbol := range product.Symbols {
			if _, ok := groupedPromotionPackages[product.LoanPackageId]; !ok {
				groupedPromotionPackages[product.LoanPackageId] = make(map[string]bool)
			}
			groupedPromotionPackages[product.LoanPackageId][symbol] = true
		}
	}
	accountLoanPackagesByAccountNo, err := u.getUserAccountLoanPackages(ctx, custodyCode, accountNo)
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	result := make(map[string][]entity.LoanPackageWithCampaignProduct, len(accountLoanPackagesByAccountNo))
	campaignMapByLoanPackageAndSymbol := make(map[int64]map[string]*entity.PromotionCampaign)
	for _, campaign := range campaigns {
		for _, product := range campaign.Metadata.Products {
			if _, ok := campaignMapByLoanPackageAndSymbol[product.LoanPackageId]; !ok {
				campaignMapByLoanPackageAndSymbol[product.LoanPackageId] = make(map[string]*entity.PromotionCampaign)
			}
			for _, symbol := range product.Symbols {
				campaignMapByLoanPackageAndSymbol[product.LoanPackageId][symbol] = &campaign
			}
		}
	}
	for accountNo, accountLoanPackages := range accountLoanPackagesByAccountNo {
		result[accountNo] = make([]entity.LoanPackageWithCampaignProduct, 0)
		loanPackages := make([]entity.AccountLoanPackage, 0)
		for _, loanPackage := range accountLoanPackages {
			if groupedPromotionPackages[loanPackage.Id] == nil {
				continue
			}
			products := funcs.Filter(loanPackage.LoanProducts, func(product entity.LoanProduct, _ int) bool {
				return symbol == "" || product.Symbol == symbol
			})
			for _, product := range products {
				minProduct := getMinimumInterestRateProduct(accountLoanPackages, groupedPromotionPackages, product.Symbol)
				if groupedPromotionPackages[loanPackage.Id][product.Symbol] && minProduct != nil && minProduct.Id == product.Id {
					loanPackages = append(loanPackages, loanPackage)
					break
				}
			}
		}
		for _, loanPackage := range loanPackages {
			campaignProducts := make([]entity.CampaignWithProduct, 0)
			products := funcs.Filter(loanPackage.LoanProducts, func(product entity.LoanProduct, _ int) bool {
				return symbol == "" || product.Symbol == symbol
			})
			for _, product := range products {
				if campaign := campaignMapByLoanPackageAndSymbol[loanPackage.Id][product.Symbol]; campaign != nil {
					campaignProducts = append(
						campaignProducts, entity.CampaignWithProduct{
							Product:  product,
							Campaign: entity.Campaign{Name: campaign.Name, Tag: campaign.Tag, Description: campaign.Description},
						},
					)
				}

			}
			result[accountNo] = append(
				result[accountNo], entity.LoanPackageWithCampaignProduct{
					Id:                       loanPackage.Id,
					Name:                     loanPackage.Name,
					Type:                     loanPackage.Type,
					BrokerFirmBuyingFeeRate:  loanPackage.BrokerFirmBuyingFeeRate,
					BrokerFirmSellingFeeRate: loanPackage.BrokerFirmSellingFeeRate,
					TransferFee:              loanPackage.TransferFee,
					Description:              loanPackage.Description,
					BasketId:                 loanPackage.BasketId,
					CampaignProducts:         campaignProducts,
				},
			)
		}
	}
	return result, nil
}

func (u *useCase) calculateMinimumInterestRateLoanPackage(
	ctx context.Context,
	symbol string,
	loanPackagesByBasketId map[int64]entity.FinancialProductLoanPackage,
	applicableLoanBaskets []entity.MarginBasket,
) (*entity.AccountLoanPackageWithSymbol, error) {
	if len(applicableLoanBaskets) == 0 {
		return nil, nil
	}
	validProducts, err := u.findValidProductsBySymbol(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("promotionLoanPackageUseCase calculateMinimumInterestRateLoanPackage %w", err)
	}
	container := promotionLoanPackageContainer{
		minimumInterestRateProduct:     entity.LoanProduct{InterestRate: decimal.NewFromInt(1)},
		minimumInterestRateLoanPackage: nil,
	}
	for _, loanBasket := range applicableLoanBaskets {
		// v2 loan package
		if len(loanBasket.Symbols) > 0 {
			calculatePromotionV2LoanPackage(&container, symbol, loanPackagesByBasketId, loanBasket)
		}
		// v3 loan package
		if len(loanBasket.LoanProducts) == 0 {
			continue
		}
		if err := calculatePromotionV3LoanPackage(
			validProducts, &container, symbol, loanPackagesByBasketId, loanBasket,
		); err != nil {
			return nil, fmt.Errorf("promotionLoanPackageUseCase calculateMinimumInterestRateLoanPackage %w", err)
		}
	}
	return container.minimumInterestRateLoanPackage, nil
}

func calculatePromotionV2LoanPackage(
	container *promotionLoanPackageContainer,
	symbol string,
	loanPackagesByBasketId map[int64]entity.FinancialProductLoanPackage,
	loanBasket entity.MarginBasket,
) {
	loanPackage, ok := loanPackagesByBasketId[loanBasket.Id]
	if !ok {
		return
	}
	if loanPackage.InterestRate.GreaterThanOrEqual(container.minimumInterestRateProduct.InterestRate) {
		return
	}
	container.minimumInterestRateProduct = entity.LoanProduct{
		Symbol:                   symbol,
		InitialRate:              loanPackage.InitialRate,
		InitialRateForWithdraw:   loanPackage.InitialRateForWithdraw,
		MaintenanceRate:          loanPackage.MaintenanceRate,
		LiquidRate:               loanPackage.LiquidRate,
		InterestRate:             loanPackage.InterestRate,
		PreferentialPeriod:       loanPackage.PreferentialPeriod,
		PreferentialInterestRate: loanPackage.PreferentialInterestRate,
		Term:                     loanPackage.Term,
		AllowExtendLoanTerm:      loanPackage.AllowExtendLoanTerm,
		AllowEarlyPayment:        loanPackage.AllowEarlyPayment,
	}
	container.minimumInterestRateLoanPackage = &entity.AccountLoanPackageWithSymbol{
		AccountLoanPackage: entity.AccountLoanPackage{
			Id:                       loanPackage.Id,
			Name:                     loanPackage.Name,
			Type:                     loanPackage.LoanType,
			LoanProducts:             []entity.LoanProduct{container.minimumInterestRateProduct},
			BrokerFirmBuyingFeeRate:  loanPackage.BuyingFeeRate,
			BrokerFirmSellingFeeRate: loanPackage.BrokerFirmSellingFeeRate,
			TransferFee:              loanPackage.TransferFee,
			Description:              loanPackage.Description,
			BasketId:                 loanPackage.LoanBasketId,
		},
		Symbol: symbol,
	}
}

func calculatePromotionV3LoanPackage(
	validProducts map[int64]entity.MarginProduct,
	container *promotionLoanPackageContainer,
	symbol string,
	loanPackagesByBasketId map[int64]entity.FinancialProductLoanPackage,
	loanBasket entity.MarginBasket,
) error {
	for _, loanProduct := range loanBasket.LoanProducts {
		if _, ok := validProducts[loanProduct.Id]; !ok {
			continue
		}
		loanPackage, ok := loanPackagesByBasketId[loanBasket.Id]
		if !ok {
			continue
		}
		propagatedProduct := propagateMarginProduct(loanProduct)
		if propagatedProduct.InterestRate.GreaterThanOrEqual(container.minimumInterestRateProduct.InterestRate) {
			continue
		}
		container.minimumInterestRateProduct = propagatedProduct
		container.minimumInterestRateLoanPackage = &entity.AccountLoanPackageWithSymbol{
			AccountLoanPackage: entity.AccountLoanPackage{
				Id:                       loanPackage.Id,
				Name:                     loanPackage.Name,
				Type:                     loanPackage.LoanType,
				LoanProducts:             []entity.LoanProduct{propagatedProduct},
				BrokerFirmBuyingFeeRate:  loanPackage.BuyingFeeRate,
				BrokerFirmSellingFeeRate: loanPackage.BrokerFirmSellingFeeRate,
				TransferFee:              loanPackage.TransferFee,
				Description:              loanPackage.Description,
				BasketId:                 loanPackage.LoanBasketId,
			},
			Symbol: symbol,
		}
	}
	return nil
}

func (u *useCase) prepareLoanBasketData(ctx context.Context, symbol string, promotionPackage entity.PromotionLoanPackage) (map[int64]entity.FinancialProductLoanPackage, []entity.MarginBasket, error) {
	// get the configured promotion loan packages that have the symbol
	promotionSymbolPackages := make([]int64, 0)
	for _, promotion := range promotionPackage.LoanProducts {
		if slices.Contains(promotion.RetailSymbols, symbol) {
			promotionSymbolPackages = append(promotionSymbolPackages, promotion.LoanPackageId)
		}
	}
	// if no promotion loan package has the symbol, return nil
	if len(promotionSymbolPackages) == 0 {
		return nil, nil, nil
	}
	loanPackages, err := u.financialProductRepo.GetLoanPackageDetails(ctx, promotionSymbolPackages)
	if err != nil {
		return nil, nil, err
	}
	loanBasketIds := funcs.Map(
		loanPackages, func(loanPackage entity.FinancialProductLoanPackage) int64 {
			return loanPackage.LoanBasketId
		},
	)
	loanBaskets, err := u.financialProductRepo.GetMarginBasketsByIds(ctx, loanBasketIds)
	if err != nil {
		return nil, nil, err
	}
	loanPackagesByBasketId := make(map[int64]entity.FinancialProductLoanPackage, len(loanBaskets))
	for _, loanPackage := range loanPackages {
		existedPackage, ok := loanPackagesByBasketId[loanPackage.LoanBasketId]
		if !ok {
			loanPackagesByBasketId[loanPackage.LoanBasketId] = loanPackage
			continue
		}
		if loanPackage.InterestRate.LessThan(existedPackage.InterestRate) {
			loanPackagesByBasketId[loanPackage.LoanBasketId] = loanPackage
		}
	}
	return loanPackagesByBasketId, loanBaskets, nil
}

func (u *useCase) prepareBasketData(ctx context.Context, symbol string, products []entity.PromotionCampaignProduct) (map[int64]entity.FinancialProductLoanPackage, []entity.MarginBasket, error) {
	// get the configured promotion loan packages that have the symbol
	promotionSymbolPackages := make([]int64, 0)
	for _, promotion := range products {
		if slices.Contains(promotion.RetailSymbols, symbol) {
			promotionSymbolPackages = append(promotionSymbolPackages, promotion.LoanPackageId)
		}
	}
	// if no promotion loan package has the symbol, return nil
	if len(promotionSymbolPackages) == 0 {
		return nil, nil, nil
	}
	loanPackages, err := u.financialProductRepo.GetLoanPackageDetails(ctx, promotionSymbolPackages)
	if err != nil {
		return nil, nil, err
	}
	loanBasketIds := funcs.Map(
		loanPackages, func(loanPackage entity.FinancialProductLoanPackage) int64 {
			return loanPackage.LoanBasketId
		},
	)
	loanBaskets, err := u.financialProductRepo.GetMarginBasketsByIds(ctx, loanBasketIds)
	if err != nil {
		return nil, nil, err
	}
	loanPackagesByBasketId := make(map[int64]entity.FinancialProductLoanPackage, len(loanBaskets))
	for _, loanPackage := range loanPackages {
		existedPackage, ok := loanPackagesByBasketId[loanPackage.LoanBasketId]
		if !ok {
			loanPackagesByBasketId[loanPackage.LoanBasketId] = loanPackage
			continue
		}
		if loanPackage.InterestRate.LessThan(existedPackage.InterestRate) {
			loanPackagesByBasketId[loanPackage.LoanBasketId] = loanPackage
		}
	}
	return loanPackagesByBasketId, loanBaskets, nil
}

func propagateMarginProduct(marginProduct entity.MarginProduct) entity.LoanProduct {
	minimumInterestRatePolicy := entity.FinancialProductLoanPolicy{InterestRate: decimal.NewFromInt(1)}
	for _, financialLoanPolicy := range marginProduct.LoanPolicies {
		if financialLoanPolicy.LoanPolicy.InterestRate.LessThan(minimumInterestRatePolicy.InterestRate) {
			minimumInterestRatePolicy = financialLoanPolicy.LoanPolicy
		}
	}
	return entity.LoanProduct{
		Id:                       strconv.FormatInt(marginProduct.Id, 10),
		Name:                     marginProduct.Name,
		Symbol:                   marginProduct.Symbol,
		InitialRate:              marginProduct.LoanRate.InitialRate,
		InitialRateForWithdraw:   marginProduct.LoanRate.InitialRateForWithdraw,
		MaintenanceRate:          marginProduct.LoanRate.MaintenanceRate,
		LiquidRate:               marginProduct.LoanRate.LiquidRate,
		InterestRate:             minimumInterestRatePolicy.InterestRate,
		PreferentialPeriod:       minimumInterestRatePolicy.PreferentialPeriod,
		PreferentialInterestRate: minimumInterestRatePolicy.PreferentialInterestRate,
		Term:                     minimumInterestRatePolicy.Term,
		AllowExtendLoanTerm:      minimumInterestRatePolicy.AllowExtendLoanTerm,
		AllowEarlyPayment:        minimumInterestRatePolicy.AllowEarlyPayment,
	}
}

func (u *useCase) findValidProductsBySymbol(ctx context.Context, symbol string) (map[int64]entity.MarginProduct, error) {
	allProductWithProvidedSymbol, err := u.financialProductRepo.GetLoanProducts(
		ctx, entity.MarginProductFilter{Symbol: symbol},
	)
	if err != nil {
		return nil, err
	}
	return funcs.AssociateBy(
		allProductWithProvidedSymbol, func(item entity.MarginProduct) int64 {
			return item.Id
		},
	), nil
}

func isMarginUser(loanAccounts []entity.AccountLoanPackage) bool {
	for _, loanAccount := range loanAccounts {
		if loanAccount.Type == "M" {
			return true
		}
	}
	return false
}
