package suggested_offer_config

import (
	"context"
	"fmt"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/suggested_offer_config/repository"
	"financing-offer/pkg/optional"
)

type UseCase interface {
	GetAll(ctx context.Context) ([]entity.SuggestedOfferConfig, error)
	GetActiveSuggestedOfferConfig(ctx context.Context) (optional.Optional[entity.SuggestedOfferConfig], error)
	GetById(ctx context.Context, id int64) (entity.SuggestedOfferConfig, error)
	Update(ctx context.Context, suggestedOfferConfig entity.SuggestedOfferConfig) (entity.SuggestedOfferConfig, error)
	Create(ctx context.Context, suggestedOfferConfig entity.SuggestedOfferConfig) (entity.SuggestedOfferConfig, error)
	UpdateStatus(ctx context.Context, id int64, status entity.SuggestedOfferConfigStatus, updater string) (entity.SuggestedOfferConfig, error)
}

type useCase struct {
	repository repository.SuggestedOfferConfigRepository
}

func (u *useCase) GetAll(ctx context.Context) ([]entity.SuggestedOfferConfig, error) {
	res, err := u.repository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("suggestedOfferConfigUseCase GetAll %w", err)
	}
	return res, nil
}

func (u *useCase) GetActiveSuggestedOfferConfig(ctx context.Context) (optional.Optional[entity.SuggestedOfferConfig], error) {
	res, err := u.repository.GetActiveSuggestedOfferConfig(ctx)
	if err != nil {
		return optional.None[entity.SuggestedOfferConfig](), fmt.Errorf("suggestedOfferConfigUseCase GetActiveSuggestedOfferConfig %w", err)
	}
	return res, nil
}

func (u *useCase) GetById(ctx context.Context, id int64) (entity.SuggestedOfferConfig, error) {
	res, err := u.repository.GetById(ctx, id)
	if err != nil {
		return entity.SuggestedOfferConfig{}, fmt.Errorf("suggestedOfferConfigUseCase GetById %w", err)
	}
	return res, nil
}

func (u *useCase) Create(ctx context.Context, suggestedOfferConfig entity.SuggestedOfferConfig) (entity.SuggestedOfferConfig, error) {
	errorTemplate := "suggestedOfferConfigUseCase Create %w"
	res, err := u.repository.Create(ctx, suggestedOfferConfig)
	if err != nil {
		return entity.SuggestedOfferConfig{}, fmt.Errorf(errorTemplate, err)
	}
	return res, nil
}

func (u *useCase) Update(ctx context.Context, suggestedOfferConfig entity.SuggestedOfferConfig) (entity.SuggestedOfferConfig, error) {
	errorTemplate := "suggestedOfferConfigUseCase Update %w"
	res, err := u.repository.Update(ctx, suggestedOfferConfig)
	if err != nil {
		return entity.SuggestedOfferConfig{}, fmt.Errorf(errorTemplate, err)
	}
	return res, nil
}

func (u *useCase) UpdateStatus(ctx context.Context, id int64, status entity.SuggestedOfferConfigStatus, updater string) (entity.SuggestedOfferConfig, error) {
	errorTemplate := "suggestedOfferConfigUseCase UpdateStatus %w"
	if status == entity.SuggestedOfferConfigStatusActive {
		currentActive, err := u.repository.GetActiveSuggestedOfferConfig(ctx)
		if err != nil {
			return entity.SuggestedOfferConfig{}, fmt.Errorf(errorTemplate, err)
		}
		// If there is an active suggested offer config, return an error
		if currentActive.IsPresent() {
			return entity.SuggestedOfferConfig{}, apperrors.ErrorExistActiveSuggestedOfferConfig
		}
	}
	res, err := u.repository.UpdateStatus(ctx, id, status, updater)
	if err != nil {
		return entity.SuggestedOfferConfig{}, fmt.Errorf(errorTemplate, err)
	}
	return res, nil
}

func NewUseCase(repository repository.SuggestedOfferConfigRepository) UseCase {
	return &useCase{
		repository: repository,
	}
}
