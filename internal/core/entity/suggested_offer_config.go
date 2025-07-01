package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type SuggestedOfferConfig struct {
	Id            int64                      `json:"id"`
	Name          string                     `json:"name"`
	Value         decimal.Decimal            `json:"value"`
	ValueType     ValueType                  `json:"valueType"`
	Status        SuggestedOfferConfigStatus `json:"status"`
	CreatedBy     string                     `json:"createdBy"`
	LastUpdatedBy string                     `json:"lastUpdatedBy"`
	CreatedAt     time.Time                  `json:"createdAt"`
	UpdatedAt     time.Time                  `json:"updatedAt"`
}

type SuggestedOfferConfigStatus string

const (
	SuggestedOfferConfigStatusActive   SuggestedOfferConfigStatus = "ACTIVE"
	SuggestedOfferConfigStatusInactive SuggestedOfferConfigStatus = "INACTIVE"
)

func (s SuggestedOfferConfigStatus) String() string {
	return string(s)
}

func SuggestionsOfferConfigStatusFromString(s string) SuggestedOfferConfigStatus {
	switch s {
	case "ACTIVE":
		return SuggestedOfferConfigStatusActive
	case "INACTIVE":
		return SuggestedOfferConfigStatusInactive
	default:
		return SuggestedOfferConfigStatusInactive
	}
}

type ValueType string

const (
	ValueTypeInterestRate ValueType = "INTEREST_RATE"
	ValueTypeLoanRate     ValueType = "LOAN_RATE"
)

func (v ValueType) String() string {
	return string(v)
}

func ValueTypeFromString(s string) ValueType {
	switch s {
	case "INTEREST_RATE":
		return ValueTypeInterestRate
	case "LOAN_RATE":
		return ValueTypeLoanRate
	default:
		return ValueTypeInterestRate
	}
}
