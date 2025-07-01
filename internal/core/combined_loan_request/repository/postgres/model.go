package postgres

import (
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

type CombinedLoanRequest struct {
	model.LoanPackageRequest
	model.LoanPackageOffer
	model.Symbol
	model.LoanContract
	model.Investor
	AdminAssignedLoanPackageIds string `alias:"r.admin_assigned_loan_package_ids"`
	ActivatedLoanPackageIds     string `alias:"r.activated_loan_package_ids"`
	PackageCreatedTime          string `alias:"r.package_created_time"`
	Statuses                    string `alias:"r.statuses"`
	CancelledReasons            string `alias:"r.cancelled_reasons"`
}
