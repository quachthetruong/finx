package entity

import (
	"time"
)

type SchedulerJob struct {
	Id           int64     `json:"id"`
	JobType      JobType   `json:"jobType"`
	JobStatus    JobStatus `json:"jobStatus"`
	TriggerBy    string    `json:"triggerBy"`
	TrackingData string    `json:"trackingData"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type JobStatus string

const (
	JobStatusSuccess JobStatus = "SUCCESS"
	JobStatusFail    JobStatus = "FAIL"
)

type JobType string

const (
	JobTypeDeclineHighRiskLoanRequest JobType = "JobTypeDeclineHighRiskLoanRequest"
)
