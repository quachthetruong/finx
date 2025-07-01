package repository

import (
	"context"

	"financing-offer/internal/core/entity"
	"financing-offer/pkg/querymod"
)

type LoanContractPersistenceRepository interface {
	Create(ctx context.Context, loanContract entity.LoanContract) (entity.LoanContract, error)
	BulkCreate(ctx context.Context, loanContracts []entity.LoanContract) error
	GetById(ctx context.Context, id int64, opts ...querymod.GetOption) (entity.LoanContract, error)
	GetInvestorActiveContract(ctx context.Context, investorId string, symbolId int64) (entity.LoanContract, error)
}
