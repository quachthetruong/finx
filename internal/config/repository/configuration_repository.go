package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type ConfigurationPersistenceRepository interface {
	GetPromotionConfiguration(ctx context.Context) (entity.PromotionLoanPackage, error)
	SetPromotionConfiguration(ctx context.Context, promotionLoanPackage entity.PromotionLoanPackage, updater string) error
	SetLoanRateConfiguration(ctx context.Context, loanRate entity.LoanRateConfiguration, updater string) error
	GetLoanRateConfiguration(ctx context.Context) (entity.LoanRateConfiguration, error)
	SetMarginPoolConfiguration(ctx context.Context, marginPool entity.MarginPoolConfiguration, updater string) error
	GetMarginPoolConfiguration(ctx context.Context) (entity.MarginPoolConfiguration, error)
	SetSubmissionDefault(ctx context.Context, defaultValue entity.SubmissionDefault, updater string) error
	GetSubmissionDefault(ctx context.Context) (entity.SubmissionDefault, error)
}
