// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package field_purchaser

import (
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
	r.GET("", h.readFieldPurchaser, auth.AuthorizedFieldPurchaserMobile())
	r.POST("/order", h.createOrder, auth.AuthorizedFieldPurchaserMobile())
	r.GET("/:id", h.detailFieldPurchaser, auth.AuthorizedFieldPurchaserMobile())
	r.GET("/order/:id", h.detailOrder, auth.AuthorizedFieldPurchaserMobile())
	r.PUT("/order/:id", h.updateOrder, auth.AuthorizedFieldPurchaserMobile())
}

// get purchase order list for field purchaser app
func (h *Handler) readFieldPurchaser(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PurchaseOrder
	var total int64
	var warehouseID int64

	warehouse := ctx.QueryParam("warehouse")

	if warehouse != "" {
		if warehouseID, e = common.Decrypt(warehouse); e != nil {
			return ctx.Serve(e)
		}
	}

	if data, total, e = repository.GetPurchaseOrdersFieldPurchaser(rq, warehouseID); e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}

// createOrder : function to create new data to table field_purchase_order
func (h *Handler) createOrder(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createOrderRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = SaveOrder(r)

	return ctx.Serve(e)
}

// detailFieldPurchaser : function to get detailed data by id
func (h *Handler) detailFieldPurchaser(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetPurchaseOrderFieldPurchaser("id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// detailOrder : function to get detailed data by id
func (h *Handler) detailOrder(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetFieldPurchaseOrder("id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// updateOrder : function to update requested data based on parameters
func (h *Handler) updateOrder(c echo.Context) (e error) {
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
