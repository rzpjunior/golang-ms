package edenlabs

import (
	"math"
)

type (
	// Pagination
	Paginator struct {
		Page     int  `json:"page"`
		PerPage  int  `json:"per_page"`
		Offset   int  `json:"offset,omitempty"`
		Limit    int  `json:"limit,omitempty"`
		NumPages int  `json:"total_pages"`
		Start    int  `json:"start,omitempty"`
		End      int  `json:"end,omitempty"`
		HasPrev  bool `json:"has_prev,omitempty"`
		HasNext  bool `json:"has_next,omitempty"`
	}
)

func NewPaginator(ctx *Context) (*Paginator, error) {
	page := ctx.GetPage()
	perPage := ctx.GetPerPage()

	if page < 1 {
		return nil, ErrorValidation("page", "Invalid page number")
	}

	p := &Paginator{
		Page:    page,
		PerPage: perPage,
		Offset:  (page - 1) * perPage,
		Limit:   perPage,
	}

	p.Start = p.Offset
	p.End = p.Start + perPage

	return p, nil
}

func (p *Paginator) Json(totalItems int64) (page *Paginator) {
	if p.Page > 1 {
		p.HasPrev = true
	}

	if p.End < int(totalItems) {
		p.HasNext = true
	}

	p.NumPages = int(math.Ceil(float64(totalItems) / float64(p.PerPage)))

	if p.Page > p.NumPages {
		return
	}

	page = p
	return
}
