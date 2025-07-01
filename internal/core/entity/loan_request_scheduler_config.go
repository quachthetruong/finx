package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type LoanRequestSchedulerConfig struct {
	ID              int64           `json:"id"`
	MaximumLoanRate decimal.Decimal `json:"maximumLoanRate"`
	AffectedFrom    time.Time       `json:"affectedFrom"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}
