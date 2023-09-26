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

type CourierHandler struct {
	Option          global.HandlerOptions
	ServicesCourier service.ICourierService
}

// URLMapping implements router.RouteHandlers
func (h *CourierHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCourier = service.NewCourierService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
}

func (h CourierHandler) Get(c echo.Context) (err error) {
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
	vehicleProfileID := ctx.GetParamInt("vehicle_profile_id")
	emergencyMode := ctx.GetParamInt("emergency_mode")

	var couriers []dto.CourierResponse
	var total int64
	couriers, total, err = h.ServicesCourier.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, int64(vehicleProfileID), int64(emergencyMode))
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(couriers, total, page)

	return ctx.Serve(err)
}

func (h CourierHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var courier dto.CourierResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	courier, err = h.ServicesCourier.GetDetail(ctx.Request().Context(), id, "", 0)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = courier

	return ctx.Serve(err)
}
