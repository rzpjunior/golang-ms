package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/service"
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

	r.GET("/app", h.readApp, cMiddleware.Authorized("public"))
	r.GET("/glossary", h.readGlossary, cMiddleware.Authorized("public"))

	r.POST("/delivery-fee", h.getDeliveryFee, cMiddleware.Authorized("public"))
}

func (h ConfigHandler) readApp(c echo.Context) (err error) {
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

	var menus []dto.ApplicationConfigResponse
	var total int64
	menus, total, err = h.ServiceConfig.GetAppConfig(ctx.Request().Context(), int32(application), field, attribute, value)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(menus, total, page)
	return ctx.Serve(err)
}

func (h ConfigHandler) readGlossary(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	valueInt := ctx.GetParamInt("valueInt")
	table := ctx.GetParamString("table")
	attribute := ctx.GetParamString("attribute")
	valueName := ctx.GetParamString("valueName")

	var menus []dto.GlossaryResponse
	var total int64
	menus, total, err = h.ServiceConfig.GetGlossary(ctx.Request().Context(), table, attribute, valueInt, valueName)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(menus, total, page)
	return ctx.Serve(err)
}

func (h ConfigHandler) getDeliveryFee(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetDeliveryFee

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	var deliveryFee *dto.ResponseGetDeliveryFee
	deliveryFee, err = h.ServiceConfig.GetDeliveryFee(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.Data(deliveryFee)
	return ctx.Serve(err)
}
