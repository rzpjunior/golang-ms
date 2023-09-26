// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package payment

import (
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("sp_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("sp_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.POST("", h.create, auth.Authorized("sp_crt"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("sp_can"))
	r.PUT("/cancel/active/:id", h.cancelActive, auth.Authorized("sp_can_active"))
	r.PUT("/payment_proof/:id", h.addPaymentProof, auth.Authorized("sp_crt_active"))
	r.POST("/bulk", h.bulkCreate, auth.Authorized("sp_crt"))
	r.POST("/bulk/active", h.bulkCreateActive, auth.Authorized("sp_crt_active"))
	r.PUT("/bulk/confirm", h.bulkConfirm, auth.Authorized("sp_cnf"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var createdByFilter []int64

	createdByStr := ctx.QueryParam("created_by")
	if createdByStr != "" {
		createdByArr := strings.Split(createdByStr, ",")

		for _, v := range createdByArr {
			userId, _ := common.Decrypt(v)
			createdByFilter = append(createdByFilter, userId)
		}

	}

	var data []*model.SalesPayment
	var total int64

	if data, total, e = repository.GetSalesPayments(rq, createdByFilter...); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.SalesPayment
	var total int64

	if data, total, e = repository.GetFilterSalesPayments(rq); e == nil {
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

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetSalesPayment("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// cancel : function to unarchive requested data based on parameters
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Cancel(r)
			}
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) cancelActive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelActiveRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = CancelActive(r)
			}
		}
	}

	return ctx.Serve(e)
}

// bulkCreate : function to create several sales payment at the same time
func (h *Handler) bulkCreate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r bulkPaymentRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = BulkCreatePayment(r)
		}
	}

	return ctx.Serve(e)
}

// bulkCreate : function to create several active sales payment at the same time
func (h *Handler) bulkCreateActive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r bulkCreateActivePaymentRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = BulkCreateActivePayment(r)
		}
	}

	return ctx.Serve(e)
}

// bulkConfirm : function to create several sales payment at the same time
func (h *Handler) bulkConfirm(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r bulkConfirmPaymentRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = BulkConfirmPayment(r)
		}
	}

	return ctx.Serve(e)
}

// addPaymentProof : function to Add Payment Proof on Sales Payment
func (h *Handler) addPaymentProof(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r addPaymentProofRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = AddPaymentProof(r)

	return ctx.Serve(e)
}
