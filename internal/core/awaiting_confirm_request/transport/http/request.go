package http

import (
	"time"

	"financing-offer/internal/core"
	"financing-offer/internal/core/entity"
	"financing-offer/pkg/optional"
)

type GetAllAwaitingConfirmRequestRequest struct {
	Paging                 core.Paging
	Symbols                []string  `form:"symbols"`
	FlowTypes              []string  `form:"flowTypes"`
	AccountNumbers         []string  `form:"accountNumbers"`
	InvestorId             string    `form:"investorId"`
	Ids                    []int64   `form:"ids"`
	StartDate              time.Time `form:"startDate"`
	EndDate                time.Time `form:"endDate"`
	LatestUpdateCategories []string  `form:"latestUpdateCategories"`
	AssetType              string    `form:"assetType"`
	CustodyCode            string    `form:"custodyCode"`
	CustodyCodes           []string  `form:"custodyCodes"`
}

func (r GetAllAwaitingConfirmRequestRequest) toFilter() entity.AwaitingConfirmRequestFilter {
	return entity.AwaitingConfirmRequestFilter{
		Paging:                 r.Paging,
		Symbols:                r.Symbols,
		FlowTypes:              r.FlowTypes,
		AccountNumbers:         r.AccountNumbers,
		InvestorId:             optional.FromValueNonZero(r.InvestorId),
		Ids:                    r.Ids,
		StartDate:              optional.FromValueNonZero(r.StartDate),
		EndDate:                optional.FromValueNonZero(r.EndDate),
		LatestUpdateCategories: r.LatestUpdateCategories,
		AssetType:              optional.FromValueNonZero(r.AssetType),
		CustodyCode:            optional.FromValueNonZero(r.CustodyCode),
		CustodyCodes:           r.CustodyCodes,
	}
}
