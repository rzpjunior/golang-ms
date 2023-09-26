// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"github.com/labstack/echo/v4"
)

// Handler collection handler for auth.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.Authorized("prd_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("prd_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))

	r.POST("", h.create, auth.Authorized("prd_crt"))
	r.POST("/print", h.postPrint, auth.Authorized("prd_prt_lbl"))

	r.PUT("/:id", h.update, auth.Authorized("prd_upd"))
	r.PUT("/archive/:id", h.archive, auth.Authorized("prd_arc"))
	r.PUT("/unarchive/:id", h.unarchive, auth.Authorized("prd_urc"))
	r.PUT("/salability/:id", h.salability, auth.Authorized("prd_upd_sal"))
	r.PUT("/purchasability/:id", h.purchasability, auth.Authorized("prd_upd_pur"))
	r.PUT("/storability/:id", h.storability, auth.Authorized("prd_upd_sto"))
	r.PUT("/packable/:id", h.packable, auth.Authorized("prd_set_pac"))
	r.PUT("/unpackable/:id", h.unpackable, auth.Authorized("prd_set_upc"))
	r.PUT("/fragile/:id", h.fragile, auth.Authorized("prd_set_frg"))
	r.PUT("/notfragile/:id", h.notFragile, auth.Authorized("prd_set_ufg"))
}

func (h *Handler) postPrint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r printRequest
	var resp ResponsePrint
	req := make(map[string]interface{})

	if _, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			req["lp"] = r

			file := util.SendPrint(req, "read/label_product")
			resp.LinkPrint = file
			resp.TotalPrint = r.TotalPrint
			ctx.ResponseData = resp

		}
	}

	return ctx.Serve(e)
}

func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	tagProduct := ctx.QueryParam("tagproduct")

	data, total, e := repository.GetProducts(ctx.RequestQuery(), tagProduct)
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	data, total, e := repository.GetFilterProducts(ctx.RequestQuery())
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
		if ctx.ResponseData, e = repository.GetProduct("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

//update : function to update data based on input
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

//salability : function to update salability data based on input
func (h *Handler) salability(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r salabilityRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Salability(r)
			}
		}
	}

	return ctx.Serve(e)
}

//purchasability : function to update purchasability data based on input
func (h *Handler) purchasability(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r purchasabilityRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Purchasability(r)
			}
		}
	}

	return ctx.Serve(e)
}

//storability : function to update storability data based on input
func (h *Handler) storability(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r storabilityRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Storability(r)
			}
		}
	}

	return ctx.Serve(e)
}

//packable : function to update packability data based on input
func (h *Handler) packable(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r packableRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Packable(r)
			}
		}
	}

	return ctx.Serve(e)
}

//unpackable : function to update packability data based on input
func (h *Handler) unpackable(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r unpackableRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Unpackable(r)
			}
		}
	}

	return ctx.Serve(e)
}

//fragile : function to update fragile data based on input
func (h *Handler) fragile(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r fragileRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Fragile(r)

	return ctx.Serve(e)
}

//notFragile : function to update fragile data based on input
func (h *Handler) notFragile(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r notFragileRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = NotFragile(r)

	return ctx.Serve(e)
}

type ResponsePrint struct {
	LinkPrint  string  `json:"link_print"`
	TotalPrint float64 `json:"total_print"`
}
