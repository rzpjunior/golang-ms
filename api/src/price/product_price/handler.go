package product_price

import (
	"time"

	"git.edenfarm.id/cuxs/common/now"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("prd_prc_rdl"))
	r.GET("/template", h.getProductPriceUpdateTemplate, auth.Authorized("prd_prc_upd"))
	r.GET("/template/shadow", h.getShadowPriceUpdateTemplate, auth.Authorized("prd_prc_upd"))
	r.GET("/export", h.exportTemplate, auth.Authorized("prd_prc_exp"))

	r.PUT("/template/update", h.update, auth.Authorized("prd_prc_upd"))
	r.PUT("/template/shadow/update", h.updateShadow, auth.Authorized("prd_prc_upd"))

	//r.GET("/print", h.getProductPrice, auth.Authorized("prd_prc_rdl"))
}

//read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Price
	var total int64
	tagProduct := ctx.QueryParam("tagproduct")

	if data, total, e = repository.GetPrices(rq, tagProduct); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// update: update product price based on xlxs file
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = Update(r)
		}
	}

	return ctx.Serve(e)
}

// update: update shadow price based on xlxs file
func (h *Handler) updateShadow(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r shadowRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = UpdateShadow(r)
		}
	}

	return ctx.Serve(e)
}

// getProductPriceUpdateTemplate: Download template regards on update product price
func (h *Handler) getProductPriceUpdateTemplate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var backdate time.Time

	isExport := ctx.QueryParam("export") == "1"
	backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
	data, total, e := getProductPrice(rq, "")
	if e == nil {
		if isExport {
			var file string
			if file, e = getProductPriceUpdateXls(backdate, data); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}

// getProductPriceUpdateTemplate: Download template regards on product price(export product price)
func (h *Handler) exportTemplate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var backdate time.Time

	isExport := ctx.QueryParam("export") == "1"
	backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
	data, total, e := getProductPrice(rq, "")

	if e == nil {
		if isExport {
			var file string
			if file, e = exportTemplateXls(backdate, data); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}

// getProductPriceUpdateTemplate: Download template regards on shadow product price
func (h *Handler) getShadowPriceUpdateTemplate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var backdate time.Time

	isExport := ctx.QueryParam("export") == "1"
	backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
	tagProduct := ctx.QueryParam("tagproduct")
	data, total, e := getProductPrice(rq, tagProduct)
	if e == nil {
		if isExport {
			var file string
			if file, e = getShadowPriceUpdateXls(backdate, data); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}
