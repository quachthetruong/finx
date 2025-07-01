package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

type PromotionCampaignPostgresRepository struct {
	getDbFunc database.GetDbFunc
}

func NewPromotionCampaignRepository(getDbFunc database.GetDbFunc) *PromotionCampaignPostgresRepository {
	return &PromotionCampaignPostgresRepository{getDbFunc: getDbFunc}
}

func (r *PromotionCampaignPostgresRepository) GetAll(ctx context.Context, filter entity.GetPromotionCampaignsRequest) ([]entity.PromotionCampaign, error) {
	errTemplate := "PromotionCampaignPostgresRepository GetAll %w"
	stm := table.PromotionCampaign.SELECT(table.PromotionCampaign.AllColumns)
	if filter.Status != "" {
		stm = stm.WHERE(table.PromotionCampaign.Status.EQ(postgres.String(filter.Status)))
	}
	dest := make([]model.PromotionCampaign, 0)
	if err := stm.QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []entity.PromotionCampaign{}, nil
		}
		return nil, fmt.Errorf(errTemplate, err)
	}
	campaigns := make([]entity.PromotionCampaign, 0, len(dest))
	for i := range dest {
		c, err := MapPromotionCampaignDbToEntity(dest[i])
		if err != nil {
			return []entity.PromotionCampaign{}, fmt.Errorf(errTemplate, err)
		}
		campaigns = append(campaigns, c)
	}
	return campaigns, nil
}

func (r *PromotionCampaignPostgresRepository) Create(ctx context.Context, campaign entity.PromotionCampaign) (entity.PromotionCampaign, error) {
	errTemplate := "PromotionCampaignPostgresRepository Create %w"
	createModel, err := MapPromotionCampaignEntityToDb(campaign)
	if err != nil {
		return entity.PromotionCampaign{}, fmt.Errorf(errTemplate, err)
	}
	created := model.PromotionCampaign{}
	if err := table.PromotionCampaign.INSERT(table.PromotionCampaign.MutableColumns).
		MODEL(createModel).
		RETURNING(table.PromotionCampaign.AllColumns).
		QueryContext(ctx, r.getDbFunc(ctx), &created); err != nil {
		return entity.PromotionCampaign{}, fmt.Errorf(errTemplate, err)
	}
	result, err := MapPromotionCampaignDbToEntity(created)
	if err != nil {
		return entity.PromotionCampaign{}, fmt.Errorf(errTemplate, err)
	}
	return result, nil
}

func (r *PromotionCampaignPostgresRepository) Update(ctx context.Context, campaign entity.PromotionCampaign) (entity.PromotionCampaign, error) {
	errTemplate := "PromotionCampaignPostgresRepository Update %w"
	updateModel, err := MapPromotionCampaignEntityToDb(campaign)
	if err != nil {
		return entity.PromotionCampaign{}, fmt.Errorf(errTemplate, err)
	}
	updated := model.PromotionCampaign{}
	if err := table.PromotionCampaign.UPDATE(table.PromotionCampaign.MutableColumns).
		MODEL(updateModel).
		WHERE(table.PromotionCampaign.ID.EQ(postgres.Int64(updateModel.ID))).
		RETURNING(table.PromotionCampaign.AllColumns).
		QueryContext(ctx, r.getDbFunc(ctx), &updated); err != nil {
		return entity.PromotionCampaign{}, fmt.Errorf(errTemplate, err)
	}
	result, err := MapPromotionCampaignDbToEntity(updated)
	if err != nil {
		return entity.PromotionCampaign{}, fmt.Errorf(errTemplate, err)
	}
	return result, nil
}

func (r *PromotionCampaignPostgresRepository) GetById(ctx context.Context, id int64) (entity.PromotionCampaign, error) {
	errTemplate := "PromotionCampaignPostgresRepository GetById %w"
	var res model.PromotionCampaign
	if err := table.PromotionCampaign.SELECT(table.PromotionCampaign.AllColumns).WHERE(table.PromotionCampaign.ID.EQ(postgres.Int64(id))).QueryContext(
		ctx, r.getDbFunc(ctx), &res,
	); err != nil {
		return entity.PromotionCampaign{}, fmt.Errorf(errTemplate, err)
	}
	result, err := MapPromotionCampaignDbToEntity(res)
	if err != nil {
		return entity.PromotionCampaign{}, fmt.Errorf(errTemplate, err)
	}
	return result, nil
}
