package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type WrtHandler struct {
	Option      global.HandlerOptions
	ServicesWrt service.IWrtService
}

// URLMapping implements router.RouteHandlers
func (h *WrtHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesWrt = service.NewServiceWrt()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("edn_app"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("edn_app"))
	r.GET("/gp", h.GetGp, cMiddleware.Authorized("edn_app"))
	r.GET("/gp/:id", h.DetailGp, cMiddleware.Authorized("edn_app"))
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

	// get params filter
	search := ctx.GetParamString("search")
	regionId := ctx.GetParamString("region_id")
	offset := ctx.GetParamInt("page") - 1
	limit := ctx.GetParamInt("per_page")

	var wrts []*dto.WrtResponse
	var total int64
	wrts, err = h.ServicesWrt.Get(ctx.Request().Context(), dto.GetWrtListRequest{
		Limit:    int32(limit),
		Offset:   int32(offset),
		RegionId: regionId,
		Search:   search,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(wrts, total, page)

	return ctx.Serve(err)
}

func (h WrtHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var wrt *dto.WrtResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	wrt, err = h.ServicesWrt.GetDetailById(ctx.Request().Context(), dto.GetWrtDetailRequest{
		Id: int32(id),
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = wrt

	return ctx.Serve(err)
}

func (h WrtHandler) GetGp(c echo.Context) (err error) {
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

	var wrt []*dto.WrtGP
	var total int64
	wrt, total, err = h.ServicesWrt.GetGP(ctx.Request().Context(), dto.GetWrtListRequest{
		Limit:  int32(page.Limit),
		Offset: int32(page.Offset),
		Search: search,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(wrt, total, page)

	return ctx.Serve(err)
}

func (h WrtHandler) DetailGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var wrt *dto.WrtGP

	var id string
	id = ctx.GetParamString("id")

	wrt, err = h.ServicesWrt.GetDetaiGPlById(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = wrt

	return ctx.Serve(err)
}
