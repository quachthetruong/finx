package http

type SuggestedOfferRequest struct {
	ConfigId  int64    `json:"configId" binding:"required"`
	Symbols   []string `json:"symbols" binding:"required,min=1"`
	AccountNo string   `json:"accountNo" binding:"required"`
}
