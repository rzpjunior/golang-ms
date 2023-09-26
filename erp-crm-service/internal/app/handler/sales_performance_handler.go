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

type SalesPerformanceHandler struct {
	Option                   global.HandlerOptions
	ServicesSalesPerformance service.ISalesPerformanceService
}

// URLMapping implements router.RouteHandlers
func (h *SalesPerformanceHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesSalesPerformance = service.NewSalesPerformanceService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
}

func (h SalesPerformanceHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	territoryID := ctx.GetParamString("territory_id")
	salespersonID := ctx.GetParamString("salesperson_id")
	startDateFrom := ctx.GetParamDate("start_date_from")
	startDateTo := ctx.GetParamDate("start_date_to")

	var salesPerformances []*dto.SalesPerformanceResponse
	var total int64
	salesPerformances, total, err = h.ServicesSalesPerformance.Get(ctx.Request().Context(), page.Start, page.Limit, territoryID, salespersonID, startDateFrom, startDateTo)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(salesPerformances, total, page)

	return ctx.Serve(err)
}

func (h SalesPerformanceHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var salesPerformanceDetail *dto.SalesPerformanceDetailResponse

	id := c.Param("id")

	startDateFrom := ctx.GetParamDate("start_date_from")
	startDateTo := ctx.GetParamDate("start_date_to")
	task := ctx.GetParamInt("task")
	status := ctx.GetParamInt("status")

	salesPerformanceDetail, err = h.ServicesSalesPerformance.GetByID(ctx.Request().Context(), id, status, startDateFrom, startDateTo, task)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = salesPerformanceDetail

	return ctx.Serve(err)
}
