package repository

import (
	"context"

	"financing-offer/internal/core/entity"
	"financing-offer/pkg/optional"
)

type SuggestedOfferConfigRepository interface {
	GetAll(ctx context.Context) ([]entity.SuggestedOfferConfig, error)
	GetActiveSuggestedOfferConfig(ctx context.Context) (optional.Optional[entity.SuggestedOfferConfig], error)
	UpdateStatus(ctx context.Context, id int64, status entity.SuggestedOfferConfigStatus, updater string) (entity.SuggestedOfferConfig, error)
	Update(ctx context.Context, suggestedOfferConfig entity.SuggestedOfferConfig) (entity.SuggestedOfferConfig, error)
	Create(ctx context.Context, suggestedOfferConfig entity.SuggestedOfferConfig) (entity.SuggestedOfferConfig, error)
	GetById(ctx context.Context, id int64) (entity.SuggestedOfferConfig, error)
}
