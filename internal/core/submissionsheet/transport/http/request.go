package http

type ApproveSubmissionRequest struct {
	SubmissionId int64 `json:"submissionId"`
}

type RejectSubmissionRequest struct {
	SubmissionId int64 `json:"submissionId"`
}
