package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/symbolscore/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

var _ repository.SymbolScoreRepository = (*SymbolScoreRepository)(nil)

type SymbolScoreRepository struct {
	getDbFunc database.GetDbFunc
}

func (s *SymbolScoreRepository) GetById(ctx context.Context, id int64) (entity.SymbolScore, error) {
	var dest model.SymbolScore
	if err := table.SymbolScore.
		SELECT(table.SymbolScore.AllColumns).
		WHERE(table.SymbolScore.ID.EQ(postgres.Int64(id))).
		QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		return entity.SymbolScore{}, fmt.Errorf("SymbolScoreRepository GetById %w", err)
	}
	return MapSymbolScoreDbToEntity(dest), nil
}

func (s *SymbolScoreRepository) GetCurrentScoreForSymbol(ctx context.Context, symbolId int64) (entity.SymbolScore, error) {
	var dest []model.SymbolScore
	if err := table.SymbolScore.
		SELECT(table.SymbolScore.AllColumns).
		WHERE(
			table.SymbolScore.SymbolID.EQ(postgres.Int64(symbolId)).
				AND(table.SymbolScore.Status.EQ(postgres.String(entity.SymbolScoreStatusActive.String()))),
		).QueryContext(
		ctx, s.getDbFunc(ctx), &dest,
	); err != nil {
		return entity.SymbolScore{}, fmt.Errorf("SymbolScoreRepository GetCurrentScoreForSymbol %w", err)
	}
	var res model.SymbolScore
	for _, v := range dest {
		if v.AffectedFrom.After(res.AffectedFrom) {
			res = v
		}
	}
	return MapSymbolScoreDbToEntity(res), nil
}

func NewSymbolScoreRepository(getDbFunc database.GetDbFunc) *SymbolScoreRepository {
	return &SymbolScoreRepository{getDbFunc: getDbFunc}
}

func (s *SymbolScoreRepository) GetAll(ctx context.Context, filter entity.SymbolScoreFilter) ([]entity.SymbolScore, error) {
	stm := table.SymbolScore.SELECT(table.SymbolScore.AllColumns)
	if len(filter.Symbols) > 0 {
		stm = stm.FROM(
			table.SymbolScore.INNER_JOIN(
				table.Symbol, table.Symbol.ID.EQ(table.SymbolScore.SymbolID),
			),
		)
	}
	dest := make([]model.SymbolScore, 0)
	if err := stm.WHERE(ApplyFilter(filter)).QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []entity.SymbolScore{}, nil
		}
		return []entity.SymbolScore{}, fmt.Errorf("SymbolScoreRepository GetAll %w", err)
	}
	return MapSymbolScoresDbToEntity(dest), nil
}

func (s *SymbolScoreRepository) Update(ctx context.Context, symbolScore entity.SymbolScore) (entity.SymbolScore, error) {
	updateModel := MapSymbolScoreEntityToDb(symbolScore)
	updated := model.SymbolScore{}
	if err := table.SymbolScore.UPDATE(table.SymbolScore.MutableColumns).
		MODEL(updateModel).
		WHERE(table.SymbolScore.ID.EQ(postgres.Int64(symbolScore.Id))).
		RETURNING(table.SymbolScore.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &updated); err != nil {
		return entity.SymbolScore{}, fmt.Errorf("SymbolScoreRepository Update %w", err)
	}
	return MapSymbolScoreDbToEntity(updated), nil
}

func (s *SymbolScoreRepository) Create(ctx context.Context, symbolScore entity.SymbolScore) (entity.SymbolScore, error) {
	createModel := MapSymbolScoreEntityToDb(symbolScore)
	created := model.SymbolScore{}
	if err := table.SymbolScore.INSERT(table.SymbolScore.MutableColumns).
		MODEL(createModel).
		RETURNING(table.SymbolScore.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &created); err != nil {
		return entity.SymbolScore{}, fmt.Errorf("SymbolScoreRepository Create %w", err)
	}
	return MapSymbolScoreDbToEntity(created), nil
}
