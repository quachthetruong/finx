package entity

import (
	"time"

	"github.com/shopspring/decimal"

	"financing-offer/internal/core"
	"financing-offer/pkg/optional"
)

type LoanPackageRequest struct {
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
}

type LoanPackageFilter struct {
	core.Paging
	Symbols         []string
	AccountNumbers  []string
	InvestorId      optional.Optional[string]
	Types           []LoanPackageRequestType
	Ids             []int64
	StartDate       optional.Optional[time.Time]
	EndDate         optional.Optional[time.Time]
	LoanPercentFrom optional.Optional[decimal.Decimal]
	LoanPercentTo   optional.Optional[decimal.Decimal]
	LimitAmountFrom optional.Optional[decimal.Decimal]
	LimitAmountTo   optional.Optional[decimal.Decimal]
	Statuses        []LoanPackageRequestStatus
	AssetType       optional.Optional[AssetType]
	CustodyCode     optional.Optional[string]
	CustodyCodes    []string
}

type LoggedRequest struct {
	Id         int64     `json:"id"`
	InvestorId string    `json:"investorId"`
	SymbolId   int64     `json:"symbolId"`
	Reason     string    `json:"reason"`
	Request    string    `json:"request"`
	CreatedAt  time.Time `json:"createdAt"`
}

const (
	LoggedRequestReasonLoanRateExisted = "LOAN_RATE_EXISTED"
)

type LoanPackageRequestStatus string

const (
	LoanPackageRequestStatusPending   LoanPackageRequestStatus = "PENDING"
	LoanPackageRequestStatusConfirmed LoanPackageRequestStatus = "CONFIRMED"
)

func (l LoanPackageRequestStatus) String() string {
	return string(l)
}

func LoanPackageRequestStatusFromString(s string) LoanPackageRequestStatus {
	switch s {
	case string(LoanPackageRequestStatusPending):
		return LoanPackageRequestStatusPending
	case string(LoanPackageRequestStatusConfirmed):
		return LoanPackageRequestStatusConfirmed
	default:
		return LoanPackageRequestStatusPending
	}
}

type LoanPackageRequestType string

const (
	LoanPackageRequestTypeFlexible   LoanPackageRequestType = "FLEXIBLE"
	LoanPackageRequestTypeGuaranteed LoanPackageRequestType = "GUARANTEED"
)

func (l LoanPackageRequestType) String() string {
	return string(l)
}

func (l LoanPackageRequestType) StringNotify() string {
	switch l {
	case LoanPackageRequestTypeFlexible:
		return "Linh hoạt"
	case LoanPackageRequestTypeGuaranteed:
		return "Đảm bảo"
	default:
		return "Đảm bảo"
	}
}

func LoanPackageRequestTypeFromString(s string) LoanPackageRequestType {
	switch s {
	case string(LoanPackageRequestTypeFlexible):
		return LoanPackageRequestTypeFlexible
	case string(LoanPackageRequestTypeGuaranteed):
		return LoanPackageRequestTypeGuaranteed
	default:
		return LoanPackageRequestTypeGuaranteed
	}
}

type LoanPackageRequestConfirmedNotify struct {
	Name           string
	Symbol         string
	NumberOfOffers int64
	InvestorId     string
	AccountNo      string
	AccountNoDesc  string
	CreatedAt      time.Time
}

type LoanPackageRequestDeclinedNotify struct {
	InvestorId    string
	RequestName   string
	AccountNo     string
	AccountNoDesc string
	Symbol        string
	CreatedAt     time.Time
}

type LoanPackageDerivativeRequestDeclinedNotify struct {
	InvestorId    string
	RequestName   string
	AccountNo     string
	AccountNoDesc string
	Symbol        string
	AssetType     string
	CreatedAt     time.Time
}

type RequestOnlineConfirmationNotify struct {
	InvestorId      string
	RequestName     string
	AccountNo       string
	AccountNoDesc   string
	OfferId         int64
	OfferInterestId int64
	Symbol          string
	CreatedAt       time.Time
}

type RequestOfflineConfirmation struct {
	InvestorId    string
	RequestName   string
	AccountNo     string
	AccountNoDesc string
	Symbol        string
	CreatedAt     time.Time
}

type DerivativeRequestOfflineConfirmation struct {
	InvestorId    string
	RequestName   string
	AccountNo     string
	AccountNoDesc string
	Symbol        string
	AssetType     string
	CreatedAt     time.Time
}
