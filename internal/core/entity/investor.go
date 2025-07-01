package entity

import (
	"time"
)

type Investor struct {
	InvestorId  string    `json:"investorId"`
	CustodyCode string    `json:"custodyCode"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
