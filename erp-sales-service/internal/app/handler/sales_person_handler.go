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

type SalesPersonHandler struct {
	Option              global.HandlerOptions
	ServicesSalesPerson service.ISalesPersonService
}

// URLMapping implements router.RouteHandlers
func (h *SalesPersonHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesSalesPerson = service.NewServiceSalesPerson()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/detail", h.GetDetail, cMiddleware.Authorized())
}

func (h SalesPersonHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	salesTerritoryID := ctx.GetParamString("sales_territory_id")
	status := ctx.GetParamInt("status")
	search := ctx.GetParamString("search")

	req := dto.GetSalesPersonRequest{
		Limit:            page.Limit,
		Offset:           page.Offset,
		SalesTerritoryID: salesTerritoryID,
		Status:           int8(status),
		Search:           search,
	}

	salesPerson, total, err := h.ServicesSalesPerson.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(salesPerson, total, page)

	return ctx.Serve(err)
}

func (h SalesPersonHandler) GetDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	id := ctx.GetParamString("id")

	ctx.ResponseData, err = h.ServicesSalesPerson.GetDetail(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
