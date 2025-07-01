package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"financing-offer/internal/core/combined_loan_request/repository"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/pkg/querymod"
)

var _ repository.CombinedLoanPackageRequestPersistenceRepository = (*CombinedLoanPackageRequestPostgresRepository)(nil)

type CombinedLoanPackageRequestPostgresRepository struct {
	getDbFunc database.GetDbFunc
}

func (r *CombinedLoanPackageRequestPostgresRepository) GetAll(ctx context.Context, filter entity.CombinedLoanRequestFilter) ([]entity.CombinedLoanRequest, error) {
	dest := make([]CombinedLoanRequest, 0)
	stm := postgres.SELECT(
		table.LoanPackageRequest.AllColumns,
		table.LoanPackageOffer.AllColumns,
		table.Symbol.AllColumns,
		table.LoanContract.AllColumns,
		table.Investor.CustodyCode,
		querymod.ArrayAgg(
			postgres.NULLIF(table.LoanPackageOfferInterest.LoanID, postgres.Int(0)),
		).AS("r.admin_assigned_loan_package_ids"),
		querymod.ArrayAgg(
			postgres.CASE().
				WHEN(table.LoanPackageOfferInterest.Status.EQ(postgres.String(entity.LoanPackageOfferInterestStatusLoanPackageCreated.String()))).
				THEN(table.LoanPackageOfferInterest.LoanID),
		).AS("r.activated_loan_package_ids"),
		querymod.ArrayAgg(table.LoanContract.CreatedAt).AS("r.package_created_time"),
		querymod.ArrayAgg(table.LoanPackageOfferInterest.Status).AS("r.statuses"),
		querymod.ArrayAgg(table.LoanPackageOfferInterest.CancelledReason).AS("r.cancelled_reasons"),
	).FROM(
		table.LoanPackageRequest.LEFT_JOIN(
			table.LoanPackageOffer,
			table.LoanPackageOffer.LoanPackageRequestID.EQ(table.LoanPackageRequest.ID),
		).LEFT_JOIN(
			table.LoanPackageOfferInterest,
			table.LoanPackageOfferInterest.LoanPackageOfferID.EQ(table.LoanPackageOffer.ID),
		).LEFT_JOIN(
			table.LoanContract, table.LoanContract.LoanOfferInterestID.EQ(table.LoanPackageOfferInterest.ID),
		).INNER_JOIN(
			table.Symbol,
			table.Symbol.ID.EQ(table.LoanPackageRequest.SymbolID),
		).LEFT_JOIN(table.Investor, table.Investor.InvestorID.EQ(table.LoanPackageRequest.InvestorID)),
	).WHERE(ApplyWhere(filter)).HAVING(ApplyHaving(filter)).GROUP_BY(postgres.WRAP(groupColumns()...))
	if limit := filter.Limit(); limit > 0 {
		stm = stm.LIMIT(limit).OFFSET(filter.Offset())
	}
	if err := stm.ORDER_BY(table.LoanPackageRequest.ID.DESC()).QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		return nil, fmt.Errorf("CombinedLoanPackageRequestPostgresRepository GetAll: %w", err)
	}
	res, err := MapCombinedLoanPackageRequestsDbToEntity(dest)
	if err != nil {
		return nil, fmt.Errorf("CombinedLoanPackageRequestPostgresRepository GetAll: %w", err)
	}
	return res, nil
}

func groupColumns() []postgres.Expression {
	columns := make(
		[]postgres.Expression, 0,
		len(table.LoanPackageRequest.AllColumns)+len(table.LoanPackageOffer.AllColumns)+len(table.Symbol.AllColumns)+len(table.LoanContract.AllColumns)+len(table.LoanContract.AllColumns),
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
	for _, c := range table.LoanContract.AllColumns {
		columns = append(columns, c)
	}
	columns = append(columns, table.Investor.CustodyCode)
	return columns
}

func (r *CombinedLoanPackageRequestPostgresRepository) Count(ctx context.Context, filter entity.CombinedLoanRequestFilter) (int64, error) {
	dest := struct {
		Count int64
	}{}
	s := postgres.SELECT(
		table.LoanPackageRequest.ID,
	).FROM(
		table.LoanPackageRequest.LEFT_JOIN(
			table.LoanPackageOffer,
			table.LoanPackageOffer.LoanPackageRequestID.EQ(table.LoanPackageRequest.ID),
		).LEFT_JOIN(
			table.LoanPackageOfferInterest,
			table.LoanPackageOfferInterest.LoanPackageOfferID.EQ(table.LoanPackageOffer.ID),
		).LEFT_JOIN(
			table.LoanContract, table.LoanContract.LoanOfferInterestID.EQ(table.LoanPackageOfferInterest.ID),
		).INNER_JOIN(
			table.Symbol,
			table.Symbol.ID.EQ(table.LoanPackageRequest.SymbolID),
		).LEFT_JOIN(table.Investor, table.Investor.InvestorID.EQ(table.LoanPackageRequest.InvestorID)),
	).WHERE(ApplyWhere(filter)).GROUP_BY(table.LoanPackageRequest.ID).HAVING(ApplyHaving(filter)).AsTable("q")
	stm := postgres.SELECT(postgres.COUNT(postgres.String("*"))).FROM(s)
	if err := stm.QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("CombinedLoanPackageRequestPostgresRepository Count: %w", err)
	}
	return dest.Count, nil
}

func NewCombinedLoanPackageRequestPostgresRepository(getDbFunc database.GetDbFunc) *CombinedLoanPackageRequestPostgresRepository {
	return &CombinedLoanPackageRequestPostgresRepository{getDbFunc: getDbFunc}
}
