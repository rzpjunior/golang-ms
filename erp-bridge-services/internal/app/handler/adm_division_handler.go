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

type AdmDivisionHandler struct {
	Option              global.HandlerOptions
	ServicesAdmDivision service.IAdmDivisionService
}

// URLMapping implements router.RouteHandlers
func (h *AdmDivisionHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesAdmDivision = service.NewAdmDivisionService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
}

func (h AdmDivisionHandler) Get(c echo.Context) (err error) {
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

	regionID := ctx.GetParamInt("region_id")
	subDistrictID := ctx.GetParamInt("sub_district_id")

	var admDivisions []dto.AdmDivisionResponse
	var total int64
	admDivisions, total, err = h.ServicesAdmDivision.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, int64(regionID), int64(subDistrictID))
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(admDivisions, total, page)

	return ctx.Serve(err)
}

func (h AdmDivisionHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var admDivision dto.AdmDivisionResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	admDivision, err = h.ServicesAdmDivision.GetDetail(ctx.Request().Context(), id, "", 0, 0)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = admDivision

	return ctx.Serve(err)
}
