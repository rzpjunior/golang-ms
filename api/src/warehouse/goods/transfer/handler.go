// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package transfer

import (
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/util"
	"github.com/labstack/echo/v4"
	"strings"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("gt_rdl"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("gt_rdd"))
	r.POST("", h.create, auth.Authorized("gt_req"))
	r.PUT("/:id", h.update, auth.Authorized("gt_upd"))
	r.PUT("/commit/:id", h.commit, auth.Authorized("gt_cmt"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("gt_can"))
	r.PUT("/confirm/:id", h.confirm, auth.Authorized("gt_cnf"))
	r.GET("/print/:id", h.receivePrint, auth.Authorized("gt_prt"))
	r.PUT("/lock/:id", h.lock, auth.Authorized("gt_cmt"))
}

func (h *Handler) receivePrint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var gt *model.GoodsTransfer
	var id int64
	var session *auth.SessionData
	configs := make(map[string]string)
	req := make(map[string]interface{})
	if session, e = auth.UserSession(ctx); e == nil {
		if id, e = ctx.Decrypt("id"); e == nil {
			if gt, e = repository.GetGoodsTransfer("id", id); e != nil {
				e = echo.ErrNotFound
			} else {
				req["gt"] = gt
				req["session"] = session.Staff.ID + 56

				if config, _, e := repository.GetConfigAppsByAttribute("attribute__icontains", "company"); e == nil {
					for _, v := range config {
						configs[strings.TrimPrefix(v.Attribute, "company_")] = v.Value
					}
					configs["address"] = strings.ReplaceAll(configs["address"], "<br>", "\n")
					req["config"] = configs
				} else {
					e = echo.ErrNotFound
				}

				file := util.SendPrint(req, "read/gt")
				ctx.Files(file)

			}
		}
	}

	return ctx.Serve(e)
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.GoodsTransfer
	var total int64

	if data, total, e = repository.GetGoodsTransfers(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.GoodsTransfer
	var total int64

	if data, total, e = repository.GetFilterGoodsTransfers(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetGoodsTransfer("id", id); e != nil {
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

// update : function to unarchive requested data based on parameters
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

// update : function to commit requested data based on parameters
func (h *Handler) commit(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r commitRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)

	}
	ctx.ResponseData, e = Commit(r)

	return ctx.Serve(e)

}

// cancel : function to cancel goods receipt
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = Cancel(r)

	return ctx.Serve(e)
}

// confirm : function to confirm delivery order
func (h *Handler) confirm(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r confirmRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Confirm(r)
			}
		}
	}

	return ctx.Serve(e)
}

// lock : function to lock GT
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
