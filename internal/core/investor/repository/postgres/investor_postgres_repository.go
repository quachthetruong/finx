package postgres

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/investor/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

var _ repository.InvestorPersistenceRepository = (*InvestorPostgresRepository)(nil)

type InvestorPostgresRepository struct {
	getDbFunc database.GetDbFunc
}

func (r *InvestorPostgresRepository) CreateIfNotExist(ctx context.Context, investor entity.Investor) error {
	toCreate := MapInvestorEntityToDb(investor)
	if _, err := table.Investor.
		INSERT(table.Investor.AllColumns.Except(table.Investor.CreatedAt, table.Investor.UpdatedAt)).
		MODEL(toCreate).
		ON_CONFLICT(table.Investor.InvestorID).
		DO_NOTHING().
		ExecContext(ctx, r.getDbFunc(ctx)); err != nil {
		return fmt.Errorf("InvestorPostgresRepository CreateIfNotExist %w", err)
	}
	return nil
}

func (r *InvestorPostgresRepository) GetAllInvestorIdsForMigration(ctx context.Context) ([]string, error) {
	var investorIds []string
	if err := table.Investor.
		SELECT(table.Investor.InvestorID).
		WHERE(table.Investor.CustodyCode.EQ(postgres.String(""))).
		QueryContext(ctx, r.getDbFunc(ctx), &investorIds); err != nil {
		return nil, fmt.Errorf("InvestorPostgresRepository GetAllInvestorIdsForMigration %w", err)
	}
	return investorIds, nil
}

func (r *InvestorPostgresRepository) GetAllUniqueInvestorIdsFromRequests(ctx context.Context) ([]string, error) {
	var investorIds []string
	if err := table.LoanPackageRequest.
		SELECT(table.LoanPackageRequest.InvestorID).
		DISTINCT().
		QueryContext(ctx, r.getDbFunc(ctx), &investorIds); err != nil {
		return nil, fmt.Errorf("InvestorPostgresRepository GetAllUniqueInvestorIdsFromRequests %w", err)
	}
	return investorIds, nil
}

func (r *InvestorPostgresRepository) BulkCreate(ctx context.Context, investors []entity.Investor) error {
	toCreateModels := MapInvestorEntitiesToDb(investors)
	if _, err := table.Investor.
		INSERT(table.Investor.AllColumns.Except(table.Investor.CreatedAt, table.Investor.UpdatedAt)).
		MODELS(toCreateModels).
		ON_CONFLICT(table.Investor.InvestorID).
		DO_NOTHING().
		ExecContext(ctx, r.getDbFunc(ctx)); err != nil {
		return fmt.Errorf("InvestorPostgresRepository BulkCreate %w", err)
	}
	return nil
}

func (r *InvestorPostgresRepository) Update(ctx context.Context, investor entity.Investor) (entity.Investor, error) {
	updatedModel := MapInvestorEntityToDb(investor)
	if _, err := table.Investor.
		UPDATE(table.Investor.MutableColumns).
		MODEL(updatedModel).
		WHERE(table.Investor.InvestorID.EQ(postgres.String(investor.InvestorId))).
		ExecContext(ctx, r.getDbFunc(ctx)); err != nil {
		return entity.Investor{}, fmt.Errorf("InvestorPostgresRepository Update %w", err)
	}
	return investor, nil
}

func NewInvestorPostgresRepository(getDbFunc database.GetDbFunc) *InvestorPostgresRepository {
	return &InvestorPostgresRepository{getDbFunc: getDbFunc}
}
