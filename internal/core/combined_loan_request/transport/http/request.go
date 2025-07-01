package http

import (
	"time"

	"financing-offer/internal/core"
	"financing-offer/internal/core/entity"
	"financing-offer/pkg/optional"
)

type GetAllCombinedRequestsRequest struct {
	Paging          core.Paging
	Symbols         []string  `form:"symbols"`
	StartDate       time.Time `form:"startDate"`
	EndDate         time.Time `form:"endDate"`
	OfferDateFrom   time.Time `form:"offerDateFrom"`
	OfferDateTo     time.Time `form:"offerDateTo"`
	FlowTypes       []string  `form:"flowTypes"`
	AccountNumbers  []string  `form:"accountNumbers"`
	InvestorId      string    `form:"investorId"`
	Status          string    `form:"status"`
	AssignedLoanId  int64     `form:"assignedLoanId"`
	ActivatedLoanId int64     `form:"activatedLoanId"`
	Ids             []int64   `form:"ids"`
	AssetType       string    `form:"assetType"`
	CustodyCode     string    `form:"custodyCode"`
}

func (r GetAllCombinedRequestsRequest) toFilter() entity.CombinedLoanRequestFilter {
	return entity.CombinedLoanRequestFilter{
		Paging:          r.Paging,
		Symbols:         r.Symbols,
		StartDate:       optional.FromValueNonZero(r.StartDate),
		EndDate:         optional.FromValueNonZero(r.EndDate),
		OfferDateFrom:   optional.FromValueNonZero(r.OfferDateFrom),
		OfferDateTo:     optional.FromValueNonZero(r.OfferDateTo),
		FlowTypes:       r.FlowTypes,
		AccountNumbers:  r.AccountNumbers,
		InvestorId:      optional.FromValueNonZero(r.InvestorId),
		Status:          entity.CombinedLoanRequestStatus(r.Status),
		AssignedLoanId:  optional.FromValueNonZero(r.AssignedLoanId),
		ActivatedLoanId: optional.FromValueNonZero(r.ActivatedLoanId),
		Ids:             r.Ids,
		AssetType:       optional.FromValueNonZero(r.AssetType),
		CustodyCode:     optional.FromValueNonZero(r.CustodyCode),
	}
}
