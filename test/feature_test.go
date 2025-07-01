package test

import (
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/config"
	"financing-offer/internal/featureflag/transport/http"
	"financing-offer/internal/jwttoken"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/testhelper"
)

func TestFeatureHandler(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	cfg := do.MustInvoke[config.AppConfig](injector)
	cfg.Features = map[string]config.FeatureConfig{
		"loanRequest": {
			Enable:      false,
			InvestorIds: []string{"1", "2", "3"},
		},
	}
	do.OverrideValue[config.AppConfig](injector, cfg)
	h := do.MustInvoke[*http.FeatureHandler](injector)
	investorId := "1"
	t.Run(
		"verify feature then return true", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("GET", "/api/v1/features/verify", nil)
			ginCtx.AddParam("name", "loanRequest")
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					InvestorId: investorId,
					Sub:        investorId,
				},
			)
			h.CheckFeatureEnable(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, true, testhelper.GetBoolean(body, "data"))
		},
	)

	t.Run(
		"verify feature when feature not exist then return false", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("GET", "/api/v1/features/verify", nil)
			ginCtx.AddParam("name", "loanOffer")
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					InvestorId: investorId,
					Sub:        investorId,
				},
			)
			h.CheckFeatureEnable(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, false, testhelper.GetBoolean(body, "data"))
		},
	)
}
