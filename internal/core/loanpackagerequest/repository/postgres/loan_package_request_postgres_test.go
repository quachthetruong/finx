package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/pkg/dbtest"
)

func TestLoanPackageRequestPostgres(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		return
	}
	repo := NewLoanPackageRequestPostgresRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)

	t.Run("LockByInvestorAndSymbolSuccess", func(t *testing.T) {
		e := entity.LoanPackageRequest{
			Id:          1,
			SymbolId:    1,
			InvestorId:  "test",
			AccountNo:   "accNo",
			LoanRate:    decimal.NewFromFloat(0.3),
			LimitAmount: decimal.NewFromFloat(300000.0),
			Type:        entity.LoanPackageRequestTypeFlexible,
			Status:      entity.LoanPackageRequestStatusPending,
		}
		mock.ExpectQuery("SELECT").WillReturnRows(
			mock.NewRows(
				[]string{
					"loan_package_request.id",
					"loan_package_request.symbol_id",
					"loan_package_request.investor_id",
					"loan_package_request.account_no",
					"loan_package_request.loan_rate",
					"loan_package_request.limitAmount",
					"loan_package_request.type",
					"loan_package_request.status",
					"loan_package_request.updated_at",
					"loan_package_request.update_at",
				}).AddRow(e.Id, e.SymbolId, e.InvestorId, e.AccountNo, e.LoanRate, e.LimitAmount, e.Type, e.Status, e.CreatedAt, e.UpdatedAt),
		)
		requests, err := repo.LockAllPendingRequestByMaxPercent(context.Background(), decimal.NewFromFloat(0.3))
		assert.Nil(t, err)
		assert.Equal(t, int64(1), requests[0].Id)
		assert.Equal(t, decimal.NewFromFloat(0.3), requests[0].LoanRate)
	})

	t.Run("LockByInvestorAndSymbolFailure", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("test"))
		_, err := repo.LockAllPendingRequestByMaxPercent(context.Background(), decimal.NewFromFloat(0.3))
		assert.Equal(t, "LoanPackageRequestPostgresRepository LockAllPendingRequestByMaxPercent: jet: test", err.Error())
	})

	t.Run("UpdateStatusByLoanRequestIdsSuccess", func(t *testing.T) {
		e := entity.LoanPackageRequest{
			Id:          1,
			SymbolId:    1,
			InvestorId:  "test",
			AccountNo:   "accNo",
			LoanRate:    decimal.NewFromFloat(0.3),
			LimitAmount: decimal.NewFromFloat(300000.0),
			Type:        entity.LoanPackageRequestTypeFlexible,
			Status:      entity.LoanPackageRequestStatusPending,
		}
		mock.ExpectQuery("UPDATE").WillReturnRows(
			mock.NewRows(
				[]string{
					"loan_package_request.id",
					"loan_package_request.symbol_id",
					"loan_package_request.investor_id",
					"loan_package_request.account_no",
					"loan_package_request.loan_rate",
					"loan_package_request.limitAmount",
					"loan_package_request.type",
					"loan_package_request.status",
					"loan_package_request.updated_at",
					"loan_package_request.update_at",
				}).AddRow(e.Id, e.SymbolId, e.InvestorId, e.AccountNo, e.LoanRate, e.LimitAmount, e.Type, e.Status, e.CreatedAt, e.UpdatedAt),
		)
		requests, err := repo.UpdateStatusByLoanRequestIds(context.Background(), []int64{1}, entity.LoanPackageRequestStatusConfirmed)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), requests[0].Id)
		assert.Equal(t, decimal.NewFromFloat(0.3), requests[0].LoanRate)
	})

	t.Run("UpdateStatusByLoanRequestIdsFailure", func(t *testing.T) {
		mock.ExpectQuery("UPDATE").WillReturnError(fmt.Errorf("test"))
		_, err := repo.UpdateStatusByLoanRequestIds(context.Background(), []int64{1}, entity.LoanPackageRequestStatusConfirmed)
		assert.Equal(t, "LoanPackageRequestPostgresRepository UpdateStatusByLoanRequestIds: jet: test", err.Error())
	})

	t.Run("LockAndReturnAllPendingRequestBySymbol success", func(t *testing.T) {
		e := entity.LoanPackageRequest{
			Id:          11,
			SymbolId:    123,
			InvestorId:  "test",
			AccountNo:   "accNo",
			LoanRate:    decimal.NewFromFloat(0.3),
			LimitAmount: decimal.NewFromFloat(300000.0),
			Type:        entity.LoanPackageRequestTypeFlexible,
			Status:      entity.LoanPackageRequestStatusPending,
			AssetType:   entity.AssetTypeDerivative,
		}
		mock.ExpectQuery("SELECT").WillReturnRows(
			mock.NewRows(
				[]string{
					"loan_package_request.id",
					"loan_package_request.symbol_id",
					"loan_package_request.investor_id",
					"loan_package_request.account_no",
					"loan_package_request.loan_rate",
					"loan_package_request.limitAmount",
					"loan_package_request.type",
					"loan_package_request.status",
					"loan_package_request.updated_at",
					"loan_package_request.update_at",
					"loan_package_request.asset_type",
				}).AddRow(e.Id, e.SymbolId, e.InvestorId, e.AccountNo, e.LoanRate, e.LimitAmount, e.Type, e.Status, e.CreatedAt, e.UpdatedAt, e.AssetType),
		)
		requests, err := repo.LockAndReturnAllPendingRequestBySymbolId(context.Background(), e.SymbolId)
		assert.Nil(t, err)
		assert.Equal(t, int64(11), requests[0].Id)
		assert.Equal(t, decimal.NewFromFloat(0.3), requests[0].LoanRate)
		assert.Equal(t, entity.AssetTypeDerivative, requests[0].AssetType)
	})
	t.Run("LockAndReturnAllPendingRequestBySymbol failure", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(assert.AnError)
		_, err := repo.LockAndReturnAllPendingRequestBySymbolId(context.Background(), 123)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
