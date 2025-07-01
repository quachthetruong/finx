package entity

import (
	"time"

	"financing-offer/internal/core"
	"financing-offer/pkg/optional"
)

type CombinedLoanRequest struct {
	LoanRequest                 LoanPackageRequest        `json:"loanRequest"`
	LoanOffer                   *LoanPackageOffer         `json:"loanOffer,omitempty"`
	Symbol                      Symbol                    `json:"symbol"`
	LoanContract                *LoanContract             `json:"loanContract,omitempty"`
	Investor                    Investor                  `json:"investor"`
	AdminAssignedLoanPackageIds string                    `json:"adminAssignedLoanPackageIds"`
	ActivatedLoanPackageIds     string                    `json:"activatedLoanPackageIds"`
	PackageCreatedTime          *time.Time                `json:"packageCreatedTime,omitempty"`
	Status                      CombinedLoanRequestStatus `json:"status"`
	CancelledReason             string                    `json:"cancelledReason"`
}

type CombinedLoanRequestFilter struct {
	core.Paging
	Symbols         []string
	StartDate       optional.Optional[time.Time]
	EndDate         optional.Optional[time.Time]
	OfferDateFrom   optional.Optional[time.Time]
	OfferDateTo     optional.Optional[time.Time]
	FlowTypes       []string
	AccountNumbers  []string
	InvestorId      optional.Optional[string]
	Status          CombinedLoanRequestStatus
	AssignedLoanId  optional.Optional[int64]
	ActivatedLoanId optional.Optional[int64]
	Ids             []int64
	AssetType       optional.Optional[string]
	CustodyCode     optional.Optional[string]
}

type CombinedLoanRequestStatus string

const (
	CombinedLoanRequestStatusAwaitingOffer   CombinedLoanRequestStatus = "AWAITING_OFFER"
	CombinedLoanRequestStatusAwaitingConfirm CombinedLoanRequestStatus = "AWAITING_CONFIRM"
	CombinedLoanRequestStatusPackageCreating CombinedLoanRequestStatus = "PACKAGE_CREATING"
	CombinedLoanRequestStatusPackageCreated  CombinedLoanRequestStatus = "PACKAGE_CREATED"
	CombinedLoanRequestStatusCancelled       CombinedLoanRequestStatus = "CANCELLED"
)
