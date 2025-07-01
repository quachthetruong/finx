package postgres

import (
	"encoding/json"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

func MapSuggestedOfferDbToEntity(suggestedOffer model.SuggestedOffer) (entity.SuggestedOffer, error) {
	var symbols []string
	unmarshalErr := json.Unmarshal([]byte(suggestedOffer.Symbols), &symbols)
	if unmarshalErr != nil {
		return entity.SuggestedOffer{}, unmarshalErr
	}
	return entity.SuggestedOffer{
		Id:        suggestedOffer.ID,
		AccountNo: suggestedOffer.AccountNo,
		ConfigId:  suggestedOffer.ConfigID,
		Symbols:   symbols,
		CreatedAt: suggestedOffer.CreatedAt,
		UpdatedAt: suggestedOffer.UpdatedAt,
	}, nil
}

func MapSuggestedOfferEntityToDb(suggestedOffer entity.SuggestedOffer) (model.SuggestedOffer, error) {
	symbolsByte, marshalErr := json.Marshal(suggestedOffer.Symbols)
	if marshalErr != nil {
		return model.SuggestedOffer{}, marshalErr
	}
	configId := suggestedOffer.ConfigId
	if suggestedOffer.Config != nil {
		configId = suggestedOffer.Config.Id
	}
	return model.SuggestedOffer{
		ID:        suggestedOffer.Id,
		AccountNo: suggestedOffer.AccountNo,
		ConfigID:  configId,
		Symbols:   string(symbolsByte),
		CreatedAt: suggestedOffer.CreatedAt,
		UpdatedAt: suggestedOffer.UpdatedAt,
	}, nil
}
