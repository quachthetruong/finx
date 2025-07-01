package entity

import (
	"time"

	"financing-offer/internal/core"
	"financing-offer/pkg/optional"
)

type LoanPackageOffer struct {
	Id                   int64     `json:"id"`
	LoanPackageRequestId int64     `json:"loanPackageRequestId"`
	OfferedBy            string    `json:"offeredBy"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	ExpiredAt            time.Time `json:"expiredAt"`
	FlowType             FlowType  `json:"flowType"`

	Symbol                    *Symbol                    `json:"symbol,omitempty"`
	LoanPackageRequest        *LoanPackageRequest        `json:"loanPackageRequest,omitempty"`
	LoanPackageOfferInterests []LoanPackageOfferInterest `json:"loanPackageOfferInterests"`
}

func (o LoanPackageOffer) IsExpired() bool {
	if o.ExpiredAt.IsZero() {
		return false
	}
	return o.ExpiredAt.Before(time.Now())
}

type LoanPackageOfferFilter struct {
	core.Paging
	InvestorId          string                       `json:"investorId"`
	Symbol              optional.Optional[string]    `json:"symbol"`
	OfferInterestStatus []string                     `json:"offerInterestStatus"`
	AssetType           optional.Optional[AssetType] `json:"assetType"`
}
