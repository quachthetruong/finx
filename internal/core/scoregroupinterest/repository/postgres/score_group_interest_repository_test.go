package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/pkg/dbtest"
	testMock "financing-offer/test/mock"
)

func TestScoreGroupInterestSqlRepository(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewScoreGroupInterestSqlRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	t.Run("GetAvailablePackageBySymbolId success", func(t *testing.T) {
		e := entity.ScoreGroupInterest{
			Id:           1,
			LimitAmount:  decimal.NewFromInt(1000000000),
			LoanRate:     decimal.NewFromFloat(4.0),
			InterestRate: decimal.NewFromFloat(5.0),
			ScoreGroupId: int64(1),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		mock.ExpectQuery("SELECT").WithArgs("ACTIVE", testMock.AnyTime{}, 1, 1).WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"score_group_interest.id",
					"score_group_interest.limit_amount",
					"score_group_interest.loan_rate",
					"score_group_interest.interest_rate",
					"score_group_interest.score_group_id",
					"score_group_interest.created_at",
					"score_group_interest.updated_at",
				},
			).AddRow(e.Id, e.LimitAmount, e.LoanRate, e.InterestRate, e.ScoreGroupId, e.CreatedAt, e.UpdatedAt),
		)
		res, err := repo.GetAvailablePackageBySymbolId(context.Background(), 1)
		assert.Nil(t, err)
		assert.Equal(t, res[0].Id, e.Id)
		assert.Equal(t, res[0].ScoreGroupId, e.ScoreGroupId)
		assert.Equal(t, res[0].InterestRate, e.InterestRate)
		assert.Equal(t, res[0].LoanRate, e.LoanRate)
	})

	t.Run("GetAvailablePackageByGroupId success", func(t *testing.T) {
		e := entity.ScoreGroupInterest{
			Id:           1,
			LimitAmount:  decimal.NewFromInt(1000000000),
			LoanRate:     decimal.NewFromFloat(4.0),
			InterestRate: decimal.NewFromFloat(5.0),
			ScoreGroupId: int64(1),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"score_group_interest.id",
					"score_group_interest.limit_amount",
					"score_group_interest.loan_rate",
					"score_group_interest.interest_rate",
					"score_group_interest.score_group_id",
					"score_group_interest.created_at",
					"score_group_interest.updated_at",
				},
			).AddRow(e.Id, e.LimitAmount, e.LoanRate, e.InterestRate, e.ScoreGroupId, e.CreatedAt, e.UpdatedAt),
		)
		res, err := repo.GetAvailableScoreInterestsByScoreGroupId(context.Background(), 1)
		assert.Nil(t, err)
		assert.Equal(t, res[0].Id, e.Id)
		assert.Equal(t, res[0].ScoreGroupId, e.ScoreGroupId)
		assert.Equal(t, res[0].InterestRate, e.InterestRate)
		assert.Equal(t, res[0].LoanRate, e.LoanRate)
	})

	t.Run("GetAvailableScoreInterestsByScoreGroupId error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(fmt.Errorf("test error"))
		_, err := repo.GetAvailableScoreInterestsByScoreGroupId(context.Background(), 1)
		assert.Equal(t, "ScoreGroupInterestSqlRepository GetAvailableScoreInterestsByScoreGroupId jet: test error", err.Error())
	})

	t.Run("GetAvailablePackageBySymbolId error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").
			WithArgs("ACTIVE", testMock.AnyTime{}, 1, 1).
			WillReturnError(fmt.Errorf("test error"))
		_, err := repo.GetAvailablePackageBySymbolId(context.Background(), 1)
		assert.Equal(t, "ScoreGroupInterestSqlRepository GetAvailablePackageBySymbolId jet: test error", err.Error())
	})
}
