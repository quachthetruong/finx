package http

import (
	"strings"

	"github.com/shopspring/decimal"

	"financing-offer/internal/core"
	"financing-offer/internal/core/entity"
	"financing-offer/pkg/optional"
)

type GetLoanPackageOfferRequest struct {
	Paging              core.Paging
	Symbol              string `form:"symbol"`
	OfferInterestStatus string `form:"status"`
	AssetType           string `form:"assetType"`
}

func (r GetLoanPackageOfferRequest) toFilter() entity.LoanPackageOfferFilter {
	splitStatus := strings.Split(r.OfferInterestStatus, ",")
	statuses := make([]string, 0, len(splitStatus))
	for _, status := range splitStatus {
		if len(status) > 0 {
			statuses = append(statuses, status)
		}
	}
	assetType := optional.Some(entity.AssetTypeUnderlying)
	if r.AssetType != "" {
		assetType = optional.Some(entity.AssetType(r.AssetType))
	}
	return entity.LoanPackageOfferFilter{
		Paging:              r.Paging,
		Symbol:              optional.FromValueNonZero(r.Symbol),
		OfferInterestStatus: statuses,
		AssetType:           assetType,
	}
}

type CreateOfferUpdateRequest struct {
	Status   string `json:"status"`
	Category string `json:"category"`
	Note     string `json:"note"`
}

type AdminAssignLoanIdRequest struct {
	LoanId       int64           `json:"loanId" binding:"required"`
	InterestRate decimal.Decimal `json:"interestRate" binding:"required"`
}
