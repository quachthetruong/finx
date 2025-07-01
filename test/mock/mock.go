package mock

import (
	"context"
	"testing"
	"time"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

var _ apperrors.Service = (*ErrReporter)(nil)

type ErrReporter struct{}

func (e ErrReporter) Go(ctx context.Context, f func() error) {
	_ = f()
}

func (e ErrReporter) NotifyError(_ context.Context, _ error) error {
	return nil
}

type LoanPackageRequestEventRepository struct{}

func (l *LoanPackageRequestEventRepository) NotifyRequestDeclined(ctx context.Context, data entity.LoanPackageRequestDeclinedNotify) error {
	return nil
}

func (l *LoanPackageRequestEventRepository) NotifyDerivativeRequestDeclined(ctx context.Context, data entity.LoanPackageDerivativeRequestDeclinedNotify) error {
	return nil
}

func (l *LoanPackageRequestEventRepository) NotifyOfflineConfirmation(ctx context.Context, data entity.RequestOfflineConfirmation) error {
	return nil
}

func (l *LoanPackageRequestEventRepository) NotifyDerivativeOfflineConfirmation(ctx context.Context, data entity.DerivativeRequestOfflineConfirmation) error {
	return nil
}

func (l *LoanPackageRequestEventRepository) NotifyOnlineConfirmation(ctx context.Context, data entity.RequestOnlineConfirmationNotify) error {
	return nil
}

func (l *LoanPackageRequestEventRepository) NotifyRequestConfirmed(ctx context.Context, notifyData entity.LoanPackageRequestConfirmedNotify) error {
	return nil
}

type FinancingApiMock struct{}

func (f FinancingApiMock) GetDateAfter(date time.Time, workingDays int) (time.Time, error) {
	return date.AddDate(0, 0, workingDays), nil
}

func SeedStockExchange(t *testing.T, db database.DB, se model.StockExchange) model.StockExchange {
	res := model.StockExchange{}
	if err := table.StockExchange.
		INSERT(table.StockExchange.MutableColumns).
		RETURNING(table.StockExchange.AllColumns).
		MODEL(se).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedSymbol(t *testing.T, db database.DB, symbol model.Symbol) model.Symbol {
	res := model.Symbol{}
	if symbol.AssetType == "" {
		symbol.AssetType = model.AssetType_Underlying
	}
	if symbol.Status == "" {
		symbol.Status = "ACTIVE"
	}
	if err := table.Symbol.
		INSERT(table.Symbol.MutableColumns).
		RETURNING(table.Symbol.AllColumns).
		MODEL(symbol).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedBlacklistSymbol(t *testing.T, db database.DB, blacklistSymbol model.BlacklistSymbol) model.BlacklistSymbol {
	res := model.BlacklistSymbol{}
	if err := table.BlacklistSymbol.
		INSERT(table.BlacklistSymbol.MutableColumns).
		RETURNING(table.BlacklistSymbol.AllColumns).
		MODEL(blacklistSymbol).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedSymbolScore(t *testing.T, db database.DB, symbol model.SymbolScore) model.SymbolScore {
	res := model.SymbolScore{}
	if err := table.SymbolScore.
		INSERT(table.SymbolScore.MutableColumns).
		RETURNING(table.SymbolScore.AllColumns).
		MODEL(symbol).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedScoreGroup(t *testing.T, db database.DB, symbol model.ScoreGroup) model.ScoreGroup {
	res := model.ScoreGroup{}
	if err := table.ScoreGroup.
		INSERT(table.ScoreGroup.MutableColumns).
		RETURNING(table.ScoreGroup.AllColumns).
		MODEL(symbol).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedScoreGroupInterest(t *testing.T, db database.DB, symbol model.ScoreGroupInterest) model.ScoreGroupInterest {
	res := model.ScoreGroupInterest{}
	if err := table.ScoreGroupInterest.
		INSERT(table.ScoreGroupInterest.MutableColumns).
		RETURNING(table.ScoreGroupInterest.AllColumns).
		MODEL(symbol).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedLoanPackageRequest(t *testing.T, db database.DB, symbol model.LoanPackageRequest) model.LoanPackageRequest {
	res := model.LoanPackageRequest{}
	if symbol.AssetType == "" {
		symbol.AssetType = model.AssetType_Underlying
	}
	if err := table.LoanPackageRequest.
		INSERT(table.LoanPackageRequest.MutableColumns).
		RETURNING(table.LoanPackageRequest.AllColumns).
		MODEL(symbol).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedLoanPackageOffer(t *testing.T, db database.DB, offer model.LoanPackageOffer) model.LoanPackageOffer {
	res := model.LoanPackageOffer{}
	if err := table.LoanPackageOffer.
		INSERT(table.LoanPackageOffer.MutableColumns).
		RETURNING(table.LoanPackageOffer.AllColumns).
		MODEL(offer).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedLoanPackageOfferInterest(t *testing.T, db database.DB, offerInterest model.LoanPackageOfferInterest) model.LoanPackageOfferInterest {
	res := model.LoanPackageOfferInterest{}
	if offerInterest.AssetType == "" {
		offerInterest.AssetType = model.AssetType_Underlying
	}
	if err := table.LoanPackageOfferInterest.
		INSERT(table.LoanPackageOfferInterest.MutableColumns).
		RETURNING(table.LoanPackageOfferInterest.AllColumns).
		MODEL(offerInterest).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedOfflineOfferUpdate(t *testing.T, db database.DB, offlineOfferUpdate model.OfflineOfferUpdate) model.OfflineOfferUpdate {
	res := model.OfflineOfferUpdate{}
	if err := table.OfflineOfferUpdate.
		INSERT(table.OfflineOfferUpdate.MutableColumns).
		RETURNING(table.OfflineOfferUpdate.AllColumns).
		MODEL(offlineOfferUpdate).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedInvestor(t *testing.T, db database.DB, investor model.Investor) model.Investor {
	res := model.Investor{}
	if err := table.Investor.
		INSERT(table.Investor.AllColumns.Except(table.Investor.CreatedAt, table.Investor.UpdatedAt)).
		RETURNING(table.Investor.AllColumns).
		MODEL(investor).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedInvestorAccount(t *testing.T, db database.DB, account model.InvestorAccount) model.InvestorAccount {
	res := model.InvestorAccount{}
	if err := table.InvestorAccount.
		INSERT(
			table.InvestorAccount.AllColumns.Except(
				table.InvestorAccount.CreatedAt, table.InvestorAccount.UpdatedAt,
			),
		).
		RETURNING(table.InvestorAccount.AllColumns).
		MODEL(account).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedLoanPolicyTemplate(t *testing.T, db database.DB, template model.LoanPolicyTemplate) model.LoanPolicyTemplate {
	res := model.LoanPolicyTemplate{}
	if err := table.LoanPolicyTemplate.
		INSERT(
			table.LoanPolicyTemplate.AllColumns.Except(
				table.LoanPolicyTemplate.CreatedAt, table.LoanPolicyTemplate.UpdatedAt,
			),
		).
		RETURNING(table.LoanPolicyTemplate.AllColumns).
		MODEL(template).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedSuggestedOffer(t *testing.T, db database.DB, offer model.SuggestedOffer) model.SuggestedOffer {
	res := model.SuggestedOffer{}
	if err := table.SuggestedOffer.
		INSERT(table.SuggestedOffer.AllColumns.Except(table.SuggestedOffer.CreatedAt, table.SuggestedOffer.UpdatedAt)).
		RETURNING(table.SuggestedOffer.AllColumns).
		MODEL(offer).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedSubmissionSheetMetadata(t *testing.T, db database.DB, metadata model.SubmissionSheetMetadata) model.SubmissionSheetMetadata {
	res := model.SubmissionSheetMetadata{}
	if err := table.SubmissionSheetMetadata.
		INSERT(
			table.SubmissionSheetMetadata.AllColumns.Except(
				table.SubmissionSheetMetadata.CreatedAt, table.SubmissionSheetMetadata.UpdatedAt,
			),
		).
		RETURNING(table.SubmissionSheetMetadata.AllColumns).
		MODEL(metadata).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedSubmissionSheetDetail(t *testing.T, db database.DB, detail model.SubmissionSheetDetail) model.SubmissionSheetDetail {
	res := model.SubmissionSheetDetail{}
	if err := table.SubmissionSheetDetail.
		INSERT(
			table.SubmissionSheetDetail.AllColumns.Except(
				table.SubmissionSheetDetail.CreatedAt, table.SubmissionSheetDetail.UpdatedAt,
			),
		).
		RETURNING(table.SubmissionSheetDetail.AllColumns).
		MODEL(detail).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedSuggestedOfferConfig(t *testing.T, db database.DB, config model.SuggestedOfferConfig) model.SuggestedOfferConfig {
	res := model.SuggestedOfferConfig{}
	if err := table.SuggestedOfferConfig.
		INSERT(table.SuggestedOfferConfig.MutableColumns).
		RETURNING(table.SuggestedOfferConfig.AllColumns).
		MODEL(config).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedConfiguration(t *testing.T, db database.DB, config model.FinancialConfiguration) model.FinancialConfiguration {
	res := model.FinancialConfiguration{}
	if err := table.FinancialConfiguration.
		INSERT(
			table.FinancialConfiguration.AllColumns.Except(
				table.FinancialConfiguration.CreatedAt, table.FinancialConfiguration.UpdatedAt,
			),
		).
		RETURNING(table.FinancialConfiguration.AllColumns).
		MODEL(config).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}

func SeedPromotionCampaign(t *testing.T, db database.DB, campaign model.PromotionCampaign) model.PromotionCampaign {
	res := model.PromotionCampaign{}
	if err := table.PromotionCampaign.
		INSERT(
			table.PromotionCampaign.AllColumns.Except(
				table.PromotionCampaign.CreatedAt, table.PromotionCampaign.UpdatedAt,
			),
		).
		RETURNING(table.PromotionCampaign.AllColumns).
		MODEL(campaign).Query(db, &res); err != nil {
		t.Error(err)
	}
	return res
}
