package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/pkg/dbtest"
)

func TestInvestorAccountPostgresRepository_GetByAccountNo(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Error(err)
	}
	repo := NewInvestorAccountPostgresRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	account := entity.InvestorAccount{
		AccountNo:    "1",
		InvestorId:   "000121",
		MarginStatus: entity.MarginStatusVersion2,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	// Test cases
	t.Run("Get investor_account by investor_account no success", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(account.AccountNo).WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"investor_account.account_no",
					"investor_account.investor_id",
					"investor_account.margin_status",
					"investor_account.created_at",
					"investor_account.updated_at",
				}).AddRow(account.AccountNo, account.InvestorId, account.MarginStatus, account.CreatedAt, account.UpdatedAt),
		)
		res, err := repo.GetByAccountNo(context.Background(), account.AccountNo)
		assert.Nil(t, err)
		assert.Equal(t, account, res)
	})
	t.Run("Get investor_account version by investor_account no error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs("4").WillReturnError(assert.AnError)
		_, err := repo.GetByAccountNo(context.Background(), "4")
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestInvestorAccountPostgresRepository_Update(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Error(err)
	}
	repo := NewInvestorAccountPostgresRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	account := entity.InvestorAccount{
		AccountNo:    "1",
		MarginStatus: entity.MarginStatusVersion2,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Test cases
	t.Run("Update investor_account version success", func(t *testing.T) {
		mock.ExpectQuery("UPDATE").WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"investor_account.account_no",
					"investor_account.margin_status",
					"investor_account.created_at",
					"investor_account.updated_at",
				}).AddRow(account.AccountNo, account.MarginStatus, account.CreatedAt, account.UpdatedAt),
		)
		res, err := repo.Update(context.Background(), account)
		assert.Nil(t, err)
		assert.Equal(t, account, res)
	})
	t.Run("Update investor_account version error", func(t *testing.T) {
		mock.ExpectQuery("UPDATE").WillReturnError(assert.AnError)
		_, err := repo.Update(context.Background(), account)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestInvestorAccountPostgresRepository_Create(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Error(err)
	}
	repo := NewInvestorAccountPostgresRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	account := entity.InvestorAccount{
		AccountNo:    "1",
		MarginStatus: entity.MarginStatusVersion2,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	// Test cases
	t.Run("Create investor_account version success", func(t *testing.T) {
		mock.ExpectQuery("INSERT").WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"investor_account.account_no",
					"investor_account.margin_status",
					"investor_account.created_at",
					"investor_account.updated_at",
				}).AddRow(account.AccountNo, account.MarginStatus, account.CreatedAt, account.UpdatedAt),
		)
		res, err := repo.Create(context.Background(), account)
		assert.Nil(t, err)
		assert.Equal(t, account, res)
	})
	t.Run("Create investor_account version error", func(t *testing.T) {
		mock.ExpectQuery("INSERT").WillReturnError(assert.AnError)
		_, err := repo.Create(context.Background(), account)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
