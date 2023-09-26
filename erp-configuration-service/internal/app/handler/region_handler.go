package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type RegionHandler struct {
	Option         global.HandlerOptions
	ServicesRegion service.IRegionService
}

// URLMapping implements router.RouteHandlers
func (h *RegionHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesRegion = service.NewRegionService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
}

func (h RegionHandler) Get(c echo.Context) (err error) {
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
	status := ctx.GetParamInt("status")

	param := &dto.RegionGetRequest{
		Offset: page.Start,
		Limit:  page.Limit,
		Search: search,
		Status: int8(status),
	}

	var Regiones []*dto.RegionResponse
	var total int64
	Regiones, total, err = h.ServicesRegion.Get(ctx.Request().Context(), param)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(Regiones, total, page)

	return ctx.Serve(err)
}
