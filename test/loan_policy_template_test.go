package test

import (
	"financing-offer/internal/core/entity"
	marginOperationRepo "financing-offer/internal/core/marginoperation/repository"
	"github.com/gin-gonic/gin"
	gohttp "net/http"
	"testing"
	"time"

	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	loanPolicyRepo "financing-offer/internal/core/loanpolicytemplate/repository"
	"financing-offer/internal/core/loanpolicytemplate/transport/http"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestLoanPolicyTemplateHandler(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	do.OverrideValue[loanPolicyRepo.LoanPolicyTemplateRepository](
		injector, &mock.MockLoanPolicyTemplateRepository{},
	)
	marginOperationRepository := &mock.MockMarginOperationRepository{}
	do.OverrideValue[marginOperationRepo.MarginOperationRepository](
		injector, marginOperationRepository,
	)
	h := do.MustInvoke[*http.LoanPolicyTemplateHandler](injector)

	t.Run(
		"create success", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			marginOperationRepository.EXPECT().GetMarginPoolById(ginCtx, int64(1000)).Return(

				entity.MarginPool{

					Id:          1000,
					PoolGroupId: 1000,
					Name:        "Pool mặc định",
					Type:        "DEFAULT",
					Status:      "ACTIVE",
				}, nil,
			)
			marginOperationRepository.EXPECT().GetMarginPoolGroupsByIds(ginCtx, []int64{1000}).Return(
				[]entity.MarginPoolGroup{
					{
						Id:     1000,
						Name:   "POOL TỔNG",
						Status: "ACTIVE",
						Type:   "MARGIN",
						Source: "SOURCE_TEST",
					},
				}, nil,
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/loan-policy-template", http.CreateLoanPackageTemplateRequest{
					Name:                     "test",
					InterestRate:             decimal.NewFromFloat(0.12),
					InterestBasis:            1,
					Term:                     1,
					PoolIdRef:                1000,
					OverdueInterest:          decimal.NewFromFloat(0.12),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.12),
				},
			)
			h.Create(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusCreated, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "test", testhelper.GetString(body, "data", "name"))
		},
	)

	t.Run(
		"get all success", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       1,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test1",
					InterestRate:             decimal.NewFromFloat(0.12),
					InterestBasis:            1,
					Term:                     1,
					PoolIDRef:                1000,
					OverdueInterest:          decimal.NewFromFloat(0.12),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.12),
				},
			)
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       5,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test2",
					InterestRate:             decimal.NewFromFloat(0.12),
					InterestBasis:            1,
					Term:                     1,
					PoolIDRef:                1001,
					OverdueInterest:          decimal.NewFromFloat(0.12),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.12),
				},
			)
			marginOperationRepository.EXPECT().GetMarginPoolsByIds(ginCtx, []int64{1000, 1001}).Return(

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
						Source: "SOURCE_TEST",
					},
				}, nil,
			)

			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/loan-policy-template", nil,
			)
			h.GetAll(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "test1", testhelper.GetString(body, "data", "[0]", "name"))
			assert.Equal(t, "test2", testhelper.GetString(body, "data", "[1]", "name"))
			assert.Equal(t, "SOURCE_TEST", testhelper.GetString(body, "data", "[0]", "source"))
		},
	)

	t.Run(
		"update success", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       1,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test",
					InterestRate:             decimal.NewFromFloat(0.12),
					InterestBasis:            1,
					Term:                     1,
					PoolIDRef:                1000,
					OverdueInterest:          decimal.NewFromFloat(0.12),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.12),
				},
			)
			marginOperationRepository.EXPECT().GetMarginPoolById(ginCtx, int64(1000)).Return(

				entity.MarginPool{

					Id:          1000,
					PoolGroupId: 1000,
					Name:        "Pool mặc định",
					Type:        "DEFAULT",
					Status:      "ACTIVE",
				}, nil,
			)
			marginOperationRepository.EXPECT().GetMarginPoolGroupsByIds(ginCtx, []int64{1000}).Return(
				[]entity.MarginPoolGroup{
					{
						Id:     1000,
						Name:   "POOL TỔNG",
						Status: "ACTIVE",
						Type:   "MARGIN",
						Source: "SOURCE_TEST",
					},
				}, nil,
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"PUT", "/api/v1/loan-policy-template/1", http.UpdateLoanPackageTemplateRequest{
					Id:                       1,
					Name:                     "testNew",
					InterestRate:             decimal.NewFromFloat(0.12),
					InterestBasis:            1,
					Term:                     1,
					PoolIdRef:                1000,
					OverdueInterest:          decimal.NewFromFloat(0.12),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.12),
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: "1",
			}}
			h.Update(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "testNew", testhelper.GetString(body, "data", "name"))
		},
	)

	t.Run(
		"delete success", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       1,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test",
					InterestRate:             decimal.NewFromFloat(0.12),
					InterestBasis:            1,
					Term:                     1,
					PoolIDRef:                1000,
					OverdueInterest:          decimal.NewFromFloat(0.12),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.12),
				},
			)

			ginCtx.Request = gintest.MustMakeRequest(
				"DELETE", "/api/v1/loan-policy-template/1", nil,
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: "1",
			}}
			h.Delete(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusNoContent, result.StatusCode)
		},
	)

	t.Run(
		"getById success", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			mock.SeedLoanPolicyTemplate(
				t, db, model.LoanPolicyTemplate{
					ID:                       1,
					CreatedAt:                time.Now(),
					UpdatedAt:                time.Now(),
					UpdatedBy:                "test",
					Name:                     "test",
					InterestRate:             decimal.NewFromFloat(0.12),
					InterestBasis:            1,
					Term:                     1,
					PoolIDRef:                1000,
					OverdueInterest:          decimal.NewFromFloat(0.12),
					AllowExtendLoanTerm:      true,
					AllowEarlyPayment:        true,
					PreferentialPeriod:       1,
					PreferentialInterestRate: decimal.NewFromFloat(0.12),
				},
			)

			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/loan-policy-template/1", nil,
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: "1",
			}}
			h.GetById(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, gohttp.StatusOK, result.StatusCode)
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "test", testhelper.GetString(body, "data", "name"))
		},
	)

}
