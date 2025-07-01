package http

type SetLoanRateRequest struct {
	Ids []int64 `json:"ids"  binding:"required"`
}

type SetMarginPoolRequest struct {
	Ids []int64 `json:"ids"  binding:"required"`
}
