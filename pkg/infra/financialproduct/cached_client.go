package financialproduct

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	"financing-offer/internal/core/entity"
	"financing-offer/pkg/cache"
)

type CachedFinancialProductRepository struct {
	*Client
	cacheStore cache.Cache
}

func NewCachedFinancialProductRepository(client *Client, cacheStore cache.Cache) *CachedFinancialProductRepository {
	return &CachedFinancialProductRepository{
		Client:     client,
		cacheStore: cacheStore,
	}
}

func (r *CachedFinancialProductRepository) GetMarginBasketsByIds(ctx context.Context, ids []int64) ([]entity.MarginBasket, error) {
	var (
		errorGroup errgroup.Group
		mu         sync.Mutex
		baskets    = make([]entity.MarginBasket, 0, len(ids))
	)
	for _, id := range ids {
		errorGroup.Go(
			func() error {
				basket, err := r.GetMarginBasketDetail(ctx, id)
				if err != nil {
					return err
				}
				mu.Lock()
				baskets = append(baskets, basket)
				mu.Unlock()
				return nil
			},
		)
	}
	if err := errorGroup.Wait(); err != nil {
		return nil, fmt.Errorf("GetMarginBasketsByIds %w", err)
	}
	return baskets, nil
}

func (r *CachedFinancialProductRepository) GetMarginBasketDetail(ctx context.Context, id int64) (entity.MarginBasket, error) {
	return cache.Do[entity.MarginBasket](
		r.cacheStore, fmt.Sprintf("financial_product_margin_basket_%d", id), cache.DefaultTtl,
		func() (entity.MarginBasket, error) {
			return r.Client.GetMarginBasketDetail(ctx, id)
		},
	)
}

func (r *CachedFinancialProductRepository) GetLoanPackageDetails(ctx context.Context, loanPackageIds []int64) ([]entity.FinancialProductLoanPackage, error) {
	var (
		errorGroup   errgroup.Group
		mu           sync.Mutex
		loanPackages = make([]entity.FinancialProductLoanPackage, 0, len(loanPackageIds))
	)
	for _, id := range loanPackageIds {
		errorGroup.Go(
			func() error {
				loanPackage, err := r.GetLoanPackageDetail(ctx, id)
				if err != nil {
					return err
				}
				mu.Lock()
				loanPackages = append(loanPackages, loanPackage)
				mu.Unlock()
				return nil
			},
		)
	}
	if err := errorGroup.Wait(); err != nil {
		return nil, fmt.Errorf("GetLoanPackageDetails %w", err)
	}
	return loanPackages, nil
}

func (r *CachedFinancialProductRepository) GetLoanPackageDetail(ctx context.Context, loanPackageId int64) (entity.FinancialProductLoanPackage, error) {
	return cache.Do[entity.FinancialProductLoanPackage](
		r.cacheStore, fmt.Sprintf("financial_product_loan_package_%d", loanPackageId), cache.DefaultTtl,
		func() (entity.FinancialProductLoanPackage, error) {
			return r.Client.GetLoanPackageDetail(ctx, loanPackageId)
		},
	)
}

func (r *CachedFinancialProductRepository) GetLoanProducts(ctx context.Context, filter entity.MarginProductFilter) ([]entity.MarginProduct, error) {
	return cache.Do[[]entity.MarginProduct](
		r.cacheStore, fmt.Sprintf("financial_product_margin_products_%s", filter), cache.DefaultTtl,
		func() ([]entity.MarginProduct, error) {
			return r.Client.GetLoanProducts(ctx, filter)
		},
	)
}
