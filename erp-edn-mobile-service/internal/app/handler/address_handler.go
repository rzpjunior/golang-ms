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

type AddressHandler struct {
	Option          global.HandlerOptions
	ServicesAddress service.IAddressService
}

// URLMapping implements router.RouteHandlers
func (h *AddressHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesAddress = service.NewServiceAddress()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("edn_app"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("edn_app"))
	r.GET("/gp", h.GetGp, cMiddleware.Authorized("edn_app"))
	r.GET("/gp/:id", h.DetailGp, cMiddleware.Authorized("edn_app"))
}

func (h AddressHandler) Get(c echo.Context) (err error) {
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

	var addresses []*dto.Address
	var total int64
	addresses, total, err = h.ServicesAddress.GetAddresss(ctx.Request().Context(), dto.AddressListRequest{
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

	ctx.DataList(addresses, total, page)

	return ctx.Serve(err)
}

func (h AddressHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var customer *dto.AddressResponse

	// var id string
	id := ctx.Param("id")

	customer, err = h.ServicesAddress.GetAddressDetailById(ctx.Request().Context(), dto.AddressDetailRequest{
		Id: id,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = customer

	return ctx.Serve(err)
}

func (h AddressHandler) GetGp(c echo.Context) (err error) {
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

	var addresses []*dto.AddressGP
	var total int64
	addresses, total, err = h.ServicesAddress.GetListGp(ctx.Request().Context(), dto.GetAddressGPListRequest{
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

	ctx.DataList(addresses, total, page)

	return ctx.Serve(err)
}

func (h AddressHandler) DetailGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var address *dto.AddressGP

	var id string
	id = ctx.GetParamString("id")

	address, err = h.ServicesAddress.GetDetailGp(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = address

	return ctx.Serve(err)
}
