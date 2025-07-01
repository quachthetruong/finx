package loancontract

import (
	"financing-offer/internal/atomicity"
	loanContractRepo "financing-offer/internal/core/loancontract/repository"
	loanPackageOfferRepo "financing-offer/internal/core/loanoffer/repository"
	"financing-offer/internal/core/loanofferinterest/repository"
	loanPackageRequestRepo "financing-offer/internal/core/loanpackagerequest/repository"
)

type UseCase interface{}

type useCase struct {
	atomicExecutor                    atomicity.AtomicExecutor
	loanContractPersistenceRepository loanContractRepo.LoanContractPersistenceRepository
	loanPackageRequestRepository      loanPackageRequestRepo.LoanPackageRequestRepository
	offerInterestRepository           repository.LoanPackageOfferInterestRepository
	offerRepository                   loanPackageOfferRepo.LoanPackageOfferRepository
}

func NewUseCase(
	atomicExecutor atomicity.AtomicExecutor,
	loanContractPersistenceRepository loanContractRepo.LoanContractPersistenceRepository,
	loanPackageRequestRepository loanPackageRequestRepo.LoanPackageRequestRepository,
	offerInterestRepository repository.LoanPackageOfferInterestRepository,
	offerRepository loanPackageOfferRepo.LoanPackageOfferRepository,
) UseCase {
	return &useCase{
		atomicExecutor:                    atomicExecutor,
		loanContractPersistenceRepository: loanContractPersistenceRepository,
		loanPackageRequestRepository:      loanPackageRequestRepository,
		offerInterestRepository:           offerInterestRepository,
		offerRepository:                   offerRepository,
	}
}
