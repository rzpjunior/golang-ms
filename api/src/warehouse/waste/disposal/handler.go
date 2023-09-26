// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package disposal

import (
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/common/now"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("wd_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("wd_rdd"))
	r.GET("/export/form", h.exportForm, auth.Authorized("wd_exp_form"))
	r.POST("", h.create, auth.Authorized("wd_crt"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("wd_can"))
	r.PUT("/confirm/:id", h.confirm, auth.Authorized("wd_cnf"))
	r.GET("/print/:id", h.receivePrint, auth.Authorized("wd_prt"))
}

func (h *Handler) receivePrint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var wd *model.WasteDisposal
	var id int64
	var session *auth.SessionData
	configs := make(map[string]string)
	req := make(map[string]interface{})
	if session, e = auth.UserSession(ctx); e == nil {
		if id, e = ctx.Decrypt("id"); e == nil {
			if wd, e = repository.GetWasteDisposal("id", id); e != nil {
				e = echo.ErrNotFound
			} else {
				req["wd"] = wd
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

				file := util.SendPrint(req, "read/wd")
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

	var data []*model.WasteDisposal
	var total int64

	if data, total, e = repository.GetWasteDisposals(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetWasteDisposal("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.WasteDisposal
	var total int64

	if data, total, e = repository.GetFilterWasteDisposals(rq); e == nil {
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

//// update : function to unarchive requested data based on parameters
//func (h *Handler) update(c echo.Context) (e error) {
//	ctx := c.(*cuxs.Context)
//	var r updateRequest
//
//	if r.Session, e = auth.UserSession(ctx); e == nil {
//		if r.ID, e = ctx.Decrypt("id"); e == nil {
//			if e = ctx.Bind(&r); e != nil {
//				panic(e)
//			} else {
//				ctx.ResponseData, e = Update(r)
//			}
//		}
//	}
//
//	return ctx.Serve(e)
//}

func (h *Handler) exportForm(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var backdate time.Time
	var filter, exclude map[string]interface{}

	isExport := ctx.QueryParam("export") == "1"
	warehouseID, _ := common.Decrypt(ctx.QueryParam("warehouse_id"))
	warehouse, _ := repository.GetWarehouse("id", warehouseID)

	filter = map[string]interface{}{"warehouse_id": warehouseID, "waste_stock__gt": 0, "status": 1}

	backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
	data, total, e := repository.CheckStockData(filter, exclude)
	if e == nil {
		if isExport {
			var file string
			if file, e = ExportFormXls(backdate, data, warehouse); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}

// cancel : function to cancel data
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

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Confirm(r)
			}
		}
	}

	return ctx.Serve(e)
}
