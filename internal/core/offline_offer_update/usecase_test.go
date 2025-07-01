package offlineofferupdate

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"

	"financing-offer/internal/atomicity"
	"financing-offer/internal/core/entity"
	"financing-offer/test/mock"
)

func TestLoanOfferUpdateUseCase(t *testing.T) {
	t.Parallel()

	t.Run(
		"CreateLoanOfferUpdateProcessingSuccess", func(t *testing.T) {
			db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			if err != nil {
				t.Errorf("%v", err)
			}
			offlineOfferUpdateRepo := mock.NewMockOfflineOfferUpdatePersistenceRepository(t)
			offerRepo := mock.NewMockLoanPackageOfferRepository(t)
			offerInterestRepo := mock.NewMockLoanPackageOfferInterestRepository(t)
			useCase := NewUseCase(
				offlineOfferUpdateRepo, offerRepo, offerInterestRepo,
				&atomicity.DbAtomicExecutor{
					DB: db,
				},
			)
			sqlMock.ExpectBegin()
			sqlMock.ExpectCommit()
			offlineOfferUpdateRepo.On("Create", testifyMock.Anything, testifyMock.Anything).
				Return(
					entity.OfflineOfferUpdate{
						Id:       1,
						OfferId:  1,
						Status:   entity.OfflineOfferUpdateStatusProcessing,
						Category: "1",
						Note:     "note",
					}, nil,
				)
			offerRepo.On("InvestorGetById", testifyMock.Anything, testifyMock.Anything).
				Return(
					entity.LoanPackageOffer{
						Id:                   1,
						LoanPackageRequestId: 1,
						OfferedBy:            "admin",
						CreatedAt:            time.Now(),
						UpdatedAt:            time.Now(),
						ExpiredAt:            time.Now().Add(time.Hour * 24),
						FlowType:             entity.FLowTypeDnseOffline,
						LoanPackageRequest: &entity.LoanPackageRequest{
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
			_, err = useCase.Create(
				context.Background(), entity.OfflineOfferUpdate{
					Id:       1,
					OfferId:  1,
					Status:   entity.OfflineOfferUpdateStatusProcessing,
					Category: "1",
					Note:     "note",
				}, entity.AssetTypeUnderlying,
			)
			assert.Nil(t, err)
		},
	)

	t.Run(
		"CreateLoanOfferUpdateProcessingFail", func(t *testing.T) {
			db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			if err != nil {
				t.Errorf("%v", err)
			}
			offlineOfferUpdateRepo := mock.NewMockOfflineOfferUpdatePersistenceRepository(t)
			offerRepo := mock.NewMockLoanPackageOfferRepository(t)
			offerInterestRepo := mock.NewMockLoanPackageOfferInterestRepository(t)
			useCase := NewUseCase(
				offlineOfferUpdateRepo, offerRepo, offerInterestRepo,
				&atomicity.DbAtomicExecutor{
					DB: db,
				},
			)
			sqlMock.ExpectBegin()
			sqlMock.ExpectRollback()
			offlineOfferUpdateRepo.On("Create", testifyMock.Anything, testifyMock.Anything).
				Return(
					entity.OfflineOfferUpdate{}, errors.New("error"),
				)
			offerRepo.On("InvestorGetById", testifyMock.Anything, testifyMock.Anything).
				Return(
					entity.LoanPackageOffer{
						Id:                   1,
						LoanPackageRequestId: 1,
						OfferedBy:            "admin",
						CreatedAt:            time.Now(),
						UpdatedAt:            time.Now(),
						ExpiredAt:            time.Now().Add(time.Hour * 24),
						FlowType:             entity.FLowTypeDnseOffline,
						LoanPackageRequest: &entity.LoanPackageRequest{
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
					},
					nil,
				)
			_, err = useCase.Create(
				context.Background(), entity.OfflineOfferUpdate{
					Id:       1,
					OfferId:  1,
					Status:   entity.OfflineOfferUpdateStatusProcessing,
					Category: "1",
					Note:     "note",
				}, entity.AssetTypeUnderlying,
			)
			assert.Equal(t, "offline offer useCase Create error", err.Error())
		},
	)

	t.Run(
		"CreateLoanOfferUpdateCancelSuccess", func(t *testing.T) {
			db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			if err != nil {
				t.Errorf("%v", err)
			}
			offlineOfferUpdateRepo := mock.NewMockOfflineOfferUpdatePersistenceRepository(t)
			offerRepo := mock.NewMockLoanPackageOfferRepository(t)
			offerInterestRepo := mock.NewMockLoanPackageOfferInterestRepository(t)
			useCase := NewUseCase(
				offlineOfferUpdateRepo, offerRepo, offerInterestRepo,
				&atomicity.DbAtomicExecutor{
					DB: db,
				},
			)
			sqlMock.ExpectBegin()
			sqlMock.ExpectCommit()
			offlineOfferUpdateRepo.On("Create", testifyMock.Anything, testifyMock.Anything).
				Return(
					entity.OfflineOfferUpdate{
						Id:       1,
						OfferId:  1,
						Status:   entity.OfflineOfferUpdateStatusRejected,
						Category: "1",
						Note:     "note",
					}, nil,
				)
			offerRepo.On("InvestorGetById", testifyMock.Anything, testifyMock.Anything).
				Return(
					entity.LoanPackageOffer{
						Id:                   1,
						LoanPackageRequestId: 1,
						OfferedBy:            "admin",
						CreatedAt:            time.Now(),
						UpdatedAt:            time.Now(),
						ExpiredAt:            time.Now().Add(time.Hour * 24),
						FlowType:             entity.FLowTypeDnseOffline,
						LoanPackageRequest: &entity.LoanPackageRequest{
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
					},
					nil,
				)
			offerInterestRepo.On(
				"CancelByOfferId", testifyMock.Anything, testifyMock.Anything, testifyMock.Anything,
				testifyMock.Anything,
			).
				Return(nil)
			_, err = useCase.Create(
				context.Background(), entity.OfflineOfferUpdate{
					Id:       1,
					OfferId:  1,
					Status:   entity.OfflineOfferUpdateStatusRejected,
					Category: "1",
					Note:     "note",
				}, entity.AssetTypeUnderlying,
			)
			assert.Nil(t, err)
		},
	)

	t.Run(
		"CreateLoanOfferUpdateCancelFail", func(t *testing.T) {
			db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			if err != nil {
				t.Errorf("%v", err)
			}
			offlineOfferUpdateRepo := mock.NewMockOfflineOfferUpdatePersistenceRepository(t)
			offerRepo := mock.NewMockLoanPackageOfferRepository(t)
			offerInterestRepo := mock.NewMockLoanPackageOfferInterestRepository(t)
			useCase := NewUseCase(
				offlineOfferUpdateRepo, offerRepo, offerInterestRepo,
				&atomicity.DbAtomicExecutor{
					DB: db,
				},
			)
			sqlMock.ExpectBegin()
			sqlMock.ExpectRollback()
			offlineOfferUpdateRepo.On("Create", testifyMock.Anything, testifyMock.Anything).
				Return(
					entity.OfflineOfferUpdate{
						Id:       1,
						OfferId:  1,
						Status:   entity.OfflineOfferUpdateStatusRejected,
						Category: "1",
						Note:     "note",
					}, nil,
				)
			offerRepo.On("InvestorGetById", testifyMock.Anything, testifyMock.Anything).
				Return(
					entity.LoanPackageOffer{
						Id:                   1,
						LoanPackageRequestId: 1,
						OfferedBy:            "admin",
						CreatedAt:            time.Now(),
						UpdatedAt:            time.Now(),
						ExpiredAt:            time.Now().Add(time.Hour * 24),
						FlowType:             entity.FLowTypeDnseOffline,
						LoanPackageRequest: &entity.LoanPackageRequest{
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
					},
					nil,
				)
			offerInterestRepo.On(
				"CancelByOfferId", testifyMock.Anything, testifyMock.Anything, testifyMock.Anything,
				testifyMock.Anything,
			).
				Return(errors.New("error"))
			_, err = useCase.Create(
				context.Background(), entity.OfflineOfferUpdate{
					Id:       1,
					OfferId:  1,
					Status:   entity.OfflineOfferUpdateStatusRejected,
					Category: "1",
					Note:     "note",
				}, entity.AssetTypeUnderlying,
			)
			assert.Equal(t, "offline offer useCase Create error", err.Error())
		},
	)
}
