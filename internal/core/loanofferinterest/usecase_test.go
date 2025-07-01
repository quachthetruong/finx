package loanofferinterest

import (
	"context"
	"financing-offer/internal/config"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/atomicity"
	"financing-offer/internal/core/entity"
	"financing-offer/test/mock"
)

func TestLoanPackageOfferInterestUseCase(t *testing.T) {
	t.Parallel()

	t.Run(
		"AdminAssignLoanIdByOfferId", func(t *testing.T) {
			db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			if err != nil {
				t.Errorf("%v", err)
			}
			loanPackageOfferRepository := mock.NewMockLoanPackageOfferRepository(t)
			loanPackageOfferInterestRepository := mock.NewMockLoanPackageOfferInterestRepository(t)
			loanPackageRequestEventRepository := mock.NewMockLoanPackageOfferInterestEventRepository(t)
			symbolRepo := mock.NewMockSymbolRepository(t)
			loanContractRepo := mock.NewMockLoanContractPersistenceRepository(t)
			financialProductRepo := mock.NewMockFinancialProductRepository(t)
			submissionSheetRepo := mock.NewMockSubmissionSheetRepository(t)
			policyTemplateRepo := mock.NewMockLoanPolicyTemplateRepository(t)
			appConfig := config.AppConfig{}
			useCase := NewUseCase(
				loanPackageOfferInterestRepository,
				&atomicity.DbAtomicExecutor{
					DB: db,
				},
				loanPackageOfferRepository,
				loanContractRepo,
				financialProductRepo,
				loanPackageRequestEventRepository,
				mock.ErrReporter{},
				symbolRepo,
				submissionSheetRepo,
				policyTemplateRepo,
				appConfig,
			)
			sqlMock.ExpectBegin()
			sqlMock.ExpectCommit()
			loanPackageOffer := entity.LoanPackageOffer{
				Id:                   1,
				LoanPackageRequestId: 1,
				OfferedBy:            "system",
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
				ExpiredAt:            time.Now(),
				FlowType:             entity.FLowTypeDnseOffline,
				LoanPackageRequest: &entity.LoanPackageRequest{
					Id:                 1,
					SymbolId:           1,
					InvestorId:         "1",
					AccountNo:          "1",
					LoanRate:           decimal.Decimal{},
					LimitAmount:        decimal.Decimal{},
					Type:               "",
					Status:             "",
					GuaranteedDuration: 0,
					AssetType:          entity.AssetTypeDerivative,
					InitialRate:        decimal.Decimal{},
					ContractSize:       0,
					CreatedAt:          time.Time{},
					UpdatedAt:          time.Time{},
					Investor:           entity.Investor{},
				},
			}
			loanPackageInterest := entity.LoanPackageOfferInterest{
				Id:                   1,
				LoanPackageOfferId:   1,
				ScoreGroupInterestId: 1,
				LimitAmount:          decimal.NewFromInt(1),
				LoanRate:             decimal.NewFromInt(1),
				InterestRate:         decimal.NewFromInt(1),
				Status:               entity.LoanPackageOfferInterestStatusPending,
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
				LoanID:               0,
				CancelledBy:          "",
				CancelledAt:          time.Time{},
				Term:                 0,
				FeeRate:              decimal.NewFromInt(1),
				CancelledReason:      "",
				AssetType:            entity.AssetTypeDerivative,
				InitialRate:          decimal.NewFromInt(1),
				ContractSize:         0,
			}
			loanPackageOfferRepository.
				On("FindByIdWithRequest", testifyMock.Anything, int64(1), testifyMock.Anything).Return(
				loanPackageOffer, nil,
			)

			loanPackageOfferInterestRepository.On("GetByOfferIdWithLock", testifyMock.Anything, int64(1)).Return(
				[]entity.LoanPackageOfferInterest{
					loanPackageInterest,
				}, nil,
			)

			financialProductRepo.
				On("GetLoanPackageDerivative", testifyMock.Anything, testifyMock.Anything).Return(
				entity.FinancialProductLoanPackageDerivative{
					Id:          1,
					Name:        "BTCUSDT",
					InitialRate: decimal.NewFromInt(1),
				}, nil,
			)

			loanPackageOfferInterestRepository.
				On("Update", testifyMock.Anything, testifyMock.Anything).
				Return(loanPackageInterest, nil)

			financialProductRepo.
				On(
					"AssignLoanPackageOrGetLoanPackageAccountId", testifyMock.Anything, testifyMock.Anything,
					testifyMock.Anything, testifyMock.Anything, testifyMock.Anything,
				).
				Return(
					int64(1), nil,
				)

			loanContractRepo.
				On("Create", testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).
				Return(entity.LoanContract{}, nil)

			symbolRepo.
				On("GetById", testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).
				Return(
					entity.Symbol{
						Id:     1,
						Symbol: "BTCUSDT",
					}, nil,
				)
			financialProductRepo.
				On("GetAllAccountDetail", testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).
				Return(
					[]entity.FinancialAccountDetail{
						{
							Id:          "",
							AccountNo:   "1",
							AccountType: "1",
						},
					}, nil,
				)
			loanPackageRequestEventRepository.On(
				"NotifyDerivativeLoanPackageOfferReady", testifyMock.Anything, testifyMock.Anything,
			).
				Return(nil)

			err = useCase.AdminAssignLoanIdByOfferId(context.Background(), 1, 1)
			assert.Nil(t, err)
		},
	)

	t.Run(
		"InvestorConfirmLoanPackageInterest success", func(t *testing.T) {
			db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			if err != nil {
				t.Errorf("%v", err)
			}
			loanPackageOfferRepository := mock.NewMockLoanPackageOfferRepository(t)
			loanPackageOfferInterestRepository := mock.NewMockLoanPackageOfferInterestRepository(t)
			loanPackageRequestEventRepository := mock.NewMockLoanPackageOfferInterestEventRepository(t)
			symbolRepo := mock.NewMockSymbolRepository(t)
			loanContractRepo := mock.NewMockLoanContractPersistenceRepository(t)
			financialProductRepo := mock.NewMockFinancialProductRepository(t)
			submissionSheetRepo := mock.NewMockSubmissionSheetRepository(t)
			policyTemplateRepo := mock.NewMockLoanPolicyTemplateRepository(t)
			appConfig := config.AppConfig{}
			useCase := NewUseCase(
				loanPackageOfferInterestRepository,
				&atomicity.DbAtomicExecutor{
					DB: db,
				},
				loanPackageOfferRepository,
				loanContractRepo,
				financialProductRepo,
				loanPackageRequestEventRepository,
				mock.ErrReporter{},
				symbolRepo,
				submissionSheetRepo,
				policyTemplateRepo,
				appConfig,
			)
			sqlMock.ExpectBegin()
			sqlMock.ExpectCommit()
			loanPackageRequest := entity.LoanPackageRequest{
				Id:                 1,
				SymbolId:           1,
				InvestorId:         "1",
				AccountNo:          "1",
				LoanRate:           decimal.Decimal{},
				LimitAmount:        decimal.Decimal{},
				Type:               "",
				Status:             "",
				GuaranteedDuration: 0,
				AssetType:          entity.AssetTypeUnderlying,
				InitialRate:        decimal.Decimal{},
				ContractSize:       0,
				CreatedAt:          time.Time{},
				UpdatedAt:          time.Time{},
				Investor:           entity.Investor{},
			}
			loanPackageOffer := entity.LoanPackageOffer{
				Id:                   1,
				LoanPackageRequestId: 1,
				OfferedBy:            "system",
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
				ExpiredAt:            time.Now().AddDate(1, 0, 0),
				FlowType:             entity.FlowTypeDnseOnline,
				LoanPackageRequest:   &loanPackageRequest,
			}
			loanPackageInterest := entity.LoanPackageOfferInterest{
				Id:                   1,
				LoanPackageOfferId:   1,
				ScoreGroupInterestId: 1,
				LimitAmount:          decimal.NewFromInt(1),
				LoanRate:             decimal.NewFromInt(1),
				InterestRate:         decimal.NewFromInt(1),
				Status:               entity.LoanPackageOfferInterestStatusPending,
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
				LoanID:               0,
				CancelledBy:          "",
				CancelledAt:          time.Time{},
				Term:                 0,
				FeeRate:              decimal.NewFromInt(1),
				CancelledReason:      "",
				AssetType:            entity.AssetTypeUnderlying,
				InitialRate:          decimal.NewFromInt(1),
				ContractSize:         0,
				LoanPackageOffer:     &loanPackageOffer,
			}
			symbol := entity.Symbol{
				Id:        1,
				Symbol:    "BTC",
				AssetType: entity.AssetTypeUnderlying,
				Status:    entity.SymbolStatusActive,
			}
			financialProductDetail := entity.FinancialAccountDetail{
				Id:          "1",
				AccountNo:   "1",
				AccountType: "1",
			}
			loanPackageOfferInterestRepository.EXPECT().GetByIds(testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).Return([]entity.LoanPackageOfferInterest{loanPackageInterest}, nil)
			loanPackageOfferRepository.EXPECT().FindByIdWithRequest(testifyMock.Anything, loanPackageOffer.Id).Return(loanPackageOffer, nil)

			loanPackageOfferInterestRepository.EXPECT().CancelByOfferId(testifyMock.Anything, testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).Return(nil)
			financialProductRepo.EXPECT().AssignLoanPackageOrGetLoanPackageAccountId(testifyMock.Anything, testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).Return(int64(1), nil)
			loanContractRepo.EXPECT().BulkCreate(testifyMock.Anything, testifyMock.Anything).Return(nil)
			loanPackageOfferInterestRepository.EXPECT().UpdateStatus(testifyMock.Anything, testifyMock.Anything, entity.LoanPackageOfferInterestStatusLoanPackageCreated).Return(nil)

			symbolRepo.EXPECT().GetById(testifyMock.Anything, testifyMock.Anything).Return(symbol, nil)
			financialProductRepo.EXPECT().GetAllAccountDetail(testifyMock.Anything, testifyMock.Anything).Return([]entity.FinancialAccountDetail{financialProductDetail}, nil)
			loanPackageRequestEventRepository.EXPECT().NotifyLoanPackageOfferReady(testifyMock.Anything, testifyMock.Anything).Return(nil)

			err = useCase.InvestorConfirmLoanPackageInterest(context.Background(), []int64{1}, "1")
			assert.Nil(t, err)
		})

	t.Run("InvestorConfirmLoanPackageInterest verifyLoanOfferLines ErrorLoanPackageOfferInterestIsCreating", func(t *testing.T) {
		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
		if err != nil {
			t.Errorf("%v", err)
		}
		loanPackageOfferRepository := mock.NewMockLoanPackageOfferRepository(t)
		loanPackageOfferInterestRepository := mock.NewMockLoanPackageOfferInterestRepository(t)
		loanPackageRequestEventRepository := mock.NewMockLoanPackageOfferInterestEventRepository(t)
		symbolRepo := mock.NewMockSymbolRepository(t)
		loanContractRepo := mock.NewMockLoanContractPersistenceRepository(t)
		financialProductRepo := mock.NewMockFinancialProductRepository(t)
		submissionSheetRepo := mock.NewMockSubmissionSheetRepository(t)
		policyTemplateRepo := mock.NewMockLoanPolicyTemplateRepository(t)
		appConfig := config.AppConfig{}
		useCase := NewUseCase(
			loanPackageOfferInterestRepository,
			&atomicity.DbAtomicExecutor{
				DB: db,
			},
			loanPackageOfferRepository,
			loanContractRepo,
			financialProductRepo,
			loanPackageRequestEventRepository,
			mock.ErrReporter{},
			symbolRepo,
			submissionSheetRepo,
			policyTemplateRepo,
			appConfig,
		)
		sqlMock.ExpectBegin()
		sqlMock.ExpectRollback()
		loanPackageRequest := entity.LoanPackageRequest{
			Id:                 1,
			SymbolId:           1,
			InvestorId:         "1",
			AccountNo:          "1",
			LoanRate:           decimal.Decimal{},
			LimitAmount:        decimal.Decimal{},
			Type:               "",
			Status:             "",
			GuaranteedDuration: 0,
			AssetType:          entity.AssetTypeUnderlying,
			InitialRate:        decimal.Decimal{},
			ContractSize:       0,
			CreatedAt:          time.Time{},
			UpdatedAt:          time.Time{},
			Investor:           entity.Investor{},
		}
		loanPackageOffer := entity.LoanPackageOffer{
			Id:                   1,
			LoanPackageRequestId: 1,
			OfferedBy:            "system",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			ExpiredAt:            time.Now().AddDate(1, 0, 0),
			FlowType:             entity.FlowTypeDnseOnline,
			LoanPackageRequest:   &loanPackageRequest,
		}
		loanPackageInterest := entity.LoanPackageOfferInterest{
			Id:                   1,
			LoanPackageOfferId:   1,
			ScoreGroupInterestId: 1,
			LimitAmount:          decimal.NewFromInt(1),
			LoanRate:             decimal.NewFromInt(1),
			InterestRate:         decimal.NewFromInt(1),
			Status:               entity.LoanPackageOfferInterestStatusCreatingLoanPackage,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			LoanID:               0,
			CancelledBy:          "",
			CancelledAt:          time.Time{},
			Term:                 0,
			FeeRate:              decimal.NewFromInt(1),
			CancelledReason:      "",
			AssetType:            entity.AssetTypeUnderlying,
			InitialRate:          decimal.NewFromInt(1),
			ContractSize:         0,
			LoanPackageOffer:     &loanPackageOffer,
		}
		loanPackageOfferInterestRepository.EXPECT().GetByIds(testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).Return([]entity.LoanPackageOfferInterest{loanPackageInterest}, nil)
		loanPackageOfferRepository.EXPECT().FindByIdWithRequest(testifyMock.Anything, loanPackageOffer.Id).Return(loanPackageOffer, nil)

		err = useCase.InvestorConfirmLoanPackageInterest(context.Background(), []int64{1}, "1")
		assert.ErrorIs(t, err, apperrors.ErrorLoanPackageOfferInterestIsCreating)
	})
}
