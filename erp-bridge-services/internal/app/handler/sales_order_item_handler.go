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

type SalesOrderItemHandler struct {
	Option                 global.HandlerOptions
	ServicesSalesOrderItem service.ISalesOrderItemService
}

// URLMapping implements router.RouteHandlers
func (h *SalesOrderItemHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesSalesOrderItem = service.NewSalesOrderItemService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
}

func (h SalesOrderItemHandler) Get(c echo.Context) (err error) {
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
	salesOrderID := ctx.GetParamInt("sales_order_id")
	itemID := ctx.GetParamInt("item_id")

	var territories []dto.SalesOrderItemResponse
	var total int64
	territories, total, err = h.ServicesSalesOrderItem.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, int64(salesOrderID), int64(itemID))
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(territories, total, page)

	return ctx.Serve(err)
}

func (h SalesOrderItemHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var salesOrderItem dto.SalesOrderItemResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	salesOrderItem, err = h.ServicesSalesOrderItem.GetDetail(ctx.Request().Context(), id, "")
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = salesOrderItem

	return ctx.Serve(err)
}
