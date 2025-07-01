package submissionsheet

import (
	"context"
	"financing-offer/internal/apperrors"
	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"financing-offer/test/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestLoanPackageRequestUseCase_AdminApproveSubmission(t *testing.T) {
	t.Parallel()
	appConfig := config.AppConfig{
		LoanRequest: config.LoanRequestConfig{
			ExpireDays: 7,
		},
	}
	loanPackageRequestRepo := mock.NewMockLoanPackageRequestRepository(t)
	loanPackageOfferRepository := mock.NewMockLoanPackageOfferRepository(t)
	loanPackageOfferInterestRepository := mock.NewMockLoanPackageOfferInterestRepository(t)
	loanPackageRequestEventRepository := mock.NewMockLoanPackageRequestEventRepository(t)
	symbolRepo := mock.NewMockSymbolRepository(t)
	financialProductRepo := mock.NewMockFinancialProductRepository(t)
	financingRepo := mock.NewMockFinancingRepository(t)
	loanPolicyTemplateRepo := mock.NewMockLoanPolicyTemplateRepository(t)
	submissionSheetRepo := mock.NewMockSubmissionSheetRepository(t)
	marginOperationRepo := mock.NewMockMarginOperationRepository(t)
	atomicExecutor := mock.NewMockAtomicExecutorExecutePassthrough(t)
	errorService := mock.ErrReporter{}
	useCase := NewUseCase(submissionSheetRepo, atomicExecutor, loanPolicyTemplateRepo, marginOperationRepo, financialProductRepo, loanPackageRequestRepo, loanPackageOfferRepository, loanPackageOfferInterestRepository, financingRepo, appConfig, errorService, loanPackageRequestEventRepository, symbolRepo)

	t.Run(
		"AdminApproveSubmission_RejectAndSendOtherProposal_success", func(t *testing.T) {
			symbol := entity.Symbol{
				Id:     1,
				Symbol: "ABC",
			}
			request := entity.LoanPackageRequest{
				Id:         1,
				Status:     entity.LoanPackageRequestStatusPending,
				AssetType:  entity.AssetTypeUnderlying,
				SymbolId:   symbol.Id,
				InvestorId: "investorId",
			}
			confirmedRequest := entity.LoanPackageRequest{
				Id:         1,
				Status:     entity.LoanPackageRequestStatusConfirmed,
				AssetType:  entity.AssetTypeUnderlying,
				SymbolId:   symbol.Id,
				InvestorId: "investorId",
			}
			submissionSheet := entity.SubmissionSheet{
				Metadata: entity.SubmissionSheetMetadata{
					Id:                   1,
					Status:               entity.SubmissionSheetStatusSubmitted,
					LoanPackageRequestId: request.Id,
					FlowType:             entity.FlowTypeDnseOnline,
				},
				Detail: entity.SubmissionSheetDetail{
					Id: 1,
					LoanPolicies: []entity.LoanPolicySnapShot{
						{
							Term:        30,
							InitialRate: decimal.NewFromFloat(0.2),
						},
					},
				},
			}
			expireDate := time.Date(2021, 0, 0, 0, 0, 0, 0, time.Local)
			offer := entity.LoanPackageOffer{
				Id:                   1,
				LoanPackageRequestId: request.Id,
				ExpiredAt:            expireDate,
			}
			offerInterests := []entity.LoanPackageOfferInterest{
				{
					LoanPackageOfferId:      offer.Id,
					SubmissionSheetDetailId: submissionSheet.Detail.Id,
				},
			}
			accounts := []entity.FinancialAccountDetail{
				{
					AccountNo: "abc",
				},
			}
			submissionSheetRepo.EXPECT().GetById(testifyMock.Anything, submissionSheet.Metadata.Id).Return(submissionSheet, nil).Once()
			loanPackageRequestRepo.EXPECT().GetById(testifyMock.Anything, request.Id, testifyMock.Anything).Return(request, nil).Once()
			financingRepo.EXPECT().GetDateAfter(testifyMock.Anything, appConfig.LoanRequest.ExpireDays).Return(expireDate, nil).Once()
			submissionSheetRepo.EXPECT().UpdateMetadataStatusById(testifyMock.Anything, submissionSheet.Metadata.Id, entity.SubmissionSheetStatusApproved).Return(nil).Once()
			loanPackageRequestRepo.EXPECT().UpdateStatusById(testifyMock.Anything, request.Id, entity.LoanPackageRequestStatusConfirmed).Return(confirmedRequest, nil).Once()
			loanPackageOfferRepository.EXPECT().Create(testifyMock.Anything, testifyMock.Anything).Return(offer, nil).Once()
			loanPackageOfferInterestRepository.EXPECT().BulkCreate(testifyMock.Anything, testifyMock.Anything).Return(offerInterests, nil).Once()
			symbolRepo.EXPECT().GetById(testifyMock.Anything, request.SymbolId).Return(symbol, nil).Once()
			financialProductRepo.EXPECT().GetAllAccountDetail(testifyMock.Anything, request.InvestorId).Return(accounts, nil).Once()
			loanPackageRequestEventRepository.EXPECT().NotifyOnlineConfirmation(testifyMock.Anything, testifyMock.Anything).Return(nil).Once()
			err := useCase.AdminApproveSubmission(context.Background(), 1)
			assert.Nil(t, err)
		},
	)

	t.Run(
		"AdminApproveSubmission_Approve_success", func(t *testing.T) {
			symbol := entity.Symbol{
				Id:     1,
				Symbol: "ABC",
			}
			request := entity.LoanPackageRequest{
				Id:         1,
				Status:     entity.LoanPackageRequestStatusPending,
				AssetType:  entity.AssetTypeUnderlying,
				SymbolId:   symbol.Id,
				InvestorId: "investorId",
			}
			confirmedRequest := entity.LoanPackageRequest{
				Id:         1,
				Status:     entity.LoanPackageRequestStatusConfirmed,
				AssetType:  entity.AssetTypeUnderlying,
				SymbolId:   symbol.Id,
				InvestorId: "investorId",
			}
			submissionSheet := entity.SubmissionSheet{
				Metadata: entity.SubmissionSheetMetadata{
					Id:                   1,
					Status:               entity.SubmissionSheetStatusSubmitted,
					LoanPackageRequestId: request.Id,
					FlowType:             entity.FlowTypeDnseOnline,
				},
				Detail: entity.SubmissionSheetDetail{
					Id: 1,
					LoanPolicies: []entity.LoanPolicySnapShot{
						{
							Term:        30,
							InitialRate: decimal.NewFromFloat(0.2),
						},
					},
				},
			}
			expireDate := time.Date(2021, 0, 0, 0, 0, 0, 0, time.Local)
			offer := entity.LoanPackageOffer{
				Id:                   1,
				LoanPackageRequestId: request.Id,
				ExpiredAt:            expireDate,
			}
			offerInterests := []entity.LoanPackageOfferInterest{
				{
					LoanPackageOfferId:      offer.Id,
					SubmissionSheetDetailId: submissionSheet.Detail.Id,
				},
			}
			accounts := []entity.FinancialAccountDetail{
				{
					AccountNo: "abc",
				},
			}
			submissionSheetRepo.EXPECT().GetById(testifyMock.Anything, submissionSheet.Metadata.Id).Return(submissionSheet, nil).Once()
			loanPackageRequestRepo.EXPECT().GetById(testifyMock.Anything, request.Id, testifyMock.Anything).Return(request, nil).Once()
			financingRepo.EXPECT().GetDateAfter(testifyMock.Anything, appConfig.LoanRequest.ExpireDays).Return(expireDate, nil).Once()
			submissionSheetRepo.EXPECT().UpdateMetadataStatusById(testifyMock.Anything, submissionSheet.Metadata.Id, entity.SubmissionSheetStatusApproved).Return(nil).Once()
			loanPackageRequestRepo.EXPECT().UpdateStatusById(testifyMock.Anything, request.Id, entity.LoanPackageRequestStatusConfirmed).Return(confirmedRequest, nil).Once()
			loanPackageOfferRepository.EXPECT().Create(testifyMock.Anything, testifyMock.Anything).Return(offer, nil).Once()
			loanPackageOfferInterestRepository.EXPECT().BulkCreate(testifyMock.Anything, testifyMock.Anything).Return(offerInterests, nil).Once()
			symbolRepo.EXPECT().GetById(testifyMock.Anything, request.SymbolId).Return(symbol, nil).Once()
			financialProductRepo.EXPECT().GetAllAccountDetail(testifyMock.Anything, request.InvestorId).Return(accounts, nil).Once()
			loanPackageRequestEventRepository.EXPECT().NotifyOnlineConfirmation(testifyMock.Anything, testifyMock.Anything).Return(nil).Once()
			err := useCase.AdminApproveSubmission(context.Background(), 1)
			assert.Nil(t, err)
		},
	)
}

func TestLoanPackageRequestUseCase_AdminRejectSubmission(t *testing.T) {
	t.Parallel()
	appConfig := config.AppConfig{
		LoanRequest: config.LoanRequestConfig{
			ExpireDays: 7,
		},
	}
	loanPackageRequestRepo := mock.NewMockLoanPackageRequestRepository(t)
	loanPackageOfferRepository := mock.NewMockLoanPackageOfferRepository(t)
	loanPackageOfferInterestRepository := mock.NewMockLoanPackageOfferInterestRepository(t)
	loanPackageRequestEventRepository := mock.NewMockLoanPackageRequestEventRepository(t)
	symbolRepo := mock.NewMockSymbolRepository(t)
	financialProductRepo := mock.NewMockFinancialProductRepository(t)
	financingRepo := mock.NewMockFinancingRepository(t)
	loanPolicyTemplateRepo := mock.NewMockLoanPolicyTemplateRepository(t)
	submissionSheetRepo := mock.NewMockSubmissionSheetRepository(t)
	marginOperationRepo := mock.NewMockMarginOperationRepository(t)
	atomicExecutor := mock.NewMockAtomicExecutorExecutePassthrough(t)
	errorService := mock.ErrReporter{}
	useCase := NewUseCase(submissionSheetRepo, atomicExecutor, loanPolicyTemplateRepo, marginOperationRepo, financialProductRepo, loanPackageRequestRepo, loanPackageOfferRepository, loanPackageOfferInterestRepository, financingRepo, appConfig, errorService, loanPackageRequestEventRepository, symbolRepo)

	t.Run(
		"AdminRejectSubmission_success", func(t *testing.T) {
			submissionId := int64(1)
			submissionSheetRepo.EXPECT().GetById(testifyMock.Anything, submissionId).Return(entity.SubmissionSheet{
				Metadata: entity.SubmissionSheetMetadata{
					Status: entity.SubmissionSheetStatusSubmitted,
				},
			}, nil).Once()
			submissionSheetRepo.EXPECT().UpdateMetadataStatusById(testifyMock.Anything, submissionId, entity.SubmissionSheetStatusRejected).Return(nil).Once()
			err := useCase.AdminRejectSubmission(context.Background(), submissionId)
			assert.Nil(t, err)
		})

	t.Run(
		"AdminRejectSubmission_GetById_error", func(t *testing.T) {
			submissionId := int64(1)
			submissionSheetRepo.EXPECT().GetById(testifyMock.Anything, submissionId).Return(entity.SubmissionSheet{}, assert.AnError).Once()
			err := useCase.AdminRejectSubmission(context.Background(), submissionId)
			assert.ErrorIs(t, err, assert.AnError)
		})

	t.Run(
		"AdminRejectSubmission_success", func(t *testing.T) {
			submissionId := int64(1)
			submissionSheetRepo.EXPECT().GetById(testifyMock.Anything, submissionId).Return(entity.SubmissionSheet{
				Metadata: entity.SubmissionSheetMetadata{
					Status: entity.SubmissionSheetStatusRejected,
				},
			}, nil).Once()
			err := useCase.AdminRejectSubmission(context.Background(), submissionId)
			assert.ErrorIs(t, err, apperrors.ErrorInvalidCurrentSubmissionStatus)
		})

	t.Run(
		"AdminRejectSubmission_UpdateMetadataStatusById_error", func(t *testing.T) {
			submissionId := int64(1)
			submissionSheetRepo.EXPECT().GetById(testifyMock.Anything, submissionId).Return(entity.SubmissionSheet{
				Metadata: entity.SubmissionSheetMetadata{
					Status: entity.SubmissionSheetStatusSubmitted,
				},
			}, nil).Once()
			submissionSheetRepo.EXPECT().UpdateMetadataStatusById(testifyMock.Anything, submissionId, entity.SubmissionSheetStatusRejected).Return(assert.AnError).Once()
			err := useCase.AdminRejectSubmission(context.Background(), submissionId)
			assert.ErrorIs(t, err, assert.AnError)
		})
}
