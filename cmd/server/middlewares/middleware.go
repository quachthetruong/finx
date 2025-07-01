package middlewares

import (
	"log/slog"

	"financing-offer/internal/config"
	flexOpenApiRepo "financing-offer/internal/core/flex/repository"
	"financing-offer/internal/featureflag"
)

type Middleware struct {
	Logger             *slog.Logger
	Config             config.AppConfig
	FeatureFlagUseCase featureflag.UseCase
	FlexRepo           flexOpenApiRepo.FlexOpenApiRepository
}
