package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type InvestorAccountRepository interface {
	GetByAccountNo(ctx context.Context, accountNo string) (entity.InvestorAccount, error)
	Update(ctx context.Context, account entity.InvestorAccount) (entity.InvestorAccount, error)
	Create(ctx context.Context, account entity.InvestorAccount) (entity.InvestorAccount, error)
}
