// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package merchant

import (
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("main_olt_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("main_olt_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.PUT("/:id", h.update, auth.Authorized("main_olt_upd"))
	r.PUT("/archive/:id", h.archive, auth.Authorized("main_olt_arc"))
	r.PUT("/unarchive/:id", h.unarchive, auth.Authorized("main_olt_urc"))
	// r.GET("/get/custom", h.getByQuery, auth.Authorized("main_olt_rdl"))
	r.PUT("/tag/:id", h.updateTag, auth.Authorized("main_olt_upd_cust_tag"))
	r.PUT("/phone_number/:id", h.updatePhoneNumber, auth.Authorized("main_olt_upd_pho_num"))
	r.POST("/suspension", h.suspension, auth.Authorized("main_olt_sus"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Merchant
	var total int64

	if data, total, e = repository.GetMerchants(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetMerchant("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Merchant
	var total int64

	if data, total, e = repository.GetFilterMerchants(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// update : function to update data
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

// archive : function to set status of active data into archive
func (h *Handler) archive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r archiveRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Archive(r)
			}
		}
	}

	return ctx.Serve(e)
}

// unarchive : function to set status of archive data into active
func (h *Handler) unarchive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r unarchiveRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Unarchive(r)
			}
		}
	}

	return ctx.Serve(e)
}

// func (h *Handler) getByQuery(c echo.Context) (e error) {
// 	ctx := c.(*cuxs.Context)
// 	rq := ctx.RequestQuery()
// 	var productID int64
// 	param := ctx.QueryParams()
// 	productName := param.Get("product_name")
// 	iD, _ := strconv.Atoi(param.Get("id"))
// 	productCode := param.Get("product_code")
// 	if iD > 0 {
// 		productID, _ = common.Decrypt(iD)
// 	}

// 	data, total, e := GetMerchantByQuery(rq, productName, productCode, productID)
// 	if e == nil {
// 		ctx.Data(data, total)
// 	}

// 	return ctx.Serve(e)
// }

// updateTag : function to update data tag customer
func (h *Handler) updateTag(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateTagRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UpdateTag(r)
			}
		}
	}

	return ctx.Serve(e)
}

// updatePhoneNumber : function to update phone number of merchant
func (h *Handler) updatePhoneNumber(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updatePhoneRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UpdatePhoneNumber(r)
			}
		}
	}

	return ctx.Serve(e)
}

// suspension : function to suspend or un-suspend of merchant
func (h *Handler) suspension(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r suspensionRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = Suspension(r)
		}
	}

	return ctx.Serve(e)
}
