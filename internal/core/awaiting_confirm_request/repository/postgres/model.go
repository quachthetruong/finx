package postgres

import (
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

type AwaitingConfirmRequest struct {
	model.LoanPackageRequest
	model.LoanPackageOffer
	model.LatestOfferUpdate
	model.Symbol
	model.Investor
	LoanPackageIds string `alias:"r.loan_package_ids"`
}
