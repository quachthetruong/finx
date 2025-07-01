package http

type CreateLoanContractRequest struct {
	LoanPackageOfferId         int64 `json:"loanPackageOfferId" binding:"required"`
	LoanPackageOfferInterestId int64 `json:"loanPackageOfferInterestId" binding:"required"`
}

type AssignLoanContractRequest struct {
	LoanPackageId int64 `json:"loanPackageId" binding:"required"`
}
