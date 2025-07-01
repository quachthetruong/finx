package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type LoanPackageOfferRepository interface {
	FindAllForInvestorWithRequestAndLine(ctx context.Context, filter entity.LoanPackageOfferFilter) ([]entity.LoanPackageOffer, error)
	FindByIdWithRequest(ctx context.Context, id int64) (entity.LoanPackageOffer, error)
	Create(ctx context.Context, loanPackageOffer entity.LoanPackageOffer) (entity.LoanPackageOffer, error)
	InvestorGetById(ctx context.Context, id int64) (entity.LoanPackageOffer, error)
	GetExpiredOffers(ctx context.Context) ([]entity.LoanPackageOffer, error)
	BulkCreate(ctx context.Context, loanPackageOffers []entity.LoanPackageOffer) ([]entity.LoanPackageOffer, error)
}
