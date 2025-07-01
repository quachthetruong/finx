package awaitingconfirmrequest

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"financing-offer/internal/core"
	"financing-offer/internal/core/awaiting_confirm_request/repository"
	"financing-offer/internal/core/entity"
)

type UseCase interface {
	GetAll(ctx context.Context, filter entity.AwaitingConfirmRequestFilter) ([]entity.AwaitingConfirmRequest, core.PagingMetaData, error)
}

var _ UseCase = (*useCase)(nil)

type useCase struct {
	repository repository.AwaitingConfirmRequestPersistenceRepository
}

func (u *useCase) GetAll(ctx context.Context, filter entity.AwaitingConfirmRequestFilter) ([]entity.AwaitingConfirmRequest, core.PagingMetaData, error) {
	var (
		pagingMetaData = core.PagingMetaData{PageSize: filter.Size, PageNumber: filter.Number}
		eg             errgroup.Group
		requests       []entity.AwaitingConfirmRequest
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
		return nil, core.PagingMetaData{}, fmt.Errorf("awaitingConfirmRequestUseCase GetAll: %w", err)
	}
	return requests, pagingMetaData, nil
}

func NewUseCase(repository repository.AwaitingConfirmRequestPersistenceRepository) UseCase {
	return &useCase{repository: repository}
}
