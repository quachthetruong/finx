package postgres

import (
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

func MapScoreGroupsDbToEntity(scoreGroups []model.ScoreGroup) []entity.ScoreGroup {
	dest := make([]entity.ScoreGroup, 0, len(scoreGroups))
	for _, scoreGroup := range scoreGroups {
		dest = append(dest, MapScoreGroupDbToEntity(scoreGroup))
	}
	return dest
}

func MapScoreGroupDbToEntity(scoreGroup model.ScoreGroup) entity.ScoreGroup {
	return entity.ScoreGroup{
		Id:        scoreGroup.ID,
		Code:      scoreGroup.Code,
		MinScore:  scoreGroup.MinScore,
		MaxScore:  scoreGroup.MaxScore,
		CreatedAt: scoreGroup.CreatedAt,
		UpdatedAt: scoreGroup.UpdatedAt,
	}
}

func MapScoreGroupEntityToDb(scoreGroup entity.ScoreGroup) model.ScoreGroup {
	return model.ScoreGroup{
		ID:        scoreGroup.Id,
		Code:      scoreGroup.Code,
		MinScore:  scoreGroup.MinScore,
		MaxScore:  scoreGroup.MaxScore,
		CreatedAt: scoreGroup.CreatedAt,
		UpdatedAt: scoreGroup.UpdatedAt,
	}
}
