package http

import (
	"github.com/shopspring/decimal"

	"financing-offer/internal/core/entity"
)

type CreateSuggestedOfferConfigRequest struct {
	Name      string          `json:"name" binding:"required"`
	Value     decimal.Decimal `json:"value" binding:"required"`
	ValueType string          `json:"valueType" binding:"required,oneof=INTEREST_RATE LOAN_RATE"`
}

func (r *CreateSuggestedOfferConfigRequest) toEntity() entity.SuggestedOfferConfig {
	return entity.SuggestedOfferConfig{
		Name:      r.Name,
		Value:     r.Value,
		ValueType: entity.ValueTypeFromString(r.ValueType),
		Status:    entity.SuggestedOfferConfigStatusInactive,
	}
}

type UpdateSuggestedOfferConfigStatusRequest struct {
	Status entity.SuggestedOfferConfigStatus `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}
