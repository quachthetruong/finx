package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type MarginOperationRepository interface {
	GetMarginPoolById(ctx context.Context, marginPoolId int64) (entity.MarginPool, error)
	GetMarginPoolsByIds(ctx context.Context, marginPoolIds []int64) ([]entity.MarginPool, error)
	GetMarginPoolGroupsByIds(ctx context.Context, marginPoolGroupIds []int64) ([]entity.MarginPoolGroup, error)
}
