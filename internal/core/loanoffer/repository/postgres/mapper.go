package postgres

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/volatiletech/null/v9"

	"financing-offer/internal/core/entity"
	loanContractPostgres "financing-offer/internal/core/loancontract/repository/postgres"
	loanPackageOfferInterestPostgres "financing-offer/internal/core/loanofferinterest/repository/postgres"
	loanPackageRequestPostgres "financing-offer/internal/core/loanpackagerequest/repository/postgres"
	symbolPostgres "financing-offer/internal/core/symbol/repository/postgres"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/pkg/querymod"
)

type LoanPackageOfferWithRequest struct {
	model.LoanPackageOffer
	Symbol                    model.Symbol
	LoanPackageRequest        model.LoanPackageRequest
	LoanPackageOfferInterests []struct {
		model.LoanPackageOfferInterest
		LoanContract model.LoanContract
	}
}

func MapLoanPackageOfferWithRequestDbToEntity(o LoanPackageOfferWithRequest) entity.LoanPackageOffer {
	offerLines := make([]entity.LoanPackageOfferInterest, 0, len(o.LoanPackageOfferInterests))
	for _, l := range o.LoanPackageOfferInterests {
		i := loanPackageOfferInterestPostgres.MapLoanPackageOfferInterestDbToEntity(l.LoanPackageOfferInterest)
		if contract := loanContractPostgres.MapLoanContractDbToEntity(l.LoanContract); contract.Id != 0 {
			i.LoanContract = &contract
		}
		offerLines = append(offerLines, i)
	}
	res := entity.LoanPackageOffer{
		Id:                   o.ID,
		LoanPackageRequestId: o.LoanPackageRequestID,
		OfferedBy:            o.OfferedBy,
		CreatedAt:            o.CreatedAt,
		UpdatedAt:            o.UpdatedAt,
		ExpiredAt:            o.ExpiredAt.Time,
		FlowType:             entity.FlowTypeFromString(o.FlowType),

		LoanPackageOfferInterests: offerLines,
	}
	if symbol := symbolPostgres.MapSymbolDbToEntity(o.Symbol); symbol.Id != 0 {
		res.Symbol = &symbol
	}
	if request := loanPackageRequestPostgres.MapLoanPackageRequestDbToEntity(o.LoanPackageRequest); request.Id != 0 {
		res.LoanPackageRequest = &request
	}
	return res
}

func MapLoanPackageOfferWithRequestsDbToEntity(offers []LoanPackageOfferWithRequest) []entity.LoanPackageOffer {
	result := make([]entity.LoanPackageOffer, 0, len(offers))
	for _, o := range offers {
		result = append(result, MapLoanPackageOfferWithRequestDbToEntity(o))
	}
	return result
}

func MapLoanPackageOffersEntityToDb(offers []entity.LoanPackageOffer) []model.LoanPackageOffer {
	result := make([]model.LoanPackageOffer, 0, len(offers))
	for _, o := range offers {
		result = append(result, MapLoanPackageOfferEntityToDb(o))
	}
	return result
}

func MapLoanPackageOffersDbToEntity(offers []model.LoanPackageOffer) []entity.LoanPackageOffer {
	result := make([]entity.LoanPackageOffer, 0, len(offers))
	for _, o := range offers {
		result = append(result, MapLoanPackageOfferDbToEntity(o))
	}
	return result
}

func MapLoanPackageOfferEntityToDb(o entity.LoanPackageOffer) model.LoanPackageOffer {
	res := model.LoanPackageOffer{
		ID:                   o.Id,
		LoanPackageRequestID: o.LoanPackageRequestId,
		OfferedBy:            o.OfferedBy,
		CreatedAt:            o.CreatedAt,
		UpdatedAt:            o.UpdatedAt,
		FlowType:             o.FlowType.String(),
	}
	if !o.ExpiredAt.IsZero() {
		res.ExpiredAt = null.TimeFrom(o.ExpiredAt)
	}
	return res
}

func MapLoanPackageOfferDbToEntity(o model.LoanPackageOffer) entity.LoanPackageOffer {
	return entity.LoanPackageOffer{
		Id:                   o.ID,
		LoanPackageRequestId: o.LoanPackageRequestID,
		OfferedBy:            o.OfferedBy,
		CreatedAt:            o.CreatedAt,
		UpdatedAt:            o.UpdatedAt,
		ExpiredAt:            o.ExpiredAt.Time,
		FlowType:             entity.FlowTypeFromString(o.FlowType),
	}
}

func ApplyFilter(filter entity.LoanPackageOfferFilter) postgres.BoolExpression {
	expr := postgres.Bool(true)
	if filter.InvestorId != "" {
		expr = expr.AND(table.LoanPackageRequest.InvestorID.EQ(postgres.String(filter.InvestorId)))
	}
	if filter.Symbol.IsPresent() {
		expr = expr.AND(table.Symbol.Symbol.EQ(postgres.String(filter.Symbol.Get())))
	}
	if len(filter.OfferInterestStatus) > 0 {
		expr = expr.AND(
			table.LoanPackageOfferInterest.Status.IN(querymod.In(filter.OfferInterestStatus)...).
				OR(table.LoanPackageOfferInterest.Status.IS_NULL()),
		)
	}
	if filter.AssetType.IsPresent() {
		expr = expr.AND(table.LoanPackageRequest.AssetType.EQ(postgres.NewEnumValue(filter.AssetType.Get().String())))
	}
	return expr
}
