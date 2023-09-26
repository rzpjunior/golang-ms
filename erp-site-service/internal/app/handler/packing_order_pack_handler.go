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

type PackingOrderPackHandler struct {
	Option                   global.HandlerOptions
	ServicesPackingOrderPack service.IPackingOrderPackService
}

// URLMapping implements router.RouteHandlers
func (h *PackingOrderPackHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesPackingOrderPack = service.NewPackingOrderPackService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("pc_rdl"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("pc_rdd"))
	r.PUT("/update/:id", h.Update, cMiddleware.Authorized("pc_upd"))
	r.PUT("/print/:id", h.Print, cMiddleware.Authorized("pc_prt"))
	r.DELETE("/dispose/:id", h.Dispose, cMiddleware.Authorized("pc_del"))
}

func (h PackingOrderPackHandler) Get(c echo.Context) (err error) {
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
	itemID := ctx.GetParamString("item_id")
	packingDateFrom := ctx.GetParamDate("packing_date_from")
	packingDateTo := ctx.GetParamDate("packing_date_to")

	var packingOrders []*dto.PackingOrderItemResponse
	var total int64
	packingOrders, total, err = h.ServicesPackingOrderPack.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, siteID, itemID, packingDateFrom, packingDateTo)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(packingOrders, total, page)

	return ctx.Serve(err)
}

func (h PackingOrderPackHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	packType := ctx.GetParamFloat64("pack_type")
	itemID := ctx.GetParamString("item_id")

	var packingOrder *dto.PackingOrderItemResponse
	packingOrder, err = h.ServicesPackingOrderPack.GetDetail(ctx.Request().Context(), id, packType, itemID)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = packingOrder

	return ctx.Serve(err)
}

func (h PackingOrderPackHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.PackingOrderItemRequestUpdate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesPackingOrderPack.Update(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h PackingOrderPackHandler) Print(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.PackingOrderItemRequestPrint
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesPackingOrderPack.Print(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h PackingOrderPackHandler) Dispose(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.PackingOrderItemRequestDispose
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesPackingOrderPack.Dispose(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}
