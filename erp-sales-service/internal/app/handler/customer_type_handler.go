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

type CustomerTypeHandler struct {
	Option               global.HandlerOptions
	ServicesCustomerType service.ICustomerTypeService
}

// URLMapping implements router.RouteHandlers
func (h *CustomerTypeHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCustomerType = service.NewServiceCustomerType()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/detail", h.GetDetail, cMiddleware.Authorized())
}

func (h CustomerTypeHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	req := dto.GetCustomerTypeRequest{
		Limit:  page.Limit,
		Offset: page.Offset,
	}

	customerType, total, err := h.ServicesCustomerType.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(customerType, total, page)

	return ctx.Serve(err)
}

func (h CustomerTypeHandler) GetDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	id := ctx.GetParamString("id")

	ctx.ResponseData, err = h.ServicesCustomerType.GetDetail(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
