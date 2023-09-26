// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package agent

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

// Handler collection handler for auth.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.Authorized("agt_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("agt_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.GET("/shipping/address/:id", h.getShippingAddress, auth.Authorized("agt_rdd_shp_adr"))

	r.POST("", h.create, auth.Authorized("agt_crt"))
	r.POST("/shipping/address", h.createShippingAddress, auth.Authorized("agt_crt_shp_adr"))

	r.PUT("/:id", h.update, auth.Authorized("agt_upd"))
	r.PUT("/archive/:id", h.archive, auth.Authorized("agt_arc"))
	r.PUT("/unarchive/:id", h.unarchive, auth.Authorized("agt_urc"))
	r.PUT("/tag/:id", h.updateTagCustomer, auth.Authorized("agt_upd_cust_tag"))
	r.PUT("/salesperson/:id", h.updateSalesperson, auth.Authorized("agt_upd_sps"))
	r.PUT("/archetype/:id", h.updateArchetype, auth.Authorized("agt_cvt_arc"))
	r.PUT("/shipping/address/:id", h.updateShippingAddress, auth.Authorized("agt_upd_shp_adr"))
	r.PUT("/phonenumber/:id", h.updatePhoneNumber, auth.Authorized("agt_upd_pho_num"))
}

func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	data, total, e := repository.GetBranchs(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	data, total, e := repository.GetFilterBranchs(ctx.RequestQuery())
	if e == nil {
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

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, _, e = repository.GetBranchsByMerchantId(id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// update : function to update data based on input
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

// updateTagCustomer : function to update tag customer data based on input
func (h *Handler) updateTagCustomer(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateTagRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UpdateTagCustomer(r)
			}
		}
	}

	return ctx.Serve(e)
}

// updateSalesperson : function to update archetype data based on input
func (h *Handler) updateSalesperson(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateSalespersonRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UpdateSalesperson(r)
			}
		}
	}

	return ctx.Serve(e)
}

// updateArchetype : function to update archetype data based on input
func (h *Handler) updateArchetype(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateArchetypeRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UpdateArchetype(r)
			}
		}
	}

	return ctx.Serve(e)
}

// createShippingAddress : function to create new shipping address
func (h *Handler) createShippingAddress(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createShippingAddressRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = SaveShippingAddress(r)
		}
	}

	return ctx.Serve(e)
}

// updateShippingAddress : function to update shipping address data based on input
func (h *Handler) updateShippingAddress(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateShippingAddressRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UpdateShippingAddress(r)
			}
		}
	}

	return ctx.Serve(e)
}

// getShippingAddress : function to get detailed data of shipping address by id
func (h *Handler) getShippingAddress(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetBranch("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// updatePhoneNumber : function to update phone number data based on input
func (h *Handler) updatePhoneNumber(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updatePhoneNumber

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UpdatePhoneNumber(r)
			}
		}
	}

	return ctx.Serve(e)
}
