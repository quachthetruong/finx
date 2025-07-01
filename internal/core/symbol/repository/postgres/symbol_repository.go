package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/enum"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

type SymbolRepository struct {
	getDbFunc database.GetDbFunc
}

func NewSymbolRepository(getDbFunc database.GetDbFunc) *SymbolRepository {
	return &SymbolRepository{getDbFunc: getDbFunc}
}

func (s *SymbolRepository) GetAll(ctx context.Context, filter entity.SymbolFilter) ([]entity.Symbol, error) {
	stm := table.Symbol.SELECT(table.Symbol.AllColumns)
	if filter.StockExchangeCode.IsPresent() {
		stm = stm.FROM(
			table.Symbol.INNER_JOIN(
				table.StockExchange, table.StockExchange.ID.EQ(table.Symbol.StockExchangeID),
			),
		)
	}
	stm = stm.WHERE(ApplyFilter(filter))
	if limit := filter.Limit(); limit > 0 {
		stm = stm.LIMIT(limit).OFFSET(filter.Offset())
	}
	if orderClause := ApplySort(filter); len(orderClause) > 0 {
		stm = stm.ORDER_BY(orderClause...)
	}
	dest := make([]model.Symbol, 0)
	if err := stm.QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []entity.Symbol{}, nil
		}
		return []entity.Symbol{}, fmt.Errorf("SymbolRepository GetAll %w", err)
	}
	return MapSymbolsDbToEntity(dest), nil
}

func (s *SymbolRepository) Count(ctx context.Context, filter entity.SymbolFilter) (int64, error) {
	stm := table.Symbol.SELECT(postgres.COUNT(table.Symbol.ID))
	if filter.StockExchangeCode.IsPresent() {
		stm = stm.FROM(
			table.Symbol.INNER_JOIN(
				table.StockExchange, table.StockExchange.ID.EQ(table.Symbol.StockExchangeID),
			),
		)
	}
	stm = stm.WHERE(ApplyFilter(filter))
	dest := struct {
		Count int64
	}{}
	if err := stm.QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("SymbolRepository Count %w", err)
	}
	return dest.Count, nil
}

func (s *SymbolRepository) GetById(ctx context.Context, id int64) (entity.Symbol, error) {
	var res model.Symbol
	if err := table.Symbol.SELECT(table.Symbol.AllColumns).WHERE(table.Symbol.ID.EQ(postgres.Int64(id))).QueryContext(
		ctx, s.getDbFunc(ctx), &res,
	); err != nil {
		return entity.Symbol{}, fmt.Errorf("SymbolRepository GetById %w", err)
	}
	return MapSymbolDbToEntity(res), nil
}

func (s *SymbolRepository) Create(ctx context.Context, symbol entity.Symbol) (entity.Symbol, error) {
	createModel := MapSymbolEntityToDb(symbol)
	created := model.Symbol{}
	if err := table.Symbol.INSERT(table.Symbol.MutableColumns).
		MODEL(createModel).
		RETURNING(table.Symbol.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &created); err != nil {
		return entity.Symbol{}, fmt.Errorf("SymbolRepository Create %w", err)
	}
	return MapSymbolDbToEntity(created), nil
}

func (s *SymbolRepository) Update(ctx context.Context, symbol entity.Symbol) (entity.Symbol, error) {
	updateModel := MapSymbolEntityToDb(symbol)
	updated := model.Symbol{}
	if err := table.Symbol.UPDATE(table.Symbol.MutableColumns).
		MODEL(updateModel).
		WHERE(table.Symbol.ID.EQ(postgres.Int64(updateModel.ID))).
		RETURNING(table.Symbol.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &updated); err != nil {
		return entity.Symbol{}, fmt.Errorf("SymbolRepository Update %w", err)
	}
	return MapSymbolDbToEntity(updated), nil
}

func (s *SymbolRepository) GetBySymbol(ctx context.Context, symbol string) (entity.Symbol, error) {
	var res model.Symbol
	if err := table.Symbol.
		SELECT(table.Symbol.AllColumns).
		WHERE(table.Symbol.Symbol.EQ(postgres.String(symbol))).
		QueryContext(
			ctx, s.getDbFunc(ctx), &res,
		); err != nil {
		return entity.Symbol{}, fmt.Errorf("SymbolRepository GetBySymbol %w", err)
	}
	return MapSymbolDbToEntity(res), nil
}

func (s *SymbolRepository) GetSymbolWithActiveBlacklist(ctx context.Context, symbol string) (entity.Symbol, error) {
	var res model.Symbol
	if err := table.Symbol.SELECT(table.Symbol.AllColumns).FROM(
		table.Symbol.LEFT_JOIN(
			table.BlacklistSymbol, table.Symbol.ID.EQ(table.BlacklistSymbol.SymbolID),
		),
	).WHERE(
		table.BlacklistSymbol.Status.EQ(enum.Blacklistsymbolstatus.Active).
			AND(table.BlacklistSymbol.AffectedFrom.LT_EQ(postgres.TimestampExp(postgres.NOW()))).
			AND(
				table.BlacklistSymbol.AffectedTo.IS_NULL().
					OR(table.BlacklistSymbol.AffectedTo.GT_EQ(postgres.TimestampExp(postgres.NOW()))),
			).AND(table.Symbol.Symbol.EQ(postgres.String(symbol))),
	).
		QueryContext(ctx, s.getDbFunc(ctx), &res); err != nil {
		return entity.Symbol{}, fmt.Errorf("SymbolRepository GetSymbolWithActiveBlacklist %w", err)
	}

	return MapSymbolDbToEntity(res), nil
}
