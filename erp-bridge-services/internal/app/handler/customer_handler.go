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
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
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
	// get params filters
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")
	customerTypeId := ctx.GetParamInt("customer_type_id")

	var customers []dto.CustomerResponse
	var total int64
	customers, total, err = h.ServicesCustomer.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, int64(customerTypeId))
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(customers, total, page)

	return ctx.Serve(err)
}

func (h CustomerHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var customer dto.CustomerResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	customer, err = h.ServicesCustomer.GetDetail(ctx.Request().Context(), id, "", "")
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = customer

	return ctx.Serve(err)
}
