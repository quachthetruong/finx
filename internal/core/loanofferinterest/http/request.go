package loanOfferInterestHttp

import (
	"strings"

	"financing-offer/internal/core"
	"financing-offer/internal/core/entity"
)

type GetAllLoanOfferInterestRequest struct {
	Paging core.Paging
	Status string `form:"status"`
}

func (r GetAllLoanOfferInterestRequest) toFilter() entity.OfferInterestFilter {
	splitStatus := strings.Split(r.Status, ",")
	statuses := make([]entity.LoanPackageOfferInterestStatus, 0, len(splitStatus))
	for _, status := range splitStatus {
		if s := entity.LoanPackageOfferInterestStatusFromString(status); s != "" {
			statuses = append(statuses, s)
		}
	}
	return entity.OfferInterestFilter{
		Paging:   r.Paging,
		Statuses: statuses,
	}
}

type InvestorConfirmLoanPackageOfferInterestRequest struct {
	LoanPackageOfferInterestIds []int64 `json:"loanPackageOfferInterestIds" binding:"required"`
}

type CreateAssignedLoanOfferInterestLoanContractRequest struct {
	LoanPackageAccountId int64                              `json:"loanPackageAccountId"`
	LoanProductIdRef     int64                              `json:"loanProductIdRef"`
	LoanPackage          entity.FinancialProductLoanPackage `json:"loanPackage"`
}
