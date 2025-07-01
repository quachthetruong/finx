package entity

import (
	"time"

	"financing-offer/internal/core"
	"financing-offer/pkg/optional"
)

type AwaitingConfirmRequest struct {
	LoanRequest    LoanPackageRequest  `json:"loanRequest"`
	LoanOffer      LoanPackageOffer    `json:"loanOffer"`
	Symbol         Symbol              `json:"symbol"`
	LatestUpdate   *OfflineOfferUpdate `json:"latestUpdate,omitempty"`
	Investor       Investor            `json:"investor"`
	LoanPackageIds string              `json:"loanPackageIds"`
}

type AwaitingConfirmRequestFilter struct {
	core.Paging
	Symbols                []string
	FlowTypes              []string
	AccountNumbers         []string
	InvestorId             optional.Optional[string]
	Ids                    []int64
	StartDate              optional.Optional[time.Time]
	EndDate                optional.Optional[time.Time]
	LatestUpdateCategories []string
	AssetType              optional.Optional[string]
	CustodyCode            optional.Optional[string]
	CustodyCodes           []string
}
