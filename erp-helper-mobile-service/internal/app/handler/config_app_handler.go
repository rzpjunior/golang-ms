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

type ConfigAppHandler struct {
	Option            global.HandlerOptions
	ServicesConfigApp service.IConfigAppService
}

// URLMapping implements router.RouteHandlers
func (h *ConfigAppHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesConfigApp = service.NewServiceConfigApp()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.GetConfigApp, cMiddleware.AuthorizedHelperMobile())
}

func (h ConfigAppHandler) GetConfigApp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.GetConfigAppRequest

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	req.Id = ctx.GetParamInt("id")
	req.Application = ctx.GetParamInt("application")
	req.Field = ctx.GetParamString("field")
	req.Attribute = ctx.GetParamString("attribute")
	req.Value = ctx.GetParamString("value")

	glossaries, total, err := h.ServicesConfigApp.GetConfigApp(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(glossaries, total, page)

	return ctx.Serve(err)
}
