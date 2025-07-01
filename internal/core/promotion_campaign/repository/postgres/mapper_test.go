package postgres

import (
	"encoding/json"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMapPromotionCampaignEntityToDb(t *testing.T) {
	createdAt := time.Now()
	updatedAt := time.Now()
	metadata := entity.PromotionCampaignMetadata{
		Products: []entity.PromotionCampaignProduct{
			{
				Symbols:       []string{"HPG", "DGW"},
				LoanPackageId: 1,
				RetailSymbols: []string{"HPG", "DGW"},
			},
		},
	}
	metadataJson, _ := json.Marshal(metadata)

	entity := entity.PromotionCampaign{
		Id:          1,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		UpdatedBy:   "user1",
		Name:        "Campaign1",
		Tag:         "Tag1",
		Description: "Description1",
		Status:      entity.Active,
		Metadata:    metadata,
	}

	expected := model.PromotionCampaign{
		ID:          1,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		UpdatedBy:   "user1",
		Name:        "Campaign1",
		Tag:         "Tag1",
		Description: "Description1",
		Status:      "ACTIVE",
		Metadata:    string(metadataJson),
	}

	result, err := MapPromotionCampaignEntityToDb(entity)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestMapPromotionCampaignDbToEntity(t *testing.T) {
	createdAt := time.Now()
	updatedAt := time.Now()
	metadata := entity.PromotionCampaignMetadata{
		Products: []entity.PromotionCampaignProduct{
			{
				Symbols:       []string{"HPG", "DGW"},
				LoanPackageId: 1,
				RetailSymbols: []string{"HPG", "DGW"},
			},
		},
	}
	metadataJson, _ := json.Marshal(metadata)

	dbModel := model.PromotionCampaign{
		ID:          1,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		UpdatedBy:   "user1",
		Name:        "Campaign1",
		Tag:         "Tag1",
		Description: "Description1",
		Status:      "ACTIVE",
		Metadata:    string(metadataJson),
	}

	expected := entity.PromotionCampaign{
		Id:          1,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		UpdatedBy:   "user1",
		Name:        "Campaign1",
		Tag:         "Tag1",
		Description: "Description1",
		Status:      entity.Active,
		Metadata:    metadata,
	}

	result, err := MapPromotionCampaignDbToEntity(dbModel)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
