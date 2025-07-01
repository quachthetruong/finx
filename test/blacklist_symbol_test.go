package test

import (
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v9"

	"financing-offer/internal/core/blacklistsymbol/transport/http"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestBlacklistSymbolHandler(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	h := do.MustInvoke[*http.BlacklistSymbolHandler](injector)

	t.Run(
		"create blacklist symbol handler", func(t *testing.T) {
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
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/symbols/1/blacklist_symbols", http.BlacklistSymbolRequest{
					AffectedFrom: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					AffectedTo:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
					Status:       entity.BlacklistSymbolStatusActive,
				},
			)
			ginCtx.Params = []gin.Param{
				{
					Key:   "id",
					Value: "1",
				},
			}
			h.Create(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, int64(1), testhelper.GetInt(body, "data", "symbolId"))
		},
	)

	t.Run(
		"create blacklist symbol invalid id", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/symbols/1/blacklist_symbols", http.BlacklistSymbolRequest{},
			)
			ginCtx.Params = []gin.Param{
				{
					Key:   "id",
					Value: "vietnam",
				},
			}
			h.Create(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "id invalid", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"create blacklist symbol invalid payload", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/symbols/1/blacklist_symbols",
				struct{ AffectedTo string }{AffectedTo: "vietnam"},
			)
			ginCtx.Params = []gin.Param{
				{
					Key:   "id",
					Value: "1",
				},
			}

			h.Create(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "invalid payload", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"create blacklist symbol error overlap", func(t *testing.T) {
			defer truncateData()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 50,
					MaxScore: 100,
				},
			)
			sym := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "BID",
					AssetType:       "UNDERLYING",
				},
			)

			mock.SeedBlacklistSymbol(
				t, db, model.BlacklistSymbol{
					SymbolID:     sym.ID,
					AffectedFrom: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					AffectedTo:   null.TimeFrom(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					Status:       model.Blacklistsymbolstatus_Active,
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/symbols/1/blacklist_symbols",
				http.BlacklistSymbolRequest{
					AffectedFrom: time.Date(2021, 12, 1, 2, 1, 3, 4, time.UTC),
					AffectedTo:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					Status:       entity.BlacklistSymbolStatusActive,
				},
			)
			ginCtx.Params = []gin.Param{
				{
					Key:   "id",
					Value: "1",
				},
			}
			ginCtx.Keys = map[string]interface{}{
				"UserInformation": "encapital",
				"requestId":       "1234567890",
			}
			h.Create(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "affected time overlap", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"update blacklist symbol", func(t *testing.T) {
			defer truncateData()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 50,
					MaxScore: 100,
				},
			)
			sym := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "BID",
					AssetType:       "UNDERLYING",
				},
			)

			blSym := mock.SeedBlacklistSymbol(
				t, db, model.BlacklistSymbol{
					SymbolID:     sym.ID,
					AffectedFrom: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					AffectedTo:   null.TimeFrom(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					Status:       model.Blacklistsymbolstatus_Active,
				},
			)

			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"PATCH", "/api/v1/blacklist_symbols/1", http.BlacklistSymbolRequest{
					AffectedFrom: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					AffectedTo:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
					Status:       entity.BlacklistSymbolStatusInactive,
				},
			)
			ginCtx.Params = []gin.Param{
				{
					Key:   "id",
					Value: strconv.FormatInt(blSym.ID, 10),
				},
			}
			h.Update(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, blSym.SymbolID, testhelper.GetInt(body, "data", "symbolId"))
			assert.Equal(t, "INACTIVE", testhelper.GetString(body, "data", "status"))
		},
	)

	t.Run(
		"update blacklist symbol invalid id", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"PATCH", "/api/v1/blacklist_symbols/1", http.BlacklistSymbolRequest{},
			)
			ginCtx.Params = []gin.Param{
				{
					Key:   "id",
					Value: "vietnam",
				},
			}
			h.Update(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "id invalid", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"Test update blacklist symbol invalid payload", func(t *testing.T) {
			defer truncateData()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 50,
					MaxScore: 100,
				},
			)
			sym := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "BID",
					AssetType:       "UNDERLYING",
				},
			)

			mock.SeedBlacklistSymbol(
				t, db, model.BlacklistSymbol{
					SymbolID:     sym.ID,
					AffectedFrom: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					AffectedTo:   null.TimeFrom(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					Status:       model.Blacklistsymbolstatus_Active,
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"PATCH", "/api/v1/blacklist_symbols/1",
				struct{ AffectedTo string }{AffectedTo: "vietnam"},
			)
			ginCtx.Params = []gin.Param{
				{
					Key:   "id",
					Value: "1",
				},
			}
			h.Update(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "invalid payload", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"Test update blacklist symbol not found", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"PATCH", "/api/v1/blacklist_symbols/1", http.BlacklistSymbolRequest{},
			)
			ginCtx.Params = []gin.Param{
				{
					Key:   "id",
					Value: "1",
				},
			}
			h.Update(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "blacklist symbol id: 1 not found", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"Test update blacklist symbol error overlap", func(t *testing.T) {
			defer truncateData()
			se := mock.SeedStockExchange(
				t, db, model.StockExchange{
					Code:     "HOSE",
					MinScore: 50,
					MaxScore: 100,
				},
			)
			sym := mock.SeedSymbol(
				t, db, model.Symbol{
					StockExchangeID: se.ID,
					Symbol:          "BID",
					AssetType:       "UNDERLYING",
				},
			)

			blSym := mock.SeedBlacklistSymbol(
				t, db, model.BlacklistSymbol{
					SymbolID:     sym.ID,
					AffectedFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					AffectedTo:   null.TimeFrom(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
					Status:       model.Blacklistsymbolstatus_Active,
				},
			)
			mock.SeedBlacklistSymbol(
				t, db, model.BlacklistSymbol{
					SymbolID:     sym.ID,
					AffectedFrom: time.Date(2021, 12, 1, 0, 0, 0, 0, time.UTC),
					AffectedTo:   null.TimeFrom(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					Status:       model.Blacklistsymbolstatus_Active,
				},
			)

			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"PATCH", "/api/v1/blacklist_symbols/1",
				http.BlacklistSymbolRequest{
					AffectedFrom: time.Date(2021, 11, 1, 2, 1, 3, 4, time.UTC),
					AffectedTo:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					Status:       entity.BlacklistSymbolStatusActive,
				},
			)
			ginCtx.Params = []gin.Param{
				{
					Key:   "id",
					Value: strconv.FormatInt(blSym.ID, 10),
				},
			}
			ginCtx.Keys = map[string]interface{}{
				"UserInformation": "encapital",
				"requestId":       "1234567890",
			}
			h.Update(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "affected time overlap", testhelper.GetString(body, "error"))
		},
	)
}
