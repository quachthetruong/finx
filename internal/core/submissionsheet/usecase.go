package submissionsheet

import (
	"context"
	"financing-offer/internal/config"
	financingRepo "financing-offer/internal/core/financing/repository"
	loanPackageOfferRepo "financing-offer/internal/core/loanoffer/repository"
	loanOfferInterestRepo "financing-offer/internal/core/loanofferinterest/repository"
	loanPackageRequestRepo "financing-offer/internal/core/loanpackagerequest/repository"
	symbolRepo "financing-offer/internal/core/symbol/repository"
	"fmt"
	"github.com/shopspring/decimal"
	"time"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/atomicity"
	"financing-offer/internal/core/entity"
	financingProductRepo "financing-offer/internal/core/financialproduct/repository"
	loanPolicyTemplateRepo "financing-offer/internal/core/loanpolicytemplate/repository"
	marginOperationRepo "financing-offer/internal/core/marginoperation/repository"
	"financing-offer/internal/core/submissionsheet/repository"
	"financing-offer/internal/funcs"
)

type UseCase interface {
	Upsert(ctx context.Context, submissionSheetRequest entity.SubmissionSheetShorten) (entity.SubmissionSheet, error)
	Update(ctx context.Context, submissionSheetRequest entity.SubmissionSheetShorten, loanRate entity.LoanRate, loanPolicyTemplates []entity.AggregateLoanPolicyTemplate) (entity.SubmissionSheet, error)
	Create(ctx context.Context, submissionSheetRequest entity.SubmissionSheetShorten, loanRate entity.LoanRate, loanPolicyTemplates []entity.AggregateLoanPolicyTemplate) (entity.SubmissionSheet, error)
	GetLatestByRequestId(ctx context.Context, id int64) (entity.SubmissionSheet, error)
	AdminApproveSubmission(ctx context.Context, submissionId int64) error
	AdminRejectSubmission(ctx context.Context, submissionId int64) error
}

type submissionSheetUseCase struct {
	repository                         repository.SubmissionSheetRepository
	atomicExecutor                     atomicity.AtomicExecutor
	loanPolicyRepository               loanPolicyTemplateRepo.LoanPolicyTemplateRepository
	marginOperationRepository          marginOperationRepo.MarginOperationRepository
	financialProductRepository         financingProductRepo.FinancialProductRepository
	loanPackageRequestRepository       loanPackageRequestRepo.LoanPackageRequestRepository
	loanPackageOfferRepository         loanPackageOfferRepo.LoanPackageOfferRepository
	loanPackageOfferInterestRepository loanOfferInterestRepo.LoanPackageOfferInterestRepository
	financingRepository                financingRepo.FinancingRepository
	appConfig                          config.AppConfig
	errorService                       apperrors.Service
	loanPackageRequestEventRepository  loanPackageRequestRepo.LoanPackageRequestEventRepository
	symbolRepository                   symbolRepo.SymbolRepository
}

// Create new submission sheet if submission sheet is not existed or the latest submission sheet is rejected by odoo
// Update submission sheet in other cases
func (u *submissionSheetUseCase) Upsert(ctx context.Context, submissionSheetRequest entity.SubmissionSheetShorten) (entity.SubmissionSheet, error) {
	errorTemplate := "submissionSheetUseCase Upsert %w"
	loanRate, err := u.financialProductRepository.GetLoanRateDetail(ctx, submissionSheetRequest.Detail.LoanPackageRateId)
	if err != nil {
		return entity.SubmissionSheet{}, fmt.Errorf(errorTemplate, err)
	}
	loanPolicyTemplates, err := u.GetTemplates(ctx, submissionSheetRequest.Detail.LoanPolicies)
	if err != nil {
		return entity.SubmissionSheet{}, fmt.Errorf(errorTemplate, err)
	}
	existedSubmissions, err := u.repository.GetMetadataByRequestId(ctx, submissionSheetRequest.Metadata.LoanPackageRequestId)
	if err != nil {
		return entity.SubmissionSheet{}, fmt.Errorf(errorTemplate, err)
	}
	if len(existedSubmissions) == 0 {
		res, createErr := u.Create(ctx, submissionSheetRequest, loanRate, loanPolicyTemplates)
		if createErr != nil {
			return entity.SubmissionSheet{}, fmt.Errorf(errorTemplate, createErr)
		}
		return res, nil
	}

	latestSubmission := existedSubmissions[0]
	if latestSubmission.Status == entity.SubmissionSheetStatusRejected {
		res, createErr := u.Create(ctx, submissionSheetRequest, loanRate, loanPolicyTemplates)
		if createErr != nil {
			return entity.SubmissionSheet{}, fmt.Errorf(errorTemplate, createErr)
		}
		return res, nil
	}

	if latestSubmission.Id != submissionSheetRequest.Metadata.Id {
		return entity.SubmissionSheet{}, apperrors.ErrorSubmissionIsNotTheLatest
	}
	res, updateErr := u.Update(ctx, submissionSheetRequest, loanRate, loanPolicyTemplates)
	if updateErr != nil {
		return entity.SubmissionSheet{},
			fmt.Errorf(errorTemplate, updateErr)
	}

	return res, nil
}

func (u *submissionSheetUseCase) Update(ctx context.Context, submissionSheetRequest entity.SubmissionSheetShorten, loanRate entity.LoanRate, loanPolicyTemplates []entity.AggregateLoanPolicyTemplate) (entity.SubmissionSheet, error) {
	submissionSheet := entity.SubmissionSheet{}
	txErr := u.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			submissionSheetMetadataReq := entity.SubmissionSheetMetadata{
				Id:                   submissionSheetRequest.Metadata.Id,
				FlowType:             submissionSheetRequest.Metadata.FlowType,
				ActionType:           submissionSheetRequest.Metadata.ActionType,
				ProposeType:          submissionSheetRequest.Metadata.ProposeType,
				LoanPackageRequestId: submissionSheetRequest.Metadata.LoanPackageRequestId,
				Creator:              submissionSheetRequest.Metadata.Creator,
				Status:               submissionSheetRequest.Metadata.Status,
			}
			submissionSheetMetadata, err := u.repository.UpdateMetadata(tc, submissionSheetMetadataReq)
			if err != nil {
				return fmt.Errorf("submissionSheetUseCase Update %w", err)
			}
			loanPolicySnapShots := make([]entity.LoanPolicySnapShot, 0)
			for i, aggregateLoanPolicy := range loanPolicyTemplates {
				loanPolicySnapShots = append(loanPolicySnapShots, aggregateLoanPolicy.ToSnapShotModel(submissionSheetRequest.Detail.LoanPolicies[i]))
			}
			submissionSheetRequestDetailReq := entity.SubmissionSheetDetail{
				SubmissionSheetId: submissionSheetMetadata.Id,
				LoanRate:          loanRate,
				LoanPolicies:      loanPolicySnapShots,
				FirmBuyingFee:     submissionSheetRequest.Detail.FirmBuyingFee,
				FirmSellingFee:    submissionSheetRequest.Detail.FirmSellingFee,
				TransferFee:       submissionSheetRequest.Detail.TransferFee,
				Comment:           submissionSheetRequest.Detail.Comment,
			}
			submissionSheetDetail, err := u.repository.UpdateDetail(tc, submissionSheetRequestDetailReq)
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

func (u *submissionSheetUseCase) Create(ctx context.Context, submissionSheetRequest entity.SubmissionSheetShorten, loanRate entity.LoanRate, loanPolicyTemplates []entity.AggregateLoanPolicyTemplate) (entity.SubmissionSheet, error) {
	submissionSheet := entity.SubmissionSheet{}
	txErr := u.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			submissionSheetMetadataReq := entity.SubmissionSheetMetadata{
				FlowType:             submissionSheetRequest.Metadata.FlowType,
				ActionType:           submissionSheetRequest.Metadata.ActionType,
				ProposeType:          submissionSheetRequest.Metadata.ProposeType,
				LoanPackageRequestId: submissionSheetRequest.Metadata.LoanPackageRequestId,
				Creator:              submissionSheetRequest.Metadata.Creator,
				Status:               submissionSheetRequest.Metadata.Status,
				CreatedAt:            time.Now(),
			}
			submissionSheetMetadata, err := u.repository.CreateMetadata(ctx, submissionSheetMetadataReq)
			if err != nil {
				return fmt.Errorf("submissionSheetUseCase Create %w", err)
			}
			loanPolicySnapShots := make([]entity.LoanPolicySnapShot, 0)
			for i, aggregateLoanPolicy := range loanPolicyTemplates {
				loanPolicySnapShots = append(loanPolicySnapShots, aggregateLoanPolicy.ToSnapShotModel(submissionSheetRequest.Detail.LoanPolicies[i]))
			}
			submissionSheetRequestDetailReq := entity.SubmissionSheetDetail{
				SubmissionSheetId: submissionSheetMetadata.Id,
				LoanRate:          loanRate,
				LoanPolicies:      loanPolicySnapShots,
				FirmBuyingFee:     submissionSheetRequest.Detail.FirmBuyingFee,
				FirmSellingFee:    submissionSheetRequest.Detail.FirmSellingFee,
				TransferFee:       submissionSheetRequest.Detail.TransferFee,
				Comment:           submissionSheetRequest.Detail.Comment,
			}
			submissionSheetDetail, err := u.repository.CreateDetail(ctx, submissionSheetRequestDetailReq)
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

func (u *submissionSheetUseCase) GetTemplates(ctx context.Context, loanPolicies []entity.LoanPolicyShorten) ([]entity.AggregateLoanPolicyTemplate, error) {
	if len(loanPolicies) == 0 {
		return []entity.AggregateLoanPolicyTemplate{}, apperrors.ErrMissingLoanPolicyTemplate
	}
	loanPolicyIds := make([]int64, 0)
	for _, loanPolicy := range loanPolicies {
		loanPolicyIds = append(loanPolicyIds, loanPolicy.LoanPolicyTemplateId)
	}
	loanPolicyTemplates, err := u.loanPolicyRepository.GetByIds(ctx, loanPolicyIds) // loanPolicyTemplates, err := u.loanPolicyRepository.GetByIds(ctx, loanPolicyTemplateIds)
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
		loanPolicyTemplateSnapShots = append(loanPolicyTemplateSnapShots,
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
			})
	}
	return loanPolicyTemplateSnapShots, nil
}

func (u *submissionSheetUseCase) GetMarginPools(ctx context.Context, marginPoolIds []int64) (map[int64]entity.MarginPoolGroup, error) {
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

func (u *submissionSheetUseCase) GetLatestByRequestId(ctx context.Context, id int64) (entity.SubmissionSheet, error) {
	submissionSheet, err := u.repository.GetLatestByRequestId(ctx, id)
	if err != nil {
		return entity.SubmissionSheet{}, fmt.Errorf("submissionSheetUseCase GetByRequestId %w", err)
	}
	return submissionSheet, nil
}

func (u *submissionSheetUseCase) AdminApproveSubmission(ctx context.Context, submissionId int64) error {
	errorTemplate := "submissionSheetUseCase AdminApproveSubmission %w"
	submissionSheet, err := u.repository.GetById(ctx, submissionId)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	if submissionSheet.Metadata.Status != entity.SubmissionSheetStatusSubmitted {
		return fmt.Errorf(errorTemplate, apperrors.ErrorInvalidCurrentSubmissionStatus)
	}
	request, err := u.loanPackageRequestRepository.GetById(ctx, submissionSheet.Metadata.LoanPackageRequestId, entity.LoanPackageFilter{})
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	if request.Status != entity.LoanPackageRequestStatusPending {
		return fmt.Errorf(errorTemplate, apperrors.ErrInvalidRequestStatus)
	}
	offerExpireTime, err := u.financingRepository.GetDateAfter(time.Now(), u.appConfig.LoanRequest.ExpireDays)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	firstPolicy := submissionSheet.Detail.LoanPolicies[0] //all policy share same attributes such as term, interest rate, etc
	var (
		offer               entity.LoanPackageOffer
		acceptOfferInterest entity.LoanPackageOfferInterest
	)
	txErr := u.atomicExecutor.Execute(
		ctx, func(tc context.Context) error {
			err := u.repository.UpdateMetadataStatusById(
				tc, submissionSheet.Metadata.Id, entity.SubmissionSheetStatusApproved,
			)
			if err != nil {
				return err
			}
			request, err = u.loanPackageRequestRepository.UpdateStatusById(tc, request.Id, entity.LoanPackageRequestStatusConfirmed)
			if err != nil {
				return err
			}
			offer, err = u.loanPackageOfferRepository.Create(
				ctx, entity.LoanPackageOffer{
					LoanPackageRequestId: request.Id,
					OfferedBy:            submissionSheet.Metadata.Creator,
					FlowType:             submissionSheet.Metadata.FlowType,
					ExpiredAt:            offerExpireTime,
				},
			)
			if err != nil {
				return err
			}
			acceptOfferInterest = entity.LoanPackageOfferInterest{
				LoanPackageOfferId:      offer.Id,
				SubmissionSheetDetailId: submissionSheet.Detail.Id,
				LoanID:                  0,
				Status:                  entity.LoanPackageOfferInterestStatusPending,
				AssetType:               request.AssetType,
				LimitAmount:             request.LimitAmount,
				ContractSize:            request.ContractSize,
				InitialRate:             request.InitialRate,
				InterestRate:            firstPolicy.InterestRate,
				LoanRate:                decimal.NewFromInt(1).Sub(submissionSheet.Detail.LoanRate.InitialRate),
				FeeRate:                 submissionSheet.Detail.FirmBuyingFee,
				Term:                    int(firstPolicy.Term),
			}
			offerInterests := []entity.LoanPackageOfferInterest{acceptOfferInterest}
			if submissionSheet.Metadata.ActionType == entity.RejectAndSendOtherProposal {
				cancelOfferInterest := entity.LoanPackageOfferInterest{
					LoanPackageOfferId: offer.Id,
					LimitAmount:        request.LimitAmount,
					LoanRate:           request.LoanRate,
					InterestRate:       decimal.Zero,
					Status:             entity.LoanPackageOfferInterestStatusCancelled,
					CancelledBy:        submissionSheet.Metadata.Creator,
					CancelledAt:        time.Now(),
					FeeRate:            decimal.Zero,
					CancelledReason:    entity.LoanPackageOfferCancelledReasonAlternativeOption,
					AssetType:          request.AssetType,
				}
				offerInterests = append(offerInterests, cancelOfferInterest)
			}

			_, err = u.loanPackageOfferInterestRepository.BulkCreate(ctx, offerInterests)
			if err != nil {
				return err
			}
			return nil
		})

	if txErr != nil {
		return fmt.Errorf(errorTemplate, txErr)
	}
	u.errorService.Go(
		ctx, func() error {
			return u.notifyRequestOnlineConfirmation(
				atomicity.WithIgnoreTx(ctx), request, acceptOfferInterest.Id, offer.Id,
			)
		},
	)

	return nil
}

func (u *submissionSheetUseCase) AdminRejectSubmission(ctx context.Context, submissionId int64) error {
	errorTemplate := "submissionSheetUseCase AdminRejectSubmission %w"
	submissionSheet, err := u.repository.GetById(ctx, submissionId)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	if submissionSheet.Metadata.Status != entity.SubmissionSheetStatusSubmitted {
		return fmt.Errorf(errorTemplate, apperrors.ErrorInvalidCurrentSubmissionStatus)
	}
	err = u.repository.UpdateMetadataStatusById(
		ctx, submissionId, entity.SubmissionSheetStatusRejected,
	)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	return nil
}

func (u *submissionSheetUseCase) notifyRequestOnlineConfirmation(
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

func NewUseCase(
	repository repository.SubmissionSheetRepository,
	atomicExecutor atomicity.AtomicExecutor,
	loanPolicyRepository loanPolicyTemplateRepo.LoanPolicyTemplateRepository,
	marginOperationRepository marginOperationRepo.MarginOperationRepository,
	financialProductRepository financingProductRepo.FinancialProductRepository,
	loanPackageRequestRepository loanPackageRequestRepo.LoanPackageRequestRepository,
	loanPackageOfferRepository loanPackageOfferRepo.LoanPackageOfferRepository,
	loanOfferInterestRepository loanOfferInterestRepo.LoanPackageOfferInterestRepository,
	financingRepository financingRepo.FinancingRepository,
	appConfig config.AppConfig,
	errorService apperrors.Service,
	loanPackageRequestEventRepository loanPackageRequestRepo.LoanPackageRequestEventRepository,
	symbolRepository symbolRepo.SymbolRepository,
) UseCase {
	return &submissionSheetUseCase{
		repository:                         repository,
		atomicExecutor:                     atomicExecutor,
		loanPolicyRepository:               loanPolicyRepository,
		marginOperationRepository:          marginOperationRepository,
		financialProductRepository:         financialProductRepository,
		loanPackageRequestRepository:       loanPackageRequestRepository,
		loanPackageOfferRepository:         loanPackageOfferRepository,
		loanPackageOfferInterestRepository: loanOfferInterestRepository,
		financingRepository:                financingRepository,
		appConfig:                          appConfig,
		errorService:                       errorService,
		loanPackageRequestEventRepository:  loanPackageRequestEventRepository,
		symbolRepository:                   symbolRepository,
	}
}
