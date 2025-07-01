package postgres

import (
	"github.com/go-jet/jet/v2/postgres"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/pkg/querymod"
)

func MapScoreGroupInterestDbToEntity(score model.ScoreGroupInterest) entity.ScoreGroupInterest {
	return entity.ScoreGroupInterest{
		Id:           score.ID,
		LimitAmount:  score.LimitAmount,
		LoanRate:     score.LoanRate,
		InterestRate: score.InterestRate,
		ScoreGroupId: score.ScoreGroupID,
		CreatedAt:    score.CreatedAt,
		UpdatedAt:    score.UpdatedAt,
	}
}

func MapScoreGroupInterestEntityToDb(score entity.ScoreGroupInterest) model.ScoreGroupInterest {
	return model.ScoreGroupInterest{
		ID:           score.Id,
		LimitAmount:  score.LimitAmount,
		LoanRate:     score.LoanRate,
		InterestRate: score.InterestRate,
		ScoreGroupID: score.ScoreGroupId,
		CreatedAt:    score.CreatedAt,
		UpdatedAt:    score.UpdatedAt,
	}
}

func MapScoreGroupInterestsDbToEntity(scoreGroupInterests []model.ScoreGroupInterest) []entity.ScoreGroupInterest {
	dest := make([]entity.ScoreGroupInterest, 0, len(scoreGroupInterests))
	for _, scoreGroupInterest := range scoreGroupInterests {
		dest = append(dest, MapScoreGroupInterestDbToEntity(scoreGroupInterest))
	}
	return dest
}

func ApplyFilter(filter entity.ScoreGroupInterestFilter) postgres.BoolExpression {
	expr := postgres.Bool(true)
	if filter.Score.IsPresent() {
		expr = expr.AND(table.ScoreGroup.MinScore.LT_EQ(postgres.Int32(filter.Score.Get()))).
			AND(table.ScoreGroup.MaxScore.GT_EQ(postgres.Int32(filter.Score.Get())))
	}
	if len(filter.Ids) > 0 {
		expr = expr.AND(table.ScoreGroupInterest.ID.IN(querymod.In(filter.Ids)...))
	}
	return expr
}
