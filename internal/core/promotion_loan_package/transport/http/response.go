package http

import "financing-offer/internal/core/entity"

type GetPromotionLoanPackageBySymbolResponse struct {
	Symbol      string                                   `json:"symbol"`
	CustodyCode string                                   `json:"custodyCode"`
	AccountNos  []entity.AccountLoanPackageWithAccountNo `json:"accountNos"`
}

type PromotionPackagesWithAccountNo struct {
	AccountNo string                                `json:"accountNo"`
	Symbols   []entity.AccountLoanPackageWithSymbol `json:"symbols"`
}

type PromotionLoanPackages struct {
	AccountNo    string                                  `json:"accountNo"`
	LoanPackages []entity.LoanPackageWithCampaignProduct `json:"loanPackages"`
}

type GetPromotionLoanPackageBySymbolV2Response struct {
	Symbol       string                                   `json:"symbol"`
	CustodyCode  string                                   `json:"custodyCode"`
	LoanPackages []entity.AccountLoanPackageWithAccountNo `json:"loanPackages"`
}
