package test

import (
	"database/sql"
	"encoding/json"
	"financing-offer/internal/core/entity"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/core/promotion_campaign/transport/http"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/jwttoken"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/mock"
	"financing-offer/test/testhelper"
)

func TestPromotionCampaignHandler(t *testing.T) {
	t.Parallel()

	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	h := do.MustInvoke[*http.PromotionCampaignHandler](injector)
	t.Run(
		"create promotion campaign success", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "test@dnse.com",
				},
			)
			ginCtx.Request = gintest.MustMakeRequest(
				"POST", "/api/v1/promotion-campaigns", http.CreatePromotionCampaignRequest{
					Name:        "name",
					Tag:         "tag",
					Description: "description",
				},
			)
			h.Create(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())

			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "name", testhelper.GetString(body, "data", "name"))
			assert.Equal(t, "tag", testhelper.GetString(body, "data", "tag"))
			assert.Equal(t, "description", testhelper.GetString(body, "data", "description"))
			assert.Equal(t, "ACTIVE", testhelper.GetString(body, "data", "status"))
		},
	)

	t.Run(
		"create promotion campaign invalid payload", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("POST", "/api/v1/promotion-campaigns", "")
			h.Create(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "invalid payload", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"create promotion campaign error", func(t *testing.T) {
			defer func() {
				db, tearDownDb, truncateData = dbtest.NewDb(t)
				do.OverrideValue[*sql.DB](injector, db)
				truncateData()
				injector := testhelper.NewInjector(testhelper.WithDb(db))
				h = do.MustInvoke[*http.PromotionCampaignHandler](injector)
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
				"POST", "/api/v1/promotion-campaigns", http.CreatePromotionCampaignRequest{
					Name:        "name",
					Tag:         "tag",
					Description: "description",
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
		"update promotion campaign success", func(t *testing.T) {
			defer truncateData()
			metadata := entity.PromotionCampaignMetadata{
				Products: []entity.PromotionCampaignProduct{
					{
						Symbols:       []string{"HPG", "DGW"},
						LoanPackageId: 1,
						RetailSymbols: []string{"HPG", "DGW"},
					},
				},
			}
			campaign := mock.SeedPromotionCampaign(
				t, db, model.PromotionCampaign{
					Name:        "name 1",
					Status:      "ACTIVE",
					UpdatedBy:   "admin",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Tag:         "tag",
					Description: "description",
					Metadata:    "{}",
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Set(
				appcontext.UserInformation, &jwttoken.AdminClaims{
					Sub: "test@dnse.com",
				},
			)
			ginCtx.AddParam("id", strconv.FormatInt(campaign.ID, 10))
			ginCtx.Request = gintest.MustMakeRequest(
				"PATCH", fmt.Sprintf("/api/v1/promotion-campaigns/%v", campaign.ID), entity.PromotionCampaign{
					Name:        "name 2",
					Description: "description 2",
					Metadata:    metadata,
				},
			)
			h.Update(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "name 2", testhelper.GetString(body, "data", "name"))
			assert.Equal(t, "description 2", testhelper.GetString(body, "data", "description"))
		},
	)

	t.Run(
		"update promotion campaign id invalid", func(t *testing.T) {
			defer truncateData()
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest("PATCH", "/api/v1/promotion-campaigns/1", "")
			h.Update(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "id invalid", testhelper.GetString(body, "error"))
		},
	)

	t.Run(
		"get promotion campaigns success", func(t *testing.T) {
			defer truncateData()
			metadata := entity.PromotionCampaignMetadata{
				Products: []entity.PromotionCampaignProduct{
					{
						Symbols:       []string{"HPG", "DGW"},
						LoanPackageId: 1,
						RetailSymbols: []string{"HPG", "DGW"},
					},
				},
			}
			metadataJson, _ := json.Marshal(metadata)
			campaign := mock.SeedPromotionCampaign(
				t, db, model.PromotionCampaign{
					ID:          1,
					Name:        "name 1",
					Status:      "ACTIVE",
					UpdatedBy:   "admin",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Tag:         "tag",
					Description: "description",
					Metadata:    string(metadataJson),
				},
			)
			campaign2 := mock.SeedPromotionCampaign(
				t, db, model.PromotionCampaign{
					ID:          2,
					Name:        "name 2",
					Status:      "ACTIVE",
					UpdatedBy:   "admin",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Tag:         "tag",
					Description: "description",
					Metadata:    string(metadataJson),
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/promotion-campaigns", nil,
			)
			h.GetAll(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, campaign.Name, testhelper.GetString(body, "data", "[0]", "name"))
			assert.Equal(t, campaign2.Name, testhelper.GetString(body, "data", "[1]", "name"))
		},
	)

	t.Run(
		"get promotion campaigns success with status = ACTIVE", func(t *testing.T) {
			defer truncateData()
			metadata := entity.PromotionCampaignMetadata{
				Products: []entity.PromotionCampaignProduct{
					{
						Symbols:       []string{"HPG", "DGW"},
						LoanPackageId: 1,
						RetailSymbols: []string{"HPG", "DGW"},
					},
				},
			}
			metadataJson, _ := json.Marshal(metadata)
			campaign := mock.SeedPromotionCampaign(
				t, db, model.PromotionCampaign{
					ID:          1,
					Name:        "name 1",
					Status:      "ACTIVE",
					UpdatedBy:   "admin",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Tag:         "tag",
					Description: "description",
					Metadata:    string(metadataJson),
				},
			)
			_ = mock.SeedPromotionCampaign(
				t, db, model.PromotionCampaign{
					ID:          2,
					Name:        "name 2",
					Status:      "INACTIVE",
					UpdatedBy:   "admin",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Tag:         "tag",
					Description: "description",
					Metadata:    string(metadataJson),
				},
			)
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/promotion-campaigns?status=ACTIVE", nil,
			)
			h.GetAll(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, campaign.Name, testhelper.GetString(body, "data", "[0]", "name"))
			assert.Equal(t, 1, testhelper.GetArrayLength(body, "data"))
		},
	)

	t.Run(
		"get promotion campaigns error", func(t *testing.T) {
			defer func() {
				db, tearDownDb, truncateData = dbtest.NewDb(t)
				do.OverrideValue[*sql.DB](injector, db)
				truncateData()
				injector := testhelper.NewInjector(testhelper.WithDb(db))
				h = do.MustInvoke[*http.PromotionCampaignHandler](injector)
			}()

			if err := db.Close(); err != nil {
				t.Fatal(err)
			}
			ginCtx, _, recorder := gintest.GetTestContext()
			ginCtx.Request = gintest.MustMakeRequest(
				"GET", "/api/v1/promotion-campaigns", nil,
			)
			h.GetAll(ginCtx)
			result := recorder.Result()
			defer assert.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, "an error happened, please try again later", testhelper.GetString(body, "error"))
		},
	)

}
