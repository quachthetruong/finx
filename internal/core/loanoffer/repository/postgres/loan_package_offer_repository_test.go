package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/pkg/dbtest"
)

func TestLoanRequestSchedulerConfigRepository(t *testing.T) {
	t.Parallel()
	db, mock, _ := dbtest.New()
	repo := NewLoanPackageOfferPostgresRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)

	t.Run("BulkCreateSuccess", func(t *testing.T) {
		loanOffers := []entity.LoanPackageOffer{
			{
				Id:                   1,
				LoanPackageRequestId: 1,
				OfferedBy:            "admin",
				FlowType:             entity.FLowTypeDnseOffline,
				ExpiredAt:            time.Now(),
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
			{
				Id:                   2,
				LoanPackageRequestId: 1,
				OfferedBy:            "admin",
				FlowType:             entity.FlowTypeDnseOnline,
				ExpiredAt:            time.Now(),
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
		}
		mock.ExpectQuery("INSERT").WillReturnRows(
			mock.NewRows(
				[]string{
					"loan_package_offer.id",
					"loan_package_offer.loan_package_request_id",
					"loan_package_offer.offer_by",
					"loan_package_offer.flow_type",
					"loan_package_offer.expired_at",
					"loan_package_offer.created_at",
					"loan_package_offer.updated_at",
				}).AddRow(loanOffers[0].Id, loanOffers[0].LoanPackageRequestId, loanOffers[0].OfferedBy,
				loanOffers[0].FlowType, loanOffers[0].ExpiredAt, loanOffers[0].CreatedAt, loanOffers[0].UpdatedAt,
			).AddRow(loanOffers[1].Id, loanOffers[1].LoanPackageRequestId, loanOffers[1].OfferedBy,
				loanOffers[1].FlowType, loanOffers[1].ExpiredAt, loanOffers[1].CreatedAt, loanOffers[1].UpdatedAt),
		)
		loanOffers, err := repo.BulkCreate(context.Background(), loanOffers)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(loanOffers))
		assert.Equal(t, int64(1), loanOffers[0].Id)
		assert.Equal(t, int64(2), loanOffers[1].Id)
	})

	t.Run("GetLoanRequestSchedulerConfigFailure", func(t *testing.T) {
		loanOffers := []entity.LoanPackageOffer{
			{
				Id:                   1,
				LoanPackageRequestId: 1,
				OfferedBy:            "admin",
				FlowType:             entity.FLowTypeDnseOffline,
				ExpiredAt:            time.Now(),
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
		}
		mock.ExpectQuery("INSERT").WillReturnError(
			fmt.Errorf("error"),
		)
		_, err := repo.BulkCreate(context.Background(), loanOffers)
		assert.Equal(t, "LoanPackageOfferPostgresRepository BulkCreate jet: error", err.Error())
	})
}
