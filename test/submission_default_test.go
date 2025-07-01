package test

import (
	"encoding/json"
	"financing-offer/internal/appcontext"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/submission_default/transport/http"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/jwttoken"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSubmissionDefaultHandler_GetSubmissionDefault(t *testing.T) {
	t.Parallel()
	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	submissionDefaultHandler := do.MustInvoke[*http.SubmissionDefaultHandler](injector)

	t.Run(
		"test get submission default success", func(t *testing.T) {
			defer truncateData()
			bytes, err := json.Marshal(
				entity.SubmissionDefault{
					FirmBuyingFeeRate:        0.1,
					FirmSellingFeeRate:       0.2,
					TransferFee:              0.5,
					AllowedOverdueLoanInDays: 3,
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "submissionDefault",
					Value:         string(bytes),
					LastUpdatedBy: "admin",
				},
			)
			ctx, _, recorder := gintest.GetTestContext()
			ctx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "admin",
				},
			)
			submissionDefaultHandler.GetSubmissionDefault(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 0.1, testhelper.GetFloat(body, "data", "firmBuyingFeeRate"))
			assert.Equal(t, 0.2, testhelper.GetFloat(body, "data", "firmSellingFeeRate"))
			assert.Equal(t, 0.5, testhelper.GetFloat(body, "data", "transferFee"))
			assert.Equal(t, int64(3), testhelper.GetInt(body, "data", "allowedOverdueLoanInDays"))
		},
	)
}

func TestNewSubmissionDefaultHandler_SetSubmissionDefault(t *testing.T) {
	t.Parallel()
	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	submissionDefaultHandler := do.MustInvoke[*http.SubmissionDefaultHandler](injector)

	t.Run(
		"test set submission default success", func(t *testing.T) {
			defer truncateData()
			bytes, err := json.Marshal(
				entity.SubmissionDefault{
					FirmBuyingFeeRate:        0.1,
					FirmSellingFeeRate:       0.2,
					TransferFee:              0.5,
					AllowedOverdueLoanInDays: 3,
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "submissionDefault",
					Value:         string(bytes),
					LastUpdatedBy: "admin",
				},
			)
			ctx, _, recorder := gintest.GetTestContext()
			ctx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "admin",
				},
			)
			ctx.Request = gintest.MustMakeRequest(
				"POST", "/v1/submission-default",
				entity.SubmissionDefault{
					FirmBuyingFeeRate:        0.8,
					FirmSellingFeeRate:       0.8,
					TransferFee:              0.8,
					AllowedOverdueLoanInDays: 8,
				},
			)
			submissionDefaultHandler.SetSubmissionDefault(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 0.8, testhelper.GetFloat(body, "data", "firmBuyingFeeRate"))
			assert.Equal(t, 0.8, testhelper.GetFloat(body, "data", "firmSellingFeeRate"))
			assert.Equal(t, 0.8, testhelper.GetFloat(body, "data", "transferFee"))
			assert.Equal(t, int64(8), testhelper.GetInt(body, "data", "allowedOverdueLoanInDays"))
		},
	)

}
