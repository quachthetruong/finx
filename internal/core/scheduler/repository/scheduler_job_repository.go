package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type SchedulerJobRepository interface {
	Create(ctx context.Context, job entity.SchedulerJob) error
}
