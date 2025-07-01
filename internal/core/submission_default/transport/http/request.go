package http

type SetSubmissionDefaultRequest struct {
	FirmSellingFeeRate       float64 `json:"firmSellingFeeRate"  binding:"required"`
	FirmBuyingFeeRate        float64 `json:"firmBuyingFeeRate"  binding:"required"`
	TransferFee              float64 `json:"transferFee"  binding:"required"`
	AllowedOverdueLoanInDays int64   `json:"allowedOverdueLoanInDays"  binding:"required"`
}
