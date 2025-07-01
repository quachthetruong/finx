package financialproduct

import (
	"fmt"

	"financing-offer/internal/core/entity"
)

type FinancialAccountsResponse struct {
	Accounts []entity.FinancialAccountDetail `json:"accounts"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
}

type AssignLoanIdAccountNoRequest struct {
	LoanPackageId string `json:"loanPackageId"`
	AccountNo     string `json:"accountNo"`
}

type AssignLoanPackageResponse struct {
	Id            int64  `json:"Id"`
	AccountNo     string `json:"AccountNo"`
	LoanPackageId int64  `json:"LoanPackageId"`
}

type LoanPackageDetailsResponse struct {
	Data []entity.FinancialProductLoanPackage `json:"data"`
}

type ListResponse[T any] struct {
	Data  []T `json:"data"`
	Total int `json:"total"`
	Start int `json:"start"`
	End   int `json:"end"`
}
