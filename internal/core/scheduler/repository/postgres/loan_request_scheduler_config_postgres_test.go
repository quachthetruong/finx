package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/pkg/dbtest"
)

func TestLoanRequestSchedulerConfigRepository(t *testing.T) {
	t.Parallel()
	db, mock, _ := dbtest.New()
	repo := NewLoanRequestSchedulerConfigRepo(
		func(ctx context.Context) database.DB {
			return db
		},
	)

	t.Run("GetAllLoanRequestSchedulerConfigSuccess", func(t *testing.T) {
		affectedFrom := time.Now()
		e := entity.LoanRequestSchedulerConfig{
			ID:              1,
			MaximumLoanRate: decimal.NewFromFloat(0.2),
			AffectedFrom:    affectedFrom,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		mock.ExpectQuery("SELECT").WillReturnRows(
			mock.NewRows(
				[]string{
					"loan_request_scheduler_config.id",
					"loan_request_scheduler_config.maximum_loan_rate",
					"loan_request_scheduler_config.affected_from",
					"loan_request_scheduler_config.created_at",
					"loan_request_scheduler_config.updated_at",
				}).AddRow(e.ID, e.MaximumLoanRate, e.AffectedFrom, e.CreatedAt, e.UpdatedAt),
		)
		configs, err := repo.GetAll(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, int64(1), configs[0].ID)
		assert.Equal(t, decimal.NewFromFloat(0.2), configs[0].MaximumLoanRate)
		assert.Equal(t, affectedFrom, configs[0].AffectedFrom)
	})
	t.Run("GetLoanRequestSchedulerConfigFailure", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(
			fmt.Errorf("error"),
		)
		_, err := repo.GetAll(context.Background())
		assert.Equal(t, "LoanRequestSchedulerConfigRepository GetAll: jet: error", err.Error())
	})
	t.Run("GetCurrentLoanRequestSchedulerConfigSuccess", func(t *testing.T) {
		affectedFrom := time.Now()
		e := entity.LoanRequestSchedulerConfig{
			ID:              1,
			MaximumLoanRate: decimal.NewFromFloat(0.2),
			AffectedFrom:    affectedFrom,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		mock.ExpectQuery("SELECT").WillReturnRows(
			mock.NewRows(
				[]string{
					"loan_request_scheduler_config.id",
					"loan_request_scheduler_config.maximum_loan_rate",
					"loan_request_scheduler_config.affected_from",
					"loan_request_scheduler_config.created_at",
					"loan_request_scheduler_config.updated_at",
				}).AddRow(e.ID, e.MaximumLoanRate, e.AffectedFrom, e.CreatedAt, e.UpdatedAt),
		)
		config, err := repo.GetCurrentConfig(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, int64(1), config.ID)
		assert.Equal(t, decimal.NewFromFloat(0.2), config.MaximumLoanRate)
		assert.Equal(t, affectedFrom, config.AffectedFrom)
	})

	t.Run("GetCurrentLoanRequestSchedulerConfigFailure", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(
			fmt.Errorf("error"),
		)
		_, err := repo.GetCurrentConfig(context.Background())
		assert.Equal(t, "LoanRequestSchedulerConfigRepository GetCurrentConfig: jet: error", err.Error())
	})

	t.Run("CreateLoanRequestSchedulerConfigFailure", func(t *testing.T) {
		mock.ExpectQuery("INSERT").WillReturnError(
			fmt.Errorf("error"),
		)
		_, err := repo.Create(context.Background(), entity.LoanRequestSchedulerConfig{
			ID:              1,
			MaximumLoanRate: decimal.NewFromFloat(0.2),
			AffectedFrom:    time.Now(),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		})
		assert.Equal(t, "LoanRequestSchedulerConfigRepository Create: jet: error", err.Error())
	})

	t.Run("CreateLoanRequestSchedulerConfigSuccess", func(t *testing.T) {
		affectedFrom := time.Now()
		e := entity.LoanRequestSchedulerConfig{
			ID:              1,
			MaximumLoanRate: decimal.NewFromFloat(0.2),
			AffectedFrom:    affectedFrom,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		mock.ExpectQuery("INSERT").WillReturnRows(
			mock.NewRows(
				[]string{
					"loan_request_scheduler_config.id",
					"loan_request_scheduler_config.maximum_loan_rate",
					"loan_request_scheduler_config.affected_from",
					"loan_request_scheduler_config.created_at",
					"loan_request_scheduler_config.updated_at",
				}).AddRow(e.ID, e.MaximumLoanRate, e.AffectedFrom, e.CreatedAt, e.UpdatedAt),
		)
		config, err := repo.Create(context.Background(), entity.LoanRequestSchedulerConfig{
			ID:              1,
			MaximumLoanRate: decimal.NewFromFloat(0.2),
			AffectedFrom:    time.Now(),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), config.ID)
		assert.Equal(t, decimal.NewFromFloat(0.2), config.MaximumLoanRate)
		assert.Equal(t, affectedFrom, config.AffectedFrom)
	})
}
