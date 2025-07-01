package suggested_offer

import (
	"context"
	"errors"
	"fmt"

	"financing-offer/internal/core/entity"
	orderServiceRepo "financing-offer/internal/core/orderservice/repository"
	"financing-offer/internal/core/suggested_offer/repository"
	configRepository "financing-offer/internal/core/suggested_offer_config/repository"
)

type UseCase interface {
	CreateSuggestedOffer(ctx context.Context, investorId string, custodyCode string, suggestedOffer entity.SuggestedOffer) (entity.SuggestedOffer, error)
}

type useCase struct {
	repository       repository.SuggestedOfferRepository
	configRepository configRepository.SuggestedOfferConfigRepository
	eventRepository  repository.SuggestedOfferEventRepository
	orderServiceRepo orderServiceRepo.OrderServiceRepository
}

func (u *useCase) CreateSuggestedOffer(ctx context.Context, investorId string, custodyCode string, suggestedOffer entity.SuggestedOffer) (entity.SuggestedOffer, error) {
	errorTemplate := "suggestedOfferConfigUseCase CreateSuggestedOffer %w"
	config, err := u.configRepository.GetById(ctx, suggestedOffer.ConfigId)
	if err != nil {
		return entity.SuggestedOffer{}, fmt.Errorf(errorTemplate, err)
	}
	if config.Status != entity.SuggestedOfferConfigStatusActive {
		return entity.SuggestedOffer{}, fmt.Errorf(errorTemplate, errors.New("inactive program"))
	}
	_, err = u.orderServiceRepo.GetAccountByAccountNoAndCustodyCode(ctx, custodyCode, suggestedOffer.AccountNo)
	if err != nil {
		return entity.SuggestedOffer{}, fmt.Errorf(errorTemplate, err)
	}
	createdOffer, err := u.repository.Create(ctx, suggestedOffer)
	if err != nil {
		return entity.SuggestedOffer{}, fmt.Errorf(errorTemplate, err)
	}
	err = u.eventRepository.NotifySuggestedOfferCreated(ctx, investorId, config, createdOffer)
	if err != nil {
		return entity.SuggestedOffer{}, fmt.Errorf(errorTemplate, err)
	}
	createdOffer.Config = &config
	return createdOffer, nil
}

func NewUseCase(
	repository repository.SuggestedOfferRepository,
	configRepository configRepository.SuggestedOfferConfigRepository,
	eventRepository repository.SuggestedOfferEventRepository,
	orderServiceRepo orderServiceRepo.OrderServiceRepository,
) UseCase {
	return &useCase{
		repository:       repository,
		configRepository: configRepository,
		eventRepository:  eventRepository,
		orderServiceRepo: orderServiceRepo,
	}
}
