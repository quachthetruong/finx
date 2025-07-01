package entity

import (
	"time"

	"github.com/shopspring/decimal"

	"financing-offer/pkg/optional"
)

type ScoreGroupInterest struct {
	Id           int64           `json:"id"`
	LimitAmount  decimal.Decimal `json:"limitAmount"`
	LoanRate     decimal.Decimal `json:"loanRate"`
	InterestRate decimal.Decimal `json:"interestRate"`
	ScoreGroupId int64           `json:"scoreGroupId"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

type ScoreGroupInterestFilter struct {
	Ids   []int64
	Score optional.Optional[int32]
}
