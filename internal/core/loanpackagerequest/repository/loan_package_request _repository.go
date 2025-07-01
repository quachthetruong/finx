package repository

import (
	"context"

	"github.com/shopspring/decimal"

	"financing-offer/internal/core/entity"
	"financing-offer/pkg/querymod"
)

type LoanPackageRequestRepository interface {
	GetAll(ctx context.Context, filter entity.LoanPackageFilter) ([]entity.LoanPackageRequest, error)
	GetAllUnderlyingRequests(ctx context.Context, filter entity.UnderlyingLoanPackageFilter) ([]entity.UnderlyingLoanPackageRequest, error)
	Count(ctx context.Context, filter entity.LoanPackageFilter) (int64, error)
	GetById(ctx context.Context, id int64, filter entity.LoanPackageFilter, opts ...querymod.GetOption) (entity.LoanPackageRequest, error)
	Create(ctx context.Context, loanPackageRequest entity.LoanPackageRequest) (entity.LoanPackageRequest, error)
	Update(ctx context.Context, loanPackageRequest entity.LoanPackageRequest) (entity.LoanPackageRequest, error)
	Delete(ctx context.Context, id int64) error
	SaveLoggedRequest(ctx context.Context, request entity.LoggedRequest) (entity.LoggedRequest, error)
	LockAllPendingRequestByMaxPercent(ctx context.Context, maximumLoanRate decimal.Decimal) ([]entity.LoanPackageRequest, error)
	UpdateStatusByLoanRequestIds(ctx context.Context, loanRequestIds []int64, status entity.LoanPackageRequestStatus) ([]entity.LoanPackageRequest, error)
	LockAndReturnAllPendingRequestBySymbolId(ctx context.Context, symbolId int64) ([]entity.LoanPackageRequest, error)
	UpdateStatusById(ctx context.Context, id int64, status entity.LoanPackageRequestStatus) (entity.LoanPackageRequest, error)
}

type LoanPackageRequestEventRepository interface {
	NotifyRequestDeclined(ctx context.Context, data entity.LoanPackageRequestDeclinedNotify) error
	NotifyDerivativeRequestDeclined(ctx context.Context, data entity.LoanPackageDerivativeRequestDeclinedNotify) error
	NotifyOfflineConfirmation(ctx context.Context, data entity.RequestOfflineConfirmation) error
	NotifyDerivativeOfflineConfirmation(ctx context.Context, data entity.DerivativeRequestOfflineConfirmation) error
	NotifyOnlineConfirmation(ctx context.Context, data entity.RequestOnlineConfirmationNotify) error
}
