package postgres

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"

	"financing-offer/internal/core/awaiting_confirm_request/repository"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/internal/database/dbmodels/finoffer/public/view"
)

var _ repository.AwaitingConfirmRequestPersistenceRepository = (*AwaitingConfirmRequestPostgresRepository)(nil)

type AwaitingConfirmRequestPostgresRepository struct {
	getDbFunc database.GetDbFunc
}

func (r *AwaitingConfirmRequestPostgresRepository) GetAll(ctx context.Context, filter entity.AwaitingConfirmRequestFilter) ([]entity.AwaitingConfirmRequest, error) {
	dest := make([]AwaitingConfirmRequest, 0)
	stm := postgres.SELECT(
		table.LoanPackageRequest.AllColumns,
		table.LoanPackageOffer.AllColumns,
		table.Symbol.AllColumns,
		view.LatestOfferUpdate.AllColumns,
		table.Investor.CustodyCode,
		postgres.Func(
			"string_agg",
			postgres.NULLIF(postgres.CAST(table.LoanPackageOfferInterest.LoanID).AS_TEXT(), postgres.String("0")),
			postgres.String(", "),
		).AS("r.loan_package_ids"),
	).
		FROM(
			table.LoanPackageRequest.INNER_JOIN(
				table.LoanPackageOffer,
				table.LoanPackageOffer.LoanPackageRequestID.EQ(table.LoanPackageRequest.ID),
			).INNER_JOIN(
				table.LoanPackageOfferInterest,
				table.LoanPackageOfferInterest.LoanPackageOfferID.EQ(table.LoanPackageOffer.ID).
					AND(table.LoanPackageOfferInterest.Status.EQ(postgres.String(entity.LoanPackageOfferInterestStatusPending.String()))),
			).INNER_JOIN(
				table.Symbol,
				table.Symbol.ID.EQ(table.LoanPackageRequest.SymbolID),
			).LEFT_JOIN(
				view.LatestOfferUpdate,
				view.LatestOfferUpdate.OfferID.EQ(table.LoanPackageOffer.ID),
			).LEFT_JOIN(table.Investor, table.Investor.InvestorID.EQ(table.LoanPackageRequest.InvestorID)),
		).WHERE(ApplyFilter(filter)).GROUP_BY(postgres.WRAP(groupColumns()...))
	if limit := filter.Limit(); limit > 0 {
		stm = stm.LIMIT(limit).OFFSET(filter.Offset())
	}
	if orderClause := ApplySort(filter); len(orderClause) > 0 {
		stm = stm.ORDER_BY(orderClause...)
	}
	if err := stm.QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		return nil, fmt.Errorf("AwaitingConfirmRequestPostgresRepository GetAll %w", err)
	}
	return MapAwaitingConfirmRequestsDbToEntity(dest), nil
}

func groupColumns() []postgres.Expression {
	columns := make(
		[]postgres.Expression, 0,
		len(table.LoanPackageRequest.AllColumns)+len(table.LoanPackageOffer.AllColumns)+len(view.LatestOfferUpdate.AllColumns)+len(table.Symbol.AllColumns),
	)
	for _, c := range table.LoanPackageRequest.AllColumns {
		columns = append(columns, c)
	}
	for _, c := range table.LoanPackageOffer.AllColumns {
		columns = append(columns, c)
	}
	for _, c := range table.Symbol.AllColumns {
		columns = append(columns, c)
	}
	for _, c := range view.LatestOfferUpdate.AllColumns {
		columns = append(columns, c)
	}
	columns = append(columns, table.Investor.CustodyCode)

	return columns
}

func (r *AwaitingConfirmRequestPostgresRepository) Count(ctx context.Context, filter entity.AwaitingConfirmRequestFilter) (int64, error) {
	dest := struct {
		Count int64
	}{}
	stm := postgres.SELECT(postgres.COUNT(postgres.DISTINCT(table.LoanPackageOffer.ID))).
		FROM(
			table.LoanPackageRequest.INNER_JOIN(
				table.LoanPackageOffer,
				table.LoanPackageOffer.LoanPackageRequestID.EQ(table.LoanPackageRequest.ID),
			).INNER_JOIN(
				table.LoanPackageOfferInterest,
				table.LoanPackageOfferInterest.LoanPackageOfferID.EQ(table.LoanPackageOffer.ID).
					AND(table.LoanPackageOfferInterest.Status.EQ(postgres.String(entity.LoanPackageOfferInterestStatusPending.String()))),
			).INNER_JOIN(
				table.Symbol,
				table.Symbol.ID.EQ(table.LoanPackageRequest.SymbolID),
			).LEFT_JOIN(
				view.LatestOfferUpdate,
				view.LatestOfferUpdate.OfferID.EQ(table.LoanPackageOffer.ID),
			).LEFT_JOIN(table.Investor, table.Investor.InvestorID.EQ(table.LoanPackageRequest.InvestorID)),
		).WHERE(ApplyFilter(filter))
	if err := stm.QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		return 0, fmt.Errorf("AwaitingConfirmRequestPostgresRepository Count %w", err)
	}
	return dest.Count, nil
}

func (r *AwaitingConfirmRequestPostgresRepository) CountForStatistic(ctx context.Context) (int64, error) {
	dest := struct {
		Count int64
	}{}
	stm := postgres.SELECT(postgres.COUNT(postgres.DISTINCT(table.LoanPackageOffer.ID))).
		FROM(
			table.LoanPackageRequest.INNER_JOIN(
				table.LoanPackageOffer,
				table.LoanPackageOffer.LoanPackageRequestID.EQ(table.LoanPackageRequest.ID),
			).INNER_JOIN(
				table.LoanPackageOfferInterest,
				table.LoanPackageOfferInterest.LoanPackageOfferID.EQ(table.LoanPackageOffer.ID).
					AND(table.LoanPackageOfferInterest.Status.EQ(postgres.String(entity.LoanPackageOfferInterestStatusPending.String()))),
			).INNER_JOIN(
				table.Symbol,
				table.Symbol.ID.EQ(table.LoanPackageRequest.SymbolID),
			).LEFT_JOIN(
				view.LatestOfferUpdate,
				view.LatestOfferUpdate.OfferID.EQ(table.LoanPackageOffer.ID),
			),
		)
	if err := stm.QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		return 0, fmt.Errorf("AwaitingConfirmRequestPostgresRepository Count %w", err)
	}
	return dest.Count, nil
}

func NewAwaitingConfirmRequestPostgresRepository(getDbFunc database.GetDbFunc) *AwaitingConfirmRequestPostgresRepository {
	return &AwaitingConfirmRequestPostgresRepository{getDbFunc: getDbFunc}
}
