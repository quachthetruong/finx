package http

import (
	"time"

	"github.com/shopspring/decimal"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/core"
	"financing-offer/internal/core/entity"
	"financing-offer/pkg/optional"
)

type GetLoanPackageRequest struct {
	Paging          core.Paging
	Symbols         []string        `form:"symbols"`
	Types           []string        `form:"types"`
	InvestorId      string          `form:"investorId"`
	Statuses        []string        `form:"statuses"`
	Ids             []int64         `form:"ids"`
	StartDate       time.Time       `form:"startDate"`
	EndDate         time.Time       `form:"endDate"`
	LoanPercentFrom decimal.Decimal `form:"loanPercentFrom"`
	LoanPercentTo   decimal.Decimal `form:"loanPercentTo"`
	LimitAmountFrom decimal.Decimal `form:"limitAmountFrom"`
	LimitAmountTo   decimal.Decimal `form:"limitAmountTo"`
	AccountNumbers  []string        `form:"accountNumbers"`
	AssetType       string          `form:"assetType"`
	CustodyCode     string          `form:"custodyCode"`
	CustodyCodes    []string        `form:"custodyCodes"`
}

func (r GetLoanPackageRequest) toEntity() entity.LoanPackageFilter {
	types := make([]entity.LoanPackageRequestType, 0)
	for _, t := range r.Types {
		types = append(types, entity.LoanPackageRequestType(t))
	}
	statuses := make([]entity.LoanPackageRequestStatus, len(r.Statuses))
	for _, s := range r.Statuses {
		statuses = append(statuses, entity.LoanPackageRequestStatus(s))
	}
	assetType := optional.FromValueNonZero(entity.AssetTypeUnderlying)
	if r.AssetType != "" {
		assetType = optional.FromValueNonZero(entity.AssetType(r.AssetType))
	}
	return entity.LoanPackageFilter{
		Paging:          r.Paging,
		Symbols:         r.Symbols,
		InvestorId:      optional.FromValueNonZero(r.InvestorId),
		Ids:             r.Ids,
		AccountNumbers:  r.AccountNumbers,
		StartDate:       optional.FromValueNonZero(r.StartDate),
		EndDate:         optional.FromValueNonZero(r.EndDate),
		LoanPercentFrom: optional.FromValueNonZero(r.LoanPercentFrom),
		LoanPercentTo:   optional.FromValueNonZero(r.LoanPercentTo),
		LimitAmountFrom: optional.FromValueNonZero(r.LimitAmountFrom),
		LimitAmountTo:   optional.FromValueNonZero(r.LimitAmountTo),
		Types:           types,
		Statuses:        statuses,
		AssetType:       assetType,
		CustodyCode:     optional.FromValueNonZero(r.CustodyCode),
		CustodyCodes:    r.CustodyCodes,
	}
}

type CreateLoanPackageRequestUnderlyingRequest struct {
	SymbolId           int64                         `json:"symbolId" binding:"gte=0"`
	LoanRate           decimal.Decimal               `json:"loanRate" binding:"required"`
	LimitAmount        decimal.Decimal               `json:"limitAmount" binding:"required"`
	AccountNo          string                        `json:"accountNo" binding:"required"`
	Type               entity.LoanPackageRequestType `json:"type,omitempty" binding:"required,oneof=FLEXIBLE GUARANTEED"`
	GuaranteedDuration int                           `json:"guaranteedDuration"`
}

func (r CreateLoanPackageRequestUnderlyingRequest) toEntity(investorId string) entity.LoanPackageRequest {
	return entity.LoanPackageRequest{
		InvestorId:         investorId,
		AccountNo:          r.AccountNo,
		SymbolId:           r.SymbolId,
		LoanRate:           r.LoanRate,
		LimitAmount:        r.LimitAmount,
		Type:               r.Type,
		Status:             entity.LoanPackageRequestStatusPending,
		GuaranteedDuration: r.GuaranteedDuration,
		AssetType:          entity.AssetTypeUnderlying,
	}
}

type CreateLoanPackageRequestDerivativeRequest struct {
	SymbolId     int64           `json:"symbolId" binding:"gte=0"`
	InitialRate  decimal.Decimal `json:"initialRate" binding:"required"`
	ContractSize int64           `json:"contractSize" binding:"required"`
	AccountNo    string          `json:"accountNo" binding:"required"`
}

func (r CreateLoanPackageRequestDerivativeRequest) toEntity(investorId string) (entity.LoanPackageRequest, error) {
	if r.InitialRate.LessThan(decimal.NewFromFloat(0.03)) ||
		r.InitialRate.GreaterThan(decimal.NewFromFloat(0.3)) ||
		r.ContractSize <= 0 {
		return entity.LoanPackageRequest{}, apperrors.ErrInvalidInput("invalid initial rate or contract size")
	}
	return entity.LoanPackageRequest{
		InvestorId:   investorId,
		AccountNo:    r.AccountNo,
		SymbolId:     r.SymbolId,
		InitialRate:  r.InitialRate,
		ContractSize: r.ContractSize,
		Type:         entity.LoanPackageRequestTypeFlexible,
		Status:       entity.LoanPackageRequestStatusPending,
		AssetType:    entity.AssetTypeDerivative,
	}, nil
}

type LoggedRequestRequest struct {
	Request CreateLoanPackageRequestUnderlyingRequest `json:"request"`
}

type ConfirmLoanPackageRequestRequest struct {
	LoanId    int64  `json:"loanId"`
	OfferedBy string `json:"offeredBy"`
}

type CancelLoanPackageRequestRequest struct {
	LoanIds   []int64 `json:"loanIds"`
	OfferedBy string  `json:"offeredBy"`
}

type ApproveSubmissionRequest struct {
	SubmissionId int64 `json:"submissionId"`
}

type RejectSubmissionRequest struct {
	SubmissionId int64 `json:"submissionId"`
}
