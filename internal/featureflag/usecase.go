package featureflag

import (
	"slices"

	"financing-offer/internal/config"
)

type UseCase interface {
	IsFeatureEnable(featureName string, investorId string) (bool, error)
}

type useCase struct {
	features map[string]config.FeatureConfig
}

func (u *useCase) IsFeatureEnable(featureName string, investorId string) (bool, error) {
	feature, ok := u.features[featureName]
	if !ok {
		return false, nil
	}
	if feature.Enable {
		return true, nil
	}
	return slices.Contains(feature.InvestorIds, investorId), nil
}

func NewUseCase(features map[string]config.FeatureConfig) UseCase {
	return &useCase{
		features: features,
	}
}
