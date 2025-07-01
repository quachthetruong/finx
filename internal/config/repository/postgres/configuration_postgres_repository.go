package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"financing-offer/internal/config/repository"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	string_helper "financing-offer/pkg/string-helper"
)

const promotionLoanPackageAttributeName = "promotionLoanPackage"
const loanRateAttributeName = "loanRate"
const marginPoolAttributeName = "marginPool"
const submissionDefaultAttributeName = "submissionDefault"

var _ repository.ConfigurationPersistenceRepository = &ConfigurationPostgresRepository{}

type ConfigurationPostgresRepository struct {
	getDbFunc database.GetDbFunc
}

func NewConfigurationPostgresRepository(getDbFunc database.GetDbFunc) *ConfigurationPostgresRepository {
	return &ConfigurationPostgresRepository{
		getDbFunc: getDbFunc,
	}
}

func (r *ConfigurationPostgresRepository) SetPromotionConfiguration(ctx context.Context, promotionLoanPackage entity.PromotionLoanPackage, updater string) error {
	loanPackages, err := json.Marshal(promotionLoanPackage)
	if err != nil {
		return fmt.Errorf(
			"ConfigurationPostgresRepository SetPromotionConfiguration: %w", err,
		)
	}
	loanPackagesString := string(loanPackages)
	insertModel := model.FinancialConfiguration{
		Attribute:     promotionLoanPackageAttributeName,
		Value:         loanPackagesString,
		LastUpdatedBy: updater,
	}
	if _, err := table.FinancialConfiguration.
		INSERT(table.FinancialConfiguration.MutableColumns).
		MODEL(insertModel).
		ON_CONFLICT(table.FinancialConfiguration.Attribute).
		DO_UPDATE(
			postgres.SET(
				table.FinancialConfiguration.Value.SET(postgres.Json(loanPackagesString)),
				table.FinancialConfiguration.LastUpdatedBy.SET(postgres.String(updater)),
			),
		).ExecContext(
		ctx, r.getDbFunc(ctx),
	); err != nil {
		return fmt.Errorf("ConfigurationPostgresRepository SetPromotionConfiguration: %w", err)
	}
	return nil
}

func (r *ConfigurationPostgresRepository) GetPromotionConfiguration(ctx context.Context) (entity.PromotionLoanPackage, error) {
	dest := model.FinancialConfiguration{}
	if err := table.FinancialConfiguration.
		SELECT(table.FinancialConfiguration.AllColumns).
		WHERE(table.FinancialConfiguration.Attribute.EQ(postgres.String(promotionLoanPackageAttributeName))).
		QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return entity.PromotionLoanPackage{
				LoanProducts: []entity.PromotionLoanProduct{},
			}, nil
		}
		return entity.PromotionLoanPackage{}, fmt.Errorf(
			"ConfigurationPostgresRepository GetPromotionConfiguration: %w", err,
		)
	}
	loanPackages := entity.PromotionLoanPackage{}
	if err := json.Unmarshal(string_helper.StringToBytes(dest.Value), &loanPackages); err != nil {
		return entity.PromotionLoanPackage{}, fmt.Errorf(
			"ConfigurationPostgresRepository GetPromotionConfiguration: %w", err,
		)
	}
	return loanPackages, nil
}

func (r *ConfigurationPostgresRepository) SetLoanRateConfiguration(ctx context.Context, loanRate entity.LoanRateConfiguration, updater string) error {
	errTemplate := "ConfigurationPostgresRepository SetLoanRateConfiguration: %w"
	value, err := json.Marshal(loanRate)
	if err != nil {
		return fmt.Errorf(
			errTemplate, err,
		)
	}
	insertModel := model.FinancialConfiguration{
		Attribute:     loanRateAttributeName,
		Value:         string(value),
		LastUpdatedBy: updater,
	}
	if _, err := table.FinancialConfiguration.
		INSERT(table.FinancialConfiguration.MutableColumns).
		MODEL(insertModel).
		ON_CONFLICT(table.FinancialConfiguration.Attribute).
		DO_UPDATE(
			postgres.SET(
				table.FinancialConfiguration.Value.SET(postgres.Json(string(value))),
				table.FinancialConfiguration.LastUpdatedBy.SET(postgres.String(updater)),
			),
		).ExecContext(
		ctx, r.getDbFunc(ctx),
	); err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	return nil
}

func (r *ConfigurationPostgresRepository) GetLoanRateConfiguration(ctx context.Context) (entity.LoanRateConfiguration, error) {
	errTemplate := "ConfigurationPostgresRepository GetLoanRateConfiguration: %w"
	dest := model.FinancialConfiguration{}
	if err := table.FinancialConfiguration.
		SELECT(table.FinancialConfiguration.AllColumns).
		WHERE(table.FinancialConfiguration.Attribute.EQ(postgres.String(loanRateAttributeName))).
		QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return entity.LoanRateConfiguration{}, nil
		}
		return entity.LoanRateConfiguration{}, fmt.Errorf(
			errTemplate, err,
		)
	}
	loanRate := entity.LoanRateConfiguration{}
	if err := json.Unmarshal(string_helper.StringToBytes(dest.Value), &loanRate); err != nil {
		return entity.LoanRateConfiguration{}, fmt.Errorf(
			errTemplate, err,
		)
	}
	return loanRate, nil
}

func (r *ConfigurationPostgresRepository) SetMarginPoolConfiguration(ctx context.Context, marginPool entity.MarginPoolConfiguration, updater string) error {
	errTemplate := "ConfigurationPostgresRepository SetMarginPoolConfiguration: %w"
	value, err := json.Marshal(marginPool)
	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	insertModel := model.FinancialConfiguration{
		Attribute:     marginPoolAttributeName,
		Value:         string(value),
		LastUpdatedBy: updater,
	}
	if _, err := table.FinancialConfiguration.
		INSERT(table.FinancialConfiguration.MutableColumns).
		MODEL(insertModel).
		ON_CONFLICT(table.FinancialConfiguration.Attribute).
		DO_UPDATE(
			postgres.SET(
				table.FinancialConfiguration.Value.SET(postgres.Json(string(value))),
				table.FinancialConfiguration.LastUpdatedBy.SET(postgres.String(updater)),
			),
		).ExecContext(
		ctx, r.getDbFunc(ctx),
	); err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	return nil
}

func (r *ConfigurationPostgresRepository) GetMarginPoolConfiguration(ctx context.Context) (entity.MarginPoolConfiguration, error) {
	errTemplate := "ConfigurationPostgresRepository GetMarginPoolConfiguration: %w"
	dest := model.FinancialConfiguration{}
	if err := table.FinancialConfiguration.
		SELECT(table.FinancialConfiguration.AllColumns).
		WHERE(table.FinancialConfiguration.Attribute.EQ(postgres.String(marginPoolAttributeName))).
		QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return entity.MarginPoolConfiguration{}, nil
		}
		return entity.MarginPoolConfiguration{}, fmt.Errorf(errTemplate, err)
	}
	marginPool := entity.MarginPoolConfiguration{}
	if err := json.Unmarshal(string_helper.StringToBytes(dest.Value), &marginPool); err != nil {
		return entity.MarginPoolConfiguration{}, fmt.Errorf(errTemplate, err)
	}
	return marginPool, nil
}

func (r *ConfigurationPostgresRepository) SetSubmissionDefault(ctx context.Context, defaultValues entity.SubmissionDefault, updater string) error {
	errTemplate := "ConfigurationPostgresRepository SetSubmissionDefault: %w"
	value, err := json.Marshal(defaultValues)
	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	insertModel := model.FinancialConfiguration{
		Attribute:     submissionDefaultAttributeName,
		Value:         string(value),
		LastUpdatedBy: updater,
	}
	if _, err := table.FinancialConfiguration.
		INSERT(table.FinancialConfiguration.MutableColumns).
		MODEL(insertModel).
		ON_CONFLICT(table.FinancialConfiguration.Attribute).
		DO_UPDATE(
			postgres.SET(
				table.FinancialConfiguration.Value.SET(postgres.Json(string(value))),
				table.FinancialConfiguration.LastUpdatedBy.SET(postgres.String(updater)),
			),
		).ExecContext(
		ctx, r.getDbFunc(ctx),
	); err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	return nil
}

func (r *ConfigurationPostgresRepository) GetSubmissionDefault(ctx context.Context) (entity.SubmissionDefault, error) {
	errTemplate := "ConfigurationPostgresRepository GetSubmissionDefault: %w"
	dest := model.FinancialConfiguration{}
	if err := table.FinancialConfiguration.
		SELECT(table.FinancialConfiguration.AllColumns).
		WHERE(table.FinancialConfiguration.Attribute.EQ(postgres.String(submissionDefaultAttributeName))).
		QueryContext(ctx, r.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return entity.SubmissionDefault{}, nil
		}
		return entity.SubmissionDefault{}, fmt.Errorf(errTemplate, err)
	}
	defaultValue := entity.SubmissionDefault{}
	if err := json.Unmarshal(string_helper.StringToBytes(dest.Value), &defaultValue); err != nil {
		return entity.SubmissionDefault{}, fmt.Errorf(errTemplate, err)
	}
	return defaultValue, nil
}
