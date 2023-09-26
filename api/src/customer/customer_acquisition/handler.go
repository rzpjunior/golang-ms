package customer_acquisition

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("ca_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("ca_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.CustomerAcquisition
	var total int64

	if data, total, e = repository.GetCustomerAcquisitions(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.CustomerAcquisition
	var total int64

	if data, total, e = repository.GetFilterCustomerAcquisition(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetCustomerAcquisition("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}
