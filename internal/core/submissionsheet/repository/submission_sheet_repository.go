package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type SubmissionSheetRepository interface {
	GetDetailById(ctx context.Context, id int64) (entity.SubmissionSheetDetail, error)
	GetLatestByRequestId(ctx context.Context, requestId int64) (entity.SubmissionSheet, error)
	GetMetadataByRequestId(ctx context.Context, requestId int64) ([]entity.SubmissionSheetMetadata, error)
	CreateMetadata(ctx context.Context, submissionSheetMetadata entity.SubmissionSheetMetadata) (entity.SubmissionSheetMetadata, error)
	CreateDetail(ctx context.Context, submissionSheetDetail entity.SubmissionSheetDetail) (entity.SubmissionSheetDetail, error)
	UpdateMetadata(ctx context.Context, submissionSheetMetadata entity.SubmissionSheetMetadata) (entity.SubmissionSheetMetadata, error)
	UpdateDetail(ctx context.Context, submissionSheetDetail entity.SubmissionSheetDetail) (entity.SubmissionSheetDetail, error)
	GetById(ctx context.Context, id int64) (entity.SubmissionSheet, error)
	UpdateMetadataStatusById(ctx context.Context, id int64, status entity.SubmissionSheetStatus) error
}
