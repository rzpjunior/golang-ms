package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/service"
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

	r.GET("", h.Get, cMiddleware.AuthorizedHelperMobile())
	r.GET("/detail", h.GetDetail, cMiddleware.AuthorizedHelperMobile())
}

func (h WrtHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	siteId := ctx.Request().Context().Value(constants.KeySiteID).(string)

	req := dto.GetWrtRequest{
		Limit:    page.Limit,
		Offset:   page.Offset,
		Search:   ctx.GetParamString("search"),
		RegionId: ctx.GetParamInt("region_id"),
		SiteId:   siteId,
	}

	wrt, total, err := h.ServicesWrt.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(wrt, total, page)

	return ctx.Serve(err)
}

func (h WrtHandler) GetDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	id := ctx.GetParamString("id")

	ctx.ResponseData, err = h.ServicesWrt.GetDetail(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
