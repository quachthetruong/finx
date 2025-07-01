package loanofferinterest

import (
	"context"
	"financing-offer/internal/config"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/atomicity"
	"financing-offer/internal/core"
	"financing-offer/internal/core/entity"
	financialProductRepo "financing-offer/internal/core/financialproduct/repository"
	loanContractRepo "financing-offer/internal/core/loancontract/repository"
	loanPackageOfferRepo "financing-offer/internal/core/loanoffer/repository"
	"financing-offer/internal/core/loanofferinterest/repository"
	loanPolicyRepo "financing-offer/internal/core/loanpolicytemplate/repository"
	submissionSheetRepo "financing-offer/internal/core/submissionsheet/repository"
	symbolRepo "financing-offer/internal/core/symbol/repository"
	"financing-offer/internal/funcs"
	"financing-offer/pkg/querymod"
)

type UseCase interface {
	GetAll(ctx context.Context, filter entity.OfferInterestFilter) ([]entity.LoanPackageOfferInterest, core.PagingMetaData, error)
	InvestorConfirmLoanPackageInterest(ctx context.Context, ids []int64, investorId string) error
	AdminAssignLoanIdByOfferId(ctx context.Context, offerId int64, loanId int64) error
	AdminCancelLoanPackageInterestByOfferId(ctx context.Context, offerId int64, canceler string) error
	InvestorCancelLoanPackageInterest(ctx context.Context, id int64, investorId string) error
	SyncLoanPackageData(ctx context.Context) (int, error)
	CreateAssignedLoanOfferInterestLoanContract(
		ctx context.Context,
		loanPackageOfferInterestId,
		loanPackageAccountId,
		loanProductIdRef int64,
		loanPackage entity.FinancialProductLoanPackage,
	) (entity.LoanContract, error)
}

type useCase struct {
	atomicExecutor             atomicity.AtomicExecutor
	repository                 repository.LoanPackageOfferInterestRepository
	loanOfferRepository        loanPackageOfferRepo.LoanPackageOfferRepository
	loanContractRepository     loanContractRepo.LoanContractPersistenceRepository
	financialProductRepository financialProductRepo.FinancialProductRepository
	eventRepository            repository.LoanPackageOfferInterestEventRepository
	errorService               apperrors.Service
	symbolRepository           symbolRepo.SymbolRepository
	submissionSheetRepository  submissionSheetRepo.SubmissionSheetRepository
	policyTemplateRepository   loanPolicyRepo.LoanPolicyTemplateRepository
	appConfig                  config.AppConfig
}

func (u *useCase) CreateAssignedLoanOfferInterestLoanContract(ctx context.Context, loanPackageOfferInterestId, loanPackageAccountId, loanProductIdRef int64, loanPackage entity.FinancialProductLoanPackage) (entity.LoanContract, error) {
	var (
		createdLoanContract entity.LoanContract
		request             entity.LoanPackageRequest
	)
	transactionFunc := func(ctx context.Context) error {
		loanOfferInterest, err := u.repository.GetById(ctx, loanPackageOfferInterestId)
		if err != nil {
			return err
		}
		offerWithRequest, err := u.loanOfferRepository.FindByIdWithRequest(ctx, loanOfferInterest.LoanPackageOfferId)
		if err != nil {
			return err
		}
		request = *offerWithRequest.LoanPackageRequest
		if !slices.Contains(
			loanOfferInterest.Status.NextStatuses(offerWithRequest.FlowType),
			entity.LoanPackageOfferInterestStatusLoanPackageCreated,
		) {
			return apperrors.ErrInvalidLoanPackageOfferInterestStatus
		}

		loanOfferInterest.Status = entity.LoanPackageOfferInterestStatusLoanPackageCreated
		loanOfferInterest.LoanID = loanPackage.Id
		loanOfferInterest.FeeRate = loanPackage.BuyingFeeRate
		loanOfferInterest.InterestRate = loanPackage.InterestRate
		loanOfferInterest.LoanRate = decimal.NewFromInt(1).Sub(loanPackage.InitialRate)
		loanOfferInterest.Term = loanPackage.Term

		if _, err = u.repository.Update(ctx, loanOfferInterest); err != nil {
			return err
		}
		loanContract := entity.LoanContract{
			LoanOfferInterestId:  loanOfferInterest.Id,
			SymbolId:             offerWithRequest.LoanPackageRequest.SymbolId,
			InvestorId:           offerWithRequest.LoanPackageRequest.InvestorId,
			AccountNo:            offerWithRequest.LoanPackageRequest.AccountNo,
			LoanId:               loanPackage.Id,
			LoanProductIdRef:     loanProductIdRef,
			LoanPackageAccountId: loanPackageAccountId,
		}
		createdLoanContract, err = u.loanContractRepository.Create(ctx, loanContract)
		if err != nil {
			return err
		}
		return nil
	}
	if err := u.atomicExecutor.Execute(ctx, transactionFunc); err != nil {
		return entity.LoanContract{}, fmt.Errorf(
			"loanContractUseCase CreateAssignedLoanOfferInterestLoanContract %w", err,
		)
	}
	u.errorService.Go(
		ctx, func() error {
			return u.NotifyLoanPackageReady(
				atomicity.WithIgnoreTx(ctx), request, AssignedLoanPackageAccount{
					LoanPackageOfferInterestId: loanPackageOfferInterestId,
					LoanPackageId:              loanPackage.Id,
					InterestRate:               loanPackage.InterestRate,
					LoanRate:                   decimal.NewFromInt(1).Sub(loanPackage.InitialRate),
					InitialRate:                loanPackage.InitialRate,
				},
			)
		},
	)
	return createdLoanContract, nil
}

func (u *useCase) AdminCancelLoanPackageInterestByOfferId(ctx context.Context, offerId int64, canceler string) error {
	txErr := u.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			offer, err := u.loanOfferRepository.FindByIdWithRequest(tc, offerId)
			if err != nil {
				return err
			}
			offerLines, err := u.repository.GetByOfferIdWithLock(tc, offerId)
			if err != nil {
				return err
			}
			if len(offerLines) == 0 {
				return apperrors.ErrorNotFoundLoanPackageOfferInterest
			}
			offerLine := offerLines[0]
			if !slices.Contains(
				offerLine.Status.NextStatuses(offer.FlowType), entity.LoanPackageOfferInterestStatusCancelled,
			) {
				return apperrors.ErrInvalidLoanPackageOfferInterestStatus
			}
			offerLine.Status = entity.LoanPackageOfferInterestStatusCancelled
			offerLine.CancelledBy = canceler
			offerLine.CancelledAt = time.Now()
			offerLine.CancelledReason = entity.LoanPackageOfferCancelledReasonAdmin
			if _, err := u.repository.Update(tc, offerLine); err != nil {
				return err
			}
			return nil
		},
	)
	if txErr != nil {
		return fmt.Errorf("loanOfferInterestUseCase AdminCancelLoanPackageOfferInterest %w", txErr)
	}
	return nil
}

func (u *useCase) InvestorCancelLoanPackageInterest(ctx context.Context, id int64, investorId string) error {
	txErr := u.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			offerLine, err := u.repository.GetById(tc, id, querymod.WithLock())
			if err != nil {
				return err
			}
			offer, err := u.loanOfferRepository.FindByIdWithRequest(tc, offerLine.LoanPackageOfferId)
			if err != nil {
				return err
			}
			if !slices.Contains(
				offerLine.Status.NextStatuses(offer.FlowType), entity.LoanPackageOfferInterestStatusCancelled,
			) {
				return apperrors.ErrInvalidLoanPackageOfferInterestStatus
			}
			request := offer.LoanPackageRequest
			if request.InvestorId != investorId {
				return apperrors.ErrInvestorNotAllowed
			}
			offerLine.Status = entity.LoanPackageOfferInterestStatusCancelled
			offerLine.CancelledBy = investorId
			offerLine.CancelledAt = time.Now()
			offerLine.CancelledReason = entity.LoanPackageOfferCancelledReasonInvestor
			if _, err := u.repository.Update(tc, offerLine); err != nil {
				return err
			}
			return nil
		},
	)
	if txErr != nil {
		return fmt.Errorf("loanOfferInterestUseCase InvestorCancelLoanPackageOfferInterest %w", txErr)
	}
	return nil
}

func (u *useCase) AdminAssignLoanIdByOfferId(ctx context.Context, offerId int64, loanId int64) error {
	var (
		assignedOfferLine entity.LoanPackageOfferInterest
		request           entity.LoanPackageRequest
	)
	txErr := u.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			offer, offerLine, err := u.getAndVerifyOfferWithRequest(tc, offerId)
			if err != nil {
				return err
			}
			request = *offer.LoanPackageRequest
			assignedOfferLine, err = u.fillUnderlyingLoanPackageToOfferLine(tc, loanId, offerLine)
			if err != nil {
				return err
			}
			loanPackageAccountId, err := u.financialProductRepository.AssignLoanPackageOrGetLoanPackageAccountId(
				ctx, request.AccountNo, loanId, request.AssetType,
			)
			if err != nil {
				return err
			}
			contract := entity.LoanContract{
				LoanOfferInterestId:  assignedOfferLine.Id,
				SymbolId:             request.SymbolId,
				InvestorId:           request.InvestorId,
				AccountNo:            request.AccountNo,
				LoanId:               loanId,
				LoanPackageAccountId: loanPackageAccountId,
			}
			if request.Type == entity.LoanPackageRequestTypeGuaranteed {
				contract.GuaranteedEndAt = time.Now().Add(time.Duration(request.GuaranteedDuration) * 24 * time.Hour)
			}
			if _, err := u.loanContractRepository.Create(tc, contract); err != nil {
				return err
			}
			return nil
		},
	)
	if txErr != nil {
		return fmt.Errorf("loanOfferInterestUseCase AdminAssignLoanId %w", txErr)
	}
	u.errorService.Go(
		ctx, func() error {
			return u.NotifyLoanPackageReady(
				atomicity.WithIgnoreTx(ctx), request, AssignedLoanPackageAccount{
					LoanPackageOfferInterestId: assignedOfferLine.Id,
					LoanPackageId:              assignedOfferLine.LoanID,
					InterestRate:               assignedOfferLine.InterestRate,
					LoanRate:                   assignedOfferLine.LoanRate,
					InitialRate:                assignedOfferLine.InitialRate,
				},
			)
		},
	)
	return nil
}

func (u *useCase) getAndVerifyOfferWithRequest(atomicContext context.Context, offerId int64) (entity.LoanPackageOffer, entity.LoanPackageOfferInterest, error) {
	offer, err := u.loanOfferRepository.FindByIdWithRequest(atomicContext, offerId)
	if err != nil {
		return offer, entity.LoanPackageOfferInterest{}, err
	}
	if offer.FlowType != entity.FLowTypeDnseOffline {
		return offer, entity.LoanPackageOfferInterest{}, apperrors.ErrInvalidFlowType(offer.FlowType)
	}
	offerLines, err := u.repository.GetByOfferIdWithLock(atomicContext, offerId)
	if err != nil {
		return offer, entity.LoanPackageOfferInterest{}, err
	}
	if len(offerLines) == 0 {
		return offer, entity.LoanPackageOfferInterest{}, apperrors.ErrorNotFoundLoanPackageOfferInterest
	}
	return offer, offerLines[0], nil
}

func (u *useCase) fillUnderlyingLoanPackageToOfferLine(
	ctx context.Context,
	loanId int64,
	offerLine entity.LoanPackageOfferInterest,
) (entity.LoanPackageOfferInterest, error) {
	if !slices.Contains(
		offerLine.Status.NextStatuses(entity.FLowTypeDnseOffline),
		entity.LoanPackageOfferInterestStatusLoanPackageCreated,
	) {
		return entity.LoanPackageOfferInterest{}, apperrors.ErrInvalidLoanPackageOfferInterestStatus
	}
	switch offerLine.AssetType {
	case entity.AssetTypeUnderlying:
		loanPackage, err := u.financialProductRepository.GetLoanPackageDetail(ctx, loanId)
		if err != nil {
			return entity.LoanPackageOfferInterest{}, err
		}
		offerLine.Status = entity.LoanPackageOfferInterestStatusLoanPackageCreated
		offerLine.LoanID = loanPackage.Id
		offerLine.FeeRate = loanPackage.BuyingFeeRate
		offerLine.InterestRate = loanPackage.InterestRate
		offerLine.LoanRate = decimal.NewFromInt(1).Sub(loanPackage.InitialRate)
		offerLine.Term = loanPackage.Term
	case entity.AssetTypeDerivative:
		loanPackage, err := u.financialProductRepository.GetLoanPackageDerivative(ctx, loanId)
		if err != nil {
			return entity.LoanPackageOfferInterest{}, err
		}
		offerLine.Status = entity.LoanPackageOfferInterestStatusLoanPackageCreated
		offerLine.LoanID = loanPackage.Id
		offerLine.InitialRate = loanPackage.InitialRate
	default:
		return entity.LoanPackageOfferInterest{}, apperrors.ErrMismatchAssetType
	}
	if _, err := u.repository.Update(ctx, offerLine); err != nil {
		return entity.LoanPackageOfferInterest{}, err
	}
	return offerLine, nil
}

func (u *useCase) NotifyLoanPackageReady(
	ctx context.Context,
	request entity.LoanPackageRequest,
	assignedLoanPackageAccount AssignedLoanPackageAccount,
) error {
	symbol, err := u.symbolRepository.GetById(ctx, request.SymbolId)
	if err != nil {
		return err
	}
	// only include accountNo if investor has more than 1 account
	accountNoDesc := ""
	accounts, err := u.financialProductRepository.GetAllAccountDetail(ctx, request.InvestorId)
	if err != nil {
		return err
	}
	if len(accounts) > 1 {
		for _, account := range accounts {
			if account.AccountNo == request.AccountNo {
				accountNoDesc = account.AccountTypeName
				break
			}
		}
	}
	switch request.AssetType {
	case entity.AssetTypeUnderlying:
		return u.eventRepository.NotifyLoanPackageOfferReady(
			ctx, entity.LoanPackageOfferReadyNotify{
				InvestorId:    request.InvestorId,
				RequestName:   fmt.Sprintf("%s-%d", symbol.Symbol, request.Id),
				AccountNo:     request.AccountNo,
				AccountNoDesc: accountNoDesc,
				Symbol:        symbol.Symbol,
				LoanRate:      assignedLoanPackageAccount.LoanRate,
				LoanType:      request.Type,
				InterestRate:  assignedLoanPackageAccount.InterestRate,
				LoanPackageId: assignedLoanPackageAccount.LoanPackageId,
				CreatedAt:     time.Now(),
			},
		)
	case entity.AssetTypeDerivative:
		return u.eventRepository.NotifyDerivativeLoanPackageOfferReady(
			ctx, entity.DerivativeLoanPackageOfferReadyNotify{
				InvestorId:    request.InvestorId,
				RequestName:   fmt.Sprintf("%s-%d", symbol.Symbol, request.Id),
				AccountNo:     request.AccountNo,
				AccountNoDesc: accountNoDesc,
				Symbol:        symbol.Symbol,
				LoanPackageId: assignedLoanPackageAccount.LoanPackageId,
				AssetType:     request.AssetType.String(),
				CreatedAt:     time.Now(),
			},
		)
	default:
		return apperrors.ErrMismatchAssetType
	}
}

func (u *useCase) InvestorConfirmLoanPackageInterest(ctx context.Context, ids []int64, investorId string) error {
	if len(ids) == 0 {
		return nil
	}
	offerLines, err := u.repository.GetByIds(ctx, ids)
	if err != nil {
		return err
	}
	offer, existedMoLoanPackage, err := u.verifyLoanOfferLines(ctx, offerLines, investorId)
	if err != nil {
		return err
	}
	request := *offer.LoanPackageRequest
	if existedMoLoanPackage {
		return u.assignExistedLoanPackages(ctx, request, offer, offerLines, investorId)
	} else {
		return u.createAndAssignNewLoanPackage(ctx, request, offerLines[0])
	}
}

func (u *useCase) createAndAssignNewLoanPackage(
	ctx context.Context,
	request entity.LoanPackageRequest,
	offerLine entity.LoanPackageOfferInterest,
) error {
	if err := u.repository.UpdateStatus(
		ctx, []int64{offerLine.Id}, entity.LoanPackageOfferInterestStatusCreatingLoanPackage,
	); err != nil {
		return err
	}
	submissionSheet, err := u.submissionSheetRepository.GetDetailById(ctx, offerLine.SubmissionSheetDetailId)
	if err != nil {
		return err
	}

	symbol, err := u.symbolRepository.GetById(ctx, request.SymbolId)
	if err != nil {
		return err
	}

	if err := u.eventRepository.CreateMarginLoanPackage(
		ctx, entity.AssignmentState{
			Submission: entity.Submission{
				Symbol:                     symbol.Symbol,
				LoanPackageOfferInterestId: offerLine.Id,
				LoanPackageRequestId:       request.Id,
				AccountNo:                  request.AccountNo,
				LoanRate:                   submissionSheet.LoanRate,
				Templates: funcs.Map(
					submissionSheet.LoanPolicies, func(policy entity.LoanPolicySnapShot) entity.TemplateWithProductRate {
						return entity.TemplateWithProductRate{
							LoanPolicySnapShot:       policy,
							AllowedOverdueLoanInDays: policy.AllowedOverdueLoanInDays,
							ProductRate:              policy.InitialRate,
							ProductRateForWithdraw:   policy.InitialRateForWithdraw,
						}
					},
				),
				FirmBuyingFeeRate:  submissionSheet.FirmBuyingFee,
				FirmSellingFeeRate: submissionSheet.FirmSellingFee,
				TransferFee:        submissionSheet.TransferFee,
				ProductCategoryId:  u.appConfig.ProductCategoryId,
			},
		},
	); err != nil {
		return err
	}
	return nil
}

func (u *useCase) assignExistedLoanPackages(
	ctx context.Context,
	request entity.LoanPackageRequest,
	offer entity.LoanPackageOffer,
	offerLines []entity.LoanPackageOfferInterest,
	investorId string,
) error {
	ids := funcs.Map(
		offerLines, func(offerLine entity.LoanPackageOfferInterest) int64 {
			return offerLine.Id
		},
	)
	txFunc := func(ctx context.Context) error {
		if err := u.repository.UpdateStatus(
			ctx, ids, entity.LoanPackageOfferInterestStatusLoanPackageCreated,
		); err != nil {
			return err
		}
		if err := u.repository.CancelByOfferId(
			ctx, offer.Id, investorId, entity.LoanPackageOfferCancelledReasonInvestor,
		); err != nil {
			return err
		}
		assignedLoanPackages := make([]AssignedLoanPackageAccount, 0, len(offerLines))
		assignedLoanPackages, err := u.assignMultipleLoanPackages(
			ctx, request.AccountNo, offerLines, request.AssetType,
		)
		if err != nil {
			return err
		}
		if len(assignedLoanPackages) > 0 {
			if err := u.loanContractRepository.BulkCreate(
				ctx, prepareContracts(request, investorId, assignedLoanPackages),
			); err != nil {
				return err
			}
		}
		for _, assignedLoanPackage := range assignedLoanPackages {
			assignedLoanPackage := assignedLoanPackage
			u.errorService.Go(
				ctx, func() error {
					return u.NotifyLoanPackageReady(
						atomicity.WithIgnoreTx(ctx), request, assignedLoanPackage,
					)
				},
			)
		}
		return nil
	}
	if txErr := u.atomicExecutor.Execute(ctx, txFunc); txErr != nil {
		return fmt.Errorf("loanOfferInterestUseCase assignExistedLoanPackages %w", txErr)
	}
	return nil
}

func prepareContracts(
	request entity.LoanPackageRequest,
	investorId string,
	assignedLoanPackages []AssignedLoanPackageAccount,
) []entity.LoanContract {
	contractsToCreate := make([]entity.LoanContract, 0, len(assignedLoanPackages))
	guaranteedEndTime := time.Time{}
	if request.Type == entity.LoanPackageRequestTypeGuaranteed {
		guaranteedEndTime = time.Now().Add(time.Duration(request.GuaranteedDuration) * 24 * time.Hour)
	}
	for _, assignedLoanPackage := range assignedLoanPackages {
		contractsToCreate = append(
			contractsToCreate, entity.LoanContract{
				LoanOfferInterestId:  assignedLoanPackage.LoanPackageOfferInterestId,
				SymbolId:             request.SymbolId,
				InvestorId:           investorId,
				AccountNo:            request.AccountNo,
				LoanId:               assignedLoanPackage.LoanPackageId,
				LoanPackageAccountId: assignedLoanPackage.CreatedLoanPackageAccountId,
				GuaranteedEndAt:      guaranteedEndTime,
			},
		)
	}
	return contractsToCreate
}

func (u *useCase) verifyLoanOfferLines(ctx context.Context, offerLines []entity.LoanPackageOfferInterest, investorId string) (loanOffer entity.LoanPackageOffer, existedMoLoanPackage bool, err error) {
	offerId := offerLines[0].LoanPackageOfferId
	offer, err := u.loanOfferRepository.FindByIdWithRequest(ctx, offerId)
	if err != nil {
		return offer, false, err
	}
	if offer.LoanPackageRequest.InvestorId != investorId {
		return offer, false, apperrors.ErrInvalidInvestorId
	}
	if offer.IsExpired() {
		return offer, false, apperrors.ErrOfferExpired
	}
	for _, offerLine := range offerLines {
		if offerLine.Status == entity.LoanPackageOfferInterestStatusCreatingLoanPackage {
			return offer, false, apperrors.ErrorLoanPackageOfferInterestIsCreating
		}
		if !slices.Contains(
			offerLine.Status.NextStatuses(offer.FlowType),
			entity.LoanPackageOfferInterestStatusLoanPackageCreated,
		) {
			return offer, false, apperrors.ErrInvalidLoanPackageOfferInterestStatus
		}
		if offerLine.LoanPackageOfferId != offerId {
			return offer, false, apperrors.ErrorOfferInterestMismatch
		}
	}

	if offerLines[0].LoanID == 0 && offerLines[0].SubmissionSheetDetailId != 0 {
		return offer, false, nil
	}
	return offer, true, nil
}

func (u *useCase) assignMultipleLoanPackages(
	ctx context.Context,
	accountNo string,
	offerLines []entity.LoanPackageOfferInterest,
	assetType entity.AssetType,
) ([]AssignedLoanPackageAccount, error) {
	assignedLoanPackageAccounts := make([]AssignedLoanPackageAccount, 0, len(offerLines))
	var (
		eg errgroup.Group
		mu sync.Mutex
	)
	for _, offerLine := range offerLines {
		offerLine := offerLine
		eg.Go(
			func() error {
				loanPackageAccountId, err := u.financialProductRepository.AssignLoanPackageOrGetLoanPackageAccountId(
					ctx, accountNo, offerLine.LoanID, assetType,
				)
				if err != nil {
					return err
				}
				mu.Lock()
				assignedLoanPackageAccounts = append(
					assignedLoanPackageAccounts, AssignedLoanPackageAccount{
						LoanPackageOfferInterestId:  offerLine.Id,
						LoanPackageId:               offerLine.LoanID,
						CreatedLoanPackageAccountId: loanPackageAccountId,
						InterestRate:                offerLine.InterestRate,
						LoanRate:                    offerLine.LoanRate,
					},
				)
				mu.Unlock()
				return nil
			},
		)
	}
	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("loanOfferInterestUseCase assignMultipleLoanPackages %w", err)
	}
	return assignedLoanPackageAccounts, nil
}

func (u *useCase) GetAll(ctx context.Context, filter entity.OfferInterestFilter) ([]entity.LoanPackageOfferInterest, core.PagingMetaData, error) {
	var (
		eg             errgroup.Group
		pagingMetaData = core.PagingMetaData{PageSize: filter.Size, PageNumber: filter.Number}
		res            []entity.LoanPackageOfferInterest
	)
	eg.Go(
		func() error {
			foundEntities, err := u.repository.GetWithFilter(ctx, filter)
			res = foundEntities
			return err
		},
	)
	eg.Go(
		func() error {
			count, err := u.repository.CountWithFilter(ctx, filter)
			pagingMetaData.Total = count
			pagingMetaData.TotalPages = filter.TotalPages(count)
			return err
		},
	)
	if err := eg.Wait(); err != nil {
		return res, pagingMetaData, fmt.Errorf("loanPackageRequestUseCase GetAll %w", err)
	}
	return res, pagingMetaData, nil
}

// SyncLoanPackageData syncs loan package data to loan offer interest (currently only support underlying asset type)
func (u *useCase) SyncLoanPackageData(ctx context.Context) (int, error) {
	lines, err := u.repository.GetRequestBasedLoanOfferInterests(ctx)
	if err != nil {
		return 0, fmt.Errorf("loanOfferInterestUseCase MigrateWithWithLoanPackageData %w", err)
	}
	succeededLineCount := 0
	for _, line := range lines {
		loanPackage, err := u.financialProductRepository.GetLoanPackageDetail(ctx, line.LoanID)
		if err != nil {
			_ = u.errorService.NotifyError(ctx, err)
			continue
		}
		line.FeeRate = loanPackage.BuyingFeeRate
		line.InterestRate = loanPackage.InterestRate
		line.LoanRate = decimal.NewFromInt(1).Sub(loanPackage.InitialRate)
		line.Term = loanPackage.Term
		if _, err := u.repository.Update(ctx, line); err != nil {
			_ = u.errorService.NotifyError(ctx, err)
			continue
		}
		succeededLineCount++
	}
	return succeededLineCount, nil
}

func NewUseCase(
	repository repository.LoanPackageOfferInterestRepository,
	atomicExecutor atomicity.AtomicExecutor,
	loanOfferRepository loanPackageOfferRepo.LoanPackageOfferRepository,
	loanContractRepository loanContractRepo.LoanContractPersistenceRepository,
	financialProductRepository financialProductRepo.FinancialProductRepository,
	eventRepository repository.LoanPackageOfferInterestEventRepository,
	errorService apperrors.Service,
	symbolRepository symbolRepo.SymbolRepository,
	submissionSheetRepository submissionSheetRepo.SubmissionSheetRepository,
	policyTemplateRepository loanPolicyRepo.LoanPolicyTemplateRepository,
	appConfig config.AppConfig,
) UseCase {
	return &useCase{
		atomicExecutor:             atomicExecutor,
		repository:                 repository,
		loanOfferRepository:        loanOfferRepository,
		loanContractRepository:     loanContractRepository,
		financialProductRepository: financialProductRepository,
		eventRepository:            eventRepository,
		errorService:               errorService,
		symbolRepository:           symbolRepository,
		submissionSheetRepository:  submissionSheetRepository,
		policyTemplateRepository:   policyTemplateRepository,
		appConfig:                  appConfig,
	}
}
