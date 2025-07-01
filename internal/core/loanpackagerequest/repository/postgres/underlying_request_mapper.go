package postgres

import (
	"financing-offer/internal/core/entity"
	investorPostgres "financing-offer/internal/core/investor/repository/postgres"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/pkg/querymod"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/shopspring/decimal"
)

type UnderlyingRequest struct {
	model.LoanPackageRequest
	model.Investor
	SubmissionId        *int64  `alias:"s.submission_id"`
	SubmissionStatus    *string `alias:"s.submission_status"`
	SubmissionCreator   *string `alias:"s.submission_creator"`
	SubmissionCreatedAt *string `alias:"s.submission_created_at"`
}

func MapUnderlyingRequestDbToEntity(loanPackageRequest UnderlyingRequest) entity.UnderlyingLoanPackageRequest {
	return entity.UnderlyingLoanPackageRequest{
		Id:                  loanPackageRequest.ID,
		SymbolId:            loanPackageRequest.SymbolID,
		InvestorId:          loanPackageRequest.LoanPackageRequest.InvestorID,
		AccountNo:           loanPackageRequest.AccountNo,
		LoanRate:            loanPackageRequest.LoanRate,
		LimitAmount:         loanPackageRequest.LimitAmount,
		Type:                entity.LoanPackageRequestTypeFromString(loanPackageRequest.Type),
		Status:              entity.LoanPackageRequestStatusFromString(loanPackageRequest.Status),
		CreatedAt:           loanPackageRequest.LoanPackageRequest.CreatedAt,
		UpdatedAt:           loanPackageRequest.LoanPackageRequest.UpdatedAt,
		GuaranteedDuration:  int(loanPackageRequest.GuaranteedDuration),
		AssetType:           entity.AssetType(loanPackageRequest.AssetType),
		InitialRate:         loanPackageRequest.InitialRate,
		ContractSize:        loanPackageRequest.ContractSize,
		Investor:            investorPostgres.MapInvestorDbToEntity(loanPackageRequest.Investor),
		SubmissionId:        loanPackageRequest.SubmissionId,
		SubmissionStatus:    loanPackageRequest.SubmissionStatus,
		SubmissionCreator:   loanPackageRequest.SubmissionCreator,
		SubmissionCreatedAt: loanPackageRequest.SubmissionCreatedAt,
	}
}

func ApplyUnderlyingFilter(filter entity.UnderlyingLoanPackageFilter, otherExpressions []postgres.BoolExpression) postgres.BoolExpression {
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
	if filter.CustodyCode.IsPresent() {
		expr = expr.AND(table.Investor.CustodyCode.EQ(postgres.String(filter.CustodyCode.Get())))
	}
	if len(filter.CustodyCodes) > 0 {
		expr = expr.AND(table.Investor.CustodyCode.IN(querymod.In(filter.CustodyCodes)...))
	}
	for _, otherExpression := range otherExpressions {
		expr = expr.AND(otherExpression)
	}
	return expr
}
