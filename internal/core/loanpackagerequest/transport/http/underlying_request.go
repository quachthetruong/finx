package http

import (
	"time"

	"financing-offer/internal/core/entity"
	"financing-offer/pkg/optional"
	"github.com/shopspring/decimal"
)

type GetUnderlyingLoanPackageRequest struct {
	Symbols            []string        `form:"symbols"`
	Types              []string        `form:"types"`
	InvestorId         string          `form:"investorId"`
	Statuses           []string        `form:"statuses"`
	Ids                []int64         `form:"ids"`
	StartDate          time.Time       `form:"startDate"`
	EndDate            time.Time       `form:"endDate"`
	LoanPercentFrom    decimal.Decimal `form:"loanPercentFrom"`
	LoanPercentTo      decimal.Decimal `form:"loanPercentTo"`
	LimitAmountFrom    decimal.Decimal `form:"limitAmountFrom"`
	LimitAmountTo      decimal.Decimal `form:"limitAmountTo"`
	AccountNumbers     []string        `form:"accountNumbers"`
	CustodyCode        string          `form:"custodyCode"`
	CustodyCodes       []string        `form:"custodyCodes"`
	SubmissionStatuses []string        `form:"submissionStatuses"`
}

func (r GetUnderlyingLoanPackageRequest) toEntity() entity.UnderlyingLoanPackageFilter {
	types := make([]entity.LoanPackageRequestType, 0)
	for _, t := range r.Types {
		types = append(types, entity.LoanPackageRequestType(t))
	}
	statuses := make([]entity.LoanPackageRequestStatus, len(r.Statuses))
	for _, s := range r.Statuses {
		statuses = append(statuses, entity.LoanPackageRequestStatus(s))
	}
	submissonStatuses := make([]entity.SubmissionSheetStatus, len(r.SubmissionStatuses))
	for _, s := range r.SubmissionStatuses {
		submissonStatuses = append(submissonStatuses, entity.SubmissionSheetStatus(s))
	}
	return entity.UnderlyingLoanPackageFilter{
		Symbols:            r.Symbols,
		InvestorId:         optional.FromValueNonZero(r.InvestorId),
		Ids:                r.Ids,
		AccountNumbers:     r.AccountNumbers,
		StartDate:          optional.FromValueNonZero(r.StartDate),
		EndDate:            optional.FromValueNonZero(r.EndDate),
		LoanPercentFrom:    optional.FromValueNonZero(r.LoanPercentFrom),
		LoanPercentTo:      optional.FromValueNonZero(r.LoanPercentTo),
		LimitAmountFrom:    optional.FromValueNonZero(r.LimitAmountFrom),
		LimitAmountTo:      optional.FromValueNonZero(r.LimitAmountTo),
		Types:              types,
		Statuses:           statuses,
		CustodyCode:        optional.FromValueNonZero(r.CustodyCode),
		CustodyCodes:       r.CustodyCodes,
		SubmissionStatuses: submissonStatuses,
	}
}
