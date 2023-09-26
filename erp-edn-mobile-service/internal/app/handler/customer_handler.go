package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"github.com/labstack/echo/v4"
)

type CustomerHandler struct {
	Option           global.HandlerOptions
	ServicesCustomer service.ICustomerService
}

// URLMapping implements router.RouteHandlers
func (h *CustomerHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCustomer = service.NewServiceCustomer()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("edn_app"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("edn_app"))
	r.GET("/gp", h.GetGp, cMiddleware.Authorized("edn_app"))
	r.GET("/gp/:id", h.DetailGp, cMiddleware.Authorized("edn_app"))
	r.POST("", h.CreateGP, cMiddleware.Authorized("edn_app"))
	r.GET("/overduemitra", h.GetOverdueMitra, cMiddleware.Authorized("edn_app"))
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

	// get params filter
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")
	offset := ctx.GetParamInt("page") - 1
	limit := ctx.GetParamInt("per_page")

	var customerAcquisitions []*dto.CustomerResponse
	var total int64
	customerAcquisitions, total, err = h.ServicesCustomer.GetCustomers(ctx.Request().Context(), dto.CustomerListRequest{
		Limit:   int32(limit),
		Offset:  int32(offset),
		Status:  int32(status),
		Search:  search,
		OrderBy: orderBy,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(customerAcquisitions, total, page)

	return ctx.Serve(err)
}

func (h CustomerHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var customerAcquisition *dto.CustomerResponse

	var id string
	id = ctx.Param("id")

	customerAcquisition, err = h.ServicesCustomer.GetCustomerDetailById(ctx.Request().Context(), dto.CustomerDetailRequest{
		Id: id,
	})

	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = customerAcquisition

	return ctx.Serve(err)
}

func (h CustomerHandler) GetGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")

	var customers []*dto.CustomerGP
	var total int64
	customers, total, err = h.ServicesCustomer.GetListGp(ctx.Request().Context(), dto.GetCustomerGPListRequest{
		Limit:   int32(page.Limit),
		Offset:  int32(page.Offset),
		Status:  int32(status),
		Search:  search,
		OrderBy: orderBy,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(customers, total, page)

	return ctx.Serve(err)
}

func (h CustomerHandler) DetailGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var customer *dto.CustomerGP

	var id string
	id = ctx.GetParamString("id")

	customer, err = h.ServicesCustomer.GetDetailGp(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = customer

	return ctx.Serve(err)
}

func (h CustomerHandler) CreateGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *bridgeService.CreateCustomerGPResponse

	var req dto.CreateCustomerGPRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesCustomer.CreateGP(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h CustomerHandler) GetOverdueMitra(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	offset := ctx.GetParamInt("page") - 1
	limit := ctx.GetParamInt("per_page")

	var customerAcquisitions []*dto.CustomerResponse
	var total int64
	customerAcquisitions, total, err = h.ServicesCustomer.GetOverdueMitra(ctx.Request().Context(), dto.CustomerListRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(customerAcquisitions, total, page)

	return ctx.Serve(err)
}
