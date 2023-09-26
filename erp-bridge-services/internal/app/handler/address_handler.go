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

type AddressHandler struct {
	Option          global.HandlerOptions
	ServicesAddress service.IAddressService
}

// URLMapping implements router.RouteHandlers
func (h *AddressHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesAddress = service.NewAddressService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
}

// Get List
func (h AddressHandler) Get(c echo.Context) (err error) {
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

	archetypeID := ctx.GetParamInt("archetype_id")
	admDivisionID := ctx.GetParamInt("adm_division_id")
	siteID := ctx.GetParamInt("site_id")
	salespersonID := ctx.GetParamInt("salesperson_id")
	territoryID := ctx.GetParamInt("territory_id")
	taxScheduleID := ctx.GetParamInt("tax_schedule_id")

	var addresses []dto.AddressResponse
	var total int64
	addresses, total, err = h.ServicesAddress.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, int64(archetypeID), int64(admDivisionID), int64(siteID), int64(salespersonID), int64(territoryID), int64(taxScheduleID))
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(addresses, total, page)

	return ctx.Serve(err)
}

// Detail
func (h AddressHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var address dto.AddressResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	address, err = h.ServicesAddress.GetDetail(ctx.Request().Context(), id, "")
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = address

	return ctx.Serve(err)
}
