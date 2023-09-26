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

type WrtHandler struct {
	Option      global.HandlerOptions
	ServicesWrt service.IWrtService
}

// URLMapping implements router.RouteHandlers
func (h *WrtHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesWrt = service.NewWrtService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("wrt_rdl"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("wrt_rdd"))
}

func (h WrtHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	// get params filters
	regionID := ctx.GetParamInt("region_id")
	search := ctx.GetParamString("search")

	var Wrtes []dto.WrtResponse
	var total int64
	Wrtes, total, err = h.ServicesWrt.Get(ctx.Request().Context(), page.Start, page.Limit, int64(regionID), search)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(Wrtes, total, page)

	return ctx.Serve(err)
}

func (h WrtHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var Wrt dto.WrtResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	Wrt, err = h.ServicesWrt.GetDetail(ctx.Request().Context(), id, "")
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = Wrt

	return ctx.Serve(err)
}
