package entity

import (
	"time"
)

type ScoreGroup struct {
	Id        int64     `json:"id"`
	Code      string    `json:"code"`
	MinScore  int32     `json:"minScore" binding:"lte=100,gte=0"`
	MaxScore  int32     `json:"maxScore" binding:"lte=100,gte=0"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
