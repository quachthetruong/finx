package test

import (
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/core/suggested_offer_config/transport/http"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/jwttoken"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestSuggestedOfferConfigHandler(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	h := do.MustInvoke[*http.SuggestedOfferConfigHandler](injector)
	t.Run(
		"create suggested offer config success", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "test@dnse.com",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/suggested-offer-configs", http.CreateSuggestedOfferConfigRequest{
					Name:      "Test 1",
					Value:     decimal.NewFromFloat(6.99),
					ValueType: "LOAN_RATE",
				},
			)
			h.Create(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "Test 1", testhelper.GetString(body, "data", "name"))
			assert.Equal(t, 6.99, testhelper.GetFloat(body, "data", "value"))
			assert.Equal(t, "LOAN_RATE", testhelper.GetString(body, "data", "valueType"))
			assert.Equal(t, "INACTIVE", testhelper.GetString(body, "data", "status"))
		},
	)

	t.Run(
		"create suggested offer config invalid payload", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("POST", "/api/v1/suggested-offer-configs", "")
			h.Create(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "invalid payload", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"create suggested offer config error", func(t *testing.T) {
			defer func() {
				db, tearDownDb, truncateData = dbtest.NewDb(t)
				do.OverrideValue[*sql.DB](injector, db)
				truncateData()
				injector := testhelper.NewInjector(testhelper.WithDb(db))
				h = do.MustInvoke[*http.SuggestedOfferConfigHandler](injector)
			}()
			if err := db.Close(); err != nil {
				t.Fatal(err)
			}
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "test@dnse.com",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/suggested-offer-configs", http.CreateSuggestedOfferConfigRequest{
					Name:      "Test 1",
					Value:     decimal.NewFromFloat(6.99),
					ValueType: "LOAN_RATE",
				},
			)
			h.Create(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "an error happened, please try again later", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"update suggested offer config success", func(t *testing.T) {
			defer truncateData()
			suggestedOfferConfig := mock.SeedSuggestedOfferConfig(
				t, db, model.SuggestedOfferConfig{
					Name:          "Test",
					Value:         decimal.NewFromFloat(0.99),
					ValueType:     "INTEREST_RATE",
					Status:        "ACTIVE",
					CreatedBy:     "admin",
					LastUpdatedBy: "admin",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "test@dnse.com",
				},
			)
			ginCtx.AddParam("id", strconv.FormatInt(suggestedOfferConfig.ID, 10))
			ginCtx.Request = gintest.MustMakeRequest(
				"PATCH", "/api/v1/suggested-offer-configs/1", http.CreateSuggestedOfferConfigRequest{
					Name:      "Test 2",
					Value:     decimal.NewFromFloat(0.88),
					ValueType: "LOAN_RATE",
				},
			)
			h.Update(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "Test 2", testhelper.GetString(body, "data", "name"))
			assert.Equal(t, 0.88, testhelper.GetFloat(body, "data", "value"))
			assert.Equal(t, "LOAN_RATE", testhelper.GetString(body, "data", "valueType"))
			assert.Equal(t, "ACTIVE", testhelper.GetString(body, "data", "status"))
		},
	)

	t.Run(
		"update suggested offer config id invalid", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("PATCH", "/api/v1/suggested-offer-configs/1", "")
			h.Update(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "id invalid", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"update suggested offer config invalid payload", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			suggestedOfferConfig := mock.SeedSuggestedOfferConfig(
				t, db, model.SuggestedOfferConfig{
					Name:          "Test",
					Value:         decimal.NewFromFloat(0.99),
					ValueType:     "INTEREST_RATE",
					Status:        "ACTIVE",
					CreatedBy:     "admin",
					LastUpdatedBy: "admin",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				},
			)
			ginCtx.AddParam("id", strconv.FormatInt(suggestedOfferConfig.ID, 10))
			ginCtx.Request = gintest.MustMakeRequest("PATCH", "/api/v1/suggested-offer-configs/1", "")
			h.Update(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "invalid payload", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"get suggested offer config success", func(t *testing.T) {
			defer truncateData()
			suggestedOfferConfig := mock.SeedSuggestedOfferConfig(
				t, db, model.SuggestedOfferConfig{
					Name:      "Test",
					Value:     decimal.NewFromFloat(0.99),
					ValueType: "INTEREST_RATE",
					Status:    "ACTIVE",
					CreatedBy: "admin",
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/active-suggested-offer-config", nil,
			)
			h.Get(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, suggestedOfferConfig.Name, testhelper.GetString(body, "data", "name"))
			assert.Equal(t, suggestedOfferConfig.Value.InexactFloat64(), testhelper.GetFloat(body, "data", "value"))
			assert.Equal(t, suggestedOfferConfig.ValueType, testhelper.GetString(body, "data", "valueType"))
			assert.Equal(t, suggestedOfferConfig.Status, testhelper.GetString(body, "data", "status"))
		},
	)

	t.Run(
		"get suggested offer config error", func(t *testing.T) {
			defer func() {
				db, tearDownDb, truncateData = dbtest.NewDb(t)
				do.OverrideValue[*sql.DB](injector, db)
				truncateData()
				injector := testhelper.NewInjector(testhelper.WithDb(db))
				h = do.MustInvoke[*http.SuggestedOfferConfigHandler](injector)
			}()

			if err := db.Close(); err != nil {
				t.Fatal(err)
			}
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/active-suggested-offer-config", nil,
			)
			h.Get(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "an error happened, please try again later", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"get suggested offer config not found", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/active-suggested-offer-config", nil,
			)
			h.Get(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Empty(t, testhelper.GetString(body, "data", "name"))
		},
	)

	t.Run(
		"get all suggested offer config success", func(t *testing.T) {
			defer truncateData()
			suggestedOfferConfig := mock.SeedSuggestedOfferConfig(
				t, db, model.SuggestedOfferConfig{
					Name:      "Test",
					Value:     decimal.NewFromFloat(0.99),
					ValueType: "INTEREST_RATE",
					Status:    "ACTIVE",
					CreatedBy: "admin",
				},
			)
			suggestedOfferConfig2 := mock.SeedSuggestedOfferConfig(
				t, db, model.SuggestedOfferConfig{
					Name:      "Test 2",
					Value:     decimal.NewFromFloat(0.99),
					ValueType: "LOAN_RATE",
					Status:    "INACTIVE",
					CreatedBy: "admin",
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/suggested-offer-configs", nil,
			)
			h.GetAll(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, suggestedOfferConfig.Name, testhelper.GetString(body, "data", "[0]", "name"))
			assert.Equal(t, suggestedOfferConfig2.Name, testhelper.GetString(body, "data", "[1]", "name"))
		},
	)

	t.Run(
		"get all suggested offer config error", func(t *testing.T) {
			defer func() {
				db, tearDownDb, truncateData = dbtest.NewDb(t)
				do.OverrideValue[*sql.DB](injector, db)
				truncateData()
				injector := testhelper.NewInjector(testhelper.WithDb(db))
				h = do.MustInvoke[*http.SuggestedOfferConfigHandler](injector)
			}()

			if err := db.Close(); err != nil {
				t.Fatal(err)
			}
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/suggested-offer-configs", nil,
			)
			h.GetAll(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "an error happened, please try again later", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"get by id suggested offer config success", func(t *testing.T) {
			defer truncateData()
			suggestedOfferConfig := mock.SeedSuggestedOfferConfig(
				t, db, model.SuggestedOfferConfig{
					Name:      "Test",
					Value:     decimal.NewFromFloat(0.99),
					ValueType: "INTEREST_RATE",
					Status:    "ACTIVE",
					CreatedBy: "admin",
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.AddParam("id", strconv.FormatInt(suggestedOfferConfig.ID, 10))
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/suggested-offer-configs/1", nil,
			)
			h.GetById(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, suggestedOfferConfig.Name, testhelper.GetString(body, "data", "name"))
			assert.Equal(t, suggestedOfferConfig.Value.InexactFloat64(), testhelper.GetFloat(body, "data", "value"))
			assert.Equal(t, suggestedOfferConfig.ValueType, testhelper.GetString(body, "data", "valueType"))
			assert.Equal(t, suggestedOfferConfig.Status, testhelper.GetString(body, "data", "status"))
		},
	)

	t.Run(
		"get by id suggested offer config id invalid", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("GET", "/api/v1/suggested-offer-configs/1", nil)
			h.GetById(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "id invalid", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"get by id suggested offer config error not found", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.AddParam("id", "1")
			ginCtx.Request = gintest.MustMakeRequest("GET", "/api/v1/suggested-offer-configs/1", nil)
			h.GetById(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "not found resources", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"update suggested offer config status active success", func(t *testing.T) {
			defer truncateData()
			suggestedOfferConfig := mock.SeedSuggestedOfferConfig(
				t, db, model.SuggestedOfferConfig{
					Name:      "Test 2",
					Value:     decimal.NewFromFloat(0.59),
					ValueType: "INTEREST_RATE",
					Status:    "INACTIVE",
					CreatedBy: "test",
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.AddParam("id", strconv.FormatInt(suggestedOfferConfig.ID, 10))
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "admin@dnse.com",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"PATCH", "/api/v1/suggested-offer-configs/1/status", http.UpdateSuggestedOfferConfigStatusRequest{
					Status: "ACTIVE",
				},
			)
			h.UpdateStatus(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "ACTIVE", testhelper.GetString(body, "data", "status"))
			assert.Equal(t, "Test 2", testhelper.GetString(body, "data", "name"))
			assert.Equal(t, "admin@dnse.com", testhelper.GetString(body, "data", "lastUpdatedBy"))
		},
	)

	t.Run(
		"update suggested offer config status inactive success", func(t *testing.T) {
			defer truncateData()
			suggestedOfferConfig := mock.SeedSuggestedOfferConfig(
				t, db, model.SuggestedOfferConfig{
					Name:      "Test 2",
					Value:     decimal.NewFromFloat(0.59),
					ValueType: "INTEREST_RATE",
					Status:    "ACTIVE",
					CreatedBy: "test",
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.AddParam("id", strconv.FormatInt(suggestedOfferConfig.ID, 10))
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "admin@dnse.com",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"PATCH", "/api/v1/suggested-offer-configs/1/status", http.UpdateSuggestedOfferConfigStatusRequest{
					Status: "INACTIVE",
				},
			)
			h.UpdateStatus(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "INACTIVE", testhelper.GetString(body, "data", "status"))
			assert.Equal(t, "Test 2", testhelper.GetString(body, "data", "name"))
			assert.Equal(t, "admin@dnse.com", testhelper.GetString(body, "data", "lastUpdatedBy"))
		},
	)

	t.Run(
		"update suggested offer config status id invalid", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("PATCH", "/api/v1/suggested-offer-configs/1/status", "")
			h.UpdateStatus(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "id invalid", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"update suggested offer config status invalid payload", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			suggestedOfferConfig := mock.SeedSuggestedOfferConfig(
				t, db, model.SuggestedOfferConfig{
					Name:      "Test 2",
					Value:     decimal.NewFromFloat(0.59),
					ValueType: "INTEREST_RATE",
					Status:    "ACTIVE",
					CreatedBy: "admin",
				},
			)
			ginCtx.AddParam("id", strconv.FormatInt(suggestedOfferConfig.ID, 10))
			ginCtx.Request = gintest.MustMakeRequest("PATCH", "/api/v1/suggested-offer-configs/1/status", "")
			h.UpdateStatus(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "invalid payload", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"update suggested offer config status error not found", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.AddParam("id", "1")
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "test@dnse.com",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"PATCH", "/api/v1/suggested-offer-configs/1/status", http.UpdateSuggestedOfferConfigStatusRequest{
					Status: "ACTIVE",
				},
			)
			h.UpdateStatus(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "not found resources", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"update suggested offer config status error exist active config", func(t *testing.T) {
			defer truncateData()
			mock.SeedSuggestedOfferConfig(
				t, db, model.SuggestedOfferConfig{
					Name:      "Test 1",
					Value:     decimal.NewFromFloat(0.69),
					ValueType: "LOAN_RATE",
					Status:    "ACTIVE",
					CreatedBy: "admin",
				},
			)
			suggestedOfferConfig := mock.SeedSuggestedOfferConfig(
				t, db, model.SuggestedOfferConfig{
					Name:      "Test 2",
					Value:     decimal.NewFromFloat(0.59),
					ValueType: "INTEREST_RATE",
					Status:    "INACTIVE",
					CreatedBy: "admin",
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.AddParam("id", strconv.FormatInt(suggestedOfferConfig.ID, 10))
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "admin@gmail.com",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"PATCH", "/api/v1/suggested-offer-configs/1/status", http.UpdateSuggestedOfferConfigStatusRequest{
					Status: "ACTIVE",
				},
			)
			h.UpdateStatus(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "exist active suggested offer config", testhelper.GetString(body, "error"))
		},
	)
}
