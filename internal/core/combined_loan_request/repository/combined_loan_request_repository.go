package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type CombinedLoanPackageRequestPersistenceRepository interface {
	GetAll(ctx context.Context, filter entity.CombinedLoanRequestFilter) ([]entity.CombinedLoanRequest, error)
	Count(ctx context.Context, filter entity.CombinedLoanRequestFilter) (int64, error)
}
