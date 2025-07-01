package repository

import (
	"financing-offer/internal/core/entity"
)

type OdooServiceRepository interface {
	SendLoanApprovalRequest(loanApprovalRequest entity.LoanApprovalRequest) error
}
