package flex_api

import (
	"context"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/config"
)

func TestClient_IsHOActive(t *testing.T) {
	defer gock.Off()
	flexOpenApiConfig := config.FlexOpenApiConfig{
		Url:      "http://flex-api",
		Username: "admin",
		Password: "123456",
	}
	client := NewClient(flexOpenApiConfig)
	t.Run(
		"check when HO is active", func(t *testing.T) {
			gock.New(flexOpenApiConfig.Url).Get("/system/checkHOStatus").
				Reply(200).
				JSON(
					IsHOActiveResponse{
						ErrorCode: "ok",
						HOStatus:  "1",
					},
				)
			isActive, err := client.IsHOActive(context.Background())
			assert.Nil(t, err)
			assert.Equal(t, true, isActive)
		})
	t.Run(
		"check when HO is not active", func(t *testing.T) {
			gock.New(flexOpenApiConfig.Url).Get("/system/checkHOStatus").
				Reply(200).
				JSON(
					IsHOActiveResponse{
						ErrorCode: "ok",
						HOStatus:  "0",
					},
				)
			isActive, err := client.IsHOActive(context.Background())
			assert.Nil(t, err)
			assert.Equal(t, false, isActive)
		})
	t.Run(
		"check when return error", func(t *testing.T) {
			gock.New(flexOpenApiConfig.Url).Get("/system/checkHOStatus").
				Reply(200).
				JSON(
					IsHOActiveResponse{
						ErrorCode: "error",
						HOStatus:  "0",
					},
				)
			isActive, err := client.IsHOActive(context.Background())
			assert.ErrorContains(t, err, "got error code error")
			assert.Equal(t, false, isActive)
		})
}
