package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type ConfigHandler struct {
	Option        global.HandlerOptions
	ServiceConfig service.IConfigService
}

// URLMapping declare endpoint with handler function.
func (h *ConfigHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceConfig = service.NewConfigService()
	cMiddleware := middleware.NewMiddleware()

	r.GET("/app", h.GetApp, cMiddleware.Authorized("purchaser_app"))
	r.GET("/glossary", h.GetGlossary, cMiddleware.Authorized("purchaser_app"))
}

func (h ConfigHandler) GetApp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	application := ctx.GetParamInt("application")
	field := ctx.GetParamString("field")
	attribute := ctx.GetParamString("attribute")
	value := ctx.GetParamString("value")

	var configApps []dto.ConfigAppResponse
	var total int64
	configApps, total, err = h.ServiceConfig.GetConfigApp(ctx.Request().Context(), int32(application), field, attribute, value)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(configApps, total, page)

	return ctx.Serve(err)
}

func (h ConfigHandler) GetGlossary(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	valueInt := ctx.GetParamInt("value_int")
	table := ctx.GetParamString("table")
	attribute := ctx.GetParamString("attribute")
	valueName := ctx.GetParamString("value_name")

	var glossaries []dto.GlossaryResponse
	var total int64
	glossaries, total, err = h.ServiceConfig.GetGlossary(ctx.Request().Context(), table, attribute, valueInt, valueName)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(glossaries, total, page)

	return ctx.Serve(err)
}
