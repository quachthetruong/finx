package postgres

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/offline_offer_update/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

var _ repository.OfflineOfferUpdatePersistenceRepository = &OfflineOfferUpdatePostgresRepository{}

type OfflineOfferUpdatePostgresRepository struct {
	getDbFunc database.GetDbFunc
}

func (r *OfflineOfferUpdatePostgresRepository) GetByOfferId(ctx context.Context, offerId int64, assetType entity.AssetType) ([]entity.OfflineOfferUpdate, error) {
	updates := make([]model.OfflineOfferUpdate, 0)
	if err := table.OfflineOfferUpdate.
		SELECT(table.OfflineOfferUpdate.AllColumns).
		FROM(
			table.OfflineOfferUpdate.INNER_JOIN(table.LoanPackageOffer, table.LoanPackageOffer.ID.EQ(table.OfflineOfferUpdate.OfferID)).
				INNER_JOIN(table.LoanPackageRequest, table.LoanPackageOffer.LoanPackageRequestID.EQ(table.LoanPackageRequest.ID)),
		).
		WHERE(
			table.OfflineOfferUpdate.OfferID.EQ(postgres.Int64(offerId)).
				AND(table.LoanPackageRequest.AssetType.EQ(postgres.NewEnumValue(assetType.String()))),
		).
		ORDER_BY(table.OfflineOfferUpdate.ID.DESC()).
		QueryContext(ctx, r.getDbFunc(ctx), &updates); err != nil {
		return nil, fmt.Errorf("OfflineOfferUpdatePostgresRepository GetByOfferId %w", err)
	}
	return MapOfflineOfferUpdatesDbToEntity(updates), nil
}

func (r *OfflineOfferUpdatePostgresRepository) Create(ctx context.Context, offlineOfferUpdate entity.OfflineOfferUpdate) (entity.OfflineOfferUpdate, error) {
	created := model.OfflineOfferUpdate{}
	if err := table.OfflineOfferUpdate.
		INSERT(table.OfflineOfferUpdate.MutableColumns).
		MODEL(MapOfflineOfferUpdateEntityToDb(offlineOfferUpdate)).RETURNING(table.OfflineOfferUpdate.AllColumns).QueryContext(
		ctx, r.getDbFunc(ctx), &created,
	); err != nil {
		return entity.OfflineOfferUpdate{}, fmt.Errorf("OfflineOfferUpdatePostgresRepository Create %w", err)
	}
	return MapOfflineOfferUpdateDbToEntity(created), nil
}

func NewOfflineOfferUpdatePostgresRepository(getDbFunc database.GetDbFunc) *OfflineOfferUpdatePostgresRepository {
	return &OfflineOfferUpdatePostgresRepository{getDbFunc: getDbFunc}
}
