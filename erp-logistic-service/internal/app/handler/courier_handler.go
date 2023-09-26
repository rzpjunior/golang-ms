package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type CourierHandler struct {
	Option          global.HandlerOptions
	ServicesCourier service.ICourierService
}

// URLMapping implements router.RouteHandlers
func (h *CourierHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCourier = service.NewServiceCourier()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/detail", h.GetDetail, cMiddleware.Authorized())
}

func (h CourierHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get param filters
	name := ctx.GetParamString("name")
	courierVendor := ctx.GetParamString("courier_vendor_id")

	req := dto.GetCourierRequest{
		Limit:           page.Limit,
		Offset:          page.Offset,
		Name:            name,
		CourierVendorID: courierVendor,
	}

	courier, total, err := h.ServicesCourier.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(courier, total, page)

	return ctx.Serve(err)
}

func (h CourierHandler) GetDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	id := ctx.GetParamString("id")

	ctx.ResponseData, err = h.ServicesCourier.GetDetail(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
