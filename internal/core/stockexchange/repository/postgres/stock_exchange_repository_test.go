package postgres

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/pkg/dbtest"
)

func TestStockExchangeRepository_Create(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewStockExchangeRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	t.Run(
		"create success", func(t *testing.T) {
			e := entity.StockExchange{
				Id:           0,
				Code:         "CLONE",
				ScoreGroupId: int64(1),
				MinScore:     78,
				MaxScore:     100,
			}
			mock.ExpectQuery("INSERT INTO public.stock_exchange").WillReturnRows(
				sqlmock.NewRows(
					[]string{
						"stock_exchange.id",
						"stock_exchange.code",
						"stock_exchange.min_score",
						"stock_exchange.max_score",
						"stock_exchange.created_at",
						"stock_exchange.updated_at",
						"stock_exchange.score_group_id",
					},
				).AddRow(1, e.Code, e.MinScore, e.MaxScore, e.CreatedAt, e.UpdatedAt, e.ScoreGroupId),
			)
			created, err := repo.Create(context.Background(), e)
			assert.Nil(t, err)
			assert.Equal(t, e.Code, created.Code)
			assert.Equal(t, e.MinScore, created.MinScore)
			assert.Equal(t, e.ScoreGroupId, created.ScoreGroupId)
		},
	)
	t.Run(
		"create error", func(t *testing.T) {
			mock.ExpectQuery("INSERT INTO public.stock_exchange").WillReturnError(assert.AnError)
			_, err := repo.Create(context.Background(), entity.StockExchange{})
			assert.NotNil(t, err)
		},
	)
}

func TestStockExchangeRepository_GetAll(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewStockExchangeRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	t.Run(
		"get all success", func(t *testing.T) {
			e := entity.StockExchange{
				Id:           0,
				Code:         "CLONE",
				MinScore:     78,
				MaxScore:     100,
				ScoreGroupId: int64(1),
			}
			mock.ExpectQuery("SELECT stock_exchange.id").WillReturnRows(
				sqlmock.NewRows(
					[]string{
						"stock_exchange.id",
						"stock_exchange.code",
						"stock_exchange.min_score",
						"stock_exchange.max_score",
						"stock_exchange.created_at",
						"stock_exchange.updated_at",
						"stock_exchange.score_group_id",
					},
				).AddRow(1, e.Code, e.MinScore, e.MaxScore, e.CreatedAt, e.UpdatedAt, e.ScoreGroupId),
			)
			res, err := repo.GetAll(context.Background())
			assert.Nil(t, err)
			assert.Equal(t, e.Code, res[0].Code)
			assert.Equal(t, e.MinScore, res[0].MinScore)
			assert.Equal(t, e.ScoreGroupId, res[0].ScoreGroupId)
		},
	)
	t.Run(
		"get all error", func(t *testing.T) {
			mock.ExpectQuery("SELECT stock_exchange.id").WillReturnError(assert.AnError)
			_, err := repo.GetAll(context.Background())
			assert.NotNil(t, err)
		},
	)
}

func TestStockExchangeRepository_Update(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewStockExchangeRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	t.Run(
		"update success", func(t *testing.T) {
			e := entity.StockExchange{
				Id:           0,
				Code:         "CLONE",
				MinScore:     78,
				MaxScore:     100,
				ScoreGroupId: int64(1),
			}
			mock.ExpectQuery("UPDATE public.stock_exchange").WillReturnRows(
				sqlmock.NewRows(
					[]string{
						"stock_exchange.id",
						"stock_exchange.code",
						"stock_exchange.min_score",
						"stock_exchange.max_score",
						"stock_exchange.created_at",
						"stock_exchange.updated_at",
						"stock_exchange.score_group_id",
					},
				).AddRow(1, e.Code, e.MinScore, e.MaxScore, e.CreatedAt, e.UpdatedAt, e.ScoreGroupId),
			)
			updated, err := repo.Update(context.Background(), e)
			assert.Nil(t, err)
			assert.Equal(t, e.Code, updated.Code)
			assert.Equal(t, e.MinScore, updated.MinScore)
			assert.Equal(t, e.ScoreGroupId, updated.ScoreGroupId)
		},
	)
	t.Run(
		"update error", func(t *testing.T) {
			mock.ExpectQuery("UPDATE public.stock_exchange").WillReturnError(assert.AnError)
			_, err := repo.Update(context.Background(), entity.StockExchange{})
			assert.NotNil(t, err)
		},
	)
}

func TestStockExchangeRepository_Delete(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewStockExchangeRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	t.Run(
		"delete success", func(t *testing.T) {
			mock.ExpectExec("DELETE FROM public.stock_exchange").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
			err := repo.Delete(context.Background(), 1)
			assert.Nil(t, err)
		},
	)
	t.Run(
		"delete error", func(t *testing.T) {
			mock.ExpectExec("DELETE FROM public.stock_exchange").WillReturnError(assert.AnError)
			err := repo.Delete(context.Background(), 1)
			assert.NotNil(t, err)
		},
	)
}
