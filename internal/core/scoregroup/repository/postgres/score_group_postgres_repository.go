package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/scoregroup/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

var _ repository.ScoreGroupRepository = (*ScoreGroupRepository)(nil)

type ScoreGroupRepository struct {
	getDbFunc database.GetDbFunc
}

func (s *ScoreGroupRepository) GetById(ctx context.Context, id int64) (entity.ScoreGroup, error) {
	var dest model.ScoreGroup
	if err := table.ScoreGroup.
		SELECT(table.ScoreGroup.AllColumns).
		WHERE(table.ScoreGroup.ID.EQ(postgres.Int64(id))).
		QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		return entity.ScoreGroup{}, fmt.Errorf("ScoreGroupRepository GetById %w", err)
	}
	return MapScoreGroupDbToEntity(dest), nil
}

func NewScoreGroupRepository(getDbFunc database.GetDbFunc) *ScoreGroupRepository {
	return &ScoreGroupRepository{
		getDbFunc: getDbFunc,
	}
}

func (s *ScoreGroupRepository) Delete(ctx context.Context, id int64) error {
	if _, err := table.ScoreGroup.DELETE().WHERE(table.ScoreGroup.ID.EQ(postgres.Int64(id))).ExecContext(
		ctx, s.getDbFunc(ctx),
	); err != nil {
		return fmt.Errorf("ScoreGroupRepository Delete %w", err)
	}
	return nil
}

func (s *ScoreGroupRepository) GetAll(ctx context.Context) ([]entity.ScoreGroup, error) {
	dest := make([]model.ScoreGroup, 0)
	if err := table.ScoreGroup.SELECT(table.ScoreGroup.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []entity.ScoreGroup{}, nil
		}
		return []entity.ScoreGroup{}, fmt.Errorf("ScoreGroupRepository GetAll %w", err)
	}
	return MapScoreGroupsDbToEntity(dest), nil
}

func (s *ScoreGroupRepository) Create(ctx context.Context, scoreGroup entity.ScoreGroup) (entity.ScoreGroup, error) {
	createModel := MapScoreGroupEntityToDb(scoreGroup)
	created := model.ScoreGroup{}
	if err := table.ScoreGroup.INSERT(table.ScoreGroup.MutableColumns).
		MODEL(createModel).
		RETURNING(table.ScoreGroup.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &created); err != nil {
		return entity.ScoreGroup{}, fmt.Errorf("ScoreGroupRepository Create %w", err)
	}
	return MapScoreGroupDbToEntity(created), nil
}

func (s *ScoreGroupRepository) Update(ctx context.Context, scoreGroup entity.ScoreGroup) (entity.ScoreGroup, error) {
	updateModel := MapScoreGroupEntityToDb(scoreGroup)
	updated := model.ScoreGroup{}
	if err := table.ScoreGroup.UPDATE(table.ScoreGroup.MutableColumns).
		MODEL(updateModel).
		WHERE(table.ScoreGroup.ID.EQ(postgres.Int64(scoreGroup.Id))).
		RETURNING(table.ScoreGroup.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &updated); err != nil {
		return entity.ScoreGroup{}, fmt.Errorf("ScoreGroupRepository Update %w", err)
	}
	return MapScoreGroupDbToEntity(updated), nil
}
