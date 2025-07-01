package scoregroup

import (
	"context"
	"fmt"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/scoregroup/repository"
	scoreGroupInterestRepository "financing-offer/internal/core/scoregroupinterest/repository"
)

type UseCase interface {
	GetAll(ctx context.Context) ([]entity.ScoreGroup, error)
	Create(ctx context.Context, scoreGroup entity.ScoreGroup) (entity.ScoreGroup, error)
	Update(ctx context.Context, scoreGroup entity.ScoreGroup) (entity.ScoreGroup, error)
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (entity.ScoreGroup, error)
	GetAvailablePackagesById(ctx context.Context, scoreGroupId int64) ([]entity.ScoreGroupInterest, error)
}

func NewUseCase(
	scoreGroupRepo repository.ScoreGroupRepository,
	scoreGroupInterestRepository scoreGroupInterestRepository.ScoreGroupInterestRepository,
) UseCase {
	return &scoreGroupUseCase{scoreGroupRepository: scoreGroupRepo, scoreGroupInterestRepository: scoreGroupInterestRepository}
}

type scoreGroupUseCase struct {
	scoreGroupRepository         repository.ScoreGroupRepository
	scoreGroupInterestRepository scoreGroupInterestRepository.ScoreGroupInterestRepository
}

func (u *scoreGroupUseCase) GetById(ctx context.Context, id int64) (entity.ScoreGroup, error) {
	res, err := u.scoreGroupRepository.GetById(ctx, id)
	if err != nil {
		return res, fmt.Errorf("scoreGroupUseCase GetById %w", err)
	}
	return res, nil
}

func (u *scoreGroupUseCase) GetAll(ctx context.Context) ([]entity.ScoreGroup, error) {
	res, err := u.scoreGroupRepository.GetAll(ctx)
	if err != nil {
		return res, fmt.Errorf("scoreGroupUseCase GetAll %w", err)
	}
	return res, nil
}

func (u *scoreGroupUseCase) Create(ctx context.Context, scoreGroup entity.ScoreGroup) (entity.ScoreGroup, error) {
	res, err := u.scoreGroupRepository.Create(ctx, scoreGroup)
	if err != nil {
		return res, fmt.Errorf("scoreGroupUseCase Create %w", err)
	}
	return res, nil
}

func (u *scoreGroupUseCase) Update(ctx context.Context, scoreGroup entity.ScoreGroup) (entity.ScoreGroup, error) {
	res, err := u.scoreGroupRepository.Update(ctx, scoreGroup)
	if err != nil {
		return res, fmt.Errorf("scoreGroupUseCase Update %w", err)
	}
	return res, nil
}

func (u *scoreGroupUseCase) Delete(ctx context.Context, id int64) error {
	err := u.scoreGroupRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("scoreGroupUseCase Delete %w", err)
	}
	return nil
}

func (u *scoreGroupUseCase) GetAvailablePackagesById(ctx context.Context, id int64) ([]entity.ScoreGroupInterest, error) {
	availablePackages, err := u.scoreGroupInterestRepository.GetAvailableScoreInterestsByScoreGroupId(ctx, id)
	if err != nil {
		return availablePackages, fmt.Errorf("scoreGroupUseCase GetAvailablePackagesById %w", err)
	}
	return availablePackages, nil
}
