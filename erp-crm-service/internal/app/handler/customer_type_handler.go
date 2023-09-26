package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type CustomerTypeHandler struct {
	Option               global.HandlerOptions
	ServicesCustomerType service.ICustomerTypeService
}

// URLMapping implements router.RouteHandlers
func (h *CustomerTypeHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCustomerType = service.NewCustomerTypeService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/detail", h.GetDetail, cMiddleware.Authorized())
}

func (h CustomerTypeHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	search := ctx.GetParamString("search")
	status := ctx.GetParamInt("status")

	param := &dto.CustomerTypeGetListRequest{
		Offset: page.Start,
		Limit:  page.Limit,
		Search: search,
		Status: int8(status),
	}

	var customerTypes []*dto.CustomerTypeResponse
	var total int64
	customerTypes, total, err = h.ServicesCustomerType.Get(ctx.Request().Context(), param)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(customerTypes, total, page)

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
