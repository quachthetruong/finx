package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/enum"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

type BlackListSymbolRepository struct {
	getDbFunc database.GetDbFunc
}

func NewBlackListSymbolRepository(getDbFunc database.GetDbFunc) *BlackListSymbolRepository {
	return &BlackListSymbolRepository{
		getDbFunc: getDbFunc,
	}
}

func (b *BlackListSymbolRepository) GetAll(ctx context.Context, filter entity.BlacklistSymbolFilter) ([]entity.BlacklistSymbol, error) {
	stm := table.BlacklistSymbol.SELECT(table.BlacklistSymbol.AllColumns)
	if filter.Symbol.IsPresent() {
		stm = stm.FROM(
			table.BlacklistSymbol.INNER_JOIN(
				table.Symbol, table.BlacklistSymbol.SymbolID.EQ(table.Symbol.ID),
			),
		).WHERE(table.Symbol.Symbol.EQ(postgres.String(filter.Symbol.Get())))
	}
	dest := make([]model.BlacklistSymbol, 0)
	if err := stm.QueryContext(ctx, b.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []entity.BlacklistSymbol{}, nil
		}
		return nil, fmt.Errorf("BlackListSymbolRepository GetAll %w", err)
	}
	return MapBlacklistSymbolsDbToEntity(dest), nil
}

func (b *BlackListSymbolRepository) GetById(ctx context.Context, id int64) (entity.BlacklistSymbol, error) {
	stm := table.BlacklistSymbol.SELECT(table.BlacklistSymbol.AllColumns).WHERE(table.BlacklistSymbol.ID.EQ(postgres.Int64(id)))
	dest := model.BlacklistSymbol{}
	if err := stm.QueryContext(ctx, b.getDbFunc(ctx), &dest); err != nil {
		return entity.BlacklistSymbol{}, fmt.Errorf("BlackListSymbolRepository GetById %w", err)
	}
	return MapBlacklistSymbolDbToEntity(dest), nil
}

func (b *BlackListSymbolRepository) Create(ctx context.Context, symbol entity.BlacklistSymbol) (entity.BlacklistSymbol, error) {
	createModel := MapBlacklistSymbolEntityToDb(symbol)
	created := model.BlacklistSymbol{}
	if err := table.BlacklistSymbol.INSERT(table.BlacklistSymbol.MutableColumns).
		MODEL(createModel).
		RETURNING(table.BlacklistSymbol.AllColumns).
		QueryContext(ctx, b.getDbFunc(ctx), &created); err != nil {
		return entity.BlacklistSymbol{}, fmt.Errorf("BlackListSymbolRepository Create %w", err)
	}
	return MapBlacklistSymbolDbToEntity(created), nil
}

func (b *BlackListSymbolRepository) Update(ctx context.Context, symbol entity.BlacklistSymbol) (entity.BlacklistSymbol, error) {
	updateModel := MapBlacklistSymbolEntityToDb(symbol)
	updated := model.BlacklistSymbol{}
	if err := table.BlacklistSymbol.UPDATE(table.BlacklistSymbol.MutableColumns).
		MODEL(updateModel).
		WHERE(table.BlacklistSymbol.ID.EQ(postgres.Int64(symbol.Id))).
		RETURNING(table.BlacklistSymbol.AllColumns).
		QueryContext(ctx, b.getDbFunc(ctx), &updated); err != nil {
		return entity.BlacklistSymbol{}, fmt.Errorf("BlackListSymbolRepository Update %w", err)
	}
	return MapBlacklistSymbolDbToEntity(updated), nil
}

func (b *BlackListSymbolRepository) GetByAffectTime(ctx context.Context, symbolId int64, affectedFrom time.Time, affectedTo time.Time) ([]entity.BlacklistSymbol, error) {
	stm := postgres.SelectStatement(nil)
	if affectedTo.IsZero() {
		stm = table.BlacklistSymbol.SELECT(table.BlacklistSymbol.AllColumns).
			WHERE(table.BlacklistSymbol.SymbolID.EQ(postgres.Int64(symbolId)).
				AND(table.BlacklistSymbol.Status.EQ(enum.Blacklistsymbolstatus.Active)).
				AND(table.BlacklistSymbol.AffectedTo.IS_NULL().
					OR(table.BlacklistSymbol.AffectedTo.GT_EQ(postgres.TimestampT(affectedFrom))),
				))
	} else {
		stm = table.BlacklistSymbol.SELECT(table.BlacklistSymbol.AllColumns).
			WHERE(table.BlacklistSymbol.SymbolID.EQ(postgres.Int64(symbolId)).
				AND(table.BlacklistSymbol.Status.EQ(enum.Blacklistsymbolstatus.Active)).
				AND(table.BlacklistSymbol.AffectedFrom.LT_EQ(postgres.TimestampT(affectedTo))).
				AND(table.BlacklistSymbol.AffectedTo.IS_NULL().
					OR(table.BlacklistSymbol.AffectedTo.GT_EQ(postgres.TimestampT(affectedFrom))),
				))
	}

	dest := make([]model.BlacklistSymbol, 0)
	if err := stm.QueryContext(ctx, b.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []entity.BlacklistSymbol{}, nil
		}
		return nil, fmt.Errorf("BlackListSymbolRepository GetByAffectTime %w", err)
	}
	return MapBlacklistSymbolsDbToEntity(dest), nil
}
