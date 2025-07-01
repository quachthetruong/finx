package testhelper

import (
	"testing"

	"github.com/go-jet/jet/v2/postgres"

	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

func GetLoanOfferInterest(t *testing.T, db database.DB, id int64) model.LoanPackageOfferInterest {
	res := model.LoanPackageOfferInterest{}
	if err := table.LoanPackageOfferInterest.SELECT(table.LoanPackageOfferInterest.AllColumns).WHERE(table.LoanPackageOfferInterest.ID.EQ(postgres.Int64(id))).Query(
		db, &res,
	); err != nil {
		t.Error(err)
	}
	return res
}

func GetLoanContract(t *testing.T, db database.DB, id int64) model.LoanContract {
	res := model.LoanContract{}
	if err := table.LoanContract.SELECT(table.LoanContract.AllColumns).WHERE(table.LoanContract.ID.EQ(postgres.Int64(id))).Query(
		db, &res,
	); err != nil {
		t.Error(err)
	}
	return res
}
