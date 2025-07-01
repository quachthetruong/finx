package marginoperation

import (
	"context"
	"financing-offer/internal/apperrors"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/marginoperation/repository"
	"financing-offer/internal/funcs"
	"fmt"
)

type UseCase interface {
	VerifyMarginPoolId(ctx context.Context, id int64) error
	GetPoolGroupMapByPoolIds(ctx context.Context, poolIds []int64) (map[int64]entity.MarginPoolGroup, error)
}

type marginPoolUseCase struct {
	marginPoolRepository repository.MarginOperationRepository
}

func NewUseCase(marginPoolRepository repository.MarginOperationRepository) UseCase {
	return &marginPoolUseCase{
		marginPoolRepository: marginPoolRepository,
	}
}

func (u *marginPoolUseCase) VerifyMarginPoolId(ctx context.Context, id int64) error {
	errTemplate := "marginPoolUseCase VerifyMarginPoolId %w"
	_, err := u.marginPoolRepository.GetMarginPoolById(ctx, id)
	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	return nil
}

func (u *marginPoolUseCase) GetPoolGroupMapByPoolIds(ctx context.Context, poolIds []int64) (map[int64]entity.MarginPoolGroup, error) {
	errTemplate := "marginPoolUseCase GetPoolGroupMapByPoolId %w"
	marginPools, err := u.marginPoolRepository.GetMarginPoolsByIds(ctx, poolIds)
	if err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}
	poolGroupIds := funcs.Map(marginPools, func(m entity.MarginPool) int64 { return m.PoolGroupId })
	marginPoolGroups, err := u.marginPoolRepository.GetMarginPoolGroupsByIds(ctx, poolGroupIds)
	if err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}
	marginPoolGroupMapped := funcs.AssociateBy[entity.MarginPoolGroup, int64](
		marginPoolGroups, func(group entity.MarginPoolGroup) int64 {
			return group.Id
		},
	)
	marginPoolMapped := make(map[int64]entity.MarginPoolGroup)
	for _, pool := range marginPools {
		group, ok := marginPoolGroupMapped[pool.PoolGroupId]
		if !ok {
			return nil, fmt.Errorf(errTemplate, apperrors.ErrorInvalidPoolId)
		}
		marginPoolMapped[pool.Id] = group
	}
	return marginPoolMapped, nil
}
