package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/optional"
)

func TestSuggestedOfferConfigRepository_GetAll(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewSuggestedOfferConfigRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	listEntity := []entity.SuggestedOfferConfig{
		{
			Id:            1,
			Name:          "Test 1",
			Value:         decimal.NewFromFloat(6.99),
			ValueType:     entity.ValueTypeInterestRate,
			Status:        "ACTIVE",
			CreatedBy:     "admin@gmail.com",
			LastUpdatedBy: "",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Id:            2,
			Value:         decimal.NewFromFloat(5.99),
			ValueType:     entity.ValueTypeInterestRate,
			Name:          "Test 2",
			Status:        "ACTIVE",
			CreatedBy:     "admin@gmail.com",
			LastUpdatedBy: "",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	// Test case
	t.Run("get all success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"suggested_offer_config.id",
			"suggested_offer_config.value_type",
			"suggested_offer_config.name",
			"suggested_offer_config.value",
			"suggested_offer_config.status",
			"suggested_offer_config.created_by",
			"suggested_offer_config.last_updated_by",
			"suggested_offer_config.created_at",
			"suggested_offer_config.updated_at",
		})
		for _, data := range listEntity {
			rows.AddRow(
				data.Id,
				data.ValueType,
				data.Name,
				data.Value,
				data.Status,
				data.CreatedBy,
				data.LastUpdatedBy,
				data.CreatedAt,
				data.UpdatedAt,
			)
		}
		mock.ExpectQuery("SELECT .* FROM public.suggested_offer_config").WillReturnRows(rows)
		res, err := repo.GetAll(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, listEntity, res)
	})

	t.Run("get all error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM public.suggested_offer_config").WillReturnError(assert.AnError)
		res, err := repo.GetAll(context.Background())
		assert.Nil(t, res)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestSuggestedOfferConfigRepository_GetActiveSuggestedOfferConfig(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewSuggestedOfferConfigRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	suggestedOfferConfig := entity.SuggestedOfferConfig{
		Id:            1,
		Name:          "Test 1",
		Value:         decimal.NewFromFloat(6.99),
		ValueType:     entity.ValueTypeInterestRate,
		Status:        "ACTIVE",
		CreatedBy:     "",
		LastUpdatedBy: "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	// Test case
	t.Run("get active suggested offer config success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"suggested_offer_config.id",
			"suggested_offer_config.value_type",
			"suggested_offer_config.name",
			"suggested_offer_config.value",
			"suggested_offer_config.status",
			"suggested_offer_config.created_by",
			"suggested_offer_config.last_updated_by",
			"suggested_offer_config.created_at",
			"suggested_offer_config.updated_at",
		})
		rows.AddRow(
			suggestedOfferConfig.Id,
			suggestedOfferConfig.ValueType,
			suggestedOfferConfig.Name,
			suggestedOfferConfig.Value,
			suggestedOfferConfig.Status,
			suggestedOfferConfig.CreatedBy,
			suggestedOfferConfig.LastUpdatedBy,
			suggestedOfferConfig.CreatedAt,
			suggestedOfferConfig.UpdatedAt,
		)
		mock.ExpectQuery("SELECT .* FROM public.suggested_offer_config").WillReturnRows(rows)
		res, err := repo.GetActiveSuggestedOfferConfig(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, optional.Some(suggestedOfferConfig), res)
	})
	t.Run("get active suggested offer config error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM public.suggested_offer_config").WillReturnError(assert.AnError)
		res, err := repo.GetActiveSuggestedOfferConfig(context.Background())
		assert.Empty(t, res)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("get active suggested offer config no rows", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM public.suggested_offer_config").WillReturnError(qrm.ErrNoRows)
		res, err := repo.GetActiveSuggestedOfferConfig(context.Background())
		assert.Equal(t, optional.None[entity.SuggestedOfferConfig](), res)
		assert.Nil(t, err)
	})
}

func TestSuggestedOfferConfigRepository_UpdateStatus(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewSuggestedOfferConfigRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	suggestedOfferConfig := entity.SuggestedOfferConfig{
		Id:            1,
		Name:          "Test 1",
		Value:         decimal.NewFromFloat(6.99),
		ValueType:     entity.ValueTypeInterestRate,
		Status:        "ACTIVE",
		CreatedBy:     "",
		LastUpdatedBy: "admin",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	// Test case
	t.Run("update suggested offer config status success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"suggested_offer_config.id",
			"suggested_offer_config.value_type",
			"suggested_offer_config.name",
			"suggested_offer_config.value",
			"suggested_offer_config.status",
			"suggested_offer_config.created_by",
			"suggested_offer_config.last_updated_by",
			"suggested_offer_config.created_at",
			"suggested_offer_config.updated_at",
		})
		rows.AddRow(
			suggestedOfferConfig.Id,
			suggestedOfferConfig.ValueType,
			suggestedOfferConfig.Name,
			suggestedOfferConfig.Value,
			suggestedOfferConfig.Status,
			suggestedOfferConfig.CreatedBy,
			suggestedOfferConfig.LastUpdatedBy,
			suggestedOfferConfig.CreatedAt,
			suggestedOfferConfig.UpdatedAt,
		)
		mock.ExpectQuery("UPDATE").WillReturnRows(rows)
		res, err := repo.UpdateStatus(context.Background(), suggestedOfferConfig.Id, entity.SuggestedOfferConfigStatusActive, suggestedOfferConfig.LastUpdatedBy)
		assert.Nil(t, err)
		assert.Equal(t, suggestedOfferConfig, res)
	})
	t.Run("update suggested offer config status error", func(t *testing.T) {
		mock.ExpectQuery("UPDATE").WillReturnError(assert.AnError)
		res, err := repo.UpdateStatus(context.Background(), suggestedOfferConfig.Id, entity.SuggestedOfferConfigStatusActive, suggestedOfferConfig.LastUpdatedBy)
		assert.Empty(t, res)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestSuggestedOfferConfigRepository_Update(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewSuggestedOfferConfigRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	suggestedOfferConfig := entity.SuggestedOfferConfig{
		Id:            1,
		Name:          "Test 1",
		Value:         decimal.NewFromFloat(6.99),
		ValueType:     entity.ValueTypeInterestRate,
		Status:        "ACTIVE",
		CreatedBy:     "",
		LastUpdatedBy: "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	// Test case
	t.Run("update suggested offer config success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"suggested_offer_config.id",
			"suggested_offer_config.value_type",
			"suggested_offer_config.name",
			"suggested_offer_config.value",
			"suggested_offer_config.status",
			"suggested_offer_config.created_by",
			"suggested_offer_config.last_updated_by",
			"suggested_offer_config.created_at",
			"suggested_offer_config.updated_at",
		})
		rows.AddRow(
			suggestedOfferConfig.Id,
			suggestedOfferConfig.ValueType,
			suggestedOfferConfig.Name,
			suggestedOfferConfig.Value,
			suggestedOfferConfig.Status,
			suggestedOfferConfig.CreatedBy,
			suggestedOfferConfig.LastUpdatedBy,
			suggestedOfferConfig.CreatedAt,
			suggestedOfferConfig.UpdatedAt,
		)
		mock.ExpectQuery("UPDATE").WillReturnRows(rows)
		res, err := repo.Update(context.Background(), suggestedOfferConfig)
		assert.Nil(t, err)
		assert.Equal(t, suggestedOfferConfig, res)
	})
	t.Run("update suggested offer config error", func(t *testing.T) {
		mock.ExpectQuery("UPDATE").WillReturnError(assert.AnError)
		res, err := repo.Update(context.Background(), suggestedOfferConfig)
		assert.Empty(t, res)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestSuggestedOfferConfigRepository_Create(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewSuggestedOfferConfigRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	suggestedOfferConfig := entity.SuggestedOfferConfig{
		Id:            1,
		Name:          "Test 1",
		Value:         decimal.NewFromFloat(6.99),
		ValueType:     entity.ValueTypeInterestRate,
		Status:        "ACTIVE",
		CreatedBy:     "",
		LastUpdatedBy: "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	// Test case
	t.Run("create suggested offer config success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"suggested_offer_config.id",
			"suggested_offer_config.value_type",
			"suggested_offer_config.name",
			"suggested_offer_config.value",
			"suggested_offer_config.status",
			"suggested_offer_config.created_by",
			"suggested_offer_config.last_updated_by",
			"suggested_offer_config.created_at",
			"suggested_offer_config.updated_at",
		})
		rows.AddRow(
			suggestedOfferConfig.Id,
			suggestedOfferConfig.ValueType,
			suggestedOfferConfig.Name,
			suggestedOfferConfig.Value,
			suggestedOfferConfig.Status,
			suggestedOfferConfig.CreatedBy,
			suggestedOfferConfig.LastUpdatedBy,
			suggestedOfferConfig.CreatedAt,
			suggestedOfferConfig.UpdatedAt,
		)
		mock.ExpectQuery("INSERT").WillReturnRows(rows)
		res, err := repo.Create(context.Background(), suggestedOfferConfig)
		assert.Nil(t, err)
		assert.Equal(t, suggestedOfferConfig, res)
	})
	t.Run("create suggested offer config error", func(t *testing.T) {
		mock.ExpectQuery("INSERT").WillReturnError(assert.AnError)
		res, err := repo.Create(context.Background(), suggestedOfferConfig)
		assert.Empty(t, res)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestSuggestedOfferConfigRepository_GetById(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewSuggestedOfferConfigRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	suggestedOfferConfig := entity.SuggestedOfferConfig{
		Id:            1,
		Name:          "Test 1",
		Value:         decimal.NewFromFloat(6.99),
		ValueType:     entity.ValueTypeInterestRate,
		Status:        "ACTIVE",
		CreatedBy:     "",
		LastUpdatedBy: "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	// Test case
	t.Run("get suggested offer config by id success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"suggested_offer_config.id",
			"suggested_offer_config.value_type",
			"suggested_offer_config.name",
			"suggested_offer_config.value",
			"suggested_offer_config.status",
			"suggested_offer_config.created_by",
			"suggested_offer_config.last_updated_by",
			"suggested_offer_config.created_at",
			"suggested_offer_config.updated_at",
		})
		rows.AddRow(
			suggestedOfferConfig.Id,
			suggestedOfferConfig.ValueType,
			suggestedOfferConfig.Name,
			suggestedOfferConfig.Value,
			suggestedOfferConfig.Status,
			suggestedOfferConfig.CreatedBy,
			suggestedOfferConfig.LastUpdatedBy,
			suggestedOfferConfig.CreatedAt,
			suggestedOfferConfig.UpdatedAt,
		)
		mock.ExpectQuery("SELECT .* FROM public.suggested_offer_config").WillReturnRows(rows)
		res, err := repo.GetById(context.Background(), suggestedOfferConfig.Id)
		assert.Nil(t, err)
		assert.Equal(t, suggestedOfferConfig, res)
	})
	t.Run("get suggested offer config by id error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM public.suggested_offer_config").WillReturnError(assert.AnError)
		res, err := repo.GetById(context.Background(), suggestedOfferConfig.Id)
		assert.Empty(t, res)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
