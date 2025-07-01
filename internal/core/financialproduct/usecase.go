package financialproduct

import (
	"context"
	configRepo "financing-offer/internal/config/repository"
	"fmt"

	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/financialproduct/repository"
	repositomarginOperationRepoy "financing-offer/internal/core/marginoperation/repository"
	"financing-offer/internal/funcs"
)

type UseCase interface {
	GetLoanRates(ctx context.Context) ([]entity.LoanRate, error)
	GetMarginPools(ctx context.Context) ([]entity.MarginPool, error)
}

type useCase struct {
	cfg                          config.AppConfig
	financialProductRepository   repository.FinancialProductRepository
	marginOperationRepository    repositomarginOperationRepoy.MarginOperationRepository
	configurationPersistenceRepo configRepo.ConfigurationPersistenceRepository
}

func (u *useCase) GetLoanRates(ctx context.Context) ([]entity.LoanRate, error) {
	errTemplate := "GetLoanRates %w"
	loanRate, err := u.configurationPersistenceRepo.GetLoanRateConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}
	loanRates, err := u.financialProductRepository.GetLoanRatesByIds(ctx, loanRate.Ids)
	if err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}
	return loanRates, nil
}

func (u *useCase) GetMarginPools(ctx context.Context) ([]entity.MarginPool, error) {
	errTemplate := "GetMarginPools %w"
	marginPool, err := u.configurationPersistenceRepo.GetMarginPoolConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}
	if len(marginPool.Ids) == 0 {
		return []entity.MarginPool{}, nil
	}
	marginPools, err := u.marginOperationRepository.GetMarginPoolsByIds(ctx, marginPool.Ids)
	if err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}
	poolGroupIds := funcs.Map(marginPools, func(m entity.MarginPool) int64 { return m.PoolGroupId })
	marginPoolGroups, err := u.marginOperationRepository.GetMarginPoolGroupsByIds(ctx, poolGroupIds)
	if err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}
	mapped := make(map[int64]entity.MarginPoolGroup)
	for _, group := range marginPoolGroups {
		mapped[group.Id] = group
	}
	for i, pool := range marginPools {
		if group, ok := mapped[pool.PoolGroupId]; ok {
			marginPools[i].Group = group
		}
	}
	return marginPools, nil
}

func NewUseCase(cfg config.AppConfig,
	financialProductRepository repository.FinancialProductRepository,
	marginOperationRepository repositomarginOperationRepoy.MarginOperationRepository,
	configurationPersistenceRepo configRepo.ConfigurationPersistenceRepository,
) UseCase {
	return &useCase{
		cfg:                          cfg,
		financialProductRepository:   financialProductRepository,
		marginOperationRepository:    marginOperationRepository,
		configurationPersistenceRepo: configurationPersistenceRepo,
	}
}
