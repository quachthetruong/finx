package postgres

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/loancontract/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/pkg/querymod"
)

var _ repository.LoanContractPersistenceRepository = (*LoanContractRepository)(nil)

type LoanContractRepository struct {
	getDbFunc database.GetDbFunc
}

func (r *LoanContractRepository) GetInvestorActiveContract(ctx context.Context, investorId string, symbolId int64) (entity.LoanContract, error) {
	res := model.LoanContract{}
	stm := postgres.SELECT(table.LoanContract.AllColumns).
		FROM(
			table.LoanContract.
				INNER_JOIN(
					table.LoanPackageOfferInterest,
					table.LoanContract.LoanOfferInterestID.EQ(table.LoanPackageOfferInterest.ID),
				),
		).
		WHERE(
			table.LoanContract.InvestorID.EQ(postgres.String(investorId)).
				AND(
					table.LoanContract.SymbolID.EQ(postgres.Int64(symbolId)).
						AND(
							table.LoanPackageOfferInterest.Status.EQ(postgres.String(entity.LoanPackageOfferInterestStatusLoanPackageCreated.String())).
								OR(table.LoanPackageOfferInterest.Status.EQ(postgres.String(entity.LoanPackageOfferInterestStatusSigned.String()))),
						),
				),
		)
	if err := stm.QueryContext(ctx, r.getDbFunc(ctx), &res); err != nil {
		return entity.LoanContract{}, fmt.Errorf("LoanContractRepository GetInvestorActiveContract %w", err)
	}
	return MapLoanContractDbToEntity(res), nil
}

func (r *LoanContractRepository) GetById(ctx context.Context, id int64, opts ...querymod.GetOption) (entity.LoanContract, error) {
	getQm := querymod.GetQm{}
	for _, opt := range opts {
		opt(&getQm)
	}
	loanContract := model.LoanContract{}
	stm := table.LoanContract.
		SELECT(table.LoanContract.AllColumns).
		WHERE(table.LoanContract.ID.EQ(postgres.Int64(id)))
	if getQm.ForUpdate {
		stm = stm.FOR(postgres.UPDATE())
	}
	if err := stm.QueryContext(ctx, r.getDbFunc(ctx), &loanContract); err != nil {
		return entity.LoanContract{}, fmt.Errorf("LoanContractRepository GetById %w", err)
	}
	return MapLoanContractDbToEntity(loanContract), nil
}

func (r *LoanContractRepository) Create(ctx context.Context, loanContract entity.LoanContract) (entity.LoanContract, error) {
	toCreate := MapLoanContractEntityToDb(loanContract)
	created := model.LoanContract{}
	if err := table.LoanContract.
		INSERT(table.LoanContract.MutableColumns).
		MODEL(toCreate).
		RETURNING(table.LoanContract.AllColumns).
		QueryContext(ctx, r.getDbFunc(ctx), &created); err != nil {
		return entity.LoanContract{}, fmt.Errorf("LoanContractRepository Create %w", err)
	}
	return MapLoanContractDbToEntity(created), nil
}

func (r *LoanContractRepository) BulkCreate(ctx context.Context, loanContracts []entity.LoanContract) error {
	toCreateModels := MapLoanContractsEntityToDb(loanContracts)
	if _, err := table.LoanContract.
		INSERT(table.LoanContract.MutableColumns).
		MODELS(toCreateModels).
		ExecContext(ctx, r.getDbFunc(ctx)); err != nil {
		return fmt.Errorf("LoanContractRepository Create %w", err)
	}
	return nil
}

func NewLoanContractRepository(getDbFunc database.GetDbFunc) *LoanContractRepository {
	return &LoanContractRepository{getDbFunc: getDbFunc}
}
