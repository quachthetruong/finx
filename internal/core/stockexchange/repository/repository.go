package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type StockExchangeRepository interface {
	GetAll(ctx context.Context) ([]entity.StockExchange, error)
	Create(ctx context.Context, stockExchange entity.StockExchange) (entity.StockExchange, error)
	Update(ctx context.Context, stockExchange entity.StockExchange) (entity.StockExchange, error)
	Delete(ctx context.Context, id int64) error
	GetBySymbolId(ctx context.Context, symbolId int64) (entity.StockExchange, error)
	GetById(ctx context.Context, id int64) (entity.StockExchange, error)
}
