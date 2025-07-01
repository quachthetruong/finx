package postgres

import (
	"github.com/go-jet/jet/v2/postgres"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/pkg/querymod"
)

func MapSymbolScoreDbToEntity(symbolScore model.SymbolScore) entity.SymbolScore {
	return entity.SymbolScore{
		Id:           symbolScore.ID,
		SymbolId:     symbolScore.SymbolID,
		Score:        symbolScore.Score,
		AffectedFrom: symbolScore.AffectedFrom,
		Status:       entity.SymbolScoreStatusFromString(symbolScore.Status),
		Type:         entity.SymbolScoreTypeFromString(symbolScore.Type),
		Creator:      symbolScore.Creator,
		CreatedAt:    symbolScore.CreatedAt,
		UpdatedAt:    symbolScore.UpdatedAt,
	}
}

func MapSymbolScoresDbToEntity(symbolScores []model.SymbolScore) []entity.SymbolScore {
	res := make([]entity.SymbolScore, 0, len(symbolScores))
	for _, v := range symbolScores {
		res = append(res, MapSymbolScoreDbToEntity(v))
	}
	return res
}

func MapSymbolScoreEntityToDb(symbolScore entity.SymbolScore) model.SymbolScore {
	return model.SymbolScore{
		ID:           symbolScore.Id,
		SymbolID:     symbolScore.SymbolId,
		Score:        symbolScore.Score,
		AffectedFrom: symbolScore.AffectedFrom,
		Status:       symbolScore.Status.String(),
		Type:         symbolScore.Type.String(),
		Creator:      symbolScore.Creator,
		CreatedAt:    symbolScore.CreatedAt,
		UpdatedAt:    symbolScore.UpdatedAt,
	}
}

func ApplyFilter(filter entity.SymbolScoreFilter) postgres.BoolExpression {
	expr := postgres.Bool(true)
	if len(filter.Symbols) > 0 {
		expr = expr.AND(table.Symbol.Symbol.IN(querymod.In(filter.Symbols)...))
	}
	if filter.Status.IsPresent() {
		expr = expr.AND(table.SymbolScore.Status.EQ(postgres.String(filter.Status.Get().String())))
	}
	if filter.Type.IsPresent() {
		expr = expr.AND(table.SymbolScore.Type.EQ(postgres.String(filter.Type.Get().String())))
	}
	return expr
}
