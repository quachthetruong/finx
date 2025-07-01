package postgres

import (
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

func MapStockExchangesDbToEntity(stockExchanges []model.StockExchange) []entity.StockExchange {
	res := make([]entity.StockExchange, 0, len(stockExchanges))
	for _, v := range stockExchanges {
		res = append(res, MapStockExchangeDbToEntity(v))
	}
	return res
}

func MapStockExchangeDbToEntity(stockExchange model.StockExchange) entity.StockExchange {
	res := entity.StockExchange{
		Id:        stockExchange.ID,
		Code:      stockExchange.Code,
		MinScore:  stockExchange.MinScore,
		MaxScore:  stockExchange.MaxScore,
		CreatedAt: stockExchange.CreatedAt,
		UpdatedAt: stockExchange.UpdatedAt,
	}
	if stockExchange.ScoreGroupID != nil {
		res.ScoreGroupId = *stockExchange.ScoreGroupID
	}
	return res
}

func MapStockExchangeEntityToDb(stockExchange entity.StockExchange) model.StockExchange {
	return model.StockExchange{
		ID:           stockExchange.Id,
		Code:         stockExchange.Code,
		ScoreGroupID: &stockExchange.ScoreGroupId,
		MinScore:     stockExchange.MinScore,
		MaxScore:     stockExchange.MaxScore,
		CreatedAt:    stockExchange.CreatedAt,
		UpdatedAt:    stockExchange.UpdatedAt,
	}
}
