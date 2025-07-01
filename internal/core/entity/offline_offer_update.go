package entity

import (
	"time"
)

type OfflineOfferUpdate struct {
	Id        int64                    `json:"id"`
	OfferId   int64                    `json:"offerId"`
	Status    OfflineOfferUpdateStatus `json:"status"`
	Category  string                   `json:"category"`
	Note      string                   `json:"note"`
	CreatedBy string                   `json:"createdBy"`
	CreatedAt time.Time                `json:"createdAt"`
}

type OfflineOfferUpdateStatus string

const (
	OfflineOfferUpdateStatusProcessing OfflineOfferUpdateStatus = "PROCESSING"
	OfflineOfferUpdateStatusRejected   OfflineOfferUpdateStatus = "REJECTED"
	OfflineOfferUpdateStatusApproved   OfflineOfferUpdateStatus = "APPROVED"
)
