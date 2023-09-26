package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type SalesPriceLevelHandler struct {
	Option                  global.HandlerOptions
	ServicesSalesPriceLevel service.ISalesPriceLevelService
}

// URLMapping implements router.RouteHandlers
func (h *SalesPriceLevelHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesSalesPriceLevel = service.NewServiceSalesPriceLevel()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/detail", h.GetDetail, cMiddleware.Authorized())
}

func (h SalesPriceLevelHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	regionID := ctx.GetParamString("region_id")
	customerTypeID := ctx.GetParamString("customer_type_id")

	req := &dto.GetSalesPriceLevelRequest{
		Limit:          int64(page.Limit),
		Offset:         int64(page.Offset),
		RegionID:       regionID,
		CustomerTypeID: customerTypeID,
	}

	salesPriceLevel, total, err := h.ServicesSalesPriceLevel.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(salesPriceLevel, total, page)

	return ctx.Serve(err)
}

func (h SalesPriceLevelHandler) GetDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	id := ctx.GetParamString("id")

	ctx.ResponseData, err = h.ServicesSalesPriceLevel.GetDetail(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
