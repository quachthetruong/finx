package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type SuggestedOfferRepository interface {
	Create(ctx context.Context, suggestedOffer entity.SuggestedOffer) (entity.SuggestedOffer, error)
}

type SuggestedOfferEventRepository interface {
	NotifySuggestedOfferCreated(ctx context.Context, investorId string, config entity.SuggestedOfferConfig, createdOffer entity.SuggestedOffer) error
}
