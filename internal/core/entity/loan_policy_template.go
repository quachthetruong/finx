package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type LoanPolicyTemplate struct {
	Id                       int64           `json:"id"`
	CreatedAt                time.Time       `json:"createdAt"`
	UpdatedAt                time.Time       `json:"updatedAt"`
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
}

type AggregateLoanPolicyTemplate struct {
	Id                       int64           `json:"id"`
	CreatedAt                time.Time       `json:"createdAt"`
	UpdatedAt                time.Time       `json:"updatedAt"`
	UpdatedBy                string          `json:"updatedBy"`
	Name                     string          `json:"name"`
	Source                   string          `json:"source"` // Add source for front-end
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

type TemplateWithProductRate struct {
	LoanPolicySnapShot
	AllowedOverdueLoanInDays int64           `json:"allowedOverdueLoanInDays"`
	ProductRate              decimal.Decimal `json:"productRate"`
	ProductRateForWithdraw   decimal.Decimal `json:"productRateForWithdraw"`
}

type LoanPolicyShorten struct {
	LoanPolicyTemplateId     int64           `json:"loanPolicyTemplateId"`
	AllowedOverdueLoanInDays int64           `json:"allowedOverdueLoanInDays"`
	InitialRateForWithdraw   decimal.Decimal `json:"initialRateForWithdraw"`
	InitialRate              decimal.Decimal `json:"initialRate"`
	Source                   string          `json:"source"`
}

type LoanPolicySnapShot struct {
	LoanPolicyTemplateId     int64           `json:"loanPolicyTemplateId"`
	AllowedOverdueLoanInDays int64           `json:"allowedOverdueLoanInDays"`
	InitialRateForWithdraw   decimal.Decimal `json:"initialRateForWithdraw"`
	InitialRate              decimal.Decimal `json:"initialRate"`
	CreatedAt                time.Time       `json:"createdAt"`
	UpdatedAt                time.Time       `json:"updatedAt"`
	UpdatedBy                string          `json:"updatedBy"`
	LoanPolicyTemplateName   string          `json:"loanPolicyTemplateName"`
	Source                   string          `json:"source"`
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

func (template LoanPolicyTemplate) ToAggregateModel(poolGroup MarginPoolGroup) AggregateLoanPolicyTemplate {
	return AggregateLoanPolicyTemplate{
		Id:                       template.Id,
		CreatedAt:                template.CreatedAt,
		UpdatedAt:                template.UpdatedAt,
		UpdatedBy:                template.UpdatedBy,
		Name:                     template.Name,
		InterestRate:             template.InterestRate,
		InterestBasis:            template.InterestBasis,
		Term:                     template.Term,
		PoolIdRef:                template.PoolIdRef,
		OverdueInterest:          template.OverdueInterest,
		AllowExtendLoanTerm:      template.AllowExtendLoanTerm,
		AllowEarlyPayment:        template.AllowEarlyPayment,
		PreferentialPeriod:       template.PreferentialPeriod,
		PreferentialInterestRate: template.PreferentialInterestRate,
		Source:                   poolGroup.Source,
	}
}

func (template AggregateLoanPolicyTemplate) ToSnapShotModel(shortTemplate LoanPolicyShorten) LoanPolicySnapShot {
	return LoanPolicySnapShot{
		AllowedOverdueLoanInDays: shortTemplate.AllowedOverdueLoanInDays,
		LoanPolicyTemplateId:     shortTemplate.LoanPolicyTemplateId,
		InitialRateForWithdraw:   shortTemplate.InitialRateForWithdraw,
		InitialRate:              shortTemplate.InitialRate,
		CreatedAt:                template.CreatedAt,
		UpdatedAt:                template.UpdatedAt,
		UpdatedBy:                template.UpdatedBy,
		LoanPolicyTemplateName:   template.Name,
		InterestRate:             template.InterestRate,
		InterestBasis:            template.InterestBasis,
		Term:                     template.Term,
		PoolIdRef:                template.PoolIdRef,
		OverdueInterest:          template.OverdueInterest,
		AllowExtendLoanTerm:      template.AllowExtendLoanTerm,
		AllowEarlyPayment:        template.AllowEarlyPayment,
		PreferentialPeriod:       template.PreferentialPeriod,
		PreferentialInterestRate: template.PreferentialInterestRate,
		Source:                   template.Source,
	}
}
