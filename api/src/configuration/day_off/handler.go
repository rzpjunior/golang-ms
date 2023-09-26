package day_off

import (
	"time"

	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read)
	r.GET("/:id", h.detail)
	r.POST("", h.create, auth.Authorized("dyo_crt"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.DayOff
	var total int64

	if data, total, e = repository.GetDayOffs(rq); e == nil {
		if len(data) == 0 {
			date, _ := time.Parse("2006-01-02", rq.Conditions[0]["off_date"])
			if int(date.Weekday()) == 0 {
				dayOff := new(model.DayOff)
				dayOff.OffDate = date
				dayOff.Note = "Sunday"
				data = append(data, dayOff)
			}
		}

		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetDayOff("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// create : insert day off
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Save(r)

	return ctx.Serve(e)
}
