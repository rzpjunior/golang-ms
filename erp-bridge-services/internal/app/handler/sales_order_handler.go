package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/service"
	"github.com/labstack/echo/v4"
)

type SalesOrderHandler struct {
	Option             global.HandlerOptions
	ServicesSalesOrder service.ISalesOrderService
}

// URLMapping implements router.RouteHandlers
func (h *SalesOrderHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesSalesOrder = service.NewSalesOrderService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
}

func (h SalesOrderHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	// get params filters
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")
	addressID := ctx.GetParamInt("address_id")
	customerID := ctx.GetParamInt("customer_id")
	salespersonID := ctx.GetParamInt("salesperson_id")
	orderDateFrom := ctx.GetParamDate("order_date_from")
	orderDateTo := ctx.GetParamDate("order_date_to")

	var salesOrders []dto.SalesOrderResponse
	var total int64
	salesOrders, total, err = h.ServicesSalesOrder.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, int64(addressID), int64(customerID), int64(salespersonID), orderDateFrom, orderDateTo)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(salesOrders, total, page)

	return ctx.Serve(err)
}

func (h SalesOrderHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var salesOrder dto.SalesOrderResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	salesOrder, err = h.ServicesSalesOrder.GetDetail(ctx.Request().Context(), id, "")
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = salesOrder

	return ctx.Serve(err)
}
