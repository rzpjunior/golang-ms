package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type SalesPriceLevelHandler struct {
	Option             global.HandlerOptions
	ServicesPriceLevel service.ISalesPriceLevelService
}

// URLMapping implements router.RouteHandlers
func (h *SalesPriceLevelHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesPriceLevel = service.NewServiceSalesPriceLevel()

	cMiddleware := middleware.NewMiddleware()
	// GP Integrated
	r.GET("", h.GetGp, cMiddleware.Authorized("edn_app"))
	r.GET("/:id", h.DetailGp, cMiddleware.Authorized("edn_app"))
}

func (h SalesPriceLevelHandler) GetGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	gnlRegion := ctx.GetParamString("region_id")
	gnlCustTypeId := ctx.GetParamString("customer_type_id")
	prclevel := ctx.GetParamString("price_level")
	offset := ctx.GetParamInt("page") - 1
	limit := ctx.GetParamInt("per_page")

	var pl []*dto.SalesPriceLevel
	var total int64
	pl, total, err = h.ServicesPriceLevel.GetGP(ctx.Request().Context(), dto.GetSalesPriceLevelListRequest{
		Limit:      int32(limit),
		Offset:     int32(offset),
		RegionID:   gnlRegion,
		CustTypeID: gnlCustTypeId,
		PriceLevel: prclevel,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(pl, total, page)

	return ctx.Serve(err)
}

func (h SalesPriceLevelHandler) DetailGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var pl *dto.SalesPriceLevel

	var id string
	id = ctx.Param("id")

	pl, err = h.ServicesPriceLevel.GetDetaiGPlById(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = pl

	return ctx.Serve(err)
}
