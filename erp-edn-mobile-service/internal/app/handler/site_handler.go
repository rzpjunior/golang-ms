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

type SiteHandler struct {
	Option       global.HandlerOptions
	ServicesSite service.ISiteService
}

// URLMapping implements router.RouteHandlers
func (h *SiteHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesSite = service.NewServiceSite()

	cMiddleware := middleware.NewMiddleware()
	// r.GET("", h.Get, cMiddleware.Authorized("edn_app"))
	// r.GET("/:id", h.Detail, cMiddleware.Authorized("edn_app"))
	r.GET("", h.GetGp, cMiddleware.Authorized("edn_app"))
	r.GET("/:id", h.DetailGp, cMiddleware.Authorized("edn_app"))
}

func (h SiteHandler) Get(c echo.Context) (err error) {
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
	offset := ctx.GetParamInt("page") - 1
	limit := ctx.GetParamInt("per_page")

	var sites []*dto.SiteResponse
	var total int64
	sites, err = h.ServicesSite.GetSites(ctx.Request().Context(), dto.SiteListRequest{
		Limit:   int32(limit),
		Offset:  int32(offset),
		Status:  int32(status),
		Search:  search,
		OrderBy: orderBy,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(sites, total, page)

	return ctx.Serve(err)
}

func (h SiteHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var site *dto.SiteResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	site, err = h.ServicesSite.GetSiteDetailById(ctx.Request().Context(), dto.SiteDetailRequest{
		Id: int32(id),
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = site

	return ctx.Serve(err)
}

func (h SiteHandler) GetGp(c echo.Context) (err error) {
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

	var site []*dto.SiteGP
	var total int64
	site, total, err = h.ServicesSite.GetGP(ctx.Request().Context(), dto.SiteListRequest{
		Limit:  int32(page.Limit),
		Offset: int32(page.Offset),
		Search: search,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(site, total, page)

	return ctx.Serve(err)
}

func (h SiteHandler) DetailGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var site *dto.SiteGP

	var id string
	id = ctx.GetUriParamString("id")

	site, err = h.ServicesSite.GetDetaiGPlById(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = site

	return ctx.Serve(err)
}
