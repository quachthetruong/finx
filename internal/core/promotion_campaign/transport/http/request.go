package http

import "financing-offer/internal/core/entity"

type CreatePromotionCampaignRequest struct {
	Name        string `json:"name" binding:"required"`
	Tag         string `json:"tag"  binding:"required"`
	Description string `json:"description"  binding:"required"`
}

func (r CreatePromotionCampaignRequest) toEntity() entity.PromotionCampaign {
	return entity.PromotionCampaign{
		Name:        r.Name,
		Tag:         r.Tag,
		Description: r.Description,
		Status:      entity.Active,
	}
}
