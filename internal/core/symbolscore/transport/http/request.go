package http

import (
	"time"

	"financing-offer/internal/core/entity"
)

type CreateSymbolScoreRequest struct {
	SymbolId     int64                  `json:"symbolId" binding:"gte=0"`
	Score        int32                  `json:"score" binding:"gte=0,lte=100"`
	AffectedFrom time.Time              `json:"affectedFrom" binding:"required"`
	Type         entity.SymbolScoreType `json:"type" binding:"required,oneof=MANUAL SYSTEM"`
}

func (r CreateSymbolScoreRequest) toEntity(id int64) entity.SymbolScore {
	return entity.SymbolScore{
		Id:           id,
		SymbolId:     r.SymbolId,
		Score:        r.Score,
		AffectedFrom: r.AffectedFrom,
		Status:       entity.SymbolScoreStatusActive,
		Type:         r.Type,
	}
}
