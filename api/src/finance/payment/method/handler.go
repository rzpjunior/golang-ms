// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package method

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
	r.GET("", h.read, auth.Authorized("pym_rdl"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("pym_rdd"))
	r.GET("/field_purchaser", h.readFieldPurchaser, auth.AuthorizedFieldPurchaserMobile())

	// r.POST("", h.create, auth.Authorized("py_mtd_crt"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PaymentMethod
	var total int64

	if data, total, e = repository.GetPaymentMethods(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PaymentMethod
	var total int64

	if data, total, e = repository.GetFilterPaymentMethods(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// // create : function to create new data based on input
// func (h *Handler) create(c echo.Context) (e error) {
// 	ctx := c.(*cuxs.Context)
// 	var r createRequest

// 	if r.Session, e = auth.UserSession(ctx); e == nil {
// 		if e = ctx.Bind(&r); e == nil {
// 			ctx.ResponseData, e = Save(r)
// 		}
// 	}

// 	return ctx.Serve(e)
// }

// detail : function to get requested data based on parameters
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		ctx.ResponseData, e = repository.GetPaymentMethod("id", id)
	}

	return ctx.Serve(e)
}

// readFieldPurchaser : function to get requested data based on parameters
func (h *Handler) readFieldPurchaser(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PaymentMethod
	var total int64

	if data, total, e = repository.GetPaymentMethodsFieldPurchaser(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}
