package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/service"
	"github.com/labstack/echo/v4"
)

type CourierVendorHandler struct {
	Option                global.HandlerOptions
	ServicesCourierVendor service.ICourierVendorService
}

// URLMapping implements router.RouteHandlers
func (h *CourierVendorHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCourierVendor = service.NewCourierVendorService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
}

func (h CourierVendorHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filters
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")
	siteID := ctx.GetParamInt("site_id")

	var courierVendors []dto.CourierVendorResponse
	var total int64
	courierVendors, total, err = h.ServicesCourierVendor.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, int64(siteID))
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(courierVendors, total, page)

	return ctx.Serve(err)
}

func (h CourierVendorHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var courierVendor dto.CourierVendorResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	courierVendor, err = h.ServicesCourierVendor.GetDetail(ctx.Request().Context(), id, "")
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = courierVendor

	return ctx.Serve(err)
}
