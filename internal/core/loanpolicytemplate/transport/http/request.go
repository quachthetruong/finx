package http

import (
	"github.com/shopspring/decimal"

	"financing-offer/internal/core/entity"
)

type CreateLoanPackageTemplateRequest struct {
	CreatedBy                string          `json:"createdBy"`
	Name                     string          `json:"name"`
	InterestRate             decimal.Decimal `json:"interestRate"`
	InterestBasis            int16           `json:"interestBasis"`
	Term                     int32           `json:"term"`
	PoolIdRef                int64           `json:"poolIdRef"`
	OverdueInterest          decimal.Decimal `json:"overdueInterest"`
	AllowExtendLoanTerm      bool            `json:"allowExtendLoanTerm"`
	AllowEarlyPayment        bool            `json:"allowEarlyPayment"`
	PreferentialPeriod       int32           `json:"preferentialPeriod"`
	PreferentialInterestRate decimal.Decimal `json:"preferentialInterestRate"`
}

func (r CreateLoanPackageTemplateRequest) toEntity() entity.LoanPolicyTemplate {
	return entity.LoanPolicyTemplate{
		UpdatedBy:                r.CreatedBy,
		Name:                     r.Name,
		InterestRate:             r.InterestRate,
		InterestBasis:            r.InterestBasis,
		Term:                     r.Term,
		PoolIdRef:                r.PoolIdRef,
		OverdueInterest:          r.OverdueInterest,
		AllowExtendLoanTerm:      r.AllowExtendLoanTerm,
		AllowEarlyPayment:        r.AllowEarlyPayment,
		PreferentialPeriod:       r.PreferentialPeriod,
		PreferentialInterestRate: r.PreferentialInterestRate,
	}
}

type UpdateLoanPackageTemplateRequest struct {
	Id                       int64           `json:"id"`
	UpdatedBy                string          `json:"updatedBy"`
	Name                     string          `json:"name"`
	InterestRate             decimal.Decimal `json:"interestRate"`
	InterestBasis            int16           `json:"interestBasis"`
	Term                     int32           `json:"term"`
	PoolIdRef                int64           `json:"poolIdRef"`
	OverdueInterest          decimal.Decimal `json:"overdueInterest"`
	AllowExtendLoanTerm      bool            `json:"allowExtendLoanTerm"`
	AllowEarlyPayment        bool            `json:"allowEarlyPayment"`
	PreferentialPeriod       int32           `json:"preferentialPeriod"`
	PreferentialInterestRate decimal.Decimal `json:"preferentialInterestRate"`
	ProductCategoryId        int32           `json:"productCategoryId"`
	AllowedOverdueLoanInDays int32           `json:"allowedOverdueLoanInDays"`
}

func (r UpdateLoanPackageTemplateRequest) toEntity(id int64) entity.LoanPolicyTemplate {
	return entity.LoanPolicyTemplate{
		Id:                       id,
		UpdatedBy:                r.UpdatedBy,
		Name:                     r.Name,
		InterestRate:             r.InterestRate,
		InterestBasis:            r.InterestBasis,
		Term:                     r.Term,
		PoolIdRef:                r.PoolIdRef,
		OverdueInterest:          r.OverdueInterest,
		AllowExtendLoanTerm:      r.AllowExtendLoanTerm,
		AllowEarlyPayment:        r.AllowEarlyPayment,
		PreferentialPeriod:       r.PreferentialPeriod,
		PreferentialInterestRate: r.PreferentialInterestRate,
	}
}
