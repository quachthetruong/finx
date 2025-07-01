package entity

import (
	"time"
)

type StockExchange struct {
	Id           int64     `json:"id"`
	Code         string    `json:"code"`
	ScoreGroupId int64     `json:"scoreGroupId"`
	MinScore     int32     `json:"minScore"`
	MaxScore     int32     `json:"maxScore"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
