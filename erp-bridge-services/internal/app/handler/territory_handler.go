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

type TerritoryHandler struct {
	Option            global.HandlerOptions
	ServicesTerritory service.ITerritoryService
}

// URLMapping implements router.RouteHandlers
func (h *TerritoryHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesTerritory = service.NewTerritoryService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
}

func (h TerritoryHandler) Get(c echo.Context) (err error) {
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
	salespersonID := ctx.GetParamInt("salesperson_id")
	CustomerTypeID := ctx.GetParamInt("customer_type_id")
	subDistrictID := ctx.GetParamInt("sub_district_id")

	var territories []dto.TerritoryResponse
	var total int64
	territories, total, err = h.ServicesTerritory.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, int64(regionID), int64(salespersonID), int64(CustomerTypeID), int64(subDistrictID))
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(territories, total, page)

	return ctx.Serve(err)
}

func (h TerritoryHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var territory dto.TerritoryResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	code := ctx.GetParamString("code")
	salespersonID := ctx.GetParamInt("salesperson_id")

	territory, err = h.ServicesTerritory.GetDetail(ctx.Request().Context(), id, code, int64(salespersonID))
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = territory

	return ctx.Serve(err)
}
