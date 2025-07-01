package postgres

import (
	"github.com/go-jet/jet/v2/postgres"

	"financing-offer/internal/core"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

func MapSymbolsDbToEntity(symbols []model.Symbol) []entity.Symbol {
	res := make([]entity.Symbol, 0, len(symbols))
	for _, v := range symbols {
		res = append(res, MapSymbolDbToEntity(v))
	}
	return res
}

func MapSymbolDbToEntity(symbol model.Symbol) entity.Symbol {
	return entity.Symbol{
		Id:              symbol.ID,
		StockExchangeId: symbol.StockExchangeID,
		Symbol:          symbol.Symbol,
		AssetType:       entity.AssetType(symbol.AssetType),
		Status:          entity.SymbolStatusFromString(symbol.Status),
		LastUpdatedBy:   symbol.LastUpdatedBy,
		CreatedAt:       symbol.CreatedAt,
		UpdatedAt:       symbol.UpdatedAt,
	}
}

func MapSymbolEntityToDb(symbol entity.Symbol) model.Symbol {
	return model.Symbol{
		ID:              symbol.Id,
		StockExchangeID: symbol.StockExchangeId,
		Symbol:          symbol.Symbol,
		AssetType:       model.AssetType(symbol.AssetType),
		Status:          symbol.Status.String(),
		LastUpdatedBy:   symbol.LastUpdatedBy,
		CreatedAt:       symbol.CreatedAt,
		UpdatedAt:       symbol.UpdatedAt,
	}
}

func ApplySort(filter entity.SymbolFilter) []postgres.OrderByClause {
	expr := make([]postgres.OrderByClause, 0, len(filter.Sort))
	for _, s := range filter.Sort {
		var column postgres.Column
		for _, c := range table.Symbol.AllColumns {
			if c.Name() == s.ColumnName {
				column = c
				break
			}
		}
		if column == nil {
			continue
		}
		if s.Direction == core.DirectionAsc {
			expr = append(expr, column.ASC())
		} else {
			expr = append(expr, column.DESC())
		}
	}
	return expr
}

func ApplyFilter(filter entity.SymbolFilter) postgres.BoolExpression {
	expr := postgres.Bool(true)
	if filter.AssetType.IsPresent() {
		expr = expr.AND(table.Symbol.AssetType.EQ(postgres.NewEnumValue(filter.AssetType.Get())))
	}
	if filter.StockExchangeCode.IsPresent() {
		expr = expr.AND(table.StockExchange.Code.EQ(postgres.String(filter.StockExchangeCode.Get())))
	}
	return expr
}
