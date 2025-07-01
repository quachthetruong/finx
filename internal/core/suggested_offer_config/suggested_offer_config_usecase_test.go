package suggested_offer_config

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/core/entity"
	"financing-offer/pkg/optional"
	"financing-offer/test/mock"
)

func TestSuggestedOfferConfigUseCase_GetAll(t *testing.T) {
	t.Parallel()

	t.Run("get all suggested offer config success", func(t *testing.T) {
		repository := mock.NewMockSuggestedOfferConfigRepository(t)
		useCase := NewUseCase(repository)
		suggestedOfferConfigs := []entity.SuggestedOfferConfig{
			{
				Id:            1,
				Name:          "Test 1",
				Value:         decimal.NewFromFloat(6.99),
				ValueType:     entity.ValueTypeInterestRate,
				Status:        entity.SuggestedOfferConfigStatusActive,
				CreatedBy:     "",
				LastUpdatedBy: "",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}
		repository.EXPECT().GetAll(mock2.Anything).Return(suggestedOfferConfigs, nil)
		res, err := useCase.GetAll(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, suggestedOfferConfigs, res)
	})

	t.Run("get all suggested offer config error", func(t *testing.T) {
		repository := mock.NewMockSuggestedOfferConfigRepository(t)
		useCase := NewUseCase(repository)
		repository.EXPECT().GetAll(mock2.Anything).Return(nil, assert.AnError)
		_, err := useCase.GetAll(context.Background())
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestSuggestedOfferConfigUseCase_GetActiveSuggestedOfferConfig(t *testing.T) {
	t.Parallel()

	t.Run("get active suggested offer config success", func(t *testing.T) {
		repository := mock.NewMockSuggestedOfferConfigRepository(t)
		useCase := NewUseCase(repository)
		suggestedOfferConfig := optional.Some(entity.SuggestedOfferConfig{
			Id:            1,
			Name:          "Test 1",
			Value:         decimal.NewFromFloat(6.99),
			ValueType:     entity.ValueTypeInterestRate,
			Status:        entity.SuggestedOfferConfigStatusActive,
			CreatedBy:     "",
			LastUpdatedBy: "",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
		repository.EXPECT().GetActiveSuggestedOfferConfig(mock2.Anything).Return(suggestedOfferConfig, nil)
		res, err := useCase.GetActiveSuggestedOfferConfig(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, suggestedOfferConfig, res)
	})

	t.Run("get active suggested offer config error", func(t *testing.T) {
		repository := mock.NewMockSuggestedOfferConfigRepository(t)
		useCase := NewUseCase(repository)
		repository.EXPECT().GetActiveSuggestedOfferConfig(mock2.Anything).Return(optional.None[entity.SuggestedOfferConfig](), assert.AnError)
		res, err := useCase.GetActiveSuggestedOfferConfig(context.Background())
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, res)
	})
}

func TestSuggestedOfferConfigUseCase_GetById(t *testing.T) {
	t.Parallel()

	t.Run("get by id suggested offer config success", func(t *testing.T) {
		repository := mock.NewMockSuggestedOfferConfigRepository(t)
		useCase := NewUseCase(repository)
		suggestedOfferConfig := entity.SuggestedOfferConfig{
			Id:            1,
			Name:          "Test 1",
			Value:         decimal.NewFromFloat(6.99),
			ValueType:     entity.ValueTypeInterestRate,
			Status:        entity.SuggestedOfferConfigStatusActive,
			CreatedBy:     "",
			LastUpdatedBy: "",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		repository.EXPECT().GetById(mock2.Anything, suggestedOfferConfig.Id).Return(suggestedOfferConfig, nil)
		res, err := useCase.GetById(context.Background(), 1)
		assert.Nil(t, err)
		assert.Equal(t, suggestedOfferConfig, res)
	})

	t.Run("get by id suggested offer config error", func(t *testing.T) {
		repository := mock.NewMockSuggestedOfferConfigRepository(t)
		useCase := NewUseCase(repository)
		repository.EXPECT().GetById(mock2.Anything, int64(1)).Return(entity.SuggestedOfferConfig{}, assert.AnError)
		_, err := useCase.GetById(context.Background(), 1)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestSuggestedOfferConfigUseCase_Create(t *testing.T) {
	t.Parallel()

	t.Run("create suggested offer config success", func(t *testing.T) {
		repository := mock.NewMockSuggestedOfferConfigRepository(t)
		useCase := NewUseCase(repository)
		suggestedOfferConfig := entity.SuggestedOfferConfig{
			Id:            1,
			Name:          "Test 1",
			Value:         decimal.NewFromFloat(6.99),
			ValueType:     entity.ValueTypeInterestRate,
			Status:        entity.SuggestedOfferConfigStatusActive,
			CreatedBy:     "",
			LastUpdatedBy: "",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		repository.EXPECT().Create(mock2.Anything, suggestedOfferConfig).Return(suggestedOfferConfig, nil)
		res, err := useCase.Create(context.Background(), suggestedOfferConfig)
		assert.Nil(t, err)
		assert.Equal(t, suggestedOfferConfig, res)
	})

	t.Run("create suggested offer config error", func(t *testing.T) {
		repository := mock.NewMockSuggestedOfferConfigRepository(t)
		useCase := NewUseCase(repository)
		suggestedOfferConfig := entity.SuggestedOfferConfig{
			Id:            1,
			Name:          "Test 1",
			Value:         decimal.NewFromFloat(6.99),
			ValueType:     entity.ValueTypeInterestRate,
			Status:        entity.SuggestedOfferConfigStatusInactive,
			CreatedBy:     "",
			LastUpdatedBy: "",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		repository.EXPECT().Create(mock2.Anything, suggestedOfferConfig).Return(entity.SuggestedOfferConfig{}, assert.AnError)
		_, err := useCase.Create(context.Background(), suggestedOfferConfig)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestSuggestedOfferConfigUseCase_Update(t *testing.T) {
	t.Parallel()

	t.Run("update suggested offer config success", func(t *testing.T) {
		repository := mock.NewMockSuggestedOfferConfigRepository(t)
		useCase := NewUseCase(repository)
		suggestedOfferConfig := entity.SuggestedOfferConfig{
			Id:        1,
			Name:      "Test 1",
			Value:     decimal.NewFromFloat(6.99),
			ValueType: entity.ValueTypeInterestRate,
			Status:    entity.SuggestedOfferConfigStatusInactive,
		}
		repository.EXPECT().Update(mock2.Anything, suggestedOfferConfig).Return(suggestedOfferConfig, nil)
		res, err := useCase.Update(context.Background(), suggestedOfferConfig)
		assert.Nil(t, err)
		assert.Equal(t, suggestedOfferConfig, res)
	})

	t.Run("update suggested offer config with active status fail update", func(t *testing.T) {
		repository := mock.NewMockSuggestedOfferConfigRepository(t)
		useCase := NewUseCase(repository)
		suggestedOfferConfig := entity.SuggestedOfferConfig{
			Id:        1,
			Name:      "Test 1",
			Value:     decimal.NewFromFloat(6.99),
			ValueType: entity.ValueTypeInterestRate,
			Status:    entity.SuggestedOfferConfigStatusActive,
		}
		repository.EXPECT().Update(mock2.Anything, suggestedOfferConfig).Return(entity.SuggestedOfferConfig{}, assert.AnError)
		_, err := useCase.Update(context.Background(), suggestedOfferConfig)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestSuggestedOfferConfigUseCase_UpdateStatus(t *testing.T) {
	t.Parallel()

	t.Run("update status suggested offer config with active status success return error", func(t *testing.T) {
		repository := mock.NewMockSuggestedOfferConfigRepository(t)
		useCase := NewUseCase(repository)
		activeSuggestedOfferConfig := entity.SuggestedOfferConfig{
			Id:        2,
			Name:      "Test 2",
			Value:     decimal.NewFromFloat(6.99),
			ValueType: entity.ValueTypeInterestRate,
			Status:    entity.SuggestedOfferConfigStatusActive,
		}
		suggestedOfferConfig := entity.SuggestedOfferConfig{
			Id:            1,
			Name:          "Test 1",
			Value:         decimal.NewFromFloat(6.99),
			ValueType:     entity.ValueTypeInterestRate,
			Status:        entity.SuggestedOfferConfigStatusActive,
			LastUpdatedBy: "admin",
		}
		repository.EXPECT().GetActiveSuggestedOfferConfig(mock2.Anything).Return(optional.Some(activeSuggestedOfferConfig), nil)
		res, err := useCase.UpdateStatus(context.Background(), suggestedOfferConfig.Id, suggestedOfferConfig.Status, suggestedOfferConfig.LastUpdatedBy)
		assert.ErrorIs(t, err, apperrors.ErrorExistActiveSuggestedOfferConfig)
		assert.Empty(t, res)
	})

	t.Run("update status suggested offer config with active status fail GetActiveSuggestedOfferConfig", func(t *testing.T) {
		repository := mock.NewMockSuggestedOfferConfigRepository(t)
		useCase := NewUseCase(repository)
		repository.EXPECT().GetActiveSuggestedOfferConfig(mock2.Anything).Return(optional.None[entity.SuggestedOfferConfig](), assert.AnError)
		_, err := useCase.UpdateStatus(context.Background(), 1, entity.SuggestedOfferConfigStatusActive, "admin")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("update status suggested offer config with active status fail UpdateStatus", func(t *testing.T) {
		repository := mock.NewMockSuggestedOfferConfigRepository(t)
		useCase := NewUseCase(repository)
		suggestedOfferConfig := entity.SuggestedOfferConfig{
			Id:        1,
			Name:      "Test 1",
			Value:     decimal.NewFromFloat(6.99),
			ValueType: entity.ValueTypeInterestRate,
			Status:    entity.SuggestedOfferConfigStatusActive,
		}
		repository.EXPECT().GetActiveSuggestedOfferConfig(mock2.Anything).Return(optional.None[entity.SuggestedOfferConfig](), nil)
		repository.EXPECT().UpdateStatus(mock2.Anything, suggestedOfferConfig.Id, entity.SuggestedOfferConfigStatusActive, "admin").Return(entity.SuggestedOfferConfig{}, assert.AnError)
		_, err := useCase.UpdateStatus(context.Background(), suggestedOfferConfig.Id, suggestedOfferConfig.Status, "admin")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("update status suggested offer config with active status and none active config success", func(t *testing.T) {
		repository := mock.NewMockSuggestedOfferConfigRepository(t)
		useCase := NewUseCase(repository)
		suggestedOfferConfig := entity.SuggestedOfferConfig{
			Id:        1,
			Name:      "Test 1",
			Value:     decimal.NewFromFloat(6.99),
			ValueType: entity.ValueTypeInterestRate,
			Status:    entity.SuggestedOfferConfigStatusActive,
		}
		repository.EXPECT().GetActiveSuggestedOfferConfig(mock2.Anything).Return(optional.None[entity.SuggestedOfferConfig](), nil)
		repository.EXPECT().UpdateStatus(mock2.Anything, suggestedOfferConfig.Id, entity.SuggestedOfferConfigStatusActive, "admin").Return(suggestedOfferConfig, nil)
		res, err := useCase.UpdateStatus(context.Background(), suggestedOfferConfig.Id, suggestedOfferConfig.Status, "admin")
		assert.Nil(t, err)
		assert.Equal(t, suggestedOfferConfig, res)
	})
}
