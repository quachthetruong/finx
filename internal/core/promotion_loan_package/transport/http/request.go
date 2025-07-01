package http

import (
	"financing-offer/internal/core/entity"
	"financing-offer/internal/funcs"
)

type GetPromotionLoanPackageBySymbolRequest struct {
	AccountNo string `form:"accountNo"`
}

type GetInvestorPromotionLoanPackagesRequest struct {
	AccountNo string `form:"accountNo"`
}

type GetPromotionLoanPackagesRequest struct {
	AccountNo string `form:"accountNo"`
	Symbol    string `form:"symbol"`
}

type GetPublicLoanPackagesRequest struct {
	Symbol string `form:"symbol"`
}

type SetPromotionLoanPackagesRequest struct {
	LoanProducts []PromotionLoanProductRequest `json:"loanProducts"  binding:"required,dive"`
}

type PromotionLoanProductRequest struct {
	LoanPackageId    int64    `json:"loanPackageId"  binding:"required"`
	RetailSymbols    []string `json:"retailSymbols" binding:"required"`
	NonRetailSymbols []string `json:"nonRetailSymbols"  binding:"required"`
}

func (r SetPromotionLoanPackagesRequest) toEntity() entity.PromotionLoanPackage {
	return entity.PromotionLoanPackage{
		LoanProducts: func() []entity.PromotionLoanProduct {
			var loanProducts []entity.PromotionLoanProduct
			for _, loanProduct := range r.LoanProducts {
				loanProducts = append(loanProducts, loanProduct.toEntity())
			}
			return loanProducts
		}(),
	}
}

func (r PromotionLoanProductRequest) toEntity() entity.PromotionLoanProduct {
	return entity.PromotionLoanProduct{
		LoanPackageId:    r.LoanPackageId,
		Symbols:          funcs.UniqueElements(r.RetailSymbols, r.NonRetailSymbols),
		RetailSymbols:    r.RetailSymbols,
		NonRetailSymbols: r.NonRetailSymbols,
	}

}
