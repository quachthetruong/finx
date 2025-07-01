package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/pkg/dbtest"
)

func TestSchedulerJobRepository(t *testing.T) {
	t.Parallel()
	db, mock, _ := dbtest.New()
	repo := NewSchedulerJobRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)

	t.Run("TestCreateSchedulerJobSuccess", func(t *testing.T) {
		e := entity.SchedulerJob{
			Id:           1,
			JobType:      entity.JobTypeDeclineHighRiskLoanRequest,
			JobStatus:    entity.JobStatusSuccess,
			TrackingData: "{}",
			TriggerBy:    "system",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		mock.ExpectExec("INSERT").
			WithArgs(e.JobType, e.JobStatus, e.TriggerBy, e.TrackingData).
			WillReturnResult(
				sqlmock.NewResult(1, 1))
		err := repo.Create(context.Background(), e)
		assert.Nil(t, err)
	})
	t.Run("GetLoanRequestSchedulerConfigFailure", func(t *testing.T) {
		e := entity.SchedulerJob{
			Id:           1,
			JobType:      entity.JobTypeDeclineHighRiskLoanRequest,
			JobStatus:    entity.JobStatusSuccess,
			TrackingData: "{}",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		mock.ExpectExec("INSERT").WillReturnError(
			fmt.Errorf("error"),
		)
		err := repo.Create(context.Background(), e)
		assert.Equal(t, "SchedulerJobRepository Create: error", err.Error())
	})
}
