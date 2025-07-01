package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type LoanRequestSchedulerConfigRepository interface {
	GetAll(ctx context.Context) ([]entity.LoanRequestSchedulerConfig, error)
	Create(ctx context.Context, config entity.LoanRequestSchedulerConfig) (entity.LoanRequestSchedulerConfig, error)
	GetCurrentConfig(ctx context.Context) (entity.LoanRequestSchedulerConfig, error)
}
