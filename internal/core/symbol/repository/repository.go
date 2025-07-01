package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type SymbolRepository interface {
	GetAll(ctx context.Context, filter entity.SymbolFilter) ([]entity.Symbol, error)
	GetById(ctx context.Context, id int64) (entity.Symbol, error)
	GetBySymbol(ctx context.Context, symbol string) (entity.Symbol, error)
	Create(ctx context.Context, symbol entity.Symbol) (entity.Symbol, error)
	Update(ctx context.Context, symbol entity.Symbol) (entity.Symbol, error)
	Count(ctx context.Context, filter entity.SymbolFilter) (int64, error)
	GetSymbolWithActiveBlacklist(ctx context.Context, symbol string) (entity.Symbol, error)
}
