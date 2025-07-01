package repository

import (
	"context"

	"financing-offer/internal/core/entity"
	"financing-offer/pkg/querymod"
)

type LoanPackageOfferInterestRepository interface {
	GetWithFilter(ctx context.Context, filter entity.OfferInterestFilter) ([]entity.LoanPackageOfferInterest, error)
	CountWithFilter(ctx context.Context, filter entity.OfferInterestFilter) (int64, error)
	BulkCreate(ctx context.Context, loanPackageOfferInterests []entity.LoanPackageOfferInterest) ([]entity.LoanPackageOfferInterest, error)
	Create(ctx context.Context, loanPackageOfferInterest entity.LoanPackageOfferInterest) (entity.LoanPackageOfferInterest, error)
	GetById(ctx context.Context, id int64, opts ...querymod.GetOption) (entity.LoanPackageOfferInterest, error)
	GetByIds(ctx context.Context, ids []int64, opts ...querymod.GetOption) ([]entity.LoanPackageOfferInterest, error)
	UpdateStatus(ctx context.Context, ids []int64, status entity.LoanPackageOfferInterestStatus) error
	Update(ctx context.Context, offerInterest entity.LoanPackageOfferInterest) (entity.LoanPackageOfferInterest, error)
	CancelByOfferId(ctx context.Context, offerId int64, cancelledBy string, cancelledReason entity.CancelledReason) error
	CancelExpiredOfferInterests(ctx context.Context, offerIds []int64) error
	GetRequestBasedLoanOfferInterests(ctx context.Context) ([]entity.LoanPackageOfferInterest, error)
	GetByOfferIdWithLock(ctx context.Context, offerId int64) ([]entity.LoanPackageOfferInterest, error)
}

type LoanPackageOfferInterestEventRepository interface {
	NotifyLoanPackageOfferReady(ctx context.Context, data entity.LoanPackageOfferReadyNotify) error
	NotifyDerivativeLoanPackageOfferReady(ctx context.Context, data entity.DerivativeLoanPackageOfferReadyNotify) error
	CreateMarginLoanPackage(ctx context.Context, state entity.AssignmentState) error
}
