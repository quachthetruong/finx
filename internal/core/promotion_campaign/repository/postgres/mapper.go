package postgres

import (
	"encoding/json"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

func MapPromotionCampaignEntityToDb(c entity.PromotionCampaign) (model.PromotionCampaign, error) {
	res := model.PromotionCampaign{
		ID:          c.Id,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		UpdatedBy:   c.UpdatedBy,
		Name:        c.Name,
		Tag:         c.Tag,
		Description: c.Description,
		Status:      c.Status.String(),
	}
	metadata, err := json.Marshal(c.Metadata)
	if err != nil {
		return model.PromotionCampaign{}, err
	}
	res.Metadata = string(metadata)
	return res, nil
}

func MapPromotionCampaignDbToEntity(o model.PromotionCampaign) (entity.PromotionCampaign, error) {
	res := entity.PromotionCampaign{
		Id:          o.ID,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
		UpdatedBy:   o.UpdatedBy,
		Name:        o.Name,
		Tag:         o.Tag,
		Description: o.Description,
		Status:      entity.PromotionCampaignStatusFromString(o.Status),
	}
	var metadata entity.PromotionCampaignMetadata
	err := json.Unmarshal([]byte(o.Metadata), &metadata)
	if err != nil {
		return entity.PromotionCampaign{}, err
	}
	res.Metadata = metadata
	return res, nil
}
