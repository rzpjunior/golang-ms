// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package branch

import (
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
	r.GET("", h.read, auth.Authorized("olt_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("olt_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.POST("", h.create, auth.Authorized("olt_crt"))
	r.PUT("/:id", h.update, auth.Authorized("olt_upd"))
	r.PUT("/salesperson/:id", h.updateSalesPerson, auth.Authorized("olt_upd_sps"))
	r.PUT("/archetype/:id", h.convertArchetype, auth.Authorized("olt_cvt_arc"))
	r.PUT("/archive/:id", h.archive, auth.Authorized("olt_arc"))
	r.PUT("/unarchive/:id", h.unarchive, auth.Authorized("olt_urc"))
	r.GET("/export", h.exportTemplate, auth.Authorized("sps_blk_exp"))
	r.POST("/salesperson/bulk", h.uploadTemplateUpdateBulkSalesperson, auth.Authorized("sps_blk_upl"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Branch
	var total int64

	if data, total, e = repository.GetBranchs(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetBranch("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Branch
	var total int64

	if data, total, e = repository.GetFilterBranchs(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

//create : function to create new data based on input
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

// updateSalesPerson : function to update data sales person
func (h *Handler) updateSalesPerson(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updatesalespersonRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UpdateSalesPerson(r)
			}
		}
	}

	return ctx.Serve(e)
}

// convertArchetype : function to convert archetype
func (h *Handler) convertArchetype(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r convertarchetypeRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = ConvertArchetype(r)
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

// read : function to get requested data based on parameters
func (h *Handler) exportTemplate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	isExport := ctx.QueryParam("export") == "1"
	cond := make(map[string]interface{})

	salesGroupID, _ := common.Decrypt(ctx.QueryParam("sales_group_id"))
	salesPersonID, _ := common.Decrypt(ctx.QueryParam("salesperson_id"))
	if salesGroupID != 0 {
		cond["sg.id = "] = salesGroupID
	}
	if salesPersonID != 0 {
		cond["s.id = "] = salesPersonID
	}

	data, e := getBranchFilterBySalesperson(cond)
	if e == nil {
		if isExport {
			var file string
			if file, e = getBranchFilterBySalespersonXls(data); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// uploadTemplate: create sales assignment based on xlxs file
func (h *Handler) uploadTemplateUpdateBulkSalesperson(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateBulkSalespersonReq

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			e = UpdateBulkSalesPerson(r)
		}
	}

	return ctx.Serve(e)
}
