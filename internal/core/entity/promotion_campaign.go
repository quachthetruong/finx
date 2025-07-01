package entity

import "time"

type GetPromotionCampaignsRequest struct {
	Status string `form:"status"`
}

type PromotionCampaign struct {
	Id          int64                     `json:"id"`
	CreatedAt   time.Time                 `json:"createdAt"`
	UpdatedAt   time.Time                 `json:"updatedAt"`
	UpdatedBy   string                    `json:"updatedBy"`
	Name        string                    `json:"name"`
	Tag         string                    `json:"tag"`
	Description string                    `json:"description"`
	Status      PromotionCampaignStatus   `json:"status"`
	Metadata    PromotionCampaignMetadata `json:"metadata"`
}

type PromotionCampaignMetadata struct {
	Products []PromotionCampaignProduct `json:"products"`
}

type PromotionCampaignProduct struct {
	Symbols       []string `json:"symbols"`
	LoanPackageId int64    `json:"loanPackageId"`
	RetailSymbols []string `json:"retailSymbols"`
}
