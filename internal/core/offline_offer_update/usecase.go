package offlineofferupdate

import (
	"context"
	"fmt"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/atomicity"
	"financing-offer/internal/core/entity"
	loanPackageOfferRepo "financing-offer/internal/core/loanoffer/repository"
	loanPackageOfferInterestRepo "financing-offer/internal/core/loanofferinterest/repository"
	"financing-offer/internal/core/offline_offer_update/repository"
)

type UseCase interface {
	Create(ctx context.Context, offlineOfferUpdate entity.OfflineOfferUpdate, assetType entity.AssetType) (entity.OfflineOfferUpdate, error)
	GetByOfferId(ctx context.Context, offerId int64, assetType entity.AssetType) ([]entity.OfflineOfferUpdate, error)
}

type useCase struct {
	repository              repository.OfflineOfferUpdatePersistenceRepository
	offerRepository         loanPackageOfferRepo.LoanPackageOfferRepository
	offerInterestRepository loanPackageOfferInterestRepo.LoanPackageOfferInterestRepository
	atomicExecutor          atomicity.AtomicExecutor
}

func (u *useCase) Create(ctx context.Context, offlineOfferUpdate entity.OfflineOfferUpdate, assetType entity.AssetType) (entity.OfflineOfferUpdate, error) {
	var (
		res   entity.OfflineOfferUpdate
		offer entity.LoanPackageOffer
		err   error
	)
	offer, err = u.offerRepository.InvestorGetById(ctx, offlineOfferUpdate.OfferId)
	if err != nil {
		return res, err
	}
	if offer.LoanPackageRequest != nil && offer.LoanPackageRequest.AssetType != assetType {
		return res, apperrors.AssetTypeDoesNotMatch
	}
	txErr := u.atomicExecutor.Execute(
		ctx, func(ctx context.Context) error {
			res, err = u.repository.Create(ctx, offlineOfferUpdate)
			if err != nil {
				return err
			}
			if offlineOfferUpdate.Status == entity.OfflineOfferUpdateStatusRejected {
				err = u.offerInterestRepository.CancelByOfferId(
					ctx, offlineOfferUpdate.OfferId, offlineOfferUpdate.CreatedBy,
					entity.LoanPackageOfferCancelledReasonInvestor,
				)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
	if txErr != nil {
		return res, fmt.Errorf("offline offer useCase Create %w", txErr)
	}
	return res, nil
}

func (u *useCase) GetByOfferId(ctx context.Context, offerId int64, assetType entity.AssetType) ([]entity.OfflineOfferUpdate, error) {
	res, err := u.repository.GetByOfferId(ctx, offerId, assetType)
	if err != nil {
		return nil, fmt.Errorf("offline offer useCase GetByOfferId %w", err)
	}
	return res, nil
}

func NewUseCase(
	repository repository.OfflineOfferUpdatePersistenceRepository,
	offerRepository loanPackageOfferRepo.LoanPackageOfferRepository,
	offerInterestRepository loanPackageOfferInterestRepo.LoanPackageOfferInterestRepository,
	atomicExecutor atomicity.AtomicExecutor,
) UseCase {
	return &useCase{
		repository:              repository,
		offerRepository:         offerRepository,
		offerInterestRepository: offerInterestRepository,
		atomicExecutor:          atomicExecutor,
	}
}
