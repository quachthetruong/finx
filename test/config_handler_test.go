package test

import (
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"financing-offer/internal/config"
	"financing-offer/internal/config/transport/http"
	"financing-offer/pkg/dbtest"
	"financing-offer/pkg/gintest"
	"financing-offer/test/testhelper"
)

func TestConfigHandler(t *testing.T) {
	t.Parallel()
	db, tearDownDb, truncateData := dbtest.NewDb(t)
	defer tearDownDb()
	injector := testhelper.NewInjector(testhelper.WithDb(db))
	configurationHandler := do.MustInvoke[*http.ConfigHandler](injector)
	cfg := do.MustInvoke[config.AppConfig](injector)
	t.Run(
		"test configuration handler", func(t *testing.T) {
			defer truncateData()
			ctx, _, recorder := gintest.GetTestContext()
			configurationHandler.GetConfiguration(ctx)
			result := recorder.Result()
			defer require.Nil(t, result.Body.Close())
			body := gintest.ExtractBody(result.Body)
			assert.Equal(t, cfg.LoanRequest.GuaranteeFeeRate, testhelper.GetFloat(body, "data", "guaranteeFeeRate"))
			assert.Equal(
				t, int64(cfg.LoanRequest.MaxGuaranteedDuration),
				testhelper.GetInt(body, "data", "maxGuaranteedDuration"),
			)
		},
	)
}
