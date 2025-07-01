package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type OrderServiceRepository interface {
	GetAllAccountLoanPackages(ctx context.Context, accountNo string) ([]entity.AccountLoanPackage, error)
	GetAccountByAccountNoAndCustodyCode(ctx context.Context, custodyCode string, accountNo string) (entity.OrderServiceAccount, error)
}
