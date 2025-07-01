package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type ScoreGroupInterestRepository interface {
	GetAll(ctx context.Context, filter entity.ScoreGroupInterestFilter) ([]entity.ScoreGroupInterest, error)
	GetById(ctx context.Context, id int64) (entity.ScoreGroupInterest, error)
	Create(ctx context.Context, symbol entity.ScoreGroupInterest) (entity.ScoreGroupInterest, error)
	Update(ctx context.Context, symbol entity.ScoreGroupInterest) (entity.ScoreGroupInterest, error)
	Delete(ctx context.Context, id int64) (bool, error)
	GetAvailablePackageBySymbolId(ctx context.Context, symbolId int64) ([]entity.ScoreGroupInterest, error)
	GetAvailableScoreInterestsByScoreGroupId(ctx context.Context, scoreGroupId int64) ([]entity.ScoreGroupInterest, error)
}
