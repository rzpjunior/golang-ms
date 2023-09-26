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

type GmapsHandler struct {
	Option       global.HandlerOptions
	ServiceGmaps service.IGmapsService
}

// URLMapping declare endpoint with handler function.
func (h *GmapsHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceGmaps = service.NewGmapsService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("/autocomplete", h.GetAutoComplete, cMiddleware.Authorized("public"))
	r.POST("/geocode", h.GetGeocode, cMiddleware.Authorized("public"))
}

func (h GmapsHandler) GetAutoComplete(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.GetAutoCompleteRequest

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceGmaps.GetAutoComplete(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h GmapsHandler) GetGeocode(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.GetGeocodeRequest

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceGmaps.GetGeocode(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)

}
