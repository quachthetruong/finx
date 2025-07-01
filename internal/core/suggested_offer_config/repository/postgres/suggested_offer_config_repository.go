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
	"financing-offer/pkg/optional"
)

type SuggestedOfferConfigPostgresRepository struct {
	getDbFunc database.GetDbFunc
}

func NewSuggestedOfferConfigRepository(getDbFunc database.GetDbFunc) *SuggestedOfferConfigPostgresRepository {
	return &SuggestedOfferConfigPostgresRepository{getDbFunc: getDbFunc}
}

func (s *SuggestedOfferConfigPostgresRepository) GetAll(ctx context.Context) ([]entity.SuggestedOfferConfig, error) {
	dest := make([]model.SuggestedOfferConfig, 0)
	if err := table.SuggestedOfferConfig.
		SELECT(table.SuggestedOfferConfig.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		return nil, fmt.Errorf("SuggestedOfferConfigRepository GetAll %w", err)
	}
	return MapSuggestedOfferConfigsDbToEntity(dest), nil
}

func (s *SuggestedOfferConfigPostgresRepository) GetActiveSuggestedOfferConfig(ctx context.Context) (optional.Optional[entity.SuggestedOfferConfig], error) {
	dest := model.SuggestedOfferConfig{}
	if err := table.SuggestedOfferConfig.
		SELECT(table.SuggestedOfferConfig.AllColumns).
		WHERE(table.SuggestedOfferConfig.Status.EQ(postgres.String(entity.SuggestedOfferConfigStatusActive.String()))).
		QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return optional.None[entity.SuggestedOfferConfig](), nil
		}
		return optional.None[entity.SuggestedOfferConfig](), fmt.Errorf("SuggestedOfferConfigRepository GetActiveSuggestedOfferConfig %w", err)
	}
	return optional.Some(MapSuggestedOfferConfigDbToEntity(dest)), nil
}

func (s *SuggestedOfferConfigPostgresRepository) UpdateStatus(ctx context.Context, id int64, status entity.SuggestedOfferConfigStatus, updater string) (entity.SuggestedOfferConfig, error) {
	dest := model.SuggestedOfferConfig{}
	if err := table.SuggestedOfferConfig.
		UPDATE().
		SET(
			table.SuggestedOfferConfig.Status.SET(postgres.String(status.String())),
			table.SuggestedOfferConfig.LastUpdatedBy.SET(postgres.String(updater)),
		).
		WHERE(table.SuggestedOfferConfig.ID.EQ(postgres.Int64(id))).
		RETURNING(table.SuggestedOfferConfig.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		return entity.SuggestedOfferConfig{}, fmt.Errorf("SuggestedOfferConfigRepository UpdateStatus %w", err)
	}
	return MapSuggestedOfferConfigDbToEntity(dest), nil
}

func (s *SuggestedOfferConfigPostgresRepository) Update(ctx context.Context, suggestedOfferConfig entity.SuggestedOfferConfig) (entity.SuggestedOfferConfig, error) {
	updatedModel := MapSuggestedOfferConfigEntityToDb(suggestedOfferConfig)
	dest := model.SuggestedOfferConfig{}
	if err := table.SuggestedOfferConfig.
		UPDATE(table.SuggestedOfferConfig.MutableColumns.Except(table.SuggestedOfferConfig.Status)).
		MODEL(updatedModel).
		WHERE(table.SuggestedOfferConfig.ID.EQ(postgres.Int64(suggestedOfferConfig.Id))).
		RETURNING(table.SuggestedOfferConfig.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		return entity.SuggestedOfferConfig{}, fmt.Errorf("SuggestedOfferConfigRepository Update %w", err)
	}
	return MapSuggestedOfferConfigDbToEntity(dest), nil
}

func (s *SuggestedOfferConfigPostgresRepository) Create(ctx context.Context, suggestedOfferConfig entity.SuggestedOfferConfig) (entity.SuggestedOfferConfig, error) {
	dest := model.SuggestedOfferConfig{}
	toCreate := MapSuggestedOfferConfigEntityToDb(suggestedOfferConfig)
	if err := table.SuggestedOfferConfig.
		INSERT(table.SuggestedOfferConfig.MutableColumns).
		MODEL(toCreate).
		RETURNING(table.SuggestedOfferConfig.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		return entity.SuggestedOfferConfig{}, fmt.Errorf("SuggestedOfferConfigRepository Create %w", err)
	}
	return MapSuggestedOfferConfigDbToEntity(dest), nil
}

func (s *SuggestedOfferConfigPostgresRepository) GetById(ctx context.Context, id int64) (entity.SuggestedOfferConfig, error) {
	dest := model.SuggestedOfferConfig{}
	if err := table.SuggestedOfferConfig.
		SELECT(table.SuggestedOfferConfig.AllColumns).
		WHERE(table.SuggestedOfferConfig.ID.EQ(postgres.Int64(id))).
		QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		return entity.SuggestedOfferConfig{}, fmt.Errorf("SuggestedOfferConfigRepository GetById %w", err)
	}
	return MapSuggestedOfferConfigDbToEntity(dest), nil
}
