package postgres

import (
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

func MapLoanPolicyTemplatesDbToEntities(policies []model.LoanPolicyTemplate) []entity.LoanPolicyTemplate {
	dest := make([]entity.LoanPolicyTemplate, 0, len(policies))
	for _, p := range policies {
		dest = append(dest, MapLoanPolicyTemplateDbToEntity(p))
	}
	return dest
}

func MapLoanPolicyTemplateDbToEntity(loanPolicy model.LoanPolicyTemplate) entity.LoanPolicyTemplate {
	return entity.LoanPolicyTemplate{
		Id:                       loanPolicy.ID,
		CreatedAt:                loanPolicy.CreatedAt,
		UpdatedAt:                loanPolicy.UpdatedAt,
		UpdatedBy:                loanPolicy.UpdatedBy,
		Name:                     loanPolicy.Name,
		InterestRate:             loanPolicy.InterestRate,
		InterestBasis:            loanPolicy.InterestBasis,
		Term:                     loanPolicy.Term,
		PoolIdRef:                loanPolicy.PoolIDRef,
		OverdueInterest:          loanPolicy.OverdueInterest,
		AllowExtendLoanTerm:      loanPolicy.AllowExtendLoanTerm,
		AllowEarlyPayment:        loanPolicy.AllowEarlyPayment,
		PreferentialPeriod:       loanPolicy.PreferentialPeriod,
		PreferentialInterestRate: loanPolicy.PreferentialInterestRate,
	}
}

func MapLoanPolicyTemplateEntityToDb(loanPolicy entity.LoanPolicyTemplate) model.LoanPolicyTemplate {
	return model.LoanPolicyTemplate{
		ID:                       loanPolicy.Id,
		CreatedAt:                loanPolicy.CreatedAt,
		UpdatedAt:                loanPolicy.UpdatedAt,
		UpdatedBy:                loanPolicy.UpdatedBy,
		Name:                     loanPolicy.Name,
		InterestRate:             loanPolicy.InterestRate,
		InterestBasis:            loanPolicy.InterestBasis,
		Term:                     loanPolicy.Term,
		PoolIDRef:                loanPolicy.PoolIdRef,
		OverdueInterest:          loanPolicy.OverdueInterest,
		AllowExtendLoanTerm:      loanPolicy.AllowExtendLoanTerm,
		AllowEarlyPayment:        loanPolicy.AllowEarlyPayment,
		PreferentialPeriod:       loanPolicy.PreferentialPeriod,
		PreferentialInterestRate: loanPolicy.PreferentialInterestRate,
	}
}
