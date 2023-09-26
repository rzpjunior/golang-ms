package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type SalesAssignmentHandler struct {
	Option                  global.HandlerOptions
	ServicesSalesAssignment service.ISalesAssignmentService
}

// URLMapping implements router.RouteHandlers
func (h *SalesAssignmentHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesSalesAssignment = service.NewSalesAssignmentService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
	r.GET("/export", h.Export, cMiddleware.Authorized())
	r.POST("/import", h.Import, cMiddleware.Authorized())
	r.PUT("/cancel/:id", h.CancelBatch, cMiddleware.Authorized())
	r.PUT("/cancel/item/:id", h.CancelItem, cMiddleware.Authorized())
}

func (h SalesAssignmentHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")
	territoryID := ctx.GetParamString("territory_id")
	startDateFrom := ctx.GetParamDate("start_date_from")
	startDateTo := ctx.GetParamDate("start_date_to")
	endDateFrom := ctx.GetParamDate("end_date_from")
	endDateTo := ctx.GetParamDate("end_date_to")

	var salesAssignments []*dto.SalesAssignmentResponse
	var total int64
	salesAssignments, total, err = h.ServicesSalesAssignment.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, territoryID, startDateFrom, startDateTo, endDateFrom, endDateTo)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(salesAssignments, total, page)

	return ctx.Serve(err)
}

func (h SalesAssignmentHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var salesAssignment dto.SalesAssignmentResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	finishDateFrom := ctx.GetParamDate("finish_date_from")
	finishDateTo := ctx.GetParamDate("finish_date_to")
	taskType := ctx.GetParamInt("task_type")
	status := ctx.GetParamInt("status")
	search := ctx.GetParamString("search")

	salesAssignment, err = h.ServicesSalesAssignment.GetByID(ctx.Request().Context(), id, status, search, taskType, finishDateFrom, finishDateTo)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = salesAssignment

	return ctx.Serve(err)
}

func (h SalesAssignmentHandler) Export(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params territory
	territoryID := ctx.GetParamString("territory_id")

	ctx.ResponseData, err = h.ServicesSalesAssignment.Export(ctx.Request().Context(), territoryID)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h SalesAssignmentHandler) Import(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.SalesAssignmentImportRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	err = h.ServicesSalesAssignment.Import(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h SalesAssignmentHandler) CancelBatch(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesSalesAssignment.CancelBatch(ctx.Request().Context(), id)

	return ctx.Serve(err)
}

func (h SalesAssignmentHandler) CancelItem(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesSalesAssignment.CancelItem(ctx.Request().Context(), id)

	return ctx.Serve(err)
}
