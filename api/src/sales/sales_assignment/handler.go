// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sales_assignment

import (
	"time"

	"git.edenfarm.id/project-version2/api/src/auth"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/export", h.exportTemplate, auth.Authorized("sla_exp"))
	r.POST("/upload", h.uploadTemplate, auth.Authorized("sla_upl"))
	r.GET("", h.read, auth.Authorized("sla_rdl"))
	r.GET("/item", h.detail, auth.Authorized("sla_rdd"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("sla_can"))
	r.PUT("/cancel/item/:id", h.cancelItem, auth.Authorized("sla_can"))
	r.GET("/submission", h.readListSubmission, auth.Authorized("sla_sub_rdl"))
	r.GET("/submission/:id", h.readDetailSubmission, auth.Authorized("sla_sub_rdd"))
	r.GET("/visit_tracker", h.readListVisitTracker, auth.Authorized("slp_rdd"))
	r.GET("/visit_tracker/:id", h.readVisitTrackerDetail, auth.Authorized("slp_rdd"))
}

// exportTemplate: Download template based on selected sales group
func (h *Handler) exportTemplate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	isExport := ctx.QueryParam("export") == "1"
	cond := make(map[string]interface{})

	salesGroupID, _ := common.Decrypt(ctx.QueryParam("sales_group_id"))

	if salesGroupID != 0 {
		cond["sg.id = "] = salesGroupID
	}

	data, e := getBranchBySalesGroup(cond)
	if e == nil {
		if isExport {
			var file string
			if file, e = getBranchBySalesGroupXls(data); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}
	return ctx.Serve(e)
}

// uploadTemplate: create sales assignment based on xlxs file
func (h *Handler) uploadTemplate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = Save(r)
		}
	}

	return ctx.Serve(e)
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var data []*model.SalesAssignment
	var total int64

	if data, total, e = repository.GetSalesAssignments(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by sales assignment id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var data []*model.SalesAssignmentItem
	var total int64

	if data, total, e = repository.GetSalesAssignmentItems(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// cancel : function to set status of active data into cancel
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

// archiveItem : function to set status of active data into cancel
func (h *Handler) cancelItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelItemRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = CancelItem(r)
			}
		}
	}

	return ctx.Serve(e)
}

// readListSubmission : function to get submission list
func (h *Handler) readListSubmission(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var total int64

	taskType := ctx.QueryParam("task")

	// task : Customer Acquisition
	if taskType == "3" {
		var data []*model.CustomerAcquisition

		if data, total, e = repository.GetSubmissionCA(rq); e == nil {
			ctx.Data(data, total)
		}

		// task : Visit or Follow Up
	} else {
		var data []*model.SalesAssignmentItem
		if taskType == "" {
			if data, total, e = repository.GetSubmissionVisitAndFollowUp(rq); e == nil {
				ctx.Data(data, total)
			}
		} else {
			if data, total, e = repository.GetSubmissionSA(rq, taskType); e == nil {
				ctx.Data(data, total)
			}
		}
	}

	return ctx.Serve(e)
}

// readDetailSubmission : function to get detail submission
func (h *Handler) readDetailSubmission(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64
	taskType := ctx.QueryParam("task")

	// task : Customer Acquisition
	if taskType == "3" {
		if id, e = ctx.Decrypt("id"); e == nil {
			if ctx.ResponseData, e = repository.GetSubmissionCustomerAcquisitionDetail("id", id); e != nil {
				e = echo.ErrNotFound
			}
		}
		// task : Visit or Follow Up
	} else {
		if id, e = ctx.Decrypt("id"); e == nil {
			if ctx.ResponseData, e = repository.GetSubmissionSalesAssignmentDetail("id", id); e != nil {
				e = echo.ErrNotFound
			}
		}
	}

	return ctx.Serve(e)
}

// readListVisitTracker : function to get visit tracker list
func (h *Handler) readListVisitTracker(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var total, salesGroupID, salesPersonID int64

	fromDateStr := ctx.QueryParam("fromdate")
	toDateStr := ctx.QueryParam("todate")
	salesGroupIdStr := ctx.QueryParam("salesgroup_id")
	salesPersonIdStr := ctx.QueryParam("salesperson_id")
	if salesGroupIdStr != "" {
		salesGroupID, e = common.Decrypt(salesGroupIdStr)
		if e != nil {
			return ctx.Serve(e)
		}
	}
	if salesPersonIdStr != "" {
		salesPersonID, e = common.Decrypt(salesPersonIdStr)
		if e != nil {
			return ctx.Serve(e)
		}
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")

	var data []*model.SalesAssignmentItem

	layout := "2006-01-02"
	fromDate, e := time.ParseInLocation(layout, fromDateStr, loc)
	if e != nil {
		return ctx.Serve(e)
	}
	toDate, e := time.ParseInLocation(layout, toDateStr, loc)
	if e != nil {
		return ctx.Serve(e)
	}

	if data, total, e = repository.GetSalesAssignmentItemsGroup(rq, fromDate, toDate, salesGroupID, salesPersonID); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readVisitTrackerDetail : function to get submission detail tracker
func (h *Handler) readVisitTrackerDetail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetStaff("id", id); e != nil {
			return ctx.Serve(e)
		}
	}

	fromDateStr := ctx.QueryParam("fromdate")
	toDateStr := ctx.QueryParam("todate")
	loc, _ := time.LoadLocation("Asia/Jakarta")

	layout := "2006-01-02"
	fromDate, e := time.ParseInLocation(layout, fromDateStr, loc)
	if e != nil {
		return ctx.Serve(e)
	}
	toDate, e := time.ParseInLocation(layout, toDateStr, loc)
	if e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetSalesAssignmentItemsTracker(id, fromDate, toDate); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}
