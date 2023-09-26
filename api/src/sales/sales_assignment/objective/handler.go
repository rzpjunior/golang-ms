package objective

import (
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("sla_obj_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("sla_obj_rdl"))
	r.GET("/filter", h.readFilter, auth.Authorized("sla_obj_rdl"))
	r.POST("", h.create, auth.Authorized("sla_obj_crt"))
	r.PUT("/:id", h.update, auth.Authorized("sla_obj_upd"))
	r.PUT("/archive/:id", h.archive, auth.Authorized("sla_obj_arc"))
	r.PUT("/unarchive/:id", h.unarchive, auth.Authorized("sla_obj_urc"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.SalesAssignmentObjective
	var total int64

	if data, total, e = repository.GetSalesAssignmentObjectives(rq); e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.SalesAssignmentObjective
	var total int64

	if data, total, e = repository.GetFilterSalesAssignmentObjectives(rq); e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}

//create : function to create new data based on input
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

//update : function to update data based on input
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Update(r)

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetSalesAssignmentObjective("id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// archive : function to set status of active data into archive
func (h *Handler) archive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r archiveRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Archive(r)

	return ctx.Serve(e)
}

// unarchive : function to set status of archive data into active
func (h *Handler) unarchive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r unarchiveRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = UnArchive(r)

	return ctx.Serve(e)
}
