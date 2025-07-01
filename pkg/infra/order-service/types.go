package order_service

import "financing-offer/internal/core/entity"

type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type LoanPackagesResponse struct {
	LoanPackages []entity.AccountLoanPackage `json:"loanPackages"`
}
