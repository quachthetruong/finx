package configuration

import (
	"context"
	"financing-offer/internal/config"
	"financing-offer/internal/config/repository"
	"financing-offer/internal/core/entity"
	"fmt"
)

type UseCase interface {
	SetLoanRate(ctx context.Context, loanRate entity.LoanRateConfiguration, updater string) (entity.LoanRateConfiguration, error)
	GetLoanRate(ctx context.Context) (entity.LoanRateConfiguration, error)
	SetMarginPool(ctx context.Context, marginPool entity.MarginPoolConfiguration, updater string) (entity.MarginPoolConfiguration, error)
	GetMarginPool(ctx context.Context) (entity.MarginPoolConfiguration, error)
}

type useCase struct {
	cfg                          config.AppConfig
	configurationPersistenceRepo repository.ConfigurationPersistenceRepository
}

func NewUseCase(
	cfg config.AppConfig,
	configurationPersistenceRepo repository.ConfigurationPersistenceRepository,
) UseCase {
	return &useCase{
		cfg:                          cfg,
		configurationPersistenceRepo: configurationPersistenceRepo,
	}
}

func (u *useCase) SetLoanRate(ctx context.Context, loanRate entity.LoanRateConfiguration, updater string) (entity.LoanRateConfiguration, error) {
	err := u.configurationPersistenceRepo.SetLoanRateConfiguration(ctx, loanRate, updater)
	if err != nil {
		return entity.LoanRateConfiguration{}, fmt.Errorf("SetLoanRateConfiguration %w", err)
	}
	return loanRate, nil
}

func (u *useCase) GetLoanRate(ctx context.Context) (entity.LoanRateConfiguration, error) {
	result, err := u.configurationPersistenceRepo.GetLoanRateConfiguration(ctx)
	if err != nil {
		return entity.LoanRateConfiguration{}, fmt.Errorf("GetLoanRate %w", err)
	}
	return result, nil
}

func (u *useCase) SetMarginPool(ctx context.Context, marginPool entity.MarginPoolConfiguration, updater string) (entity.MarginPoolConfiguration, error) {
	err := u.configurationPersistenceRepo.SetMarginPoolConfiguration(ctx, marginPool, updater)
	if err != nil {
		return entity.MarginPoolConfiguration{}, fmt.Errorf("SetMarginPool %w", err)
	}
	return marginPool, nil
}

func (u *useCase) GetMarginPool(ctx context.Context) (entity.MarginPoolConfiguration, error) {
	result, err := u.configurationPersistenceRepo.GetMarginPoolConfiguration(ctx)
	if err != nil {
		return entity.MarginPoolConfiguration{}, fmt.Errorf("GetMarginPool %w", err)
	}
	return result, nil
}
