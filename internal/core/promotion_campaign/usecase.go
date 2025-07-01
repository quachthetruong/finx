package promotion_campaign

import (
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/promotion_campaign/repository"
	"fmt"
	"golang.org/x/net/context"
)

type UseCase interface {
	GetAll(ctx context.Context, filter entity.GetPromotionCampaignsRequest) ([]entity.PromotionCampaign, error)
	GetById(ctx context.Context, id int64) (entity.PromotionCampaign, error)
	Update(ctx context.Context, symbol entity.PromotionCampaign) (entity.PromotionCampaign, error)
	Create(ctx context.Context, symbol entity.PromotionCampaign) (entity.PromotionCampaign, error)
}

type useCase struct {
	repository repository.PromotionCampaignRepository
}

func NewUseCase(
	repository repository.PromotionCampaignRepository,
) UseCase {
	return &useCase{
		repository: repository,
	}
}

func (s *useCase) GetAll(ctx context.Context, filter entity.GetPromotionCampaignsRequest) ([]entity.PromotionCampaign, error) {
	errorWrapMsg := "UseCase GetAll %w"
	res, err := s.repository.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf(errorWrapMsg, err)
	}
	return res, nil
}

func (s *useCase) Update(ctx context.Context, campaign entity.PromotionCampaign) (entity.PromotionCampaign, error) {
	errorWrapMsg := "UseCase Update %w"
	res, err := s.repository.Update(ctx, campaign)
	if err != nil {
		return entity.PromotionCampaign{}, fmt.Errorf(errorWrapMsg, err)
	}
	return res, nil
}

func (s *useCase) Create(ctx context.Context, campaign entity.PromotionCampaign) (entity.PromotionCampaign, error) {
	res, err := s.repository.Create(ctx, campaign)
	if err != nil {
		return res, fmt.Errorf("UseCase Create %w", err)
	}
	return res, nil
}

func (s *useCase) GetById(ctx context.Context, id int64) (entity.PromotionCampaign, error) {
	res, err := s.repository.GetById(ctx, id)
	if err != nil {
		return entity.PromotionCampaign{}, fmt.Errorf("UseCase GetById %w", err)
	}
	return res, nil
}
