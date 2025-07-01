package entity

import (
	"time"

	"financing-offer/pkg/optional"
)

type BlacklistSymbol struct {
	Id           int64                 `json:"id"`
	SymbolId     int64                 `json:"symbolId"`
	AffectedFrom time.Time             `json:"affectedFrom"`
	AffectedTo   time.Time             `json:"affectedTo"`
	Status       BlacklistSymbolStatus `json:"status"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
}

type BlacklistSymbolFilter struct {
	Symbol optional.Optional[string] `json:"symbol"`
}

type BlacklistSymbolStatus string

const (
	BlacklistSymbolStatusActive   BlacklistSymbolStatus = "ACTIVE"
	BlacklistSymbolStatusInactive BlacklistSymbolStatus = "INACTIVE"
)

func BlacklistSymbolStatusFromString(s string) BlacklistSymbolStatus {
	switch s {
	case string(BlacklistSymbolStatusActive):
		return BlacklistSymbolStatusActive
	case string(BlacklistSymbolStatusInactive):
		return BlacklistSymbolStatusInactive
	default:
		return BlacklistSymbolStatusInactive
	}
}

func (b BlacklistSymbolStatus) String() string {
	return string(b)
}
