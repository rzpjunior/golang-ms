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

type SalesAssignmentSubmissionHandler struct {
	Option                            global.HandlerOptions
	ServicesCustomerAcquisition       service.ICustomerAcquisitionService
	ServicesSalesAssignmentSubmission service.ISalesAssignmentSubmissionService
	ServicesSalesAssignmentItem       service.ISalesAssignmentItemService
}

// URLMapping implements router.RouteHandlers
func (h *SalesAssignmentSubmissionHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCustomerAcquisition = service.NewCustomerAcquisitionService()
	h.ServicesSalesAssignmentSubmission = service.NewSalesAssignmentSubmissionService()
	h.ServicesSalesAssignmentItem = service.NewSalesAssignmentItemService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
}

func (h SalesAssignmentSubmissionHandler) Get(c echo.Context) (err error) {
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
	salespersonID := ctx.GetParamString("salesperson_id")
	submitDateFrom := ctx.GetParamDate("submit_date_from")
	submitDateTo := ctx.GetParamDate("submit_date_to")
	outOfRoute := ctx.GetParamInt("out_of_route")
	task := ctx.GetParamInt("task")

	if task == 3 {
		var customerAcquisitions []*dto.CustomerAcquisitionResponse
		var total int64
		customerAcquisitions, total, err = h.ServicesCustomerAcquisition.GetSubmissions(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, territoryID, salespersonID, submitDateTo, submitDateFrom)
		if err != nil {
			h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
			return ctx.Serve(err)
		}

		ctx.DataList(customerAcquisitions, total, page)
	} else {
		var salesAssignmentSubmissions []*dto.SalesAssignmentSubmissionResponse
		var total int64
		salesAssignmentSubmissions, total, err = h.ServicesSalesAssignmentSubmission.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, territoryID, salespersonID, submitDateFrom, submitDateTo, task, outOfRoute)
		if err != nil {
			h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
			return ctx.Serve(err)
		}

		ctx.DataList(salesAssignmentSubmissions, total, page)
	}

	return ctx.Serve(err)
}

func (h SalesAssignmentSubmissionHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	task := ctx.GetParamInt("task")

	var salesAssignmentItem dto.SalesAssignmentItemResponse
	var customerAcquisition dto.CustomerAcquisitionResponse

	if task == 3 {
		// task = 3, get data from customer acquisition
		customerAcquisition, err = h.ServicesCustomerAcquisition.GetByID(ctx.Request().Context(), id)
		if err != nil {
			h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
			return ctx.Serve(err)
		}

		ctx.ResponseData = customerAcquisition
	} else {
		// task = 1, get data from sales assignment item
		salesAssignmentItem, err = h.ServicesSalesAssignmentItem.GetByID(ctx.Request().Context(), id)
		if err != nil {
			h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
			return ctx.Serve(err)
		}

		ctx.ResponseData = salesAssignmentItem
	}

	return ctx.Serve(err)
}
