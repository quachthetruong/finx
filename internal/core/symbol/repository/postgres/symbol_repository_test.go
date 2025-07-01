package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/optional"
)

func TestSymbolRepository_Create(t *testing.T) {
	t.Parallel()
	e := entity.Symbol{
		StockExchangeId: 1,
		Symbol:          "BER",
	}
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewSymbolRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)

	t.Run(
		"create success", func(t *testing.T) {
			mock.ExpectQuery("INSERT INTO public.symbol").WillReturnRows(
				sqlmock.NewRows(
					[]string{
						"symbol.id",
						"symbol.stock_exchange_id",
						"symbol.symbol",
						"symbol.created_at",
					},
				).AddRow(1, e.StockExchangeId, e.Symbol, e.CreatedAt),
			)
			created, err := repo.Create(context.Background(), e)
			assert.Nil(t, err)
			assert.Equal(t, e.StockExchangeId, created.StockExchangeId)
			assert.Equal(t, e.Symbol, created.Symbol)
		},
	)
	t.Run(
		"create error", func(t *testing.T) {
			mock.ExpectQuery("INSERT INTO public.symbol").WillReturnError(
				assert.AnError,
			)
			_, err := repo.Create(context.Background(), e)
			assert.NotNil(t, err)
		},
	)
}

func TestSymbolRepository_GetById(t *testing.T) {
	t.Parallel()
	e := entity.Symbol{
		StockExchangeId: 2,
		Symbol:          "DIO",
	}
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewSymbolRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)

	t.Run(
		"get by id success", func(t *testing.T) {
			mock.ExpectQuery("SELECT symbol.id").WillReturnRows(
				sqlmock.NewRows(
					[]string{
						"symbol.id",
						"symbol.stock_exchange_id",
						"symbol.symbol",
						"symbol.created_at",
					},
				).AddRow(1, e.StockExchangeId, e.Symbol, e.CreatedAt),
			)
			res, err := repo.GetById(context.Background(), 1)
			assert.Nil(t, err)
			assert.Equal(t, e.StockExchangeId, res.StockExchangeId)
			assert.Equal(t, e.Symbol, res.Symbol)
		},
	)
	t.Run(
		"get by id error", func(t *testing.T) {
			mock.ExpectQuery("SELECT symbol.id").WillReturnError(
				assert.AnError,
			)
			_, err := repo.GetById(context.Background(), 1)
			assert.NotNil(t, err)
		},
	)
}

func TestSymbolRepository_GetAll(t *testing.T) {
	t.Parallel()
	ee := []entity.Symbol{
		{
			StockExchangeId: 1,
			Symbol:          "DIO",
		},
		{
			StockExchangeId: 2,
			Symbol:          "IJJ",
		},
		{
			StockExchangeId: 2,
			Symbol:          "IUD",
		},
	}
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewSymbolRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)

	t.Run(
		"get all success", func(t *testing.T) {
			rows := sqlmock.NewRows(
				[]string{
					"symbol.id",
					"symbol.stock_exchange_id",
					"symbol.symbol",
					"symbol.created_at",
				},
			)
			for i, e := range ee {
				rows.AddRow(i+1, e.StockExchangeId, e.Symbol, e.CreatedAt)
			}
			mock.ExpectQuery("SELECT symbol.id").WillReturnRows(rows)
			res, err := repo.GetAll(context.Background(), entity.SymbolFilter{})
			assert.Nil(t, err)
			assert.Equal(t, len(ee), len(res))
			assert.Equal(t, ee[0].StockExchangeId, res[0].StockExchangeId)
			assert.Equal(t, ee[2].Symbol, res[2].Symbol)
		},
	)
	t.Run(
		"get all success with filter", func(t *testing.T) {
			rows := sqlmock.NewRows(
				[]string{
					"symbol.id",
					"symbol.stock_exchange_id",
					"symbol.symbol",
					"symbol.created_at",
				},
			)
			for i, e := range ee {
				if i == 0 {
					continue
				}
				rows.AddRow(i+1, e.StockExchangeId, e.Symbol, e.CreatedAt)
			}
			pageSize := uint(20)
			pageNumber := uint(34)
			mock.ExpectQuery("SELECT symbol.id").WithArgs(
				true, "SHA", pageSize, pageSize*(pageNumber-1),
			).WillReturnRows(rows)
			res, err := repo.GetAll(
				context.Background(), entity.SymbolFilter{
					Paging: core.Paging{
						Size:   pageSize,
						Number: pageNumber,
					},
					StockExchangeCode: optional.Some("SHA"),
				},
			)
			assert.Nil(t, err)
			assert.Equal(t, 2, len(res))
			assert.Equal(t, ee[1].StockExchangeId, res[0].StockExchangeId)
			assert.Equal(t, ee[2].Symbol, res[1].Symbol)
		},
	)
	t.Run(
		"get all error", func(t *testing.T) {
			mock.ExpectQuery("SELECT symbol.id").WillReturnError(
				qrm.ErrNoRows,
			)
			_, err := repo.GetAll(context.Background(), entity.SymbolFilter{})
			assert.Nil(t, err)
		},
	)
}

func TestSymbolRepository_Update(t *testing.T) {
	t.Parallel()
	e := entity.Symbol{
		StockExchangeId: 3,
		Symbol:          "BER",
	}
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewSymbolRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)

	t.Run(
		"update success", func(t *testing.T) {
			mock.ExpectQuery("UPDATE public.symbol").WillReturnRows(
				sqlmock.NewRows(
					[]string{
						"symbol.id",
						"symbol.stock_exchange_id",
						"symbol.symbol",
						"symbol.created_at",
					},
				).AddRow(1, e.StockExchangeId, e.Symbol, e.CreatedAt),
			)
			updated, err := repo.Update(context.Background(), e)
			assert.Nil(t, err)
			assert.Equal(t, e.StockExchangeId, updated.StockExchangeId)
			assert.Equal(t, e.Symbol, updated.Symbol)
		},
	)
	t.Run(
		"update error", func(t *testing.T) {
			mock.ExpectQuery("UPDATE public.symbol").WillReturnError(
				assert.AnError,
			)
			_, err := repo.Update(context.Background(), e)
			assert.NotNil(t, err)
		},
	)
}

func TestSymbolRepository_GetSymbolWithActiveBlacklist(t *testing.T) {
	t.Parallel()
	symbol := entity.Symbol{
		Id:              12,
		StockExchangeId: 1,
		Symbol:          "FPT",
		AssetType:       "UNDERLYING",
		Status:          "ACTIVE",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Scores:          nil,
	}
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewSymbolRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	t.Run("get symbol with active blacklist success", func(t *testing.T) {
		mock.ExpectQuery("SELECT symbol.*").
			WillReturnRows(
				sqlmock.NewRows(
					[]string{
						"symbol.id",
						"symbol.stock_exchange_id",
						"symbol.symbol",
						"symbol.asset_type",
						"symbol.status",
						"symbol.created_at",
						"symbol.updated_at",
					},
				).AddRow(symbol.Id, symbol.StockExchangeId, symbol.Symbol, symbol.AssetType, symbol.Status, symbol.CreatedAt, symbol.UpdatedAt),
			)
		res, err := repo.GetSymbolWithActiveBlacklist(context.Background(), symbol.Symbol)
		assert.Nil(t, err)
		assert.Equal(t, symbol, res)
	})
	t.Run("get symbol with active blacklist error", func(t *testing.T) {
		mock.ExpectQuery("SELECT symbol.*").
			WillReturnError(assert.AnError)
		res, err := repo.GetSymbolWithActiveBlacklist(context.Background(), symbol.Symbol)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, res)
	})
}

func TestSymbolRepository_GetBySymbol(t *testing.T) {
	t.Parallel()
	symbol := entity.Symbol{
		Id:              12,
		StockExchangeId: 1,
		Symbol:          "FPT",
		AssetType:       "UNDERLYING",
		Status:          "ACTIVE",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Scores:          nil,
	}
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewSymbolRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	t.Run("get symbol by symbol success", func(t *testing.T) {
		mock.ExpectQuery("SELECT symbol.*").WithArgs(symbol.Symbol).WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"symbol.id",
					"symbol.stock_exchange_id",
					"symbol.symbol",
					"symbol.asset_type",
					"symbol.status",
					"symbol.created_at",
					"symbol.updated_at",
				},
			).AddRow(symbol.Id, symbol.StockExchangeId, symbol.Symbol, symbol.AssetType, symbol.Status, symbol.CreatedAt, symbol.UpdatedAt),
		)
		res, err := repo.GetBySymbol(context.Background(), symbol.Symbol)
		assert.Nil(t, err)
		assert.Equal(t, symbol, res)
	})
	t.Run("get symbol by symbol error", func(t *testing.T) {
		mock.ExpectQuery("SELECT symbol.*").WithArgs(symbol.Symbol).WillReturnError(assert.AnError)
		res, err := repo.GetBySymbol(context.Background(), symbol.Symbol)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, res)
	})
}
