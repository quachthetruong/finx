package postgres

import (
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

func MapOfflineOfferUpdateDbToEntity(offlineOfferUpdate model.OfflineOfferUpdate) entity.OfflineOfferUpdate {
	return entity.OfflineOfferUpdate{
		Id:        offlineOfferUpdate.ID,
		OfferId:   offlineOfferUpdate.OfferID,
		Status:    entity.OfflineOfferUpdateStatus(offlineOfferUpdate.Status),
		Category:  offlineOfferUpdate.Category,
		Note:      offlineOfferUpdate.Note,
		CreatedBy: offlineOfferUpdate.CreatedBy,
		CreatedAt: offlineOfferUpdate.CreatedAt,
	}
}

func MapOfflineOfferUpdatesDbToEntity(offlineOfferUpdates []model.OfflineOfferUpdate) []entity.OfflineOfferUpdate {
	offlineOfferUpdatesEntity := make([]entity.OfflineOfferUpdate, 0, len(offlineOfferUpdates))
	for _, offlineOfferUpdate := range offlineOfferUpdates {
		offlineOfferUpdatesEntity = append(
			offlineOfferUpdatesEntity, MapOfflineOfferUpdateDbToEntity(offlineOfferUpdate),
		)
	}
	return offlineOfferUpdatesEntity
}

func MapOfflineOfferUpdateEntityToDb(offlineOfferUpdate entity.OfflineOfferUpdate) model.OfflineOfferUpdate {
	return model.OfflineOfferUpdate{
		ID:        offlineOfferUpdate.Id,
		OfferID:   offlineOfferUpdate.OfferId,
		Status:    string(offlineOfferUpdate.Status),
		Category:  offlineOfferUpdate.Category,
		Note:      offlineOfferUpdate.Note,
		CreatedBy: offlineOfferUpdate.CreatedBy,
		CreatedAt: offlineOfferUpdate.CreatedAt,
	}
}

func MapLatestOfferUpdateDbToEntity(offlineOfferUpdate model.LatestOfferUpdate) entity.OfflineOfferUpdate {
	res := entity.OfflineOfferUpdate{}
	if offlineOfferUpdate.ID != nil {
		res.Id = *offlineOfferUpdate.ID
	}
	if offlineOfferUpdate.OfferID != nil {
		res.OfferId = *offlineOfferUpdate.OfferID
	}
	if offlineOfferUpdate.Status != nil {
		res.Status = entity.OfflineOfferUpdateStatus(*offlineOfferUpdate.Status)
	}
	if offlineOfferUpdate.Category != nil {
		res.Category = *offlineOfferUpdate.Category
	}
	if offlineOfferUpdate.Note != nil {
		res.Note = *offlineOfferUpdate.Note
	}
	if offlineOfferUpdate.CreatedBy != nil {
		res.CreatedBy = *offlineOfferUpdate.CreatedBy
	}
	if offlineOfferUpdate.CreatedAt != nil {
		res.CreatedAt = *offlineOfferUpdate.CreatedAt
	}
	return res
}
