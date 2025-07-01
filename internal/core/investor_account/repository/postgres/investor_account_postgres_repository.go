package postgres

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

type InvestorAccountPostgresRepository struct {
	getDbFunc database.GetDbFunc
}

func NewInvestorAccountPostgresRepository(getDbFunc database.GetDbFunc) *InvestorAccountPostgresRepository {
	return &InvestorAccountPostgresRepository{getDbFunc: getDbFunc}
}

func (r *InvestorAccountPostgresRepository) GetByAccountNo(ctx context.Context, accountNo string) (entity.InvestorAccount, error) {
	errorTemplate := "GetByAccountNo %w"
	account := model.InvestorAccount{}
	if err := table.InvestorAccount.
		SELECT(table.InvestorAccount.AllColumns).
		WHERE(table.InvestorAccount.AccountNo.EQ(postgres.String(accountNo))).
		QueryContext(ctx, r.getDbFunc(ctx), &account); err != nil {
		return entity.InvestorAccount{}, fmt.Errorf(errorTemplate, err)
	}
	return MapInvestorAccountDbToEntity(account), nil
}

func (r *InvestorAccountPostgresRepository) Update(ctx context.Context, account entity.InvestorAccount) (entity.InvestorAccount, error) {
	errorTemplate := "Update %w"
	updateModel := MapInvestorAccountEntityToDb(account)
	updated := model.InvestorAccount{}
	if err := table.InvestorAccount.
		UPDATE(table.InvestorAccount.MutableColumns).
		MODEL(updateModel).
		WHERE(table.InvestorAccount.AccountNo.EQ(postgres.String(account.AccountNo))).
		RETURNING(table.InvestorAccount.AllColumns).
		QueryContext(ctx, r.getDbFunc(ctx), &updated); err != nil {
		return entity.InvestorAccount{}, fmt.Errorf(errorTemplate, err)
	}
	return MapInvestorAccountDbToEntity(updated), nil
}

func (r *InvestorAccountPostgresRepository) Create(ctx context.Context, account entity.InvestorAccount) (entity.InvestorAccount, error) {
	errorTemplate := "Create %w"
	insertModel := MapInvestorAccountEntityToDb(account)
	created := model.InvestorAccount{}
	if err := table.InvestorAccount.
		INSERT(table.InvestorAccount.AllColumns.Except(table.InvestorAccount.CreatedAt, table.InvestorAccount.UpdatedAt)).
		MODEL(insertModel).
		RETURNING(table.InvestorAccount.AllColumns).
		QueryContext(ctx, r.getDbFunc(ctx), &created); err != nil {
		return entity.InvestorAccount{}, fmt.Errorf(errorTemplate, err)
	}
	return MapInvestorAccountDbToEntity(created), nil
}
