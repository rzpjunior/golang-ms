package menu

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read)
	r.GET("/:id", h.detail)
	r.GET("/filter", h.readFilter)
	r.GET("/user/:id", h.menuByUserID)
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Menu
	var total int64

	if data, total, e = repository.GetMenus(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Menu
	var total int64

	if data, total, e = repository.GetFilterMenus(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetMenu("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// menuByUserID : function to get menu by user id
func (h *Handler) menuByUserID(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var data []*model.Menu
	if id, e = ctx.Decrypt("id"); e == nil {
		if data, e = repository.GetMenusByUserID(id); e == nil {
			ctx.Data(data)
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}
