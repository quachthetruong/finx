package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/go-jet/jet/v2/postgres"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/scheduler/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

var _ repository.LoanRequestSchedulerConfigRepository = (*LoanRequestSchedulerConfigRepository)(nil)

type LoanRequestSchedulerConfigRepository struct {
	getDbFunc database.GetDbFunc
}

func (r *LoanRequestSchedulerConfigRepository) GetAll(ctx context.Context) ([]entity.LoanRequestSchedulerConfig, error) {
	des := make([]model.LoanRequestSchedulerConfig, 0)
	err := table.LoanRequestSchedulerConfig.SELECT(table.LoanRequestSchedulerConfig.AllColumns).
		QueryContext(ctx, r.getDbFunc(ctx), &des)
	if err != nil {
		return []entity.LoanRequestSchedulerConfig{}, fmt.Errorf("LoanRequestSchedulerConfigRepository GetAll: %w", err)
	}
	return MapLoanRequestSchedulerConfigsDbToEntity(des), nil
}

func (r *LoanRequestSchedulerConfigRepository) Create(ctx context.Context, config entity.LoanRequestSchedulerConfig) (entity.LoanRequestSchedulerConfig, error) {
	createdModel := MapLoanRequestSchedulerConfigEntityToDb(config)
	created := model.LoanRequestSchedulerConfig{}
	err := table.LoanRequestSchedulerConfig.INSERT(table.LoanRequestSchedulerConfig.MutableColumns).
		MODEL(createdModel).
		RETURNING(table.LoanRequestSchedulerConfig.AllColumns).
		QueryContext(ctx, r.getDbFunc(ctx), &created)
	if err != nil {
		return entity.LoanRequestSchedulerConfig{}, fmt.Errorf("LoanRequestSchedulerConfigRepository Create: %w", err)
	}
	return MapLoanRequestSchedulerConfigDbToEntity(created), nil
}

func (r *LoanRequestSchedulerConfigRepository) GetCurrentConfig(ctx context.Context) (entity.LoanRequestSchedulerConfig, error) {
	var des model.LoanRequestSchedulerConfig
	err := table.LoanRequestSchedulerConfig.SELECT(table.LoanRequestSchedulerConfig.AllColumns).
		WHERE(table.LoanRequestSchedulerConfig.AffectedFrom.LT_EQ(postgres.TimestampT(time.Now()))).
		ORDER_BY(table.LoanRequestSchedulerConfig.AffectedFrom.DESC()).
		LIMIT(1).
		QueryContext(ctx, r.getDbFunc(ctx), &des)
	if err != nil {
		return entity.LoanRequestSchedulerConfig{}, fmt.Errorf("LoanRequestSchedulerConfigRepository GetCurrentConfig: %w", err)
	}
	return MapLoanRequestSchedulerConfigDbToEntity(des), nil
}

func NewLoanRequestSchedulerConfigRepo(getDbFunction database.GetDbFunc) *LoanRequestSchedulerConfigRepository {
	return &LoanRequestSchedulerConfigRepository{getDbFunc: getDbFunction}
}
