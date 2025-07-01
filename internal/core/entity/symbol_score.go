package entity

import (
	"time"

	"financing-offer/pkg/optional"
)

type SymbolScore struct {
	Id           int64             `json:"id"`
	SymbolId     int64             `json:"symbolId"`
	Score        int32             `json:"score" binding:"lte=100,gte=0"`
	AffectedFrom time.Time         `json:"affectedFrom"`
	Status       SymbolScoreStatus `json:"status" binding:"oneof=ACTIVE INACTIVE"`
	Type         SymbolScoreType   `json:"type" binding:"oneof=MANUAL SYSTEM"`
	Creator      string            `json:"creator"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
}

type SymbolScoreFilter struct {
	Symbols []string                             `json:"symbol"`
	Status  optional.Optional[SymbolScoreStatus] `json:"status"`
	Type    optional.Optional[SymbolScoreType]   `json:"type"`
}

type SymbolScoreStatus string

const (
	SymbolScoreStatusActive   SymbolScoreStatus = "ACTIVE"
	SymbolScoreStatusInactive SymbolScoreStatus = "INACTIVE"
)

func (s SymbolScoreStatus) String() string {
	return string(s)
}

func SymbolScoreStatusFromString(s string) SymbolScoreStatus {
	switch s {
	case string(SymbolScoreStatusActive):
		return SymbolScoreStatusActive
	case string(SymbolScoreStatusInactive):
		return SymbolScoreStatusInactive
	default:
		return SymbolScoreStatusActive
	}
}

type SymbolScoreType string

const (
	SymbolScoreTypeManual SymbolScoreType = "MANUAL"
	SymbolScoreTypeSystem SymbolScoreType = "SYSTEM"
)

func (s SymbolScoreType) String() string {
	return string(s)
}

func SymbolScoreTypeFromString(s string) SymbolScoreType {
	switch s {
	case string(SymbolScoreTypeManual):
		return SymbolScoreTypeManual
	case string(SymbolScoreTypeSystem):
		return SymbolScoreTypeSystem
	default:
		return SymbolScoreTypeManual
	}
}
