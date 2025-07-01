package core

import (
	"math"
)

type Direction string

const (
	DirectionAsc  Direction = "ASC"
	DirectionDesc Direction = "DESC"
)

type Order struct {
	Direction  Direction
	ColumnName string
}

type Orders []Order

func (oo *Orders) Contain(columnName string) bool {
	for _, o := range *oo {
		if o.ColumnName == columnName {
			return true
		}
	}
	return false
}

// Add merge given orders with skipping duplicated items
func (oo *Orders) Add(orders ...Order) {
	for i := range orders {
		order := &orders[i]
		if !oo.Contain(order.ColumnName) {
			*oo = append(*oo, *order)
		}
	}
}

// Paging request
type Paging struct {
	Sort   Orders
	Size   uint
	Number uint
}

func (p *Paging) Orders() Orders {
	return p.Sort
}

func (p *Paging) Limit() int64 {
	return int64(p.Size)
}

func (p *Paging) Offset() int64 {
	return int64((p.Number - 1) * p.Size)
}

func (p *Paging) TotalPages(totalRecords int64) uint {
	if p.Size == 0 {
		return 1
	}
	return uint(math.Ceil(float64(totalRecords) / float64(p.Size)))
}

type PagingMetaData struct {
	Total      int64 `json:"total"`
	PageSize   uint  `json:"pageSize"`
	PageNumber uint  `json:"pageNumber"`
	TotalPages uint  `json:"totalPages"`
}
