package mo_service

import (
	"context"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/config"
)

func TestClient_GetMarginPoolsByIds(t *testing.T) {
	defer gock.Off()
	moServiceConfig := config.MoServiceConfig{
		Url:   "http://dnse-mo-service",
		Token: "financing-product-token",
	}
	client := NewClient(moServiceConfig)
	t.Run(
		"get margin pools by ids success", func(t *testing.T) {
			gock.New(moServiceConfig.Url).Get("/pools").HeaderPresent("Authorization").MatchParam(
				"ids", "887,889",
			).Reply(200).BodyString(
				`
						{
							"data": [
								{
									"id": 887,
									"poolGroupId": 1000,
									"poolLimit": 30000000000,
									"name": "Pool high margin",
									"type": "RESERVE",
									"createdDate": "2022-12-06T07:03:59.894121Z",
									"modifiedDate": "2022-12-06T07:03:59.894152Z",
									"status": "ACTIVE"
								},
								{
									"id": 889,
									"poolGroupId": 1001,
									"poolLimit": 2000500000,
									"name": "Pool",
									"type": "RESERVE",
									"createdDate": "2022-12-06T07:03:59.894121Z",
									"modifiedDate": "2022-12-06T07:03:59.894152Z",
									"status": "ACTIVE"
								}		
							],
							"total": 2,
							"start": 0,
							"end": 2
						}`,
			)
			pools, err := client.GetMarginPoolsByIds(context.Background(), []int64{887, 889})
			assert.Nil(t, err)
			assert.Equal(t, 2, len(pools))
			assert.Equal(t, int64(889), pools[1].Id)
		},
	)

	t.Run(
		"get margin pools by ids failed", func(t *testing.T) {
			gock.New(moServiceConfig.Url).Get("/pools").HeaderPresent("Authorization").MatchParam(
				"ids", "1",
			).Reply(400).BodyString("{}")
			_, err := client.GetMarginPoolsByIds(context.Background(), []int64{1})
			assert.Error(t, err)
			dest := apperrors.AppError{}
			assert.ErrorAs(t, err, &dest)
			assert.Equal(t, 400, dest.Code)
		},
	)
}

func TestClient_GetMarginPoolsById(t *testing.T) {
	defer gock.Off()
	moServiceConfig := config.MoServiceConfig{
		Url:   "http://dnse-mo-service",
		Token: "financing-product-token",
	}
	client := NewClient(moServiceConfig)
	t.Run(
		"get margin pool by id success", func(t *testing.T) {
			gock.New(moServiceConfig.Url).Get("/pools/887").HeaderPresent("Authorization").Reply(200).BodyString(
				`
					{
						"id": 887,
						"poolGroupId": 1000,
						"poolLimit": 30000000000,
						"name": "Pool high margin",
						"type": "RESERVE",
						"createdDate": "2022-12-06T07:03:59.894121Z",
						"modifiedDate": "2022-12-06T07:03:59.894152Z",
						"status": "ACTIVE"
					}`,
			)
			pool, err := client.GetMarginPoolById(context.Background(), int64(887))
			assert.Nil(t, err)
			assert.Equal(t, int64(887), pool.Id)
		},
	)

	t.Run(
		"get margin pools by id failed", func(t *testing.T) {
			gock.New(moServiceConfig.Url).Get("/pools").HeaderPresent("Authorization").MatchParam(
				"ids", "1",
			).Reply(400).BodyString("{}")
			_, err := client.GetMarginPoolsByIds(context.Background(), []int64{1})
			assert.Error(t, err)
			dest := apperrors.AppError{}
			assert.ErrorAs(t, err, &dest)
			assert.Equal(t, 400, dest.Code)
		},
	)
}

func TestClient_GetMarginPoolGroupsByIds(t *testing.T) {
	defer gock.Off()
	moServiceConfig := config.MoServiceConfig{
		Url:   "http://dnse-mo-service",
		Token: "financing-product-token",
	}
	client := NewClient(moServiceConfig)
	t.Run(
		"get margin pool groups by ids success", func(t *testing.T) {
			gock.New(moServiceConfig.Url).Get("/pool-groups").HeaderPresent("Authorization").MatchParam(
				"ids", "1000,1001",
			).Reply(200).BodyString(
				`
						{
							"data": [
								{
									"id": 1000,
									"type": "MARGIN",
									"name": "POOL TỔNG",
									"source": "DNSE",
									"totalCash": 3000000000000,
									"allocatedCash": 34000000001,
									"createdDate": "2022-10-31T11:45:04.105368Z",
									"modifiedDate": "2022-10-31T11:45:04.105387Z",
									"status": "ACTIVE"
								},
								{
									"id": 1001,
									"type": "MARGIN",
									"name": "POOL BT3",
									"source": "C",
									"totalCash": 1000000000000,
									"allocatedCash": 400000001,
									"createdDate": "2022-12-06T07:04:18.224204Z",
									"modifiedDate": "2022-12-06T07:04:18.224217Z",
									"status": "ACTIVE"
								}
							],
							"total": 2,
							"start": 0,
							"end": 2
						}`,
			)
			poolGroups, err := client.GetMarginPoolGroupsByIds(context.Background(), []int64{1000, 1001})
			assert.Nil(t, err)
			assert.Equal(t, 2, len(poolGroups))
			assert.Equal(t, int64(1001), poolGroups[1].Id)
		},
	)

	t.Run(
		"get margin pool groups by ids failed", func(t *testing.T) {
			gock.New(moServiceConfig.Url).Get("/pool-groups").HeaderPresent("Authorization").MatchParam(
				"ids", "1",
			).Reply(400).BodyString("{}")
			_, err := client.GetMarginPoolGroupsByIds(context.Background(), []int64{1})
			assert.Error(t, err)
			dest := apperrors.AppError{}
			assert.ErrorAs(t, err, &dest)
			assert.Equal(t, 400, dest.Code)
		},
	)
}
