package entity

import (
	"time"

	"github.com/shopspring/decimal"

	"financing-offer/internal/core"
)

type LoanPackageOfferInterest struct {
	Id                      int64                          `json:"id"`
	LoanPackageOfferId      int64                          `json:"loanPackageOfferId"`
	SubmissionSheetDetailId int64                          `json:"submissionSheetDetailId"`
	ScoreGroupInterestId    int64                          `json:"scoreGroupInterestId"`
	LimitAmount             decimal.Decimal                `json:"limitAmount"`
	LoanRate                decimal.Decimal                `json:"loanRate"`
	InterestRate            decimal.Decimal                `json:"interestRate"`
	Status                  LoanPackageOfferInterestStatus `json:"status"`
	CreatedAt               time.Time                      `json:"createdAt"`
	UpdatedAt               time.Time                      `json:"updatedAt"`
	LoanID                  int64                          `json:"loanId"`
	CancelledBy             string                         `json:"cancelledBy"`
	CancelledAt             time.Time                      `json:"cancelledAt"`
	Term                    int                            `json:"term"`
	FeeRate                 decimal.Decimal                `json:"feeRate"`
	CancelledReason         CancelledReason                `json:"cancelledReason"`
	AssetType               AssetType                      `json:"assetType"`
	InitialRate             decimal.Decimal                `json:"initialRate"`
	ContractSize            int64                          `json:"contractSize"`

	LoanContract     *LoanContract     `json:"loanContract,omitempty"`
	LoanPackageOffer *LoanPackageOffer `json:"loanPackageOffer,omitempty"`
}

type OfferInterestFilter struct {
	core.Paging
	Statuses []LoanPackageOfferInterestStatus `json:"statuses"`
}

type LoanPackageOfferInterestStatus string

const (
	LoanPackageOfferInterestStatusPending             LoanPackageOfferInterestStatus = "PENDING"
	LoanPackageOfferInterestStatusCancelled           LoanPackageOfferInterestStatus = "CANCELLED"
	LoanPackageOfferInterestStatusSigned              LoanPackageOfferInterestStatus = "SIGNED"
	LoanPackageOfferInterestStatusLoanPackageCreated  LoanPackageOfferInterestStatus = "PACKAGE_CREATED"
	LoanPackageOfferInterestStatusCreatingLoanPackage LoanPackageOfferInterestStatus = "PACKAGE_CREATING"
)

func (s LoanPackageOfferInterestStatus) String() string {
	return string(s)
}

func (s LoanPackageOfferInterestStatus) NextStatuses(flowType FlowType) []LoanPackageOfferInterestStatus {
	switch flowType {
	case FLowTypeDnseOffline, FlowTypeDnseOnline:
		{
			switch s {
			case LoanPackageOfferInterestStatusPending:
				return []LoanPackageOfferInterestStatus{LoanPackageOfferInterestStatusCancelled, LoanPackageOfferInterestStatusLoanPackageCreated, LoanPackageOfferInterestStatusCreatingLoanPackage}
			case LoanPackageOfferInterestStatusCreatingLoanPackage:
				return []LoanPackageOfferInterestStatus{LoanPackageOfferInterestStatusLoanPackageCreated}
			default:
				return []LoanPackageOfferInterestStatus{}
			}
		}
	default:
		return []LoanPackageOfferInterestStatus{}
	}
}

func LoanPackageOfferInterestStatusFromString(str string) LoanPackageOfferInterestStatus {
	switch str {
	case LoanPackageOfferInterestStatusPending.String():
		return LoanPackageOfferInterestStatusPending
	case LoanPackageOfferInterestStatusCancelled.String():
		return LoanPackageOfferInterestStatusCancelled
	case LoanPackageOfferInterestStatusSigned.String():
		return LoanPackageOfferInterestStatusSigned
	case LoanPackageOfferInterestStatusLoanPackageCreated.String():
		return LoanPackageOfferInterestStatusLoanPackageCreated
	case LoanPackageOfferInterestStatusCreatingLoanPackage.String():
		return LoanPackageOfferInterestStatusCreatingLoanPackage
	}
	return ""
}

type LoanPackageOfferReadyNotify struct {
	InvestorId      string
	RequestName     string
	AccountNo       string
	AccountNoDesc   string
	Symbol          string
	OfferId         int64
	OfferInterestId int64
	LoanRate        decimal.Decimal
	LoanType        LoanPackageRequestType
	InterestRate    decimal.Decimal
	LoanPackageId   int64
	CreatedAt       time.Time
}

type DerivativeLoanPackageOfferReadyNotify struct {
	InvestorId      string
	RequestName     string
	AccountNo       string
	AccountNoDesc   string
	Symbol          string
	OfferId         int64
	OfferInterestId int64
	LoanPackageId   int64
	AssetType       string
	CreatedAt       time.Time
}

type CancelledReason string

func (c CancelledReason) String() string {
	return string(c)
}

const (
	LoanPackageOfferCancelledReasonUnknown           CancelledReason = "UNKNOWN"
	LoanPackageOfferCancelledReasonExpired           CancelledReason = "EXPIRED"
	LoanPackageOfferCancelledReasonInvestor          CancelledReason = "INVESTOR"
	LoanPackageOfferCancelledReasonAdmin             CancelledReason = "ADMIN"
	LoanPackageOfferCancelledReasonAlternativeOption CancelledReason = "ALTERNATIVE_OPTION"
	LoanPackageOfferCancelledReasonHighLoanRate      CancelledReason = "HIGH_LOAN_RATE"
)
