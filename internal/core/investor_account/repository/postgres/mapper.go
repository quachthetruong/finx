package postgres

import (
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

func MapInvestorAccountDbToEntity(accountDb model.InvestorAccount) entity.InvestorAccount {
	return entity.InvestorAccount{
		AccountNo:    accountDb.AccountNo,
		InvestorId:   accountDb.InvestorID,
		MarginStatus: entity.MarginStatusFromString(accountDb.MarginStatus),
		CreatedAt:    accountDb.CreatedAt,
		UpdatedAt:    accountDb.UpdatedAt,
	}
}

func MapInvestorAccountEntityToDb(account entity.InvestorAccount) model.InvestorAccount {
	return model.InvestorAccount{
		AccountNo:    account.AccountNo,
		InvestorID:   account.InvestorId,
		MarginStatus: account.MarginStatus.String(),
		CreatedAt:    account.CreatedAt,
		UpdatedAt:    account.UpdatedAt,
	}
}
