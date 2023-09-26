// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoice

import (
	"strings"

	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"

	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("pi_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("pi_rdd"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("pi_can"))
	r.PUT("/:id", h.update, auth.Authorized("pi_upd"))
	r.PUT("/:id/tax-invoice", h.addTaxInvoice, auth.Authorized("pi_add_tax_invoice"))
	r.POST("", h.create, auth.Authorized("pi_crt"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PurchaseInvoice
	var total int64

	ataDateStr := ctx.QueryParam("ata_date")
	ataDateArr := strings.Split(ataDateStr, "|")

	if data, total, e = repository.GetPurchaseInvoices(rq, ataDateArr...); e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetPurchaseInvoice("id", id); e != nil {
		return ctx.Serve(e)
	}

	return ctx.Serve(e)
}

// create : function to create new data based on input
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

// update : function to update delivery order
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)

	}

	if r.PurchaseInvoice, e = repository.GetPurchaseInvoice("id", r.ID); e != nil {
		return ctx.Serve(e)
	}

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Update(r)

	return ctx.Serve(e)
}

// cancel : function to cancel purchase invoice
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Cancel(r)

	return ctx.Serve(e)
}

// addTaxInvoice : function to addTaxInvoice delivery order
func (h *Handler) addTaxInvoice(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r addTaxInvoiceRequest

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if r.PurchaseInvoice, e = repository.GetPurchaseInvoice("id", r.ID); e != nil {
		return ctx.Serve(e)
	}

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = AddTaxInvoice(r)

	return ctx.Serve(e)
}
