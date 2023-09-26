// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package debit_note

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("dn_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("dn_rdd"))
	r.GET("/print/:id", h.receivePrint, auth.Authorized("dn_prt"))

}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if ctx.ResponseData, e = repository.GetDebitNote("id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.DebitNote
	var total int64

	if data, total, e = repository.GetDebitNotes(rq); e != nil {
		return ctx.Serve(e)
	}
	ctx.Data(data, total)
	return ctx.Serve(e)

}

func (h *Handler) receivePrint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var dn *model.DebitNote
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
	if dn, e = repository.GetDebitNote("id", id); e != nil {
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
	req["dn"] = dn
	configs["address"] = strings.ReplaceAll(configs["address"], "<br>", "\n")
	req["config"] = configs

	file := util.SendPrint(req, "read/dn")
	ctx.Files(file)

	return ctx.Serve(e)
}
