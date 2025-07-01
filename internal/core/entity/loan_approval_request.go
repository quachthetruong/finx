package entity

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

type LoanApprovalRequest struct {
	LoanRequestId      int64           `json:"loan_request_id"`
	CreateAt           time.Time       `json:"create_at"`
	SubmissionId       int64           `json:"submission_id"`
	SubmissionBy       string          `json:"submission_by"`
	SubmissionCreateAt time.Time       `json:"submission_create_at"`
	InvestorId         string          `json:"investor_id"`
	AccountNo          string          `json:"account_no"`
	Symbol             string          `json:"symbol"`
	LoanRate           decimal.Decimal `json:"loan_rate"`
	Source             string          `json:"source"`
	InterestRate       decimal.Decimal `json:"interest_rate"`
	BuyingFee          decimal.Decimal `json:"buying_fee"`
	Term               int32           `json:"term"`
	Description        string          `json:"description"`
	CategoryId         int64           `json:"category_id"`
}

func (loanApprovalRequest *LoanApprovalRequest) ToOdooFormat() []map[string]any {
	return []map[string]any{
		{
			"loan_request_id":      strconv.FormatInt(loanApprovalRequest.LoanRequestId, 10),
			"create_at":            loanApprovalRequest.CreateAt.UTC().Format(time.DateTime),
			"submission_id":        strconv.FormatInt(loanApprovalRequest.SubmissionId, 10),
			"submission_by":        loanApprovalRequest.SubmissionBy,
			"submission_create_at": loanApprovalRequest.SubmissionCreateAt.UTC().Format(time.DateTime),
			"investor_id":          loanApprovalRequest.InvestorId,
			"account_no":           loanApprovalRequest.AccountNo,
			"symbol":               loanApprovalRequest.Symbol,
			"loan_rate":            loanApprovalRequest.LoanRate.String(),
			"source":               loanApprovalRequest.Source,
			"interest_rate":        loanApprovalRequest.InterestRate.String(),
			"buying_fee":           loanApprovalRequest.BuyingFee.String(),
			"term":                 strconv.FormatInt(int64(loanApprovalRequest.Term), 10),
			"description":          fmt.Sprintf("%s-%s", "test", loanApprovalRequest.Description),
		},
	}
}
