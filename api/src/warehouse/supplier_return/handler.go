// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier_return

import (
	"net/http"
	"strings"

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
	r.GET("", h.read, auth.Authorized("supr_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("supr_rdd"))
	r.POST("", h.create, auth.Authorized("supr_crt"))
	r.PUT("/:id", h.update, auth.Authorized("supr_upd"))
	r.PUT("/confirm/:id", h.confirm, auth.Authorized("supr_cnf"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("supr_can"))
	r.GET("/print/:id", h.receivePrint, auth.Authorized("supr_prt"))

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

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if ctx.ResponseData, e = repository.GetSupplierReturn("id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.SupplierReturn
	var total int64

	if data, total, e = repository.GetSupplierReturns(rq); e != nil {
		return ctx.Serve(e)
	}
	ctx.Data(data, total)
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

// cancel : function to confirm data
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

func (h *Handler) receivePrint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var sr *model.SupplierReturn
	var id int64
	configs := make(map[string]string)
	req := make(map[string]interface{})
	var config []*model.ConfigApp
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		e = echo.NewHTTPError(http.StatusUnprocessableEntity)
		return ctx.Serve(e)
	}

	if id, e = ctx.Decrypt("id"); e != nil {
		e = echo.NewHTTPError(http.StatusUnprocessableEntity)
		return ctx.Serve(e)
	}
	if sr, e = repository.GetSupplierReturn("id", id); e != nil {
		e = echo.NewHTTPError(http.StatusUnprocessableEntity)
		return ctx.Serve(e)
	}

	if config, _, e = repository.GetConfigAppsByAttribute("attribute__icontains", "company"); e != nil {
		e = echo.NewHTTPError(http.StatusUnprocessableEntity)
		return ctx.Serve(e)
	}
	for _, v := range config {
		configs[strings.TrimPrefix(v.Attribute, "company_")] = v.Value
	}
	req["session"] = session.Staff.ID + 56
	req["sr"] = sr
	configs["address"] = strings.ReplaceAll(configs["address"], "<br>", "\n")
	req["config"] = configs

	file := util.SendPrint(req, "read/sr")
	ctx.Files(file)

	// delta print
	sr.DeltaPrint = sr.DeltaPrint + 1
	sr.Save("DeltaPrint")

	return ctx.Serve(e)
}
