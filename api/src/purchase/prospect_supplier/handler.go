package prospect_supplier

import (
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("pro_sup_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("pro_sup_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.POST("", h.create)
	r.PUT("/register/:id", h.register, auth.Authorized("pro_sup_reg"))
	r.PUT("/decline/:id", h.decline, auth.Authorized("pro_sup_dec"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.ProspectSupplier
	var total int64

	if data, total, e = repository.GetProspectSuppliers(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get requested data based on parameters
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		ctx.ResponseData, e = repository.GetProspectSupplier("id", id)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.ProspectSupplier
	var total int64

	if data, total, e = repository.GetFilterProspectSuppliers(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

//register : function to register data based on input
func (h *Handler) register(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r registerRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Register(r)
			}
		}
	}

	return ctx.Serve(e)
}

//decline : function to decline data based on input
func (h *Handler) decline(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r declineRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Decline(r)
			}
		}
	}

	return ctx.Serve(e)
}

//create : function to create new data based on input
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	//if r.Session, e = auth.UserSession(ctx); e == nil {
	if e = ctx.Bind(&r); e == nil {
		ctx.ResponseData, e = Save(r)
	}
	//}

	return ctx.Serve(e)
}
