package entity

import "github.com/shopspring/decimal"

type AccountLoanPackage struct {
	Id                       int64           `json:"id"`
	Name                     string          `json:"name"`
	Type                     string          `json:"type"`
	BrokerFirmBuyingFeeRate  decimal.Decimal `json:"brokerFirmBuyingFeeRate"`
	BrokerFirmSellingFeeRate decimal.Decimal `json:"brokerFirmSellingFeeRate"`
	TransferFee              decimal.Decimal `json:"transferFee"`
	Description              string          `json:"description"`
	LoanProducts             []LoanProduct   `json:"loanProducts"`
	BasketId                 int64           `json:"basketId"`
}

type LoanProduct struct {
	Id                       string          `json:"id"`
	Name                     string          `json:"name"`
	Symbol                   string          `json:"symbol"`
	InitialRate              decimal.Decimal `json:"initialRate"`
	InitialRateForWithdraw   decimal.Decimal `json:"initialRateForWithdraw"`
	MaintenanceRate          decimal.Decimal `json:"maintenanceRate"`
	LiquidRate               decimal.Decimal `json:"liquidRate"`
	InterestRate             decimal.Decimal `json:"interestRate"`
	PreferentialPeriod       int             `json:"preferentialPeriod"`
	PreferentialInterestRate decimal.Decimal `json:"preferentialInterestRate"`
	Term                     int             `json:"term"`
	AllowExtendLoanTerm      bool            `json:"allowExtendLoanTerm"`
	AllowEarlyPayment        bool            `json:"allowEarlyPayment"`
}

type AccountLoanPackageWithSymbol struct {
	AccountLoanPackage
	Symbol string `json:"symbol"`
}

type AccountLoanPackageWithAccountNo struct {
	AccountLoanPackage
	AccountNo string `json:"accountNo"`
}

type LoanPackageWithCampaignProduct struct {
	Id                       int64                 `json:"id"`
	Name                     string                `json:"name"`
	Type                     string                `json:"type"`
	BrokerFirmBuyingFeeRate  decimal.Decimal       `json:"brokerFirmBuyingFeeRate"`
	BrokerFirmSellingFeeRate decimal.Decimal       `json:"brokerFirmSellingFeeRate"`
	TransferFee              decimal.Decimal       `json:"transferFee"`
	Description              string                `json:"description"`
	BasketId                 int64                 `json:"basketId"`
	CampaignProducts         []CampaignWithProduct `json:"campaignProducts"`
}

type Campaign struct {
	Name        string `json:"name"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
}

type CampaignWithProduct struct {
	Product  LoanProduct `json:"product"`
	Campaign Campaign    `json:"campaign"`
}
