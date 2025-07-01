package entity

import (
	"financing-offer/pkg/optional"
	"github.com/shopspring/decimal"
	"time"
)

type UnderlyingLoanPackageRequest struct {
	Id                 int64                    `json:"id"`
	SymbolId           int64                    `json:"symbolId"`
	InvestorId         string                   `json:"investorId"`
	AccountNo          string                   `json:"accountNo"`
	LoanRate           decimal.Decimal          `json:"loanRate"`
	LimitAmount        decimal.Decimal          `json:"limitAmount"`
	Type               LoanPackageRequestType   `json:"type"`
	Status             LoanPackageRequestStatus `json:"status"`
	GuaranteedDuration int                      `json:"guaranteedDuration"`
	AssetType          AssetType                `json:"assetType"`
	InitialRate        decimal.Decimal          `json:"initialRate"`
	ContractSize       int64                    `json:"contractSize"`
	CreatedAt          time.Time                `json:"createdAt"`
	UpdatedAt          time.Time                `json:"updatedAt"`

	Investor Investor `json:"investor"`

	SubmissionId        *int64  `json:"submissionId"`
	SubmissionStatus    *string `json:"submissionStatus"`
	SubmissionCreator   *string `json:"submissionCreator"`
	SubmissionCreatedAt *string `json:"submissionCreatedAt"`
}

type UnderlyingLoanPackageFilter struct {
	Symbols            []string
	AccountNumbers     []string
	InvestorId         optional.Optional[string]
	Types              []LoanPackageRequestType
	Ids                []int64
	StartDate          optional.Optional[time.Time]
	EndDate            optional.Optional[time.Time]
	LoanPercentFrom    optional.Optional[decimal.Decimal]
	LoanPercentTo      optional.Optional[decimal.Decimal]
	LimitAmountFrom    optional.Optional[decimal.Decimal]
	LimitAmountTo      optional.Optional[decimal.Decimal]
	Statuses           []LoanPackageRequestStatus
	CustodyCode        optional.Optional[string]
	CustodyCodes       []string
	SubmissionStatuses []SubmissionSheetStatus
}
