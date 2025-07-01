package http

import "financing-offer/internal/core/entity"

type VerifyInvestorAccountMarginStatusRequest struct {
	InvestorId string `json:"investorId" binding:"required"`
}

func (r VerifyInvestorAccountMarginStatusRequest) toEntity(accountNo string) entity.InvestorAccount {
	return entity.InvestorAccount{
		AccountNo:  accountNo,
		InvestorId: r.InvestorId,
	}
}
