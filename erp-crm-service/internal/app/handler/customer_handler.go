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

type CustomerHandler struct {
	Option           global.HandlerOptions
	ServicesCustomer service.ICustomerService
}

// URLMapping implements router.RouteHandlers
func (h *CustomerHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCustomer = service.NewCustomerService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.GetDetail, cMiddleware.Authorized())
}

func (h CustomerHandler) Get(c echo.Context) (err error) {
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
	customerType := ctx.GetParamString("customer_type")

	param := &dto.CustomerGetListRequest{
		Offset:       page.Start,
		Limit:        page.Limit,
		Search:       search,
		Status:       int8(status),
		CustomerType: customerType,
	}

	var customers []*dto.CustomerResponseGet
	var total int64
	customers, total, err = h.ServicesCustomer.Get(ctx.Request().Context(), param)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(customers, total, page)

	return ctx.Serve(err)
}

func (h CustomerHandler) GetDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	param := &dto.CustomerRequestGetDetail{
		ID: id,
	}

	ctx.ResponseData, err = h.ServicesCustomer.GetDetailComplex(ctx.Request().Context(), param)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
