package postgres

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/volatiletech/null/v9"

	"financing-offer/internal/core"
	"financing-offer/internal/core/entity"
	loanContractPostgres "financing-offer/internal/core/loancontract/repository/postgres"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/internal/database/mapper"
)

type OfferInterestWithContract struct {
	model.LoanPackageOfferInterest
	LoanPackageOffer model.LoanPackageOffer
	LoanContract     model.LoanContract
}

func MapOfferInterestWithContractDbToEntity(l OfferInterestWithContract) entity.LoanPackageOfferInterest {
	loanContract := loanContractPostgres.MapLoanContractDbToEntity(l.LoanContract)
	loanOffer := mapper.MapLoanPackageOfferDbToEntity(l.LoanPackageOffer)
	res := MapLoanPackageOfferInterestDbToEntity(l.LoanPackageOfferInterest)
	if l.ScoreGroupInterestID != nil {
		res.ScoreGroupInterestId = *l.ScoreGroupInterestID
	}
	if loanContract.Id != 0 {
		res.LoanContract = &loanContract
	}
	if loanOffer.Id != 0 {
		res.LoanPackageOffer = &loanOffer
	}
	return res
}

func MapOfferInterestWithContractsDbToEntity(ll []OfferInterestWithContract) []entity.LoanPackageOfferInterest {
	result := make([]entity.LoanPackageOfferInterest, 0, len(ll))
	for _, l := range ll {
		result = append(result, MapOfferInterestWithContractDbToEntity(l))
	}
	return result
}

func MapLoanPackageOfferInterestsDbToEntity(ll []model.LoanPackageOfferInterest) []entity.LoanPackageOfferInterest {
	result := make([]entity.LoanPackageOfferInterest, 0, len(ll))
	for _, l := range ll {
		result = append(result, MapLoanPackageOfferInterestDbToEntity(l))
	}
	return result
}

func MapLoanPackageOfferInterestDbToEntity(l model.LoanPackageOfferInterest) entity.LoanPackageOfferInterest {
	res := entity.LoanPackageOfferInterest{
		Id:                 l.ID,
		LoanPackageOfferId: l.LoanPackageOfferID,
		LimitAmount:        l.LimitAmount,
		LoanRate:           l.LoanRate,
		InterestRate:       l.InterestRate,
		Status:             entity.LoanPackageOfferInterestStatusFromString(l.Status),
		CreatedAt:          l.CreatedAt,
		UpdatedAt:          l.UpdatedAt,
		LoanID:             l.LoanID,
		CancelledBy:        l.CancelledBy,
		CancelledAt:        l.CancelledAt.Time,
		Term:               int(l.Term),
		FeeRate:            l.FeeRate,
		CancelledReason:    entity.CancelledReason(l.CancelledReason),
		AssetType:          entity.AssetType(l.AssetType),
		ContractSize:       l.ContractSize,
		InitialRate:        l.InitialRate,
	}
	if l.ScoreGroupInterestID != nil {
		res.ScoreGroupInterestId = *l.ScoreGroupInterestID
	}
	if l.SubmissionSheetDetailID != nil {
		res.SubmissionSheetDetailId = *l.SubmissionSheetDetailID
	}

	return res
}

func MapLoanPackageOfferInterestEntityToDb(l entity.LoanPackageOfferInterest) model.LoanPackageOfferInterest {
	res := model.LoanPackageOfferInterest{
		ID:                 l.Id,
		LoanPackageOfferID: l.LoanPackageOfferId,
		LimitAmount:        l.LimitAmount,
		LoanRate:           l.LoanRate,
		InterestRate:       l.InterestRate,
		Status:             l.Status.String(),
		CreatedAt:          l.CreatedAt,
		UpdatedAt:          l.UpdatedAt,
		LoanID:             l.LoanID,
		CancelledBy:        l.CancelledBy,
		Term:               int32(l.Term),
		FeeRate:            l.FeeRate,
		CancelledReason:    l.CancelledReason.String(),
		AssetType:          model.AssetType(l.AssetType),
		ContractSize:       l.ContractSize,
		InitialRate:        l.InitialRate,
	}
	if !l.CancelledAt.IsZero() {
		res.CancelledAt = null.TimeFrom(l.CancelledAt)
	}
	if l.ScoreGroupInterestId != 0 {
		res.ScoreGroupInterestID = &l.ScoreGroupInterestId
	}
	if l.SubmissionSheetDetailId != 0 {
		res.SubmissionSheetDetailID = &l.SubmissionSheetDetailId
	}
	return res
}

func ApplySort(filter entity.OfferInterestFilter) []postgres.OrderByClause {
	expr := make([]postgres.OrderByClause, 0, len(filter.Sort))
	for _, s := range filter.Sort {
		var column postgres.Column
		for _, c := range table.LoanPackageOfferInterest.AllColumns {
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

func ApplyFilter(filter entity.OfferInterestFilter) postgres.BoolExpression {
	stm := postgres.Bool(true)
	if len(filter.Statuses) > 0 {
		statuesString := make([]postgres.Expression, 0, len(filter.Statuses))
		for _, status := range filter.Statuses {
			statuesString = append(statuesString, postgres.String(string(status)))
		}
		stm = stm.AND(table.LoanPackageOfferInterest.Status.IN(statuesString...))
	}
	return stm
}
