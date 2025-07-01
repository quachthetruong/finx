package postgres

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/submissionsheet/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

var _ repository.SubmissionSheetRepository = (*SubmissionSheetPostgresRepository)(nil)

type SubmissionSheetPostgresRepository struct {
	getDbFunc database.GetDbFunc
}

func (r *SubmissionSheetPostgresRepository) GetDetailById(ctx context.Context, id int64) (entity.SubmissionSheetDetail, error) {
	errorTemplate := "SubmissionSheetPostgresRepository GetDetailById %w"
	submissionSheetDetail := model.SubmissionSheetDetail{}
	err := table.SubmissionSheetDetail.
		SELECT(table.SubmissionSheetDetail.AllColumns).
		WHERE(table.SubmissionSheetDetail.ID.EQ(postgres.Int64(id))).QueryContext(
		ctx, r.getDbFunc(ctx), &submissionSheetDetail,
	)
	if err != nil {
		return entity.SubmissionSheetDetail{}, fmt.Errorf(errorTemplate, err)
	}
	res, err := MapSubmissionSheetDetailDbToEntity(submissionSheetDetail)
	if err != nil {
		return entity.SubmissionSheetDetail{}, fmt.Errorf(errorTemplate, err)
	}
	return res, nil
}

func (r *SubmissionSheetPostgresRepository) GetLatestByRequestId(ctx context.Context, requestId int64) (entity.SubmissionSheet, error) {
	errorTemplate := "SubmissionSheetPostgresRepository GetByRequestId %w"
	submissionSheetMetadata := model.SubmissionSheetMetadata{}
	submissionSheetDetail := model.SubmissionSheetDetail{}
	err := table.SubmissionSheetMetadata.
		SELECT(table.SubmissionSheetMetadata.AllColumns).
		WHERE(table.SubmissionSheetMetadata.LoanPackageRequestID.EQ(postgres.Int64(requestId))).
		ORDER_BY(table.SubmissionSheetMetadata.CreatedAt.DESC()).
		LIMIT(1).
		QueryContext(
			ctx, r.getDbFunc(ctx), &submissionSheetMetadata,
		)
	if err != nil {
		return entity.SubmissionSheet{}, fmt.Errorf(errorTemplate, err)
	}
	err = table.SubmissionSheetDetail.
		SELECT(table.SubmissionSheetDetail.AllColumns).
		WHERE(table.SubmissionSheetDetail.SubmissionSheetID.EQ(postgres.Int64(submissionSheetMetadata.ID))).QueryContext(
		ctx, r.getDbFunc(ctx), &submissionSheetDetail,
	)
	if err != nil {
		return entity.SubmissionSheet{}, fmt.Errorf(errorTemplate, err)
	}
	metadata := MapSubmissionSheetMetadataDbToEntity(submissionSheetMetadata)
	detail, err := MapSubmissionSheetDetailDbToEntity(submissionSheetDetail)
	if err != nil {
		return entity.SubmissionSheet{}, fmt.Errorf(errorTemplate, err)
	}
	return entity.SubmissionSheet{
		Metadata: metadata,
		Detail:   detail,
	}, nil
}

func (r *SubmissionSheetPostgresRepository) GetMetadataByRequestId(ctx context.Context, requestId int64) ([]entity.SubmissionSheetMetadata, error) {
	errorTemplate := "SubmissionSheetPostgresRepository GetMetadataByRequestId %w"
	dest := make([]model.SubmissionSheetMetadata, 0)
	err := table.SubmissionSheetMetadata.
		SELECT(table.SubmissionSheetMetadata.AllColumns).
		WHERE(table.SubmissionSheetMetadata.LoanPackageRequestID.EQ(postgres.Int64(requestId))).
		ORDER_BY(table.SubmissionSheetMetadata.CreatedAt.DESC()).
		QueryContext(
			ctx, r.getDbFunc(ctx), &dest,
		)
	if err != nil {
		return []entity.SubmissionSheetMetadata{}, fmt.Errorf(errorTemplate, err)
	}
	entities := make([]entity.SubmissionSheetMetadata, 0, len(dest))
	for _, metadata := range dest {
		entities = append(entities, MapSubmissionSheetMetadataDbToEntity(metadata))
	}
	return entities, nil
}

func (r *SubmissionSheetPostgresRepository) GetById(ctx context.Context, id int64) (entity.SubmissionSheet, error) {
	errorTemplate := "SubmissionSheetPostgresRepository GetById %w"
	submissionSheetMetadata := model.SubmissionSheetMetadata{}
	submissionSheetDetail := model.SubmissionSheetDetail{}
	err := table.SubmissionSheetMetadata.
		SELECT(table.SubmissionSheetMetadata.AllColumns).
		WHERE(table.SubmissionSheetMetadata.ID.EQ(postgres.Int64(id))).QueryContext(
		ctx, r.getDbFunc(ctx), &submissionSheetMetadata,
	)
	if err != nil {
		return entity.SubmissionSheet{}, fmt.Errorf(errorTemplate, err)
	}
	err = table.SubmissionSheetDetail.
		SELECT(table.SubmissionSheetDetail.AllColumns).
		WHERE(table.SubmissionSheetDetail.SubmissionSheetID.EQ(postgres.Int64(id))).QueryContext(
		ctx, r.getDbFunc(ctx), &submissionSheetDetail,
	)
	if err != nil {
		return entity.SubmissionSheet{}, fmt.Errorf(errorTemplate, err)
	}
	metadata := MapSubmissionSheetMetadataDbToEntity(submissionSheetMetadata)
	detail, err := MapSubmissionSheetDetailDbToEntity(submissionSheetDetail)
	if err != nil {
		return entity.SubmissionSheet{}, fmt.Errorf(errorTemplate, err)
	}
	return entity.SubmissionSheet{
		Metadata: metadata,
		Detail:   detail,
	}, nil
}

func (r *SubmissionSheetPostgresRepository) UpdateMetadataStatusById(ctx context.Context, id int64, status entity.SubmissionSheetStatus) error {
	errorTemplate := "SubmissionSheetPostgresRepository UpdateMetadataStatusById %w"
	_, err := table.SubmissionSheetMetadata.
		UPDATE(table.SubmissionSheetMetadata.Status).
		SET(postgres.String(status.String())).
		WHERE(table.SubmissionSheetMetadata.ID.EQ(postgres.Int64(id))).
		ExecContext(ctx, r.getDbFunc(ctx))
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	return nil
}

func (r *SubmissionSheetPostgresRepository) CreateMetadata(ctx context.Context, submissionSheetMetadata entity.SubmissionSheetMetadata) (entity.SubmissionSheetMetadata, error) {
	createModel := MapSubmissionSheetMetadataEntityToDb(submissionSheetMetadata)
	created := model.SubmissionSheetMetadata{}
	err := table.SubmissionSheetMetadata.INSERT(table.SubmissionSheetMetadata.MutableColumns).MODEL(createModel).RETURNING(table.SubmissionSheetMetadata.AllColumns).QueryContext(
		ctx, r.getDbFunc(ctx), &created,
	)
	if err != nil {
		return entity.SubmissionSheetMetadata{}, fmt.Errorf("LoanPackageOfferPostgresRepository Create %w", err)
	}
	return MapSubmissionSheetMetadataDbToEntity(created), nil
}

func (r *SubmissionSheetPostgresRepository) CreateDetail(ctx context.Context, submissionSheetDetail entity.SubmissionSheetDetail) (entity.SubmissionSheetDetail, error) {
	errorTemplate := "SubmissionSheetPostgresRepository CreateDetail %w"
	createModel, err := MapSubmissionSheetDetailEntityToDb(submissionSheetDetail)
	if err != nil {
		return entity.SubmissionSheetDetail{}, fmt.Errorf(errorTemplate, err)
	}
	created := model.SubmissionSheetDetail{}
	err = table.SubmissionSheetDetail.INSERT(table.SubmissionSheetDetail.MutableColumns).MODEL(createModel).RETURNING(table.SubmissionSheetDetail.AllColumns).QueryContext(
		ctx, r.getDbFunc(ctx), &created,
	)
	if err != nil {
		return entity.SubmissionSheetDetail{}, fmt.Errorf(errorTemplate, err)
	}
	res, err := MapSubmissionSheetDetailDbToEntity(created)
	if err != nil {
		return entity.SubmissionSheetDetail{}, fmt.Errorf(errorTemplate, err)
	}
	return res, nil
}

func (r *SubmissionSheetPostgresRepository) UpdateMetadata(ctx context.Context, submissionSheetMetadata entity.SubmissionSheetMetadata) (entity.SubmissionSheetMetadata, error) {
	updateModel := MapSubmissionSheetMetadataEntityToDb(submissionSheetMetadata)
	updated := model.SubmissionSheetMetadata{}
	err := table.SubmissionSheetMetadata.
		UPDATE(table.SubmissionSheetMetadata.MutableColumns).
		MODEL(updateModel).
		WHERE(table.SubmissionSheetMetadata.ID.EQ(postgres.Int64(submissionSheetMetadata.Id))).
		RETURNING(table.SubmissionSheetMetadata.AllColumns).QueryContext(
		ctx, r.getDbFunc(ctx), &updated,
	)
	if err != nil {
		return entity.SubmissionSheetMetadata{}, fmt.Errorf("LoanPackageOfferPostgresRepository Create %w", err)
	}
	return MapSubmissionSheetMetadataDbToEntity(updated), nil
}

func (r *SubmissionSheetPostgresRepository) UpdateDetail(ctx context.Context, submissionSheetDetail entity.SubmissionSheetDetail) (entity.SubmissionSheetDetail, error) {
	errorTemplate := "SubmissionSheetPostgresRepository UpdateDetail %w"
	updateModel, err := MapSubmissionSheetDetailEntityToDb(submissionSheetDetail)
	if err != nil {
		return entity.SubmissionSheetDetail{}, fmt.Errorf(errorTemplate, err)
	}
	updated := model.SubmissionSheetDetail{}
	err = table.SubmissionSheetDetail.
		UPDATE(table.SubmissionSheetDetail.MutableColumns).
		MODEL(updateModel).
		WHERE(table.SubmissionSheetDetail.SubmissionSheetID.EQ(postgres.Int64(submissionSheetDetail.SubmissionSheetId))).
		RETURNING(table.SubmissionSheetDetail.AllColumns).QueryContext(
		ctx, r.getDbFunc(ctx), &updated)
	if err != nil {
		return entity.SubmissionSheetDetail{}, fmt.Errorf("LoanPackageOfferPostgresRepository Create %w", err)
	}
	res, err := MapSubmissionSheetDetailDbToEntity(updated)
	if err != nil {
		return entity.SubmissionSheetDetail{}, fmt.Errorf(errorTemplate, err)
	}
	return res, nil
}

func NewSubmissionSheetPostgresRepository(getDbFunc database.GetDbFunc) *SubmissionSheetPostgresRepository {
	return &SubmissionSheetPostgresRepository{
		getDbFunc: getDbFunc,
	}
}
