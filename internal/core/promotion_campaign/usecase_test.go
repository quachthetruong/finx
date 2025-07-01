package promotion_campaign

import (
	"context"
	"financing-offer/internal/core/entity"
	"financing-offer/test/mock"
	"github.com/stretchr/testify/assert"
	testify "github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestPromotionCampaignUseCase_GetAll(t *testing.T) {
	t.Parallel()

	t.Run("get all success", func(t *testing.T) {
		repository := mock.NewMockPromotionCampaignRepository(t)
		useCase := NewUseCase(repository)
		campaigns := []entity.PromotionCampaign{
			{
				Id:          1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				UpdatedBy:   "kiennt",
				Name:        "name",
				Tag:         "5.99*",
				Status:      entity.Active,
				Description: "description",
				Metadata: entity.PromotionCampaignMetadata{
					Products: []entity.PromotionCampaignProduct{
						{
							Symbols:       []string{"HPG", "DGW"},
							LoanPackageId: 1,
							RetailSymbols: []string{"HPG", "DGW"},
						},
					},
				},
			},
		}
		repository.EXPECT().GetAll(testify.Anything, testify.Anything).Return(campaigns, nil)
		res, err := useCase.GetAll(context.Background(), entity.GetPromotionCampaignsRequest{})
		assert.Nil(t, err)
		assert.Equal(t, campaigns, res)
	})

	t.Run("get all error", func(t *testing.T) {
		repository := mock.NewMockPromotionCampaignRepository(t)
		useCase := NewUseCase(repository)
		repository.EXPECT().GetAll(testify.Anything, testify.Anything).Return(nil, assert.AnError)
		_, err := useCase.GetAll(context.Background(), entity.GetPromotionCampaignsRequest{})
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestPromotionCampaignUseCase_GetById(t *testing.T) {
	t.Parallel()

	t.Run("get by id success", func(t *testing.T) {
		repository := mock.NewMockPromotionCampaignRepository(t)
		useCase := NewUseCase(repository)
		campaign := entity.PromotionCampaign{
			Id:          1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UpdatedBy:   "kiennt",
			Name:        "name",
			Tag:         "5.99*",
			Status:      entity.Active,
			Description: "description",
			Metadata: entity.PromotionCampaignMetadata{
				Products: []entity.PromotionCampaignProduct{
					{
						Symbols:       []string{"HPG", "DGW"},
						LoanPackageId: 1,
						RetailSymbols: []string{"HPG", "DGW"},
					},
				},
			},
		}
		repository.EXPECT().GetById(testify.Anything, campaign.Id).Return(campaign, nil)
		res, err := useCase.GetById(context.Background(), 1)
		assert.Nil(t, err)
		assert.Equal(t, campaign, res)
	})

	t.Run("get by id error", func(t *testing.T) {
		repository := mock.NewMockPromotionCampaignRepository(t)
		useCase := NewUseCase(repository)
		repository.EXPECT().GetById(testify.Anything, int64(1)).Return(entity.PromotionCampaign{}, assert.AnError)
		_, err := useCase.GetById(context.Background(), 1)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestPromotionCampaignUseCase_Create(t *testing.T) {
	t.Parallel()

	t.Run("create promotion campaign success", func(t *testing.T) {
		repository := mock.NewMockPromotionCampaignRepository(t)
		useCase := NewUseCase(repository)
		campaign := entity.PromotionCampaign{
			Id:          1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UpdatedBy:   "kiennt",
			Name:        "name",
			Tag:         "5.99*",
			Status:      entity.Active,
			Description: "description",
			Metadata: entity.PromotionCampaignMetadata{
				Products: []entity.PromotionCampaignProduct{
					{
						Symbols:       []string{"HPG", "DGW"},
						LoanPackageId: 1,
						RetailSymbols: []string{"HPG", "DGW"},
					},
				},
			},
		}
		repository.EXPECT().Create(testify.Anything, campaign).Return(campaign, nil)
		res, err := useCase.Create(context.Background(), campaign)
		assert.Nil(t, err)
		assert.Equal(t, campaign, res)
	})

	t.Run("create promotion campaign error", func(t *testing.T) {
		repository := mock.NewMockPromotionCampaignRepository(t)
		useCase := NewUseCase(repository)
		campaign := entity.PromotionCampaign{
			Id:          1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UpdatedBy:   "kiennt",
			Name:        "name",
			Tag:         "5.99*",
			Status:      entity.Active,
			Description: "description",
			Metadata: entity.PromotionCampaignMetadata{
				Products: []entity.PromotionCampaignProduct{
					{
						Symbols:       []string{"HPG", "DGW"},
						LoanPackageId: 1,
						RetailSymbols: []string{"HPG", "DGW"},
					},
				},
			},
		}
		repository.EXPECT().Create(testify.Anything, campaign).Return(entity.PromotionCampaign{}, assert.AnError)
		_, err := useCase.Create(context.Background(), campaign)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestPromotionCampaignUseCase_Update(t *testing.T) {
	t.Parallel()

	t.Run("update promotion campaign success", func(t *testing.T) {
		repository := mock.NewMockPromotionCampaignRepository(t)
		useCase := NewUseCase(repository)
		campaign := entity.PromotionCampaign{
			Id:          1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UpdatedBy:   "kiennt",
			Name:        "name",
			Tag:         "5.99*",
			Status:      entity.Active,
			Description: "description",
			Metadata: entity.PromotionCampaignMetadata{
				Products: []entity.PromotionCampaignProduct{
					{
						Symbols:       []string{"HPG", "DGW"},
						LoanPackageId: 1,
						RetailSymbols: []string{"HPG", "DGW"},
					},
				},
			},
		}
		repository.EXPECT().Update(testify.Anything, campaign).Return(campaign, nil)
		res, err := useCase.Update(context.Background(), campaign)
		assert.Nil(t, err)
		assert.Equal(t, campaign, res)
	})

	t.Run("promotion campaign fail", func(t *testing.T) {
		repository := mock.NewMockPromotionCampaignRepository(t)
		useCase := NewUseCase(repository)
		campaign := entity.PromotionCampaign{
			Id:          1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UpdatedBy:   "kiennt",
			Name:        "name",
			Tag:         "5.99*",
			Status:      entity.Active,
			Description: "description",
			Metadata: entity.PromotionCampaignMetadata{
				Products: []entity.PromotionCampaignProduct{
					{
						Symbols:       []string{"HPG", "DGW"},
						LoanPackageId: 1,
						RetailSymbols: []string{"HPG", "DGW"},
					},
				},
			},
		}
		repository.EXPECT().Update(testify.Anything, campaign).Return(entity.PromotionCampaign{}, assert.AnError)
		_, err := useCase.Update(context.Background(), campaign)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
