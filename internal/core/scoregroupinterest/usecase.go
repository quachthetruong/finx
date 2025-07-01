package scoregroupinterest

import (
	"context"
	"fmt"
	"slices"

	"github.com/shopspring/decimal"

	"financing-offer/internal/core/entity"
	loanPackageRequestRepo "financing-offer/internal/core/loanpackagerequest/repository"
	"financing-offer/internal/core/scoregroupinterest/repository"
	symbolScoreRepo "financing-offer/internal/core/symbolscore/repository"
)

type UseCase interface {
	Create(ctx context.Context, groupRole entity.ScoreGroupInterest) (entity.ScoreGroupInterest, error)
	Update(ctx context.Context, groupRole entity.ScoreGroupInterest) (entity.ScoreGroupInterest, error)
	GetAll(ctx context.Context) ([]entity.ScoreGroupInterest, error)
	GetById(ctx context.Context, id int64) (entity.ScoreGroupInterest, error)
	Delete(ctx context.Context, id int64) (bool, error)
	GetForLoanPackageRequest(ctx context.Context, loanPackageRequestId int64) ([]entity.ScoreGroupInterest, error)
}

type useCase struct {
	loanPackageRequestRepository loanPackageRequestRepo.LoanPackageRequestRepository
	repository                   repository.ScoreGroupInterestRepository
	symbolScoreRepository        symbolScoreRepo.SymbolScoreRepository
}

func NewUseCase(
	repository repository.ScoreGroupInterestRepository,
	loanPackageRequestRepository loanPackageRequestRepo.LoanPackageRequestRepository,
	symbolScoreRepository symbolScoreRepo.SymbolScoreRepository,
) UseCase {
	return &useCase{
		repository:                   repository,
		loanPackageRequestRepository: loanPackageRequestRepository,
		symbolScoreRepository:        symbolScoreRepository,
	}
}

func (u *useCase) Create(ctx context.Context, groupRole entity.ScoreGroupInterest) (entity.ScoreGroupInterest, error) {
	res, err := u.repository.Create(ctx, groupRole)
	if err != nil {
		return res, fmt.Errorf("scoreGroupInterestUseCase Create %w", err)
	}
	return res, nil
}

func (u *useCase) Update(ctx context.Context, groupRole entity.ScoreGroupInterest) (entity.ScoreGroupInterest, error) {
	res, err := u.repository.Update(ctx, groupRole)
	if err != nil {
		return res, fmt.Errorf("scoreGroupInterestUseCase Update %w", err)
	}
	return res, nil
}

func (u *useCase) GetAll(ctx context.Context) ([]entity.ScoreGroupInterest, error) {
	res, err := u.repository.GetAll(ctx, entity.ScoreGroupInterestFilter{})
	if err != nil {
		return res, fmt.Errorf("scoreGroupInterestUseCase GetAll %w", err)
	}
	return res, nil
}

func (u *useCase) Delete(ctx context.Context, id int64) (bool, error) {
	res, err := u.repository.Delete(ctx, id)
	if err != nil {
		return res, fmt.Errorf("scoreGroupInterestUseCase Delete %w", err)
	}
	return res, nil
}

func (u *useCase) GetById(ctx context.Context, id int64) (entity.ScoreGroupInterest, error) {
	res, err := u.repository.GetById(ctx, id)
	if err != nil {
		return res, fmt.Errorf("scoreGroupInterestUseCase GetById %w", err)
	}
	return res, nil
}

func (u *useCase) GetForLoanPackageRequest(ctx context.Context, loanPackageRequestId int64) ([]entity.ScoreGroupInterest, error) {
	request, err := u.loanPackageRequestRepository.GetById(ctx, loanPackageRequestId, entity.LoanPackageFilter{})
	if err != nil {
		return nil, fmt.Errorf("scoreGroupInterestUseCase GetForLoanPackageRequest %w", err)
	}
	scoreGroupInterests, err := u.repository.GetAvailablePackageBySymbolId(ctx, request.SymbolId)
	if err != nil {
		return nil, fmt.Errorf("scoreGroupInterestUseCase GetForLoanPackageRequest %w", err)
	}
	slices.SortFunc(
		scoreGroupInterests, func(a, b entity.ScoreGroupInterest) int {
			return a.LimitAmount.Cmp(b.LimitAmount)
		},
	)
	res := make([]entity.ScoreGroupInterest, 0, len(scoreGroupInterests))
	smallestGreaterThanRequestAmount := decimal.Zero
	for _, v := range scoreGroupInterests {
		if v.LimitAmount.LessThanOrEqual(request.LimitAmount) {
			res = append(res, v)
		} else if smallestGreaterThanRequestAmount.IsZero() {
			smallestGreaterThanRequestAmount = v.LimitAmount
			res = append(res, v)
		} else if smallestGreaterThanRequestAmount.Equal(v.LimitAmount) {
			res = append(res, v)
		}
	}
	slices.Reverse(res)
	return res, nil
}
