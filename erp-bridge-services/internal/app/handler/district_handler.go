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

type DistrictHandler struct {
	Option           global.HandlerOptions
	ServicesDistrict service.IDistrictService
}

// URLMapping implements router.RouteHandlers
func (h *DistrictHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesDistrict = service.NewDistrictService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
}

func (h DistrictHandler) Get(c echo.Context) (err error) {
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

	var districts []*dto.DistrictResponse
	var total int64
	districts, total, err = h.ServicesDistrict.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(districts, total, page)

	return ctx.Serve(err)
}

func (h DistrictHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var district *dto.DistrictResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	district, err = h.ServicesDistrict.GetDetail(ctx.Request().Context(), id, "")
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = district

	return ctx.Serve(err)
}
