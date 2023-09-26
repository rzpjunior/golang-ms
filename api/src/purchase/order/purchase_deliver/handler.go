// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchase_deliver

import (
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.readPurchaseDeliver, auth.AuthorizedFieldPurchaserMobile())
	r.GET("/:id", h.detailPurchaseDeliver, auth.AuthorizedFieldPurchaserMobile())
	r.POST("/signature", h.signPurchaseDeliver, auth.AuthorizedFieldPurchaserMobile())
	r.PUT("/print/:id", h.printCount, auth.AuthorizedFieldPurchaserMobile())
	r.POST("/consolidated", h.consolidatedPurchaseDeliver, auth.AuthorizedFieldPurchaserMobile())
	r.GET("/consolidated/:id", h.detailConsolidatedPurchaseDeliver, auth.AuthorizedFieldPurchaserMobile())
	r.GET("/consolidated", h.readConsolidatedPurchaseDeliver, auth.AuthorizedFieldPurchaserMobile())
	r.PUT("/consolidated/print/:id", h.printCountConsolidate, auth.AuthorizedFieldPurchaserMobile())
	r.GET("/filter", h.filterPurchaseDeliver, auth.AuthorizedFieldPurchaserMobile())
	r.POST("/consolidated/signature", h.signConsolidatedPurchaseDeliver, auth.AuthorizedFieldPurchaserMobile())
}

// readPurchaseDeliver : function to get requested data based on parameters
func (h *Handler) readPurchaseDeliver(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PurchaseDeliver
	var total int64
	var staffID int64

	staff := ctx.QueryParam("staff")

	if staff != "" {
		if staffID, e = common.Decrypt(staff); e != nil {
			return ctx.Serve(e)
		}
	}

	if data, total, e = repository.GetPurchaseDelivers(rq, staffID); e != nil {
		ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}

// detailPurchaseDeliver : function to get detailed data by id
func (h *Handler) detailPurchaseDeliver(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetPurchaseDeliver("id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// signPurchaseDeliver : function to sign purchase deliver
func (h *Handler) signPurchaseDeliver(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r signRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Sign(r)

	return ctx.Serve(e)
}

// printCount : function to count  copy of print
func (h *Handler) printCount(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r printRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Print(r)

	return ctx.Serve(e)
}

// consolidatedPurchaseDeliver : function to consolidate purchase_deliver
func (h *Handler) consolidatedPurchaseDeliver(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r consolidateRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Consolidate(r)

	return ctx.Serve(e)
}

// get detail consolidated purchase deliver
func (h *Handler) detailConsolidatedPurchaseDeliver(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetConsolidatedPurchaseDeliver("id", id); e != nil {
		return ctx.Serve(e)
	}

	return ctx.Serve(e)
}

// get consolidated purchase deliver list
func (h *Handler) readConsolidatedPurchaseDeliver(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.ConsolidatedPurchaseDeliver
	var total int64
	var warehouseID int64

	warehouse := ctx.QueryParam("warehouse")

	if warehouse != "" {
		if warehouseID, e = common.Decrypt(warehouse); e != nil {
			return ctx.Serve(e)
		}
	}

	if data, total, e = repository.GetConsolidatedPurchaseDelivers(rq, warehouseID); e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}

// printCountConsolidate : function to count  copy of print
func (h *Handler) printCountConsolidate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r printConsolidateRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = PrintConsolidate(r)

	return ctx.Serve(e)
}

// get filter purchase deliver
func (h *Handler) filterPurchaseDeliver(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PurchaseDeliver
	var total int64

	codeStr := ctx.QueryParam("code")
	codeArr := strings.Split(codeStr, "|")

	if data, total, e = repository.GetFilterPurchaseDelivers(rq, codeArr); e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}

// signConsolidatedPurchaseDeliver : function to sign consolidated purchase deliver
func (h *Handler) signConsolidatedPurchaseDeliver(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r signConsolidateRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = SignConsolidate(r)

	return ctx.Serve(e)
}
