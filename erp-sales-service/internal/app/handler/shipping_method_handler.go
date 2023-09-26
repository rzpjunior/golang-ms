package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type ShippingMethodHandler struct {
	Option                 global.HandlerOptions
	ServicesShippingMethod service.IShippingMethodService
}

// URLMapping implements router.RouteHandlers
func (h *ShippingMethodHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesShippingMethod = service.NewServiceShippingMethod()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.Get, cMiddleware.Authorized())
}

func (h ShippingMethodHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	typeMethod := ctx.GetParamString("type")
	req := &dto.GetShippingMethodRequest{
		Limit:  int64(page.Limit),
		Offset: int64(page.Offset),
		Type:   typeMethod,
	}

	salesPerson, total, err := h.ServicesShippingMethod.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(salesPerson, total, page)

	return ctx.Serve(err)
}
