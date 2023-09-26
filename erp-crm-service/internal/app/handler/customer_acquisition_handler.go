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

type CustomerAcquisitionHandler struct {
	Option                      global.HandlerOptions
	ServicesCustomerAcquisition service.ICustomerAcquisitionService
}

// URLMapping implements router.RouteHandlers
func (h *CustomerAcquisitionHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCustomerAcquisition = service.NewCustomerAcquisitionService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
}

func (h CustomerAcquisitionHandler) Get(c echo.Context) (err error) {
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

	var customerAcquisitions []*dto.CustomerAcquisitionResponse
	var total int64
	customerAcquisitions, total, err = h.ServicesCustomerAcquisition.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, territoryID, salespersonID, submitDateFrom, submitDateTo)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(customerAcquisitions, total, page)

	return ctx.Serve(err)
}

func (h CustomerAcquisitionHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var customerAcquisition dto.CustomerAcquisitionResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	customerAcquisition, err = h.ServicesCustomerAcquisition.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = customerAcquisition

	return ctx.Serve(err)
}
