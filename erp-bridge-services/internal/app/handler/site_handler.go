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

type SiteHandler struct {
	Option       global.HandlerOptions
	ServicesSite service.ISiteService
}

// URLMapping implements router.RouteHandlers
func (h *SiteHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesSite = service.NewSiteService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
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

	// get params filters
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")

	var sites []dto.SiteResponse
	var total int64
	sites, total, err = h.ServicesSite.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(sites, total, page)

	return ctx.Serve(err)
}

func (h SiteHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var site dto.SiteResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	site, err = h.ServicesSite.GetDetail(ctx.Request().Context(), id, "")
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = site

	return ctx.Serve(err)
}
