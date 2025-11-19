package store

import (
	"math"
)

// type OrderBy struct {}

type SearchPagination struct {
	CurrPage  int // given as params
	PerPage   int // given as params

	TotalPage int // provided from data query
	TotalData int // calculated with finalize()
}
func (pg *SearchPagination) Normalize() {
	if pg.CurrPage <= 0 {
		pg.CurrPage = 1
	}
	if pg.PerPage <= 0 {
		pg.PerPage = 10
	}
}
func (pg *SearchPagination) Finalize() {
	if pg.TotalData == 0 {
		pg.TotalPage = 1
	} else {
		tpF := float64(pg.TotalData) / float64(pg.PerPage)
		remains := math.Mod(tpF, 1.0)
		if remains > 0 { tpF = tpF + 1 }
		
		if tpF < 1 { tpF = 1 }

		pg.TotalPage = int(tpF)
	}
}
