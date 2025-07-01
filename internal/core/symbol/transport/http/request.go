package http

import (
	"financing-offer/internal/core"
	"financing-offer/internal/core/entity"
	"financing-offer/pkg/optional"
)

type SymbolRequest struct {
	StockExchangeId int64            `json:"stockExchangeId" binding:"required"`
	Symbol          string           `json:"symbol" form:"symbol" binding:"required"`
	AssetType       entity.AssetType `json:"assetType,omitempty" form:"assetType" binding:"required,oneof=UNDERLYING DERIVATIVE"`
}

func (r SymbolRequest) toEntity(id int64) entity.Symbol {
	return entity.Symbol{
		Id:              id,
		StockExchangeId: r.StockExchangeId,
		Symbol:          r.Symbol,
		AssetType:       r.AssetType,
	}
}

type GetSymbolsRequest struct {
	Paging            core.Paging
	StockExchangeCode string `form:"stockExchangeCode"`
	AssetType         string `form:"assetType"`
}

func (r *GetSymbolsRequest) toFilter() entity.SymbolFilter {
	return entity.SymbolFilter{
		Paging:            r.Paging,
		StockExchangeCode: optional.FromValueNonZero(r.StockExchangeCode),
		AssetType:         optional.FromValueNonZero(r.AssetType),
	}
}

type UpdateSymbolStatusRequest struct {
	Status entity.SymbolStatus `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}
