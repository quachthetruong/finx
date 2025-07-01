package repository

import (
	"context"
	"time"

	"financing-offer/internal/core/entity"
)

type BlackListSymbolRepository interface {
	GetAll(ctx context.Context, filter entity.BlacklistSymbolFilter) ([]entity.BlacklistSymbol, error)
	Create(ctx context.Context, symbol entity.BlacklistSymbol) (entity.BlacklistSymbol, error)
	Update(ctx context.Context, symbol entity.BlacklistSymbol) (entity.BlacklistSymbol, error)
	GetById(ctx context.Context, id int64) (entity.BlacklistSymbol, error)
	GetByAffectTime(ctx context.Context, symbolId int64, affectedFrom time.Time, affectedTo time.Time) ([]entity.BlacklistSymbol, error)
}
