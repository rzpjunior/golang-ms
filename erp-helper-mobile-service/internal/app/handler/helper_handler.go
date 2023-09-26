package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type HelperHandler struct {
	Option         global.HandlerOptions
	ServicesHelper service.IHelperService
}

// URLMapping implements router.RouteHandlers
func (h *HelperHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesHelper = service.NewServiceHelper()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.Get, cMiddleware.AuthorizedHelperMobile())
}

func (h HelperHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	req := &dto.GetHelperRequest{
		Limit:  page.PerPage,
		Offset: page.Page - 1,
		SiteId: ctx.GetParamString("site_id"),
		Role:   ctx.GetParamString("role"),
		Name:   ctx.GetParamString("name"),
		Type:   ctx.GetParamString("type"),
	}

	ctx.ResponseData, err = h.ServicesHelper.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
