package supplier

import (
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("sup_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("sup_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.POST("", h.create, auth.Authorized("sup_crt"))
	r.PUT("/:id", h.update, auth.Authorized("sup_upd"))
	r.PUT("/archive/:id", h.archive, auth.Authorized("sup_arc"))
	r.PUT("/unarchive/:id", h.unarchive, auth.Authorized("sup_urc"))
	r.POST("/field_purchaser", h.createFieldPurchaserApp, auth.AuthorizedFieldPurchaserMobile())
	r.GET("/field_purchaser", h.readFieldPurchaserApp, auth.AuthorizedFieldPurchaserMobile())
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Supplier
	var total int64

	if data, total, e = repository.GetSuppliers(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get requested data based on parameters
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		ctx.ResponseData, e = repository.GetSupplierDetail("id", id)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Supplier
	var total int64

	if data, total, e = repository.GetFilterSuppliers(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

//create : function to create new data based on input
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = Save(r)
		}
	}

	return ctx.Serve(e)
}

//update : function to update data based on input
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Update(r)
			}
		}
	}

	return ctx.Serve(e)
}

// archive : function to archive requested data based on parameters
func (h *Handler) archive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r archiveRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Archive(r)
			}
		}
	}

	return ctx.Serve(e)
}

// unarchive : function to unarchive requested data based on parameters
func (h *Handler) unarchive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r unarchiveRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UnArchive(r)
			}
		}
	}

	return ctx.Serve(e)
}

//createFieldPurchaserApp : function to create new data based on input
func (h *Handler) createFieldPurchaserApp(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createFieldPurchaserRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = SaveFieldPurchaser(r)

	return ctx.Serve(e)
}

// readFieldPurchaserApp : function to get requested data based on parameters for field purchaser apps
func (h *Handler) readFieldPurchaserApp(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Supplier
	var total int64

	if data, total, e = repository.GetSuppliersInFieldPurchaserApps(rq); e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}
