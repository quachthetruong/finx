package investor_account

import (
	"context"
	"testing"
	"time"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/core/entity"
	"financing-offer/test/mock"
)

func TestInvestorAccountUseCase(t *testing.T) {
	t.Parallel()

	t.Run(
		"VerifyAndUpdateInvestorAccountVersion verify investor_account version 2 success", func(t *testing.T) {
			investorAccountRepository := mock.NewMockInvestorAccountRepository(t)
			orderServiceRepository := mock.NewMockOrderServiceRepository(t)
			useCase := NewUseCase(
				investorAccountRepository,
				orderServiceRepository,
			)
			account := entity.InvestorAccount{
				AccountNo:    "1",
				InvestorId:   "1",
				MarginStatus: entity.MarginStatusVersion2,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			loanPackage := entity.AccountLoanPackage{
				Id:                       1,
				Name:                     "Loan Package 1",
				Type:                     "Type 1",
				BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.02),
				BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.02),
				TransferFee:              decimal.NewFromInt(200000),
				Description:              "This is a description for Loan Package 1",
				LoanProducts:             []entity.LoanProduct{},
				BasketId:                 1,
			}
			investorAccountRepository.EXPECT().GetByAccountNo(mock2.Anything, "1").Return(account, nil)
			orderServiceRepository.EXPECT().GetAllAccountLoanPackages(
				mock2.Anything, "1",
			).Return([]entity.AccountLoanPackage{loanPackage}, nil)
			res, err := useCase.VerifyAndUpdateInvestorAccountMarginStatus(context.Background(), account)
			assert.Nil(t, err)
			assert.Equal(t, account, res)
		},
	)

	t.Run(
		"VerifyAndUpdateInvestorAccountMarginStatus verify investor_account version 3 success", func(t *testing.T) {
			investorAccountRepository := mock.NewMockInvestorAccountRepository(t)
			orderServiceRepository := mock.NewMockOrderServiceRepository(t)
			useCase := NewUseCase(
				investorAccountRepository,
				orderServiceRepository,
			)
			account := entity.InvestorAccount{
				AccountNo:    "1",
				InvestorId:   "1",
				MarginStatus: entity.MarginStatusVersion3,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			investorAccountRepository.EXPECT().GetByAccountNo(mock2.Anything, "1").Return(account, nil)
			res, err := useCase.VerifyAndUpdateInvestorAccountMarginStatus(context.Background(), account)
			assert.Nil(t, err)
			assert.Equal(t, account, res)
		},
	)

	t.Run(
		"VerifyAndUpdateInvestorAccountMarginStatus get investor_account margin status from db error",
		func(t *testing.T) {
			investorAccountRepository := mock.NewMockInvestorAccountRepository(t)
			orderServiceRepository := mock.NewMockOrderServiceRepository(t)
			useCase := NewUseCase(
				investorAccountRepository,
				orderServiceRepository,
			)
			account := entity.InvestorAccount{}
			investorAccountRepository.EXPECT().GetByAccountNo(mock2.Anything, "1").Return(account, assert.AnError)
			_, err := useCase.VerifyAndUpdateInvestorAccountMarginStatus(
				context.Background(), entity.InvestorAccount{
					AccountNo: "1",
				},
			)
			assert.ErrorIs(t, err, assert.AnError)
		},
	)

	t.Run(
		"VerifyAndUpdateInvestorAccountMarginStatus get all investor_account loan package error", func(t *testing.T) {
			investorAccountRepository := mock.NewMockInvestorAccountRepository(t)
			orderServiceRepository := mock.NewMockOrderServiceRepository(t)
			useCase := NewUseCase(
				investorAccountRepository,
				orderServiceRepository,
			)
			account := entity.InvestorAccount{
				AccountNo:    "1",
				InvestorId:   "1",
				MarginStatus: entity.MarginStatusVersion2,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			investorAccountRepository.EXPECT().GetByAccountNo(mock2.Anything, "1").Return(
				entity.InvestorAccount{}, qrm.ErrNoRows,
			)
			orderServiceRepository.EXPECT().GetAllAccountLoanPackages(
				mock2.Anything, "1",
			).Return([]entity.AccountLoanPackage{}, assert.AnError)
			_, err := useCase.VerifyAndUpdateInvestorAccountMarginStatus(context.Background(), account)
			assert.ErrorIs(t, err, assert.AnError)
		},
	)

	t.Run(
		"VerifyAndUpdateInvestorAccountMarginStatus get investor_account margin status from db no row then create",
		func(t *testing.T) {
			investorAccountRepository := mock.NewMockInvestorAccountRepository(t)
			orderServiceRepository := mock.NewMockOrderServiceRepository(t)
			useCase := NewUseCase(
				investorAccountRepository,
				orderServiceRepository,
			)
			account := entity.InvestorAccount{
				AccountNo:    "1",
				InvestorId:   "1",
				MarginStatus: entity.MarginStatusVersion2,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			loanPackage := entity.AccountLoanPackage{
				Id:                       1,
				Name:                     "Loan Package 1",
				Type:                     "M",
				BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.02),
				BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.02),
				TransferFee:              decimal.NewFromFloat(300),
				Description:              "This is a description for Loan Package 1",
				LoanProducts: []entity.LoanProduct{
					{
						Id:   "1-SCB",
						Name: "Loan Product 1",
					},
				},
				BasketId: 1,
			}
			investorAccountRepository.EXPECT().GetByAccountNo(mock2.Anything, "1").Return(
				entity.InvestorAccount{}, qrm.ErrNoRows,
			)
			orderServiceRepository.EXPECT().GetAllAccountLoanPackages(
				mock2.Anything, "1",
			).Return([]entity.AccountLoanPackage{loanPackage}, nil)
			investorAccountRepository.EXPECT().Create(mock2.Anything, mock2.Anything).Return(account, nil)
			res, err := useCase.VerifyAndUpdateInvestorAccountMarginStatus(context.Background(), account)
			assert.Nil(t, err)
			assert.Equal(t, account, res)
		},
	)

	t.Run(
		"VerifyAndUpdateInvestorAccountMarginStatus get investor_account margin status from db no row then create error",
		func(t *testing.T) {
			investorAccountRepository := mock.NewMockInvestorAccountRepository(t)
			orderServiceRepository := mock.NewMockOrderServiceRepository(t)
			useCase := NewUseCase(
				investorAccountRepository,
				orderServiceRepository,
			)
			account := entity.InvestorAccount{
				AccountNo:    "1",
				InvestorId:   "1",
				MarginStatus: entity.MarginStatusVersion2,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			loanPackage := entity.AccountLoanPackage{
				Id:                       1,
				Name:                     "Loan Package 1",
				Type:                     "M",
				BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.02),
				BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.02),
				TransferFee:              decimal.NewFromFloat(300),
				Description:              "This is a description for Loan Package 1",
				LoanProducts: []entity.LoanProduct{
					{
						Id:   "1-SCB",
						Name: "Loan Product 1",
					},
				},
				BasketId: 1,
			}
			investorAccountRepository.EXPECT().GetByAccountNo(mock2.Anything, "1").Return(
				entity.InvestorAccount{}, qrm.ErrNoRows,
			)
			orderServiceRepository.EXPECT().GetAllAccountLoanPackages(
				mock2.Anything, "1",
			).Return([]entity.AccountLoanPackage{loanPackage}, nil)
			investorAccountRepository.EXPECT().Create(mock2.Anything, mock2.Anything).Return(
				entity.InvestorAccount{}, assert.AnError,
			)
			_, err := useCase.VerifyAndUpdateInvestorAccountMarginStatus(context.Background(), account)
			assert.ErrorIs(t, err, assert.AnError)
		},
	)

	t.Run(
		"VerifyAndUpdateInvestorAccountMarginStatus get investor_account margin status from db version 2 then update to version 3",
		func(t *testing.T) {
			investorAccountRepository := mock.NewMockInvestorAccountRepository(t)
			orderServiceRepository := mock.NewMockOrderServiceRepository(t)
			useCase := NewUseCase(
				investorAccountRepository,
				orderServiceRepository,
			)
			account := entity.InvestorAccount{
				AccountNo:    "1",
				InvestorId:   "1",
				MarginStatus: entity.MarginStatusVersion2,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			loanPackage := entity.AccountLoanPackage{
				Id:                       1,
				Name:                     "Loan Package 1",
				Type:                     "M",
				BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.02),
				BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.02),
				TransferFee:              decimal.NewFromFloat(200000),
				Description:              "This is a description for Loan Package 1",
				LoanProducts: []entity.LoanProduct{
					{
						Id:   "1-120",
						Name: "Loan Product 1",
					},
				},
				BasketId: 1,
			}
			investorAccountRepository.EXPECT().GetByAccountNo(mock2.Anything, "1").Return(account, nil)
			orderServiceRepository.EXPECT().GetAllAccountLoanPackages(
				mock2.Anything, "1",
			).Return([]entity.AccountLoanPackage{loanPackage}, nil)
			account.MarginStatus = entity.MarginStatusVersion3
			investorAccountRepository.EXPECT().Update(mock2.Anything, mock2.Anything).Return(account, nil)
			res, err := useCase.VerifyAndUpdateInvestorAccountMarginStatus(context.Background(), account)
			assert.Nil(t, err)
			assert.Equal(t, account.MarginStatus, res.MarginStatus)
		},
	)

	t.Run(
		"VerifyAndUpdateInvestorAccountMarginStatus get investor_account margin status from db version 2 then update to version 3 error",
		func(t *testing.T) {
			investorAccountRepository := mock.NewMockInvestorAccountRepository(t)
			orderServiceRepository := mock.NewMockOrderServiceRepository(t)
			useCase := NewUseCase(
				investorAccountRepository,
				orderServiceRepository,
			)
			account := entity.InvestorAccount{
				AccountNo:    "1",
				InvestorId:   "1",
				MarginStatus: entity.MarginStatusVersion2,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			loanPackage := entity.AccountLoanPackage{
				Id:                       1,
				Name:                     "Loan Package 1",
				Type:                     "M",
				BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.02),
				BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.02),
				TransferFee:              decimal.NewFromFloat(200000),
				Description:              "This is a description for Loan Package 1",
				LoanProducts: []entity.LoanProduct{
					{
						Id:   "1-1204",
						Name: "Loan Product 1",
					},
				},
				BasketId: 1,
			}
			investorAccountRepository.EXPECT().GetByAccountNo(mock2.Anything, "1").Return(account, nil)
			orderServiceRepository.EXPECT().GetAllAccountLoanPackages(
				mock2.Anything, "1",
			).Return([]entity.AccountLoanPackage{loanPackage}, nil)
			account.MarginStatus = entity.MarginStatusVersion3
			investorAccountRepository.EXPECT().Update(mock2.Anything, mock2.Anything).Return(
				entity.InvestorAccount{}, assert.AnError,
			)
			_, err := useCase.VerifyAndUpdateInvestorAccountMarginStatus(context.Background(), account)
			assert.ErrorIs(t, err, assert.AnError)
		},
	)

	t.Run(
		"VerifyAndUpdateInvestorAccountMarginStatus checkMarginStatus non margin account error", func(t *testing.T) {
			investorAccountRepository := mock.NewMockInvestorAccountRepository(t)
			orderServiceRepository := mock.NewMockOrderServiceRepository(t)
			useCase := NewUseCase(
				investorAccountRepository,
				orderServiceRepository,
			)
			account := entity.InvestorAccount{
				AccountNo:    "1",
				InvestorId:   "1",
				MarginStatus: entity.MarginStatusVersion2,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			loanPackage := entity.AccountLoanPackage{
				Id:                       1,
				Name:                     "Loan Package 1",
				Type:                     "N",
				BrokerFirmBuyingFeeRate:  decimal.NewFromFloat(0.02),
				BrokerFirmSellingFeeRate: decimal.NewFromFloat(0.02),
				TransferFee:              decimal.NewFromFloat(200000),
				Description:              "This is a description for Loan Package 1",
				LoanProducts: []entity.LoanProduct{
					{
						Id:   "1-SCB",
						Name: "Loan Product 1",
					},
				},
				BasketId: 1,
			}
			investorAccountRepository.EXPECT().GetByAccountNo(mock2.Anything, "1").Return(account, nil)
			orderServiceRepository.EXPECT().GetAllAccountLoanPackages(
				mock2.Anything, "1",
			).Return([]entity.AccountLoanPackage{loanPackage}, nil)
			_, err := useCase.VerifyAndUpdateInvestorAccountMarginStatus(context.Background(), account)
			assert.ErrorIs(t, err, apperrors.ErrNonMarginAccount)
		},
	)

	t.Run(
		"VerifyAndUpdateInvestorAccountMarginStatus checkMarginStatus get all account loan package empty",
		func(t *testing.T) {
			investorAccountRepository := mock.NewMockInvestorAccountRepository(t)
			orderServiceRepository := mock.NewMockOrderServiceRepository(t)
			useCase := NewUseCase(
				investorAccountRepository,
				orderServiceRepository,
			)
			account := entity.InvestorAccount{
				AccountNo:    "1",
				InvestorId:   "1",
				MarginStatus: entity.MarginStatusVersion2,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			investorAccountRepository.EXPECT().GetByAccountNo(mock2.Anything, "1").Return(account, nil)
			orderServiceRepository.EXPECT().GetAllAccountLoanPackages(
				mock2.Anything, "1",
			).Return([]entity.AccountLoanPackage{}, nil)
			_, err := useCase.VerifyAndUpdateInvestorAccountMarginStatus(context.Background(), account)
			assert.ErrorIs(t, err, apperrors.ErrAccountNoInvalid)
		},
	)
}
