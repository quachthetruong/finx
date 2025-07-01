package investor

import (
	"context"
	"fmt"

	"financing-offer/internal/core/entity"
	financialProductRepo "financing-offer/internal/core/financialproduct/repository"
	"financing-offer/internal/core/investor/repository"
)

type UseCase interface {
	FillInvestorIdsFromRequests(ctx context.Context) (int, error)
}

type useCase struct {
	investorPersistenceRepository repository.InvestorPersistenceRepository
	financialProductRepository    financialProductRepo.FinancialProductRepository
}

func (u *useCase) FillInvestorIdsFromRequests(ctx context.Context) (int, error) {
	errorTemplate := "investorUseCase FillInvestorIdsFromRequests %w"
	investorIds, err := u.investorPersistenceRepository.GetAllUniqueInvestorIdsFromRequests(ctx)
	if err != nil {
		return 0, fmt.Errorf(errorTemplate, err)
	}
	investorsToCreate := make([]entity.Investor, 0, len(investorIds))
	for _, investorId := range investorIds {
		investorsToCreate = append(investorsToCreate, entity.Investor{InvestorId: investorId})
	}
	if err := u.investorPersistenceRepository.BulkCreate(ctx, investorsToCreate); err != nil {
		return 0, fmt.Errorf(errorTemplate, err)
	}
	filledInvestorIdsCount, err := u.fillCustodyCodeFromExternalSystem(ctx)
	if err != nil {
		return 0, fmt.Errorf(errorTemplate, err)
	}
	return filledInvestorIdsCount, nil
}

func (u *useCase) fillCustodyCodeFromExternalSystem(ctx context.Context) (int, error) {
	errorTemplate := "investorUseCase fillCustodyCodeFromExternalSystem %w"
	existedInvestorIds, err := u.investorPersistenceRepository.GetAllInvestorIdsForMigration(ctx)
	if err != nil {
		return 0, fmt.Errorf(errorTemplate, err)
	}
	for _, investorId := range existedInvestorIds {
		investorDetails, err := u.financialProductRepository.GetAllAccountDetail(ctx, investorId)
		if err != nil {
			return 0, fmt.Errorf(errorTemplate, err)
		}
		if len(investorDetails) == 0 {
			continue
		}
		custodyCode := investorDetails[0].Custody
		if _, err := u.investorPersistenceRepository.Update(
			ctx, entity.Investor{InvestorId: investorId, CustodyCode: custodyCode},
		); err != nil {
			return 0, fmt.Errorf(errorTemplate, err)
		}
	}
	return len(existedInvestorIds), nil
}

func NewUseCase(investorPersistenceRepository repository.InvestorPersistenceRepository, financialProductRepository financialProductRepo.FinancialProductRepository) UseCase {
	return &useCase{investorPersistenceRepository: investorPersistenceRepository, financialProductRepository: financialProductRepository}
}
