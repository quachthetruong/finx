package test

import (
	"encoding/json"
	"fmt"
	"log/slog"
	http2 "net/http"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/core/entity"
	orderServiceRepo2 "financing-offer/internal/core/orderservice/repository"
	"financing-offer/internal/core/suggested_offer/transport/http"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/event"
	"financing-offer/internal/jwttoken"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/pkg/shutdown"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestSuggestedOfferHandler_CreateOffer(t *testing.T) {
	t.Parallel()
	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	tasks, _ := shutdown.NewShutdownTasks(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	kafkaPublisher := mock.NewMockPublisher(t)
	kafkaPublisher.EXPECT().Publish(testifyMock.Anything).Return(nil)
	orderServiceRepo := mock.NewMockOrderServiceRepository(t)
	orderServiceRepo.EXPECT().GetAccountByAccountNoAndCustodyCode(testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).Return(entity.OrderServiceAccount{}, nil)
	do.OverrideValue[event.Publisher](injector, kafkaPublisher)
	do.OverrideValue[*shutdown.Tasks](injector, tasks)
	do.OverrideValue[orderServiceRepo2.OrderServiceRepository](injector, orderServiceRepo)
	handler := do.MustInvoke[*http.SuggestedOfferHandler](injector)

	t.Run("create offer with invalid payload", func(t *testing.T) {
		defer truncateData()
		ginCtx, _, recorder := gintest.GetTestContext()
		ginCtx.Request = gintest.MustMakeRequest(
			"POST",
			"/api/v1/suggested-offers",
			json.RawMessage(`{"configId": 1, "symbols":"bar", "accountNo": "1234"}`),
		)
		ginCtx.Params = gin.Params{{
			Key:   "id",
			Value: "1",
		}}
		handler.CreateOffer(ginCtx)
		result := recorder.Result()
		gintest.ExtractBody(result.Body)
		assert.Nil(t, result.Body.Close())
		assert.Equal(t, http2.StatusBadRequest, result.StatusCode)
	})

	t.Run("create offer with invalid user", func(t *testing.T) {
		defer truncateData()
		ginCtx, _, recorder := gintest.GetTestContext()
		ginCtx.Request = gintest.MustMakeRequest(
			"POST",
			"/api/v1/suggested-offers",
			json.RawMessage(`{"configId": 1, "symbols":["ACB"], "accountNo": "1234"}`),
		)
		ginCtx.Params = gin.Params{{
			Key:   "id",
			Value: "1",
		}}
		handler.CreateOffer(ginCtx)
		result := recorder.Result()
		gintest.ExtractBody(result.Body)
		assert.Nil(t, result.Body.Close())
		assert.Equal(t, http2.StatusUnauthorized, result.StatusCode)
	})

	t.Run("create offer with config not found", func(t *testing.T) {
		defer truncateData()
		mock.SeedSuggestedOfferConfig(t, db, model.SuggestedOfferConfig{
			Name:      "config",
			Value:     decimal.NewFromFloat(0.2),
			ValueType: "INTEREST_RATE",
			Status:    "ACTIVE",
		})
		ginCtx, _, recorder := gintest.GetTestContext()
		ginCtx.Request = gintest.MustMakeRequest(
			"POST",
			"/api/v1/suggested-offers",
			json.RawMessage(`{"configId": 10, "symbols":["ACB"], "accountNo": "1234"}`),
		)
		ginCtx.Set(
			appcontext.UserInformation, &jwttoken.AdminClaims{
				InvestorId: "0001000115",
				Sub:        "0001000115",
			},
		)
		handler.CreateOffer(ginCtx)
		result := recorder.Result()
		gintest.ExtractBody(result.Body)
		assert.Nil(t, result.Body.Close())
		assert.Equal(t, http2.StatusNotFound, result.StatusCode)
	})

	t.Run("create offer with config success", func(t *testing.T) {
		defer truncateData()
		config := mock.SeedSuggestedOfferConfig(t, db, model.SuggestedOfferConfig{
			Name:      "config",
			Value:     decimal.NewFromFloat(0.2),
			ValueType: "INTEREST_RATE",
			Status:    "ACTIVE",
		})
		ginCtx, _, recorder := gintest.GetTestContext()
		ginCtx.Request = gintest.MustMakeRequest(
			"POST",
			"/api/v1/suggested-offers",
			json.RawMessage(fmt.Sprintf("{\"configId\": %s, \"symbols\":[\"ACB\"], \"accountNo\": \"1234\"}", strconv.FormatInt(config.ID, 10))),
		)
		ginCtx.Set(
			appcontext.UserInformation, &jwttoken.AdminClaims{
				InvestorId: "0001000115",
				Sub:        "0001000115",
			},
		)
		handler.CreateOffer(ginCtx)
		result := recorder.Result()
		gintest.ExtractBody(result.Body)
		assert.Nil(t, result.Body.Close())
		assert.Equal(t, http2.StatusOK, result.StatusCode)
	})
}
