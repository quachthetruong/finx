package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type LoanPolicyTemplateRepository interface {
	GetAll(ctx context.Context) ([]entity.LoanPolicyTemplate, error)
	Create(ctx context.Context, template entity.LoanPolicyTemplate) (entity.LoanPolicyTemplate, error)
	Update(ctx context.Context, template entity.LoanPolicyTemplate) (entity.LoanPolicyTemplate, error)
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (entity.LoanPolicyTemplate, error)
	GetByIds(ctx context.Context, ids []int64) ([]entity.LoanPolicyTemplate, error)
}
