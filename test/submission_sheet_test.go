package test

import (
	configRepo "financing-offer/internal/config/repository"
	loanPackageRequestRepo "financing-offer/internal/core/loanpackagerequest/repository"
	odooServiceRepo "financing-offer/internal/core/odoo_service/repository"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"github.com/gin-gonic/gin"
	"github.com/go-jet/jet/v2/postgres"
	mock2 "github.com/stretchr/testify/mock"
	gohttp "net/http"
	"strconv"
	"testing"
	"time"

	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core/entity"
	financialProductRepo "financing-offer/internal/core/financialproduct/repository"
	financingApiRepository "financing-offer/internal/core/financing/repository"
	marginOperationRepo "financing-offer/internal/core/marginoperation/repository"
	submissionSheetRepo "financing-offer/internal/core/submissionsheet/repository"
	"financing-offer/internal/core/submissionsheet/transport/http"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestSubmissionSheetHandler(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	do.OverrideValue[submissionSheetRepo.SubmissionSheetRepository](
		injector, &mock.MockSubmissionSheetRepository{},
	)
	financialProductRepository := &mock.MockFinancialProductRepository{}
	do.OverrideValue[financialProductRepo.FinancialProductRepository](
		injector, financialProductRepository,
	)
	marginOperationRepository := &mock.MockMarginOperationRepository{}
	do.OverrideValue[marginOperationRepo.MarginOperationRepository](
		injector, marginOperationRepository,
	)
	do.OverrideValue[financingApiRepository.FinancingRepository](injector, mock.FinancingApiMock{})
	h := do.MustInvoke[*http.SubmissionSheetHandler](injector)

	t.Run(
		"create new draft", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			loanId := int64(7055)
			financialProductRepository.EXPECT().GetLoanRateDetail(ginCtx, loanId).Return(
				entity.LoanRate{
					Id:                     loanId,
					Name:                   "ky quy 50%",
					InitialRate:            decimal.NewFromFloat(0.5),
					InitialRateForWithdraw: decimal.NewFromFloat(0.51),
					MaintenanceRate:        decimal.NewFromFloat(0.4),
					LiquidRate:             decimal.NewFromFloat(0.3),
				}, nil,
			)
			marginOperationRepository.EXPECT().GetMarginPoolsByIds(ginCtx, []int64{1001, 1000}).Return(

				[]entity.MarginPool{
					{
						Id:          1000,
						PoolGroupId: 1000,
						Name:        "Pool mặc định",
						Type:        "DEFAULT",
						Status:      "ACTIVE",
					},
					{
						Id:          1001,
						PoolGroupId: 1000,
						Name:        "Pool high margin",
						Type:        "RESERVE",
						Status:      "ACTIVE",
					},
				}, nil,
			)
			marginOperationRepository.EXPECT().GetMarginPoolGroupsByIds(ginCtx, []int64{1000, 1000}).Return(
				[]entity.MarginPoolGroup{
					{
						Id:     1000,
						Name:   "POOL TỔNG",
						Status: "ACTIVE",
						Type:   "MARGIN",
						Source: "DNSE",
					},
				}, nil,
			)

			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 0,
					MaxScore: 100,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "BID",
					AssetType:       "UNDERLYING",
				},
			)
			request := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       1,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test1",
					InterestRate:             decimal.NewFromFloat(0.12),
					InterestBasis:            31,
					Term:                     1,
					PoolIDRef:                1001,
					OverdueInterest:          decimal.NewFromFloat(0.66),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.43),
				},
			)
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       5,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test2",
					InterestRate:             decimal.NewFromFloat(0.53),
					InterestBasis:            30,
					Term:                     2,
					PoolIDRef:                1000,
					OverdueInterest:          decimal.NewFromFloat(0.12),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.442),
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/submission-sheets", entity.SubmissionSheetShorten{
					Metadata: entity.SubmissionSheetMetadata{
						ActionType:           "APPROVE",
						FlowType:             "002",
						ProposeType:          "NEW_LOAN_PACKAGE",
						LoanPackageRequestId: request.ID,
						Status:               "DRAFT",
					},
					Detail: entity.SubmissionSheetDetailShorten{
						LoanPackageRateId: 7055,
						FirmBuyingFee:     decimal.NewFromFloat(0.234),
						FirmSellingFee:    decimal.NewFromFloat(0.236),
						TransferFee:       decimal.NewFromFloat(234),
						LoanPolicies: []entity.LoanPolicyShorten{
							{
								LoanPolicyTemplateId:   1,
								InitialRateForWithdraw: decimal.NewFromFloat(0.4606),
								InitialRate:            decimal.NewFromFloat(0.47),
							},
							{
								LoanPolicyTemplateId:   5,
								InitialRateForWithdraw: decimal.NewFromFloat(0.0294),
								InitialRate:            decimal.NewFromFloat(0.03),
							},
						},
						Comment: "description",
					},
				},
			)
			h.Upsert(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "DRAFT", testhelper.GetString(body, "data", "metadata", "status"))
			assert.Equal(t, "test2", testhelper.GetString(body, "data", "detail", "loanPolicies", "[1]", "loanPolicyTemplateName"))
		},
	)

	t.Run(
		"update existed draft", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			loanId := int64(7055)
			financialProductRepository.EXPECT().GetLoanRateDetail(ginCtx, loanId).Return(
				entity.LoanRate{
					Id:                     loanId,
					Name:                   "ky quy 50%",
					InitialRate:            decimal.NewFromFloat(0.5),
					InitialRateForWithdraw: decimal.NewFromFloat(0.51),
					MaintenanceRate:        decimal.NewFromFloat(0.4),
					LiquidRate:             decimal.NewFromFloat(0.3),
				}, nil,
			)
			marginOperationRepository.EXPECT().GetMarginPoolsByIds(ginCtx, []int64{1001, 1000}).Return(

				[]entity.MarginPool{
					{
						Id:          1000,
						PoolGroupId: 1000,
						Name:        "Pool mặc định",
						Type:        "DEFAULT",
						Status:      "ACTIVE",
					},
					{
						Id:          1001,
						PoolGroupId: 1000,
						Name:        "Pool high margin",
						Type:        "RESERVE",
						Status:      "ACTIVE",
					},
				}, nil,
			)

			marginOperationRepository.EXPECT().GetMarginPoolGroupsByIds(ginCtx, []int64{1000, 1000}).Return(
				[]entity.MarginPoolGroup{
					{
						Id:     1000,
						Name:   "POOL TỔNG",
						Status: "ACTIVE",
						Type:   "MARGIN",
						Source: "DNSE",
					},
				}, nil,
			)
			//
			//submissionSheetMetadta:=mock.SeedSubmissionSheetMetadata(
			//	t,db,model.SubmissionSheetMetadata{})

			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 0,
					MaxScore: 100,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "BID",
					AssetType:       "UNDERLYING",
				},
			)
			request := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       1,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test1",
					InterestRate:             decimal.NewFromFloat(0.12),
					InterestBasis:            31,
					Term:                     1,
					PoolIDRef:                1001,
					OverdueInterest:          decimal.NewFromFloat(0.66),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.43),
				},
			)
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       5,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test2",
					InterestRate:             decimal.NewFromFloat(0.53),
					InterestBasis:            30,
					Term:                     2,
					PoolIDRef:                1000,
					OverdueInterest:          decimal.NewFromFloat(0.12),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.442),
				},
			)

			submissionSheetMetadata := mock.SeedSubmissionSheetMetadata(t, db, model.SubmissionSheetMetadata{
				ID:                   1,
				LoanPackageRequestID: request.ID,
				Status:               "DRAFT",
				ActionType:           "REJECT_AND_SEND_OTHER_PROPOSAL",
				FlowType:             "002",
				ProposeType:          "NEW_LOAN_PACKAGE",
			})
			submissionSheetDetail := mock.SeedSubmissionSheetDetail(t, db, model.SubmissionSheetDetail{
				ID:                1,
				SubmissionSheetID: submissionSheetMetadata.ID,
				LoanRate:          "{\"id\": 7071, \"name\": \"ky quy 30%\", \"liquidRate\": 0.2, \"initialRate\": 0.3, \"interestRate\": 0, \"maintenanceRate\": 0.25, \"initialRateForWithdraw\": 0.31}",
				FirmBuyingFee:     decimal.NewFromFloat(0.234),
				FirmSellingFee:    decimal.NewFromFloat(0.236),
				TransferFee:       decimal.NewFromFloat(234),
				LoanPolicies:      "[\n  {\n    \"loanPolicyTemplateName\": \"truong.quach\",\n    \"term\": 3,\n    \"source\": \"C\",\n    \"createdAt\": \"2024-03-16T16:26:23.875309Z\",\n    \"poolIdRef\": 1001,\n    \"updatedAt\": \"2024-03-16T16:26:23.875309Z\",\n    \"updatedBy\": \"dnse.admin@dnse.com.vn\",\n    \"initialRate\": 0.7,\n    \"interestRate\": 0.245,\n    \"interestBasis\": 24,\n    \"overdueInterest\": 0.435,\n    \"allowEarlyPayment\": true,\n    \"preferentialPeriod\": 365,\n    \"allowExtendLoanTerm\": true,\n    \"loanPolicyTemplateId\": 1,\n    \"initialRateForWithdraw\": 0.69,\n    \"preferentialInterestRate\": 0.325\n  }\n]",
				Comment:           "test lan 1",
			})
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/submission-sheets", entity.SubmissionSheetShorten{
					Metadata: entity.SubmissionSheetMetadata{
						Id:                   submissionSheetMetadata.ID,
						ActionType:           "APPROVE",
						FlowType:             "002",
						ProposeType:          "NEW_LOAN_PACKAGE",
						LoanPackageRequestId: request.ID,
						Status:               "DRAFT",
					},
					Detail: entity.SubmissionSheetDetailShorten{
						Id:                submissionSheetDetail.ID,
						SubmissionSheetId: submissionSheetMetadata.ID,
						LoanPackageRateId: 7055,
						FirmBuyingFee:     decimal.NewFromFloat(0.234),
						FirmSellingFee:    decimal.NewFromFloat(0.236),
						TransferFee:       decimal.NewFromFloat(234),
						LoanPolicies: []entity.LoanPolicyShorten{
							{
								LoanPolicyTemplateId:   1,
								InitialRateForWithdraw: decimal.NewFromFloat(0.4606),
								InitialRate:            decimal.NewFromFloat(0.47),
							},
							{
								LoanPolicyTemplateId:   5,
								InitialRateForWithdraw: decimal.NewFromFloat(0.0294),
								InitialRate:            decimal.NewFromFloat(0.03),
							},
						},
						Comment: "description",
					},
				},
			)
			h.Upsert(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "DRAFT", testhelper.GetString(body, "data", "metadata", "status"))
			assert.Equal(t, "test2", testhelper.GetString(body, "data", "detail", "loanPolicies", "[1]", "loanPolicyTemplateName"))
		},
	)
}

func TestLoanPackageRequestHandler_AdminApproveSubmissionSheet(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))

	do.OverrideValue[loanPackageRequestRepo.LoanPackageRequestEventRepository](
		injector, &mock.LoanPackageRequestEventRepository{},
	)
	financialProductRepository := &mock.MockFinancialProductRepository{}
	do.OverrideValue[financialProductRepo.FinancialProductRepository](
		injector, financialProductRepository,
	)
	configurationRepoMock := mock.NewMockConfigurationPersistenceRepository(t)
	marginOperationRepository := &mock.MockMarginOperationRepository{}
	odooServiceRepository := &mock.MockOdooServiceRepository{}
	do.OverrideValue[marginOperationRepo.MarginOperationRepository](
		injector, marginOperationRepository,
	)
	do.OverrideValue[financingApiRepository.FinancingRepository](injector, mock.FinancingApiMock{})
	do.OverrideValue[configRepo.ConfigurationPersistenceRepository](injector, configurationRepoMock)
	do.OverrideValue[odooServiceRepo.OdooServiceRepository](injector, odooServiceRepository)
	handler := do.MustInvoke[*http.SubmissionSheetHandler](injector)

	t.Run(
		"admin approve submission", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 0,
					MaxScore: 100,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "BID",
					AssetType:       "UNDERLYING",
				},
			)

			request := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)

			mock.SeedSubmissionSheetMetadata(
				t, db, model.SubmissionSheetMetadata{
					ID:                   1,
					LoanPackageRequestID: request.ID,
					Status:               "SUBMITTED",
					ActionType:           "REJECT_AND_SEND_OTHER_PROPOSAL",
					FlowType:             "002",
					ProposeType:          "NEW_LOAN_PACKAGE",
					CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
				},
			)

			submissionSheetMetadata2 := mock.SeedSubmissionSheetMetadata(
				t, db, model.SubmissionSheetMetadata{
					ID:                   2,
					LoanPackageRequestID: request.ID,
					Status:               "SUBMITTED",
					ActionType:           "REJECT_AND_SEND_OTHER_PROPOSAL",
					FlowType:             "002",
					ProposeType:          "NEW_LOAN_PACKAGE",
					CreatedAt:            time.Date(2021, 2, 1, 0, 0, 0, 0, time.Local),
				},
			)
			mock.SeedSubmissionSheetDetail(
				t, db, model.SubmissionSheetDetail{
					ID:                1,
					SubmissionSheetID: submissionSheetMetadata2.ID,
					LoanRate:          "{\"id\": 7071, \"name\": \"ky quy 30%\", \"liquidRate\": 0.2, \"initialRate\": 0.3, \"interestRate\": 0, \"maintenanceRate\": 0.25, \"initialRateForWithdraw\": 0.31}",
					FirmBuyingFee:     decimal.NewFromFloat(0.234),
					FirmSellingFee:    decimal.NewFromFloat(0.236),
					TransferFee:       decimal.NewFromFloat(234),
					LoanPolicies:      "[\n  {\n    \"loanPolicyTemplateName\": \"truong.quach\",\n    \"term\": 3,\n    \"source\": \"C\",\n    \"createdAt\": \"2024-03-16T16:26:23.875309Z\",\n    \"poolIdRef\": 1001,\n    \"updatedAt\": \"2024-03-16T16:26:23.875309Z\",\n    \"updatedBy\": \"dnse.admin@dnse.com.vn\",\n    \"initialRate\": 0.7,\n    \"interestRate\": 0.245,\n    \"interestBasis\": 24,\n    \"overdueInterest\": 0.435,\n    \"allowEarlyPayment\": true,\n    \"preferentialPeriod\": 365,\n    \"allowExtendLoanTerm\": true,\n    \"loanPolicyTemplateId\": 1,\n    \"initialRateForWithdraw\": 0.69,\n    \"preferentialInterestRate\": 0.325\n  }\n]",
					Comment:           "test lan 1",
				},
			)
			financialProductRepository.EXPECT().GetAllAccountDetail(mock2.Anything, request.InvestorID).Return(
				[]entity.FinancialAccountDetail{
					{
						AccountNo: request.AccountNo,
					},
				}, nil,
			)

			ginCtx.Request = gintest.MustMakeRequest("POST", "/:id/approve", nil)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(submissionSheetMetadata2.ID, 10),
			}}
			handler.AdminApproveSubmissionSheet(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "success", testhelper.GetString(body, "data"))

			var submissionMetadata model.SubmissionSheetMetadata
			err := table.SubmissionSheetMetadata.SELECT(table.SubmissionSheetMetadata.AllColumns).WHERE(table.SubmissionSheetMetadata.ID.EQ(postgres.Int64(submissionSheetMetadata2.ID))).Query(db, &submissionMetadata)
			assert.Nil(t, err)
			assert.Equal(t, "APPROVED", submissionMetadata.Status)
		})

	t.Run(
		"admin approve submission error", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("POST", "/:id/approve", nil)
			handler.AdminApproveSubmissionSheet(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			assert.Equal(t, gohttp.StatusBadRequest, result.StatusCode)
		},
	)
}

func TestLoanPackageRequestHandler_AdminRejectSubmissionSheet(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))

	do.OverrideValue[loanPackageRequestRepo.LoanPackageRequestEventRepository](
		injector, &mock.LoanPackageRequestEventRepository{},
	)
	financialProductRepository := &mock.MockFinancialProductRepository{}
	do.OverrideValue[financialProductRepo.FinancialProductRepository](
		injector, financialProductRepository,
	)
	configurationRepoMock := mock.NewMockConfigurationPersistenceRepository(t)
	marginOperationRepository := &mock.MockMarginOperationRepository{}
	odooServiceRepository := &mock.MockOdooServiceRepository{}
	do.OverrideValue[marginOperationRepo.MarginOperationRepository](
		injector, marginOperationRepository,
	)
	do.OverrideValue[financingApiRepository.FinancingRepository](injector, mock.FinancingApiMock{})
	do.OverrideValue[configRepo.ConfigurationPersistenceRepository](injector, configurationRepoMock)
	do.OverrideValue[odooServiceRepo.OdooServiceRepository](injector, odooServiceRepository)
	handler := do.MustInvoke[*http.SubmissionSheetHandler](injector)

	t.Run(
		"admin reject submission", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 0,
					MaxScore: 100,
				},
			)
			symbol := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "BID",
					AssetType:       "UNDERLYING",
				},
			)

			request := mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:   symbol.ID,
					InvestorID: "0001000115",
					AccountNo:  "0001000115",
					LoanRate:   decimal.NewFromFloat(0.7),
					Type:       "FLEXIBLE",
					Status:     "PENDING",
					AssetType:  "UNDERLYING",
				},
			)

			mock.SeedSubmissionSheetMetadata(
				t, db, model.SubmissionSheetMetadata{
					ID:                   1,
					LoanPackageRequestID: request.ID,
					Status:               "REJECTED",
					ActionType:           "REJECT_AND_SEND_OTHER_PROPOSAL",
					FlowType:             "002",
					ProposeType:          "NEW_LOAN_PACKAGE",
					CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
				},
			)

			submissionSheetMetadata2 := mock.SeedSubmissionSheetMetadata(
				t, db, model.SubmissionSheetMetadata{
					ID:                   2,
					LoanPackageRequestID: request.ID,
					Status:               "SUBMITTED",
					ActionType:           "REJECT_AND_SEND_OTHER_PROPOSAL",
					FlowType:             "002",
					ProposeType:          "NEW_LOAN_PACKAGE",
					CreatedAt:            time.Date(2021, 2, 1, 0, 0, 0, 0, time.Local),
				},
			)
			mock.SeedSubmissionSheetDetail(
				t, db, model.SubmissionSheetDetail{
					ID:                1,
					SubmissionSheetID: submissionSheetMetadata2.ID,
					LoanRate:          "{\"id\": 7071, \"name\": \"ky quy 30%\", \"liquidRate\": 0.2, \"initialRate\": 0.3, \"interestRate\": 0, \"maintenanceRate\": 0.25, \"initialRateForWithdraw\": 0.31}",
					FirmBuyingFee:     decimal.NewFromFloat(0.234),
					FirmSellingFee:    decimal.NewFromFloat(0.236),
					TransferFee:       decimal.NewFromFloat(234),
					LoanPolicies:      "[\n  {\n    \"loanPolicyTemplateName\": \"truong.quach\",\n    \"term\": 3,\n    \"source\": \"C\",\n    \"createdAt\": \"2024-03-16T16:26:23.875309Z\",\n    \"poolIdRef\": 1001,\n    \"updatedAt\": \"2024-03-16T16:26:23.875309Z\",\n    \"updatedBy\": \"dnse.admin@dnse.com.vn\",\n    \"initialRate\": 0.7,\n    \"interestRate\": 0.245,\n    \"interestBasis\": 24,\n    \"overdueInterest\": 0.435,\n    \"allowEarlyPayment\": true,\n    \"preferentialPeriod\": 365,\n    \"allowExtendLoanTerm\": true,\n    \"loanPolicyTemplateId\": 1,\n    \"initialRateForWithdraw\": 0.69,\n    \"preferentialInterestRate\": 0.325\n  }\n]",
					Comment:           "test lan 1",
				},
			)
			financialProductRepository.EXPECT().GetAllAccountDetail(mock2.Anything, request.InvestorID).Return(
				[]entity.FinancialAccountDetail{
					{
						AccountNo: request.AccountNo,
					},
				}, nil,
			)

			ginCtx.Request = gintest.MustMakeRequest("POST", "/:id/reject", nil)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(submissionSheetMetadata2.ID, 10),
			}}
			handler.AdminRejectSubmissionSheet(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "success", testhelper.GetString(body, "data"))

			var submissionMetadata model.SubmissionSheetMetadata
			err := table.SubmissionSheetMetadata.SELECT(table.SubmissionSheetMetadata.AllColumns).WHERE(table.SubmissionSheetMetadata.ID.EQ(postgres.Int64(submissionSheetMetadata2.ID))).Query(db, &submissionMetadata)
			assert.Nil(t, err)
			assert.Equal(t, "REJECTED", submissionMetadata.Status)
		})

	t.Run(
		"admin reject submission error", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("POST", "/:id/reject", nil)
			handler.AdminRejectSubmissionSheet(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			assert.Equal(t, gohttp.StatusBadRequest, result.StatusCode)
		},
	)
}
