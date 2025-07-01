package entity

import "time"

type InvestorAccount struct {
	AccountNo    string       `json:"accountNo"`
	InvestorId   string       `json:"investorId"`
	MarginStatus MarginStatus `json:"marginStatus"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

type MarginStatus string

const (
	MarginStatusVersion2 MarginStatus = "v2"
	MarginStatusVersion3 MarginStatus = "v3"
)

func (v MarginStatus) String() string {
	return string(v)
}

func MarginStatusFromString(s string) MarginStatus {
	switch s {
	case "v2":
		return MarginStatusVersion2
	case "v3":
		return MarginStatusVersion3
	default:
		return ""
	}
}
