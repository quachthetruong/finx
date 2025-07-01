package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type SubmissionSheetMetadata struct {
	Id                   int64                 `json:"id"`
	LoanPackageRequestId int64                 `json:"loanPackageRequestId"`
	Creator              string                `json:"creator"`
	Status               SubmissionSheetStatus `json:"status"`
	FlowType             FlowType              `json:"flowType"`
	ActionType           ActionType            `json:"actionType"`
	ProposeType          ProposeType           `json:"proposeType"`
	CreatedAt            time.Time             `json:"createdAt"`
	UpdatedAt            time.Time             `json:"updatedAt"`
}

type SubmissionSheetDetail struct {
	Id                int64                `json:"id"`
	LoanRate          LoanRate             `json:"loanRate"`
	SubmissionSheetId int64                `json:"submissionSheetId"`
	FirmSellingFee    decimal.Decimal      `json:"firmSellingFee"`
	FirmBuyingFee     decimal.Decimal      `json:"firmBuyingFee"`
	TransferFee       decimal.Decimal      `json:"transferFee"`
	LoanPolicies      []LoanPolicySnapShot `json:"loanPolicies"`
	Comment           string               `json:"comment"`
}

type SubmissionSheet struct {
	Metadata SubmissionSheetMetadata `json:"metadata"`
	Detail   SubmissionSheetDetail   `json:"detail"`
}

type SubmissionSheetDetailShorten struct {
	Id                int64               `json:"id"`
	LoanPackageRateId int64               `json:"loanPackageRateId"`
	SubmissionSheetId int64               `json:"submissionSheetId"`
	FirmSellingFee    decimal.Decimal     `json:"firmSellingFee"`
	FirmBuyingFee     decimal.Decimal     `json:"firmBuyingFee"`
	TransferFee       decimal.Decimal     `json:"transferFee"`
	LoanPolicies      []LoanPolicyShorten `json:"loanPolicies"`
	Comment           string              `json:"comment"`
}

type SubmissionSheetShorten struct {
	Metadata SubmissionSheetMetadata      `json:"metadata"`
	Detail   SubmissionSheetDetailShorten `json:"detail"`
}

type SubmissionSheetStatus string

const (
	SubmissionSheetStatusSubmitted SubmissionSheetStatus = "SUBMITTED"
	SubmissionSheetStatusDraft     SubmissionSheetStatus = "DRAFT"
	SubmissionSheetStatusApproved  SubmissionSheetStatus = "APPROVED"
	SubmissionSheetStatusRejected  SubmissionSheetStatus = "REJECTED"
)

func (s SubmissionSheetStatus) String() string {
	return string(s)
}

func SubmissionSheetStatusFromString(str string) SubmissionSheetStatus {
	switch str {
	case string(SubmissionSheetStatusSubmitted):
		return SubmissionSheetStatusSubmitted
	case string(SubmissionSheetStatusDraft):
		return SubmissionSheetStatusDraft
	case string(SubmissionSheetStatusApproved):
		return SubmissionSheetStatusApproved
	case string(SubmissionSheetStatusRejected):
		return SubmissionSheetStatusRejected
	default:
		return SubmissionSheetStatusDraft
	}
}

func (s SubmissionSheetShorten) ToSubmissionSheet(loanRate LoanRate, loanPolicySnapShots []LoanPolicySnapShot) SubmissionSheet {
	return SubmissionSheet{
		Metadata: SubmissionSheetMetadata{
			Id:                   s.Metadata.Id,
			FlowType:             s.Metadata.FlowType,
			ActionType:           s.Metadata.ActionType,
			ProposeType:          s.Metadata.ProposeType,
			LoanPackageRequestId: s.Metadata.LoanPackageRequestId,
			Creator:              s.Metadata.Creator,
			Status:               s.Metadata.Status,
		},
		Detail: SubmissionSheetDetail{
			LoanRate:       loanRate,
			LoanPolicies:   loanPolicySnapShots,
			FirmBuyingFee:  s.Detail.FirmBuyingFee,
			FirmSellingFee: s.Detail.FirmSellingFee,
			TransferFee:    s.Detail.TransferFee,
			Comment:        s.Detail.Comment,
		},
	}
}

type SubmissionDefault struct {
	FirmSellingFeeRate       float64 `json:"firmSellingFeeRate"`
	FirmBuyingFeeRate        float64 `json:"firmBuyingFeeRate"`
	TransferFee              float64 `json:"transferFee"`
	AllowedOverdueLoanInDays int64   `json:"allowedOverdueLoanInDays"`
}
