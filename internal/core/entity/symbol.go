package entity

import (
	"time"

	"financing-offer/internal/core"
	"financing-offer/pkg/optional"
)

type Symbol struct {
	Id              int64        `json:"id"`
	StockExchangeId int64        `json:"stockExchangeId"`
	Symbol          string       `json:"symbol"`
	AssetType       AssetType    `json:"assetType"`
	Status          SymbolStatus `json:"status"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
	LastUpdatedBy   string       `json:"lastUpdatedBy"`

	Scores []SymbolScore `json:"scores"`
}

type SymbolFilter struct {
	core.Paging
	StockExchangeCode optional.Optional[string] `json:"stockExchangeCode"`
	AssetType         optional.Optional[string] `json:"assetType"`
}

type SymbolStatus string

const (
	SymbolStatusActive   SymbolStatus = "ACTIVE"
	SymbolStatusInactive SymbolStatus = "INACTIVE"
)

func (s SymbolStatus) String() string {
	return string(s)
}

func SymbolStatusFromString(s string) SymbolStatus {
	switch s {
	case "ACTIVE":
		return SymbolStatusActive
	case "INACTIVE":
		return SymbolStatusInactive
	default:
		return SymbolStatusActive
	}
}
