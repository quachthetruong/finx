package http

import (
	"time"

	"financing-offer/internal/core/entity"
)

type BlacklistSymbolRequest struct {
	AffectedFrom time.Time                    `json:"affectedFrom"`
	AffectedTo   time.Time                    `json:"affectedTo"`
	Status       entity.BlacklistSymbolStatus `json:"status"`
}
