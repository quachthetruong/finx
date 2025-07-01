package test

import (
	"log/slog"
	"os"
	"testing"

	"financing-offer/pkg/dbtest"
)

func TestMain(m *testing.M) {
	slog.Info("Start integration test..")
	cleanUpPostgres := dbtest.InitPostgres()
	code := m.Run()
	slog.Info("End integration test..")
	cleanUpPostgres()
	os.Exit(code)
}
