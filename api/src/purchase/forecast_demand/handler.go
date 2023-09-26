// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package forecast_demand

import (
	"time"

	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("frc_dmd_rdl"))
	r.PUT("", h.update, auth.Authorized("frc_dmd_upd"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var total int64

	isExport := ctx.QueryParam("export") == "1"

	if isExport {
		var data []orm.Params
		var arrDate []string

		if data, total, arrDate, e = repository.GetForecastDemandsForExport(rq); e == nil {
			var file string
			if file, e = DownloadForecastDemandXls(time.Now(), data, arrDate); e == nil {
				ctx.Files(file)
			}
		}
	} else {
		var data []*model.ForecastDemand
		if data, total, e = repository.GetForecastDemands(rq); e == nil {
			ctx.Data(data, total)
		}

	}

	return ctx.Serve(e)
}

// update : function to update data based on parameters
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			e = Update(r)
		}
	}

	return ctx.Serve(e)
}
