package postgres

import (
	"github.com/volatiletech/null/v9"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

func MapLoanContractDbToEntity(l model.LoanContract) entity.LoanContract {
	return entity.LoanContract{
		Id:                   l.ID,
		LoanOfferInterestId:  l.LoanOfferInterestID,
		SymbolId:             l.SymbolID,
		InvestorId:           l.InvestorID,
		AccountNo:            l.AccountNo,
		LoanId:               l.LoanID,
		CreatedAt:            l.CreatedAt,
		UpdatedAt:            l.UpdatedAt,
		GuaranteedEndAt:      l.GuaranteedEndAt.Time,
		LoanPackageAccountId: l.LoanPackageAccountID,
		LoanProductIdRef:     l.LoanProductIDRef,
	}
}

func MapLoanContractEntityToDb(l entity.LoanContract) model.LoanContract {
	res := model.LoanContract{
		ID:                   l.Id,
		LoanOfferInterestID:  l.LoanOfferInterestId,
		SymbolID:             l.SymbolId,
		InvestorID:           l.InvestorId,
		AccountNo:            l.AccountNo,
		LoanID:               l.LoanId,
		CreatedAt:            l.CreatedAt,
		UpdatedAt:            l.UpdatedAt,
		LoanPackageAccountID: l.LoanPackageAccountId,
		LoanProductIDRef:     l.LoanProductIdRef,
	}
	if !l.GuaranteedEndAt.IsZero() {
		res.GuaranteedEndAt = null.TimeFrom(l.GuaranteedEndAt)
	}
	return res
}

func MapLoanContractsEntityToDb(l []entity.LoanContract) []model.LoanContract {
	res := make([]model.LoanContract, len(l))
	for i, v := range l {
		res[i] = MapLoanContractEntityToDb(v)
	}
	return res
}
