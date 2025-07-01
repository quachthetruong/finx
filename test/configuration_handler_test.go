package test

import (
	"encoding/json"
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/core/configuration/transport/http"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/jwttoken"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestConfigurationForAdminHandler(t *testing.T) {
	t.Parallel()
	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	do.OverrideValue(injector, &mock.EmptyCache{})
	configurationHandler := do.MustInvoke[*http.ConfigurationHandler](injector)

	t.Run(
		"test set loan rate configuration and success", func(t *testing.T) {
			defer truncateData()
			bytes, err := json.Marshal(
				entity.LoanRateConfiguration{
					Ids: []int64{1, 2},
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "loanRate",
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
				"POST", "/v1/configurations/loan-rate", entity.LoanRateConfiguration{
					Ids: []int64{1, 2},
				},
			)
			configurationHandler.SetLoanRate(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(
				t, []string{"1", "2"},
				testhelper.GetArrayString(body, "data", "ids"),
			)
		},
	)

	t.Run(
		"test get loan rate configuration and success", func(t *testing.T) {
			defer truncateData()
			bytes, err := json.Marshal(
				entity.LoanRateConfiguration{
					Ids: []int64{1, 2},
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "loanRate",
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
			configurationHandler.GetLoanRate(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, int64(1), testhelper.GetInt(body, "data", "ids", "[0]"))
			assert.Equal(t, int64(2), testhelper.GetInt(body, "data", "ids", "[1]"))
		},
	)

	t.Run(
		"test set margin pool configuration and success", func(t *testing.T) {
			defer truncateData()
			bytes, err := json.Marshal(
				entity.MarginPoolConfiguration{
					Ids: []int64{1, 2},
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "marginPool",
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
				"POST", "/v1/configurations/margin-pool", entity.MarginPoolConfiguration{
					Ids: []int64{1, 2},
				},
			)
			configurationHandler.SetMarginPool(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(
				t, []string{"1", "2"},
				testhelper.GetArrayString(body, "data", "ids"),
			)
		},
	)

	t.Run(
		"test get margin pool configuration and success", func(t *testing.T) {
			defer truncateData()
			bytes, err := json.Marshal(
				entity.MarginPoolConfiguration{
					Ids: []int64{1, 2},
				},
			)
			require.NoError(t, err)
			mock.SeedConfiguration(
				t, db, model.FinancialConfiguration{
					Attribute:     "marginPool",
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
			configurationHandler.GetMarginPool(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, int64(1), testhelper.GetInt(body, "data", "ids", "[0]"))
			assert.Equal(t, int64(2), testhelper.GetInt(body, "data", "ids", "[1]"))
		},
	)
}
