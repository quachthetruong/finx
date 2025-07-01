package entity

type PromotionLoanPackage struct {
	LoanProducts []PromotionLoanProduct `json:"loanProducts"`
}

type LoanRateConfiguration struct {
	Ids []int64 `json:"ids"`
}

type MarginPoolConfiguration struct {
	Ids []int64 `json:"ids"`
}

type PromotionLoanProduct struct {
	LoanPackageId    int64    `json:"loanPackageId"`
	Symbols          []string `json:"symbols"`
	RetailSymbols    []string `json:"retailSymbols"`
	NonRetailSymbols []string `json:"nonRetailSymbols"`
}

func (p PromotionLoanProduct) AllSymbols() []string {
	return append(p.RetailSymbols, p.NonRetailSymbols...)
}
