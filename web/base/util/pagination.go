package util

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"noname001/app/base/store"
)

type WebPaginationParams struct {
	CurrPage         int
	PerPage          int
	TotalPage        int
	TotalData        int

	SearchPagination *store.SearchPagination

	MidLinksToShow   int

	LinkURI          string
	QueryParams      map[string]string
	Anchor           string
}

type WebPagination struct {
	ShownData string
	TotalData string

	FirstLink WebPaginationLink
	LastLink  WebPaginationLink
	MidLinks  []WebPaginationLink
	PrevLink  WebPaginationLink
	NextLink  WebPaginationLink
}

type WebPaginationLink struct {
	Name     string
	Link     string
	IsActive bool
}

func (util *Util) ExtractPaginationFromQueryParams(ctx *fasthttp.RequestCtx) (int, int) {
	currPage := ctx.QueryArgs().GetUintOrZero("page")
	perPage  := ctx.QueryArgs().GetUintOrZero("perpage")

	if currPage <= 0 { currPage = 1 }
	if perPage <= 0  { perPage = 10 }

	return currPage, perPage
}

func (util *Util) EmptyWebPaginationParams() (*WebPaginationParams) {
	return &WebPaginationParams{}
}

func (util *Util) GenerateWebPagination(params *WebPaginationParams) (*WebPagination) {
	var (
		firstPage      int = 1
		lastPage       int = params.TotalPage

		currPage       int = params.CurrPage
		perPage        int = params.PerPage
		totalData      int = params.TotalData

		midLinksToShow int = params.MidLinksToShow

		linkURI     string            = params.LinkURI
		queryParams map[string]string = params.QueryParams
		anchor      string            = params.Anchor

		prevPage, nextPage int
		midPoint, currStart, currEnd int
		midPages []int
	)

	// === DEFAULTS ===
	if sp := params.SearchPagination; sp != nil {
		currPage  = sp.CurrPage
		perPage   = sp.PerPage
		lastPage  = sp.TotalPage
		totalData = sp.TotalData
	}

	if midLinksToShow == 0 {
		midLinksToShow = 5
	}
	// === DEFAULTS ===

	prevPage = currPage - 1
	nextPage = currPage + 1
	
	if prevPage < firstPage { prevPage = firstPage }
	if nextPage > lastPage  { nextPage = lastPage  }

	if lastPage > midLinksToShow {
		midPoint  = midLinksToShow / 2
		currStart = currPage - (midPoint + 1)
		currEnd   = currPage + (midPoint)

		if currStart < 0 {
			currStart = 0
			currEnd   = 0 + midLinksToShow
		}
		if currEnd > lastPage {
			currStart = lastPage - midLinksToShow
			currEnd   = lastPage
		}

		midPages = make([]int, 0, midLinksToShow)
		for i := range (midLinksToShow) {
			midPages = append(midPages, currStart + i + 1)
		}
	} else {
		midPages = make([]int, 0, lastPage)
		for i := range (lastPage) {
			midPages = append(midPages, i + 1)
		}
	}

	anchorage := ""
	if anchor != "" { anchorage = fmt.Sprintf("#%s", anchor) }

	currQP := ""
	if queryParams != nil {
		for qpKey, qpVal := range queryParams {
			currQP += fmt.Sprintf("&%s=%s", qpKey, qpVal)
		}
	}

	pg := &WebPagination{
		TotalData: fmt.Sprintf("%d", totalData),
	}
	pg.FirstLink = WebPaginationLink{
		Link: fmt.Sprintf("%s?page=%d&perpage=%d%s%s", linkURI, firstPage, perPage, currQP, anchorage),
		IsActive: currPage == firstPage,
	}
	pg.PrevLink = WebPaginationLink{
		Link: fmt.Sprintf("%s?page=%d&perpage=%d%s%s", linkURI, prevPage, perPage, currQP, anchorage),
		IsActive: currPage == prevPage,
	}
	pg.MidLinks = make([]WebPaginationLink, 0, len(midPages))
	for _, pgNum := range midPages {
		pg.MidLinks = append(pg.MidLinks, WebPaginationLink{
			Name: fmt.Sprintf("%d", pgNum),
			Link: fmt.Sprintf("%s?page=%d&perpage=%d%s%s", linkURI, pgNum, perPage, currQP, anchorage),
			IsActive: currPage == pgNum,
		})
	}
	pg.NextLink = WebPaginationLink{
		Link: fmt.Sprintf("%s?page=%d&perpage=%d%s%s", linkURI, nextPage, perPage, currQP, anchorage),
		IsActive: currPage == nextPage,
	}
	pg.LastLink = WebPaginationLink{
		Link: fmt.Sprintf("%s?page=%d&perpage=%d%s%s", linkURI, lastPage, perPage, currQP, anchorage),
		IsActive: currPage == lastPage,
	}

	return pg
}
