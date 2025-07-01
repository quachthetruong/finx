package config

import (
	"financing-offer/internal/config/repository"
)

type UseCase interface {
	GetConfigurations() (map[string]any, error)
}

type useCase struct {
	cfg                          AppConfig
	configurationPersistenceRepo repository.ConfigurationPersistenceRepository
}

func (u *useCase) GetConfigurations() (map[string]any, error) {
	res := map[string]any{
		"guaranteeFeeRate":            u.cfg.LoanRequest.GuaranteeFeeRate,
		"maxGuaranteedDuration":       u.cfg.LoanRequest.MaxGuaranteedDuration,
		"minimumAppVersion":           u.cfg.LoanRequest.MinimumAppVersion,
		"minimumAppVersionDerivative": u.cfg.LoanRequest.MinimumAppVersionDerivative,
	}
	return res, nil
}

func NewUseCase(
	cfg AppConfig,
	configurationPersistenceRepo repository.ConfigurationPersistenceRepository,
) UseCase {
	return &useCase{
		cfg:                          cfg,
		configurationPersistenceRepo: configurationPersistenceRepo,
	}
}
