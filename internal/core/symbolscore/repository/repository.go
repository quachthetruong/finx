package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type SymbolScoreRepository interface {
	GetAll(ctx context.Context, filter entity.SymbolScoreFilter) ([]entity.SymbolScore, error)
	Update(ctx context.Context, symbolScore entity.SymbolScore) (entity.SymbolScore, error)
	Create(ctx context.Context, symbolScore entity.SymbolScore) (entity.SymbolScore, error)
	GetCurrentScoreForSymbol(ctx context.Context, symbolId int64) (entity.SymbolScore, error)
	GetById(ctx context.Context, id int64) (entity.SymbolScore, error)
}
