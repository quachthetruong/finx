package postgres

import (
	"github.com/volatiletech/null/v9"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

func MapBlacklistSymbolsDbToEntity(blacklists []model.BlacklistSymbol) []entity.BlacklistSymbol {
	res := make([]entity.BlacklistSymbol, 0, len(blacklists))
	for _, v := range blacklists {
		res = append(res, MapBlacklistSymbolDbToEntity(v))
	}
	return res
}

func MapBlacklistSymbolDbToEntity(blacklist model.BlacklistSymbol) entity.BlacklistSymbol {
	e := entity.BlacklistSymbol{
		Id:           blacklist.ID,
		SymbolId:     blacklist.SymbolID,
		AffectedFrom: blacklist.AffectedFrom,
		Status:       entity.BlacklistSymbolStatusFromString(blacklist.Status.String()),
		CreatedAt:    blacklist.CreatedAt,
		UpdatedAt:    blacklist.UpdatedAt,
	}
	if blacklist.AffectedTo.IsValid() {
		e.AffectedTo = blacklist.AffectedTo.Time
	}
	return e
}

func MapBlacklistSymbolEntityToDb(blacklist entity.BlacklistSymbol) model.BlacklistSymbol {
	affectTo := null.Time{}
	if !blacklist.AffectedTo.IsZero() {
		affectTo = null.TimeFrom(blacklist.AffectedTo)
	}
	return model.BlacklistSymbol{
		ID:           blacklist.Id,
		SymbolID:     blacklist.SymbolId,
		AffectedFrom: blacklist.AffectedFrom,
		AffectedTo:   affectTo,
		Status:       model.Blacklistsymbolstatus(blacklist.Status),
		CreatedAt:    blacklist.CreatedAt,
		UpdatedAt:    blacklist.UpdatedAt,
	}
}
