// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receipt

import (
	"strings"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("gr_rdl"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("gr_rdd"))
	r.GET("/item", h.readItem, auth.Authorized("gr_rdl"))
	r.GET("/item/filter", h.readItemFilter, auth.Authorized("filter_rdl"))
	r.POST("", h.create, auth.Authorized("gr_crt"))
	r.PUT("/:id", h.update, auth.Authorized("gr_upd"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("gr_can"))
	r.PUT("/confirm/:id", h.confirm, auth.Authorized("gr_cnf"))
	r.GET("/print/:id", h.receivePrint, auth.Authorized("gr_prt"))
	r.PUT("/lock/:id", h.lock, auth.Authorized("gr_upd"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.GoodsReceipt
	var total int64

	if data, total, e = repository.GetGoodsReceipts(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.GoodsReceipt
	var total int64

	if data, total, e = repository.GetFilterGoodsReceipts(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetGoodsReceiptWithProductGroup("id", id); e != nil {
			e = echo.ErrNotFound
		}
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

// update : function to unarchive requested data based on parameters
func (h *Handler) update(c echo.Context) (e error) {
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

// cancel : function to cancel datas
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

// confirm : function to confirm data
func (h *Handler) confirm(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r confirmRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = Confirm(r)

	return ctx.Serve(e)
}

// readItem : function to get requested item data based on parameters
func (h *Handler) readItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.GoodsReceiptItem
	var total int64

	if data, total, e = repository.GetGoodsReceiptItems(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readItem : function to get requested item data based on parameters with filtered permission
func (h *Handler) readItemFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.GoodsReceiptItem
	var total int64

	if data, total, e = repository.GetFilterGoodsReceiptItems(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

func (h *Handler) receivePrint(c echo.Context) (e error) {

	ctx := c.(*cuxs.Context)
	var gr *model.GoodsReceipt
	var id int64
	configs := make(map[string]string)
	req := make(map[string]interface{})

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if gr, e = repository.GetGoodsReceiptWithProductGroup("id", id); e != nil {
		e = echo.ErrNotFound
	} else {
		req["gr"] = gr
		if config, _, e := repository.GetConfigAppsByAttribute("attribute__icontains", "company"); e == nil {
			for _, v := range config {
				configs[strings.TrimPrefix(v.Attribute, "company_")] = v.Value
			}
			configs["address"] = strings.ReplaceAll(configs["address"], "<br>", "\n")
			req["config"] = configs
		} else {
			e = echo.ErrNotFound
		}

		file := util.SendPrint(req, "read/gr")
		ctx.Files(file)
	}

	return ctx.Serve(e)
}

// lock : function to lock GR
func (h *Handler) lock(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r lockRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = Lock(r)

	return ctx.Serve(e)
}
