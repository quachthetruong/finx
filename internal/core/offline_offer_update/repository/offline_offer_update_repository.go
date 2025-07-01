package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type OfflineOfferUpdatePersistenceRepository interface {
	GetByOfferId(ctx context.Context, offerId int64, assetType entity.AssetType) ([]entity.OfflineOfferUpdate, error)
	Create(ctx context.Context, offlineOfferUpdate entity.OfflineOfferUpdate) (entity.OfflineOfferUpdate, error)
}
