package mapper

import (
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

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

func MapLoanPackageOffersDbToEntity(offers []model.LoanPackageOffer) []entity.LoanPackageOffer {
	result := make([]entity.LoanPackageOffer, 0, len(offers))
	for _, o := range offers {
		result = append(result, MapLoanPackageOfferDbToEntity(o))
	}
	return result
}
