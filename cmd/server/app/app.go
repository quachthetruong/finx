package app

import (
	"financing-offer/assets"
	"log/slog"

	"github.com/samber/do"

	"financing-offer/cmd/server/middlewares"
	"financing-offer/internal/config"
	flexOpenApiRepo "financing-offer/internal/core/flex/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/di"
	"financing-offer/internal/featureflag"
	"financing-offer/pkg/environment"
	"financing-offer/pkg/shutdown"
)

type Application struct {
	Config     config.AppConfig
	Logger     *slog.Logger
	Injector   *do.Injector
	Tasks      *shutdown.Tasks
	Middleware middlewares.Middleware
}

func Run(logger *slog.Logger, tasks *shutdown.Tasks) error {
	cfg, err := config.InitConfig[config.AppConfig](assets.EmbeddedFiles)
	if err != nil {
		return err
	}

	getDbFunc, atomicExecutor, err := database.New(cfg.Db, tasks)
	if err != nil {
		return err
	}
	env := environment.Development
	if cfg.Env == config.EnvProduction {
		env = environment.Production
	}

	injector := di.NewInjector(logger)

	do.ProvideValue(injector, env)
	do.ProvideValue(injector, getDbFunc)
	do.ProvideValue(injector, cfg)
	do.ProvideValue(injector, atomicExecutor)
	do.ProvideValue(injector, tasks)

	application := &Application{
		Config:   cfg,
		Logger:   logger,
		Injector: injector,
		Tasks:    tasks,
		Middleware: middlewares.Middleware{
			Logger:             logger,
			Config:             cfg,
			FeatureFlagUseCase: do.MustInvoke[featureflag.UseCase](injector),
			FlexRepo:           do.MustInvoke[flexOpenApiRepo.FlexOpenApiRepository](injector),
		},
	}
	if err := application.StartScheduler(); err != nil {
		return err
	}

	return application.ServeHTTP()
}
