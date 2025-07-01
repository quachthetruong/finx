package loanpackagerequest

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"

	"financing-offer/internal/atomicity"
	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"financing-offer/test/mock"
)

func TestLoanRequestUseCase(t *testing.T) {
	t.Parallel()

	t.Run(
		"SystemDeclineRiskLoanRequestsSuccess", func(t *testing.T) {
			db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			if err != nil {
				t.Errorf("%v", err)
			}
			loanPackageRequestRepo := mock.NewMockLoanPackageRequestRepository(t)
			scoreGroupInterestRepo := mock.NewMockScoreGroupInterestRepository(t)
			loanPackageOfferRepository := mock.NewMockLoanPackageOfferRepository(t)
			loanPackageOfferInterestRepository := mock.NewMockLoanPackageOfferInterestRepository(t)
			loanPackageRequestEventRepository := mock.NewMockLoanPackageRequestEventRepository(t)
			symbolRepo := mock.NewMockSymbolRepository(t)
			loanContractRepo := mock.NewMockLoanContractPersistenceRepository(t)
			financialProductRepo := mock.NewMockFinancialProductRepository(t)
			financingRepo := mock.NewMockFinancingRepository(t)
			schedulerJobRepo := mock.NewMockSchedulerJobRepository(t)
			investorRepo := mock.NewMockInvestorPersistenceRepository(t)
			loanPolicyTemplateRepo := mock.NewMockLoanPolicyTemplateRepository(t)
			submissionSheetRepo := mock.NewMockSubmissionSheetRepository(t)
			marginOperationRepo := mock.NewMockMarginOperationRepository(t)
			configurationRepo := mock.NewMockConfigurationPersistenceRepository(t)
			odooServiceRepo := mock.NewMockOdooServiceRepository(t)
			useCase := NewUseCase(
				loanPackageRequestRepo,
				&atomicity.DbAtomicExecutor{
					DB: db,
				},
				scoreGroupInterestRepo,
				loanPackageOfferRepository,
				loanPackageOfferInterestRepository,
				loanPackageRequestEventRepository,
				symbolRepo,
				loanContractRepo,
				financialProductRepo,
				config.AppConfig{},
				loanPolicyTemplateRepo,
				slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				financingRepo,
				schedulerJobRepo,
				mock.ErrReporter{},
				investorRepo,
				submissionSheetRepo,
				marginOperationRepo,
				configurationRepo,
				odooServiceRepo,
			)
			sqlMock.ExpectBegin()
			sqlMock.ExpectCommit()
			loanPackageRequestRepo.On(
				"LockAllPendingRequestByMaxPercent", testifyMock.Anything, decimal.NewFromFloat(0.3),
				testifyMock.Anything,
			).
				Return(
					[]entity.LoanPackageRequest{
						{
							Id:          1,
							SymbolId:    1,
							InvestorId:  "test",
							AccountNo:   "accNo",
							LoanRate:    decimal.NewFromFloat(0.3),
							LimitAmount: decimal.NewFromFloat(300000.0),
							Type:        entity.LoanPackageRequestTypeFlexible,
							Status:      entity.LoanPackageRequestStatusPending,
							AssetType:   entity.AssetTypeUnderlying,
						},
					}, nil,
				)
			loanPackageRequestRepo.On(
				"UpdateStatusByLoanRequestIds", testifyMock.Anything, []int64{1},
				entity.LoanPackageRequestStatusConfirmed,
			).
				Return(
					[]entity.LoanPackageRequest{
						{
							Id:          1,
							SymbolId:    1,
							InvestorId:  "test",
							AccountNo:   "accNo",
							LoanRate:    decimal.NewFromFloat(0.3),
							LimitAmount: decimal.NewFromFloat(300000.0),
							Type:        entity.LoanPackageRequestTypeFlexible,
							Status:      entity.LoanPackageRequestStatusPending,
						},
					}, nil,
				)

			loanPackageOfferRepository.On("BulkCreate", testifyMock.Anything, testifyMock.Anything).
				Return(
					[]entity.LoanPackageOffer{
						{
							Id:                   1,
							LoanPackageRequestId: 1,
							OfferedBy:            "system",
							CreatedAt:            time.Now(),
							UpdatedAt:            time.Now(),
							ExpiredAt:            time.Now(),
							FlowType:             entity.FlowTypeDnseOnline,
						},
					}, nil,
				)

			schedulerJobRepo.On("Create", testifyMock.Anything, testifyMock.Anything).
				Return(nil)

			symbolRepo.On("GetById", testifyMock.Anything, int64(1)).Return(
				entity.Symbol{
					Id:              1,
					StockExchangeId: 1,
					Symbol:          "ACB",
				}, nil,
			)

			financialProductRepo.On("GetAllAccountDetail", testifyMock.Anything, "test").
				Return(
					[]entity.FinancialAccountDetail{
						{
							Id:         "1",
							InvestorId: "1",
						},
					}, nil,
				)
			loanPackageRequestEventRepository.On("NotifyRequestDeclined", testifyMock.Anything, testifyMock.Anything).
				Return(nil)

			err = useCase.SystemDeclineRiskLoanRequests(context.Background(), decimal.NewFromFloat(0.3))
			assert.Nil(t, err)
		},
	)

	t.Run(
		"SystemDeclineRiskLoanRequests_UpdateStatusByLoanIdsFail", func(t *testing.T) {
			db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			if err != nil {
				t.Errorf("%v", err)
			}
			loanPackageRequestRepo := mock.NewMockLoanPackageRequestRepository(t)
			scoreGroupInterestRepo := mock.NewMockScoreGroupInterestRepository(t)
			loanPackageOfferRepository := mock.NewMockLoanPackageOfferRepository(t)
			loanPackageOfferInterestRepository := mock.NewMockLoanPackageOfferInterestRepository(t)
			loanPackageRequestEventRepository := mock.NewMockLoanPackageRequestEventRepository(t)
			symbolRepo := mock.NewMockSymbolRepository(t)
			loanContractRepo := mock.NewMockLoanContractPersistenceRepository(t)
			financialProductRepo := mock.NewMockFinancialProductRepository(t)
			financingRepo := mock.NewMockFinancingRepository(t)
			schedulerJobRepo := mock.NewMockSchedulerJobRepository(t)
			investorRepo := mock.NewMockInvestorPersistenceRepository(t)
			loanPolicyTemplateRepo := mock.NewMockLoanPolicyTemplateRepository(t)
			submissionSheetRepo := mock.NewMockSubmissionSheetRepository(t)
			marginOperationRepo := mock.NewMockMarginOperationRepository(t)
			configurationRepo := mock.NewMockConfigurationPersistenceRepository(t)
			odooServiceRepo := mock.NewMockOdooServiceRepository(t)
			useCase := NewUseCase(
				loanPackageRequestRepo,
				&atomicity.DbAtomicExecutor{
					DB: db,
				},
				scoreGroupInterestRepo,
				loanPackageOfferRepository,
				loanPackageOfferInterestRepository,
				loanPackageRequestEventRepository,
				symbolRepo,
				loanContractRepo,
				financialProductRepo,
				config.AppConfig{},
				loanPolicyTemplateRepo,
				slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				financingRepo,
				schedulerJobRepo,
				mock.ErrReporter{},
				investorRepo,
				submissionSheetRepo,
				marginOperationRepo,
				configurationRepo,
				odooServiceRepo,
			)
			sqlMock.ExpectBegin()
			sqlMock.ExpectRollback()
			loanPackageRequestRepo.On(
				"LockAllPendingRequestByMaxPercent", testifyMock.Anything, decimal.NewFromFloat(0.3),
				testifyMock.Anything,
			).
				Return(
					[]entity.LoanPackageRequest{
						{
							Id:          1,
							SymbolId:    1,
							InvestorId:  "test",
							AccountNo:   "accNo",
							LoanRate:    decimal.NewFromFloat(0.3),
							LimitAmount: decimal.NewFromFloat(300000.0),
							Type:        entity.LoanPackageRequestTypeFlexible,
							Status:      entity.LoanPackageRequestStatusPending,
						},
					}, nil,
				)
			loanPackageRequestRepo.On(
				"UpdateStatusByLoanRequestIds", testifyMock.Anything, []int64{1},
				entity.LoanPackageRequestStatusConfirmed,
			).
				Return([]entity.LoanPackageRequest{}, fmt.Errorf("test"))
			schedulerJobRepo.On(
				"Create", testifyMock.Anything, entity.SchedulerJob{
					JobType:      entity.JobTypeDeclineHighRiskLoanRequest,
					JobStatus:    entity.JobStatusFail,
					TriggerBy:    "system",
					TrackingData: "{\"error\":\"SystemDeclineRiskLoanRequests systemDeclineLoanRequestIds UpdateStatusByLoanRequestIds test\"}",
				},
			).Return(nil)
			err = useCase.SystemDeclineRiskLoanRequests(context.Background(), decimal.NewFromFloat(0.3))
			assert.Equal(
				t,
				"SystemDeclineRiskLoanRequests SystemDeclineRiskLoanRequests systemDeclineLoanRequestIds UpdateStatusByLoanRequestIds test",
				err.Error(),
			)
		},
	)

	t.Run(
		"SystemDeclineRiskLoanRequests_BulkCreateLoanOfferFail", func(t *testing.T) {
			db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			if err != nil {
				t.Errorf("%v", err)
			}
			loanPackageRequestRepo := mock.NewMockLoanPackageRequestRepository(t)
			scoreGroupInterestRepo := mock.NewMockScoreGroupInterestRepository(t)
			loanPackageOfferRepository := mock.NewMockLoanPackageOfferRepository(t)
			loanPackageOfferInterestRepository := mock.NewMockLoanPackageOfferInterestRepository(t)
			loanPackageRequestEventRepository := mock.NewMockLoanPackageRequestEventRepository(t)
			symbolRepo := mock.NewMockSymbolRepository(t)
			loanContractRepo := mock.NewMockLoanContractPersistenceRepository(t)
			financialProductRepo := mock.NewMockFinancialProductRepository(t)
			financingRepo := mock.NewMockFinancingRepository(t)
			schedulerJobRepo := mock.NewMockSchedulerJobRepository(t)
			investorRepo := mock.NewMockInvestorPersistenceRepository(t)
			loanPolicyTemplateRepo := mock.NewMockLoanPolicyTemplateRepository(t)
			submissionSheetRepo := mock.NewMockSubmissionSheetRepository(t)
			marginOperationRepo := mock.NewMockMarginOperationRepository(t)
			configurationRepo := mock.NewMockConfigurationPersistenceRepository(t)
			odooServiceRepo := mock.NewMockOdooServiceRepository(t)
			useCase := NewUseCase(
				loanPackageRequestRepo,
				&atomicity.DbAtomicExecutor{
					DB: db,
				},
				scoreGroupInterestRepo,
				loanPackageOfferRepository,
				loanPackageOfferInterestRepository,
				loanPackageRequestEventRepository,
				symbolRepo,
				loanContractRepo,
				financialProductRepo,
				config.AppConfig{},
				loanPolicyTemplateRepo,
				slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				financingRepo,
				schedulerJobRepo,
				mock.ErrReporter{},
				investorRepo,
				submissionSheetRepo,
				marginOperationRepo,
				configurationRepo,
				odooServiceRepo,
			)
			sqlMock.ExpectBegin()
			sqlMock.ExpectRollback()
			loanPackageRequestRepo.On(
				"LockAllPendingRequestByMaxPercent", testifyMock.Anything, decimal.NewFromFloat(0.3),
				testifyMock.Anything,
			).
				Return(
					[]entity.LoanPackageRequest{
						{
							Id:          1,
							SymbolId:    1,
							InvestorId:  "test",
							AccountNo:   "accNo",
							LoanRate:    decimal.NewFromFloat(0.3),
							LimitAmount: decimal.NewFromFloat(300000.0),
							Type:        entity.LoanPackageRequestTypeFlexible,
							Status:      entity.LoanPackageRequestStatusPending,
						},
					}, nil,
				)
			loanPackageRequestRepo.On(
				"UpdateStatusByLoanRequestIds", testifyMock.Anything, []int64{1},
				entity.LoanPackageRequestStatusConfirmed,
			).
				Return(
					[]entity.LoanPackageRequest{
						{
							Id:          1,
							SymbolId:    1,
							InvestorId:  "test",
							AccountNo:   "accNo",
							LoanRate:    decimal.NewFromFloat(0.3),
							LimitAmount: decimal.NewFromFloat(300000.0),
							Type:        entity.LoanPackageRequestTypeFlexible,
							Status:      entity.LoanPackageRequestStatusPending,
						},
					}, nil,
				)

			loanPackageOfferRepository.On("BulkCreate", testifyMock.Anything, testifyMock.Anything).
				Return([]entity.LoanPackageOffer{}, fmt.Errorf("test"))
			schedulerJobRepo.On(
				"Create", testifyMock.Anything, entity.SchedulerJob{
					JobType:      entity.JobTypeDeclineHighRiskLoanRequest,
					JobStatus:    entity.JobStatusFail,
					TriggerBy:    "system",
					TrackingData: "{\"error\":\"SystemDeclineRiskLoanRequests systemDeclineLoanRequestIds BulkCreate LoanOffer test\"}",
				},
			).
				Return(nil)

			err = useCase.SystemDeclineRiskLoanRequests(context.Background(), decimal.NewFromFloat(0.3))
			assert.Equal(
				t,
				"SystemDeclineRiskLoanRequests SystemDeclineRiskLoanRequests systemDeclineLoanRequestIds BulkCreate LoanOffer test",
				err.Error(),
			)
		},
	)

	t.Run(
		"SystemDeclineRiskLoanRequestsSuccess when AssetType = DERIVATIVE", func(t *testing.T) {
			db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			if err != nil {
				t.Errorf("%v", err)
			}
			loanPackageRequestRepo := mock.NewMockLoanPackageRequestRepository(t)
			scoreGroupInterestRepo := mock.NewMockScoreGroupInterestRepository(t)
			loanPackageOfferRepository := mock.NewMockLoanPackageOfferRepository(t)
			loanPackageOfferInterestRepository := mock.NewMockLoanPackageOfferInterestRepository(t)
			loanPackageRequestEventRepository := mock.NewMockLoanPackageRequestEventRepository(t)
			symbolRepo := mock.NewMockSymbolRepository(t)
			loanContractRepo := mock.NewMockLoanContractPersistenceRepository(t)
			financialProductRepo := mock.NewMockFinancialProductRepository(t)
			financingRepo := mock.NewMockFinancingRepository(t)
			schedulerJobRepo := mock.NewMockSchedulerJobRepository(t)
			investorRepo := mock.NewMockInvestorPersistenceRepository(t)
			loanPolicyTemplateRepo := mock.NewMockLoanPolicyTemplateRepository(t)
			submissionSheetRepo := mock.NewMockSubmissionSheetRepository(t)
			marginOperationRepo := mock.NewMockMarginOperationRepository(t)
			configurationRepo := mock.NewMockConfigurationPersistenceRepository(t)
			odooServiceRepo := mock.NewMockOdooServiceRepository(t)
			useCase := NewUseCase(
				loanPackageRequestRepo,
				&atomicity.DbAtomicExecutor{
					DB: db,
				},
				scoreGroupInterestRepo,
				loanPackageOfferRepository,
				loanPackageOfferInterestRepository,
				loanPackageRequestEventRepository,
				symbolRepo,
				loanContractRepo,
				financialProductRepo,
				config.AppConfig{},
				loanPolicyTemplateRepo,
				slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				financingRepo,
				schedulerJobRepo,
				mock.ErrReporter{},
				investorRepo,
				submissionSheetRepo,
				marginOperationRepo,
				configurationRepo,
				odooServiceRepo,
			)
			sqlMock.ExpectBegin()
			sqlMock.ExpectCommit()
			loanPackageRequestRepo.On(
				"LockAllPendingRequestByMaxPercent", testifyMock.Anything, decimal.NewFromFloat(0.3),
				testifyMock.Anything,
			).
				Return(
					[]entity.LoanPackageRequest{
						{
							Id:          1,
							SymbolId:    1,
							InvestorId:  "test",
							AccountNo:   "accNo",
							LoanRate:    decimal.NewFromFloat(0.3),
							LimitAmount: decimal.NewFromFloat(300000.0),
							Type:        entity.LoanPackageRequestTypeFlexible,
							Status:      entity.LoanPackageRequestStatusPending,
							AssetType:   entity.AssetTypeDerivative,
						},
					}, nil,
				)
			loanPackageRequestRepo.On(
				"UpdateStatusByLoanRequestIds", testifyMock.Anything, []int64{1},
				entity.LoanPackageRequestStatusConfirmed,
			).
				Return(
					[]entity.LoanPackageRequest{
						{
							Id:          1,
							SymbolId:    1,
							InvestorId:  "test",
							AccountNo:   "accNo",
							LoanRate:    decimal.NewFromFloat(0.3),
							LimitAmount: decimal.NewFromFloat(300000.0),
							Type:        entity.LoanPackageRequestTypeFlexible,
							Status:      entity.LoanPackageRequestStatusPending,
						},
					}, nil,
				)

			loanPackageOfferRepository.On("BulkCreate", testifyMock.Anything, testifyMock.Anything).
				Return(
					[]entity.LoanPackageOffer{
						{
							Id:                   1,
							LoanPackageRequestId: 1,
							OfferedBy:            "system",
							CreatedAt:            time.Now(),
							UpdatedAt:            time.Now(),
							ExpiredAt:            time.Now(),
							FlowType:             entity.FlowTypeDnseOnline,
						},
					}, nil,
				)

			schedulerJobRepo.On("Create", testifyMock.Anything, testifyMock.Anything).
				Return(nil)

			symbolRepo.On("GetById", testifyMock.Anything, int64(1)).Return(
				entity.Symbol{
					Id:              1,
					StockExchangeId: 1,
					Symbol:          "ACB",
				}, nil,
			)

			financialProductRepo.On("GetAllAccountDetail", testifyMock.Anything, "test").
				Return(
					[]entity.FinancialAccountDetail{
						{
							Id:         "1",
							InvestorId: "1",
						},
					}, nil,
				)
			loanPackageRequestEventRepository.On(
				"NotifyDerivativeRequestDeclined", testifyMock.Anything, testifyMock.Anything,
			).
				Return(nil)

			err = useCase.SystemDeclineRiskLoanRequests(context.Background(), decimal.NewFromFloat(0.3))
			assert.Nil(t, err)
		},
	)
}

func TestLoanPackageRequestUseCase_GetAllUnderlyingRequests(t *testing.T) {
	t.Parallel()
	loanPackageRequestRepo := mock.NewMockLoanPackageRequestRepository(t)
	scoreGroupInterestRepo := mock.NewMockScoreGroupInterestRepository(t)
	loanPackageOfferRepository := mock.NewMockLoanPackageOfferRepository(t)
	loanPackageOfferInterestRepository := mock.NewMockLoanPackageOfferInterestRepository(t)
	loanPackageRequestEventRepository := mock.NewMockLoanPackageRequestEventRepository(t)
	symbolRepo := mock.NewMockSymbolRepository(t)
	loanContractRepo := mock.NewMockLoanContractPersistenceRepository(t)
	financialProductRepo := mock.NewMockFinancialProductRepository(t)
	financingRepo := mock.NewMockFinancingRepository(t)
	schedulerJobRepo := mock.NewMockSchedulerJobRepository(t)
	investorRepo := mock.NewMockInvestorPersistenceRepository(t)
	loanPolicyTemplateRepo := mock.NewMockLoanPolicyTemplateRepository(t)
	submissionSheetRepo := mock.NewMockSubmissionSheetRepository(t)
	marginOperationRepo := mock.NewMockMarginOperationRepository(t)
	configurationRepo := mock.NewMockConfigurationPersistenceRepository(t)
	odooServiceRepo := mock.NewMockOdooServiceRepository(t)
	atomicExecutor := mock.NewMockAtomicExecutorExecutePassthrough(t)
	useCase := NewUseCase(
		loanPackageRequestRepo,
		atomicExecutor,
		scoreGroupInterestRepo,
		loanPackageOfferRepository,
		loanPackageOfferInterestRepository,
		loanPackageRequestEventRepository,
		symbolRepo,
		loanContractRepo,
		financialProductRepo,
		config.AppConfig{},
		loanPolicyTemplateRepo,
		slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		financingRepo,
		schedulerJobRepo,
		mock.ErrReporter{},
		investorRepo,
		submissionSheetRepo,
		marginOperationRepo,
		configurationRepo,
		odooServiceRepo,
	)
	t.Run(
		"GetAllUnderlyingRequests_success", func(t *testing.T) {
			underlyingRequests := []entity.UnderlyingLoanPackageRequest{
				{
					Id: 1,
				},
				{
					Id: 2,
				},
			}
			loanPackageRequestRepo.EXPECT().GetAllUnderlyingRequests(testifyMock.Anything, entity.UnderlyingLoanPackageFilter{Statuses: []entity.LoanPackageRequestStatus{entity.LoanPackageRequestStatusPending}}).Return(underlyingRequests, nil).Once()
			result, err := useCase.GetAllUnderlyingRequests(context.Background(), entity.UnderlyingLoanPackageFilter{Statuses: []entity.LoanPackageRequestStatus{entity.LoanPackageRequestStatusPending}})
			assert.Nil(t, err)
			assert.Equal(t, len(underlyingRequests), len(result))
		})

	t.Run(
		"GetAllUnderlyingRequests_error", func(t *testing.T) {

			loanPackageRequestRepo.EXPECT().GetAllUnderlyingRequests(testifyMock.Anything, entity.UnderlyingLoanPackageFilter{Statuses: []entity.LoanPackageRequestStatus{entity.LoanPackageRequestStatusPending}}).Return([]entity.UnderlyingLoanPackageRequest{}, assert.AnError).Once()
			_, err := useCase.GetAllUnderlyingRequests(context.Background(), entity.UnderlyingLoanPackageFilter{Statuses: []entity.LoanPackageRequestStatus{entity.LoanPackageRequestStatusPending}})
			assert.ErrorIs(t, err, assert.AnError)
		})
}
