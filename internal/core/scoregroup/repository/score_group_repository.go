package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type ScoreGroupRepository interface {
	GetAll(ctx context.Context) ([]entity.ScoreGroup, error)
	Create(ctx context.Context, scoreGroup entity.ScoreGroup) (entity.ScoreGroup, error)
	Update(ctx context.Context, scoreGroup entity.ScoreGroup) (entity.ScoreGroup, error)
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (entity.ScoreGroup, error)
}
