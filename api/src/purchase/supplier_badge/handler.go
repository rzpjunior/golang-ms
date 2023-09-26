// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier_badge

import (
	"strconv"

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
	r.GET("", h.read, auth.Authorized("filter_rdl"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("sub_rdd"))
	r.POST("", h.create, auth.Authorized("sub_crt"))
	r.PUT("/:id", h.update, auth.Authorized("sub_upd"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.SupplierBadge
	var total int64

	var supplierCommodityID int64
	var paramSupplierCommodityId = ctx.QueryParam("supplier_commodity_id")

	if paramSupplierCommodityId != "" {
		supplierCommodityID, e = strconv.ParseInt(paramSupplierCommodityId, 10, 64)

		if e != nil {
			return ctx.Serve(e)
		}

		if supplierCommodityID != 0 {

			supplierCommodityID, e = common.Decrypt(paramSupplierCommodityId)
	
			if e != nil {
				return ctx.Serve(e)
			}
	
		}
	}


	if data, total, e = repository.GetSupplierBadges(rq, supplierCommodityID); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.SupplierBadge
	var total int64
	var supplierCommodityID int64
	var paramSupplierCommodityId = ctx.QueryParam("supplier_commodity_id")

	if paramSupplierCommodityId != "" {
		supplierCommodityID, e = strconv.ParseInt(paramSupplierCommodityId, 10, 64)

		if e != nil {
			return ctx.Serve(e)
		}

		if supplierCommodityID != 0 {

			supplierCommodityID, e = common.Decrypt(paramSupplierCommodityId)
	
			if e != nil {
				return ctx.Serve(e)
			}
	
		}
	}

	if data, total, e = repository.GetFilterSupplierBadges(rq, supplierCommodityID); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetSupplierBadge("id", id); e != nil {
			e = echo.ErrNotFound
		}
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
