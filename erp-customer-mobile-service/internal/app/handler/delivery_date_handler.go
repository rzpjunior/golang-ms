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

type DeliveryDateHandler struct {
	Option              global.HandlerOptions
	ServiceDeliveryDate service.IDeliveryDateService
	ServiceWRT          service.IWRTService
}

// URLMapping implements router.RouteHandlers
func (h *DeliveryDateHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceDeliveryDate = service.NewDeliveryDateService()
	h.ServiceWRT = service.NewWRTService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("/get-date", h.GetDate, cMiddleware.Authorized("public"))
	r.POST("/get-wrt", h.GetWRT, cMiddleware.Authorized("public"))
}

func (h DeliveryDateHandler) GetDate(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.DeliveryDateRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var deliveryDate dto.DeliveryDateResponse
	deliveryDate, err = h.ServiceDeliveryDate.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = deliveryDate

	return ctx.Serve(err)
}

func (h DeliveryDateHandler) GetWRT(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.WrtRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var wrt []dto.WrtResponse
	wrt, err = h.ServiceWRT.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = wrt

	return ctx.Serve(err)
}
