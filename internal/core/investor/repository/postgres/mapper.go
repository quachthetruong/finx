package postgres

import (
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

func MapInvestorsDbToEntity(investorDb []model.Investor) []entity.Investor {
	investors := make([]entity.Investor, len(investorDb))
	for i, investor := range investorDb {
		investors[i] = MapInvestorDbToEntity(investor)
	}
	return investors
}

func MapInvestorDbToEntity(investor model.Investor) entity.Investor {
	return entity.Investor{
		InvestorId:  investor.InvestorID,
		CustodyCode: investor.CustodyCode,
		CreatedAt:   investor.CreatedAt,
		UpdatedAt:   investor.UpdatedAt,
	}
}

func MapInvestorEntityToDb(investor entity.Investor) model.Investor {
	return model.Investor{
		InvestorID:  investor.InvestorId,
		CustodyCode: investor.CustodyCode,
		CreatedAt:   investor.CreatedAt,
		UpdatedAt:   investor.UpdatedAt,
	}
}

func MapInvestorEntitiesToDb(investors []entity.Investor) []model.Investor {
	investorsDb := make([]model.Investor, len(investors))
	for i, investor := range investors {
		investorsDb[i] = MapInvestorEntityToDb(investor)
	}
	return investorsDb
}
