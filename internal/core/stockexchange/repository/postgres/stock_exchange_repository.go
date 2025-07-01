package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/stockexchange/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

var _ repository.StockExchangeRepository = (*StockExchangeRepository)(nil)

type StockExchangeRepository struct {
	getDbFunc database.GetDbFunc
}

func (s *StockExchangeRepository) GetById(ctx context.Context, id int64) (entity.StockExchange, error) {
	var res model.StockExchange
	if err := table.StockExchange.SELECT(table.StockExchange.AllColumns).
		WHERE(table.StockExchange.ID.EQ(postgres.Int64(id))).
		QueryContext(ctx, s.getDbFunc(ctx), &res); err != nil {
		return entity.StockExchange{}, fmt.Errorf("StockExchangeRepository GetById %w", err)
	}
	return MapStockExchangeDbToEntity(res), nil
}

func (s *StockExchangeRepository) GetBySymbolId(ctx context.Context, symbolId int64) (entity.StockExchange, error) {
	var res model.StockExchange
	if err := table.StockExchange.SELECT(table.StockExchange.AllColumns).
		FROM(table.StockExchange.INNER_JOIN(table.Symbol, table.StockExchange.ID.EQ(table.Symbol.StockExchangeID))).
		WHERE(table.Symbol.ID.EQ(postgres.Int64(symbolId))).
		QueryContext(ctx, s.getDbFunc(ctx), &res); err != nil {
		return entity.StockExchange{}, fmt.Errorf("StockExchangeRepository GetBySymbolId %w", err)
	}
	return MapStockExchangeDbToEntity(res), nil
}

func (s *StockExchangeRepository) GetAll(ctx context.Context) ([]entity.StockExchange, error) {
	stm := table.StockExchange.SELECT(table.StockExchange.AllColumns)
	dest := make([]model.StockExchange, 0)
	if err := stm.QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []entity.StockExchange{}, nil
		}
		return nil, fmt.Errorf("StockExchangeRepository GetAll %w", err)
	}
	return MapStockExchangesDbToEntity(dest), nil
}

func (s *StockExchangeRepository) Create(ctx context.Context, stockExchange entity.StockExchange) (entity.StockExchange, error) {
	createModel := MapStockExchangeEntityToDb(stockExchange)
	created := model.StockExchange{}
	if err := table.StockExchange.INSERT(table.StockExchange.MutableColumns).
		MODEL(createModel).
		RETURNING(table.StockExchange.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &created); err != nil {
		return entity.StockExchange{}, fmt.Errorf("StockExchangeRepository Create %w", err)
	}
	return MapStockExchangeDbToEntity(created), nil
}

func (s *StockExchangeRepository) Update(ctx context.Context, stockExchange entity.StockExchange) (entity.StockExchange, error) {
	updateModel := MapStockExchangeEntityToDb(stockExchange)
	updated := model.StockExchange{}
	if err := table.StockExchange.UPDATE(table.StockExchange.MutableColumns).
		MODEL(updateModel).
		WHERE(table.StockExchange.ID.EQ(postgres.Int64(updateModel.ID))).
		RETURNING(table.StockExchange.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &updated); err != nil {
		return entity.StockExchange{}, fmt.Errorf("StockExchangeRepository Update %w", err)
	}
	return MapStockExchangeDbToEntity(updated), nil
}

func (s *StockExchangeRepository) Delete(ctx context.Context, id int64) error {
	if _, err := table.StockExchange.DELETE().WHERE(table.StockExchange.ID.EQ(postgres.Int64(id))).ExecContext(
		ctx, s.getDbFunc(ctx),
	); err != nil {
		return fmt.Errorf("StockExchangeRepository Delete %w", err)
	}
	return nil
}

func NewStockExchangeRepository(getDbFunc database.GetDbFunc) *StockExchangeRepository {
	return &StockExchangeRepository{getDbFunc: getDbFunc}
}
