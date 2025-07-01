package http

import "financing-offer/internal/core/entity"

type StockExchangeRequest struct {
	Code         string `json:"code" binding:"required"`
	ScoreGroupId int64  `json:"scoreGroupId" binding:"required"`
}

func (r StockExchangeRequest) toEntity(id int64) entity.StockExchange {
	return entity.StockExchange{
		Id:           id,
		Code:         r.Code,
		ScoreGroupId: r.ScoreGroupId,
	}
}
