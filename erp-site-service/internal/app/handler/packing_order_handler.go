package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type PackingOrderHandler struct {
	Option               global.HandlerOptions
	ServicesPackingOrder service.IPackingOrderService
}

// URLMapping implements router.RouteHandlers
func (h *PackingOrderHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesPackingOrder = service.NewPackingOrderService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("pc_rdl"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("pc_rdd"))
	r.POST("", h.Generate, cMiddleware.Authorized("pc_crt"))
	r.GET("/export/:id", h.Export, cMiddleware.Authorized("pc_exp"))
}

func (h PackingOrderHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	status := ctx.GetParamInt("status")
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	siteID := ctx.GetParamString("site_id")
	deliveryDateFrom := ctx.GetParamDate("delivery_date_from")
	deliveryDateTo := ctx.GetParamDate("delivery_date_to")

	var packingOrders []*dto.PackingOrderResponse
	var total int64
	packingOrders, total, err = h.ServicesPackingOrder.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, siteID, deliveryDateFrom, deliveryDateTo)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(packingOrders, total, page)

	return ctx.Serve(err)
}

func (h PackingOrderHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var packingOrder *dto.PackingOrderResponse
	packingOrder, err = h.ServicesPackingOrder.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = packingOrder

	return ctx.Serve(err)
}

func (h PackingOrderHandler) Generate(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var req dto.PackingOrderRequestGenerate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	err = h.ServicesPackingOrder.Generate(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}
func (h PackingOrderHandler) Export(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var packingOrderExport *dto.PackingOrderResponseExport
	packingOrderExport, err = h.ServicesPackingOrder.Export(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = packingOrderExport

	return ctx.Serve(err)
}
