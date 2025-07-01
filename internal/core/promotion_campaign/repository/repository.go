package repository

import (
	"context"
	"financing-offer/internal/core/entity"
)

type PromotionCampaignRepository interface {
	GetAll(ctx context.Context, filter entity.GetPromotionCampaignsRequest) ([]entity.PromotionCampaign, error)
	GetById(ctx context.Context, id int64) (entity.PromotionCampaign, error)
	Create(ctx context.Context, campaign entity.PromotionCampaign) (entity.PromotionCampaign, error)
	Update(ctx context.Context, campaign entity.PromotionCampaign) (entity.PromotionCampaign, error)
}
