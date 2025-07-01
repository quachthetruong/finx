package loanofferinterest

import (
	"github.com/shopspring/decimal"
)

type AssignedLoanPackageAccount struct {
	LoanPackageOfferInterestId  int64
	LoanPackageId               int64
	CreatedLoanPackageAccountId int64
	InterestRate                decimal.Decimal
	LoanRate                    decimal.Decimal
	InitialRate                 decimal.Decimal
}
