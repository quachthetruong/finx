package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/optional"
)

func TestBlackListSymbolRepository_GetByAffectTime(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Error(err)
	}
	repo := NewBlackListSymbolRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	blacklistSymbols := []entity.BlacklistSymbol{
		{
			Id:           123,
			SymbolId:     2,
			AffectedFrom: time.Date(2001, 2, 3, 4, 5, 6, 7, time.UTC),
			AffectedTo:   time.Date(2011, 2, 3, 4, 5, 6, 7, time.UTC),
			Status:       entity.BlacklistSymbolStatusInactive,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Id:           12,
			SymbolId:     2,
			AffectedFrom: time.Date(2022, 2, 3, 4, 5, 6, 7, time.UTC),
			AffectedTo:   time.Date(2024, 2, 3, 4, 5, 6, 7, time.UTC),
			Status:       entity.BlacklistSymbolStatusActive,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}
	t.Run("get by affect time success", func(t *testing.T) {
		rows := sqlmock.NewRows(
			[]string{
				"blacklist_symbol.id",
				"blacklist_symbol.symbol_id",
				"blacklist_symbol.affected_from",
				"blacklist_symbol.affected_to",
				"blacklist_symbol.status",
				"blacklist_symbol.created_at",
				"blacklist_symbol.updated_at",
			})
		for _, v := range blacklistSymbols {
			rows.AddRow(
				v.Id,
				v.SymbolId,
				v.AffectedFrom,
				v.AffectedTo,
				v.Status,
				v.CreatedAt,
				v.UpdatedAt,
			)
		}
		mock.ExpectQuery("SELECT .* FROM public.blacklist_symbol").
			WithArgs(
				blacklistSymbols[0].SymbolId,
				time.Date(2025, 2, 3, 4, 5, 6, 7, time.UTC),
				time.Date(2000, 2, 3, 4, 5, 6, 7, time.UTC),
			).
			WillReturnRows(rows)
		res, err := repo.GetByAffectTime(context.Background(), blacklistSymbols[0].SymbolId,
			time.Date(2000, 2, 3, 4, 5, 6, 7, time.UTC),
			time.Date(2025, 2, 3, 4, 5, 6, 7, time.UTC),
		)
		assert.Nil(t, err)
		assert.Equal(t, len(blacklistSymbols), len(res))
		assert.Equal(t, blacklistSymbols, res)
	})

	t.Run("get by affect time error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM public.blacklist_symbol").
			WithArgs(
				blacklistSymbols[0].SymbolId,
				time.Date(2025, 2, 3, 4, 5, 6, 7, time.UTC),
				time.Date(2000, 2, 3, 4, 5, 6, 7, time.UTC),
			).WillReturnError(assert.AnError)
		res, err := repo.GetByAffectTime(context.Background(), blacklistSymbols[0].SymbolId,
			time.Date(2000, 2, 3, 4, 5, 6, 7, time.UTC),
			time.Date(2025, 2, 3, 4, 5, 6, 7, time.UTC),
		)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, res)
	})

	t.Run("get by affect time no rows", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM public.blacklist_symbol").
			WithArgs(
				blacklistSymbols[0].SymbolId,
				time.Date(2025, 2, 3, 4, 5, 6, 7, time.UTC),
				time.Date(2000, 2, 3, 4, 5, 6, 7, time.UTC),
			).WillReturnError(qrm.ErrNoRows)
		res, err := repo.GetByAffectTime(context.Background(), blacklistSymbols[0].SymbolId,
			time.Date(2000, 2, 3, 4, 5, 6, 7, time.UTC),
			time.Date(2025, 2, 3, 4, 5, 6, 7, time.UTC),
		)
		assert.Nil(t, err)
		assert.Empty(t, res)
	})

	t.Run("get by affect time success with null", func(t *testing.T) {
		rows := sqlmock.NewRows(
			[]string{
				"blacklist_symbol.id",
				"blacklist_symbol.symbol_id",
				"blacklist_symbol.affected_from",
				"blacklist_symbol.affected_to",
				"blacklist_symbol.status",
				"blacklist_symbol.created_at",
				"blacklist_symbol.updated_at",
			})
		for _, v := range blacklistSymbols {
			rows.AddRow(
				v.Id,
				v.SymbolId,
				v.AffectedFrom,
				v.AffectedTo,
				v.Status,
				v.CreatedAt,
				v.UpdatedAt,
			)
		}
		mock.ExpectQuery("SELECT .* FROM public.blacklist_symbol").
			WithArgs(
				blacklistSymbols[0].SymbolId,
				time.Date(2000, 2, 3, 4, 5, 6, 7, time.UTC),
			).
			WillReturnRows(rows)
		res, err := repo.GetByAffectTime(context.Background(), blacklistSymbols[0].SymbolId,
			time.Date(2000, 2, 3, 4, 5, 6, 7, time.UTC),
			time.Time{},
		)
		assert.Nil(t, err)
		assert.Equal(t, len(blacklistSymbols), len(res))
		assert.Equal(t, blacklistSymbols, res)
	})

	t.Run("get by affect time error with null", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM public.blacklist_symbol").
			WithArgs(
				blacklistSymbols[0].SymbolId,
				time.Date(2000, 2, 3, 4, 5, 6, 7, time.UTC),
			).WillReturnError(assert.AnError)
		res, err := repo.GetByAffectTime(context.Background(), blacklistSymbols[0].SymbolId,
			time.Date(2000, 2, 3, 4, 5, 6, 7, time.UTC),
			time.Time{},
		)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, res)
	})

	t.Run("get by affect time no rows with null", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM public.blacklist_symbol").
			WithArgs(
				blacklistSymbols[0].SymbolId,
				time.Date(3000, 2, 3, 4, 5, 6, 7, time.UTC),
			).WillReturnError(qrm.ErrNoRows)
		res, err := repo.GetByAffectTime(context.Background(), blacklistSymbols[0].SymbolId,
			time.Date(3000, 2, 3, 4, 5, 6, 7, time.UTC),
			time.Time{},
		)
		assert.Nil(t, err)
		assert.Empty(t, res)
	})
}

func TestBlackListSymbolRepository_GetAll(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Error(err)
	}
	repo := NewBlackListSymbolRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	blacklistSymbols := []entity.BlacklistSymbol{
		{
			Id:           123,
			SymbolId:     2,
			AffectedFrom: time.Date(2001, 2, 3, 4, 5, 6, 7, time.UTC),
			AffectedTo:   time.Date(2011, 2, 3, 4, 5, 6, 7, time.UTC),
			Status:       entity.BlacklistSymbolStatusInactive,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Id:           12,
			SymbolId:     2,
			AffectedFrom: time.Date(2022, 2, 3, 4, 5, 6, 7, time.UTC),
			AffectedTo:   time.Date(2024, 2, 3, 4, 5, 6, 7, time.UTC),
			Status:       entity.BlacklistSymbolStatusActive,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	t.Run("get all success", func(t *testing.T) {
		rows := sqlmock.NewRows(
			[]string{
				"blacklist_symbol.id",
				"blacklist_symbol.symbol_id",
				"blacklist_symbol.affected_from",
				"blacklist_symbol.affected_to",
				"blacklist_symbol.status",
				"blacklist_symbol.created_at",
				"blacklist_symbol.updated_at",
			})
		for _, v := range blacklistSymbols {
			rows.AddRow(
				v.Id,
				v.SymbolId,
				v.AffectedFrom,
				v.AffectedTo,
				v.Status,
				v.CreatedAt,
				v.UpdatedAt,
			)
		}
		mock.ExpectQuery("SELECT .* FROM public.blacklist_symbol").WillReturnRows(rows)
		res, err := repo.GetAll(context.Background(), entity.BlacklistSymbolFilter{})
		assert.Nil(t, err)
		assert.Equal(t, len(blacklistSymbols), len(res))
	})

	t.Run("get all success with filter", func(t *testing.T) {
		rows := sqlmock.NewRows(
			[]string{
				"blacklist_symbol.id",
				"blacklist_symbol.symbol_id",
				"blacklist_symbol.affected_from",
				"blacklist_symbol.affected_to",
				"blacklist_symbol.status",
				"blacklist_symbol.created_at",
				"blacklist_symbol.updated_at",
			})
		for _, v := range blacklistSymbols {
			rows.AddRow(
				v.Id,
				v.SymbolId,
				v.AffectedFrom,
				v.AffectedTo,
				v.Status,
				v.CreatedAt,
				v.UpdatedAt,
			)
		}
		mock.ExpectQuery("SELECT .* FROM public.blacklist_symbol").WillReturnRows(rows)
		res, err := repo.GetAll(context.Background(), entity.BlacklistSymbolFilter{
			Symbol: optional.Some("BTC"),
		})
		assert.Nil(t, err)
		assert.Equal(t, len(blacklistSymbols), len(res))
	})

	t.Run("get all error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM public.blacklist_symbol").WillReturnError(assert.AnError)
		res, err := repo.GetAll(context.Background(), entity.BlacklistSymbolFilter{})
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, res)
	})

	t.Run("get all no rows", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM public.blacklist_symbol").WillReturnError(qrm.ErrNoRows)
		res, err := repo.GetAll(context.Background(), entity.BlacklistSymbolFilter{
			Symbol: optional.Some("BTC"),
		})
		assert.Nil(t, err)
		assert.Empty(t, res)
	})
}

func TestBlackListSymbolRepository_GetById(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Error(err)
	}
	repo := NewBlackListSymbolRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	blacklistSymbol := entity.BlacklistSymbol{
		Id:           123,
		SymbolId:     2,
		AffectedFrom: time.Date(2001, 2, 3, 4, 5, 6, 7, time.UTC),
		AffectedTo:   time.Date(2011, 2, 3, 4, 5, 6, 7, time.UTC),
		Status:       entity.BlacklistSymbolStatusInactive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	// Test case
	t.Run("get by id success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM public.blacklist_symbol").WithArgs(blacklistSymbol.Id).WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"blacklist_symbol.id",
					"blacklist_symbol.symbol_id",
					"blacklist_symbol.affected_from",
					"blacklist_symbol.affected_to",
					"blacklist_symbol.status",
					"blacklist_symbol.created_at",
					"blacklist_symbol.updated_at",
				}).AddRow(blacklistSymbol.Id, blacklistSymbol.SymbolId, blacklistSymbol.AffectedFrom,
				blacklistSymbol.AffectedTo, blacklistSymbol.Status, blacklistSymbol.CreatedAt,
				blacklistSymbol.UpdatedAt),
		)
		res, err := repo.GetById(context.Background(), blacklistSymbol.Id)
		assert.Nil(t, err)
		assert.Equal(t, blacklistSymbol, res)
	})
	t.Run("get by id error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM public.blacklist_symbol").WithArgs(blacklistSymbol.Id).
			WillReturnError(assert.AnError)
		res, err := repo.GetById(context.Background(), blacklistSymbol.Id)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, res)
	})
}

func TestBlackListSymbolRepository_Create(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Error(err)
	}
	repo := NewBlackListSymbolRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	blacklistSymbol := entity.BlacklistSymbol{
		Id:           12,
		SymbolId:     3,
		AffectedFrom: time.Date(2001, 2, 3, 4, 5, 6, 7, time.UTC),
		AffectedTo:   time.Date(2011, 2, 3, 4, 5, 6, 7, time.UTC),
		Status:       entity.BlacklistSymbolStatusInactive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	// Test case
	t.Run("create BlacklistSymbol success", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO public.blacklist_symbol").WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"blacklist_symbol.id",
					"blacklist_symbol.symbol_id",
					"blacklist_symbol.affected_from",
					"blacklist_symbol.affected_to",
					"blacklist_symbol.status",
					"blacklist_symbol.created_at",
					"blacklist_symbol.updated_at",
				}).AddRow(blacklistSymbol.Id, blacklistSymbol.SymbolId, blacklistSymbol.AffectedFrom,
				blacklistSymbol.AffectedTo, blacklistSymbol.Status, blacklistSymbol.CreatedAt,
				blacklistSymbol.UpdatedAt),
		)
		res, err := repo.Create(context.Background(), blacklistSymbol)
		assert.Nil(t, err)
		assert.Equal(t, blacklistSymbol, res)
	})
	t.Run("create BlacklistSymbol error", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO public.blacklist_symbol").WillReturnError(assert.AnError)
		res, err := repo.Create(context.Background(), blacklistSymbol)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, res)
	})
}

func TestBlackListSymbolRepository_Update(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Error(err)
	}
	repo := NewBlackListSymbolRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	blacklistSymbol := entity.BlacklistSymbol{
		Id:           12,
		SymbolId:     3,
		AffectedFrom: time.Date(2001, 2, 3, 4, 5, 6, 7, time.UTC),
		AffectedTo:   time.Date(2011, 2, 3, 4, 5, 6, 7, time.UTC),
		Status:       entity.BlacklistSymbolStatusInactive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	// Test case
	t.Run("update BlacklistSymbol success", func(t *testing.T) {
		mock.ExpectQuery("UPDATE public.blacklist_symbol").WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"blacklist_symbol.id",
					"blacklist_symbol.symbol_id",
					"blacklist_symbol.affected_from",
					"blacklist_symbol.affected_to",
					"blacklist_symbol.status",
					"blacklist_symbol.created_at",
					"blacklist_symbol.updated_at",
				}).AddRow(blacklistSymbol.Id, blacklistSymbol.SymbolId, blacklistSymbol.AffectedFrom,
				blacklistSymbol.AffectedTo, entity.BlacklistSymbolStatusInactive, blacklistSymbol.CreatedAt,
				blacklistSymbol.UpdatedAt),
		)
		res, err := repo.Update(context.Background(), blacklistSymbol)
		assert.Nil(t, err)
		assert.Equal(t, entity.BlacklistSymbolStatusInactive, res.Status)
	})
	t.Run("update BlacklistSymbol error", func(t *testing.T) {
		mock.ExpectQuery("UPDATE public.blacklist_symbol").WillReturnError(assert.AnError)
		res, err := repo.Update(context.Background(), blacklistSymbol)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, res)
	})
}
