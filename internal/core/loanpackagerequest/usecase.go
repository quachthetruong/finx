package loanpackagerequest

import (
	"context"
	"encoding/json"
	configRepo "financing-offer/internal/config/repository"
	"financing-offer/internal/core"
	odooServiceRepo "financing-offer/internal/core/odoo_service/repository"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/atomicity"
	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	financialProductRepo "financing-offer/internal/core/financialproduct/repository"
	financingRepo "financing-offer/internal/core/financing/repository"
	investorRepo "financing-offer/internal/core/investor/repository"
	loanContractRepo "financing-offer/internal/core/loancontract/repository"
	loanPackageOfferRepo "financing-offer/internal/core/loanoffer/repository"
	loanPackageOfferInterestRepo "financing-offer/internal/core/loanofferinterest/repository"
	"financing-offer/internal/core/loanpackagerequest/repository"
	loanPolicyTemplateRepo "financing-offer/internal/core/loanpolicytemplate/repository"
	marginOperationRepo "financing-offer/internal/core/marginoperation/repository"
	schedulerRepo "financing-offer/internal/core/scheduler/repository"
	scoreGroupInterestRepo "financing-offer/internal/core/scoregroupinterest/repository"
	submissionSheetRepo "financing-offer/internal/core/submissionsheet/repository"
	symbolRepo "financing-offer/internal/core/symbol/repository"
	"financing-offer/internal/funcs"
	"financing-offer/pkg/optional"
	"financing-offer/pkg/querymod"
)

type UseCase interface {
	GetAll(ctx context.Context, filter entity.LoanPackageFilter) ([]entity.LoanPackageRequest, core.PagingMetaData, error)
	GetAllUnderlyingRequests(ctx context.Context, filter entity.UnderlyingLoanPackageFilter) ([]entity.UnderlyingLoanPackageRequest, error)
	InvestorGetAll(ctx context.Context, filter entity.LoanPackageFilter) ([]entity.LoanPackageRequest, error)
	GetById(ctx context.Context, id int64, filter entity.LoanPackageFilter) (entity.LoanPackageRequest, error)
	InvestorRequest(ctx context.Context, loanPackageRequest entity.LoanPackageRequest, investor entity.Investor) (entity.LoanPackageRequest, error)
	InvestorRequestDerivative(ctx context.Context, loanPackageRequest entity.LoanPackageRequest) (entity.LoanPackageRequest, error)
	Update(ctx context.Context, loanPackageRequest entity.LoanPackageRequest) (entity.LoanPackageRequest, error)
	Delete(ctx context.Context, id int64) error
	AdminConfirmLoanRequest(ctx context.Context, id int64, creator string, loanId int64) (entity.LoanPackageRequest, error)
	AdminCancelLoanRequest(ctx context.Context, id int64, creator string, loanIds []int64) (entity.LoanPackageRequest, error)
	AdminSubmitSubmission(ctx context.Context, submissionSheetRequest entity.SubmissionSheetShorten) (entity.LoanPackageRequest, error)
	SaveExistedLoanRateRequest(ctx context.Context, investorId string, loanPackageRequest entity.LoanPackageRequest) (entity.LoggedRequest, error)
	SystemDeclineRiskLoanRequests(ctx context.Context, maximumLoanRate decimal.Decimal) error
	CancelAllLoanPackageRequestBySymbolId(ctx context.Context, symbolId int64, creator string) ([]entity.LoanPackageRequest, error)
}

type loanPackageRequestUseCase struct {
	repository                         repository.LoanPackageRequestRepository
	atomicExecutor                     atomicity.AtomicExecutor
	scoreGroupInterestRepository       scoreGroupInterestRepo.ScoreGroupInterestRepository
	loanPackageOfferRepository         loanPackageOfferRepo.LoanPackageOfferRepository
	loanPackageOfferInterestRepository loanPackageOfferInterestRepo.LoanPackageOfferInterestRepository
	loanPackageRequestEventRepository  repository.LoanPackageRequestEventRepository
	symbolRepository                   symbolRepo.SymbolRepository
	contractRepository                 loanContractRepo.LoanContractPersistenceRepository
	financialProductRepository         financialProductRepo.FinancialProductRepository
	appConfig                          config.AppConfig
	logger                             *slog.Logger
	financingRepository                financingRepo.FinancingRepository
	schedulerJobRepository             schedulerRepo.SchedulerJobRepository
	errorService                       apperrors.Service
	investorRepository                 investorRepo.InvestorPersistenceRepository
	loanPolicyRepository               loanPolicyTemplateRepo.LoanPolicyTemplateRepository
	submissionSheetRepository          submissionSheetRepo.SubmissionSheetRepository
	marginOperationRepository          marginOperationRepo.MarginOperationRepository
	configurationPersistenceRepo       configRepo.ConfigurationPersistenceRepository
	odooServiceRepository              odooServiceRepo.OdooServiceRepository
}

func (u *loanPackageRequestUseCase) InvestorGetAll(ctx context.Context, filter entity.LoanPackageFilter) ([]entity.LoanPackageRequest, error) {
	res, err := u.repository.GetAll(ctx, filter)
	if err != nil {
		return res, fmt.Errorf("loanPackageRequestUseCase InvestorGetBySymbol %w", err)
	}
	return res, nil
}

func (u *loanPackageRequestUseCase) GetAll(ctx context.Context, filter entity.LoanPackageFilter) ([]entity.LoanPackageRequest, core.PagingMetaData, error) {
	var (
		loanPackageRequests []entity.LoanPackageRequest
		eg                  errgroup.Group
		pagingMetaData      = core.PagingMetaData{PageSize: filter.Size, PageNumber: filter.Number}
	)
	eg.Go(
		func() error {
			res, scopedErr := u.repository.GetAll(ctx, filter)
			loanPackageRequests = res
			return scopedErr
		},
	)
	eg.Go(
		func() error {
			res, scopedErr := u.repository.Count(ctx, filter)
			pagingMetaData.Total = res
			pagingMetaData.TotalPages = filter.TotalPages(res)
			return scopedErr
		},
	)
	if err := eg.Wait(); err != nil {
		return loanPackageRequests, pagingMetaData, fmt.Errorf("loanPackageRequestUseCase GetAll %w", err)
	}
	return loanPackageRequests, pagingMetaData, nil
}

func (u *loanPackageRequestUseCase) GetAllUnderlyingRequests(ctx context.Context, filter entity.UnderlyingLoanPackageFilter) ([]entity.UnderlyingLoanPackageRequest, error) {
	errorTemplate := "loanPackageRequestUseCase GetAllUnderlyingRequests: %w"
	res, err := u.repository.GetAllUnderlyingRequests(ctx, filter)
	if err != nil {
		return []entity.UnderlyingLoanPackageRequest{}, fmt.Errorf(errorTemplate, err)
	}

	return res, nil
}

func (u *loanPackageRequestUseCase) GetById(ctx context.Context, id int64, filter entity.LoanPackageFilter) (entity.LoanPackageRequest, error) {
	res, err := u.repository.GetById(ctx, id, filter)
	if err != nil {
		return res, fmt.Errorf("loanPackageRequestUseCase GetById %w", err)
	}
	return res, nil
}

func (u *loanPackageRequestUseCase) InvestorRequest(ctx context.Context, loanPackageRequest entity.LoanPackageRequest, investor entity.Investor) (entity.LoanPackageRequest, error) {
	errorTemplate := "loanPackageRequestUseCase InvestorRequest %w"
	if err := u.verifyAccountNumber(ctx, loanPackageRequest.InvestorId, loanPackageRequest.AccountNo); err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf(errorTemplate, err)
	}
	if loanPackageRequest.Type == entity.LoanPackageRequestTypeGuaranteed {
		if loanPackageRequest.GuaranteedDuration > u.appConfig.LoanRequest.MaxGuaranteedDuration {
			return entity.LoanPackageRequest{}, apperrors.ErrInvalidGuaranteedDuration
		}
	}
	if err := u.investorRepository.CreateIfNotExist(ctx, investor); err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf(errorTemplate, err)
	}
	res, err := u.repository.Create(ctx, loanPackageRequest)
	if err != nil {
		return res, fmt.Errorf(errorTemplate, err)
	}
	return res, nil
}

func (u *loanPackageRequestUseCase) InvestorRequestDerivative(ctx context.Context, loanPackageRequest entity.LoanPackageRequest) (entity.LoanPackageRequest, error) {
	errorTemplate := "loanPackageRequestUseCase InvestorRequestDerivative %w"
	if err := u.verifyAccountNumber(ctx, loanPackageRequest.InvestorId, loanPackageRequest.AccountNo); err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf(errorTemplate, err)
	}
	symbol, err := u.symbolRepository.GetById(ctx, loanPackageRequest.SymbolId)
	if err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf(errorTemplate, err)
	}
	if symbol.AssetType != entity.AssetTypeDerivative {
		return entity.LoanPackageRequest{}, fmt.Errorf(
			errorTemplate, apperrors.ErrMismatchAssetType,
		)
	}
	res, err := u.repository.Create(ctx, loanPackageRequest)
	if err != nil {
		return res, fmt.Errorf(errorTemplate, err)
	}
	return res, nil
}

func (u *loanPackageRequestUseCase) verifyAccountNumber(ctx context.Context, investorId, accountNo string) error {
	accounts, err := u.financialProductRepository.GetAllAccountDetail(ctx, investorId)
	if err != nil {
		return err
	}
	accountNos := funcs.Map(
		accounts, func(account entity.FinancialAccountDetail) string {
			return account.AccountNo
		},
	)
	if !slices.Contains(accountNos, accountNo) {
		return apperrors.ErrAccountNoInvalid
	}
	return nil
}

func (u *loanPackageRequestUseCase) SaveExistedLoanRateRequest(ctx context.Context, investorId string, loanPackageRequest entity.LoanPackageRequest) (entity.LoggedRequest, error) {
	encodedRequest, err := json.Marshal(loanPackageRequest)
	if err != nil {
		return entity.LoggedRequest{}, fmt.Errorf("loanPackageRequestUseCase SaveExistedLoanRateRequest %w", err)
	}
	res, err := u.repository.SaveLoggedRequest(
		ctx, entity.LoggedRequest{
			InvestorId: investorId,
			SymbolId:   loanPackageRequest.SymbolId,
			Reason:     entity.LoggedRequestReasonLoanRateExisted,
			Request:    string(encodedRequest),
		},
	)
	if err != nil {
		return res, fmt.Errorf("loanPackageRequestUseCase SaveExistedLoanRateRequest %w", err)
	}
	return res, nil
}

func (u *loanPackageRequestUseCase) AdminConfirmLoanRequest(ctx context.Context, id int64, creator string, loanId int64) (entity.LoanPackageRequest, error) {
	var (
		offer         entity.LoanPackageOffer
		offerInterest entity.LoanPackageOfferInterest
		flowType      = entity.FLowTypeDnseOffline
		moLoanPackage = optional.None[entity.FinancialProductLoanPackage]()
	)
	errorTemplate := "loanPackageRequestUseCase adminConfirmLoanRequest %w"
	if loanId > 0 {
		flowType = entity.FlowTypeDnseOnline
		financingProductLoanPackage, err := u.financialProductRepository.GetLoanPackageDetail(ctx, loanId)
		if err != nil {
			return entity.LoanPackageRequest{}, fmt.Errorf(errorTemplate, err)
		}
		moLoanPackage = optional.Some(financingProductLoanPackage)
	}
	request, err := u.getAndVerifyRequestForConfirmation(ctx, id, flowType)
	if err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf(errorTemplate, err)
	}

	offerExpireTime, err := u.financingRepository.GetDateAfter(time.Now(), u.appConfig.LoanRequest.ExpireDays)
	if err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf(errorTemplate, err)
	}

	txErr := u.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			request, err = u.repository.UpdateStatusById(tc, request.Id, entity.LoanPackageRequestStatusConfirmed)
			if err != nil {
				return err
			}

			offer, err = u.loanPackageOfferRepository.Create(
				ctx, entity.LoanPackageOffer{
					LoanPackageRequestId: request.Id,
					OfferedBy:            creator,
					FlowType:             flowType,
					ExpiredAt:            offerExpireTime,
				},
			)
			if err != nil {
				return err
			}
			offerInterestToCreate := entity.LoanPackageOfferInterest{
				LoanPackageOfferId: offer.Id,
				LoanID:             0,
				LoanRate:           request.LoanRate,
				Status:             entity.LoanPackageOfferInterestStatusPending,
				Term:               0,
				FeeRate:            decimal.Zero,
				AssetType:          request.AssetType,
				// asset type underlying request
				LimitAmount:  request.LimitAmount,
				InterestRate: decimal.Zero,
				// asset type derivative request
				ContractSize: request.ContractSize,
				InitialRate:  request.InitialRate,
			}
			if loanPackage := moLoanPackage.Get(); moLoanPackage.IsPresent() {
				offerInterestToCreate.LoanRate = decimal.NewFromInt(1).Sub(loanPackage.InitialRate)
				offerInterestToCreate.InterestRate = loanPackage.InterestRate
				offerInterestToCreate.Term = loanPackage.Term
				offerInterestToCreate.FeeRate = loanPackage.BuyingFeeRate
				offerInterestToCreate.LoanID = loanPackage.Id
			}
			offerInterest, err = u.loanPackageOfferInterestRepository.Create(ctx, offerInterestToCreate)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if txErr != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf(errorTemplate, txErr)
	}
	if flowType == entity.FlowTypeDnseOnline {
		u.errorService.Go(
			ctx, func() error {
				return u.notifyRequestOnlineConfirmation(
					atomicity.WithIgnoreTx(ctx), request, offerInterest.Id, offer.Id,
				)
			},
		)
	} else {
		u.errorService.Go(
			ctx, func() error {
				return u.notifyRequestOfflineConfirmation(atomicity.WithIgnoreTx(ctx), request)
			},
		)
	}
	return request, nil
}

func (u *loanPackageRequestUseCase) AdminSubmitSubmission(ctx context.Context, submissionSheetRequest entity.SubmissionSheetShorten) (entity.LoanPackageRequest, error) {
	errorTemplate := "loanPackageRequestUseCase AdminSubmitSubmission %w"
	if err := u.verifyActionFlowProposeType(submissionSheetRequest); err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf(errorTemplate, err)
	}
	loanRate, err := u.getAndVerifyLoanRate(ctx, submissionSheetRequest.Detail.LoanPackageRateId)
	if err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf(errorTemplate, err)
	}
	loanPolicyTemplates, err := u.GetTemplates(ctx, submissionSheetRequest.Detail.LoanPolicies)
	if err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf(errorTemplate, err)
	}
	loanPolicySnapShots := make([]entity.LoanPolicySnapShot, 0)
	for i, aggregateLoanPolicy := range loanPolicyTemplates {
		loanPolicySnapShots = append(loanPolicySnapShots, aggregateLoanPolicy.ToSnapShotModel(submissionSheetRequest.Detail.LoanPolicies[i]))
	}
	request, err := u.getAndVerifyRequestForConfirmation(
		ctx, submissionSheetRequest.Metadata.LoanPackageRequestId, submissionSheetRequest.Metadata.FlowType,
	)
	if err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf(errorTemplate, err)
	}
	symbol, err := u.symbolRepository.GetById(ctx, request.SymbolId)
	if err != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf(errorTemplate, err)
	}
	txErr := u.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			submissionSheet, err := u.UpsertSubmissionSheet(tc, submissionSheetRequest.ToSubmissionSheet(loanRate, loanPolicySnapShots))
			if err != nil {
				return err
			}
			var formattedSource string
			for index, loanPolicy := range submissionSheet.Detail.LoanPolicies {
				formattedLoanRate := loanPolicy.InterestRate.Mul(decimal.NewFromInt(100))
				if index == 0 {
					formattedSource = fmt.Sprintf("%s: %s%%", loanPolicy.Source, formattedLoanRate)
				}
				formattedSource = fmt.Sprintf("%s, %s: %s%%", formattedSource, loanPolicy.Source, formattedLoanRate)
			}
			interestRate := submissionSheet.Detail.LoanPolicies[0].InterestRate // All loan policy share the same interest rate
			term := submissionSheet.Detail.LoanPolicies[0].Term
			formattedLoanRate := decimal.NewFromInt(1).Sub(loanRate.InitialRate).Mul(decimal.NewFromInt(100))
			err = u.odooServiceRepository.SendLoanApprovalRequest(entity.LoanApprovalRequest{
				LoanRequestId:      request.Id,
				CreateAt:           request.CreatedAt,
				SubmissionId:       submissionSheet.Metadata.Id,
				SubmissionBy:       submissionSheet.Metadata.Creator,
				SubmissionCreateAt: submissionSheet.Metadata.CreatedAt,
				InvestorId:         request.InvestorId,
				AccountNo:          request.AccountNo,
				Symbol:             symbol.Symbol,
				LoanRate:           formattedLoanRate,
				Source:             formattedSource,
				InterestRate:       interestRate,
				BuyingFee:          submissionSheet.Detail.FirmBuyingFee,
				Term:               term,
				Description:        submissionSheet.Detail.Comment,
			})
			if err != nil {
				return err
			}
			return nil
		},
	)
	if txErr != nil {
		return entity.LoanPackageRequest{}, fmt.Errorf(errorTemplate, txErr)
	}

	return request, nil
}

func (u *loanPackageRequestUseCase) getAndVerifyRequestForConfirmation(ctx context.Context, id int64, flowType entity.FlowType) (entity.LoanPackageRequest, error) {
	request, err := u.repository.GetById(ctx, id, entity.LoanPackageFilter{})
	if err != nil {
		return request, err
	}
	if request.Status != entity.LoanPackageRequestStatusPending {
		return request, apperrors.ErrInvalidRequestStatus
	}
	if (request.AssetType == entity.AssetTypeDerivative) && flowType == entity.FlowTypeDnseOnline {
		return request, apperrors.ErrInvalidInput("derivative requests must be confirmed offline")
	}
	return request, nil
}

func (u *loanPackageRequestUseCase) getAndVerifyLoanRate(
	ctx context.Context,
	loanRateId int64,
) (entity.LoanRate, error) {
	loanRate, err := u.financialProductRepository.GetLoanRateDetail(ctx, loanRateId)
	if err != nil {
		return entity.LoanRate{}, err
	}
	loanRateConfig, err := u.configurationPersistenceRepo.GetLoanRateConfiguration(ctx)
	if err != nil {
		return entity.LoanRate{}, err
	}
	if !slices.Contains(loanRateConfig.Ids, loanRate.Id) {
		return entity.LoanRate{}, apperrors.ErrInvalidLoanRateId
	}
	return loanRate, nil
}

func (u *loanPackageRequestUseCase) verifyActionFlowProposeType(
	submissionSheetRequest entity.SubmissionSheetShorten,
) error {
	if submissionSheetRequest.Metadata.FlowType != entity.FlowTypeDnseOnline {
		return apperrors.ErrInvalidFlowType(submissionSheetRequest.Metadata.FlowType)
	}
	if submissionSheetRequest.Metadata.ProposeType != entity.NewLoanPackage {
		return apperrors.ErrInvalidProposeType
	}
	return nil
}

func (u *loanPackageRequestUseCase) AdminCancelLoanRequest(ctx context.Context, id int64, creator string, loanIds []int64) (entity.LoanPackageRequest, error) {
	// decline
	if len(loanIds) == 0 {
		res, err := u.AdminDeclineLoanRequestWithNoAlternativeOption(ctx, id, creator)
		if err != nil {
			return res, fmt.Errorf("loanPackageRequestUseCase AdminCancelLoanRequest %w", err)
		}
		u.errorService.Go(
			ctx, func() error {
				return u.notifyRequestDeclined(atomicity.WithIgnoreTx(ctx), res)
			},
		)
		return res, nil
	}
	// decline with alternative options
	res, notifyLoanOfferInterestId, offerId, err := u.AdminDeclineLoanRequestWithAlternativeOptions(
		ctx, id, creator, loanIds,
	)
	if err != nil {
		return res, fmt.Errorf("loanPackageRequestUseCase AdminCancelLoanRequest %w", err)
	}
	u.errorService.Go(
		ctx, func() error {
			return u.notifyRequestOnlineConfirmation(
				atomicity.WithIgnoreTx(ctx), res, notifyLoanOfferInterestId, offerId,
			)
		},
	)
	return res, nil
}

func (u *loanPackageRequestUseCase) AdminDeclineLoanRequestWithNoAlternativeOption(
	ctx context.Context,
	id int64,
	creator string,
) (entity.LoanPackageRequest, error) {
	errorTemplate := "loanPackageRequestUseCase AdminDeclineLoanRequestWithNoAlternativeOption %w"
	res := entity.LoanPackageRequest{}
	txErr := u.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			request, err := u.repository.GetById(tc, id, entity.LoanPackageFilter{}, querymod.WithLock())
			if err != nil {
				return err
			}
			if request.Status != entity.LoanPackageRequestStatusPending {
				return apperrors.ErrInvalidRequestStatus
			}

			_, err = u.loanPackageOfferRepository.Create(
				tc, entity.LoanPackageOffer{
					LoanPackageRequestId: request.Id,
					OfferedBy:            creator,
				},
			)
			if err != nil {
				return err
			}
			request.Status = entity.LoanPackageRequestStatusConfirmed
			_, err = u.repository.Update(tc, request)
			if err != nil {
				return err
			}
			res = request
			return nil
		},
	)
	if txErr != nil {
		return res, fmt.Errorf(errorTemplate, txErr)
	}
	return res, nil
}

func (u *loanPackageRequestUseCase) AdminDeclineLoanRequestWithAlternativeOptions(
	ctx context.Context,
	id int64,
	creator string,
	loanIds []int64,
) (entity.LoanPackageRequest, int64, int64, error) {
	errorTemplate := "loanPackageRequestUseCase AdminDeclineLoanRequestWithAlternativeOptions %w"
	res := entity.LoanPackageRequest{}
	createdLoanPackageOfferInterest := int64(0)
	createdOfferId := int64(0)
	financialProductLoanPackages, err := u.financialProductRepository.GetLoanPackageDetails(ctx, loanIds)
	if err != nil {
		return res, 0, 0, fmt.Errorf("AdminDeclineLoanRequestWithAlternativeOptions %w", err)
	}
	if len(financialProductLoanPackages) != len(loanIds) {
		return res, 0, 0, apperrors.ErrLoanPackageIdsInvalid
	}
	txErr := u.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			request, err := u.prepareAndPersistLoanPackageRequest(tc, id)
			expireTime, err := u.financingRepository.GetDateAfter(time.Now(), u.appConfig.LoanRequest.ExpireDays)
			if err != nil {
				return err
			}
			res = request
			newOffer, err := u.loanPackageOfferRepository.Create(
				tc, entity.LoanPackageOffer{
					LoanPackageRequestId: request.Id,
					OfferedBy:            creator,
					FlowType:             entity.FlowTypeDnseOnline,
					ExpiredAt:            expireTime,
				},
			)
			if err != nil {
				return err
			}
			createLoanOfferInterest := make([]entity.LoanPackageOfferInterest, 0, len(financialProductLoanPackages)+1)
			for _, financialProductLoanPackage := range financialProductLoanPackages {
				createLoanOfferInterest = append(
					createLoanOfferInterest, entity.LoanPackageOfferInterest{
						LoanPackageOfferId: newOffer.Id,
						LimitAmount:        request.LimitAmount,
						LoanRate:           decimal.NewFromInt(1).Sub(financialProductLoanPackage.InitialRate),
						InterestRate:       financialProductLoanPackage.InterestRate,
						Status:             entity.LoanPackageOfferInterestStatusPending,
						LoanID:             financialProductLoanPackage.Id,
						Term:               financialProductLoanPackage.Term,
						FeeRate:            financialProductLoanPackage.BuyingFeeRate,
						AssetType:          request.AssetType,
					},
				)
			}
			createLoanOfferInterest = append(
				createLoanOfferInterest, entity.LoanPackageOfferInterest{
					LoanPackageOfferId: newOffer.Id,
					LimitAmount:        request.LimitAmount,
					LoanRate:           request.LoanRate,
					InterestRate:       decimal.Zero,
					Status:             entity.LoanPackageOfferInterestStatusCancelled,
					CancelledBy:        creator,
					CancelledAt:        time.Now(),
					FeeRate:            decimal.Zero,
					CancelledReason:    entity.LoanPackageOfferCancelledReasonAlternativeOption,
					AssetType:          request.AssetType,
				},
			)
			createdLoanOfferInterests, err := u.loanPackageOfferInterestRepository.BulkCreate(
				tc, createLoanOfferInterest,
			)
			if err != nil {
				return err
			}
			if (len(createdLoanOfferInterests)) == 0 {
				return fmt.Errorf("AdminDeclineLoanRequestWithAlternativeOptions cannot create loan offer interest")
			}
			createdLoanPackageOfferInterest = createdLoanOfferInterests[0].Id
			createdOfferId = newOffer.Id
			return nil
		},
	)
	if txErr != nil {
		return res, 0, 0, fmt.Errorf(errorTemplate, txErr)
	}
	return res, createdLoanPackageOfferInterest, createdOfferId, nil
}

func (u *loanPackageRequestUseCase) prepareAndPersistLoanPackageRequest(atomicContext context.Context, id int64) (entity.LoanPackageRequest, error) {
	request, err := u.repository.GetById(atomicContext, id, entity.LoanPackageFilter{}, querymod.WithLock())
	if err != nil {
		return request, err
	}
	if request.Status != entity.LoanPackageRequestStatusPending {
		return request, apperrors.ErrInvalidRequestStatus
	}
	if request.AssetType == entity.AssetTypeDerivative {
		return request, apperrors.ErrInvalidInput("derivative requests must be confirmed offline")
	}
	request.Status = entity.LoanPackageRequestStatusConfirmed
	_, err = u.repository.Update(atomicContext, request)
	if err != nil {
		return request, err
	}
	return request, nil
}

func (u *loanPackageRequestUseCase) SystemDeclineRiskLoanRequests(ctx context.Context, maximumLoanRate decimal.Decimal) error {
	errorTemplate := "SystemDeclineRiskLoanRequests %w"
	loanRequests := make([]entity.LoanPackageRequest, 0)
	txErr := u.atomicExecutor.Execute(
		ctx, func(ctx context.Context) error {
			pendingRequests, err := u.repository.LockAllPendingRequestByMaxPercent(ctx, maximumLoanRate)
			if err != nil {
				return fmt.Errorf("SystemDeclineRiskLoanRequests cannot lock pending request %w", err)
			}
			if len(pendingRequests) == 0 {
				return nil
			}
			loanRequestIds := funcs.Map(
				pendingRequests, func(r entity.LoanPackageRequest) int64 {
					return r.Id
				},
			)
			if err = u.systemDeclineLoanRequestIds(ctx, loanRequestIds); err != nil {
				return fmt.Errorf(errorTemplate, err)
			}
			loanRequests = append(loanRequests, pendingRequests...)
			return nil
		},
	)
	if err := u.appendJobTrackingData(ctx, loanRequests, txErr); err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	if txErr != nil {
		return fmt.Errorf(errorTemplate, txErr)
	}
	for _, loanRequest := range loanRequests {
		loanRequest := loanRequest
		u.errorService.Go(
			ctx, func() error {
				return u.notifyRequestDeclined(atomicity.WithIgnoreTx(ctx), loanRequest)
			},
		)
	}
	return nil
}

func (u *loanPackageRequestUseCase) systemDeclineLoanRequestIds(ctx context.Context, loanRequestIds []int64) error {
	_, err := u.repository.UpdateStatusByLoanRequestIds(ctx, loanRequestIds, entity.LoanPackageRequestStatusConfirmed)
	if err != nil {
		return fmt.Errorf("systemDeclineLoanRequestIds UpdateStatusByLoanRequestIds %w", err)
	}
	newOffers := make([]entity.LoanPackageOffer, 0, len(loanRequestIds))
	for _, loanRequestId := range loanRequestIds {
		newOffers = append(
			newOffers, entity.LoanPackageOffer{
				LoanPackageRequestId: loanRequestId,
				OfferedBy:            "system",
			},
		)
	}
	_, err = u.loanPackageOfferRepository.BulkCreate(ctx, newOffers)
	if err != nil {
		return fmt.Errorf("systemDeclineLoanRequestIds BulkCreate LoanOffer %w", err)
	}
	return nil
}

func (u *loanPackageRequestUseCase) appendJobTrackingData(ctx context.Context, loanRequests []entity.LoanPackageRequest, err error) error {
	trackingMap := make(map[string]interface{})
	if err != nil {
		trackingMap["error"] = err.Error()
	}
	if err == nil {
		loanRequestIds := funcs.Map(
			loanRequests, func(r entity.LoanPackageRequest) int64 {
				return r.Id
			},
		)
		trackingMap["loanRequestIds"] = loanRequestIds
	}
	trackingData, marshalErr := json.Marshal(trackingMap)
	if marshalErr != nil {
		return fmt.Errorf("appendJobTrackingData cannot marshal tracking data %w", marshalErr)
	}
	jobStatus := entity.JobStatusSuccess
	if err != nil {
		jobStatus = entity.JobStatusFail
	}
	createError := u.schedulerJobRepository.Create(
		ctx, entity.SchedulerJob{
			JobType:      entity.JobTypeDeclineHighRiskLoanRequest,
			JobStatus:    jobStatus,
			TriggerBy:    "system",
			TrackingData: string(trackingData),
		},
	)
	if createError != nil {
		return fmt.Errorf("appendJobTrackingData cannot create tracking data %w", err)
	}
	return nil
}

func (u *loanPackageRequestUseCase) notifyRequestOnlineConfirmation(
	ctx context.Context,
	request entity.LoanPackageRequest,
	offerInterestId int64,
	offerId int64,
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
	return u.loanPackageRequestEventRepository.NotifyOnlineConfirmation(
		ctx, entity.RequestOnlineConfirmationNotify{
			InvestorId:      request.InvestorId,
			RequestName:     fmt.Sprintf("%s-%d", symbol.Symbol, request.Id),
			AccountNo:       request.AccountNo,
			AccountNoDesc:   accountNoDesc,
			OfferInterestId: offerInterestId,
			Symbol:          symbol.Symbol,
			CreatedAt:       time.Now(),
			OfferId:         offerId,
		},
	)
}

func (u *loanPackageRequestUseCase) notifyRequestOfflineConfirmation(ctx context.Context, request entity.LoanPackageRequest) error {
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
		return u.loanPackageRequestEventRepository.NotifyOfflineConfirmation(
			ctx, entity.RequestOfflineConfirmation{
				InvestorId:    request.InvestorId,
				RequestName:   fmt.Sprintf("%s-%d", symbol.Symbol, request.Id),
				AccountNo:     request.AccountNo,
				AccountNoDesc: accountNoDesc,
				Symbol:        symbol.Symbol,
				CreatedAt:     time.Now(),
			},
		)
	case entity.AssetTypeDerivative:
		return u.loanPackageRequestEventRepository.NotifyDerivativeOfflineConfirmation(
			ctx, entity.DerivativeRequestOfflineConfirmation{
				InvestorId:    request.InvestorId,
				RequestName:   fmt.Sprintf("%s-%d", symbol.Symbol, request.Id),
				AccountNo:     request.AccountNo,
				AccountNoDesc: accountNoDesc,
				Symbol:        symbol.Symbol,
				AssetType:     request.AssetType.String(),
				CreatedAt:     time.Now(),
			},
		)
	default:
		return apperrors.ErrMismatchAssetType
	}
}

func (u *loanPackageRequestUseCase) notifyRequestDeclined(ctx context.Context, request entity.LoanPackageRequest) error {
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
		return u.loanPackageRequestEventRepository.NotifyRequestDeclined(
			ctx, entity.LoanPackageRequestDeclinedNotify{
				InvestorId:    request.InvestorId,
				RequestName:   fmt.Sprintf("%s-%d", symbol.Symbol, request.Id),
				AccountNo:     request.AccountNo,
				AccountNoDesc: accountNoDesc,
				Symbol:        symbol.Symbol,
				CreatedAt:     time.Now(),
			},
		)
	case entity.AssetTypeDerivative:
		return u.loanPackageRequestEventRepository.NotifyDerivativeRequestDeclined(
			ctx, entity.LoanPackageDerivativeRequestDeclinedNotify{
				InvestorId:    request.InvestorId,
				RequestName:   fmt.Sprintf("%s-%d", symbol.Symbol, request.Id),
				AccountNo:     request.AccountNo,
				AccountNoDesc: accountNoDesc,
				Symbol:        symbol.Symbol,
				AssetType:     request.AssetType.String(),
				CreatedAt:     time.Now(),
			},
		)
	default:
		return apperrors.ErrMismatchAssetType
	}
}

func (u *loanPackageRequestUseCase) Update(ctx context.Context, loanPackageRequest entity.LoanPackageRequest) (entity.LoanPackageRequest, error) {
	res, err := u.repository.Update(ctx, loanPackageRequest)
	if err != nil {
		return res, fmt.Errorf("loanPackageRequestUseCase Update %w", err)
	}
	return res, nil
}

func (u *loanPackageRequestUseCase) Delete(ctx context.Context, id int64) error {
	err := u.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("loanPackageRequestUseCase Delete %w", err)
	}
	return nil
}

func (u *loanPackageRequestUseCase) CancelAllLoanPackageRequestBySymbolId(ctx context.Context, symbolId int64, creator string) ([]entity.LoanPackageRequest, error) {
	errorTemplate := "CancelAllSymbolLoanPackageRequest %w"
	loanRequests := make([]entity.LoanPackageRequest, 0)
	txErr := u.atomicExecutor.Execute(
		ctx, func(ctx context.Context) error {
			pendingRequests, err := u.repository.LockAndReturnAllPendingRequestBySymbolId(ctx, symbolId)
			if err != nil {
				return fmt.Errorf(errorTemplate, err)
			}
			for _, request := range pendingRequests {
				_, err = u.loanPackageOfferRepository.Create(
					ctx, entity.LoanPackageOffer{
						LoanPackageRequestId: request.Id,
						OfferedBy:            creator,
					},
				)
				if err != nil {
					return fmt.Errorf(errorTemplate, err)
				}
				request.Status = entity.LoanPackageRequestStatusConfirmed
				request, err = u.repository.Update(ctx, request)
				if err != nil {
					return fmt.Errorf(errorTemplate, err)
				}
				loanRequests = append(loanRequests, request)
			}
			return nil
		},
	)
	if txErr != nil {
		return nil, fmt.Errorf(errorTemplate, txErr)
	}
	for _, request := range loanRequests {
		req := request
		u.errorService.Go(
			ctx, func() error {
				return u.notifyRequestDeclined(atomicity.WithIgnoreTx(ctx), req)
			},
		)
	}
	return loanRequests, nil
}

func (u *loanPackageRequestUseCase) UpsertSubmissionSheet(ctx context.Context, submissionSheetRequest entity.SubmissionSheet) (entity.SubmissionSheet, error) {
	errorTemplate := "loanPackageRequestUseCase UpsertSubmissionSheet %w"
	existedSubmissions, err := u.submissionSheetRepository.GetMetadataByRequestId(ctx, submissionSheetRequest.Metadata.LoanPackageRequestId)
	if err != nil {
		return entity.SubmissionSheet{}, fmt.Errorf(errorTemplate, err)
	}
	if len(existedSubmissions) == 0 {
		res, createErr := u.CreateSubmissionSheet(ctx, submissionSheetRequest)
		if createErr != nil {
			return entity.SubmissionSheet{}, fmt.Errorf(errorTemplate, createErr)
		}
		return res, nil
	}

	latestSubmission := existedSubmissions[0]
	if latestSubmission.Status == entity.SubmissionSheetStatusRejected {
		res, createErr := u.CreateSubmissionSheet(ctx, submissionSheetRequest)
		if createErr != nil {
			return entity.SubmissionSheet{}, fmt.Errorf(errorTemplate, createErr)
		}
		return res, nil
	}

	if latestSubmission.Id != submissionSheetRequest.Metadata.Id {
		return entity.SubmissionSheet{}, apperrors.ErrorSubmissionIsNotTheLatest
	}
	res, updateErr := u.UpdateSubmissionSheet(ctx, submissionSheetRequest)
	if updateErr != nil {
		return entity.SubmissionSheet{},
			fmt.Errorf(errorTemplate, updateErr)
	}

	return res, nil
}

func (u *loanPackageRequestUseCase) UpdateSubmissionSheet(ctx context.Context, submissionSheetRequest entity.SubmissionSheet) (entity.SubmissionSheet, error) {
	submissionSheet := entity.SubmissionSheet{}
	txErr := u.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			submissionSheetMetadata, err := u.submissionSheetRepository.UpdateMetadata(tc, submissionSheetRequest.Metadata)
			if err != nil {
				return fmt.Errorf("submissionSheetUseCase Update %w", err)
			}
			submissionSheetRequest := submissionSheetRequest.Detail
			submissionSheetRequest.SubmissionSheetId = submissionSheetMetadata.Id
			submissionSheetDetail, err := u.submissionSheetRepository.UpdateDetail(tc, submissionSheetRequest)
			if err != nil {
				return fmt.Errorf("submissionSheetUseCase Update %w", err)
			}
			submissionSheet.Metadata = submissionSheetMetadata
			submissionSheet.Detail = submissionSheetDetail
			return nil
		},
	)
	if txErr != nil {
		return submissionSheet, fmt.Errorf("submissionSheetUseCase Update %w", txErr)
	}
	return submissionSheet, nil
}

func (u *loanPackageRequestUseCase) CreateSubmissionSheet(ctx context.Context, submissionSheetRequest entity.SubmissionSheet) (entity.SubmissionSheet, error) {
	submissionSheet := entity.SubmissionSheet{}
	txErr := u.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			submissionSheetMetadata, err := u.submissionSheetRepository.CreateMetadata(ctx, submissionSheetRequest.Metadata)
			if err != nil {
				return err
			}
			detailRequest := submissionSheetRequest.Detail
			detailRequest.SubmissionSheetId = submissionSheetMetadata.Id
			submissionSheetDetail, err := u.submissionSheetRepository.CreateDetail(ctx, detailRequest)
			if err != nil {
				return fmt.Errorf("submissionSheetUseCase Create %w", err)
			}
			submissionSheet.Metadata = submissionSheetMetadata
			submissionSheet.Detail = submissionSheetDetail
			return nil
		},
	)
	if txErr != nil {
		return submissionSheet, fmt.Errorf("submissionSheetUseCase Create %w", txErr)
	}
	return submissionSheet, nil
}

func (u *loanPackageRequestUseCase) GetTemplates(ctx context.Context, loanPolicies []entity.LoanPolicyShorten) ([]entity.AggregateLoanPolicyTemplate, error) {
	if len(loanPolicies) == 0 {
		return []entity.AggregateLoanPolicyTemplate{}, apperrors.ErrMissingLoanPolicyTemplate
	}
	loanPolicyIds := make([]int64, 0)
	for _, loanPolicy := range loanPolicies {
		loanPolicyIds = append(loanPolicyIds, loanPolicy.LoanPolicyTemplateId)
	}
	loanPolicyTemplates, err := u.loanPolicyRepository.GetByIds(
		ctx, loanPolicyIds,
	) // loanPolicyTemplates, err := u.loanPolicyRepository.GetByIds(ctx, loanPolicyTemplateIds)
	if err != nil || len(loanPolicyIds) != len(loanPolicyTemplates) {
		return []entity.AggregateLoanPolicyTemplate{}, apperrors.ErrLoanPolicyTemplateIdsInvalid
	}

	poolIds := make([]int64, 0)
	for _, policy := range loanPolicyTemplates {
		poolIds = append(poolIds, policy.PoolIdRef)
	}
	marginPools, err := u.GetMarginPools(ctx, poolIds)
	if err != nil {
		return []entity.AggregateLoanPolicyTemplate{}, fmt.Errorf("loanPolicyTemplateUseCase GetById %w", err)
	}
	if len(marginPools) == 0 {
		return []entity.AggregateLoanPolicyTemplate{}, fmt.Errorf("loanPolicyTemplateUseCase GetById: margin pool not found")
	}
	loanPolicyTemplateSnapShots := make([]entity.AggregateLoanPolicyTemplate, 0)
	for _, loanPolicy := range loanPolicyTemplates {
		loanPolicyTemplateSnapShots = append(
			loanPolicyTemplateSnapShots,
			entity.AggregateLoanPolicyTemplate{
				Id:                       loanPolicy.Id,
				Name:                     loanPolicy.Name,
				PoolIdRef:                loanPolicy.PoolIdRef,
				Source:                   marginPools[loanPolicy.PoolIdRef].Source,
				CreatedAt:                loanPolicy.CreatedAt,
				UpdatedAt:                loanPolicy.UpdatedAt,
				UpdatedBy:                loanPolicy.UpdatedBy,
				InterestRate:             loanPolicy.InterestRate,
				InterestBasis:            loanPolicy.InterestBasis,
				Term:                     loanPolicy.Term,
				OverdueInterest:          loanPolicy.OverdueInterest,
				AllowExtendLoanTerm:      loanPolicy.AllowExtendLoanTerm,
				AllowEarlyPayment:        loanPolicy.AllowEarlyPayment,
				PreferentialPeriod:       loanPolicy.PreferentialPeriod,
				PreferentialInterestRate: loanPolicy.PreferentialInterestRate,
			},
		)
	}
	return loanPolicyTemplateSnapShots, nil
}

func (u *loanPackageRequestUseCase) GetMarginPools(ctx context.Context, marginPoolIds []int64) (map[int64]entity.MarginPoolGroup, error) {
	marginPools, err := u.marginOperationRepository.GetMarginPoolsByIds(ctx, marginPoolIds)
	if err != nil {
		return nil, fmt.Errorf("GetMarginPools %w", err)
	}
	poolGroupIds := funcs.Map(marginPools, func(m entity.MarginPool) int64 { return m.PoolGroupId })
	marginPoolGroups, err := u.marginOperationRepository.GetMarginPoolGroupsByIds(ctx, poolGroupIds)
	if err != nil {
		return nil, fmt.Errorf("GetMarginPools %w", err)
	}
	marginPoolGroupMapped := funcs.AssociateBy[entity.MarginPoolGroup, int64](
		marginPoolGroups, func(group entity.MarginPoolGroup) int64 {
			return group.Id
		},
	)
	marginPoolMapped := make(map[int64]entity.MarginPoolGroup)
	for _, pool := range marginPools {
		group, ok := marginPoolGroupMapped[pool.PoolGroupId]
		if !ok {
			return nil, fmt.Errorf("GetMarginPools: margin pool group not found %d", pool.PoolGroupId)
		}
		marginPoolMapped[pool.Id] = group
	}
	return marginPoolMapped, nil
}

func NewUseCase(
	loanPackageRequestRepo repository.LoanPackageRequestRepository,
	atomicExecutor atomicity.AtomicExecutor,
	scoreGroupInterestRepository scoreGroupInterestRepo.ScoreGroupInterestRepository,
	loanPackageOfferRepository loanPackageOfferRepo.LoanPackageOfferRepository,
	loanPackageOfferInterestRepository loanPackageOfferInterestRepo.LoanPackageOfferInterestRepository,
	eventRepository repository.LoanPackageRequestEventRepository,
	symbolRepository symbolRepo.SymbolRepository,
	contractRepository loanContractRepo.LoanContractPersistenceRepository,
	financialProductRepository financialProductRepo.FinancialProductRepository,
	appConfig config.AppConfig,
	loanPolicyRepository loanPolicyTemplateRepo.LoanPolicyTemplateRepository,
	logger *slog.Logger,
	financingRepository financingRepo.FinancingRepository,
	schedulerJobRepository schedulerRepo.SchedulerJobRepository,
	errorService apperrors.Service,
	investorRepository investorRepo.InvestorPersistenceRepository,
	submissionSheetRepository submissionSheetRepo.SubmissionSheetRepository,
	marginOperationRepository marginOperationRepo.MarginOperationRepository,
	configurationPersistenceRepo configRepo.ConfigurationPersistenceRepository,
	odooServiceRepository odooServiceRepo.OdooServiceRepository,
) UseCase {
	return &loanPackageRequestUseCase{
		repository:                         loanPackageRequestRepo,
		atomicExecutor:                     atomicExecutor,
		scoreGroupInterestRepository:       scoreGroupInterestRepository,
		loanPackageOfferRepository:         loanPackageOfferRepository,
		loanPackageOfferInterestRepository: loanPackageOfferInterestRepository,
		loanPackageRequestEventRepository:  eventRepository,
		symbolRepository:                   symbolRepository,
		contractRepository:                 contractRepository,
		financialProductRepository:         financialProductRepository,
		appConfig:                          appConfig,
		loanPolicyRepository:               loanPolicyRepository,
		logger:                             logger,
		financingRepository:                financingRepository,
		schedulerJobRepository:             schedulerJobRepository,
		errorService:                       errorService,
		investorRepository:                 investorRepository,
		submissionSheetRepository:          submissionSheetRepository,
		marginOperationRepository:          marginOperationRepository,
		configurationPersistenceRepo:       configurationPersistenceRepo,
		odooServiceRepository:              odooServiceRepository,
	}
}
