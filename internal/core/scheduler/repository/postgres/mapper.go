package postgres

import (
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

func MapLoanRequestSchedulerConfigDbToEntity(config model.LoanRequestSchedulerConfig) entity.LoanRequestSchedulerConfig {
	return entity.LoanRequestSchedulerConfig{
		ID:              config.ID,
		MaximumLoanRate: config.MaximumLoanRate,
		AffectedFrom:    config.AffectedFrom,
		CreatedAt:       config.CreatedAt,
		UpdatedAt:       config.UpdatedAt,
	}
}

func MapLoanRequestSchedulerConfigEntityToDb(config entity.LoanRequestSchedulerConfig) model.LoanRequestSchedulerConfig {
	return model.LoanRequestSchedulerConfig{
		ID:              config.ID,
		MaximumLoanRate: config.MaximumLoanRate,
		AffectedFrom:    config.AffectedFrom,
		CreatedAt:       config.CreatedAt,
		UpdatedAt:       config.UpdatedAt,
	}
}

func MapLoanRequestSchedulerConfigsDbToEntity(configs []model.LoanRequestSchedulerConfig) []entity.LoanRequestSchedulerConfig {
	res := make([]entity.LoanRequestSchedulerConfig, 0, len(configs))
	for _, config := range configs {
		res = append(res, MapLoanRequestSchedulerConfigDbToEntity(config))
	}
	return res
}

func MapSchedulerJobEntityToDb(e entity.SchedulerJob) model.SchedulerJob {
	return model.SchedulerJob{
		ID:           e.Id,
		JobType:      string(e.JobType),
		JobStatus:    string(e.JobStatus),
		TriggerBy:    e.TriggerBy,
		TrackingData: e.TrackingData,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func MapSchedulerJobDbToEntity(m model.SchedulerJob) entity.SchedulerJob {
	return entity.SchedulerJob{
		Id:           m.ID,
		JobStatus:    entity.JobStatus(m.JobStatus),
		JobType:      entity.JobType(m.JobType),
		TriggerBy:    m.TriggerBy,
		TrackingData: m.TrackingData,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
