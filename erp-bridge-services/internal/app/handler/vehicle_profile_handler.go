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

type VehicleProfileHandler struct {
	Option                 global.HandlerOptions
	ServicesVehicleProfile service.IVehicleProfileService
}

// URLMapping implements router.RouteHandlers
func (h *VehicleProfileHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesVehicleProfile = service.NewVehicleProfileService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
}

func (h VehicleProfileHandler) Get(c echo.Context) (err error) {
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
	courierVendorID := ctx.GetParamInt("courier_vendor_id")

	var vehicleProfiles []dto.VehicleProfileResponse
	var total int64
	vehicleProfiles, total, err = h.ServicesVehicleProfile.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, int64(courierVendorID))
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(vehicleProfiles, total, page)

	return ctx.Serve(err)
}

func (h VehicleProfileHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var vehicleProfile dto.VehicleProfileResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	vehicleProfile, err = h.ServicesVehicleProfile.GetDetail(ctx.Request().Context(), id, "")
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = vehicleProfile

	return ctx.Serve(err)
}
