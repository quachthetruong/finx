package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/shopspring/decimal"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/loanpackagerequest/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/pkg/querymod"
)

var _ repository.LoanPackageRequestRepository = (*LoanPackageRequestPostgresRepository)(nil)

type LoanPackageRequestPostgresRepository struct {
	getDbFunc database.GetDbFunc
}

func (r *LoanPackageRequestPostgresRepository) Count(ctx context.Context, filter entity.LoanPackageFilter) (int64, error) {
	dest := struct {
		Count int64
	}{}
	stm := postgres.SELECT(postgres.COUNT(table.LoanPackageRequest.ID)).FROM(
		table.LoanPackageRequest.INNER_JOIN(
			table.Symbol, table.Symbol.ID.EQ(table.LoanPackageRequest.SymbolID),
		).LEFT_JOIN(table.Investor, table.Investor.InvestorID.EQ(table.LoanPackageRequest.InvestorID)),
	).WHERE(ApplyFilter(filter))
	if err := stm.QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		return 0, fmt.Errorf("LoanPackageRequestPostgresRepository Count: %w", err)
	}
	return dest.Count, nil
}

func (r *LoanPackageRequestPostgresRepository) GetAll(ctx context.Context, filter entity.LoanPackageFilter) ([]entity.LoanPackageRequest, error) {
	dest := make([]LoanPackageRequestWithAdditionalInfo, 0)
	stm := postgres.SELECT(table.LoanPackageRequest.AllColumns, table.Investor.AllColumns).FROM(
		table.LoanPackageRequest.INNER_JOIN(
			table.Symbol, table.Symbol.ID.EQ(table.LoanPackageRequest.SymbolID),
		).LEFT_JOIN(table.Investor, table.Investor.InvestorID.EQ(table.LoanPackageRequest.InvestorID)),
	).WHERE(ApplyFilter(filter))
	if limit := filter.Limit(); limit > 0 {
		stm = stm.OFFSET(filter.Offset()).LIMIT(limit)
	}
	if orderClause := ApplySort(filter); len(orderClause) > 0 {
		stm = stm.ORDER_BY(orderClause...)
	}
	if err := stm.QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []entity.LoanPackageRequest{}, nil
		}
		return nil, fmt.Errorf("LoanPackageRequestPostgresRepository GetAll: %w", err)
	}
	return MapLoanPackageRequestsWithAdditionalInfoDbToEntity(dest), nil
}

func (r *LoanPackageRequestPostgresRepository) GetAllUnderlyingRequests(ctx context.Context, filter entity.UnderlyingLoanPackageFilter) ([]entity.UnderlyingLoanPackageRequest, error) {
	dest := make([]UnderlyingRequest, 0)
	entities := make([]entity.UnderlyingLoanPackageRequest, 0)
	latestSubmissionCte := postgres.CTE("latest_submission")
	otherFilterExpressions := make([]postgres.BoolExpression, 0)
	submissionStatuses := make([]string, len(filter.SubmissionStatuses))

	if len(filter.SubmissionStatuses) > 0 {
		for _, t := range filter.SubmissionStatuses {
			submissionStatuses = append(submissionStatuses, t.String())
		}
		otherFilterExpressions = append(otherFilterExpressions, table.SubmissionSheetMetadata.Status.From(latestSubmissionCte).IN(querymod.In(submissionStatuses)...))
	}
	stm :=
		postgres.WITH(latestSubmissionCte.AS(
			postgres.SELECT(
				table.SubmissionSheetMetadata.LoanPackageRequestID,
				table.SubmissionSheetMetadata.ID,
				table.SubmissionSheetMetadata.Status,
				table.SubmissionSheetMetadata.CreatedAt,
				table.SubmissionSheetMetadata.Creator).
				DISTINCT(table.SubmissionSheetMetadata.LoanPackageRequestID).
				FROM(table.SubmissionSheetMetadata).
				ORDER_BY(table.SubmissionSheetMetadata.LoanPackageRequestID, table.SubmissionSheetMetadata.CreatedAt.DESC()),
		))(
			postgres.SELECT(
				table.LoanPackageRequest.AllColumns,
				table.Investor.AllColumns,
				table.SubmissionSheetMetadata.ID.From(latestSubmissionCte).AS("s.submission_id"),
				table.SubmissionSheetMetadata.Status.From(latestSubmissionCte).AS("s.submission_status"),
				table.SubmissionSheetMetadata.Creator.From(latestSubmissionCte).AS("s.submission_creator"),
				table.SubmissionSheetMetadata.CreatedAt.From(latestSubmissionCte).AS("s.submission_created_at"),
			).
				FROM(table.LoanPackageRequest.
					INNER_JOIN(
						table.Symbol,
						table.Symbol.ID.EQ(table.LoanPackageRequest.SymbolID),
					).
					LEFT_JOIN(
						table.Investor,
						table.Investor.InvestorID.EQ(table.LoanPackageRequest.InvestorID),
					).
					LEFT_JOIN(
						latestSubmissionCte,
						table.SubmissionSheetMetadata.LoanPackageRequestID.From(latestSubmissionCte).EQ(table.LoanPackageRequest.ID),
					),
				).
				WHERE(ApplyUnderlyingFilter(filter, otherFilterExpressions)).
				ORDER_BY(table.LoanPackageRequest.ID.DESC()))

	if err := stm.QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []entity.UnderlyingLoanPackageRequest{}, nil
		}
		return nil, fmt.Errorf("LoanPackageRequestPostgresRepository GetAll: %w", err)
	}

	for _, request := range dest {
		entities = append(entities, MapUnderlyingRequestDbToEntity(request))
	}

	return entities, nil
}

func (r *LoanPackageRequestPostgresRepository) GetById(ctx context.Context, id int64, filter entity.LoanPackageFilter, opts ...querymod.GetOption) (entity.LoanPackageRequest, error) {
	getQm := querymod.GetQm{}
	for _, opt := range opts {
		opt(&getQm)
	}
	var dest model.LoanPackageRequest
	stm := table.LoanPackageRequest.
		SELECT(table.LoanPackageRequest.AllColumns)
	stm = stm.WHERE(ApplyFilter(filter).AND(table.LoanPackageRequest.ID.EQ(postgres.Int64(id))))
	if getQm.ForUpdate {
		stm = stm.FOR(postgres.UPDATE())
	}
	if err := stm.QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf("LoanPackageRequestPostgresRepository GetById: %w", err)
	}
	return MapLoanPackageRequestDbToEntity(dest), nil
}

func (r *LoanPackageRequestPostgresRepository) Create(ctx context.Context, loanPackageRequest entity.LoanPackageRequest) (entity.LoanPackageRequest, error) {
	createModel := MapLoanPackageRequestEntityToDb(loanPackageRequest)
	created := model.LoanPackageRequest{}
	if err := table.LoanPackageRequest.INSERT(table.LoanPackageRequest.MutableColumns).
		MODEL(createModel).
		RETURNING(table.LoanPackageRequest.AllColumns).
		QueryContext(ctx, r.getDbFunc(ctx), &created); err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf("LoanPackageRequestPostgresRepository Create: %w", err)
	}
	return MapLoanPackageRequestDbToEntity(created), nil
}

func (r *LoanPackageRequestPostgresRepository) Update(ctx context.Context, loanPackageRequest entity.LoanPackageRequest) (entity.LoanPackageRequest, error) {
	updateModel := MapLoanPackageRequestEntityToDb(loanPackageRequest)
	updated := model.LoanPackageRequest{}
	if err := table.LoanPackageRequest.UPDATE(table.LoanPackageRequest.MutableColumns).
		MODEL(updateModel).
		WHERE(table.LoanPackageRequest.ID.EQ(postgres.Int64(loanPackageRequest.Id))).
		RETURNING(table.LoanPackageRequest.AllColumns).
		QueryContext(ctx, r.getDbFunc(ctx), &updated); err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf("LoanPackageRequestPostgresRepository Update: %w", err)
	}
	return MapLoanPackageRequestDbToEntity(updated), nil
}

func (r *LoanPackageRequestPostgresRepository) Delete(ctx context.Context, id int64) error {
	if _, err := table.LoanPackageRequest.DELETE().WHERE(table.LoanPackageRequest.ID.EQ(postgres.Int64(id))).ExecContext(
		ctx, r.getDbFunc(ctx),
	); err != nil {
		return fmt.Errorf("LoanPackageRequestPostgresRepository Delete: %w", err)
	}
	return nil
}

func (r *LoanPackageRequestPostgresRepository) LockAllPendingRequestByMaxPercent(
	ctx context.Context,
	maximumLoanRate decimal.Decimal,
) ([]entity.LoanPackageRequest, error) {
	dest := make([]model.LoanPackageRequest, 0)
	if err := table.LoanPackageRequest.
		SELECT(table.LoanPackageRequest.AllColumns).
		WHERE(
			table.LoanPackageRequest.LoanRate.
				GT_EQ(postgres.Decimal(maximumLoanRate.String())).
				AND(
					table.LoanPackageRequest.Status.
						EQ(postgres.String(entity.LoanPackageRequestStatusPending.String())),
				),
		).
		FOR(postgres.UPDATE().SKIP_LOCKED()).
		QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		return []entity.LoanPackageRequest{}, fmt.Errorf(
			"LoanPackageRequestPostgresRepository LockAllPendingRequestByMaxPercent: %w", err,
		)
	}
	return MapLoanPackageRequestsDbToEntity(dest), nil
}

func (r *LoanPackageRequestPostgresRepository) SaveLoggedRequest(ctx context.Context, request entity.LoggedRequest) (entity.LoggedRequest, error) {
	created := model.LoggedRequest{}
	if err := table.LoggedRequest.INSERT(table.LoggedRequest.MutableColumns).
		MODEL(MapLoggedRequestEntityToDb(request)).
		RETURNING(table.LoggedRequest.AllColumns).
		QueryContext(ctx, r.getDbFunc(ctx), &created); err != nil {
		return entity.LoggedRequest{}, fmt.Errorf("LoanPackageRequestPostgresRepository SaveLoggedRequest: %w", err)
	}
	return MapLoggedRequestDbToEntity(created), nil
}

func (r *LoanPackageRequestPostgresRepository) UpdateStatusByLoanRequestIds(ctx context.Context, loanRequestIds []int64, status entity.LoanPackageRequestStatus) ([]entity.LoanPackageRequest, error) {
	dest := make([]model.LoanPackageRequest, 0)
	if err := table.LoanPackageRequest.
		UPDATE(table.LoanPackageRequest.Status).
		SET(status).
		WHERE(table.LoanPackageRequest.ID.IN(querymod.In(loanRequestIds)...)).
		RETURNING(table.LoanPackageRequest.AllColumns).
		QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		return []entity.LoanPackageRequest{}, fmt.Errorf(
			"LoanPackageRequestPostgresRepository UpdateStatusByLoanRequestIds: %w", err,
		)
	}
	return MapLoanPackageRequestsDbToEntity(dest), nil
}

func (r *LoanPackageRequestPostgresRepository) LockAndReturnAllPendingRequestBySymbolId(ctx context.Context, symbolId int64) ([]entity.LoanPackageRequest, error) {
	dest := make([]model.LoanPackageRequest, 0)
	if err := table.LoanPackageRequest.
		SELECT(table.LoanPackageRequest.AllColumns).
		WHERE(
			table.LoanPackageRequest.Status.
				EQ(postgres.String(entity.LoanPackageRequestStatusPending.String())).
				AND(table.LoanPackageRequest.SymbolID.EQ(postgres.Int64(symbolId))),
		).
		FOR(postgres.UPDATE().SKIP_LOCKED()).
		QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		return []entity.LoanPackageRequest{}, fmt.Errorf(
			"LoanPackageRequestPostgresRepository LockAndReturnAllPendingRequestBySymbol: %w", err,
		)
	}
	return MapLoanPackageRequestsDbToEntity(dest), nil
}

func (r *LoanPackageRequestPostgresRepository) UpdateStatusById(ctx context.Context, id int64, status entity.LoanPackageRequestStatus) (entity.LoanPackageRequest, error) {
	updated := model.LoanPackageRequest{}
	if err := table.LoanPackageRequest.
		UPDATE(table.LoanPackageRequest.Status).
		SET(status).
		WHERE(table.LoanPackageRequest.ID.EQ(postgres.Int64(id))).
		RETURNING(table.LoanPackageRequest.AllColumns).
		QueryContext(ctx, r.getDbFunc(ctx), &updated); err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf("LoanPackageRequestPostgresRepository UpdateStatusById: %w", err)
	}
	return MapLoanPackageRequestDbToEntity(updated), nil
}

func NewLoanPackageRequestPostgresRepository(getDbFunc database.GetDbFunc) *LoanPackageRequestPostgresRepository {
	return &LoanPackageRequestPostgresRepository{
		getDbFunc: getDbFunc,
	}
}
