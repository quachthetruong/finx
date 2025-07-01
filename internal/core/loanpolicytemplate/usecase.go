package loanpolicytemplate

import (
	"context"
	"financing-offer/internal/apperrors"
	"fmt"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/loanpolicytemplate/repository"
)

type UseCase interface {
	GetAll(ctx context.Context) ([]entity.LoanPolicyTemplate, error)
	Create(ctx context.Context, template entity.LoanPolicyTemplate) (entity.LoanPolicyTemplate, error)
	Update(ctx context.Context, template entity.LoanPolicyTemplate) (entity.LoanPolicyTemplate, error)
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (entity.LoanPolicyTemplate, error)
}

func NewUseCase(
	loanPolicyTemplateRepository repository.LoanPolicyTemplateRepository,
) UseCase {
	return &loanPolicyTemplateUseCase{
		loanPolicyTemplateRepository: loanPolicyTemplateRepository,
	}
}

type loanPolicyTemplateUseCase struct {
	loanPolicyTemplateRepository repository.LoanPolicyTemplateRepository
}

func (u *loanPolicyTemplateUseCase) GetById(ctx context.Context, id int64) (entity.LoanPolicyTemplate, error) {
	loanPolicy, err := u.loanPolicyTemplateRepository.GetById(ctx, id)
	if err != nil {
		return entity.LoanPolicyTemplate{}, fmt.Errorf("loanPolicyTemplateUseCase GetById %w", err)
	}
	return loanPolicy, nil
}

func (u *loanPolicyTemplateUseCase) GetAll(ctx context.Context) ([]entity.LoanPolicyTemplate, error) {
	loanPolicies, err := u.loanPolicyTemplateRepository.GetAll(ctx)
	if err != nil {
		return []entity.LoanPolicyTemplate{}, fmt.Errorf("loanPolicyTemplateUseCase GetAll %w", err)
	}
	return loanPolicies, nil
}

func (u *loanPolicyTemplateUseCase) Create(ctx context.Context, template entity.LoanPolicyTemplate) (entity.LoanPolicyTemplate, error) {
	errTemplate := "loanPolicyTemplateUseCase Create %w"
	res, err := u.loanPolicyTemplateRepository.Create(ctx, template)
	if err != nil {
		return res, fmt.Errorf(errTemplate, err)
	}
	return res, nil
}

func (u *loanPolicyTemplateUseCase) Update(ctx context.Context, template entity.LoanPolicyTemplate) (entity.LoanPolicyTemplate, error) {
	errTemplate := "loanPolicyTemplateUseCase Update %w"
	_, err := u.loanPolicyTemplateRepository.GetById(ctx, template.Id)
	if err != nil {
		return entity.LoanPolicyTemplate{}, fmt.Errorf(errTemplate, apperrors.ErrorInvalidLoanPolicyTemplateId)
	}
	res, err := u.loanPolicyTemplateRepository.Update(ctx, template)
	if err != nil {
		return res, fmt.Errorf(errTemplate, err)
	}
	return res, nil
}

func (u *loanPolicyTemplateUseCase) Delete(ctx context.Context, id int64) error {
	errTemplate := "loanPolicyTemplateUseCase Delete %w"
	_, err := u.loanPolicyTemplateRepository.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf(errTemplate, apperrors.ErrorInvalidLoanPolicyTemplateId)
	}
	err = u.loanPolicyTemplateRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("loanPolicyTemplateUseCase Delete %w", err)
	}
	return nil
}
