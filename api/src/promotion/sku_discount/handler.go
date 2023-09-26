// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sku_discount

import (
	"strconv"
	"time"

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
	r.GET("", h.read, auth.Authorized("sku_dsc_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("sku_dsc_rdd"))
	r.POST("", h.create, auth.Authorized("sku_dsc_crt"))
	r.PUT("/archive/:id", h.archive, auth.Authorized("sku_dsc_arc"))
	r.GET("/sku_disc_data", h.readDiscountData, auth.Authorized("sku_dsc_rdl"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var (
		data                               []*model.SkuDiscount
		total                              int64
		priceSetEnc, priceSetRea, priceSet string
	)

	priceSetEnc = ctx.QueryParam("price_set.e")
	if priceSetEnc == "" {
		priceSetRea = ctx.QueryParam("price_set")
		priceSet = priceSetRea
	} else {
		priceSet = common.Encrypt(priceSetEnc)
	}

	if data, total, e = repository.GetSkuDiscounts(rq, priceSet); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get requested data based on parameters
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		ctx.ResponseData, e = repository.GetSkuDiscount("id", id)
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

	if e = ctx.Bind(&r); e == nil {
		ctx.ResponseData, e = Save(r)
	}

	return ctx.Serve(e)
}

// archive : function to change status data to 2
func (h *Handler) archive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r archiveRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e == nil {
		ctx.ResponseData, e = Archive(r)
	}

	return ctx.Serve(e)
}

// readDiscountData : function to get requested data based on parameters
func (h *Handler) readDiscountData(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var (
		priceSetID, productID, merchantID, salesOrderItemID int64
		currDate                                            time.Time
		skuDiscountItem                                     *model.SkuDiscountItem
		salesOrderItem                                      *model.SalesOrderItem
		orderChannel                                        int8
	)

	if ctx.QueryParam("price_set_id") == "" {
		e = echo.ErrNotFound
		return ctx.Serve(e)
	}

	if ctx.QueryParam("product_id") == "" {
		e = echo.ErrNotFound
		return ctx.Serve(e)
	}

	if productID, e = common.Decrypt(ctx.QueryParam("product_id")); e != nil {
		return ctx.Serve(e)
	}

	if ctx.QueryParam("merchant_id") == "" {
		e = echo.ErrNotFound
		return ctx.Serve(e)
	}

	if merchantID, e = common.Decrypt(ctx.QueryParam("merchant_id")); e != nil {
		return ctx.Serve(e)
	}

	if _, e = repository.ValidMerchant(merchantID); e != nil {
		return ctx.Serve(e)
	}

	if priceSetID, e = common.Decrypt(ctx.QueryParam("price_set_id")); e != nil {
		return ctx.Serve(e)
	}

	if _, e = repository.ValidPriceSet(priceSetID); e != nil {
		return ctx.Serve(e)
	}

	if ctx.QueryParam("order_channel") == "" {
		e = echo.ErrNotFound
		return ctx.Serve(e)
	}

	ordChannel, e := strconv.Atoi(ctx.QueryParam("order_channel"))
	if e != nil {
		e = echo.ErrNotFound
		return ctx.Serve(e)
	}

	orderChannel = int8(ordChannel)
	currDate = time.Now().Local()

	if ctx.QueryParam("so_item_id") != "" {
		if salesOrderItemID, e = common.Decrypt(ctx.QueryParam("so_item_id")); e != nil {
			return ctx.Serve(e)
		}

		if salesOrderItem, e = repository.ValidSalesOrderItem(salesOrderItemID); e != nil {
			return ctx.Serve(e)
		}

		if e = salesOrderItem.SalesOrder.Read("ID"); e != nil {
			return ctx.Serve(e)
		}

		currDate = salesOrderItem.SalesOrder.CreatedAt
		if !salesOrderItem.SalesOrder.LastUpdatedAt.IsZero() {
			currDate = salesOrderItem.SalesOrder.LastUpdatedAt
		}

		orderChannel = salesOrderItem.SalesOrder.OrderChannel
	}

	if skuDiscountItem, e = repository.GetSkuDiscountData(merchantID, priceSetID, productID, salesOrderItemID, orderChannel, currDate); e == nil {
		ctx.ResponseData = skuDiscountItem
	}

	return ctx.Serve(e)
}
