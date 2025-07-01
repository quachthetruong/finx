package test

import (
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core/scoregroup/transport/http"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestScoreGroup(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))

	scoreGroupHandler := do.MustInvoke[*http.ScoreGroupHandler](injector)

	t.Run(
		"get available package of score group success with data", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			group := mock.SeedScoreGroup(
				t, db, model.ScoreGroup{
					Code:     "A",
					MinScore: 10,
					MaxScore: 20,
				},
			)
			mock.SeedScoreGroupInterest(
				t, db, model.ScoreGroupInterest{
					LimitAmount:  decimal.NewFromFloat(10000.0),
					LoanRate:     decimal.NewFromFloat(3.0),
					InterestRate: decimal.NewFromFloat(4.0),
					ScoreGroupID: group.ID,
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "", nil,
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(group.ID, 10),
			}}
			scoreGroupHandler.GetAvailablePackages(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 1, testhelper.GetArrayLength(body, "data"))
			assert.Equal(t, int64(1), testhelper.GetInt(body, "data", "[0]", "id"))
		},
	)

	t.Run(
		"get available package of score group success with empty data", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			group := mock.SeedScoreGroup(
				t, db, model.ScoreGroup{
					Code:     "A",
					MinScore: 10,
					MaxScore: 20,
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "", nil,
			)
			ginCtx.Params = []gin.Param{{
				Key:   "id",
				Value: strconv.FormatInt(group.ID, 10),
			}}
			scoreGroupHandler.GetAvailablePackages(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, 0, testhelper.GetArrayLength(body, "data"))
		},
	)
}
