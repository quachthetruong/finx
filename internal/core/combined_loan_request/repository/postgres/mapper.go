package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-jet/jet/v2/postgres"

	"financing-offer/internal/core/entity"
	investorPostgres "financing-offer/internal/core/investor/repository/postgres"
	postgres2 "financing-offer/internal/core/loancontract/repository/postgres"
	loanPackageOfferPostgres "financing-offer/internal/core/loanoffer/repository/postgres"
	loanPackageRequestPostgres "financing-offer/internal/core/loanpackagerequest/repository/postgres"
	symbolPostgres "financing-offer/internal/core/symbol/repository/postgres"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/internal/funcs"
	"financing-offer/pkg/querymod"
)

func MapCombinedLoanPackageRequestDbToEntity(request CombinedLoanRequest) (entity.CombinedLoanRequest, error) {
	res := entity.CombinedLoanRequest{
		LoanRequest:                 loanPackageRequestPostgres.MapLoanPackageRequestDbToEntity(request.LoanPackageRequest),
		Symbol:                      symbolPostgres.MapSymbolDbToEntity(request.Symbol),
		Investor:                    investorPostgres.MapInvestorDbToEntity(request.Investor),
		AdminAssignedLoanPackageIds: strings.Join(toSliceNotNull(request.AdminAssignedLoanPackageIds), ", "),
		ActivatedLoanPackageIds:     strings.Join(toSliceNotNull(request.ActivatedLoanPackageIds), ", "),
	}
	if request.LoanPackageOffer.ID != 0 {
		offer := loanPackageOfferPostgres.MapLoanPackageOfferDbToEntity(request.LoanPackageOffer)
		res.LoanOffer = &offer
	}
	if request.LoanContract.ID != 0 {
		contract := postgres2.MapLoanContractDbToEntity(request.LoanContract)
		res.LoanContract = &contract
	}
	if packageCreatedTime := toSliceNotNull(request.PackageCreatedTime); len(packageCreatedTime) > 0 {
		t, err := time.Parse(time.DateTime, packageCreatedTime[0])
		if err != nil {
			return res, fmt.Errorf("MapCombinedLoanPackageRequestDbToEntity: %w", err)
		}
		if !t.IsZero() {
			res.PackageCreatedTime = &t
		}
	}
	res.Status = inferCombinedRequestStatus(request)
	res.CancelledReason = inferCancelledReason(res.Status, request)
	return res, nil
}

func inferCombinedRequestStatus(request CombinedLoanRequest) entity.CombinedLoanRequestStatus {
	if request.LoanPackageRequest.Status == entity.LoanPackageRequestStatusPending.String() {
		return entity.CombinedLoanRequestStatusAwaitingOffer
	} else if offerLineStatuses := toSliceNotNull(request.Statuses); len(offerLineStatuses) > 0 {
		if funcs.AnyEqual(offerLineStatuses, entity.LoanPackageOfferInterestStatusLoanPackageCreated.String()) {
			return entity.CombinedLoanRequestStatusPackageCreated
		} else if funcs.AnyEqual(offerLineStatuses, entity.LoanPackageOfferInterestStatusCreatingLoanPackage.String()) {
			return entity.CombinedLoanRequestStatusPackageCreating
		} else if funcs.AnyEqual(offerLineStatuses, entity.LoanPackageOfferInterestStatusPending.String()) {
			return entity.CombinedLoanRequestStatusAwaitingConfirm
		}
	}
	return entity.CombinedLoanRequestStatusCancelled
}

func inferCancelledReason(status entity.CombinedLoanRequestStatus, request CombinedLoanRequest) string {
	if status != entity.CombinedLoanRequestStatusCancelled {
		return ""
	}
	if cancelReasons := toSliceNotNull(request.CancelledReasons); len(cancelReasons) > 0 {
		for i := len(cancelReasons) - 1; i >= 0; i-- {
			reason := cancelReasons[i]
			if reason != entity.LoanPackageOfferCancelledReasonUnknown.String() {
				return reason
			}
		}
	}
	// admin cancel request without any alternative offers
	return entity.LoanPackageOfferCancelledReasonAdmin.String()
}

func toSliceNotNull(pgType string) []string {
	parts := strings.Split(strings.Trim(pgType, "{}"), ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.Trim(p, "\"\\")
		if trimmed != "NULL" {
			res = append(res, trimmed)
		}
	}
	return res
}

func MapCombinedLoanPackageRequestsDbToEntity(requests []CombinedLoanRequest) ([]entity.CombinedLoanRequest, error) {
	res := make([]entity.CombinedLoanRequest, 0, len(requests))
	for _, request := range requests {
		r, err := MapCombinedLoanPackageRequestDbToEntity(request)
		if err != nil {
			return nil, fmt.Errorf("MapCombinedLoanPackageRequestsDbToEntity: %w", err)
		}
		res = append(res, r)
	}
	return res, nil
}

func ApplyWhere(filter entity.CombinedLoanRequestFilter) postgres.BoolExpression {
	expr := postgres.Bool(true)
	if len(filter.Symbols) > 0 {
		expr = expr.AND(table.Symbol.Symbol.IN(querymod.In(filter.Symbols)...))
	}
	if filter.StartDate.IsPresent() {
		expr = expr.AND(table.LoanPackageRequest.CreatedAt.GT_EQ(postgres.TimestampT(filter.StartDate.Get())))
	}
	if filter.EndDate.IsPresent() {
		expr = expr.AND(table.LoanPackageRequest.CreatedAt.LT_EQ(postgres.TimestampT(filter.EndDate.Get())))
	}
	if filter.OfferDateFrom.IsPresent() {
		expr = expr.AND(table.LoanPackageOffer.CreatedAt.GT_EQ(postgres.TimestampT(filter.OfferDateFrom.Get())))
	}
	if filter.OfferDateTo.IsPresent() {
		expr = expr.AND(table.LoanPackageOffer.CreatedAt.LT_EQ(postgres.TimestampT(filter.OfferDateTo.Get())))
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
	if filter.Status == entity.CombinedLoanRequestStatusAwaitingOffer {
		expr = expr.AND(table.LoanPackageRequest.Status.EQ(postgres.String(entity.LoanPackageRequestStatusPending.String())))
	}
	if filter.Status == entity.CombinedLoanRequestStatusCancelled {
		expr = expr.AND(table.LoanPackageRequest.Status.EQ(postgres.String(entity.LoanPackageRequestStatusConfirmed.String())))
	}
	if len(filter.Ids) > 0 {
		expr = expr.AND(table.LoanPackageRequest.ID.IN(querymod.In(filter.Ids)...))
	}
	if filter.AssetType.IsPresent() {
		expr = expr.AND(table.LoanPackageRequest.AssetType.EQ(postgres.NewEnumValue(filter.AssetType.Get())))
	}
	if filter.CustodyCode.IsPresent() {
		expr = expr.AND(table.Investor.CustodyCode.EQ(postgres.String(filter.CustodyCode.Get())))
	}
	return expr
}

func ApplyHaving(filter entity.CombinedLoanRequestFilter) postgres.BoolExpression {
	expr := postgres.Bool(true)
	if filter.AssignedLoanId.IsPresent() {
		expr = expr.AND(
			postgres.Int(filter.AssignedLoanId.Get()).EQ(
				postgres.IntExp(
					querymod.ArrayAny(
						querymod.ArrayAgg(table.LoanPackageOfferInterest.LoanID),
					),
				),
			),
		)
	}
	if filter.ActivatedLoanId.IsPresent() {
		expr = expr.AND(
			postgres.Int(filter.ActivatedLoanId.Get()).EQ(
				postgres.IntExp(
					querymod.ArrayAny(
						querymod.ArrayAgg(
							postgres.CASE().
								WHEN(table.LoanPackageOfferInterest.Status.EQ(postgres.String(entity.LoanPackageOfferInterestStatusLoanPackageCreated.String()))).
								THEN(table.LoanPackageOfferInterest.LoanID),
						),
					),
				),
			),
		)
	}
	if filter.Status == entity.CombinedLoanRequestStatusAwaitingConfirm {
		expr = expr.AND(
			postgres.String(entity.LoanPackageOfferInterestStatusPending.String()).EQ(
				postgres.StringExp(
					querymod.ArrayAny(
						querymod.ArrayAgg(table.LoanPackageOfferInterest.Status),
					),
				),
			),
		)
	}
	if filter.Status == entity.CombinedLoanRequestStatusPackageCreated {
		expr = expr.AND(
			postgres.String(entity.LoanPackageOfferInterestStatusLoanPackageCreated.String()).EQ(
				postgres.StringExp(
					querymod.ArrayAny(
						querymod.ArrayAgg(table.LoanPackageOfferInterest.Status),
					),
				),
			),
		)
	}
	if filter.Status == entity.CombinedLoanRequestStatusCancelled {
		expr = expr.AND(
			postgres.String(entity.LoanPackageOfferInterestStatusCancelled.String()).EQ(
				postgres.StringExp(
					querymod.ArrayAll(
						querymod.ArrayAgg(table.LoanPackageOfferInterest.Status),
					),
				),
			).OR(
				postgres.String("{NULL}").EQ(
					postgres.StringExp(
						postgres.CAST(querymod.ArrayAgg(table.LoanPackageOfferInterest.Status)).AS_TEXT(),
					),
				),
			),
		)
	}
	return expr
}
