package entity

import (
	"time"
)

type SuggestedOffer struct {
	Id        int64                 `json:"id"`
	ConfigId  int64                 `json:"-"`
	Config    *SuggestedOfferConfig `json:"config"`
	AccountNo string                `json:"accountNo"`
	Symbols   []string              `json:"symbols"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
}
