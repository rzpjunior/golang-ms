package transfer_sku

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
	r.POST("", h.create, auth.Authorized("tfs_crt"))
	r.GET("", h.read, auth.Authorized("tfs_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("tfs_rdd"))
	r.GET("/filter/product_group", h.filter, auth.Authorized("filter_rdl"))
	r.PUT("/confirm/:id", h.confirm, auth.Authorized("tfs_cnf"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("tfs_can"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.TransferSku
	var total int64

	data, total, e = repository.GetListTransferSku(rq)

	if e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}

// create : function to create new data based on input
func (h *Handler) create(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, err = auth.UserSession(ctx); err != nil {
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&r); err != nil {
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = Save(r)

	return ctx.Serve(err)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = repository.GetTransferSku("id", id)
	if e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// confirm : function to confirm data
func (h *Handler) confirm(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r confirmRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Confirm(r)

	return ctx.Serve(e)
}

// filter : function to get requested data based on parameters with filtered query
func (h *Handler) filter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Stock
	var total int64

	data, total, e = repository.GetProductGroupTransferSku(rq)

	if e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}

// cancel : function to cancel data
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Cancel(r)

	return ctx.Serve(e)
}
