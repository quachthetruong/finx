package postgres

import (
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

func MapSuggestedOfferConfigsDbToEntity(suggestedOfferConfigs []model.SuggestedOfferConfig) []entity.SuggestedOfferConfig {
	res := make([]entity.SuggestedOfferConfig, 0, len(suggestedOfferConfigs))
	for _, v := range suggestedOfferConfigs {
		res = append(res, MapSuggestedOfferConfigDbToEntity(v))
	}
	return res
}

func MapSuggestedOfferConfigDbToEntity(suggestedOfferConfig model.SuggestedOfferConfig) entity.SuggestedOfferConfig {
	return entity.SuggestedOfferConfig{
		Id:            suggestedOfferConfig.ID,
		Name:          suggestedOfferConfig.Name,
		Value:         suggestedOfferConfig.Value,
		ValueType:     entity.ValueType(suggestedOfferConfig.ValueType),
		Status:        entity.SuggestionsOfferConfigStatusFromString(suggestedOfferConfig.Status),
		CreatedBy:     suggestedOfferConfig.CreatedBy,
		LastUpdatedBy: suggestedOfferConfig.LastUpdatedBy,
		CreatedAt:     suggestedOfferConfig.CreatedAt,
		UpdatedAt:     suggestedOfferConfig.UpdatedAt,
	}
}

func MapSuggestedOfferConfigEntityToDb(suggestedOfferConfig entity.SuggestedOfferConfig) model.SuggestedOfferConfig {
	return model.SuggestedOfferConfig{
		ID:            suggestedOfferConfig.Id,
		Name:          suggestedOfferConfig.Name,
		Value:         suggestedOfferConfig.Value,
		ValueType:     suggestedOfferConfig.ValueType.String(),
		Status:        suggestedOfferConfig.Status.String(),
		CreatedBy:     suggestedOfferConfig.CreatedBy,
		LastUpdatedBy: suggestedOfferConfig.LastUpdatedBy,
		CreatedAt:     suggestedOfferConfig.CreatedAt,
		UpdatedAt:     suggestedOfferConfig.UpdatedAt,
	}
}
