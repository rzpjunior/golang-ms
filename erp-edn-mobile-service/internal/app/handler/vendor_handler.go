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

type VendorHandler struct {
	Option         global.HandlerOptions
	ServicesVendor service.IVendorService
}

// URLMapping implements router.RouteHandlers
func (h *VendorHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesVendor = service.NewServiceVendor()

	cMiddleware := middleware.NewMiddleware()
	// r.GET("", h.Get, cMiddleware.Authorized("edn_app"))
	// r.GET("/:id", h.Detail, cMiddleware.Authorized("edn_app"))

	// GP Integrated
	r.GET("", h.GetGp, cMiddleware.Authorized("edn_app"))
	r.GET("/:id", h.DetailGp, cMiddleware.Authorized("edn_app"))
}

func (h VendorHandler) Get(c echo.Context) (err error) {
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

	var vendors []*dto.VendorResponse
	var total int64
	vendors, err = h.ServicesVendor.GetVendors(ctx.Request().Context(), dto.VendorListRequest{
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

	ctx.DataList(vendors, total, page)

	return ctx.Serve(err)
}

func (h VendorHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var vendor *dto.VendorResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	vendor, err = h.ServicesVendor.GetVendorDetailById(ctx.Request().Context(), dto.VendorDetailRequest{
		Id: int32(id),
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = vendor

	return ctx.Serve(err)
}

func (h VendorHandler) GetGp(c echo.Context) (err error) {
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

	var vendors []*dto.VendorGP
	var total int64
	vendors, total, err = h.ServicesVendor.GetListGp(ctx.Request().Context(), dto.GetVendorGPListRequest{
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

	ctx.DataList(vendors, total, page)

	return ctx.Serve(err)
}

func (h VendorHandler) DetailGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var vendor *dto.VendorGP

	var id string
	id = ctx.GetUriParamString("id")

	vendor, err = h.ServicesVendor.GetDetailGp(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = vendor

	return ctx.Serve(err)
}
