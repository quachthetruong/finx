package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type AwaitingConfirmRequestPersistenceRepository interface {
	GetAll(ctx context.Context, filter entity.AwaitingConfirmRequestFilter) ([]entity.AwaitingConfirmRequest, error)
	Count(ctx context.Context, filter entity.AwaitingConfirmRequestFilter) (int64, error)
	CountForStatistic(ctx context.Context) (int64, error)
}
