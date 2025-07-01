package http

import "financing-offer/internal/core/entity"

type ScoreGroupRequest struct {
	Code     string `json:"code" binding:"required"`
	MinScore int32  `json:"minScore" binding:"lte=100,gte=0"`
	MaxScore int32  `json:"maxScore" binding:"lte=100,gte=0"`
}

func (r ScoreGroupRequest) toEntity(id int64) entity.ScoreGroup {
	return entity.ScoreGroup{
		Id:       id,
		Code:     r.Code,
		MinScore: r.MinScore,
		MaxScore: r.MaxScore,
	}
}
