package postgres

import (
	"context"
	"encoding/json"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/pkg/dbtest"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPromotionCampaignRepository_GetAll(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewPromotionCampaignRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	metadata := entity.PromotionCampaignMetadata{
		Products: []entity.PromotionCampaignProduct{
			{
				Symbols:       []string{"HPG", "DGW"},
				LoanPackageId: 1,
				RetailSymbols: []string{"HPG", "DGW"},
			},
		},
	}
	metadataStr, err := json.Marshal(metadata)
	entities := []entity.PromotionCampaign{
		{
			Id:          1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UpdatedBy:   "kiennt",
			Name:        "name",
			Tag:         "5.99*",
			Status:      entity.Active,
			Description: "description",
			Metadata:    metadata,
		},
		{
			Id:        2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UpdatedBy: "kiennt2",
			Name:      "name2",
			Tag:       "5.99*",
			Status:    entity.Active,
			Metadata:  metadata,
		},
	}
	t.Run("get all success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"promotion_campaign.id",
			"promotion_campaign.created_at",
			"promotion_campaign.updated_at",
			"promotion_campaign.updated_by",
			"promotion_campaign.name",
			"promotion_campaign.tag",
			"promotion_campaign.description",
			"promotion_campaign.status",
			"promotion_campaign.metadata",
		})
		for _, data := range entities {
			rows.AddRow(
				data.Id,
				data.CreatedAt,
				data.UpdatedAt,
				data.UpdatedBy,
				data.Name,
				data.Tag,
				data.Description,
				data.Status,
				metadataStr,
			)
		}
		mock.ExpectQuery("SELECT .* FROM public.promotion_campaign").WillReturnRows(rows)
		res, err := repo.GetAll(context.Background(), entity.GetPromotionCampaignsRequest{})
		assert.Nil(t, err)
		assert.Equal(t, entities, res)
	})

	t.Run("get all error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM public.promotion_campaign").WillReturnError(assert.AnError)
		res, err := repo.GetAll(context.Background(), entity.GetPromotionCampaignsRequest{})
		assert.Nil(t, res)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestPromotionCampaignRepository_Update(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewPromotionCampaignRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	metadata := entity.PromotionCampaignMetadata{
		Products: []entity.PromotionCampaignProduct{
			{
				Symbols:       []string{"HPG", "DGW"},
				LoanPackageId: 1,
				RetailSymbols: []string{"HPG", "DGW"},
			},
		},
	}
	metadataStr, err := json.Marshal(metadata)
	campaign := entity.PromotionCampaign{
		Id:          1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UpdatedBy:   "kiennt",
		Name:        "name",
		Tag:         "5.99*",
		Status:      entity.Active,
		Description: "description",
		Metadata:    metadata,
	}
	t.Run("update promotion campaign success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"promotion_campaign.id",
			"promotion_campaign.created_at",
			"promotion_campaign.updated_at",
			"promotion_campaign.updated_by",
			"promotion_campaign.name",
			"promotion_campaign.tag",
			"promotion_campaign.description",
			"promotion_campaign.status",
			"promotion_campaign.metadata",
		})
		rows.AddRow(
			campaign.Id,
			campaign.CreatedAt,
			campaign.UpdatedAt,
			campaign.UpdatedBy,
			campaign.Name,
			campaign.Tag,
			campaign.Description,
			campaign.Status,
			metadataStr,
		)
		mock.ExpectQuery("UPDATE").WillReturnRows(rows)
		res, err := repo.Update(context.Background(), campaign)
		assert.Nil(t, err)
		assert.Equal(t, campaign, res)
	})
	t.Run("update promotion campaign error", func(t *testing.T) {
		mock.ExpectQuery("UPDATE").WillReturnError(assert.AnError)
		res, err := repo.Update(context.Background(), campaign)
		assert.Empty(t, res)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestPromotionCampaignRepository_Create(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewPromotionCampaignRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	metadata := entity.PromotionCampaignMetadata{
		Products: []entity.PromotionCampaignProduct{
			{
				Symbols:       []string{"HPG", "DGW"},
				LoanPackageId: 1,
				RetailSymbols: []string{"HPG", "DGW"},
			},
		},
	}
	metadataStr, err := json.Marshal(metadata)
	campaign := entity.PromotionCampaign{
		Id:          1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UpdatedBy:   "kiennt",
		Name:        "name",
		Tag:         "5.99*",
		Status:      entity.Active,
		Description: "description",
		Metadata:    metadata,
	}
	t.Run("create promotion campaign success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"promotion_campaign.id",
			"promotion_campaign.created_at",
			"promotion_campaign.updated_at",
			"promotion_campaign.updated_by",
			"promotion_campaign.name",
			"promotion_campaign.tag",
			"promotion_campaign.description",
			"promotion_campaign.status",
			"promotion_campaign.metadata",
		})
		rows.AddRow(
			campaign.Id,
			campaign.CreatedAt,
			campaign.UpdatedAt,
			campaign.UpdatedBy,
			campaign.Name,
			campaign.Tag,
			campaign.Description,
			campaign.Status,
			metadataStr,
		)
		mock.ExpectQuery("INSERT").WillReturnRows(rows)
		res, err := repo.Create(context.Background(), campaign)
		assert.Nil(t, err)
		assert.Equal(t, campaign, res)
	})
	t.Run("create promotion campaign error", func(t *testing.T) {
		mock.ExpectQuery("INSERT").WillReturnError(assert.AnError)
		res, err := repo.Create(context.Background(), campaign)
		assert.Empty(t, res)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestPromotionCampaignRepository_GetById(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewPromotionCampaignRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	metadata := entity.PromotionCampaignMetadata{
		Products: []entity.PromotionCampaignProduct{
			{
				Symbols:       []string{"HPG", "DGW"},
				LoanPackageId: 1,
				RetailSymbols: []string{"HPG", "DGW"},
			},
		},
	}
	metadataStr, err := json.Marshal(metadata)
	campaign := entity.PromotionCampaign{
		Id:          1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UpdatedBy:   "kiennt",
		Name:        "name",
		Tag:         "5.99*",
		Status:      entity.Active,
		Description: "description",
		Metadata:    metadata,
	}
	t.Run("get promotion campaign by id success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"promotion_campaign.id",
			"promotion_campaign.created_at",
			"promotion_campaign.updated_at",
			"promotion_campaign.updated_by",
			"promotion_campaign.name",
			"promotion_campaign.tag",
			"promotion_campaign.description",
			"promotion_campaign.status",
			"promotion_campaign.metadata",
		})
		rows.AddRow(
			campaign.Id,
			campaign.CreatedAt,
			campaign.UpdatedAt,
			campaign.UpdatedBy,
			campaign.Name,
			campaign.Tag,
			campaign.Description,
			campaign.Status,
			metadataStr,
		)
		mock.ExpectQuery("SELECT .* FROM public.promotion_campaign").WillReturnRows(rows)
		res, err := repo.GetById(context.Background(), campaign.Id)
		assert.Nil(t, err)
		assert.Equal(t, campaign, res)
	})
	t.Run("get promotion campaign by id error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM public.promotion_campaign").WillReturnError(assert.AnError)
		res, err := repo.GetById(context.Background(), campaign.Id)
		assert.Empty(t, res)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
