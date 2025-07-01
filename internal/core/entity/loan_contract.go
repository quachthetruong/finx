package entity

import (
	"time"
)

type LoanContract struct {
	Id                   int64     `json:"id"`
	LoanOfferInterestId  int64     `json:"loanInterestId"`
	SymbolId             int64     `json:"symbolId"`
	InvestorId           string    `json:"investorId"`
	AccountNo            string    `json:"accountNo"`
	LoanId               int64     `json:"loanId"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	GuaranteedEndAt      time.Time `json:"guaranteedEndAt"`
	LoanPackageAccountId int64     `json:"loanPackageAccountId"`
	LoanProductIdRef     int64     `json:"loanProductIdRef"`
}
