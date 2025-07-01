package submission_default

import (
	"context"
	"financing-offer/internal/config/repository"
	"financing-offer/internal/core/entity"
	"fmt"
)

type UseCase interface {
	SetSubmissionDefault(ctx context.Context, defaultValues entity.SubmissionDefault, updater string) (entity.SubmissionDefault, error)
	GetSubmissionDefault(ctx context.Context) (entity.SubmissionDefault, error)
}

type useCase struct {
	configurationPersistenceRepo repository.ConfigurationPersistenceRepository
}

func NewUseCase(
	configurationPersistenceRepo repository.ConfigurationPersistenceRepository,
) UseCase {
	return &useCase{
		configurationPersistenceRepo: configurationPersistenceRepo,
	}
}

func (u *useCase) GetSubmissionDefault(ctx context.Context) (entity.SubmissionDefault, error) {
	defaultValues, err := u.configurationPersistenceRepo.GetSubmissionDefault(ctx)
	if err != nil {
		return entity.SubmissionDefault{}, fmt.Errorf(
			"submissionSheetDefaultValueUseCase GetSubmissionDefault %w", err,
		)
	}
	return defaultValues, nil
}

func (u *useCase) SetSubmissionDefault(ctx context.Context, defaultValues entity.SubmissionDefault, updater string) (entity.SubmissionDefault, error) {
	err := u.configurationPersistenceRepo.SetSubmissionDefault(ctx, defaultValues, updater)
	if err != nil {
		return defaultValues, fmt.Errorf("submissionSheetDefaultValueUseCase SetSubmissionDefault %w", err)
	}
	return defaultValues, nil
}
