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

type CourierVendorHandler struct {
	Option                global.HandlerOptions
	ServicesCourierVendor service.ICourierVendorService
}

// URLMapping implements router.RouteHandlers
func (h *CourierVendorHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCourierVendor = service.NewServiceCourierVendor()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/detail", h.GetDetail, cMiddleware.Authorized())
}

func (h CourierVendorHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get param filter
	siteId := ctx.GetParamString("site_id")
	courierVendorName := ctx.GetParamString("name")

	req := dto.GetCourierVendorRequest{
		Limit:             page.Limit,
		Offset:            page.Offset,
		SiteId:            siteId,
		CourierVendorName: courierVendorName,
	}

	courierVendor, total, err := h.ServicesCourierVendor.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(courierVendor, total, page)

	return ctx.Serve(err)
}

func (h CourierVendorHandler) GetDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	id := ctx.GetParamString("id")

	ctx.ResponseData, err = h.ServicesCourierVendor.GetDetail(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
