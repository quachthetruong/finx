package postgres

import (
	"context"
	"fmt"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

type SuggestedOfferPostgresRepository struct {
	getDbFunc database.GetDbFunc
}

func NewSuggestedOfferRepository(getDbFunc database.GetDbFunc) *SuggestedOfferPostgresRepository {
	return &SuggestedOfferPostgresRepository{getDbFunc: getDbFunc}
}

func (s *SuggestedOfferPostgresRepository) Create(ctx context.Context, suggestedOffer entity.SuggestedOffer) (entity.SuggestedOffer, error) {
	errorTemplate := "SuggestedOfferConfigPostgresRepository Create %w"
	dest := model.SuggestedOffer{}
	toCreate, err := MapSuggestedOfferEntityToDb(suggestedOffer)
	if err != nil {
		return entity.SuggestedOffer{}, fmt.Errorf(errorTemplate, err)
	}
	if err = table.SuggestedOffer.
		INSERT(table.SuggestedOffer.MutableColumns).
		MODEL(toCreate).
		RETURNING(table.SuggestedOffer.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		return entity.SuggestedOffer{}, fmt.Errorf(errorTemplate, err)
	}
	cratedOffer, err := MapSuggestedOfferDbToEntity(dest)
	return cratedOffer, err
}
