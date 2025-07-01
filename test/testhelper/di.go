package testhelper

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"sync"

	"github.com/samber/do"
	"github.com/shopspring/decimal"

	"financing-offer/internal/apperrors/repository"
	"financing-offer/internal/atomicity"
	"financing-offer/internal/config"
	"financing-offer/internal/database"
	"financing-offer/internal/di"
	"financing-offer/pkg/environment"
	"financing-offer/pkg/shutdown"
	testConfig "financing-offer/test/config"
	"financing-offer/test/mock"
)

var configOnce sync.Once

type TestInjectorOpt func(i *do.Injector)

func WithDb(db *sql.DB) TestInjectorOpt {
	return func(i *do.Injector) {
		do.ProvideValue[*sql.DB](i, db)
	}
}

func NewInjector(opts ...TestInjectorOpt) *do.Injector {
	configOnce.Do(
		func() {
			decimal.MarshalJSONWithoutQuotes = true
		},
	)
	cfg, _ := config.InitConfig[config.AppConfig](testConfig.EmbeddedFiles)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	i := di.NewInjector(logger)
	do.ProvideValue[config.AppConfig](i, cfg)
	do.ProvideValue(i, environment.Development)
	do.ProvideValue[*shutdown.Tasks](i, &shutdown.Tasks{})
	do.Provide[*atomicity.DbAtomicExecutor](
		i, func(injector *do.Injector) (*atomicity.DbAtomicExecutor, error) {
			db := do.MustInvoke[*sql.DB](i)
			return &atomicity.DbAtomicExecutor{DB: db}, nil
		},
	)
	do.Provide[database.GetDbFunc](
		i, func(injector *do.Injector) (database.GetDbFunc, error) {
			db := do.MustInvoke[*sql.DB](i)
			return func(ctx context.Context) database.DB {
				if tx := atomicity.ContextGetTx(ctx); tx != nil {
					return tx
				}
				return db
			}, nil
		},
	)
	do.OverrideValue[repository.NotifyWebhookRepository](i, &mock.NotifyWebhookRepository{})
	for _, opt := range opts {
		opt(i)
	}
	return i
}
