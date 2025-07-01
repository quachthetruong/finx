package postgres

import (
	"context"
	"fmt"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/scheduler/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

var _ repository.SchedulerJobRepository = (*SchedulerJobRepository)(nil)

type SchedulerJobRepository struct {
	getDbFunc database.GetDbFunc
}

func (s *SchedulerJobRepository) Create(ctx context.Context, job entity.SchedulerJob) error {
	createModel := MapSchedulerJobEntityToDb(job)
	_, err := table.SchedulerJob.INSERT(table.SchedulerJob.MutableColumns).
		MODEL(createModel).
		RETURNING(table.SchedulerJob.AllColumns).
		ExecContext(ctx, s.getDbFunc(ctx))
	if err != nil {
		return fmt.Errorf("SchedulerJobRepository Create: %w", err)
	}
	return nil
}

func NewSchedulerJobRepository(getDbFunction database.GetDbFunc) *SchedulerJobRepository {
	return &SchedulerJobRepository{getDbFunc: getDbFunction}
}
