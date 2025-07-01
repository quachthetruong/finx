package postgres

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/shopspring/decimal"

	"financing-offer/internal/core"
	"financing-offer/internal/core/entity"
	investorPostgres "financing-offer/internal/core/investor/repository/postgres"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/pkg/querymod"
)

type LoanPackageRequestWithAdditionalInfo struct {
	model.LoanPackageRequest
	model.Investor
}

func MapLoanPackageRequestsWithAdditionalInfoDbToEntity(loanPackageRequests []LoanPackageRequestWithAdditionalInfo) []entity.LoanPackageRequest {
	dest := make([]entity.LoanPackageRequest, 0, len(loanPackageRequests))
	for _, loanPackageRequest := range loanPackageRequests {
		dest = append(dest, MapLoanPackageRequestWithAdditionalInfoDbToEntity(loanPackageRequest))
	}
	return dest
}

func MapLoanPackageRequestWithAdditionalInfoDbToEntity(loanPackageRequest LoanPackageRequestWithAdditionalInfo) entity.LoanPackageRequest {
	return entity.LoanPackageRequest{
		Id:                 loanPackageRequest.ID,
		SymbolId:           loanPackageRequest.SymbolID,
		InvestorId:         loanPackageRequest.LoanPackageRequest.InvestorID,
		AccountNo:          loanPackageRequest.AccountNo,
		LoanRate:           loanPackageRequest.LoanRate,
		LimitAmount:        loanPackageRequest.LimitAmount,
		Type:               entity.LoanPackageRequestTypeFromString(loanPackageRequest.Type),
		Status:             entity.LoanPackageRequestStatusFromString(loanPackageRequest.Status),
		CreatedAt:          loanPackageRequest.LoanPackageRequest.CreatedAt,
		UpdatedAt:          loanPackageRequest.LoanPackageRequest.UpdatedAt,
		GuaranteedDuration: int(loanPackageRequest.GuaranteedDuration),
		AssetType:          entity.AssetType(loanPackageRequest.AssetType),
		InitialRate:        loanPackageRequest.InitialRate,
		ContractSize:       loanPackageRequest.ContractSize,
		Investor:           investorPostgres.MapInvestorDbToEntity(loanPackageRequest.Investor),
	}
}

func MapLoanPackageRequestsDbToEntity(loanPackageRequests []model.LoanPackageRequest) []entity.LoanPackageRequest {
	dest := make([]entity.LoanPackageRequest, 0, len(loanPackageRequests))
	for _, loanPackageRequest := range loanPackageRequests {
		dest = append(dest, MapLoanPackageRequestDbToEntity(loanPackageRequest))
	}
	return dest
}

func MapLoanPackageRequestDbToEntity(loanPackageRequest model.LoanPackageRequest) entity.LoanPackageRequest {
	return entity.LoanPackageRequest{
		Id:                 loanPackageRequest.ID,
		SymbolId:           loanPackageRequest.SymbolID,
		InvestorId:         loanPackageRequest.InvestorID,
		AccountNo:          loanPackageRequest.AccountNo,
		LoanRate:           loanPackageRequest.LoanRate,
		LimitAmount:        loanPackageRequest.LimitAmount,
		Type:               entity.LoanPackageRequestTypeFromString(loanPackageRequest.Type),
		Status:             entity.LoanPackageRequestStatusFromString(loanPackageRequest.Status),
		CreatedAt:          loanPackageRequest.CreatedAt,
		UpdatedAt:          loanPackageRequest.UpdatedAt,
		GuaranteedDuration: int(loanPackageRequest.GuaranteedDuration),
		AssetType:          entity.AssetType(loanPackageRequest.AssetType),
		InitialRate:        loanPackageRequest.InitialRate,
		ContractSize:       loanPackageRequest.ContractSize,
	}
}

func MapLoanPackageRequestEntityToDb(loanPackageRequest entity.LoanPackageRequest) model.LoanPackageRequest {
	request := model.LoanPackageRequest{
		ID:                 loanPackageRequest.Id,
		SymbolID:           loanPackageRequest.SymbolId,
		InvestorID:         loanPackageRequest.InvestorId,
		AccountNo:          loanPackageRequest.AccountNo,
		LoanRate:           loanPackageRequest.LoanRate,
		LimitAmount:        loanPackageRequest.LimitAmount,
		Type:               loanPackageRequest.Type.String(),
		Status:             loanPackageRequest.Status.String(),
		CreatedAt:          loanPackageRequest.CreatedAt,
		UpdatedAt:          loanPackageRequest.UpdatedAt,
		GuaranteedDuration: int32(loanPackageRequest.GuaranteedDuration),
		AssetType:          model.AssetType(loanPackageRequest.AssetType),
		InitialRate:        loanPackageRequest.InitialRate,
		ContractSize:       loanPackageRequest.ContractSize,
	}
	return request
}

func ApplyFilter(filter entity.LoanPackageFilter) postgres.BoolExpression {
	expr := postgres.Bool(true)
	if len(filter.Symbols) > 0 {
		expr = expr.AND(table.Symbol.Symbol.IN(querymod.In(filter.Symbols)...))
	}
	if filter.InvestorId.IsPresent() {
		expr = expr.AND(table.LoanPackageRequest.InvestorID.EQ(postgres.String(filter.InvestorId.Get())))
	}
	if len(filter.Types) > 0 {
		types := make([]string, len(filter.Types))
		for _, t := range filter.Types {
			types = append(types, t.String())
		}
		expr = expr.AND(table.LoanPackageRequest.Type.IN(querymod.In(types)...))
	}
	if len(filter.Statuses) > 0 {
		statuses := make([]string, len(filter.Statuses))
		for _, t := range filter.Statuses {
			statuses = append(statuses, t.String())
		}
		expr = expr.AND(table.LoanPackageRequest.Status.IN(querymod.In(statuses)...))
	}
	if len(filter.Ids) > 0 {
		expr = expr.AND(table.LoanPackageRequest.ID.IN(querymod.In(filter.Ids)...))
	}
	if len(filter.AccountNumbers) > 0 {
		expr = expr.AND(table.LoanPackageRequest.AccountNo.IN(querymod.In(filter.AccountNumbers)...))
	}
	if filter.StartDate.IsPresent() {
		expr = expr.AND(table.LoanPackageRequest.CreatedAt.GT_EQ(postgres.TimestampT(filter.StartDate.Get())))
	}
	if filter.EndDate.IsPresent() {
		expr = expr.AND(table.LoanPackageRequest.CreatedAt.LT_EQ(postgres.TimestampT(filter.EndDate.Get())))
	}
	if filter.LoanPercentFrom.IsPresent() {
		rate := filter.LoanPercentFrom.Get().Div(decimal.NewFromInt(100))
		expr = expr.AND(table.LoanPackageRequest.LoanRate.GT_EQ(postgres.Decimal(rate.String())))
	}
	if filter.LoanPercentTo.IsPresent() {
		rate := filter.LoanPercentTo.Get().Div(decimal.NewFromInt(100))
		expr = expr.AND(table.LoanPackageRequest.LoanRate.LT_EQ(postgres.Decimal(rate.String())))
	}
	if filter.LimitAmountFrom.IsPresent() {
		expr = expr.AND(table.LoanPackageRequest.LimitAmount.GT_EQ(postgres.Decimal(filter.LimitAmountFrom.Get().String())))
	}
	if filter.LimitAmountTo.IsPresent() {
		expr = expr.AND(table.LoanPackageRequest.LimitAmount.LT_EQ(postgres.Decimal(filter.LimitAmountTo.Get().String())))
	}
	if filter.AssetType.IsPresent() {
		expr = expr.AND(table.LoanPackageRequest.AssetType.EQ(postgres.NewEnumValue(filter.AssetType.Get().String())))
	}
	if filter.CustodyCode.IsPresent() {
		expr = expr.AND(table.Investor.CustodyCode.EQ(postgres.String(filter.CustodyCode.Get())))
	}
	if len(filter.CustodyCodes) > 0 {
		expr = expr.AND(table.Investor.CustodyCode.IN(querymod.In(filter.CustodyCodes)...))
	}
	return expr
}

func ApplySort(filter entity.LoanPackageFilter) []postgres.OrderByClause {
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

func MapLoggedRequestDbToEntity(loggedRequest model.LoggedRequest) entity.LoggedRequest {
	return entity.LoggedRequest{
		Id:         loggedRequest.ID,
		InvestorId: loggedRequest.InvestorID,
		SymbolId:   loggedRequest.SymbolID,
		Reason:     loggedRequest.Reason,
		Request:    loggedRequest.Request,
		CreatedAt:  loggedRequest.CreatedAt,
	}
}

func MapLoggedRequestEntityToDb(loggedRequest entity.LoggedRequest) model.LoggedRequest {
	return model.LoggedRequest{
		ID:         loggedRequest.Id,
		InvestorID: loggedRequest.InvestorId,
		SymbolID:   loggedRequest.SymbolId,
		Reason:     loggedRequest.Reason,
		Request:    loggedRequest.Request,
		CreatedAt:  loggedRequest.CreatedAt,
	}
}
