package entity

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"

	string_helper "financing-offer/pkg/string-helper"
)

type FinancialAccountDetail struct {
	Id              string `json:"id"`
	Custody         string `json:"custody"`
	AccountNo       string `json:"accountNo"`
	AccountTypeName string `json:"accountTypeName"`
	FullName        string `json:"fullName"`
	Status          string `json:"status"`
	InvestorId      string `json:"investorId"`
	CustomerId      string `json:"customerId"`
	MarginAccount   bool   `json:"marginAccount"`
	AccountType     string `json:"accountType"`
}

type FinancialProductLoanPackage struct {
	Id                       int64           `json:"id"`
	Name                     string          `json:"name"`
	InitialRate              decimal.Decimal `json:"initialRate"`
	InterestRate             decimal.Decimal `json:"interestRate"`
	Term                     int             `json:"term"`
	BuyingFeeRate            decimal.Decimal `json:"brokerFirmBuyingFeeRate"`
	LoanBasketId             int64           `json:"loanBasketId"`
	InitialRateForWithdraw   decimal.Decimal `json:"initialRateForWithdraw"`
	MaintenanceRate          decimal.Decimal `json:"maintenanceRate"`
	LiquidRate               decimal.Decimal `json:"liquidRate"`
	PreferentialPeriod       int             `json:"preferentialPeriod"`
	PreferentialInterestRate decimal.Decimal `json:"preferentialInterestRate"`
	AllowExtendLoanTerm      bool            `json:"allowExtendLoanTerm"`
	AllowEarlyPayment        bool            `json:"allowEarlyPayment"`
	LoanType                 string          `json:"loanType"`
	BrokerFirmSellingFeeRate decimal.Decimal `json:"brokerFirmSellingFeeRate"`
	TransferFee              decimal.Decimal `json:"transferFee"`
	Description              string          `json:"description"`
}

type FinancialProductLoanPolicy struct {
	Id                       int64           `json:"id"`
	Name                     string          `json:"name"`
	Source                   string          `json:"source"`
	InterestRate             decimal.Decimal `json:"interestRate"`
	InterestBasis            int             `json:"interestBasis"`
	Term                     int             `json:"term"`
	OverdueInterest          decimal.Decimal `json:"overdueInterest"`
	AllowExtendLoanTerm      bool            `json:"allowExtendLoanTerm"`
	AllowEarlyPayment        bool            `json:"allowEarlyPayment"`
	PreferentialPeriod       int             `json:"preferentialPeriod"`
	PreferentialInterestRate decimal.Decimal `json:"preferentialInterestRate"`
	CreatedDate              time.Time       `json:"createdDate"`
	ModifiedDate             time.Time       `json:"modifiedDate"`
}

type FinancialProductLoanPackageDerivative struct {
	Id          int64           `json:"id"`
	Name        string          `json:"name"`
	InitialRate decimal.Decimal `json:"initialRate"`
}

type LoanRate struct {
	Id                     int64           `json:"id"`
	Name                   string          `json:"name"`
	InitialRate            decimal.Decimal `json:"initialRate"`
	InitialRateForWithdraw decimal.Decimal `json:"initialRateForWithdraw"`
	MaintenanceRate        decimal.Decimal `json:"maintenanceRate"`
	LiquidRate             decimal.Decimal `json:"liquidRate"`
}

type MarginPool struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	PoolGroupId int64  `json:"poolGroupId"`
	Type        string `json:"type"`
	Status      string `json:"status"`

	Group MarginPoolGroup `json:"group"`
}

type MarginPoolGroup struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Source string `json:"source"`
}

type AssignmentState struct {
	Submission Submission `json:"submission"`
}

type Submission struct {
	Symbol                     string                    `json:"symbol"`
	LoanPackageOfferInterestId int64                     `json:"loanPackageOfferInterestId"`
	LoanPackageRequestId       int64                     `json:"loanPackageRequestId"`
	AccountNo                  string                    `json:"accountNo"`
	LoanRate                   LoanRate                  `json:"loanRate"`
	Templates                  []TemplateWithProductRate `json:"templates"`
	FirmBuyingFeeRate          decimal.Decimal           `json:"firmBuyingFeeRate"`
	FirmSellingFeeRate         decimal.Decimal           `json:"firmSellingFeeRate"`
	TransferFee                decimal.Decimal           `json:"transferFee"`
	ProductCategoryId          int64                     `json:"productCategoryId"`
}

type MarginProductFilter struct {
	Symbol string
}

func (marginProductFilter MarginProductFilter) String() string {
	marshaled, _ := json.Marshal(marginProductFilter)
	return string_helper.BytesToString(marshaled)
}

type MarginProduct struct {
	Id           int64               `json:"id"`
	Name         string              `json:"name"`
	Symbol       string              `json:"symbol"`
	LoanRateId   int64               `json:"loanRateId"`
	LoanRate     LoanRate            `json:"loanRate"`
	LoanPolicies []LoanProductPolicy `json:"loanPolicies"`
}

type LoanProductPolicy struct {
	Rate            decimal.Decimal            `json:"rate"`
	RateForWithdraw decimal.Decimal            `json:"rateForWithdraw"`
	LoanPolicyId    int64                      `json:"loanPolicyId"`
	LoanPolicy      FinancialProductLoanPolicy `json:"loanPolicy"`
}

type MarginBasketFilter struct {
	Ids []int64 `json:"ids"`
}

type MarginBasket struct {
	Id             int64           `json:"id"`
	Name           string          `json:"name"`
	Symbols        []string        `json:"symbols"`
	LoanProductIds []int64         `json:"loanProductIds"`
	LoanProducts   []MarginProduct `json:"loanProducts"`
}
