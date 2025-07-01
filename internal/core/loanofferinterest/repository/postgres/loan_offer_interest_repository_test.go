package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/pkg/dbtest"
	testMock "financing-offer/test/mock"
)

func TestLoanPackageOfferInterestPostgresRepository_CancelExpiredOfferInterests(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		return
	}
	repo := NewLoanPackageOfferInterestPostgresRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	t.Run("CancelExpiredOfferInterestsSuccess", func(t *testing.T) {
		mock.ExpectExec("UPDATE").WithArgs(entity.LoanPackageOfferInterestStatusCancelled, "user", testMock.AnyTime{}, entity.LoanPackageOfferCancelledReasonInvestor, 1, entity.LoanPackageRequestStatusPending).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.CancelByOfferId(context.Background(), 1, "user", entity.LoanPackageOfferCancelledReasonInvestor)
		assert.Nil(t, err)
	})

	t.Run("CancelExpiredOfferInterestsFail", func(t *testing.T) {
		mock.ExpectExec("UPDATE").WithArgs(entity.LoanPackageOfferInterestStatusCancelled, "user", testMock.AnyTime{}, entity.LoanPackageOfferCancelledReasonInvestor, 1, entity.LoanPackageRequestStatusPending).
			WillReturnError(errors.New("test"))
		err := repo.CancelByOfferId(context.Background(), 1, "user", entity.LoanPackageOfferCancelledReasonInvestor)
		assert.Equal(t, "LoanPackageOfferInterestPostgresRepository CancelByOfferId test", err.Error())
	})
}
