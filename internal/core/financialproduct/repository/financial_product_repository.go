package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type FinancialProductRepository interface {
	GetAllAccountDetail(ctx context.Context, investorId string) (accounts []entity.FinancialAccountDetail, err error)
	GetAllAccountDetailByCustodyCode(ctx context.Context, custodyCode string) (accounts []entity.FinancialAccountDetail, err error)
	GetLoanPackageDetails(ctx context.Context, loanPackageIds []int64) ([]entity.FinancialProductLoanPackage, error)
	AssignLoanPackage(ctx context.Context, accountNo string, loanId int64, assetType entity.AssetType) (loanPackageAccountId int64, err error)
	AssignLoanPackageOrGetLoanPackageAccountId(ctx context.Context, accountNo string, loanId int64, assetType entity.AssetType) (loanPackageAccountId int64, err error)
	GetLoanPackageDetail(ctx context.Context, loanPackageId int64) (entity.FinancialProductLoanPackage, error)
	GetLoanPackageDerivative(ctx context.Context, loanPackageId int64) (entity.FinancialProductLoanPackageDerivative, error)
	GetLoanPackageAccountIdByAccountNoAndLoanPackageId(ctx context.Context, accountNo string, loanPackageId int64, assetType entity.AssetType) (int64, error)
	GetLoanRatesByIds(ctx context.Context, loanRateIds []int64) ([]entity.LoanRate, error)
	GetLoanRateDetail(ctx context.Context, loanRateId int64) (entity.LoanRate, error)
	GetLoanProducts(ctx context.Context, filter entity.MarginProductFilter) ([]entity.MarginProduct, error)
	GetMarginBasketsByIds(ctx context.Context, ids []int64) ([]entity.MarginBasket, error)
}
