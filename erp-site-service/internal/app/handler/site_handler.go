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

type SiteHandler struct {
	Option       global.HandlerOptions
	ServicesSite service.ISiteService
}

// URLMapping implements router.RouteHandlers
func (h *SiteHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesSite = service.NewServiceSite()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.GetDetail, cMiddleware.Authorized())
}

func (h SiteHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	req := dto.GetSiteRequest{
		Limit:  page.Limit,
		Offset: page.Offset,
		Search: ctx.GetParamString("search"),
	}

	site, total, err := h.ServicesSite.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(site, total, page)

	return ctx.Serve(err)
}

func (h SiteHandler) GetDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	id := c.Param("id")

	ctx.ResponseData, err = h.ServicesSite.GetDetail(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
