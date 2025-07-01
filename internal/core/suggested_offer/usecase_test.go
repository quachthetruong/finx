package suggested_offer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"

	"financing-offer/internal/core/entity"
	"financing-offer/test/mock"
)

func TestSuggestedOfferUseCase(t *testing.T) {
	t.Run("create suggested offer with invalid config id", func(t *testing.T) {
		offerRequest := entity.SuggestedOffer{
			ConfigId: 1,
		}
		repository := mock.NewMockSuggestedOfferRepository(t)
		configRepository := mock.NewMockSuggestedOfferConfigRepository(t)
		eventRepository := mock.NewMockSuggestedOfferEventRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		useCase := NewUseCase(repository, configRepository, eventRepository, orderServiceRepo)
		configRepository.EXPECT().GetById(testifyMock.Anything, offerRequest.ConfigId).Return(
			entity.SuggestedOfferConfig{}, errors.New("not found"))
		_, err := useCase.CreateSuggestedOffer(context.Background(), "investor", "custodyCode", offerRequest)
		assert.Equal(t, "suggestedOfferConfigUseCase CreateSuggestedOffer not found", err.Error())
	})

	t.Run("create suggested offer with inactive config", func(t *testing.T) {
		offerRequest := entity.SuggestedOffer{
			ConfigId: 1,
		}
		repository := mock.NewMockSuggestedOfferRepository(t)
		configRepository := mock.NewMockSuggestedOfferConfigRepository(t)
		eventRepository := mock.NewMockSuggestedOfferEventRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		useCase := NewUseCase(repository, configRepository, eventRepository, orderServiceRepo)
		configRepository.EXPECT().GetById(testifyMock.Anything, offerRequest.ConfigId).Return(
			entity.SuggestedOfferConfig{
				Status: entity.SuggestedOfferConfigStatusInactive,
			}, nil)
		_, err := useCase.CreateSuggestedOffer(context.Background(), "investor", "custodyCode", offerRequest)
		assert.Equal(t, "suggestedOfferConfigUseCase CreateSuggestedOffer inactive program", err.Error())
	})

	t.Run("create suggested offer error db", func(t *testing.T) {
		offerRequest := entity.SuggestedOffer{
			ConfigId:  1,
			AccountNo: "accountNo",
			Symbols:   []string{"ACB"},
		}
		repository := mock.NewMockSuggestedOfferRepository(t)
		configRepository := mock.NewMockSuggestedOfferConfigRepository(t)
		eventRepository := mock.NewMockSuggestedOfferEventRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		useCase := NewUseCase(repository, configRepository, eventRepository, orderServiceRepo)
		configRepository.EXPECT().GetById(testifyMock.Anything, offerRequest.ConfigId).Return(
			entity.SuggestedOfferConfig{
				Status: entity.SuggestedOfferConfigStatusActive,
			}, nil)
		repository.EXPECT().Create(testifyMock.Anything, offerRequest).Return(
			entity.SuggestedOffer{}, errors.New("duplicate key"))
		orderServiceRepo.EXPECT().GetAccountByAccountNoAndCustodyCode(testifyMock.Anything, "custodyCode", "accountNo").Return(entity.OrderServiceAccount{
			CustodyCode: "custodyCode",
			AccountNo:   "accountNo",
		}, nil)
		_, err := useCase.CreateSuggestedOffer(context.Background(), "investor", "custodyCode", offerRequest)
		assert.Equal(t, "suggestedOfferConfigUseCase CreateSuggestedOffer duplicate key", err.Error())
	})

	t.Run("create suggested offer with invalid accountNo", func(t *testing.T) {
		config := entity.SuggestedOfferConfig{
			Status: entity.SuggestedOfferConfigStatusActive,
		}
		offerRequest := entity.SuggestedOffer{
			ConfigId:  1,
			Config:    &config,
			AccountNo: "accountNo",
			Symbols:   []string{"ACB"},
		}
		repository := mock.NewMockSuggestedOfferRepository(t)
		configRepository := mock.NewMockSuggestedOfferConfigRepository(t)
		eventRepository := mock.NewMockSuggestedOfferEventRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		useCase := NewUseCase(repository, configRepository, eventRepository, orderServiceRepo)
		configRepository.EXPECT().GetById(testifyMock.Anything, offerRequest.ConfigId).Return(config, nil)
		orderServiceRepo.EXPECT().GetAccountByAccountNoAndCustodyCode(testifyMock.Anything, "custodyCode", "accountNo").Return(entity.OrderServiceAccount{}, errors.New("test"))
		_, err := useCase.CreateSuggestedOffer(context.Background(), "investor", "custodyCode", offerRequest)
		assert.Equal(t, "suggestedOfferConfigUseCase CreateSuggestedOffer test", err.Error())
	})

	t.Run("create suggested offer success", func(t *testing.T) {
		config := entity.SuggestedOfferConfig{
			Status: entity.SuggestedOfferConfigStatusActive,
		}
		offerRequest := entity.SuggestedOffer{
			ConfigId:  1,
			Config:    &config,
			AccountNo: "accountNo",
			Symbols:   []string{"ACB"},
		}
		repository := mock.NewMockSuggestedOfferRepository(t)
		configRepository := mock.NewMockSuggestedOfferConfigRepository(t)
		eventRepository := mock.NewMockSuggestedOfferEventRepository(t)
		orderServiceRepo := mock.NewMockOrderServiceRepository(t)
		useCase := NewUseCase(repository, configRepository, eventRepository, orderServiceRepo)
		configRepository.EXPECT().GetById(testifyMock.Anything, offerRequest.ConfigId).Return(config, nil)
		eventRepository.EXPECT().NotifySuggestedOfferCreated(testifyMock.Anything, "investor", config, offerRequest).Return(nil)
		repository.EXPECT().Create(testifyMock.Anything, offerRequest).Return(
			offerRequest, nil)
		orderServiceRepo.EXPECT().GetAccountByAccountNoAndCustodyCode(testifyMock.Anything, "custodyCode", "accountNo").Return(entity.OrderServiceAccount{
			CustodyCode: "custodyCode",
			AccountNo:   "accountNo",
		}, nil)
		createdOffer, err := useCase.CreateSuggestedOffer(context.Background(), "investor", "custodyCode", offerRequest)
		assert.Nil(t, err)
		assert.Equal(t, offerRequest, createdOffer)
	})
}
