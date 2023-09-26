// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchase_payment

import (
	"strconv"

	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"

	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("pp_rdl"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("pp_can"))
	r.POST("", h.create, auth.Authorized("pp_crt"))
	r.POST("/bulk", h.createBulk, auth.Authorized("pp_crt"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PurchasePayment
	var total int64

	if data, total, e = repository.GetPurchasePayments(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PurchasePayment
	var total int64

	if data, total, e = repository.GetFilterPurchasePayments(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// create : function to create new data based on input
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

// cancel : function to cancel purchase invoice
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if r.PurchasePayment, e = repository.ValidPurchasePayment(r.ID); e == nil {
				if e = ctx.Bind(&r); e == nil {
					ctx.ResponseData, e = Cancel(r)
				}
			}
		}
	}

	return ctx.Serve(e)
}

// createBulk : function to create bulky new data based on input
func (h *Handler) createBulk(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createBulkRequest
	var totalSuccess int64

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			totalSuccess, e = SaveBulk(r)
			ctx.ResponseData = strconv.Itoa(int(totalSuccess)) + " of " + strconv.Itoa(len(r.Data)) + " payment has been created successfully"
		}
	}

	return ctx.Serve(e)
}
