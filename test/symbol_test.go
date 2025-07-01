package test

import (
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/core/entity"
	loanPackageRequestRepo "financing-offer/internal/core/loanpackagerequest/repository"
	"financing-offer/internal/core/symbol/transport/http"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/jwttoken"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestSymbolHandler(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	do.OverrideValue[loanPackageRequestRepo.LoanPackageRequestEventRepository](
		injector, &mock.LoanPackageRequestEventRepository{},
	)
	h := do.MustInvoke[*http.SymbolHandler](injector)
	t.Run(
		"create symbol success", func(t *testing.T) {
			defer truncateData()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 50,
					MaxScore: 100,
				},
			)

			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/symbols", http.SymbolRequest{
					StockExchangeId: se.ID,
					Symbol:          "BID",
					AssetType:       "UNDERLYING",
				},
			)
			h.Create(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "BID", testhelper.GetString(body, "data", "symbol"))
		},
	)

	t.Run(
		"get all symbols success", func(t *testing.T) {
			defer truncateData()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 50,
					MaxScore: 100,
				},
			)
			mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "BID",
					AssetType:       "UNDERLYING",
				},
			)
			mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "FPT",
					AssetType:       "UNDERLYING",
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/symbols", nil,
			)
			h.GetAll(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "BID", testhelper.GetString(body, "data", "[0]", "symbol"))
			assert.Equal(t, "FPT", testhelper.GetString(body, "data", "[1]", "symbol"))
		},
	)
	t.Run(
		"get by symbol success", func(t *testing.T) {
			defer truncateData()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 50,
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
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("GET", "", nil)
			ginCtx.AddParam("symbol", "BID")
			h.GetBySymbol(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, symbol.Symbol, testhelper.GetString(body, "data", "symbol"))
		},
	)

	t.Run(
		"get by symbol and not found", func(t *testing.T) {
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("GET", "", nil)
			ginCtx.AddParam("symbol", "BID")
			h.GetBySymbol(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			assert.Equal(t, 404, result.StatusCode)
		},
	)

	t.Run(
		"get by id success", func(t *testing.T) {
			defer truncateData()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 50,
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
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("GET", "", nil)
			ginCtx.AddParam("id", strconv.Itoa(int(symbol.ID)))
			h.GetById(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, symbol.Symbol, testhelper.GetString(body, "data", "symbol"))
		},
	)

	t.Run(
		"update activate symbol derivative success", func(t *testing.T) {
			defer truncateData()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "DERIVATIVE",
					MinScore: 50,
					MaxScore: 100,
				},
			)
			symbolEntity := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "VN30F1Q",
					Status:          "INACTIVE",
					AssetType:       "DERIVATIVE",
				},
			)
			mock.SeedLoanPackageRequest(
				t, db, model.LoanPackageRequest{
					SymbolID:    symbolEntity.ID,
					LimitAmount: decimal.NewFromInt(1000000),
					LoanRate:    decimal.NewFromFloat(3.0),
					Status:      entity.LoanPackageRequestStatusPending.String(),
					AssetType:   "DERIVATIVE",
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "test@dnse.com",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"PATCH", "", http.UpdateSymbolStatusRequest{
					Status: "ACTIVE",
				},
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(symbolEntity.ID, 10),
			}}
			h.Update(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "VN30F1Q", testhelper.GetString(body, "data", "symbol"))
			assert.Equal(t, "ACTIVE", testhelper.GetString(body, "data", "status"))
		},
	)

	t.Run("update deactivate symbol derivative", func(t *testing.T) {
		defer truncateData()
		se := mock.SeedStockExchange(
			t, db, model.StockExchange{
				Code:     "DERIVATIVE",
				MinScore: 50,
				MaxScore: 100,
			},
		)
		symbolEntity := mock.SeedSymbol(
			t, db, model.Symbol{
				StockExchangeID: se.ID,
				Symbol:          "VN30F1Q",
				Status:          "ACTIVE",
				AssetType:       "DERIVATIVE",
			},
		)
		ginCtx, _, recorder := gintest.GetTestContext()
		ginCtx.Set(
			appcontext.UserInformation, &jwttoken.AdminClaims{
				Sub: "test@dnse.com",
			},
		)
		ginCtx.Request = gintest.MustMakeRequest(
			"PATCH", "", http.UpdateSymbolStatusRequest{
				Status: "INACTIVE",
			},
		)
		ginCtx.Params = []gin.Param{{
			Key:   "id",
			Value: strconv.FormatInt(symbolEntity.ID, 10),
		}}
		h.Update(ginCtx)
		result := recorder.Result()
		defer assert.Nil(t, result.Body.Close())

		body := gintest.ExtractBody(result.Body)
		assert.Equal(t, "VN30F1Q", testhelper.GetString(body, "data", "symbol"))
		assert.Equal(t, entity.SymbolStatusInactive.String(), testhelper.GetString(body, "data", "status"))
	})
	t.Run("update symbol derivative with invalid id", func(t *testing.T) {
		defer truncateData()
		ginCtx, _, recorder := gintest.GetTestContext()
		ginCtx.Request = gintest.MustMakeRequest(
			"PATCH", "", http.UpdateSymbolStatusRequest{
				Status: "ACTIVE",
			},
		)
		ginCtx.Params = []gin.Param{{
			Key:   "id",
			Value: "invalid",
		}}
		h.Update(ginCtx)
		result := recorder.Result()
		defer assert.Nil(t, result.Body.Close())
		body := gintest.ExtractBody(result.Body)
		assert.Equal(t, "id invalid", testhelper.GetString(body, "error"))
	})
	t.Run("update symbol derivative with invalid status", func(t *testing.T) {
		defer truncateData()
		se := mock.SeedStockExchange(
			t, db, model.StockExchange{
				Code:     "DERIVATIVE",
				MinScore: 50,
				MaxScore: 100,
			},
		)
		symbolEntity := mock.SeedSymbol(
			t, db, model.Symbol{
				StockExchangeID: se.ID,
				Symbol:          "VN30F1Q",
				Status:          "ACTIVE",
				AssetType:       "DERIVATIVE",
			},
		)
		ginCtx, _, recorder := gintest.GetTestContext()
		ginCtx.Request = gintest.MustMakeRequest(
			"PATCH", "", http.UpdateSymbolStatusRequest{
				Status: "INVALID",
			},
		)
		ginCtx.Params = []gin.Param{{
			Key:   "id",
			Value: strconv.FormatInt(symbolEntity.ID, 10),
		}}
		h.Update(ginCtx)
		result := recorder.Result()
		defer assert.Nil(t, result.Body.Close())
		body := gintest.ExtractBody(result.Body)
		assert.Equal(t, "invalid payload", testhelper.GetString(body, "error"))
	})
	t.Run("update symbol derivative not found symbol", func(t *testing.T) {
		defer truncateData()
		ginCtx, _, recorder := gintest.GetTestContext()
		ginCtx.Request = gintest.MustMakeRequest(
			"PATCH", "", http.UpdateSymbolStatusRequest{
				Status: "ACTIVE",
			},
		)
		ginCtx.Params = []gin.Param{{
			Key:   "id",
			Value: "1",
		}}
		h.Update(ginCtx)
		result := recorder.Result()
		defer assert.Nil(t, result.Body.Close())
		body := gintest.ExtractBody(result.Body)
		assert.Equal(t, "not found resources", testhelper.GetString(body, "error"))
	})
}
