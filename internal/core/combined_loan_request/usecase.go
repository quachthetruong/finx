package combinedloanrequest

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"financing-offer/internal/core"
	"financing-offer/internal/core/combined_loan_request/repository"
	"financing-offer/internal/core/entity"
)

type UseCase interface {
	GetAll(ctx context.Context, filter entity.CombinedLoanRequestFilter) ([]entity.CombinedLoanRequest, core.PagingMetaData, error)
}

type useCase struct {
	repository repository.CombinedLoanPackageRequestPersistenceRepository
}

func (u *useCase) GetAll(ctx context.Context, filter entity.CombinedLoanRequestFilter) ([]entity.CombinedLoanRequest, core.PagingMetaData, error) {
	var (
		pagingMetaData = core.PagingMetaData{PageSize: filter.Size, PageNumber: filter.Number}
		eg             errgroup.Group
		requests       []entity.CombinedLoanRequest
	)
	eg.Go(
		func() error {
			entities, err := u.repository.GetAll(ctx, filter)
			if err != nil {
				return err
			}
			requests = entities
			return nil
		},
	)
	eg.Go(
		func() error {
			count, err := u.repository.Count(ctx, filter)
			if err != nil {
				return err
			}
			pagingMetaData.Total = count
			pagingMetaData.TotalPages = filter.TotalPages(count)
			return nil
		},
	)
	if err := eg.Wait(); err != nil {
		return nil, core.PagingMetaData{}, fmt.Errorf("combinedLoanRequestUseCase GetAll: %w", err)
	}
	return requests, pagingMetaData, nil
}

func NewUseCase(repository repository.CombinedLoanPackageRequestPersistenceRepository) UseCase {
	return &useCase{repository: repository}
}
