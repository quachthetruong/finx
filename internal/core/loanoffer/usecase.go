package loanoffer

import (
	"context"
	"fmt"
	"time"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/atomicity"
	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/loanoffer/repository"
	loanOfferInterestRepo "financing-offer/internal/core/loanofferinterest/repository"
	"financing-offer/internal/funcs"
)

type UseCase interface {
	FindAllForInvestor(ctx context.Context, filter entity.LoanPackageOfferFilter) ([]entity.LoanPackageOffer, error)
	Create(ctx context.Context, loanOffer entity.LoanPackageOffer) (entity.LoanPackageOffer, error)
	InvestorCancel(ctx context.Context, investorId string, loanOfferId int64) error
	InvestorGetById(ctx context.Context, loanOfferId int64, investorId string) (entity.LoanPackageOffer, error)
	ExpireLoanOffers(ctx context.Context) error
}

type loanPackageOfferUseCase struct {
	repository                  repository.LoanPackageOfferRepository
	loanConfig                  config.LoanRequestConfig
	loanOfferInterestRepository loanOfferInterestRepo.LoanPackageOfferInterestRepository
	atomicExecutor              atomicity.AtomicExecutor
}

func (u *loanPackageOfferUseCase) FindAllForInvestor(ctx context.Context, filter entity.LoanPackageOfferFilter) ([]entity.LoanPackageOffer, error) {
	declinedRequestDisplayPeriod := -time.Duration(24*u.loanConfig.DeclinedRequestDisplayPeriod) * time.Hour
	res, err := u.repository.FindAllForInvestorWithRequestAndLine(ctx, filter)
	if err != nil {
		return res, fmt.Errorf("loanPackageOfferUseCase FindAll %w", err)
	}
	offers := make([]entity.LoanPackageOffer, 0, len(res))
	for _, offer := range res {
		if !shouldIncludeOffer(offer, declinedRequestDisplayPeriod) {
			continue
		}
		offers = append(offers, offer)
	}
	return offers, nil
}

func shouldIncludeOffer(offer entity.LoanPackageOffer, displayPeriod time.Duration) bool {
	interests := offer.LoanPackageOfferInterests
	// if request is declined, only show offer for n days
	if len(interests) == 0 && offer.CreatedAt.Before(time.Now().Add(displayPeriod)) {
		return false
	}
	if offer.IsExpired() {
		activated := false
		for _, i := range interests {
			if i.Status == entity.LoanPackageOfferInterestStatusLoanPackageCreated {
				activated = true
				break
			}
		}
		if !activated {
			return false
		}
	}
	return true
}

func (u *loanPackageOfferUseCase) InvestorCancel(ctx context.Context, investorId string, loanOfferId int64) error {
	offer, err := u.repository.FindByIdWithRequest(ctx, loanOfferId)
	if err != nil {
		return fmt.Errorf("loanPackageOfferUseCase InvestorCancel %w", err)
	}
	if investorId != offer.LoanPackageRequest.InvestorId {
		return apperrors.ErrInvalidInvestorId
	}
	if offer.IsExpired() {
		return apperrors.ErrOfferExpired
	}
	if err := u.loanOfferInterestRepository.CancelByOfferId(
		ctx, loanOfferId, investorId, entity.LoanPackageOfferCancelledReasonInvestor,
	); err != nil {
		return fmt.Errorf("loanPackageOfferUseCase InvestorCancel %w", err)
	}
	return nil
}

func (u *loanPackageOfferUseCase) Create(ctx context.Context, loanOffer entity.LoanPackageOffer) (entity.LoanPackageOffer, error) {
	res, err := u.repository.Create(ctx, loanOffer)
	if err != nil {
		return res, fmt.Errorf("loanPackageOfferUseCase Create %w", err)
	}
	return res, nil
}

func (u *loanPackageOfferUseCase) InvestorGetById(ctx context.Context, loanOfferId int64, investorId string) (entity.LoanPackageOffer, error) {
	offer, err := u.repository.InvestorGetById(ctx, loanOfferId)
	if err != nil {
		return offer, fmt.Errorf("loanPackageOfferUseCase InvestorGetById %w", err)
	}
	if offer.LoanPackageRequest.InvestorId != investorId {
		return offer, apperrors.ErrInvalidInvestorId
	}
	return offer, nil
}

func (u *loanPackageOfferUseCase) ExpireLoanOffers(ctx context.Context) error {
	offers, err := u.repository.GetExpiredOffers(ctx)
	if err != nil {
		return fmt.Errorf("loanPackageOfferUseCase ExpireLoanOffers %w", err)
	}
	expiredOfferIds := funcs.Map(
		offers, func(offer entity.LoanPackageOffer) int64 {
			return offer.Id
		},
	)
	if err := u.loanOfferInterestRepository.CancelExpiredOfferInterests(ctx, expiredOfferIds); err != nil {
		return fmt.Errorf("loanPackageOfferUseCase ExpireLoanOffers %w", err)
	}
	return nil
}

func NewUseCase(
	repository repository.LoanPackageOfferRepository,
	loanConfig config.LoanRequestConfig,
	loanOfferInterestRepository loanOfferInterestRepo.LoanPackageOfferInterestRepository,
	atomicExecutor atomicity.AtomicExecutor,
) UseCase {
	return &loanPackageOfferUseCase{
		repository:                  repository,
		loanConfig:                  loanConfig,
		loanOfferInterestRepository: loanOfferInterestRepository,
		atomicExecutor:              atomicExecutor,
	}
}
