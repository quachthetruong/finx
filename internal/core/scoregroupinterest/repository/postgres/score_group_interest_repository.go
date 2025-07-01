package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/go-jet/jet/v2/postgres"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/scoregroupinterest/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

const (
	finalSymbolScoreGroupColumnAlias = "final_symbol_score_group"
	finalSymbolScoreGroupTableAlias  = "sg"
)

var _ repository.ScoreGroupInterestRepository = (*ScoreGroupInterestSqlRepository)(nil)

type ScoreGroupInterestSqlRepository struct {
	getDbFunc database.GetDbFunc
}

func NewScoreGroupInterestSqlRepository(getDbFunc database.GetDbFunc) *ScoreGroupInterestSqlRepository {
	return &ScoreGroupInterestSqlRepository{getDbFunc: getDbFunc}
}

func (s *ScoreGroupInterestSqlRepository) GetAll(ctx context.Context, filter entity.ScoreGroupInterestFilter) ([]entity.ScoreGroupInterest, error) {
	dest := make([]model.ScoreGroupInterest, 0)
	stm := table.ScoreGroupInterest.SELECT(table.ScoreGroupInterest.AllColumns)
	if filter.Score.IsPresent() {
		stm = stm.FROM(
			table.ScoreGroup.
				INNER_JOIN(table.ScoreGroupInterest, table.ScoreGroup.ID.EQ(table.ScoreGroupInterest.ScoreGroupID)),
		)
	}
	stm = stm.WHERE(ApplyFilter(filter))
	if err := stm.QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		return nil, fmt.Errorf("get all %w", err)
	}
	roles := make([]entity.ScoreGroupInterest, 0)
	for _, v := range dest {
		roles = append(roles, MapScoreGroupInterestDbToEntity(v))
	}
	return roles, nil
}

func (s *ScoreGroupInterestSqlRepository) GetAvailablePackageBySymbolId(ctx context.Context, symbolId int64) ([]entity.ScoreGroupInterest, error) {
	dest := make([]model.ScoreGroupInterest, 0)
	stm := table.ScoreGroupInterest.SELECT(table.ScoreGroupInterest.AllColumns)
	finalSymbolScoreGroup := table.SymbolScore.SELECT(
		table.ScoreGroup.ID.AS(finalSymbolScoreGroupColumnAlias),
	).
		FROM(
			table.SymbolScore.INNER_JOIN(table.Symbol, table.Symbol.ID.EQ(table.SymbolScore.SymbolID)).
				INNER_JOIN(table.StockExchange, table.StockExchange.ID.EQ(table.Symbol.StockExchangeID)).
				INNER_JOIN(
					table.ScoreGroup, table.ScoreGroup.ID.EQ(table.StockExchange.ScoreGroupID).
						OR(table.SymbolScore.Score.BETWEEN(table.ScoreGroup.MinScore, table.ScoreGroup.MaxScore)),
				),
		).
		WHERE(
			table.SymbolScore.Status.EQ(postgres.String(entity.SymbolScoreStatusActive.String())).
				AND(table.SymbolScore.AffectedFrom.LT_EQ(postgres.TimestampT(time.Now()))).
				AND(table.Symbol.ID.EQ(postgres.Int64(symbolId))),
		).
		ORDER_BY(table.SymbolScore.AffectedFrom.DESC(), table.ScoreGroup.MaxScore.ASC()).
		LIMIT(1).
		AsTable(finalSymbolScoreGroupTableAlias)
	symbolScoreGroupColumn := postgres.IntegerColumn(finalSymbolScoreGroupColumnAlias).From(finalSymbolScoreGroup)
	stm = stm.FROM(
		finalSymbolScoreGroup.
			INNER_JOIN(table.ScoreGroupInterest, table.ScoreGroupInterest.ScoreGroupID.EQ(symbolScoreGroupColumn)),
	)
	if err := stm.QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		return nil, fmt.Errorf("ScoreGroupInterestSqlRepository GetAvailablePackageBySymbolId %w", err)
	}
	packages := make([]entity.ScoreGroupInterest, 0)
	for _, v := range dest {
		packages = append(packages, MapScoreGroupInterestDbToEntity(v))
	}
	return packages, nil
}

func (s *ScoreGroupInterestSqlRepository) GetById(ctx context.Context, id int64) (entity.ScoreGroupInterest, error) {
	db := s.getDbFunc(ctx)
	var role model.ScoreGroupInterest
	if err := table.ScoreGroupInterest.
		SELECT(table.ScoreGroupInterest.AllColumns).
		WHERE(table.ScoreGroupInterest.ID.EQ(postgres.Int64(id))).
		QueryContext(ctx, db, &role); err != nil {
		return entity.ScoreGroupInterest{}, fmt.Errorf("ScoreGroupInterestSqlRepository GetById %w", err)
	}
	return MapScoreGroupInterestDbToEntity(role), nil
}

func (s *ScoreGroupInterestSqlRepository) Create(ctx context.Context, symbol entity.ScoreGroupInterest) (entity.ScoreGroupInterest, error) {
	db := s.getDbFunc(ctx)
	groupRoleModel := MapScoreGroupInterestEntityToDb(symbol)
	created := model.ScoreGroupInterest{}
	if err := table.ScoreGroupInterest.
		INSERT(
			table.ScoreGroupInterest.MutableColumns.Except(
				table.ScoreGroupInterest.CreatedAt, table.ScoreGroupInterest.UpdatedAt,
				table.ScoreGroupInterest.LoanRate, table.ScoreGroupInterest.InterestRate,
			),
		).
		MODEL(groupRoleModel).
		RETURNING(table.ScoreGroupInterest.AllColumns).
		QueryContext(ctx, db, &created); err != nil {
		return entity.ScoreGroupInterest{}, fmt.Errorf("ScoreGroupInterestSqlRepository Create %w", err)
	}
	return MapScoreGroupInterestDbToEntity(created), nil
}

func (s *ScoreGroupInterestSqlRepository) Update(ctx context.Context, symbol entity.ScoreGroupInterest) (entity.ScoreGroupInterest, error) {
	db := s.getDbFunc(ctx)
	groupRoleModel := MapScoreGroupInterestEntityToDb(symbol)
	updated := model.ScoreGroupInterest{}
	if err := table.ScoreGroupInterest.
		UPDATE(
			table.ScoreGroupInterest.MutableColumns.
				Except(table.ScoreGroupInterest.CreatedAt, table.ScoreGroupInterest.UpdatedAt),
		).
		MODEL(groupRoleModel).
		WHERE(table.ScoreGroupInterest.ID.EQ(postgres.Int64(symbol.Id))).
		RETURNING(table.ScoreGroupInterest.AllColumns).
		QueryContext(ctx, db, &updated); err != nil {
		return entity.ScoreGroupInterest{}, fmt.Errorf("ScoreGroupInterestSqlRepository Update %w", err)
	}
	return MapScoreGroupInterestDbToEntity(updated), nil
}

func (s *ScoreGroupInterestSqlRepository) Delete(ctx context.Context, id int64) (bool, error) {
	db := s.getDbFunc(ctx)
	res, err := table.ScoreGroupInterest.DELETE().WHERE(table.ScoreGroupInterest.ID.EQ(postgres.Int64(id))).ExecContext(
		ctx, db,
	)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil && count == 1 {
			return true, nil
		}
	}
	return false, apperrors.ErrorDeleteScoreGroupInterest
}

func (s *ScoreGroupInterestSqlRepository) GetAvailableScoreInterestsByScoreGroupId(ctx context.Context, scoreGroupId int64) ([]entity.ScoreGroupInterest, error) {
	var dest []model.ScoreGroupInterest
	if err := table.ScoreGroupInterest.
		SELECT(table.ScoreGroupInterest.AllColumns).
		FROM(
			table.ScoreGroupInterest.INNER_JOIN(table.ScoreGroup, table.ScoreGroupInterest.ScoreGroupID.EQ(table.ScoreGroup.ID)),
		).
		WHERE(table.ScoreGroup.ID.EQ(postgres.Int64(scoreGroupId))).
		QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		return []entity.ScoreGroupInterest{}, fmt.Errorf("ScoreGroupInterestSqlRepository GetAvailableScoreInterestsByScoreGroupId %w", err)
	}
	return MapScoreGroupInterestsDbToEntity(dest), nil
}
