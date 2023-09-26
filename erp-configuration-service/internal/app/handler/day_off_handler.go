package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type DayOffHandler struct {
	Option        global.HandlerOptions
	ServiceDayOff service.IDayOffService
}

// URLMapping implements router.RouteHandlers
func (h *DayOffHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceDayOff = service.NewDayOffService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.POST("", h.Create, cMiddleware.Authorized("dyo_crt"))
	r.PUT("/archive/:id", h.Archive, cMiddleware.Authorized("dyo_arc"))
	r.PUT("/unarchive/:id", h.UnArchive, cMiddleware.Authorized("dyo_urc"))
}

func (h DayOffHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	status := ctx.GetParamInt("status")
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	startDate := ctx.GetParamDateTime("start_date")
	endDate := ctx.GetParamDateTime("end_date")

	var dayOffs []dto.DayOffResponse
	var total int64
	dayOffs, total, err = h.ServiceDayOff.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, startDate, endDate)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(dayOffs, total, page)

	return ctx.Serve(err)
}

func (h DayOffHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.DayOffRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServiceDayOff.Create(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h DayOffHandler) Archive(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServiceDayOff.Archive(ctx.Request().Context(), id)

	return ctx.Serve(err)
}

func (h DayOffHandler) UnArchive(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServiceDayOff.UnArchive(ctx.Request().Context(), id)

	return ctx.Serve(err)
}
