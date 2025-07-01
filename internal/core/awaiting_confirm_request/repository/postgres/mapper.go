package postgres

import (
	"github.com/go-jet/jet/v2/postgres"

	"financing-offer/internal/core"
	"financing-offer/internal/core/entity"
	investorPostgres "financing-offer/internal/core/investor/repository/postgres"
	loanPackageOfferPostgres "financing-offer/internal/core/loanoffer/repository/postgres"
	loanpackageRequestPostgres "financing-offer/internal/core/loanpackagerequest/repository/postgres"
	offlineOfferPosgres "financing-offer/internal/core/offline_offer_update/repository/postgres"
	symbolPostgres "financing-offer/internal/core/symbol/repository/postgres"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/internal/database/dbmodels/finoffer/public/view"
	"financing-offer/pkg/querymod"
)

func MapAwaitingConfirmRequestDbToEntity(awaitingConfirmRequest AwaitingConfirmRequest) entity.AwaitingConfirmRequest {
	res := entity.AwaitingConfirmRequest{
		LoanRequest:    loanpackageRequestPostgres.MapLoanPackageRequestDbToEntity(awaitingConfirmRequest.LoanPackageRequest),
		LoanOffer:      loanPackageOfferPostgres.MapLoanPackageOfferDbToEntity(awaitingConfirmRequest.LoanPackageOffer),
		Symbol:         symbolPostgres.MapSymbolDbToEntity(awaitingConfirmRequest.Symbol),
		Investor:       investorPostgres.MapInvestorDbToEntity(awaitingConfirmRequest.Investor),
		LoanPackageIds: awaitingConfirmRequest.LoanPackageIds,
	}
	if awaitingConfirmRequest.LatestOfferUpdate.OfferID != nil {
		latestUpdateEntity := offlineOfferPosgres.MapLatestOfferUpdateDbToEntity(awaitingConfirmRequest.LatestOfferUpdate)
		res.LatestUpdate = &latestUpdateEntity
	}
	return res
}

func MapAwaitingConfirmRequestsDbToEntity(awaitingConfirmRequests []AwaitingConfirmRequest) []entity.AwaitingConfirmRequest {
	awaitingConfirmRequestsEntity := make([]entity.AwaitingConfirmRequest, 0, len(awaitingConfirmRequests))
	for _, awaitingConfirmRequest := range awaitingConfirmRequests {
		awaitingConfirmRequestsEntity = append(
			awaitingConfirmRequestsEntity, MapAwaitingConfirmRequestDbToEntity(awaitingConfirmRequest),
		)
	}
	return awaitingConfirmRequestsEntity
}

func ApplyFilter(filter entity.AwaitingConfirmRequestFilter) postgres.BoolExpression {
	expr := postgres.Bool(true)
	if len(filter.Symbols) > 0 {
		expr = expr.AND(table.Symbol.Symbol.IN(querymod.In(filter.Symbols)...))
	}
	if len(filter.FlowTypes) > 0 {
		expr = expr.AND(table.LoanPackageOffer.FlowType.IN(querymod.In(filter.FlowTypes)...))
	}
	if len(filter.AccountNumbers) > 0 {
		expr = expr.AND(table.LoanPackageRequest.AccountNo.IN(querymod.In(filter.AccountNumbers)...))
	}
	if filter.InvestorId.IsPresent() {
		expr = expr.AND(table.LoanPackageRequest.InvestorID.EQ(postgres.String(filter.InvestorId.Get())))
	}
	if len(filter.Ids) > 0 {
		expr = expr.AND(table.LoanPackageRequest.ID.IN(querymod.In(filter.Ids)...))
	}
	if filter.StartDate.IsPresent() {
		expr = expr.AND(table.LoanPackageRequest.CreatedAt.GT_EQ(postgres.TimestampT(filter.StartDate.Get())))
	}
	if filter.EndDate.IsPresent() {
		expr = expr.AND(table.LoanPackageRequest.CreatedAt.LT_EQ(postgres.TimestampT(filter.EndDate.Get())))
	}
	if len(filter.LatestUpdateCategories) > 0 {
		expr = expr.AND(view.LatestOfferUpdate.Category.IN(querymod.In(filter.LatestUpdateCategories)...))
	}
	if filter.AssetType.IsPresent() {
		expr = expr.AND(table.LoanPackageRequest.AssetType.EQ(postgres.NewEnumValue(filter.AssetType.Get())))
	}
	if filter.CustodyCode.IsPresent() {
		expr = expr.AND(table.Investor.CustodyCode.EQ(postgres.String(filter.CustodyCode.Get())))
	}
	if len(filter.CustodyCodes) > 0 {
		expr = expr.AND(table.Investor.CustodyCode.IN(querymod.In(filter.CustodyCodes)...))
	}
	return expr
}

func ApplySort(filter entity.AwaitingConfirmRequestFilter) []postgres.OrderByClause {
	expr := make([]postgres.OrderByClause, 0, len(filter.Sort))
	for _, s := range filter.Sort {
		var column postgres.Column
		for _, c := range table.LoanPackageRequest.AllColumns {
			if c.Name() == s.ColumnName {
				column = c
				break
			}
		}
		if column == nil {
			continue
		}
		if s.Direction == core.DirectionAsc {
			expr = append(expr, column.ASC())
		} else {
			expr = append(expr, column.DESC())
		}
	}
	return expr
}
