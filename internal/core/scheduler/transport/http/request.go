package http

import (
	"time"

	"github.com/shopspring/decimal"
)

type LoanPackageSchedulerConfig struct {
	MaximumLoanRate decimal.Decimal `json:"maximumLoanRate" binding:"required"`
	AffectedFrom    time.Time       `json:"affectedFrom" binding:"required"`
}
