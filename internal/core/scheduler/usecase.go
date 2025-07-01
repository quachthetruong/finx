package scheduler

import (
	"context"
	"fmt"
	"log/slog"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/scheduler/repository"
)

type UseCase interface {
	GetAllLoanRequestSchedulerConfig(ctx context.Context) ([]entity.LoanRequestSchedulerConfig, error)
	CreateLoanRequestSchedulerConfig(ctx context.Context, entity entity.LoanRequestSchedulerConfig) (entity.LoanRequestSchedulerConfig, error)
	GetCurrentLoanRequestSchedulerConfig(ctx context.Context) (entity.LoanRequestSchedulerConfig, error)
}

func NewSchedulerUseCase(repo repository.LoanRequestSchedulerConfigRepository, logger *slog.Logger) UseCase {
	return &schedulerUseCase{
		repo:   repo,
		logger: logger,
	}
}

type schedulerUseCase struct {
	repo   repository.LoanRequestSchedulerConfigRepository
	logger *slog.Logger
}

func (u *schedulerUseCase) GetAllLoanRequestSchedulerConfig(ctx context.Context) ([]entity.LoanRequestSchedulerConfig, error) {
	res, err := u.repo.GetAll(ctx)
	if err != nil {
		return res, fmt.Errorf("schedulerUseCase GetAllLoanRequestSchedulerConfig %w", err)
	}
	return res, nil
}

func (u *schedulerUseCase) CreateLoanRequestSchedulerConfig(ctx context.Context, entity entity.LoanRequestSchedulerConfig) (entity.LoanRequestSchedulerConfig, error) {
	res, err := u.repo.Create(ctx, entity)
	if err != nil {
		return res, fmt.Errorf("schedulerUseCase CreateLoanRequestSchedulerConfig %w", err)
	}
	return res, nil
}

func (u *schedulerUseCase) GetCurrentLoanRequestSchedulerConfig(ctx context.Context) (entity.LoanRequestSchedulerConfig, error) {
	config, err := u.repo.GetCurrentConfig(ctx)
	if err != nil {
		return entity.LoanRequestSchedulerConfig{}, fmt.Errorf("schedulerUseCase GetCurrentLoanRequestSchedulerConfig %w", err)
	}
	return config, nil
}
